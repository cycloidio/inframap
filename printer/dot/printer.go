package dot

import (
	"bytes"
	"fmt"
	"io"

	"github.com/awalterschulze/gographviz"
	"github.com/cycloidio/inframap/graph"
	"github.com/cycloidio/inframap/provider"
	"github.com/cycloidio/inframap/provider/factory"
)

// Dot is the struct that implements
// the Printer of Dot format
type Dot struct{}

// Print prints into w the g in DOT format
func (d Dot) Print(g *graph.Graph, w io.Writer) error {
	graph := gographviz.NewGraph()
	parentName := "G"
	graph.SetName(parentName)
	graph.SetDir(true)
	graph.SetStrict(true)

	for _, n := range g.Nodes {
		// If it's nil the pv, it means we do not know it so we'll use
		// the RawProvider.
		// We do not use the inframap.Options to see if it was
		// marked as Raw as we could then see the edges distinction
		// on the output if needed
		pv, rs, _ := factory.GetProviderAndResource(n.Canonical)
		if pv == nil {
			pv = provider.RawProvider{}
		}
		shape := "ellipse"
		if pv.IsEdge(rs) {
			shape = "rectangle"
		}
		graph.AddNode(parentName, fmt.Sprintf("%q", n.Canonical), map[string]string{
			"shape": shape,
		})
	}
	for _, e := range g.Edges {
		src, _ := g.GetNodeByID(e.Source)
		tr, _ := g.GetNodeByID(e.Target)
		graph.AddEdge(fmt.Sprintf("%q", src.Canonical), fmt.Sprintf("%q", tr.Canonical), true, nil)
	}

	buff := bytes.NewBufferString(graph.String())
	io.Copy(w, buff)

	return nil
}
