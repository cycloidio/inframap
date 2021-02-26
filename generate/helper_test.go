package generate_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/cycloidio/inframap/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertEqualGraph will compare the expected graph to the actual one, the way it'll do that is by ignoring the IDs
// and only using the canonicals of the Node, same with the Edges
func assertEqualGraph(t *testing.T, expected, actual *graph.Graph, actualCfg map[string]interface{}) {

	assert.Len(t, actual.Nodes, len(expected.Nodes), "Nodes")
	assert.Len(t, actual.Edges, len(expected.Edges), "Edges")

	// nodeCans holds canonical -> graph.Node
	nodeCans := make(map[string]*graph.Node)

	// nodeIDs holds ID -> canonical
	nodeIDs := make(map[string]string)

	// edgeDirections it's a combination of "source+target" key
	edgeDirections := make(map[string]*graph.Edge)

	// edgeCans holds all the edge canonicals
	edgeCans := make(map[string]struct{})

	for _, n := range actual.Nodes {
		nodeCans[n.Canonical] = n
		nodeIDs[n.ID] = n.Canonical
	}

	for _, e := range actual.Edges {
		so, ok := nodeIDs[e.Source]
		require.True(t, ok, "The ID %q is missing", e.Source)

		ta, ok := nodeIDs[e.Target]
		require.True(t, ok, "The ID %q is missing", e.Target)

		edgeDirections[formatEdgeKey(so, ta)] = e

		for _, c := range e.Canonicals {
			edgeCans[c] = struct{}{}
		}
	}

	for _, en := range expected.Nodes {
		if an, ok := nodeCans[en.Canonical]; ok {
			en.ID = an.ID
			en.TFID = an.TFID
			en.Resource = an.Resource
			en.Weight = an.Weight
			assert.Equal(t, en, an)
		} else {
			assert.Failf(t, "Fail", "The Node with Canonical %q is missing", en.Canonical)
		}
	}

	for _, ee := range expected.Edges {
		if e, ok := edgeDirections[formatEdgeKey(ee.Source, ee.Target)]; !ok {
			if _, ok = edgeDirections[formatEdgeKey(ee.Target, ee.Source)]; ok {
				assert.Failf(t, "Fail", "The Edge with Source %q and Target %q is present but in the other direction", ee.Source, ee.Target)
			} else {
				assert.Failf(t, "Fail", "The Edge with Source %q and Target %q is missing", ee.Source, ee.Target)
			}
		} else {
			sort.Strings(e.Canonicals)
			sort.Strings(ee.Canonicals)
			assert.Equal(t, ee.Canonicals, e.Canonicals, fmt.Sprintf("Source: %s Target: %s", ee.Source, ee.Target))
		}
	}

	// As the FromHCL does not have config
	// the validation for nil has to be done before, here
	// it'll be ignored if nil
	if actualCfg != nil {
		actualCans := make([]string, 0, 0)
		for k, v := range actualCfg["resource"].(map[string]interface{}) {
			for n := range v.(map[string]interface{}) {
				actualCans = append(actualCans, fmt.Sprintf("%s.%s", k, n))
			}
		}

		expectedCans := make([]string, 0, 0)
		for k := range edgeCans {
			expectedCans = append(expectedCans, k)
		}
		for k := range nodeCans {
			expectedCans = append(expectedCans, k)
		}

		sort.Strings(actualCans)
		sort.Strings(expectedCans)

		assert.Equal(t, expectedCans, actualCans)
	}

}

// formatEdgeKey returns the "%s--%s" of s and t
func formatEdgeKey(s, t string) string {
	return fmt.Sprintf("%s--%s", s, t)
}
