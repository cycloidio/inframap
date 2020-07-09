package printer

import (
	"io"

	"github.com/cycloidio/infraview/graph"
)

type Printer interface {
	Print(g *graph.Graph, w io.Writer) error
}
