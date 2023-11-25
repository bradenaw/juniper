package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/internal/orderedhashmap"
	"github.com/bradenaw/juniper/internal/require2"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xslices"
	"github.com/bradenaw/juniper/xsort"
)

func orderedhashmapKVPairToKVPair[K any, V any](kv orderedhashmap.KVPair[uint16, int]) KVPair[uint16, int] {
	return KVPair[uint16, int]{kv.K, kv.V}
}

func FuzzBtree(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		tree := newBtree[uint16, int](xsort.OrderedLess[uint16])
		cursor := tree.Cursor()
		cursor.SeekFirst()
		oracle := orderedhashmap.NewMap[uint16, int](xsort.OrderedLess[uint16])
		oracleCursor := oracle.Cursor()

		ctr := 0

		fuzz.Operations(
			b,
			func() { // check
				t.Log(treeToString(tree))

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
				c.SeekFirst()
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

func TestSplitRoot(t *testing.T) {
	if branchFactor != 16 {
		t.Skip("test requires branchFactor 16")
	}

	tree := makeTree(t, makeLeaf([]KVPair[byte, int]{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
		{5, 5},
		{6, 6},
		{7, 7},
		{8, 8},
		{9, 9},
		{10, 10},
		{11, 11},
		{12, 12},
		{13, 13},
		{14, 14},
	}))
	require2.Equal(t, 1, numNodes(tree))
	tree.Put(15, 15)
	require2.Equal(t, 3, numNodes(tree))
	requireTreesEqual(
		t,
		tree,
		makeTree(t, makeInternal(
			makeLeaf([]KVPair[byte, int]{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
				{5, 5},
				{6, 6},
				{7, 7},
			}),
			KVPair[byte, int]{8, 8},
			makeLeaf([]KVPair[byte, int]{
				{9, 9},
				{10, 10},
				{11, 11},
				{12, 12},
				{13, 13},
				{14, 14},
				{15, 15},
			}),
		)),
	)
}

func TestMerge(t *testing.T) {
	if branchFactor != 16 {
		t.Skip("test requires branchFactor 16")
	}

	tree := makeTree(t, makeInternal(
		makeLeaf([]KVPair[byte, int]{
			{0, 0},
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
			{5, 5},
			{6, 6},
		}),
		KVPair[byte, int]{10, 10},
		makeLeaf([]KVPair[byte, int]{
			{20, 20},
			{21, 21},
			{22, 22},
			{23, 23},
			{24, 24},
			{25, 25},
			{26, 26},
		}),
	))

	require2.Equal(t, tree.root.children[0].n, minKVs)
	require2.Equal(t, tree.root.children[1].n, minKVs)
	tree.Delete(23)
	checkTree(t, tree)

	requireTreesEqual(
		t,
		tree,
		makeTree(t, makeLeaf([]KVPair[byte, int]{
			{0, 0},
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
			{5, 5},
			{6, 6},
			{10, 10},
			{20, 20},
			{21, 21},
			{22, 22},
			{24, 24},
			{25, 25},
			{26, 26},
		})),
	)
}

func TestRotateRight(t *testing.T) {
	if branchFactor != 16 {
		t.Skip("test requires branchFactor 16")
	}
	tree := makeTree(t, makeInternal(
		makeLeaf([]KVPair[byte, int]{
			{0, 0},
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
			{5, 5},
			{6, 6},
			{7, 7},
		}),
		KVPair[byte, int]{10, 10},
		makeLeaf([]KVPair[byte, int]{
			{20, 20},
			{21, 21},
			{22, 22},
			{23, 23},
			{24, 24},
			{25, 25},
			{26, 26},
		}),
	))

	require2.Equal(t, tree.root.children[0].n, minKVs+1)
	require2.Equal(t, tree.root.children[1].n, minKVs)

	tree.Delete(20)
	checkTree(t, tree)

	requireTreesEqual(
		t,
		tree,
		makeTree(t, makeInternal(
			makeLeaf([]KVPair[byte, int]{
				{0, 0},
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
				{5, 5},
				{6, 6},
			}),
			KVPair[byte, int]{7, 7},
			makeLeaf([]KVPair[byte, int]{
				{10, 10},
				{21, 21},
				{22, 22},
				{23, 23},
				{24, 24},
				{25, 25},
				{26, 26},
				{27, 27},
			}),
		)),
	)
}

func TestRotateLeft(t *testing.T) {
	if branchFactor != 16 {
		t.Skip("test requires branchFactor 16")
	}
	tree := makeTree(t, makeInternal(
		makeLeaf([]KVPair[byte, int]{
			{0, 0},
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
			{5, 5},
			{6, 6},
		}),
		KVPair[byte, int]{10, 10},
		makeLeaf([]KVPair[byte, int]{
			{20, 20},
			{21, 21},
			{22, 22},
			{23, 23},
			{24, 24},
			{25, 25},
			{26, 26},
			{27, 27},
		}),
	))

	require2.Equal(t, tree.root.children[0].n, minKVs)
	require2.Equal(t, tree.root.children[1].n, minKVs+1)

	tree.Delete(0)
	checkTree(t, tree)

	requireTreesEqual(
		t,
		tree,
		makeTree(t, makeInternal(
			makeLeaf([]KVPair[byte, int]{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
				{5, 5},
				{6, 6},
				{10, 10},
			}),
			KVPair[byte, int]{20, 20},
			makeLeaf([]KVPair[byte, int]{
				{21, 21},
				{22, 22},
				{23, 23},
				{24, 24},
				{25, 25},
				{26, 26},
				{27, 27},
			}),
		)),
	)
}

func TestMergeMulti(t *testing.T) {
	tree := newBtree[uint16, int](xsort.OrderedLess[uint16])
	i := 0
	for treeHeight(tree) < 3 {
		tree.Put(uint16(i), i)
		i++
	}
	for {
		had := false
		breadthFirst(tree, func(x *node[uint16, int]) bool {
			t.Logf("visit(%p)", x)
			if x.n > minKVs {
				tree.Delete(x.keys[0])
				had = true
				return false
			}
			return true
		})
		require2.Equal(t, treeHeight(tree), 3)
		checkTree(t, tree)
		if !had {
			break
		}
	}

	c := tree.Cursor()
	c.SeekFirst()
	expected := iterator.Collect(c.Forward())
	nNodesBefore := numNodes(tree)
	var removed uint16

	t.Logf(treeToString(tree))
	breadthFirst(tree, func(x *node[uint16, int]) bool {
		if x.leaf() {
			removed = x.keys[0]
			tree.Delete(x.keys[0])
			return false
		}
		return true
	})
	checkTree(t, tree)
	t.Logf("removed %#v", removed)
	t.Logf(treeToString(tree))

	expected = xslices.Remove(
		expected,
		xslices.IndexFunc(expected, func(pair KVPair[uint16, int]) bool {
			return pair.Key == removed
		}),
		1,
	)
	nNodesAfter := numNodes(tree)
	c = tree.Cursor()
	c.SeekFirst()
	actual := iterator.Collect(c.Forward())

	require2.SlicesEqual(t, expected, actual)
	require2.Equal(t, nNodesBefore-3, nNodesAfter)
}

func breadthFirst[K any, V any](tree *btree[K, V], visit func(*node[K, V]) bool) {
	queue := []*node[K, V]{tree.root}
	for len(queue) > 0 {
		var curr *node[K, V]
		curr, queue = queue[0], queue[1:]
		if !visit(curr) {
			return
		}
		if !curr.leaf() {
			for i := 0; i <= int(curr.n); i++ {
				queue = append(queue, curr.children[i])
			}
		}
	}
}

func requireTreesEqual(t *testing.T, a, b *btree[byte, int]) {
	eq := func() bool {
		var visit func(x, y *node[byte, int]) bool
		visit = func(x, y *node[byte, int]) bool {
			if (x == nil) != (y == nil) {
				return false
			}
			if x == nil && y == nil {
				return true
			}
			if x.n != y.n {
				return false
			}
			if x.leaf() != y.leaf() {
				return false
			}
			for i := 0; i < int(x.n); i++ {
				if x.keys[i] != y.keys[i] {
					return false
				}
				if x.values[i] != y.values[i] {
					return false
				}
			}
			if x.leaf() {
				for i := 0; i < int(x.n)+1; i++ {
					if !visit(x.children[i], y.children[i]) {
						return false
					}
				}
			}
			return true
		}
		return visit(a.root, b.root)
	}()
	if !eq {
		t.Fatalf("%s\n\n%s", treeToStringNoPtr(a), treeToStringNoPtr(b))
	}
}

func makeTree(t *testing.T, root *node[byte, int]) *btree[byte, int] {
	tree := &btree[byte, int]{
		root: root,
		less: xsort.OrderedLess[byte],
	}
	tree.size = numItems(tree)
	checkTree(t, tree)
	return tree
}

func makeInternal(items ...any) *node[byte, int] {
	x := &node[byte, int]{n: int8(len(items) / 2)}
	for i := 0; i < int(x.n)+1; i++ {
		x.children[i] = items[i*2].(*node[byte, int])
		x.children[i].parent = x
	}
	for i := 0; i < int(x.n); i++ {
		pair := items[i*2+1].(KVPair[byte, int])
		x.keys[i] = pair.Key
		x.values[i] = pair.Value
	}
	return x
}

func makeLeaf(kvs []KVPair[byte, int]) *node[byte, int] {
	x := &node[byte, int]{n: int8(len(kvs))}
	for i := range kvs {
		x.keys[i] = kvs[i].Key
		x.values[i] = kvs[i].Value
	}
	return x
}

func pairsRange(min, max byte) []KVPair[byte, int] {
	var out []KVPair[byte, int]
	for i := min; i < max; i++ {
		out = append(out, KVPair[byte, int]{byte(i), int(i)})
	}
	return out
}

func isZero[T comparable](t T) bool {
	var zero T
	return t == zero
}

func numNodes[K any, V any](tree *btree[K, V]) int {
	n := 0
	var visit func(x *node[K, V])
	visit = func(x *node[K, V]) {
		n++
		if x.leaf() {
			return
		}
		for i := 0; i < int(x.n)+1; i++ {
			visit(x.children[i])
		}
	}
	visit(tree.root)
	return n
}

func numItems[K any, V any](tree *btree[K, V]) int {
	n := 0
	var visit func(x *node[K, V])
	visit = func(x *node[K, V]) {
		n += int(x.n)
		if x.leaf() {
			return
		}
		for i := 0; i < int(x.n)+1; i++ {
			visit(x.children[i])
		}
	}
	visit(tree.root)
	return n
}

func treeHeight[K any, V any](tree *btree[K, V]) int {
	curr := tree.root
	n := 0
	for curr != nil {
		n += 1
		curr = curr.children[0]
	}
	return n
}

func checkTree[K comparable, V comparable](t *testing.T, tree *btree[K, V]) {
	foundLeaf := false
	leafDepth := 0
	var checkNode func(x *node[K, V], depth int)
	checkNode = func(x *node[K, V], depth int) {
		if x.leaf() {
			for i := 0; i < int(x.n)+1; i++ {
				require2.Nil(t, x.children[i])
			}
			if !foundLeaf {
				leafDepth = depth
				foundLeaf = true
			}
			require2.Equal(t, leafDepth, depth)
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
				checkNode(x.children[i], depth+1)
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
		require2.True(t, xslices.All(x.keys[int(x.n):], isZero[K]))
		require2.True(t, xslices.All(x.values[int(x.n):], isZero[V]))
		require2.Truef(
			t,
			xslices.All(x.children[int(x.n)+1:], isZero[*node[K, V]]),
			"%p %#v",
			x,
			x.children[int(x.n)+1:],
		)
	}

	require2.NotNil(t, tree.root)
	require2.Nil(t, tree.root.parent)
	checkNode(tree.root, 0)
}

// Returns a graphviz DOT representation of tree. (https://graphviz.org/doc/info/lang.html)
func treeToString[K any, V any](tree *btree[K, V]) string {
	return treeToStringInner(tree, func(x *node[K, V]) string { return fmt.Sprintf("%p", x) })
}

func treeToStringNoPtr[K any, V any](tree *btree[K, V]) string {
	ids := make(map[*node[K, V]]string)
	ctr := 0
	return treeToStringInner(tree, func(x *node[K, V]) string {
		id, ok := ids[x]
		if !ok {
			id = fmt.Sprintf("%d", ctr)
			ctr++
			ids[x] = id
		}
		return id
	})
}

func treeToStringInner[K any, V any](tree *btree[K, V], id func(*node[K, V]) string) string {
	var sb strings.Builder

	var logNode func(x *node[K, V])
	logNode = func(x *node[K, V]) {
		fmt.Fprintf(&sb, "\tnode%s [label=\"{%s|{", id(x), id(x))
		for i := 0; i < int(x.n); i++ {
			fmt.Fprintf(&sb, "<c%d> |%#v: %#v|", i, x.keys[i], x.values[i])
		}
		fmt.Fprintf(&sb, "<c%d> ", x.n)
		sb.WriteString("}}\"];\n")

		for i, child := range x.children {
			if child == nil {
				continue
			}
			fmt.Fprintf(&sb, "\tnode%s:c%d -> node%s;\n", id(x), i, id(child))
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
			xslices.Insert(expectedKeys[:len(expectedKeys)-1], idx, extraKey)
			xslices.Insert(expectedValues[:len(expectedKeys)-1], idx, extraValue)
			xslices.Insert(expectedChildren[:len(expectedChildren)-1], idx+1, extraChild)
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

func TestRange(t *testing.T) {
	tree := newBtree[uint16, int](xsort.OrderedLess[uint16])

	for i := 0; i < 128; i++ {
		tree.Put(uint16(i), i)
	}

	keys := func(iter iterator.Iterator[KVPair[uint16, int]]) []uint16 {
		return iterator.Collect(iterator.Map(iter, func(pair KVPair[uint16, int]) uint16 {
			return pair.Key
		}))
	}
	check := func(lower Bound[uint16], upper Bound[uint16], expected []uint16) {
		require2.SlicesEqual(t, keys(tree.Range(lower, upper)), expected)
		r := keys(tree.RangeReverse(lower, upper))
		xslices.Reverse(r)
		require2.SlicesEqual(t, r, expected)
	}

	check(
		Included(uint16(5)), Included(uint16(16)),
		[]uint16{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
	)

	check(
		Unbounded[uint16](), Excluded(uint16(4)),
		[]uint16{0, 1, 2, 3},
	)

	check(
		Excluded(uint16(123)), Unbounded[uint16](),
		[]uint16{124, 125, 126, 127},
	)
}

func TestGetContains(t *testing.T) {
	tree := newBtree[uint16, int](xsort.OrderedLess[uint16])

	for i := 0; i < 128; i++ {
		tree.Put(uint16(i*2), i*4)
	}
	for i := 0; i < 128; i++ {
		key := uint16(i * 2)
		require2.True(t, tree.Contains(key))
		require2.Equal(t, tree.Get(key), int(key)*2)
	}
	for i := 0; i <= 128; i++ {
		key := uint16(i*2 - 1)
		require2.True(t, !tree.Contains(key))
		require2.Equal(t, tree.Get(key), 0)
	}
}

func TestDelete(t *testing.T) {
	tree := newBtree[uint16, int](xsort.OrderedLess[uint16])
	for i := 0; i < 128; i++ {
		tree.Put(uint16(i)+1, i*2)
	}
	require2.Equal(t, tree.Len(), 128)
	tree.Delete(0)
	tree.Delete(129)
	require2.Equal(t, tree.Len(), 128)

	for tree.Len() > 0 {
		key := uint16(0)
		l := tree.Len()
		if tree.Len()%2 == 0 {
			key, _ = tree.First()
		} else {
			key, _ = tree.Last()
		}
		require2.True(t, tree.Contains(key))
		tree.Delete(key)
		require2.True(t, !tree.Contains(key))
		require2.Equal(t, tree.Len(), l-1)
	}
}
