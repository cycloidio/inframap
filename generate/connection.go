package generate

import (
	"github.com/cycloidio/inframap/graph"
)

type direction int

const (
	positive direction = 1
	negative direction = -1
)

// connection represents the Node and Direction of the connection(edge)
// It's always respective to the root Node in which it was created from
// An example would be: RN [-> N1] The content of the [] is what the
// connection represents.
// The Direction is respective to the RN if it's in favor it is positive( RN -> )
// if not it is negative ( RN <- )
type connection struct {
	Node *graph.Node
	// Direction is related to the
	// first node of the chain.
	Direction direction
}

// findEdgeConnections for each edge on n returns the closest connection to a Node.
// The last 'connection' on the []*connection will be a valid provider Node
func findEdgeConnections(g *graph.Graph, n *graph.Node, visited map[string]struct{}, opt Options) ([][]*connection, error) {
	res := make([][]*connection, 0)
	edges := g.GetEdgesForNode(n.ID)
	for _, e := range edges {
		// If we have already visited that Edge
		// we skip it
		if _, ok := visited[e.ID]; ok {
			continue
		}

		// Get the Node on the other
		// part of the Edge
		en, err := g.GetNodeByID(e.Source)
		direc := negative
		if e.Source == n.ID {
			en, err = g.GetNodeByID(e.Target)
			direc = positive
		}
		if err != nil {
			return nil, err
		}

		pv, rs, err := getProviderAndResource(en.Canonical, opt)
		if err != nil {
			return nil, err
		}

		visited[e.ID] = struct{}{}
		// If it's a Node we just add it
		if pv.IsNode(rs) {
			res = append(res, []*connection{
				&connection{
					Node:      en,
					Direction: direc,
				},
			})
		} else {

			// We get all the shortest path to a Node and append it
			cons, err := getShortestNodePath(g, en, visited, opt)
			if err != nil {
				return nil, err
			}
			con := []*connection{
				&connection{
					Node:      en,
					Direction: direc,
				},
			}
			res = append(res, append(con, cons...))
		}
	}
	return res, nil
}

// getShortestNodePath get the shortest path to a Node starting from n
func getShortestNodePath(g *graph.Graph, n *graph.Node, visited map[string]struct{}, opt Options) ([]*connection, error) {
	edges, err := findEdgeConnections(g, n, visited, opt)
	if err != nil {
		return nil, err
	}

	// From all the possible Edges we take the
	// shortest path
	shortestCon := make([]*connection, 0, 0)
	if len(edges) > 0 {
		shortestCon = edges[0]
		for _, cons := range edges {
			// Compare for the shortest path
			if len(shortestCon) > len(cons) {
				shortestCon = cons
			}
		}
	}
	return shortestCon, nil
}
