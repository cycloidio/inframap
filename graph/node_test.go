package graph_test

import (
	"testing"

	"github.com/cycloidio/inframap/graph"
	"github.com/stretchr/testify/assert"
)

func TestAddGroupIDs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		n := graph.Node{}
		n.AddGroupIDs("a")
		assert.Equal(t, []string{"a"}, n.GroupIDs)
		n.AddGroupIDs("b")
		assert.Equal(t, []string{"a", "b"}, n.GroupIDs)
		n.AddGroupIDs("b")
		assert.Equal(t, []string{"a", "b"}, n.GroupIDs)
		n.AddGroupIDs("a", "b", "c")
		assert.Equal(t, []string{"a", "b", "c"}, n.GroupIDs)
	})
}
