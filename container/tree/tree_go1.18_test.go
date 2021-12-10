//go:build go1.18

package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/xstd/iterator"
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
	const (
		ActPut byte = iota << 6
		ActDelete
		ActContains
		ActIterNext
	)

	f.Add([]byte{
		ActPut & 3,
	})
	f.Add([]byte{
		ActPut & 3,
		ActPut & 5,
	})
	f.Add([]byte{
		ActContains & 5,
		ActPut & 5,
		ActContains & 5,
		ActDelete & 5,
		ActContains & 5,
	})

	f.Fuzz(func(t *testing.T, b []byte) {
		tree := newTree(xsort.OrderedLess[byte])
		oracle := make(map[byte]struct{})
		var iter *iterator.Iterator[byte]
		var iterLastSeen *byte
		for i := range b {
			item := b[i] & 0b00111111
			switch b[i] & 0b11000000 {
			case ActPut:
				t.Logf("tree.Put(%#v)", item)
				tree.Put(item)
				oracle[item] = struct{}{}
			case ActDelete:
				t.Logf("tree.Delete(%#v)", item)
				tree.Delete(item)
				delete(oracle, item)
			case ActContains:
				t.Logf("tree.Contains(%#v)", item)
				_, treeOk := tree.Get(item)
				_, oracleOk := oracle[item]
				require.Equal(t, treeOk, oracleOk)
				require.Equal(t, tree.Contains(item), oracleOk)
			case ActIterNext:
				if iterLastSeen == nil {
					t.Logf("iter.Next() from beginning")
				} else {
					t.Logf("iter.Next() after %#v", *iterLastSeen)
				}

				for {
					if iter == nil {
						x := tree.Iterate()
						iter = &x
						iterLastSeen = nil
					}
					ok := (*iter).Next()
					if len(oracle) == 0 {
						require.False(t, ok)
						iter = nil
						break
					}
					var expected *byte
					for k := range oracle {
						k := k
						if iterLastSeen != nil && k <= *iterLastSeen {
							continue
						}
						if expected == nil || k < *expected {
							expected = &k
						}
					}
					require.Equal(t, expected != nil, ok)
					if ok {
						item := (*iter).Item()
						require.Equal(t, *expected, item)
						iterLastSeen = &item
						t.Logf(" -> %#v", item)
						break
					} else {
						t.Logf(" -> (end)")
						iter = nil
						iterLastSeen = nil
					}
				}
			default:
				panic("no action?")
			}
			t.Log(treeToString(tree))
		}

		t.Log("check...")
		var oracleSlice []byte
		for item := range oracle {
			oracleSlice = append(oracleSlice, item)
		}
		xsort.Slice(oracleSlice, xsort.OrderedLess[byte])
		treeSlice := iterator.Collect(tree.Iterate())
		require.Equal(t, treeSlice, oracleSlice)
	})
}

func treeToString[T any](t *tree[T]) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "tree ====\n")
	var visit func(x *node[T], prefix string, descPrefix string)
	visit = func(x *node[T], prefix string, descPrefix string) {
		fmt.Fprintf(&sb, "%s%#v\n", prefix, x.value)
		if x.left != nil {
			visit(x.left, descPrefix+"  l ", descPrefix+"    ")
		}
		if x.right != nil {
			visit(x.right, descPrefix+"  r ", descPrefix+"    ")
		}
	}
	if t.root != nil {
		visit(t.root, "", "")
	}
	fmt.Fprintf(&sb, "=========")
	return sb.String()
}
