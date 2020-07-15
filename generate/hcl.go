package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/graph"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/afero"
)

// FromHCL generates a new graph from the HCL on the path,
// it can be a file or a Module/Dir
func FromHCL(fs afero.Fs, path string, opt Options) (*graph.Graph, error) {
	parser := configs.NewParser(fs)

	g := graph.New()

	var (
		mod   *configs.Module
		diags hcl.Diagnostics
		err   error
	)

	if parser.IsConfigDir(path) {
		mod, diags = parser.LoadConfigDir(path)
	} else {
		f, dgs := parser.LoadConfigFile(path)
		if dgs.HasErrors() {
			return nil, errors.New(dgs.Error())
		}
		mod, diags = configs.NewModule([]*configs.File{f}, nil)
	}

	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
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

	for rk, rv := range mod.ManagedResources {
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

	if opt.Clean {
		g.Clean()
	}

	err = fixEdges(g, resourcesRawConfig, opt)
	if err != nil {
		return nil, err
	}

	err = mutate(g, opt)
	if err != nil {
		return nil, err
	}

	cleanHangingEdges(g, opt)

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
		case "string", "bool", "number":
			for i, vr := range attrv.Expr.Variables() {
				if i > 0 {
					continue
				}
				links[attrk] = fmt.Sprintf("${%s}", string(hclwrite.TokensForTraversal(vr).Bytes()))
			}
		case "tuple":
			aux := make([]interface{}, 0)
			for _, vr := range attrv.Expr.Variables() {
				aux = append(aux, fmt.Sprintf("${%s}", string(hclwrite.TokensForTraversal(vr).Bytes())))
			}
			// We continue to not add empty information to the config
			// so it's clean and only has required information
			if len(aux) == 0 {
				continue
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
