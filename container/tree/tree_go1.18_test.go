//go:build go1.18

package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xsort"
)

func FuzzTree(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		tree := newTree[byte, int](xsort.OrderedLess[byte])
		oracle := make(map[byte]int)
		var iter *treeIterator[byte, int]
		var iterLastSeen *byte

		// If k is non-nil, returns the greatest key in oracle less than k.
		// If k is nil, returns the last key in oracle.
		// The second return is false if there is no such key.
		oracleLastLess := func(k *byte) (byte, bool) {
			var out byte
			found := false
			for existingK := range oracle {
				if k != nil && existingK >= *k {
					continue
				}
				if !found || existingK > out {
					out = existingK
					found = true
				}
			}
			return out, found
		}
		// If k is non-nil, returns the lowest key in oracle greater than k.
		// If k is nil, returns the first key in oracle.
		// The second return is false if there is no such key.
		oracleFirstGreater := func(k *byte) (byte, bool) {
			var out byte
			found := false
			for existingK := range oracle {
				if k != nil && existingK <= *k {
					continue
				}
				if !found || existingK < out {
					out = existingK
					found = true
				}
			}
			return out, found
		}

		fuzz.Operations(
			b,
			func() { // check
				t.Log("check")
				t.Log(treeToString(tree))

				require.Equal(t, len(oracle), tree.size)
				var oracleSlice []KVPair[byte, int]
				for k, v := range oracle {
					oracleSlice = append(oracleSlice, KVPair[byte, int]{k, v})
				}
				xsort.Slice(oracleSlice, func(a, b KVPair[byte, int]) bool {
					return a.K < b.K
				})
				t.Logf("oracle: %#v", oracleSlice)

				checkTree(t, tree)

				treeSlice := iterator.Collect[KVPair[byte, int]](tree.Iterate())
				require.Equal(t, oracleSlice, treeSlice)

				if len(oracle) > 0 {
					expectedFirst, _ := oracleFirstGreater(nil)
					firstK, _ := tree.First()
					require.Equal(t, expectedFirst, firstK, "expected first")

					expectedLast, _ := oracleLastLess(nil)
					lastK, _ := tree.Last()
					require.Equal(t, expectedLast, lastK)
				}
			},
			func(k byte, v int) {
				t.Logf("tree.Put(%#v)", k)
				tree.Put(k, v)
				oracle[k] = v
			},
			func(k byte) {
				t.Logf("tree.Delete(%#v)", k)
				tree.Delete(k)
				delete(oracle, k)
			},
			func(k byte) {
				t.Logf("tree.Contains(%#v)", k)
				treeOk := tree.Contains(k)
				_, oracleOk := oracle[k]
				require.Equal(t, oracleOk, treeOk)
			},
			func(k byte) {
				t.Logf("tree.Get(%#v)", k)
				v := tree.Get(k)
				expectedV := oracle[k]
				require.Equal(t, expectedV, v)
			},
			func() {
				if iterLastSeen == nil {
					t.Logf("iter.Next() from beginning")
				} else {
					t.Logf("iter.Next() after %#v", *iterLastSeen)
				}

				if iter == nil {
					iter = tree.Iterate()
					iterLastSeen = nil
				}
				t.Logf("iter state: %#v", iter)
				item, ok := iter.Next()
				expectedK, oracleOK := oracleFirstGreater(iterLastSeen)
				require.Equal(t, oracleOK, ok)
				if ok {
					require.Equal(t, expectedK, item.K)
					require.Equal(t, oracle[expectedK], item.V)
					iterLastSeen = &item.K
					t.Logf(" -> %#v", item)
				} else {
					t.Logf(" -> (end)")
					iter = nil
					iterLastSeen = nil
				}
			},
			func() {
				if iterLastSeen == nil {
					t.Logf("iter.Prev() from end")
				} else {
					t.Logf("iter.Prev() before %#v", *iterLastSeen)
				}

				if iter == nil {
					iter = tree.Iterate()
					iterLastSeen = nil
				}
				item, ok := iter.Prev()
				expectedK, oracleOK := oracleLastLess(iterLastSeen)
				require.Equal(t, oracleOK, ok)
				if ok {
					require.Equal(t, expectedK, item.K)
					require.Equal(t, oracle[expectedK], item.V)
					iterLastSeen = &item.K
					t.Logf(" -> %#v", item)
				} else {
					t.Logf(" -> (end)")
					iter = nil
					iterLastSeen = nil
				}
			},
		)
	})
}

func checkNode[K any, V any](t *testing.T, tree *tree[K, V], curr *node[K, V]) int {
	if curr == nil {
		return 0
	}
	if curr.left != nil {
		require.True(t, tree.less(curr.left.key, curr.key))
	}
	if curr.right != nil {
		require.True(t, tree.less(curr.key, curr.right.key))
	}
	if curr.left == nil && curr.right == nil {
		require.Equalf(t, 0, curr.height, "%#v is a leaf and should have height 0", curr.key)
	} else {
		require.Equalf(
			t,
			xmath.Max(tree.leftHeight(curr), tree.rightHeight(curr))+1,
			curr.height,
			"%#v has wrong height",
			curr.key,
		)
	}
	imbalance := tree.leftHeight(curr) - tree.rightHeight(curr)
	require.LessOrEqual(t, imbalance, 1, fmt.Sprintf("%#v is imbalanced", curr.key))
	require.GreaterOrEqual(t, imbalance, -1, fmt.Sprintf("%#v is imbalanced", curr.key))

	leftSize := checkNode(t, tree, curr.left)
	rightSize := checkNode(t, tree, curr.right)
	return leftSize + rightSize + 1
}
func checkTree[K any, V any](t *testing.T, tree *tree[K, V]) {
	size := checkNode(t, tree, tree.root)
	require.Equal(t, size, tree.size)
}

func treeToString[K any, V any](t *tree[K, V]) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "tree ====\n")
	var visit func(x *node[K, V], prefix string, descPrefix string)
	seen := make(map[*node[K, V]]struct{})
	visit = func(x *node[K, V], prefix string, descPrefix string) {
		_, ok := seen[x]
		if ok {
			panic(fmt.Sprintf("cycle detected: already saw %#v", x.value))
		}
		seen[x] = struct{}{}
		fmt.Fprintf(&sb, "%s%#v h:%d\n", prefix, x.key, x.height)
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
