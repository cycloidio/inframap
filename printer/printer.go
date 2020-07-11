package printer

import (
	"io"

	"github.com/cycloidio/infraview/graph"
)

// Printer is an abstraction to Print the graph.Graph
// in different formats
type Printer interface {
	Print(g *graph.Graph, w io.Writer) error
}
