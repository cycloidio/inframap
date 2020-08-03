package dot

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/awalterschulze/gographviz"
	"github.com/cycloidio/inframap/graph"
	"github.com/cycloidio/inframap/printer"
	"github.com/cycloidio/inframap/provider"
	"github.com/cycloidio/inframap/provider/factory"
	"github.com/markbates/pkger"

	// As we require to load the assets to be used
	// we import it as empty
	_ "github.com/cycloidio/inframap/assets"
)

// Dot is the struct that implements
// the Printer of Dot format
type Dot struct{}

// Print prints into w the g in DOT format
func (d Dot) Print(g *graph.Graph, opt printer.Options, w io.Writer) error {
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
		attr := map[string]string{
			"shape": "ellipse",
		}
		if pv.IsEdge(rs) {
			attr["shape"] = "rectangle"
		}

		if opt.ShowIcons && n.Resource.Icon != "" {
			ext := path.Ext(n.Resource.Icon)
			pngIcon := fmt.Sprintf("%s.png", n.Resource.Icon[0:len(n.Resource.Icon)-len(ext)])
			assetPath := path.Join("inframap", "assets", pv.Type().String(), pngIcon)
			pathIcon := path.Join(xdg.CacheHome, assetPath)

			attr["image"] = fmt.Sprintf("%q", pathIcon)
			attr["shape"] = "plaintext"
			attr["labelloc"] = "b"
			attr["height"] = "1.15"

			// If the file does not exists on the Cache path, we have to write it,
			// if not it means it's already correct so nothing to be done
			if _, err := os.Stat(pathIcon); os.IsNotExist(err) {
				p, err := xdg.CacheFile(assetPath)
				if err != nil {
					return err
				}

				f, err := os.Create(p)
				if err != nil {
					return err
				}

				iconFile, err := pkger.Open(path.Join("/assets", "icons", pv.Type().String(), pngIcon))
				if err != nil {
					return err
				}

				if _, err = io.Copy(f, iconFile); err != nil {
					return err
				}

				f.Close()
				iconFile.Close()
			}
		}

		for _, g := range n.GroupIDs {
			clusterName := fmt.Sprintf("%q", fmt.Sprintf("cluster_%s", g))
			graph.AddSubGraph(parentName, clusterName, map[string]string{
				"label": clusterName,
			})
			graph.AddNode(clusterName, fmt.Sprintf("%q", n.Canonical), attr)
		}

		graph.AddNode(parentName, fmt.Sprintf("%q", n.Canonical), attr)
	}

	for _, e := range g.Edges {
		src, err := g.GetNodeByID(e.Source)
		if err != nil {
			return err
		}

		tr, err := g.GetNodeByID(e.Target)
		if err != nil {
			return err
		}

		graph.AddEdge(fmt.Sprintf("%q", src.Canonical), fmt.Sprintf("%q", tr.Canonical), true, nil)
	}

	buff := bytes.NewBufferString(graph.String())
	io.Copy(w, buff)

	return nil
}
