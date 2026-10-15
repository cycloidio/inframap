package dot_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/awalterschulze/gographviz"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/cycloidio/inframap/generate"
	"github.com/cycloidio/inframap/printer"
	"github.com/cycloidio/inframap/printer/dot"
)

func TestAlternativeNodeNamesPreserveGraph(t *testing.T) {
	for _, name := range []string{"web", ""} {
		t.Run(fmt.Sprintf("name=%q", name), func(t *testing.T) {
			fs := afero.NewMemMapFs()
			config := fmt.Sprintf(`
resource "aws_lb" "first" {
  name = %q
}
resource "aws_lb" "second" {
  name = %q
  depends_on = [aws_lb.first]
}
`, name, name)
			require.NoError(t, afero.WriteFile(fs, "main.tf", []byte(config), 0600))
			g, _, err := generate.FromHCL(fs, "main.tf", generate.Options{Raw: true})
			require.NoError(t, err)
			require.Len(t, g.Nodes, 2)
			require.Len(t, g.Edges, 1)

			for _, alternative := range []bool{false, true} {
				var output bytes.Buffer
				err = (dot.Dot{}).Print(g, printer.Options{AlternativeNodeNames: alternative}, &output)
				require.NoError(t, err)
				printed, err := gographviz.Read(output.Bytes())
				require.NoError(t, err)
				require.Len(t, printed.Nodes.Nodes, 2)
				require.Len(t, printed.Edges.Edges, 1)
				require.Equal(t, `"aws_lb.second"`, printed.Edges.Edges[0].Src)
				require.Equal(t, `"aws_lb.first"`, printed.Edges.Edges[0].Dst)
				for _, node := range printed.Nodes.Nodes {
					if alternative && name != "" {
						require.Equal(t, fmt.Sprintf("%q", name), node.Attrs[gographviz.Label])
					} else {
						require.NotContains(t, node.Attrs, gographviz.Label)
					}
				}
			}
		})
	}
}
