package graph_test

import (
	"testing"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEdgeReplace(t *testing.T) {
	tests := []struct {
		Name   string
		Edge   graph.Edge
		Src    string
		Rep    string
		EEdge  graph.Edge
		EError error
	}{
		{
			Name: "SuccessInSource",
			Edge: graph.Edge{ID: "1", Source: "2", Target: "3"},
			Src:  "2", Rep: "1",
			EEdge: graph.Edge{ID: "1", Source: "1", Target: "3"},
		},
		{
			Name: "SuccessInTarget",
			Edge: graph.Edge{ID: "1", Source: "2", Target: "3"},
			Src:  "3", Rep: "1",
			EEdge: graph.Edge{ID: "1", Source: "2", Target: "1"},
		},
		{
			Name: "SuccessEqual",
			Edge: graph.Edge{ID: "1", Source: "2", Target: "3"},
			Src:  "2", Rep: "2",
			EEdge: graph.Edge{ID: "1", Source: "2", Target: "3"},
		},
		{
			Name: "SuccessCircular",
			Edge: graph.Edge{ID: "1", Source: "2", Target: "3"},
			Src:  "3", Rep: "2",
			EEdge: graph.Edge{ID: "1", Source: "2", Target: "2"},
		},
		{
			Name: "ErrorInvalidFormatEdgeSourceNotFound",
			Edge: graph.Edge{ID: "1", Source: "2", Target: "3"},
			Src:  "1", Rep: "2",
			EError: errcode.ErrGraphNotFoundEdgeSource,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Edge.Replace(tt.Src, tt.Rep)
			if tt.EError == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.EEdge.Source, tt.Edge.Source)
				assert.Equal(t, tt.EEdge.Target, tt.Edge.Target)
			} else {
				assert.Error(t, err, tt.EError.Error())
			}
		})
	}
}

func TestAddCanonicals(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		e := graph.Edge{}
		e.AddCanonicals("a")
		assert.Equal(t, []string{"a"}, e.Canonicals)
		e.AddCanonicals("b")
		assert.Equal(t, []string{"a", "b"}, e.Canonicals)
		e.AddCanonicals("b")
		assert.Equal(t, []string{"a", "b"}, e.Canonicals)
		e.AddCanonicals("a", "b", "c")
		assert.Equal(t, []string{"a", "b", "c"}, e.Canonicals)
	})
}
