package generate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/factory"
	"github.com/cycloidio/inframap/graph"
	"github.com/cycloidio/inframap/provider"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states/statefile"
	uuid "github.com/satori/go.uuid"
)

// Options are the possible options
// that can be used to generate a Graph
type Options struct {
	// Raw means the RawProvider instead of the
	// specific one
	Raw bool

	// Clean means that the Nodes that do not have
	// any connection will be "removed"
	Clean bool
}

// FromState generate a graph.Graph from the tfstate applying the opt
func FromState(tfstate json.RawMessage, opt Options) (*graph.Graph, map[string]interface{}, error) {
	buf := bytes.NewBuffer(tfstate)
	file, err := statefile.Read(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("error while reading TFState: %w", err)
	}

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
			if rv.Addr.Mode != addrs.ManagedResourceMode {
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
				if len(iv.Current.DependsOn) != 0 {
					deps = append(deps, instanceCurrentDependsOnToString(iv.Current.DependsOn)...)
				}

				aux := make(map[string]interface{})
				if iv.Current.AttrsJSON != nil {
					// For TF 0.12
					err = json.Unmarshal(iv.Current.AttrsJSON, &aux)
					if err != nil {
						return nil, nil, fmt.Errorf("invalid fomrat JSON for resource %q with AttrsJSON %s: %w", string(iv.Current.AttrsJSON), rk, err)
					}
				} else {
					// For TF 0.11
					// AttrsFlat it's a map[string]string so it has to be converted
					// to map[string]interface{} to fit on the aux definition
					for k, v := range iv.Current.AttrsFlat {
						aux[k] = v
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

	if opt.Clean {
		g.Clean()
	}

	err = fixEdges(g, cfg, opt)
	if err != nil {
		return nil, nil, err
	}

	err = mutate(g, opt)
	if err != nil {
		return nil, nil, err
	}

	cleanHangingEdges(g, opt)

	endCfg, err := buildConfig(g, cfg, nodeCanIDs)
	if err != nil {
		return nil, nil, err
	}
	return g, endCfg, nil
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
			if len(edges) > 0 {
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
var reSG = regexp.MustCompile(`\$\{(?P<type>[^\.][a-z0-9-_]+)\.(?P<name>[^\.][a-z0-9-_]+)\.(?P<attr>[a-z0-9-_]+)\}`)

// fixEdges tries to fix the direction of the edges that was done based on the 'depends_on'
// to something more Provider dependent by reading the actual config.
// This would mean that in case of AWS it'll read the Config and if it's a SG it'll
// read the actual direction of it and apply it and potentially changing the Edges
// directions
func fixEdges(g *graph.Graph, cfg map[string]map[string]interface{}, opt Options) error {
	for _, n := range g.Nodes {
		pv, rs, err := getProviderAndResource(n.Canonical, opt)
		if err != nil {
			return err
		}

		if pv.IsEdge(rs) {
			edges := g.GetEdgesForNode(n.ID)
			ins, outs := pv.ResourceInOut(rs, cfg[n.ID])

			// For the ins we have to check if any of the edges Target
			// is this ID and reverse it because we want it to be the Source
			for _, in := range ins {
				isID := false
				res := reSG.FindAllStringSubmatch(in, -1)
				resMap := make(map[string]string)
				if len(res) == 0 {
					isID = true
				} else {
					for i, k := range reSG.SubexpNames() {
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
				res := reSG.FindAllStringSubmatch(out, -1)
				resMap := make(map[string]string)
				if len(res) == 0 {
					isID = true
				} else {
					for i, k := range reSG.SubexpNames() {
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

		// calculate the best Node to start with
		// the mutation
		conns[n.ID] = make([]*connection, 0, 0)
		for _, connections := range nodes {
			// We prioritize the most positive ones
			if len(connections) > 1 && ((sumConnsDirection(conns[n.ID]) <= sumConnsDirection(connections)) || len(conns[n.ID]) == 0) {
				conns[n.ID] = connections

				if bestNode == nil || ((sumConnsDirection(bestNodeConns) <= sumConnsDirection(connections)) && bestNode.Weight <= n.Weight) {
					bestNode = n
					bestNodeConns = connections
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
				return fmt.Errorf("missing edge with srcID %q and repID %q: %w", n.ID, con.Node.ID, errcode.ErrGraphRequiredEdgeBetweenNodes)
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

func instanceCurrentDependenciesToString(deps []addrs.AbsResource) []string {
	res := make([]string, 0, len(deps))
	for _, d := range deps {
		res = append(res, d.String())
	}
	return res
}

func instanceCurrentDependsOnToString(deps []addrs.Referenceable) []string {
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
			if rv.Addr.Mode != addrs.ManagedResourceMode {
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
