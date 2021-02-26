package graph

import (
	"fmt"

	"github.com/cycloidio/inframap/errcode"
)

// Graph defines the standard format of a Graph
type Graph struct {
	Edges []*Edge
	Nodes []*Node

	// nodesCans canonical -> struct{}{}
	nodesCans map[string]*Node

	// nodesIDs id -> struct{}{}
	nodesIDs map[string]*Node

	// nodesWithEdge id -> []*Edge
	nodesWithEdge map[string][]*Edge

	// edgesSourceTarget (source+target) -> struct{}{}
	// used to validate that the direction already exists
	edgesSourceTarget map[string]*Edge

	// edgesIDs id -> struct{}{}
	edgesIDs map[string]*Edge
}

// New returns a new initialized Graph
func New() *Graph {
	return &Graph{
		nodesCans:         make(map[string]*Node),
		nodesIDs:          make(map[string]*Node),
		edgesSourceTarget: make(map[string]*Edge),
		edgesIDs:          make(map[string]*Edge),

		nodesWithEdge: make(map[string][]*Edge),
	}
}

// AddEdge adds an Edge to the Graph
func (g *Graph) AddEdge(e *Edge) error {
	if e.ID == "" {
		return errcode.ErrGraphRequiredEdgeID
	}
	if e.Target == "" {
		return errcode.ErrGraphRequiredEdgeTarget
	}
	if e.Source == "" {
		return errcode.ErrGraphRequiredEdgeSource
	}

	if _, ok := g.nodesIDs[e.Target]; !ok {
		return errcode.ErrGraphNotFoundEdgeTarget
	}

	if _, ok := g.nodesIDs[e.Source]; !ok {
		return errcode.ErrGraphNotFoundEdgeSource
	}

	check := e.Source + e.Target
	if _, ok := g.edgesSourceTarget[check]; ok {
		return errcode.ErrGraphAlreadyExistsEdge
	}

	if _, ok := g.edgesIDs[e.ID]; ok {
		return errcode.ErrGraphAlreadyExistsEdgeID
	}

	g.edgesSourceTarget[check] = e
	g.edgesIDs[e.ID] = e

	g.nodesWithEdge[e.Source] = append(g.nodesWithEdge[e.Source], e)
	g.nodesWithEdge[e.Target] = append(g.nodesWithEdge[e.Target], e)

	g.Edges = append(g.Edges, e)

	return nil
}

// AddNode adds an Node to the Graph
func (g *Graph) AddNode(n *Node) error {
	if n.Canonical == "" {
		return errcode.ErrGraphRequiredNodeCanonical
	}
	if n.ID == "" {
		return errcode.ErrGraphRequiredNodeID
	}

	if _, ok := g.nodesCans[n.Canonical]; ok {
		return fmt.Errorf("with canonical %q: %w", n.Canonical, errcode.ErrGraphAlreadyExistsNode)
	}

	if _, ok := g.nodesIDs[n.ID]; ok {
		return errcode.ErrGraphAlreadyExistsNodeID
	}

	g.nodesCans[n.Canonical] = n
	g.nodesIDs[n.ID] = n

	g.Nodes = append(g.Nodes, n)

	return nil
}

// GetNodeByID returns the requested Node with the nID
func (g *Graph) GetNodeByID(nID string) (*Node, error) {
	n, ok := g.nodesIDs[nID]
	if !ok {
		return nil, errcode.ErrGraphNotFoundNode
	}
	return n, nil
}

// GetNodeByCanonical returns the requested Node with the nCan
func (g *Graph) GetNodeByCanonical(nCan string) (*Node, error) {
	n, ok := g.nodesCans[nCan]
	if !ok {
		return nil, errcode.ErrGraphNotFoundNode
	}
	return n, nil
}

// Clean removes all the Nodes that do not
// have any edge
func (g *Graph) Clean() {
	nodesToRemove := make([]int, 0)
	for i, n := range g.Nodes {
		if _, ok := g.nodesWithEdge[n.ID]; !ok {
			nodesToRemove = append(nodesToRemove, i)
		}
	}

	// For each iteration we have to decrease the next 'idx'
	// by 'i' as we removed 'i' elements
	for i, idx := range nodesToRemove {
		idx -= i
		g.removeNodeByIDX(idx)
	}
}

// GetEdgesForNode returns all the edges that have relation to this nID
func (g *Graph) GetEdgesForNode(nID string) []*Edge {
	return g.nodesWithEdge[nID]
}

// Replace will replace the srcID Node for the repID Node by removing the srcID
// and connecting all the edges from srcID to repID.
// srcID Node and repID Node have to be connected directly
func (g *Graph) Replace(srcID, repID string) error {
	srcEdges := g.GetEdgesForNode(srcID)
	srcNode, err := g.GetNodeByID(srcID)
	if err != nil {
		return err
	}
	repNode, err := g.GetNodeByID(repID)
	if err != nil {
		return err
	}

	// mutualEdge is the edge that connects this 2 Nodes
	var mutualEdge *Edge
	for _, e := range srcEdges {
		// Depending on the direction of the connection
		// we increase or decrease the Node.Weight
		if e.Source == srcID && e.Target == repID {
			repNode.Weight--
			mutualEdge = e
			break
		} else if e.Source == repID && e.Target == srcID {
			repNode.Weight++
			mutualEdge = e
			break
		}
	}

	if mutualEdge == nil {
		return fmt.Errorf("no mutual edge between srcID %q and repID %s: %w", srcID, repID, errcode.ErrGraphRequiredEdgeBetweenNodes)
	}

	for _, e := range srcEdges {
		if e.ID == mutualEdge.ID {
			continue
		}

		// Replace all the connections from the srcID to the repID
		err := e.Replace(srcID, repID)
		if err != nil {
			return err
		}

		e.AddCanonicals(append(mutualEdge.Canonicals, srcNode.Canonical)...)

		ee, okstt := g.edgesSourceTarget[e.Source+e.Target]

		// If the Edge does not exists we register it
		// If it does then we remove it as we do not want repeated edges
		if !okstt {
			g.nodesWithEdge[repID] = append(g.nodesWithEdge[repID], e)
			g.edgesSourceTarget[e.Source+e.Target] = e
		} else {
			// Before removing repeated edges we add the
			// canonicals from the edge we want to delete
			ee.AddCanonicals(e.Canonicals...)
			g.removeEdgeByID(e.ID)
		}
	}

	if err = g.RemoveNodeByID(srcID); err != nil {
		return err
	}

	return nil
}

// InvertEdge inverts the Source and Target of the eID
func (g *Graph) InvertEdge(eID string) {
	for _, e := range g.Edges {
		if e.ID == eID {
			delete(g.edgesSourceTarget, e.Source+e.Target)
			src := e.Source
			e.Source = e.Target
			e.Target = src
			g.edgesSourceTarget[e.Source+e.Target] = e
		}
	}
}

// RemoveNodeByID removes the Node with the ID and the Edges
// associated with it
func (g *Graph) RemoveNodeByID(ID string) error {
	for i, n := range g.Nodes {
		if n.ID == ID {
			idx := i
			g.removeNodeByIDX(idx)
			return nil
		}
	}
	return errcode.ErrGraphNotFoundNode
}

// removeNodeByIDX removes the idx element (via the copy) and then
// removes the last element as it's not needed. It also removes
// the Edges that where connected to this Node
func (g *Graph) removeNodeByIDX(idx int) {
	n := g.Nodes[idx]

	delete(g.nodesCans, n.Canonical)
	delete(g.nodesIDs, n.ID)
	delete(g.nodesWithEdge, n.ID)

	lenNodes := len(g.Nodes)
	copy(g.Nodes[idx:], g.Nodes[idx+1:])
	g.Nodes = g.Nodes[:lenNodes-1]

	// Remove the Edges from the Graph
RESTART:
	for _, e := range g.Edges {
		if e.Target == n.ID || e.Source == n.ID {
			g.removeEdgeByID(e.ID)
			// We restart the loop because this operation potentially
			// changes the g.Edges order/items
			goto RESTART
		}
	}
}

// removeEdgeByID removes the Edge with the ID
func (g *Graph) removeEdgeByID(ID string) {
	for i, e := range g.Edges {
		if e.ID == ID {
			delete(g.edgesSourceTarget, e.Source+e.Target)
			delete(g.edgesIDs, e.ID)

			lenEdges := len(g.Edges)
			copy(g.Edges[i:], g.Edges[i+1:])
			g.Edges = g.Edges[:lenEdges-1]

			// Remove the edge from the list of edges
			// that each node has
			sedges := g.nodesWithEdge[e.Source]
			for ii, ee := range sedges {
				if ee.ID == e.ID {
					lenEdges = len(sedges)
					copy(sedges[ii:], sedges[ii+1:])
					sedges = sedges[:lenEdges-1]
				} else if e.Target == ee.Target && e.Source == ee.Source {
					ee.AddCanonicals(e.Canonicals...)
				}
			}
			g.nodesWithEdge[e.Source] = sedges

			tedges := g.nodesWithEdge[e.Target]
			for ii, ee := range tedges {
				if ee.ID == e.ID {
					lenEdges = len(tedges)
					copy(tedges[ii:], tedges[ii+1:])
					tedges = tedges[:lenEdges-1]
				} else if e.Target == ee.Target && e.Source == ee.Source {
					ee.AddCanonicals(e.Canonicals...)
				}
			}
			g.nodesWithEdge[e.Target] = tedges
		}
	}
}
