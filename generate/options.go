package generate

// Options are the possible options
// that can be used to generate a Graph
type Options struct {
	// Raw means the RawProvider instead of the
	// specific one
	Raw bool

	// Clean means that the Nodes that do not have
	// any connection will be "removed"
	Clean bool

	// Connections toggles the Provider logic for
	// merging Edges between Nodes into one
	Connections bool
}
