//go:build go1.18

package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/maps"
	"github.com/bradenaw/juniper/slices"
	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xsort"
)

func FuzzTree(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		tree := newTree[byte, int](xsort.OrderedLess[byte])
		oracle := make(map[byte]int)
		cursor := tree.Cursor()
		var oracleCursor *byte

		// If k is non-nil, returns the greatest key in oracle less than k.
		// If k is nil, returns the last key in oracle.
		// The return is nil if there is no such key.
		oracleLastLess := func(k *byte) *byte {
			var out *byte
			for existingK := range oracle {
				existingK := existingK
				if k != nil && existingK >= *k {
					continue
				}
				if out == nil || existingK > *out {
					out = &existingK
				}
			}
			return out
		}
		// If k is non-nil, returns the lowest key in oracle greater than k.
		// If k is nil, returns the first key in oracle.
		// The return is nil if there is no such key.
		oracleFirstGreater := func(k *byte) *byte {
			var out *byte
			for existingK := range oracle {
				existingK := existingK
				if k != nil && existingK <= *k {
					continue
				}
				if out == nil || existingK < *out {
					out = &existingK
				}
			}
			return out
		}

		sortKVs := func(kvs []KVPair[byte, int]) {
			xsort.Slice(kvs, func(a, b KVPair[byte, int]) bool {
				return a.K < b.K
			})
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
				sortKVs(oracleSlice)
				t.Logf("oracle: %#v", oracleSlice)
				if oracleCursor == nil {
					t.Log("oracle cursor: nil")
				} else {
					t.Logf("oracle cursor: %#v", *oracleCursor)
				}

				checkTree(t, tree)

				c := tree.Cursor()
				treeSlice := iterator.Collect[KVPair[byte, int]](c.Forward())
				require.Equal(t, oracleSlice, treeSlice)

				if oracleCursor == nil {
					require.False(t, cursor.Ok())
				} else {
					require.True(t, cursor.Ok())
					require.Equal(t, *oracleCursor, cursor.Key())

					_, oracleOk := oracle[*oracleCursor]
					if oracleOk {
						require.Equal(t, oracle[*oracleCursor], cursor.Value())
					}
				}

				if len(oracle) > 0 {
					expectedFirst := oracleFirstGreater(nil)
					firstK, _ := tree.First()
					require.Equal(t, *expectedFirst, firstK)

					expectedLast := oracleLastLess(nil)
					lastK, _ := tree.Last()
					require.Equal(t, *expectedLast, lastK)
				}
			},
			func(k byte, v int) {
				t.Logf("tree.Put(%#v, %#v)", k, v)
				tree.Put(k, v)
				oracle[k] = v
			},
			func(k byte) {
				t.Logf("tree.Delete(%#v)", k)
				tree.Delete(k)
				delete(oracle, k)
				if oracleCursor != nil && k == *oracleCursor {
					oracleCursor = oracleFirstGreater(oracleCursor)
				}
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
				t.Log("cursor.Next()")
				cursor.Next()
				if oracleCursor != nil {
					oracleCursor = oracleFirstGreater(oracleCursor)
				}
			},
			func() {
				t.Log("cursor.Prev()")
				cursor.Prev()
				if oracleCursor != nil {
					oracleCursor = oracleLastLess(oracleCursor)
				}
			},
			func() {
				t.Log("cursor.SeekFirst()")
				cursor.SeekFirst()
				oracleCursor = oracleFirstGreater(nil)
			},
			func(k byte) {
				t.Logf("tree.SeekLastLess(%#v)", k)
				cursor.SeekLastLess(k)
				oracleCursor = oracleLastLess(&k)
			},
			func(k byte) {
				t.Logf("tree.SeekLastLessOrEqual(%#v)", k)
				cursor.SeekLastLessOrEqual(k)
				_, ok := oracle[k]
				if ok {
					oracleCursor = &k
				} else {
					oracleCursor = oracleLastLess(&k)
				}
			},
			func(k byte) {
				t.Logf("tree.SeekFirstGreaterOrEqual(%#v)", k)
				cursor.SeekFirstGreaterOrEqual(k)
				_, ok := oracle[k]
				if ok {
					oracleCursor = &k
				} else {
					oracleCursor = oracleFirstGreater(&k)
				}
			},
			func(k byte) {
				t.Logf("tree.SeekFirstGreater(%#v)", k)
				cursor.SeekFirstGreater(k)
				oracleCursor = oracleFirstGreater(&k)
			},
			func() {
				t.Log("cursor.SeekLast()")
				cursor.SeekLast()
				oracleCursor = oracleLastLess(nil)
			},
			func() {
				t.Log("cursor.Forward()")
				kvs := iterator.Collect(cursor.Forward())
				var expectedKVs []KVPair[byte, int]
				if oracleCursor != nil {
					expectedKVs = slices.Map(
						slices.Filter(
							maps.Keys(oracle),
							func(k byte) bool { return k >= *oracleCursor },
						),
						func(k byte) KVPair[byte, int] {
							return KVPair[byte, int]{k, oracle[k]}
						},
					)
					sortKVs(expectedKVs)
					if len(expectedKVs) == 0 {
						expectedKVs = nil
					}
				}
				require.Equal(t, expectedKVs, kvs)
			},
			func() {
				t.Log("cursor.Backward()")
				kvs := iterator.Collect(cursor.Backward())
				var expectedKVs []KVPair[byte, int]
				if oracleCursor != nil {
					expectedKVs = slices.Map(
						slices.Filter(
							maps.Keys(oracle),
							func(k byte) bool { return k <= *oracleCursor },
						),
						func(k byte) KVPair[byte, int] {
							return KVPair[byte, int]{k, oracle[k]}
						},
					)
					sortKVs(expectedKVs)
					slices.Reverse(expectedKVs)
					if len(expectedKVs) == 0 {
						expectedKVs = nil
					}
				}
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
