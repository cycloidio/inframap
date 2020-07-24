package printer

import (
	"io"

	"github.com/cycloidio/inframap/graph"
)

// Printer is an abstraction to Print the graph.Graph
// in different formats
type Printer interface {
	Print(g *graph.Graph, opt Options, w io.Writer) error
}
