package generate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/graph"
	"github.com/cycloidio/inframap/provider"
	"github.com/cycloidio/inframap/provider/factory"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/states/statefile"
	uuid "github.com/satori/go.uuid"
)

// FromState generate a graph.Graph from the tfstate applying the opt
func FromState(tfstate json.RawMessage, opt Options) (*graph.Graph, map[string]interface{}, error) {
	// Since TF 0.13 'depends_on' has been dropped, so we do a manual
	// replace from '"depends_on"' to '"dependencies"'
	tfstate = bytes.ReplaceAll(tfstate, []byte("\"depends_on\""), []byte("\"dependencies\""))
	err := validateTFStateVersion(tfstate)
	if err != nil {
		return nil, nil, fmt.Errorf("error while validating TFStateVersion: %w", err)
	}

	buf := bytes.NewBuffer(tfstate)

	file, err := statefile.Read(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("error while reading TFState: %w", err)
	}

	migrateVersions(tfstate, file)

	g := graph.New()

	// cfg holds the actual configuration of each element
	// it's represented as: ID -> Attrs
	cfg := make(map[string]map[string]interface{})

	// nodeCanIDs holds as key the `aws_alb.front` (graph.Node.Canonical)
	// and as value the UUID (graph.Node.ID) we give to it
	nodeCanIDs := make(map[string][]string)

	// nodeIDEdges holds as key the UUID (graph.Node.ID) and as value
	// all the edges it has, in this case it's the `depends_on` values
	// that we find on the TFState
	nodeIDEdges := make(map[string][]string)

	if !opt.Raw {
		opt, err = checkProviders(file, opt)
		if err != nil {
			return nil, nil, err
		}
	}

	for _, v := range file.State.Modules {
		for rk, rv := range v.Resources {
			// If it's not a Resource we ignore it
			if rv.Addr.Resource.Mode != addrs.ManagedResourceMode {
				continue
			}

			pv, rs, err := getProviderAndResource(rk, opt)
			if err != nil {
				if errors.Is(err, errcode.ErrProviderNotFound) {
					continue
				}
				return nil, nil, err
			}

			// If it's not a Node or Edge we ignore it
			if !pv.IsNode(rs) && !pv.IsEdge(rs) {
				continue
			}

			// The Instances is the representation of the
			// 'count' on the Instance, could also be a 'for_each'
			for id, iv := range rv.Instances {
				if id != nil && id.String() != "[0]" {
					continue
				}
				deps := make([]string, 0)
				if len(iv.Current.Dependencies) != 0 {
					deps = append(deps, instanceCurrentDependenciesToString(iv.Current.Dependencies)...)
				}

				aux := make(map[string]interface{})
				if iv.Current.AttrsJSON != nil {
					// For TF +0.12
					err = json.Unmarshal(iv.Current.AttrsJSON, &aux)
					if err != nil {
						return nil, nil, fmt.Errorf("invalid fomrat JSON for resource %q with AttrsJSON %s: %w", string(iv.Current.AttrsJSON), rk, err)
					}
				} else {
					// For TF 0.11
					// AttrsFlat it's flatmap, so we'll use the flatmap.Expand
					// to abstract the content, we'll first get the list of all
					// the unique keys so then we can Expand each key
					keys := make(map[string]struct{})
					for k := range iv.Current.AttrsFlat {
						kk := strings.Split(k, ".")
						keys[kk[0]] = struct{}{}
					}

					for k := range keys {
						aux[k] = flatmap.Expand(iv.Current.AttrsFlat, k)
					}
				}

				res, err := pv.Resource(rs)
				if err != nil {
					return nil, nil, err
				}
				n := &graph.Node{
					ID:        uuid.NewV4().String(),
					Canonical: rk,
					TFID:      aux["id"].(string),
					Resource:  *res,
				}

				err = g.AddNode(n)
				if err != nil {
					return nil, nil, err
				}

				nodeCanIDs[n.Canonical] = append(nodeCanIDs[n.Canonical], n.ID)
				nodeIDEdges[n.ID] = deps
				cfg[n.ID] = aux
			}
		}
	}

	for sourceID, edges := range nodeIDEdges {
		edgeIDs := make([]string, 0)
		for _, e := range edges {
			if IDs, ok := nodeCanIDs[e]; ok {
				for _, nid := range IDs {
					edgeIDs = append(edgeIDs, nid)
				}
			}
		}

		for _, targetID := range edgeIDs {
			err := g.AddEdge(&graph.Edge{
				ID:     uuid.NewV4().String(),
				Source: sourceID,
				Target: targetID,
			})
			if err != nil {
				return nil, nil, err
			}
		}
	}
	// call the preprocess method for each
	// TF provider in the file
	if err := preprocess(g, cfg, opt); err != nil {
		return nil, nil, err
	}

	if opt.Clean {
		g.Clean()
	}

	err = fixEdges(g, cfg, opt)
	if err != nil {
		return nil, nil, err
	}

	if opt.Connections {
		err = mutate(g, opt)
		if err != nil {
			return nil, nil, err
		}
	}

	if opt.Clean {
		err = cleanHangingEdges(g, opt)
		if err != nil {
			return nil, nil, err
		}
	}

	endCfg, err := buildConfig(g, cfg, nodeCanIDs)
	if err != nil {
		return nil, nil, err
	}
	return g, endCfg, nil
}

// migrateVersions will try to apply migrations of old
// statefile:
// * For version 3 will try to populate the Dependencies as they are
//   removed by TF as they cannot be migrated from v3->v4
func migrateVersions(src []byte, f *statefile.File) error {
	version, err := tfStateVersion(src)
	if err != nil {
		return fmt.Errorf("could not read version: %w", err)
	}

	deps := make(map[string][]addrs.ConfigResource)

	var state map[string]interface{}

	if version == 3 {
		err := json.Unmarshal(src, &state)
		if err != nil {
			return fmt.Errorf("error while reading TFState: %w", err)
		}

		for _, m := range state["modules"].([]interface{}) {
			rs, ok := m.(map[string]interface{})["resources"]
			if !ok {
				continue
			}
			for rn, rv := range rs.(map[string]interface{}) {
				// As we have replaced from 'depends_on' to 'dependencies'
				// we have to access it with 'dependencies'
				dps, ok := rv.(map[string]interface{})["dependencies"]
				if !ok && dps == nil {
					continue
				}
				for _, d := range dps.([]interface{}) {
					ds := strings.Split(d.(string), ".")
					deps[rn] = append(deps[rn], addrs.ConfigResource{
						Resource: addrs.Resource{
							Mode: addrs.ManagedResourceMode,
							Type: ds[0],
							Name: ds[1],
						},
					})
				}
			}
		}

		for _, v := range f.State.Modules {
			for rk, rv := range v.Resources {
				for _, iv := range rv.Instances {
					if d, ok := deps[rk]; ok {
						iv.Current.Dependencies = d
					}
				}
			}
		}

	}

	return nil
}

// validateTFStateVersion validates that the version is the
// one we support which is only 3 and 4
func validateTFStateVersion(b []byte) error {
	v, err := tfStateVersion(b)
	if err != nil {
		return fmt.Errorf("could not read version: %w", err)
	}

	if v != 4 && v != 3 {
		return fmt.Errorf("could not read version %d: %w", v, errcode.ErrInvalidTFStateVersion)
	}

	return nil
}

func tfStateVersion(b []byte) (uint64, error) {
	var v struct {
		Version uint64 `json:"version"`
	}

	err := json.Unmarshal(b, &v)
	if err != nil {
		return 0, fmt.Errorf("error while reading TFState version: %w", err)
	}

	return v.Version, nil
}

// cleanHangingEdges will Replace all the hanging Edges (that are Nodes now)
// with any of the closest Nodes
// The only case in which this happens is when an Edge is connected to the internet
// in which then it has N->E without another Node to connect with (which would be internet)
// it is a usecase we do not support now so we remove them
func cleanHangingEdges(g *graph.Graph, opt Options) error {
	for _, n := range g.Nodes {
		pv, rs, err := getProviderAndResource(n.Canonical, opt)
		if err != nil {
			return err
		}

		if pv.IsEdge(rs) {
			edges := g.GetEdgesForNode(n.ID)
			// If it's a hanging Edge then it should only
			// have one connection, if it has more than
			// one it means there is something else wrong
			if len(edges) == 1 {
				repID := edges[0].Source
				if edges[0].Source == n.ID {
					repID = edges[0].Target
				}
				err = g.Replace(n.ID, repID)
				if err != nil {
					return fmt.Errorf("could not replace node: %w", err)
				}
			}
		}
	}

	// Now that all the hanging Edges have been cleaned
	// we'll remove all the other Edges present if those
	// are not meant to be displayed
	if opt.Connections {
	RESTART:
		for _, n := range g.Nodes {
			pv, rs, err := getProviderAndResource(n.Canonical, opt)
			if err != nil {
				return err
			}
			if pv.IsEdge(rs) {
				if err := g.RemoveNodeByID(n.ID); err != nil {
					return err
				}
				// We restart the loop because this operations potentially
				// change the g.Nodes order/items
				goto RESTART
			}
		}
	}
	return nil
}

// buildConfig takes the Graph g and the config cfg with the nodeCanIDs and returns a config with
// the configuration of each resource mapped to the canonical they have
func buildConfig(g *graph.Graph, cfg map[string]map[string]interface{}, nodeCanIDs map[string][]string) (map[string]interface{}, error) {
	endCfg := make(map[string]interface{})

	for _, n := range g.Nodes {
		c, ok := cfg[n.ID]
		if !ok {
			return nil, fmt.Errorf("could not find config of node %q: %w", n.Canonical, errcode.ErrInvalidTFStateFile)
		}

		// path[0] == resource Type ex: aws_security_group
		// path[1] == resource Name ex: front-port80
		path := strings.Split(n.Canonical, ".")

		if _, ok := endCfg[path[0]]; !ok {
			endCfg[path[0]] = make(map[string]interface{})
		}

		if _, ok := endCfg[path[0]].(map[string]interface{})[path[1]]; ok {
			// If we have it already set, then it's not a valid config
			return nil, fmt.Errorf("repeated config node for %q: %w", n.Canonical, errcode.ErrInvalidTFStateFile)
		}

		endCfg[path[0]].(map[string]interface{})[path[1]] = c
	}

	for _, e := range g.Edges {
		for _, can := range e.Canonicals {
			pv, _, err := factory.GetProviderAndResource(can)
			if err != nil {
				return nil, fmt.Errorf("could not get provider from %s: %w", can, err)
			}
			// The IM resources do not have any config
			if pv.Type() == provider.IM {
				continue
			}
			ids, ok := nodeCanIDs[can]
			if !ok {
				return nil, fmt.Errorf("could not find the ID of the canonical %q: %w", can, errcode.ErrInvalidTFStateFile)
			}
			// We only do the ID 0 for now
			c, ok := cfg[ids[0]]
			if !ok {
				return nil, fmt.Errorf("could not find config of the Node %q: %w", can, errcode.ErrInvalidTFStateFile)
			}

			// path[0] == resource Type ex: aws_security_group
			// path[1] == resource Name ex: front-port80
			path := strings.Split(can, ".")

			if _, ok := endCfg[path[0]]; !ok {
				endCfg[path[0]] = make(map[string]interface{})
			}

			if _, ok := endCfg[path[0]].(map[string]interface{})[path[1]]; ok {
				// As a connection canonical can be shared between different
				// connections this will happen so we ignore it
				continue
			}

			endCfg[path[0]].(map[string]interface{})[path[1]] = c
		}
	}

	endCfg = map[string]interface{}{
		"resource": endCfg,
	}

	return endCfg, nil
}

// reVariable matches ${aws_security_group.front.id}
var reVariable = regexp.MustCompile(`\$\{(?P<type>[^\.][a-z0-9-_]+)\.(?P<name>[^\.][a-z0-9-_]+)\.(?P<attr>[a-z0-9-_]+)\}`)

// fixEdges tries to fix the direction of the edges that was done based on the 'depends_on'
// to something more Provider dependent by reading the actual config.
// This would mean that in case of AWS it'll read the Config and if it's a SG it'll
// read the actual direction of it and apply it and potentially changing the Edges
// directions
func fixEdges(g *graph.Graph, cfg map[string]map[string]interface{}, opt Options) error {
	// imNodes holds Node.ID => []string{NodeCan}
	imNodes := make(map[string][]string)
	for _, n := range g.Nodes {
		pv, rs, err := getProviderAndResource(n.Canonical, opt)
		if err != nil {
			return err
		}

		if pv.IsEdge(rs) {
			edges := g.GetEdgesForNode(n.ID)
			ins, outs, nodes := pv.ResourceInOutNodes(n.ID, rs, cfg)
			if opt.ExternalNodes {
				imNodes[n.ID] = nodes
			}

			// For the ins we have to check if any of the edges Target
			// is this ID and reverse it because we want it to be the Source
			for _, in := range ins {
				isID := false
				res := reVariable.FindAllStringSubmatch(in, -1)
				resMap := make(map[string]string)
				if len(res) == 0 {
					isID = true
				} else {
					for i, k := range reVariable.SubexpNames() {
						if res[0][i] != "" {
							resMap[k] = res[0][i]
						}
					}
					in = fmt.Sprintf("%s.%s", resMap["type"], resMap["name"])
				}
				for _, e := range edges {
					rep, err := g.GetNodeByID(e.Target)
					if err != nil {
						return err
					}

					if isID && rep.TFID == in {
						g.InvertEdge(e.ID)
					} else if !isID && rep.Canonical == in {
						g.InvertEdge(e.ID)
					}
				}
			}

			// For the outs we have to check if any of the edges Sources
			// is this ID and reverse it because we want it to be the Target
			for _, out := range outs {
				isID := false
				res := reVariable.FindAllStringSubmatch(out, -1)
				resMap := make(map[string]string)
				if len(res) == 0 {
					isID = true
				} else {
					for i, k := range reVariable.SubexpNames() {
						if res[0][i] != "" {
							resMap[k] = res[0][i]
						}
					}
					out = fmt.Sprintf("%s.%s", resMap["type"], resMap["name"])
				}
				for _, e := range edges {
					rep, err := g.GetNodeByID(e.Source)
					if err != nil {
						return err
					}

					if isID && rep.TFID == out {
						g.InvertEdge(e.ID)
					} else if !isID && rep.Canonical == out {
						g.InvertEdge(e.ID)
					}
				}
			}

		}
	}

	for id, nodes := range imNodes {
		edges := g.GetEdgesForNode(id)
		for _, e := range edges {
			// We only want edges that have the
			// current Node (that is Edge) as Target
			if e.Target != id {
				continue
			}

			es, err := g.GetNodeByID(e.Source)
			if err != nil {
				return err
			}
			pv, _, err := getProviderAndResource(es.Canonical, opt)
			if err != nil {
				return err
			}
			// If the source (which will be the target on the edge creation)
			// is of type IM we will not create the edge
			if pv.Type() == provider.IM {
				continue
			}
			for _, no := range nodes {
				pv, _, err := getProviderAndResource(no, opt)
				if err != nil {
					return err
				}

				res, err := pv.Resource(no)
				if err != nil {
					return err
				}

				newn := &graph.Node{
					ID:        uuid.NewV4().String(),
					Canonical: no,
					TFID:      no,
					Resource:  *res,
				}

				err = g.AddNode(newn)
				if err != nil {
					if !errors.Is(err, errcode.ErrGraphAlreadyExistsNode) {
						return err
					}

					newn, err = g.GetNodeByCanonical(newn.Canonical)
					if err != nil {
						return fmt.Errorf("could not add Node because of: %w", err)
					}
				}

				cfg[newn.ID] = map[string]interface{}{}

				err = g.AddEdge(&graph.Edge{
					ID:     uuid.NewV4().String(),
					Source: newn.ID,
					Target: e.Source,
				})
				if err != nil && !errors.Is(err, errcode.ErrGraphAlreadyExistsEdge) {
					return err
				}
			}
		}
	}
	return nil
}

// sumConnsDirection returns the total sum of all the
// directions of the conns
func sumConnsDirection(conns []*connection) int {
	var res int
	for _, c := range conns {
		res += int(c.Direction)
	}
	return res
}

// mutate will mutate the Graph by merging the Nodes that are Edges on the Provider they belong
// with the actual Nodes, at the end it'll leave a Graph with just Nodes (Provider Nodes)
func mutate(g *graph.Graph, opt Options) error {
	conns := make(map[string][]*connection)
	var bestNode *graph.Node
	var bestNodeConns []*connection
	// First of all we calculate all the shortest connections of all the Nodes
	// that are actually Nodes on the Provider.
	// From that we also get which is the bestNode to start with (most positive directions)
	// and we start with that one and the bestNodeConns
	// In case of a tie, we use the Node.Weigh which is the result
	// of all the Replaces and the Direction of those replaces
	for _, n := range g.Nodes {
		pv, rs, err := getProviderAndResource(n.Canonical, opt)
		if err != nil {
			return err
		}

		// If it's not a Node we continue
		// we only want to mutate the Nodes
		if !pv.IsNode(rs) {
			continue
		}

		// It is a Node
		nodes, err := findEdgeConnections(g, n, make(map[string]struct{}), opt)
		if err != nil {
			return fmt.Errorf("could not findEdgeConnections: %w", err)
		}

		directConnections := make(map[string]struct{})
		for _, e := range g.GetEdgesForNode(n.ID) {
			// Get the Node on the other
			// part of the Edge
			en, err := g.GetNodeByID(e.Source)
			if e.Source == n.ID {
				en, err = g.GetNodeByID(e.Target)
			}
			if err != nil {
				return err
			}

			pv, rs, err := getProviderAndResource(en.Canonical, opt)
			if err != nil {
				return err
			}

			if pv.IsNode(rs) {
				directConnections[en.ID] = struct{}{}
			}
		}
		// calculate the best Node to start with
		// the mutation
		conns[n.ID] = make([]*connection, 0, 0)
		for _, cs := range nodes {
			for i, connections := range cs {
				// We prioritize the most positive ones
				if len(connections) > 1 && ((sumConnsDirection(conns[n.ID]) <= sumConnsDirection(connections)) || len(conns[n.ID]) == 0) {
					// If the Edge that we are selecting has the end node be already a
					// direct connection to the main n, we do not want to select it
					isLastAndEmpty := (i == len(cs)-1) && bestNode == nil
					if _, ok := directConnections[connections[len(connections)-1].Node.ID]; ok && !isLastAndEmpty {
						continue
					}

					var (
						nedges, bedges int
					)

					if len(bestNodeConns) != 0 {
						// Get the number of edges that the fist Node we are trying to merge have on them
						// this will tell us "how relevant" they are, a higher value means that it's more
						// important so it is another tie breaker in favor of more important nodes first
						nedges = len(g.GetEdgesForNode(n.ID))
						bedges = len(g.GetEdgesForNode(bestNode.ID))

					}

					conns[n.ID] = connections

					if bestNode == nil || ((sumConnsDirection(bestNodeConns) <= sumConnsDirection(connections)) && bestNode.Weight <= n.Weight && bedges <= nedges) {
						bestNode = n
						bestNodeConns = connections
					}
				}
			}
		}
	}

	if bestNode == nil {
		// We have finished, all Nodes are connected to another Node
		return nil
	}

	n := bestNode
	var direc int
	var err error
	// For all the Connections of the bestNode we
	// Replace them all to make just one connection
	// between the 2 Nodes
	for i, con := range bestNodeConns {
		// If it's the last Node it means it's
		// the actual Node (not Edge) that we
		// want to join with.
		if len(bestNodeConns)-1 == i {
			direc += int(con.Direction)
			edges := g.GetEdgesForNode(n.ID)
			var edge *graph.Edge
			for _, e := range edges {
				if e.Source == con.Node.ID && e.Target == n.ID {
					// If the Node is Target means that the actual direction
					// is negative. So if we have decided that it should be
					// positive we have to invert the edge
					if direc > 0 {
						g.InvertEdge(e.ID)
					}
					edge = e
					break
				} else if e.Source == n.ID && e.Target == con.Node.ID {
					// If the Node is Target means that the actual direction
					// is positive. So if we have decided that it should be
					// negative we have to invert the edge
					if direc < 0 {
						g.InvertEdge(e.ID)
					}
					edge = e
					break
				}
			}
			if edge == nil {
				// If we are missing cyclic connection (to itself) could be related that the graph is missing
				// some node and the end result ended with a cyclic that was not cyclic at the end based
				// on how the 'mutate' works (merging by directions).
				if n.Canonical != con.Node.Canonical {
					return fmt.Errorf("missing edge with srcCan %q and repCan %q: %w", n.Canonical, con.Node.Canonical, errcode.ErrGraphRequiredEdgeBetweenNodes)
				}
			}
		} else {
			direc += int(con.Direction)
			// If the next node is in the same direction we
			// can just replace without issue
			if bestNodeConns[i+1].Direction == con.Direction {
				err = g.Replace(con.Node.ID, n.ID)
				if err != nil {
					return fmt.Errorf("could not replace edges: %w", err)
				}
			} else {
				// If the next node is in a different direction
				// then we have to merge the other way around
				err = g.Replace(con.Node.ID, bestNodeConns[i+1].Node.ID)
				if err != nil {
					return fmt.Errorf("could not replace edges: %w", err)
				}
			}
		}
	}
	// As we have mutated something, we restart again
	// to get the next best Node
	return mutate(g, opt)
}

func instanceCurrentDependenciesToString(deps []addrs.ConfigResource) []string {
	res := make([]string, 0, len(deps))
	for _, d := range deps {
		res = append(res, d.String())
	}
	return res
}

// getProviderAndResource uses factory.Options but if the opt.Raw is defined
// it'll return the RawProvider. It's a helper to not repeat all time the same logic
func getProviderAndResource(rk string, opt Options) (provider.Provider, string, error) {
	var (
		pv  provider.Provider
		rs  string
		err error
	)

	if opt.Raw {
		pv = provider.RawProvider{}
		rs = strings.Split(rk, ".")[0]
	} else {
		pv, rs, err = factory.GetProviderAndResource(rk)
	}

	return pv, rs, err
}

// checkProviders checks if we support any of the Providers from f, if not it'll set
// the opt.Raw to true so it can be used with Raw instead of returning an empty Graph
func checkProviders(f *statefile.File, opt Options) (Options, error) {
	for _, v := range f.State.Modules {
		for rk, rv := range v.Resources {
			// If it's not a Resource we ignore it
			if rv.Addr.Resource.Mode != addrs.ManagedResourceMode {
				continue
			}

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
	}

	// If we reach here means the we do not support the providers
	// of the TFState, so Raw has to be used
	opt.Raw = true

	return opt, nil
}

// preprocess will call PreProcess method of each TF provider supported in the
// config
func preprocess(g *graph.Graph, cfg map[string]map[string]interface{}, opt Options) error {
	visitedProviders := make(map[provider.Type]struct{}, 0)
	for _, node := range g.Nodes {
		pv, _, err := getProviderAndResource(node.Canonical, opt)
		if err != nil {
			// TF provider not found, no need to
			// continue or to raise an error
			continue
		}

		if _, ok := visitedProviders[pv.Type()]; ok {
			continue
		}

		edges := pv.PreProcess(cfg)
		for _, edge := range edges {
			err := g.AddEdge(&graph.Edge{
				ID:     uuid.NewV4().String(),
				Source: edge[0],
				Target: edge[1],
			})
			if err != nil {
				// If the edge already exists we can ignore it
				if errors.Is(err, errcode.ErrGraphAlreadyExistsEdge) {
					continue
				}
				return fmt.Errorf("could not add edge: %w", err)
			}
		}

		visitedProviders[pv.Type()] = struct{}{}
	}
	return nil
}
