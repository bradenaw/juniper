package tree

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/xstd/xsort"
)

func TestBasic(t *testing.T) {
	tree := newTree[int](xsort.OrderedLess[int])

	_, ok := tree.Get(5)
	require.False(t, ok)
	tree.Put(5)
	_, ok = tree.Get(5)
	require.True(t, ok)
}
