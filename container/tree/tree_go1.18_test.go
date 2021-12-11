//go:build go1.18

package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xsort"
)

func FuzzTree(f *testing.F) {
	const (
		ActPut byte = iota << 6
		ActDelete
		ActContains
		ActIterNext
	)

	f.Add([]byte{
		ActPut | 3,
	})
	f.Add([]byte{
		ActPut | 3,
		ActPut | 5,
	})
	f.Add([]byte{
		ActContains | 5,
		ActPut | 5,
		ActContains | 5,
		ActDelete | 5,
		ActContains | 5,
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
				treeOk := tree.Contains(item)
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
		checkTree(t, tree)
		var oracleSlice []byte
		for item := range oracle {
			oracleSlice = append(oracleSlice, item)
		}
		xsort.Slice(oracleSlice, xsort.OrderedLess[byte])
		treeSlice := iterator.Collect(tree.Iterate())
		require.Equal(t, treeSlice, oracleSlice)
	})
}

func checkNode[T any](t *testing.T, tree *tree[T], curr *node[T]) int {
	if curr == nil {
		return 0
	}
	if curr.left != nil {
		require.True(t, tree.less(curr.left.value, curr.value))
	}
	if curr.right != nil {
		require.True(t, tree.less(curr.value, curr.right.value))
	}
	if curr.left == nil && curr.right == nil {
		require.Equalf(t, 0, curr.height, "%#v is a leaf and should have height 0", curr.value)
	} else {
		require.Equalf(
			t,
			xmath.Max(tree.leftHeight(curr), tree.rightHeight(curr))+1,
			curr.height,
			"%#v has wrong height",
			curr.value,
		)
	}
	imbalance := tree.leftHeight(curr) - tree.rightHeight(curr)
	require.LessOrEqual(t, imbalance, 1, fmt.Sprintf("%#v is imbalanced", curr.value))
	require.GreaterOrEqual(t, imbalance, -1, fmt.Sprintf("%#v is imbalanced", curr.value))

	leftSize := checkNode(t, tree, curr.left)
	rightSize := checkNode(t, tree, curr.right)
	return leftSize + rightSize + 1
}
func checkTree[T any](t *testing.T, tree *tree[T]) {
	size := checkNode(t, tree, tree.root)
	require.Equal(t, size, tree.size)
}

func treeToString[T any](t *tree[T]) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "tree ====\n")
	var visit func(x *node[T], prefix string, descPrefix string)
	seen := make(map[*node[T]]struct{})
	visit = func(x *node[T], prefix string, descPrefix string) {
		_, ok := seen[x]
		if ok {
			panic(fmt.Sprintf("cycle detected: already saw %#v", x.value))
		}
		seen[x] = struct{}{}
		fmt.Fprintf(&sb, "%s%#v h:%d\n", prefix, x.value, x.height)
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
