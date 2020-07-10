package dot

import (
	"bytes"
	"fmt"
	"io"

	"github.com/awalterschulze/gographviz"
	"github.com/cycloidio/infraview/factory"
	"github.com/cycloidio/infraview/graph"
	"github.com/cycloidio/infraview/provider"
)

type Dot struct{}

func (d Dot) Print(g *graph.Graph, w io.Writer) error {
	graph := gographviz.NewGraph()
	parentName := "G"
	graph.SetName(parentName)
	graph.SetDir(true)
	graph.SetStrict(true)

	for _, n := range g.Nodes {
		// If it's nil the pv, it means we do not know it so we'll use
		// the RawProvider.
		// We do not use the infraview.GenerateOptions to see if it was
		// marked as Raw as we could then see the edges distiction
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
