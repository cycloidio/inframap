package graph

import "github.com/cycloidio/tfdocs/resource"

// Node defines the standard format of an Edge
type Node struct {
	// ID it'a a random UUID
	ID string

	// Canonical it's 'aws_lb.front' format
	Canonical string

	// Position holds the position of the node if any
	Position []int // x, y

	// TFID it's the internal ID it has on TF
	TFID string

	// Resource it's the information of the resource
	// it holds
	Resource resource.Resource

	// Weight is the addition of the Directions
	// of the Node
	Weight int

	// GroupIDs is the IDs of all the groups
	// it may belong to
	GroupIDs []string

	// mGroupIDs is used to know which GroupIDs
	// are already on the GroupIDs slice so we
	// do not repeat them
	mGroupIDs map[string]struct{}
}

// AddGroupIDs adds the ids to the internal list, if
// one is repeated it'll be ignored
func (n *Node) AddGroupIDs(gids ...string) {
	if n.mGroupIDs == nil {
		n.mGroupIDs = make(map[string]struct{})
	}

	for _, c := range gids {
		if _, ok := n.mGroupIDs[c]; !ok {
			n.mGroupIDs[c] = struct{}{}
			n.GroupIDs = append(n.GroupIDs, c)
		}
	}
}
