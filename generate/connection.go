package generate

import (
	"github.com/cycloidio/inframap/graph"
	"github.com/cycloidio/inframap/provider"
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

// findEdgeConnections for each edge on n returns the closests connections to a Node.
// The last 'connection' on the []*connection will be a valid provider Node.
// If there are more than one short path they will also be returend
// The [][][]*connection is:
// * First [] is for each edge of the n
// * Second [] is for the shortest connections
// * Third [] is for the connections building a shortest path to another Node
func findEdgeConnections(g *graph.Graph, n *graph.Node, visited map[string]struct{}, opt Options) ([][][]*connection, error) {
	res := make([][][]*connection, 0)
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

		// If it's IM we do not want to try to merge or anything
		// so we just continue to the next edge
		if pv.Type() == provider.IM {
			continue
		}

		visited[e.ID] = struct{}{}
		// If it's a Node we just add it
		if pv.IsNode(rs) {
			res = append(res, [][]*connection{
				{
					&connection{
						Node:      en,
						Direction: direc,
					},
				},
			})
		} else {

			// We make a copy of the visited so
			// each nested edge does not share the
			// same copy of it and produces a wrong
			// result
			aux := make(map[string]struct{})
			for k, v := range visited {
				aux[k] = v
			}

			// We get all the shortest path to a Node and append it
			cons, err := getShortestNodePath(g, en, aux, opt)
			if err != nil {
				return nil, err
			}
			// If there are no connections means it's an end Node,
			// and as it is not a Node (it'll have entered the previous
			// condition) it means it has no Node at the end so we do not
			// return it
			if len(cons) != 0 {
				for i, cs := range cons {
					// Append the first edge node
					// to the first position of the
					// connections
					con := []*connection{
						&connection{
							Node:      en,
							Direction: direc,
						},
					}
					cons[i] = append(con, cs...)
				}

				res = append(res, cons)
			}
		}
	}
	return res, nil
}

// getShortestNodePath get the shortest path to a Node starting from n
func getShortestNodePath(g *graph.Graph, n *graph.Node, visited map[string]struct{}, opt Options) ([][]*connection, error) {
	edges, err := findEdgeConnections(g, n, visited, opt)
	if err != nil {
		return nil, err
	}

	// From all the possible Edges we take the
	// shortest path
	shortestCons := make([][]*connection, 0, 0)
	if len(edges) > 0 {
		shortestCons = [][]*connection{
			edges[0][0],
		}
		for _, cs := range edges {
			for _, cons := range cs {
				// Compare for the shortest path
				if len(shortestCons[0]) > len(cons) {
					shortestCons = [][]*connection{
						cons,
					}
				} else if len(shortestCons[0]) == len(cons) {
					shortestCons = append(shortestCons, cons)
				}
			}
		}
	}
	return shortestCons, nil
}
