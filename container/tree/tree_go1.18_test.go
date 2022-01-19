//go:build go1.18

package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/internal/orderedhashmap"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xsort"
)

func orderedhashmapKVPairToKVPair[K any, V any](kv orderedhashmap.KVPair[byte, int]) KVPair[byte, int] {
	return KVPair[byte, int]{kv.K, kv.V}
}

func FuzzTree(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		tree := newTree[byte, int](xsort.OrderedLess[byte])
		cursor := tree.Cursor()

		oracle := orderedhashmap.NewMap[byte, int](xsort.OrderedLess[byte])
		oracleCursor := oracle.Cursor()

		ctr := 0
		nextV := func() int {
			ctr++
			return ctr
		}

		fuzz.Operations(
			b,
			func() { // check
				t.Log("check")
				t.Log(treeToString(tree))

				require.Equal(t, oracle.Len(), tree.size)

				oracleSlice := iterator.Collect(iterator.Map(
					oracle.Iterate(),
					orderedhashmapKVPairToKVPair[byte, int],
				))
				t.Logf("oracle: %#v", oracleSlice)

				checkTree(t, tree)

				c := tree.Cursor()
				treeSlice := iterator.Collect[KVPair[byte, int]](c.Forward())
				require.Equal(t, oracleSlice, treeSlice)

				if !oracleCursor.Ok() {
					require.False(t, cursor.Ok())
				} else {
					require.True(t, cursor.Ok())
					require.Equal(t, oracleCursor.Key(), cursor.Key())
					require.Equal(t, oracle.Get(oracleCursor.Key()), cursor.Value())
				}
			},
			func(k byte) {
				v := nextV()
				t.Logf("tree.Put(%#v, %#v)", k, v)
				tree.Put(k, v)
				oracle.Put(k, v)
			},
			func(k byte, d byte) {
				t.Logf("tree.Delete(%#v)", k)
				if oracleCursor.Ok() && oracleCursor.Key() == k {
					if d >= 128 {
						t.Log("cursor.Next()")
						cursor.Next()
						oracleCursor.Next()
					} else {
						t.Log("cursor.Prev()")
						cursor.Prev()
						oracleCursor.Prev()
					}
				}
				tree.Delete(k)
				oracle.Delete(k)
			},
			func(k byte) {
				t.Logf("tree.Contains(%#v)", k)
				treeOk := tree.Contains(k)
				oracleOk := oracle.Contains(k)
				require.Equal(t, oracleOk, treeOk)
			},
			func(k byte) {
				t.Logf("tree.Get(%#v)", k)
				v := tree.Get(k)
				expectedV := oracle.Get(k)
				require.Equal(t, expectedV, v)
			},
			func() {
				t.Logf("tree.First()")
				k, v := tree.First()
				expectedK, expectedV := oracle.First()
				require.Equal(t, expectedK, k)
				require.Equal(t, expectedV, v)
			},
			func() {
				t.Logf("tree.Last()")
				k, v := tree.Last()
				expectedK, expectedV := oracle.Last()
				require.Equal(t, expectedK, k)
				require.Equal(t, expectedV, v)
			},
			func() {
				t.Log("cursor.Next()")
				cursor.Next()
				oracleCursor.Next()
			},
			func() {
				t.Log("cursor.Prev()")
				cursor.Prev()
				oracleCursor.Prev()
			},
			func() {
				t.Log("cursor.SeekFirst()")
				cursor.SeekFirst()
				oracleCursor.SeekFirst()
			},
			func(k byte) {
				t.Logf("cursor.SeekLastLess(%#v)", k)
				cursor.SeekLastLess(k)
				oracleCursor.SeekLastLess(k)
			},
			func(k byte) {
				t.Logf("cursor.SeekLastLessOrEqual(%#v)", k)
				cursor.SeekLastLessOrEqual(k)
				oracleCursor.SeekLastLessOrEqual(k)
			},
			func(k byte) {
				t.Logf("cursor.SeekFirstGreaterOrEqual(%#v)", k)
				cursor.SeekFirstGreaterOrEqual(k)
				oracleCursor.SeekFirstGreaterOrEqual(k)
			},
			func(k byte) {
				t.Logf("cursor.SeekFirstGreater(%#v)", k)
				cursor.SeekFirstGreater(k)
				oracleCursor.SeekFirstGreater(k)
			},
			func() {
				t.Log("cursor.SeekLast()")
				cursor.SeekLast()
				oracleCursor.SeekLast()
			},
			func() {
				t.Log("cursor.Forward()")
				kvs := iterator.Collect(cursor.Forward())
				expectedKVs := iterator.Collect(iterator.Map(
					oracleCursor.Forward(),
					orderedhashmapKVPairToKVPair[byte, int],
				))
				require.Equal(t, expectedKVs, kvs)
			},
			func() {
				t.Log("cursor.Backward()")
				kvs := iterator.Collect(cursor.Backward())
				expectedKVs := iterator.Collect(iterator.Map(
					oracleCursor.Backward(),
					orderedhashmapKVPairToKVPair[byte, int],
				))
				require.Equal(t, expectedKVs, kvs)
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
		require.Equal(t, curr, curr.left.parent)
	}
	if curr.right != nil {
		require.True(t, tree.less(curr.key, curr.right.key))
		require.Equal(t, curr, curr.right.parent)
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
	if tree.root != nil {
		require.Nil(t, tree.root.parent)
	}
}

func treeToString[K any, V any](t *tree[K, V]) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "tree ====\n")
	var visit func(x *node[K, V], prefix string, descPrefix string)
	seen := make(map[*node[K, V]]struct{})
	visit = func(x *node[K, V], prefix string, descPrefix string) {
		_, ok := seen[x]
		if ok {
			panic(fmt.Sprintf("%s\ncycle detected: already saw %#v", sb.String(), x.key))
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
