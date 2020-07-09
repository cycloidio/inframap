package dot

import (
	"bytes"
	"fmt"
	"io"

	"github.com/awalterschulze/gographviz"
	"github.com/cycloidio/infraview/factory"
	"github.com/cycloidio/infraview/graph"
)

type Dot struct{}

func (d Dot) Print(g *graph.Graph, w io.Writer) error {
	graph := gographviz.NewGraph()
	parentName := "G"
	graph.SetName(parentName)
	graph.SetDir(true)
	graph.SetStrict(true)

	for _, n := range g.Nodes {
		pv, rs, _ := factory.GetProviderAndResource(n.Canonical)
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
