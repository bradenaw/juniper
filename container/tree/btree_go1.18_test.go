//go:build go1.18

package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/internal/orderedhashmap"
	"github.com/bradenaw/juniper/internal/require2"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/slices"
	"github.com/bradenaw/juniper/xsort"
)

func orderedhashmapKVPairToKVPair[K any, V any](kv orderedhashmap.KVPair[uint16, int]) KVPair[uint16, int] {
	return KVPair[uint16, int]{kv.K, kv.V}
}

func FuzzBtree(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		tree := newBtree[uint16, int](xsort.OrderedLess[uint16])
		cursor := tree.Cursor()
		oracle := orderedhashmap.NewMap[uint16, int](xsort.OrderedLess[uint16])
		oracleCursor := oracle.Cursor()

		ctr := 0

		fuzz.Operations(
			b,
			func() { // check
				t.Log(treeToString(t, tree))

				if oracleCursor.Ok() {
					t.Logf("oracleCursor @ %#v", oracleCursor.Key())
				} else {
					t.Log("oracleCursor off the edge")
				}

				checkTree(t, tree)

				require2.Equal(t, oracle.Len(), tree.Len())

				oraclePairs := iterator.Collect(
					iterator.Map(oracle.Cursor().Forward(), orderedhashmapKVPairToKVPair[uint16, int]),
				)
				xsort.Slice(oraclePairs, func(a, b KVPair[uint16, int]) bool {
					return a.Key < b.Key
				})

				c := tree.Cursor()
				treePairs := iterator.Collect(c.Forward())

				require2.SlicesEqual(t, oraclePairs, treePairs)

				require2.Equalf(t, oracleCursor.Ok(), cursor.Ok(), "cursor.Ok()")
				if oracleCursor.Ok() {
					require2.Equal(t, oracleCursor.Key(), cursor.Key())
					require2.Equal(t, oracleCursor.Value(), cursor.Value())
				}
			},
			func(k uint16) {
				v := ctr
				t.Logf("tree.Put(%#v, %#v)", k, v)
				tree.Put(k, v)
				oracle.Put(k, v)
				ctr++
			},
			func(k uint16) {
				expected := oracle.Get(k)
				t.Logf("tree.Get(%#v) -> %#v", k, expected)
				actual := tree.Get(k)
				require2.Equal(t, expected, actual)
			},
			func(k uint16) {
				t.Logf("tree.Delete(%#v)", k)
				tree.Delete(k)
				oracle.Delete(k)
			},
			func(k uint16) {
				oracleOk := oracle.Contains(k)
				t.Logf("tree.Contains(%#v) -> %t", k, oracleOk)
				treeOk := tree.Contains(k)
				require2.Equal(t, oracleOk, treeOk)
			},
			func() {
				t.Logf("tree.First()")
				k, v := tree.First()
				expectedK, expectedV := oracle.First()
				require2.Equal(t, expectedK, k)
				require2.Equal(t, expectedV, v)
			},
			func() {
				t.Logf("tree.Last()")
				k, v := tree.Last()
				expectedK, expectedV := oracle.Last()
				require2.Equal(t, expectedK, k)
				require2.Equal(t, expectedV, v)
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
			func(k uint16) {
				t.Logf("cursor.SeekLastLess(%#v)", k)
				cursor.SeekLastLess(k)
				oracleCursor.SeekLastLess(k)
			},
			func(k uint16) {
				t.Logf("cursor.SeekLastLessOrEqual(%#v)", k)
				cursor.SeekLastLessOrEqual(k)
				oracleCursor.SeekLastLessOrEqual(k)
			},
			func(k uint16) {
				t.Logf("cursor.SeekFirstGreaterOrEqual(%#v)", k)
				cursor.SeekFirstGreaterOrEqual(k)
				oracleCursor.SeekFirstGreaterOrEqual(k)
			},
			func(k uint16) {
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
					orderedhashmapKVPairToKVPair[uint16, int],
				))
				require2.SlicesEqual(t, expectedKVs, kvs)
			},
			func() {
				t.Log("cursor.Backward()")
				kvs := iterator.Collect(cursor.Backward())
				expectedKVs := iterator.Collect(iterator.Map(
					oracleCursor.Backward(),
					orderedhashmapKVPairToKVPair[uint16, int],
				))
				require2.SlicesEqual(t, expectedKVs, kvs)
			},
		)
	})
}

func isZero[T comparable](t T) bool {
	var zero T
	return t == zero
}

func checkTree[K comparable, V comparable](t *testing.T, tree *btree[K, V]) {
	var checkNode func(x *node[K, V])
	checkNode = func(x *node[K, V]) {
		if x.leaf() {
			for i := 0; i < int(x.n)+1; i++ {
				require2.Nil(t, x.children[i])
			}
		} else {
			for i := 0; i < int(x.n)+1; i++ {
				require2.NotNil(t, x.children[i])
				require2.Truef(
					t,
					x.children[i].parent == x,
					"%p ─child─> %p ─parent─> %p",
					x,
					x.children[i],
					x.children[i].parent,
				)
				checkNode(x.children[i])
			}
			for i := 0; i < int(x.n); i++ {
				left := x.children[i]
				right := x.children[i+1]
				k := x.keys[i]
				require2.True(t, tree.less(left.keys[int(left.n)-1], k))
				require2.True(t, tree.less(k, right.keys[0]))
			}
		}
		if x == tree.root {
			if tree.Len() > 0 {
				require2.GreaterOrEqual(t, int(x.n), 1)
			}
		} else {
			require2.GreaterOrEqual(t, int(x.n), minKVs)
		}
		require2.True(t, xsort.SliceIsSorted(x.keys[:int(x.n)], tree.less))
		require2.True(t, slices.All(x.keys[int(x.n):], isZero[K]))
		require2.True(t, slices.All(x.values[int(x.n):], isZero[V]))
		require2.Truef(
			t,
			slices.All(x.children[int(x.n)+1:], isZero[*node[K, V]]),
			"%p %#v",
			x,
			x.children[int(x.n)+1:],
		)
	}

	require2.NotNil(t, tree.root)
	require2.Nil(t, tree.root.parent)
	checkNode(tree.root)
}

// Returns a graphviz DOT representation of tree. (https://graphviz.org/doc/info/lang.html)
func treeToString[K any, V any](t *testing.T, tree *btree[K, V]) string {
	var sb strings.Builder

	var logNode func(x *node[K, V])
	logNode = func(x *node[K, V]) {
		fmt.Fprintf(&sb, "\tnode%p [label=\"{%p|{", x, x)
		for i := 0; i < int(x.n); i++ {
			fmt.Fprintf(&sb, "<c%d> |%#v: %#v|", i, x.keys[i], x.values[i])
		}
		fmt.Fprintf(&sb, "<c%d> ", x.n)
		sb.WriteString("}}\"];\n")

		for i, child := range x.children {
			if child == nil {
				continue
			}
			fmt.Fprintf(&sb, "\tnode%p:c%d -> node%p;\n", x, i, child)
		}

		for _, child := range x.children {
			if child == nil {
				continue
			}
			logNode(child)
		}
	}

	sb.WriteString("digraph btree {\n\tnode [shape=record];\n")
	logNode(tree.root)
	sb.WriteString("}")

	return sb.String()
}

func TestAmalgam1(t *testing.T) {
	keys := [maxKVs]byte{}
	values := [maxKVs]byte{}
	for i := range keys {
		keys[i] = byte((i + 1) * 2)
		values[i] = byte((i + 1) * 4)
	}
	children := [branchFactor]*node[byte, byte]{}
	for i := range children {
		children[i] = &node[byte, byte]{}
	}
	extraChild := &node[byte, byte]{}

	check := func(
		extraKey byte,
		extraValue byte,
	) {
		t.Run(fmt.Sprintf("extraKey=%d,extraValue=%d", extraKey, extraValue), func(t *testing.T) {
			var expectedKeys [maxKVs + 1]byte
			copy(expectedKeys[:], keys[:])
			var expectedValues [maxKVs + 1]byte
			copy(expectedValues[:], values[:])
			var expectedChildren [branchFactor + 1]*node[byte, byte]
			copy(expectedChildren[:], children[:])
			idx := xsort.Search(expectedKeys[:len(expectedKeys)-1], xsort.OrderedLess[byte], extraKey)
			slices.Insert(expectedKeys[:len(expectedKeys)-1], idx, extraKey)
			slices.Insert(expectedValues[:len(expectedKeys)-1], idx, extraValue)
			slices.Insert(expectedChildren[:len(expectedChildren)-1], idx+1, extraChild)
			require2.Truef(
				t,
				xsort.SliceIsSorted(expectedKeys[:], xsort.OrderedLess[byte]),
				"%#v",
				expectedKeys,
			)

			a := newAmalgam1(
				xsort.OrderedLess[byte],
				&keys,
				&values,
				&children,
				extraKey,
				extraValue,
				extraChild,
			)

			t.Logf("extraIdx=%d", a.extraIdx)

			var actualKeys [maxKVs + 1]byte
			var actualValues [maxKVs + 1]byte
			var actualChildren [branchFactor + 1]*node[byte, byte]
			for i := 0; i < len(actualKeys); i++ {
				actualKeys[i] = a.Key(i)
				actualValues[i] = a.Value(i)
				actualChildren[i] = a.Child(i)
			}
			actualChildren[len(actualChildren)-1] = a.Child(len(actualChildren) - 1)

			require2.Equal(t, expectedKeys, actualKeys)
			require2.Equal(t, expectedValues, actualValues)
			require2.Equal(t, expectedChildren, actualChildren)
		})
	}

	for i := 0; i < maxKVs+1; i++ {
		check(byte(i*2+1), byte(i*4+1))
	}
}
