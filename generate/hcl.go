package generate

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/graph"
	"github.com/cycloidio/inframap/provider"
	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/configs/hcl2shim"
	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/regsrc"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/afero"
)

var (
	localSourcePrefixes = []string{
		"./",
		"../",
	}
	cachePath = path.Join(xdg.CacheHome, "inframap", "modules")
)

// FromHCL generates a new graph from the HCL on the path,
// it can be a file or a Module/Dir
func FromHCL(fs afero.Fs, p string, opt Options) (*graph.Graph, error) {
	parser := configs.NewParser(fs)

	g := graph.New()

	var (
		mod   *configs.Module
		diags hcl.Diagnostics
		err   error
	)

	if parser.IsConfigDir(p) {
		mod, diags = parser.LoadConfigDir(p)
	} else {
		f, dgs := parser.LoadConfigFile(p)
		if dgs.HasErrors() {
			return nil, errors.New(dgs.Error())
		}
		mod, diags = configs.NewModule([]*configs.File{f}, nil)
	}

	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	managedResources := make(map[string]*configs.Resource)
	for rk, rv := range mod.ManagedResources {
		managedResources[rk] = rv
	}

	installedModules := make(map[string]struct{})
	calls := make([]*configs.ModuleCall, 0)
	for _, call := range mod.ModuleCalls {
		calls = append(calls, call)
	}
	p, _ = filepath.Abs(p)
	if err := moduleInstall(calls, &managedResources, p, installedModules); err != nil {
		return nil, fmt.Errorf("unable to fetch all modules: %w", err)
	}

	// nodeCanID holds as key the `aws_alb.front` (graph.Node.Canonical)
	// and as value the UUID (graph.Node.ID) we give to it
	nodeCanID := make(map[string]string)

	// nodeIDEdges holds as key the UUID (graph.Node.ID) and as value
	// all the edges it has, in this case it's the `depends_on` values
	// that we find on the TFState
	nodeIDEdges := make(map[string][]string)

	// resourcesRawConfig holds the actual configuration of each element
	// it's represented as: graph.Node.Canonical -> Attrs
	resourcesRawConfig := make(map[string]map[string]interface{})

	if !opt.Raw {
		opt, err = checkHCLProviders(mod, opt)
		if err != nil {
			return nil, err
		}
	}

	for rk, rv := range managedResources {
		pv, rs, err := getProviderAndResource(rk, opt)
		if err != nil {
			if errors.Is(err, errcode.ErrProviderNotFound) {
				continue
			}
			return nil, err
		}

		// If it's not a Node or Edge we ignore it
		if !pv.IsNode(rs) && !pv.IsEdge(rs) {
			continue
		}

		res, err := pv.Resource(rs)
		if err != nil {
			return nil, err
		}
		n := &graph.Node{
			ID:        uuid.NewV4().String(),
			Canonical: rk,
			Resource:  *res,
		}

		err = g.AddNode(n)
		if err != nil {
			return nil, err
		}

		nodeCanID[n.Canonical] = n.ID

		links := make(map[string][]string)
		body, ok := rv.Config.(*hclsyntax.Body)
		if ok {
			links = getBodyLinks(body)
			cfg := getBodyJSON(body)
			cfg[provider.HCLCanonicalKey] = rk
			resourcesRawConfig[n.ID] = cfg
		} else {
			// If it's not a hclsyntax.Body normally
			// means it's from a JSON file, the only
			// way to work with them is with rv.Config.JustAttributes()
			// and manually deal with them. For what I've tested
			// the Blocks (so egess an ingress for example) do not
			// work and fail so we should find a workaround.
			return nil, errcode.ErrGenerateFromJSON
		}

		for _, resources := range links {
			nodeIDEdges[n.ID] = append(nodeIDEdges[n.ID], resources...)
		}
	}

	for nid, resources := range nodeIDEdges {
		for _, rkattr := range resources {
			keys := strings.Split(rkattr, ".")

			rk := rkattr
			// If the values are like 'aws_security_group.front.id'
			// and we do not need the 'id' as it's not on the key
			// so we subtract it from the key
			// But it's not always the case
			if len(keys) == 3 {
				rk = strings.Join(keys[:len(keys)-1], ".")
			}

			tnid, ok := nodeCanID[rk]
			if !ok {
				continue
			}

			err := g.AddEdge(&graph.Edge{
				ID:     uuid.NewV4().String(),
				Source: nid,
				Target: tnid,
			})
			if err != nil {
				// If the edge already exists we can ignore it
				if errors.Is(err, errcode.ErrGraphAlreadyExistsEdge) {
					continue
				}
				return nil, err
			}

		}
	}

	// call the preprocess method for each
	// TF provider in the file
	if err := preprocess(g, resourcesRawConfig, opt); err != nil {
		return nil, err
	}

	if opt.Clean {
		g.Clean()
	}

	err = fixEdges(g, resourcesRawConfig, opt)
	if err != nil {
		return nil, err
	}

	if opt.Connections {
		err = mutate(g, opt)
		if err != nil {
			return nil, err
		}
	}

	if opt.Clean {
		err = cleanHangingEdges(g, opt)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

// getBodyLinks gets all the variables used and the key in which
// they where used
func getBodyLinks(b *hclsyntax.Body) map[string][]string {
	links := make(map[string][]string)
	for attrk, attrv := range b.Attributes {
		for _, vr := range attrv.Expr.Variables() {
			links[attrk] = append(links[attrk], string(hclwrite.TokensForTraversal(vr).Bytes()))
		}
	}
	for _, block := range b.Blocks {
		for attrk, attrv := range getBodyLinks(block.Body) {
			key := fmt.Sprintf("%s.%s", block.Type, attrk)
			links[key] = append(links[key], attrv...)
		}
	}

	return links
}

// getBodyJSON gets all the variables in a JSON format
// of the actual representation
func getBodyJSON(b *hclsyntax.Body) map[string]interface{} {
	links := make(map[string]interface{})
	for attrk, attrv := range b.Attributes {
		v, _ := attrv.Expr.Value(nil)
		t := v.Type().FriendlyName()
		switch t {
		case "string", "bool", "number", "dynamic":
			for i, vr := range attrv.Expr.Variables() {
				if i > 0 {
					continue
				}
				links[attrk] = fmt.Sprintf("${%s}", string(hclwrite.TokensForTraversal(vr).Bytes()))
			}
			if _, ok := links[attrk]; !ok {
				links[attrk] = hcl2shim.ConfigValueFromHCL2(v)
			}
		case "tuple":
			aux := make([]interface{}, 0)
			for _, vr := range attrv.Expr.Variables() {
				aux = append(aux, fmt.Sprintf("${%s}", string(hclwrite.TokensForTraversal(vr).Bytes())))
			}
			if len(aux) == 0 {
				i := hcl2shim.ConfigValueFromHCL2(v)
				islice, ok := i.([]interface{})
				if !ok {
					// We continue to not add empty information to the config
					// so it's clean and only has required information
					continue
				}
				aux = islice
			}
			links[attrk] = aux
		}
	}
	for _, block := range b.Blocks {
		cfg := getBodyJSON(block.Body)
		// We continue to not add empty information to the config
		// so it's clean and only has required information
		if len(cfg) == 0 {
			continue
		}
		if _, ok := links[block.Type]; !ok {
			links[block.Type] = make([]interface{}, 0)
		}
		links[block.Type] = append(links[block.Type].([]interface{}), cfg)
	}

	return links
}

// checkHCLProviders checks if we support any of the Providers from f, if not it'll set
// the opt.Raw to true so it can be used with Raw instead of returning an empty Graph
func checkHCLProviders(mod *configs.Module, opt Options) (Options, error) {
	for rk := range mod.ManagedResources {
		_, _, err := getProviderAndResource(rk, opt)
		if err != nil {
			if errors.Is(err, errcode.ErrProviderNotFound) {
				continue
			}
			return opt, err
		}

		// If we find a resource that we support the Provider
		// then we use it
		return opt, nil
	}

	// If we reach here means the we do not support the providers
	// of the HCL, so Raw has to be used
	opt.Raw = true

	return opt, nil
}

// moduleInstall will recursively walk through the module calls required by the Terraform config, it will store the downloaded module
// in $XDG_CACHE directory and stop once all the required modules have been downloaded.
func moduleInstall(calls []*configs.ModuleCall, mRes *map[string]*configs.Resource, pwd string, installedModules map[string]struct{}) error {
	// stop condition, if there is no module to
	// fetch we stop
	if len(calls) == 0 {
		return nil
	}

	call := calls[0]
	name := call.Name

	// we check if the module is already installed
	// or not
	if _, ok := installedModules[name]; ok {
		return nil
	}

	var (
		src  string = call.SourceAddr
		vers string
	)

	// we check if the module is a Terraform registry module
	// in order to get its source address from Terraform registry
	if regMod, err := regsrc.ParseModuleSource(src); err == nil {
		client := registry.NewClient(nil, nil)
		// we get the list of available module versions
		resp, err := client.ModuleVersions(regMod)
		if err != nil {
			return fmt.Errorf("unable to get module versions: %w", err)
		}

		if len(resp.Modules) < 1 {
			return fmt.Errorf("unable to find suitable versions")
		}
		meta := resp.Modules[0]

		var (
			latest *version.Version
			// match holds the version matching the
			// source constraints set in the module call
			match *version.Version
		)
		for _, vers := range meta.Versions {
			v, err := version.NewVersion(vers.Version)
			if err != nil {
				return fmt.Errorf("unable to create version from string: %w", err)
			}

			if latest == nil || v.GreaterThan(latest) {
				latest = v
			}

			if call.Version.Required.Check(v) {
				if match == nil || v.GreaterThan(match) {
					match = v
				}
			}
		}

		vers = match.String()
		// we finally get the module location, it will return
		// a string `go-getter` compliant
		src, err = client.ModuleLocation(regMod, vers)
		if err != nil {
			return fmt.Errorf("unable to fetch module location: %w", err)
		}
	}

	// since go-getter does not support yet in-memory fs,
	// we need to initialize the parser using actual fs
	// https://github.com/hashicorp/go-getter/issues/83
	pars := configs.NewParser(nil)

	// we check if the module is a local one by checking
	// its prefix "./", "../", etc.
	var isLocal bool
	for _, prefix := range localSourcePrefixes {
		if strings.HasPrefix(src, prefix) {
			isLocal = true
		}
	}

	var (
		m     *configs.Module
		diags hcl.Diagnostics
	)

	// the module is not a local one or a Terraform registry one
	// it should be handle by `go-getter`
	if !isLocal {
		dst := path.Join(cachePath, fmt.Sprintf("%s-%s", name, vers))
		// TODO: we should add a logic to invalidate
		// the cache
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			client := &getter.Client{
				Src:  src,
				Dst:  dst,
				Pwd:  pwd,
				Mode: getter.ClientModeDir,
			}
			if err := client.Get(); err != nil {
				return fmt.Errorf("unable to get remote module: %w", err)
			}
		}
		m, diags = pars.LoadConfigDir(dst)
	} else {
		m, diags = pars.LoadConfigDir(path.Join(pwd, src))
	}
	if diags.HasErrors() {
		return fmt.Errorf("unable to load config directory: %s", diags.Error())
	}

	// fill the final map of managed resources
	// using the config freshly loaded
	for rk, rv := range m.ManagedResources {
		(*mRes)[rk] = rv
	}

	// keep a trace of the imported / loaded module
	// to avoid infinite recursion
	installedModules[name] = struct{}{}

	// create the next slice of module calls to
	// check before merging it with the current we
	// still have
	next := make([]*configs.ModuleCall, 0)
	for _, call := range m.ModuleCalls {
		next = append(next, call)
	}
	calls = append(calls[1:len(calls)], next...)

	return moduleInstall(calls, mRes, pwd, installedModules)
}
