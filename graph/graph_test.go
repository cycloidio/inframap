package graph_test

import (
	"sort"
	"testing"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddNode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n1}, g.Nodes)
	})
	t.Run("RequiredNodeCanonical", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1"}

		err := g.AddNode(n1)
		assert.Error(t, err, errcode.ErrGraphRequiredNodeCanonical.Error())
	})
	t.Run("RequiredNodeID", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{Canonical: "1"}

		err := g.AddNode(n1)
		assert.Error(t, err, errcode.ErrGraphRequiredNodeID.Error())
	})
	t.Run("AlreadyExistsNode", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "1"}
		n2 := &graph.Node{ID: "2", Canonical: "1"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		assert.Error(t, err, errcode.ErrGraphAlreadyExistsNode.Error())
	})
	t.Run("AlreadyExistsNodeID", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "1"}
		n2 := &graph.Node{ID: "1", Canonical: "2"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		assert.Error(t, err, errcode.ErrGraphAlreadyExistsNodeID.Error())
	})
}

func TestRemoveNodeByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "2"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n1, n2}, g.Nodes)

		err = g.RemoveNodeByID(n1.ID)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n2}, g.Nodes)

		err = g.RemoveNodeByID(n2.ID)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{}, g.Nodes)
	})
	t.Run("SuccessWithEdges", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		e := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddEdge(e)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Edge{e}, g.Edges)

		err = g.RemoveNodeByID(n1.ID)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n2}, g.Nodes)
		assert.Equal(t, []*graph.Edge{}, g.Edges)
	})
	t.Run("SuccessWithEdges", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		n3 := &graph.Node{ID: "3", Canonical: "can3"}
		n4 := &graph.Node{ID: "4", Canonical: "can4"}
		e1 := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}
		e2 := &graph.Edge{ID: "2", Source: n2.ID, Target: n3.ID}
		e3 := &graph.Edge{ID: "3", Source: n1.ID, Target: n4.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddNode(n3)
		require.NoError(t, err)

		err = g.AddNode(n4)
		require.NoError(t, err)

		err = g.AddEdge(e1)
		require.NoError(t, err)

		err = g.AddEdge(e2)
		require.NoError(t, err)

		err = g.AddEdge(e3)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n1, n2, n3, n4}, g.Nodes)
		assert.Equal(t, []*graph.Edge{e1, e2, e3}, g.Edges)

		err = g.RemoveNodeByID(n2.ID)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n1, n3, n4}, g.Nodes)
		assert.Equal(t, []*graph.Edge{e3}, g.Edges)
	})
	t.Run("ErrGraphNotFoundNode", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.RemoveNodeByID("invalid-id")
		assert.Error(t, err, errcode.ErrGraphNotFoundNode.Error())
	})
}

func TestAddEdge(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		e := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddEdge(e)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Edge{e}, g.Edges)
	})
	t.Run("RequiredEdgeID", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		e := &graph.Edge{Source: n1.ID, Target: n2.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddEdge(e)
		assert.Error(t, err, errcode.ErrGraphRequiredEdgeID.Error())
	})
	t.Run("RequiredEdgeTarget", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		e := &graph.Edge{ID: "1", Source: n1.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddEdge(e)
		assert.Error(t, err, errcode.ErrGraphRequiredEdgeTarget.Error())
	})
	t.Run("RequiredEdgeSource", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		e := &graph.Edge{ID: "1", Target: n1.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddEdge(e)
		assert.Error(t, err, errcode.ErrGraphRequiredEdgeSource.Error())
	})
	t.Run("InvalidFormatEdgeTargetNotFound", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		e := &graph.Edge{ID: "1", Source: n1.ID, Target: "pepito"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddEdge(e)
		assert.Error(t, err, errcode.ErrGraphNotFoundEdgeTarget.Error())
	})
	t.Run("InvalidFormatEdgeSourceNotFound", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		e := &graph.Edge{ID: "1", Target: n1.ID, Source: "pepito"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddEdge(e)
		assert.Error(t, err, errcode.ErrGraphNotFoundEdgeSource.Error())
	})
	t.Run("AlreadyExistsEdge", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		e := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddEdge(e)
		require.NoError(t, err)

		err = g.AddEdge(e)
		assert.Error(t, err, errcode.ErrGraphAlreadyExistsEdge.Error())
	})
	t.Run("AlreadyExistsEdgeID", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		n3 := &graph.Node{ID: "3", Canonical: "can3"}
		e := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}
		e2 := &graph.Edge{ID: "1", Source: n1.ID, Target: n3.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddNode(n3)
		require.NoError(t, err)

		err = g.AddEdge(e)
		require.NoError(t, err)

		err = g.AddEdge(e2)
		assert.Error(t, err, errcode.ErrGraphAlreadyExistsEdgeID.Error())
	})
}

func TestGetEdgesForNode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		n3 := &graph.Node{ID: "3", Canonical: "can3"}
		e1 := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}
		e2 := &graph.Edge{ID: "2", Source: n2.ID, Target: n3.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddNode(n3)
		require.NoError(t, err)

		err = g.AddEdge(e1)
		require.NoError(t, err)

		err = g.AddEdge(e2)
		require.NoError(t, err)

		edges := g.GetEdgesForNode(n1.ID)
		assert.Equal(t, []*graph.Edge{e1}, edges)

		edges = g.GetEdgesForNode(n2.ID)
		assert.Equal(t, []*graph.Edge{e1, e2}, edges)

		edges = g.GetEdgesForNode(n3.ID)
		assert.Equal(t, []*graph.Edge{e2}, edges)
	})
}

func TestClean(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		n3 := &graph.Node{ID: "3", Canonical: "can3"}
		n4 := &graph.Node{ID: "4", Canonical: "can4"}
		e := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n4)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddNode(n3)
		require.NoError(t, err)

		err = g.AddEdge(e)
		require.NoError(t, err)

		assert.Len(t, g.Nodes, 4)

		g.Clean()

		assert.Equal(t, []*graph.Node{n1, n2}, g.Nodes)
	})
}

func TestNodeReplace(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		n3 := &graph.Node{ID: "3", Canonical: "can3"}
		e1 := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}
		e2 := &graph.Edge{ID: "2", Source: n2.ID, Target: n3.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddNode(n3)
		require.NoError(t, err)

		err = g.AddEdge(e1)
		require.NoError(t, err)

		err = g.AddEdge(e2)
		require.NoError(t, err)

		err = g.Replace("2", "3")
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n1, n3}, g.Nodes)
		assert.Equal(t, g.Edges[0].ID, e1.ID)
		assert.Equal(t, g.Edges[0].Source, n1.ID)
		assert.Equal(t, g.Edges[0].Target, n3.ID)
		assert.Equal(t, []string{"can2"}, g.Edges[0].Canonicals)

		assert.Equal(t, []*graph.Edge{g.Edges[0]}, g.GetEdgesForNode("3"))
		assert.Equal(t, []*graph.Edge{g.Edges[0]}, g.GetEdgesForNode("1"))

		assert.Len(t, g.Edges, 1)
	})
	t.Run("SuccessWithMultipleEdgesToReconnect", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		n3 := &graph.Node{ID: "3", Canonical: "can3"}
		n4 := &graph.Node{ID: "4", Canonical: "can4"}
		n5 := &graph.Node{ID: "5", Canonical: "can5"}
		e1 := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}
		e2 := &graph.Edge{ID: "2", Source: n2.ID, Target: n3.ID}
		e3 := &graph.Edge{ID: "3", Source: n2.ID, Target: n4.ID}
		e4 := &graph.Edge{ID: "4", Source: n2.ID, Target: n5.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddNode(n3)
		require.NoError(t, err)

		err = g.AddNode(n4)
		require.NoError(t, err)

		err = g.AddNode(n5)
		require.NoError(t, err)

		err = g.AddEdge(e1)
		require.NoError(t, err)

		err = g.AddEdge(e2)
		require.NoError(t, err)

		err = g.AddEdge(e3)
		require.NoError(t, err)

		err = g.AddEdge(e4)
		require.NoError(t, err)

		err = g.Replace("2", "3")
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n1, n3, n4, n5}, g.Nodes)
		assert.Equal(t, g.Edges[0].ID, e1.ID)
		assert.Equal(t, g.Edges[0].Source, n1.ID)
		assert.Equal(t, g.Edges[0].Target, n3.ID)
		assert.Equal(t, []string{"can2"}, g.Edges[0].Canonicals)

		assert.Equal(t, g.Edges[1].ID, e3.ID)
		assert.Equal(t, g.Edges[1].Source, n3.ID)
		assert.Equal(t, g.Edges[1].Target, n4.ID)
		assert.Equal(t, []string{"can2"}, g.Edges[1].Canonicals)

		assert.Equal(t, g.Edges[2].ID, e4.ID)
		assert.Equal(t, g.Edges[2].Source, n3.ID)
		assert.Equal(t, g.Edges[2].Target, n5.ID)
		assert.Equal(t, []string{"can2"}, g.Edges[2].Canonicals)

		assert.Len(t, g.GetEdgesForNode("3"), 3)

		assert.Equal(t, []*graph.Edge{g.Edges[0]}, g.GetEdgesForNode("1"))
		assert.Equal(t, g.Edges, g.GetEdgesForNode("3"))
		assert.Equal(t, []*graph.Edge{g.Edges[1]}, g.GetEdgesForNode("4"))
		assert.Equal(t, []*graph.Edge{g.Edges[2]}, g.GetEdgesForNode("5"))

		assert.Len(t, g.Edges, 3)
	})
	t.Run("SuccessWith2InternalReplaces", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		n3 := &graph.Node{ID: "3", Canonical: "can3"}
		e1 := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID, Canonicals: []string{"e1"}}
		e2 := &graph.Edge{ID: "2", Source: n2.ID, Target: n3.ID, Canonicals: []string{"e2"}}
		e3 := &graph.Edge{ID: "3", Source: n1.ID, Target: n3.ID, Canonicals: []string{"e3"}}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddNode(n3)
		require.NoError(t, err)

		err = g.AddEdge(e1)
		require.NoError(t, err)

		err = g.AddEdge(e2)
		require.NoError(t, err)

		err = g.AddEdge(e3)
		require.NoError(t, err)

		err = g.Replace("2", "1")
		require.NoError(t, err)

		assert.Equal(t, []*graph.Node{n1, n3}, g.Nodes)
		assert.Equal(t, g.Edges[0].ID, e3.ID)
		assert.Equal(t, g.Edges[0].Source, n1.ID)
		assert.Equal(t, g.Edges[0].Target, n3.ID)
		sort.Strings(g.Edges[0].Canonicals)
		assert.Equal(t, []string{"can2", "e1", "e2", "e3"}, g.Edges[0].Canonicals)

		assert.Equal(t, []*graph.Edge{g.Edges[0]}, g.GetEdgesForNode("3"))
		assert.Equal(t, []*graph.Edge{g.Edges[0]}, g.GetEdgesForNode("1"))

		assert.Len(t, g.Edges, 1)
	})
}

func TestGetNodeByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		n, err := g.GetNodeByID("1")
		require.NoError(t, err)
		assert.Equal(t, n1, n)

		n, err = g.GetNodeByID("2")
		require.NoError(t, err)
		assert.Equal(t, n2, n)
	})
	t.Run("InvalidFormatNodeNotFound", func(t *testing.T) {
		g := graph.New()

		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		_, err = g.GetNodeByID("3")
		assert.Error(t, err, errcode.ErrGraphNotFoundNode.Error())
	})
}

func TestGetNodeByCanonical(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		n, err := g.GetNodeByCanonical("can1")
		require.NoError(t, err)
		assert.Equal(t, n1, n)

		n, err = g.GetNodeByCanonical("can2")
		require.NoError(t, err)
		assert.Equal(t, n2, n)
	})
	t.Run("InvalidFormatNodeNotFound", func(t *testing.T) {
		g := graph.New()

		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		_, err = g.GetNodeByCanonical("can3")
		assert.Error(t, err, errcode.ErrGraphNotFoundNode.Error())
	})
}

func TestInvertEdge(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g := graph.New()
		n1 := &graph.Node{ID: "1", Canonical: "can1"}
		n2 := &graph.Node{ID: "2", Canonical: "can2"}
		n3 := &graph.Node{ID: "3", Canonical: "can3"}
		e1 := &graph.Edge{ID: "1", Source: n1.ID, Target: n2.ID}
		e2 := &graph.Edge{ID: "2", Source: n2.ID, Target: n3.ID}

		err := g.AddNode(n1)
		require.NoError(t, err)

		err = g.AddNode(n2)
		require.NoError(t, err)

		err = g.AddNode(n3)
		require.NoError(t, err)

		err = g.AddEdge(e1)
		require.NoError(t, err)

		err = g.AddEdge(e2)
		require.NoError(t, err)

		g.InvertEdge(e1.ID)
		require.NoError(t, err)

		assert.Equal(t, []*graph.Edge{
			&graph.Edge{
				ID:     "1",
				Source: n2.ID,
				Target: n1.ID,
			},
			&graph.Edge{
				ID:     "2",
				Source: n2.ID,
				Target: n3.ID,
			},
		}, g.Edges)
	})
}
