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
}
