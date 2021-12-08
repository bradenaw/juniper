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

func FuzzBasic(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		tree := newTree[byte](xsort.OrderedLess[byte])
		oracle := make(map[byte]struct{})
		for i := range b {
			item := b[i] & 0x7F
			switch (b[i] & 0xB0) >> 6 {
			case 0:
				tree.Put(item)
				oracle[item] = struct{}{}
			case 1:
				tree.Delete(item)
				delete(oracle, item)
			case 2:
				_, treeOk := tree.Get(item)
				_, oracleOk := oracle[item]
				require.Equal(t, treeOk, oracleOk)
			case 3:
				require.Equal(t, tree.size, len(oracle))
			}
		}
	})
}
