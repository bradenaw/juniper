//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/slices"
	"github.com/bradenaw/juniper/xsort"
)

// Maximum number of children each node can have.
const branchFactor = 16

// Maximum number of key/value pairs each node can have.
const maxKVs = branchFactor - 1

// < minKVs means we need to merge with a neighboring sibling
//             ┌───────────────────╴size of underfilled node
//             │          ┌────────╴size of sibling (any larger and t.steal() would work instead)
//             │          │      ┌─╴separator between the two in parent
//       ┌─────┴────┐   ┌─┴──┐  ┌┴┐
// thus, (minKVs - 1) + minKVs + 1   <= maxKVs
//                                   └───┬───┘
//                                       └──╴(any larger and we wouldn't be able to fit everything
//                                           into a single node)
// thus 2*minKVs <= maxKVs, so round-down is appropriate here
//
// Does not apply to the root.
const minKVs = maxKVs / 2

// Invariants:
// 1. Every node except the root has n >= minKVs.
// 2. The root has n >= 1 if the tree is non-empty.
// 3. Every node has node.n+1 children or no children.
//    - Notably, most nodes are leaves so we can do better space-wise if we can elide the children
//      array from internal nodes entirely.
type btree[K, V any] struct {
	root *node[K, V]
	less xsort.Less[K]
	size int
	// incremented when tree structure changes - used to quickly avoid reseeking cursor moving
	// through an unchanging tree
	gen int
}

func newBtree[K any, V any](less xsort.Less[K]) *btree[K, V] {
	return &btree[K, V]{
		less: less,
		root: &node[K, V]{},
		size: 0,
	}
}

// keys                 0           1           2                      n-1                        //
// values               0           1           2                      n-1                        //
// children       0           1           2          ...         n-1          n                   //
//               └┬┘         └┬┘                                             └┬┘                  //
//                │           └─╴contains keys greater than keys[0] and less  │                   //
//                │              than keys[1]                                 │                   //
//                │                                                           │                   //
//                └─────────────╴contains keys less than keys[0]              │                   //
//                                                                            │                   //
//                               contains keys greater than keys[n-1]╶────────┘                   //
type node[K any, V any] struct {
	parent   *node[K, V]
	children [branchFactor]*node[K, V]
	keys     [maxKVs]K
	values   [maxKVs]V
	// number of k/v pairs, naturally [1, maxKVs]
	n int8
}

func (x *node[K, V]) leaf() bool {
	return x.children[0] == nil
}

func (x *node[K, V]) full() bool {
	return int(x.n) == len(x.keys)
}

func (t *btree[K, V]) Len() int {
	return t.size
}

func (t *btree[K, V]) Put(k K, v V) {
	curr := t.root
	for {
		idx, inNode := t.searchNode(k, curr)
		if inNode {
			curr.values[idx] = v
			return
		}
		if curr.leaf() {
			break
		}
		curr = curr.children[idx]
	}
	if !curr.full() {
		t.insertIntoLeaf(curr, k, v)
	} else {
		t.overfill(curr, k, v, nil)
	}
	t.gen++
	t.size++
}

func (t *btree[K, V]) Get(k K) V {
	curr := t.root
	for curr != nil {
		idx, inNode := t.searchNode(k, curr)
		if inNode {
			return curr.values[idx]
		}
		curr = curr.children[idx]
	}
	var zero V
	return zero
}

func (t *btree[K, V]) Contains(k K) bool {
	curr := t.root
	for curr != nil {
		idx, inNode := t.searchNode(k, curr)
		if inNode {
			return true
		}
		curr = curr.children[idx]
	}
	return false
}

func (t *btree[K, V]) Delete(k K) {
	curr := t.root
	var idx int
	for {
		var inNode bool
		idx, inNode = t.searchNode(k, curr)
		if inNode {
			break
		}
		if curr.leaf() {
			// already at a leaf and !inNode, so k isn't in the tree
			return
		}
		curr = curr.children[idx]
	}

	t.size--
	t.gen++

	var leaf *node[K, V]
	if curr.leaf() {
		removeOne(curr.keys[:int(curr.n)], idx)
		removeOne(curr.values[:int(curr.n)], idx)
		curr.n--
		if curr.n >= minKVs || t.steal(curr) {
			return
		}
		leaf = curr
	} else {
		var replacementK K
		var replacementV V
		replacementK, replacementV, leaf = t.removeRightmost(curr.children[idx])
		curr.keys[idx] = replacementK
		curr.values[idx] = replacementV
		if leaf == nil || t.steal(leaf) {
			return
		}
	}

	if leaf != t.root {
		t.merge(leaf)
	}
}

func (t *btree[K, V]) First() (K, V) {
	if t.root.n == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV
	}
	leaf := leftmostLeaf(t.root)
	return leaf.keys[0], leaf.values[0]
}

func (t *btree[K, V]) Last() (K, V) {
	if t.root.n == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV
	}
	leaf := rightmostLeaf(t.root)
	return leaf.keys[int(leaf.n)-1], leaf.values[int(leaf.n)-1]
}

func (t *btree[K, V]) Cursor() cursor[K, V] {
	c := cursor[K, V]{t: t}
	c.SeekFirst()
	return c
}

// inserts k,v into the non-full leaf x.
func (t *btree[K, V]) insertIntoLeaf(x *node[K, V], k K, v V) {
	idx := 0
	for idx < int(x.n) {
		if t.less(k, x.keys[idx]) {
			break
		}
		idx++
	}
	insertOne(x.keys[:int(x.n)+1], idx, k)
	insertOne(x.values[:int(x.n)+1], idx, v)
	x.n++
}

// overfill adds k/v and k's right child afterK to an already-full x by splitting x into two. This
// adds a separater to x's parent, which may cause it to overflow and also need a split.
func (t *btree[K, V]) overfill(x *node[K, V], k K, v V, afterK *node[K, V]) {
	for {
		all := newAmalgam1(t.less, &x.keys, &x.values, &x.children, k, v, afterK)

		left := x
		right := &node[K, V]{}
		leaf := x.leaf()

		medianIdx := all.Len() / 2
		sepKey := all.Key(medianIdx)
		sepValue := all.Value(medianIdx)

		right.n = int8(all.Len() - medianIdx - 1)
		for i := 0; i < int(right.n); i++ {
			right.keys[i] = all.Key(medianIdx + 1 + i)
			right.values[i] = all.Value(medianIdx + 1 + i)
		}
		if !leaf {
			for i := 0; i < int(right.n)+1; i++ {
				right.children[i] = all.Child(medianIdx + 1 + i)
				right.children[i].parent = right
			}
		}

		left.n = int8(medianIdx)
		for i := int(left.n) - 1; i >= 0; i-- {
			left.keys[i] = all.Key(i)
			left.values[i] = all.Value(i)
		}
		if !leaf {
			for i := int(left.n); i >= 0; i-- {
				left.children[i] = all.Child(i)
				left.children[i].parent = left
			}
		}

		slices.Clear(left.keys[int(left.n):])
		slices.Clear(left.values[int(left.n):])
		slices.Clear(left.children[int(left.n)+1:])

		if x == t.root {
			parent := &node[K, V]{}
			parent.keys[0], parent.values[0] = sepKey, sepValue
			parent.n = 1
			parent.children[0] = left
			left.parent = parent
			parent.children[1] = right
			right.parent = parent
			t.root = parent
			return
		}

		parent := left.parent
		if !parent.full() {
			idxInParent := slices.Index(parent.children[:], left)
			insertOne(parent.keys[:int(parent.n)+1], idxInParent, sepKey)
			insertOne(parent.values[:int(parent.n)+1], idxInParent, sepValue)
			insertOne(parent.children[:int(parent.n)+2], idxInParent+1, right)
			right.parent = parent
			parent.n++
			return
		}
		x = parent
		k = sepKey
		v = sepValue
		afterK = right
	}
}

// merge merges x with one of its siblings.
//
// Assumes that it has a sibling that has n<=minKVs.
func (t *btree[K, V]) merge(x *node[K, V]) {
	left, right := t.siblings(x)
	if left != nil && left.n <= minKVs {
		t.mergeTwo(left, x)
	} else { // implies right != nil && right.n <= minKVs
		t.mergeTwo(x, right)
	}
}

// mergeTwo merges left and right together. This removes a node from the parent which may cause it
// to be underfilled as well, and will be fixed by stealing or merging.
//
// Assumes either left or right has n < minKVs and the other has n == minKVs.
//
//                       parent                                      parent                       //
//                     ┌───────────────┐                           ┌───────────────┐              //
//                     │   a   g   l   │                           │   a   l       │              //
//                     └╴•╶─╴•╶─╴•╶─╴•─┘                           └╴•╶─╴•╶─╴•╶────┘              //
//             left   ┌──────┘   └───────┐  right                 left   │                        //
//            ┌───────┴───────┐  ┌───────┴───────┐               ┌───────┴───────┐                //
//            │   c           │  │   h           │    ╶────>     │   c   g   h   │                //
//            └╴•╶─╴•╶─╴•╶─╴•╶┘  └╴•╶─╴•╶────────┘               └╴•╶─╴•╶─╴•╶─╴•╶┘                //
func (t *btree[K, V]) mergeTwo(left, right *node[K, V]) {
	parent := left.parent
	idxInParent := slices.Index(parent.children[:], left)
	sepKey := parent.keys[idxInParent]
	sepValue := parent.values[idxInParent]

	left.keys[int(left.n)] = sepKey
	copy(left.keys[int(left.n)+1:], right.keys[:int(right.n)])
	left.values[int(left.n)] = sepValue
	copy(left.values[int(left.n)+1:], right.values[:int(right.n)])
	copy(left.children[int(left.n)+1:], right.children[:int(right.n)+1])
	if !right.leaf() {
		for i := 0; i < int(right.n)+1; i++ {
			right.children[i].parent = left
		}
	}
	left.n += right.n + 1

	removeOne(parent.keys[:int(parent.n)], idxInParent)
	removeOne(parent.values[:int(parent.n)], idxInParent)
	removeOne(parent.children[:int(parent.n)+1], idxInParent+1)
	parent.n--

	// signal to cursors in right that they're lost.
	right.n = 0

	if parent == t.root {
		if parent.n == 0 {
			t.root = left
			left.parent = nil
		}
	} else if parent.n < minKVs && !t.steal(parent) {
		t.merge(parent)
	}
}

// removeRightmost finds the rightmost key and value in the subtree rooted by x and removes them.
// These are by definition in a leaf. If this caused the leaf to be underfilled, also returns the
// leaf they were removed from.
func (t *btree[K, V]) removeRightmost(x *node[K, V]) (K, V, *node[K, V]) {
	curr := rightmostLeaf(x)
	k := curr.keys[int(curr.n)-1]
	v := curr.values[int(curr.n)-1]
	var zeroK K
	curr.keys[int(curr.n)-1] = zeroK
	var zeroV V
	curr.values[int(curr.n)-1] = zeroV
	curr.n--
	var out *node[K, V]
	if curr.n < minKVs {
		out = curr
	}
	return k, v, out
}

// steal adds one k/v/child to x by taking from one of its siblings if possible. If not, returns
// false.
func (t *btree[K, V]) steal(x *node[K, V]) bool {
	left, right := t.siblings(x)
	if right != nil && right.n > minKVs {
		t.rotateLeft(x, right)
		return true
	}
	if left != nil && left.n > minKVs {
		t.rotateRight(left, x)
		return true
	}
	return false
}

// siblings returns x's immediate left and right siblings, or nil if none exists.
func (t *btree[K, V]) siblings(x *node[K, V]) (*node[K, V], *node[K, V]) {
	if x.parent == nil {
		return nil, nil
	}
	idx := slices.Index(x.parent.children[:], x)
	var left, right *node[K, V]
	if idx > 0 {
		left = x.parent.children[idx-1]
	}
	if idx < int(x.parent.n) {
		right = x.parent.children[idx+1]
	}
	return left, right
}

//                  parent                                            parent                      //
//                ┌───────────────┐                                 ┌───────────────┐             //
//                │  [g]          │                                 │  [c]          │             //
//                └╴•╶─╴•╶────────┘                                 └╴•╶─╴•╶────────┘             //
//   left    ┌──────┘   └───────┐  right               left    ┌──────┘   └───────┐  right        //
//   ┌───────┴───────┐  ┌───────┴───────┐              ┌───────┴───────┐  ┌───────┴───────┐       //
//   │   a   b  [c]  │  │   h   i       │    ╶────>    │   a   b   [ ] │  │  [g]  h   i   │       //
//   └╴•╶─╴•╶─╴•╶─[•]┘  └╴•╶─╴•╶─╴•╶────┘              └╴•╶─╴•╶─╴•╶────┘  └[•]─╴•╶─╴•╶─╴•╶┘       //
//                 │                                                        │                     //
//                 │   child                                                │   child             //
//         ┌───────┴───────┐                                        ┌───────┴───────┐             //
//         │   d   e   f   │                                        │   d   e   f   │             //
//         └╴•╶─╴•╶─╴•╶─╴•╶┘                                        └╴•╶─╴•╶─╴•╶─╴•╶┘             //
// (Changes marked with [])
//
// Assumes left and right are siblings and right is not full.
func (t *btree[K, V]) rotateRight(left *node[K, V], right *node[K, V]) {
	idxInParent := slices.Index(left.parent.children[:], left)
	oldSepK := left.parent.keys[idxInParent]
	oldSepV := left.parent.values[idxInParent]
	child := left.children[left.n]

	// copy the max key from left up to the separator
	left.parent.keys[idxInParent] = left.keys[left.n-1]
	left.parent.values[idxInParent] = left.values[left.n-1]

	// remove the max key/child from left
	var zeroK K
	left.keys[left.n-1] = zeroK
	var zeroV V
	left.values[left.n-1] = zeroV
	left.children[left.n] = nil
	left.n--

	// move the old separator to the minimum key in right
	insertOne(right.keys[:], 0, oldSepK)
	insertOne(right.values[:], 0, oldSepV)
	insertOne(right.children[:], 0, child)
	if child != nil {
		child.parent = right
	}
	right.n++
}

//                  parent                                                parent                  //
//                ┌───────────────┐                                     ┌───────────────┐         //
//                │  [c]          │                                     │  [g]          │         //
//                └╴•╶─╴•╶────────┘                                     └╴•╶─╴•╶────────┘         //
//   left    ┌──────┘   └───────┐  right                   left    ┌──────┘   └───────┐  right    //
//   ┌───────┴───────┐  ┌───────┴───────┐                  ┌───────┴───────┐  ┌───────┴───────┐   //
//   │   a   b   [ ] │  │  [g]  h   i   │      ╶────>      │   a   b  [c]  │  │ [ ]  h   i    │   //
//   └╴•╶─╴•╶─╴•╶────┘  └[•]─╴•╶─╴•╶─╴•╶┘                  └╴•╶─╴•╶─╴•╶─[•]┘  └───╴•╶─╴•╶─╴•╶─┘   //
//                        │                                              │                        //
//                        │   child                                      │   child                //
//                ┌───────┴───────┐                              ┌───────┴───────┐                //
//                │   d   e   f   │                              │   d   e   f   │                //
//                └╴•╶─╴•╶─╴•╶─╴•╶┘                              └╴•╶─╴•╶─╴•╶─╴•╶┘                //
// (Changes marked with [])
//
// Assumes left and right are siblings and left is not full.
func (t *btree[K, V]) rotateLeft(left *node[K, V], right *node[K, V]) {
	idxInParent := slices.Index(right.parent.children[:], right)
	oldSepK := right.parent.keys[idxInParent-1]
	oldSepV := right.parent.values[idxInParent-1]
	child := right.children[0]

	// copy the minimum key in right up to the separator
	right.parent.keys[idxInParent-1] = right.keys[0]
	right.parent.values[idxInParent-1] = right.values[0]

	// remove right's minimum key
	removeOne(right.keys[:], 0)
	removeOne(right.values[:], 0)
	removeOne(right.children[:], 0)
	right.n--

	// move the old separator to the maximum key in left
	left.keys[left.n] = oldSepK
	left.values[left.n] = oldSepV
	left.children[left.n+1] = child
	if child != nil {
		child.parent = left
	}
	left.n++
}

// If inNode is true, idx is the index in x.keys that k is at. If false, idx is the index of the
// child to look in.
func (t *btree[K, V]) searchNode(k K, x *node[K, V]) (idx int, inNode bool) {
	// benchmark suggests that linear search is in fact faster than binary search, at least for int
	// keys and branchFactor <= 32.
	for i := 0; i < int(x.n); i++ {
		if t.less(k, x.keys[i]) {
			return i, false
		} else if !t.less(x.keys[i], k) {
			return i, true
		}
	}
	return int(x.n), false
}

func leftmostLeaf[K any, V any](x *node[K, V]) *node[K, V] {
	curr := x
	for {
		if curr.leaf() {
			return curr
		}
		curr = curr.children[0]
	}
}

func rightmostLeaf[K any, V any](x *node[K, V]) *node[K, V] {
	curr := x
	for {
		if curr.leaf() {
			return curr
		}
		curr = curr.children[int(curr.n)]
	}
}

func removeOne[T any](a []T, idx int) {
	copy(a[idx:], a[idx+1:])
	var zero T
	a[len(a)-1] = zero
}

// Inserts x into a at index idx, shifting the rest of the elements over. Clobbers a[len(a)-1].
//
// Faster in this package than slices.Insert for use on node.{keys,values,children} since they never
// grow.
func insertOne[T any](a []T, idx int, x T) {
	copy(a[idx+1:], a[idx:])
	a[idx] = x
}

type amalgam1[K any, V any] struct {
	keys       *[maxKVs]K
	values     *[maxKVs]V
	children   *[branchFactor]*node[K, V]
	extraKey   K
	extraValue V
	extraChild *node[K, V]
	extraIdx   int
}

// Returns a cheap view that functions like a slice of all of the inputs. Assumes that both arrays
// are in sorted order, and all leftKeys are less than sepKey and sepKey is less than all rightKeys.
// extraKey may belong in any position.
//
// Example:
//      keys           extraKey
//   [a   c   e]     +    d
//  0   1   2   3           extraChild
//
//               amalgam
//           [a   c   d            e]
//          0   1   2   extraChild   3
func newAmalgam1[K any, V any](
	less xsort.Less[K],
	keys *[maxKVs]K,
	values *[maxKVs]V,
	children *[branchFactor]*node[K, V],
	extraKey K,
	extraValue V,
	extraChild *node[K, V],
) amalgam1[K, V] {
	extraIdx := func() int {
		for i := range *keys {
			if less(extraKey, keys[i]) {
				return i
			}
		}
		return len(keys)
	}()

	return amalgam1[K, V]{
		keys:       keys,
		values:     values,
		children:   children,
		extraKey:   extraKey,
		extraValue: extraValue,
		extraChild: extraChild,
		extraIdx:   extraIdx,
	}
}

func (a *amalgam1[K, V]) Len() int {
	return maxKVs + 1
}

func (a *amalgam1[K, V]) Key(i int) K {
	if i == a.extraIdx {
		return a.extraKey
	} else if i > a.extraIdx {
		i--
	}
	return a.keys[i]
}
func (a *amalgam1[K, V]) Value(i int) V {
	if i == a.extraIdx {
		return a.extraValue
	} else if i > a.extraIdx {
		i--
	}
	return a.values[i]
}
func (a *amalgam1[K, V]) Child(i int) *node[K, V] {
	if i == a.extraIdx+1 {
		return a.extraChild
	} else if i > a.extraIdx+1 {
		i--
	}
	return a.children[i]
}

type cursor[K any, V any] struct {
	t *btree[K, V]
	// Set to nil when run off the edge.
	curr *node[K, V]
	// Index of k in curr. Used to notice when k has been moved or deleted.
	i int
	// last seen gen of tree
	gen int
	k   K
}

func (c *cursor[K, V]) Next() {
	if c.lost() {
		c.SeekFirstGreater(c.k)
		return
	}
	if c.curr == nil {
		return
	}

	if c.curr.leaf() {
		c.i++
		if c.i < int(c.curr.n) {
			c.k = c.curr.keys[c.i]
			return
		}
	} else if c.i < int(c.curr.n) {
		c.curr = leftmostLeaf(c.curr.children[c.i+1])
		c.i = 0
		c.k = c.curr.keys[c.i]
		return
	}

	for {
		if c.curr.parent == nil {
			c.curr = nil
			return
		}
		idx := slices.Index(c.curr.parent.children[:], c.curr)
		c.curr = c.curr.parent
		c.i = idx
		if c.i < int(c.curr.n) {
			c.k = c.curr.keys[c.i]
			break
		}
	}
}

func (c *cursor[K, V]) Prev() {
	if c.lost() {
		c.SeekLastLess(c.k)
		return
	}
	if c.curr == nil {
		return
	}

	if c.curr.leaf() {
		c.i--
		if c.i >= 0 {
			c.k = c.curr.keys[c.i]
			return
		}
	} else if c.i >= 0 {
		c.curr = rightmostLeaf(c.curr.children[c.i])
		c.i = int(c.curr.n) - 1
		c.k = c.curr.keys[c.i]
		return
	}

	for {
		if c.curr.parent == nil {
			c.curr = nil
			return
		}
		idx := slices.Index(c.curr.parent.children[:], c.curr)
		c.curr = c.curr.parent
		c.i = idx - 1
		if c.i >= 0 {
			c.k = c.curr.keys[c.i]
			break
		}
	}
}

func (c *cursor[K, V]) Ok() bool {
	return c.curr != nil && c.refind()
}

func (c *cursor[K, V]) Key() K {
	return c.k
}

func (c *cursor[K, V]) Value() V {
	var zero V
	if !c.refind() {
		return zero
	}
	return c.curr.values[c.i]
}

func (c *cursor[K, V]) valueUnchecked() V {
	return c.curr.values[c.i]
}

func (c *cursor[K, V]) SeekFirst() {
	if c.t.root.n == 0 {
		c.curr = nil
		return
	}
	c.curr = leftmostLeaf(c.t.root)
	c.i = 0
	c.k = c.curr.keys[c.i]
	c.gen = c.t.gen
}

func (c *cursor[K, V]) SeekLastLess(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.LessOrEqual(c.t.less, k, c.k) {
		c.Prev()
	}
}

func (c *cursor[K, V]) SeekLastLessOrEqual(k K) {
	if !c.seek(k) {
		return
	}
	if c.t.less(k, c.k) {
		c.Prev()
	}
}

func (c *cursor[K, V]) SeekFirstGreaterOrEqual(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.Greater(c.t.less, k, c.k) {
		c.Next()
	}
}

func (c *cursor[K, V]) SeekFirstGreater(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.GreaterOrEqual(c.t.less, k, c.k) {
		c.Next()
	}
}

func (c *cursor[K, V]) SeekLast() {
	if c.t.root.n == 0 {
		c.curr = nil
		return
	}
	c.curr = rightmostLeaf(c.t.root)
	c.i = int(c.curr.n) - 1
	c.k = c.curr.keys[c.i]
	c.gen = c.t.gen
}

// seek moves the cursor to k or its successor or predecessor if it isn't in the tree. Returns false
// if the cursor is now invalid because the tree is empty.
func (c *cursor[K, V]) seek(k K) bool {
	c.curr, c.i, _ = c.find(k)
	if c.curr == nil {
		return false
	}
	c.k = c.curr.keys[c.i]
	c.gen = c.t.gen
	return true
}

// find looks for k in the tree. It returns the node that k appears in and the index it appears at
// and true in the final return. If k is not in the tree, then the final return is false and the
// returned node and index of a successor or predecessor of k.
//
// The returned node is nil if the tree is empty.
func (c *cursor[K, V]) find(k K) (*node[K, V], int, bool) {
	if c.t.root.n == 0 {
		return nil, 0, false
	}
	curr := c.t.root
	for {
		idx, inNode := c.t.searchNode(k, curr)
		if inNode {
			return curr, idx, true
		}
		if curr.leaf() {
			if idx == int(curr.n) {
				idx--
			}
			return curr, idx, false
		}
		curr = curr.children[idx]
	}
}

// refind ensures c.curr[c.i] == c.k if c.k is still in the tree (which could've been made false if
// the tree was modified since the cursor found its position) by reseeking. Returns false without
// modifying the cursor if c.k isn't in the tree anymore.
func (c *cursor[K, V]) refind() bool {
	if !c.lost() {
		return true
	}
	curr, i, ok := c.find(c.k)
	if !ok {
		return false
	}
	c.curr = curr
	c.i = i
	c.gen = c.t.gen
	return true
}

// lost returns true if the tree has been modified in such a way that the cursor has lost its place.
func (c *cursor[K, V]) lost() bool {
	// c.curr == nil implies the cursor is already off the edge of the tree and cannot be lost.
	//
	// Otherwise, check that the element of c.curr we're pointed at still contains the key we
	// expect, since it might've gotten shifted from e.g. deleting the element before this one.
	// Careful: x.keys[i] for i >= x.n is filled with the zero value, and c.k might happen to be the
	// zero value also. Unlinking a node during merge sets n=0, so that's handled here too.
	return c.gen != c.t.gen &&
		c.curr != nil &&
		(c.i >= int(c.curr.n) || !xsort.Equal(c.t.less, c.k, c.curr.keys[c.i]))
}

func (c *cursor[K, V]) Forward() iterator.Iterator[KVPair[K, V]] {
	return &forwardIterator[K, V]{c: *c}
}

type forwardIterator[K any, V any] struct {
	c cursor[K, V]
}

func (iter *forwardIterator[K, V]) Next() (KVPair[K, V], bool) {
	if iter.c.lost() {
		iter.c.SeekFirstGreaterOrEqual(iter.c.Key())
	}
	if iter.c.curr == nil {
		var zero KVPair[K, V]
		return zero, false
	}
	k := iter.c.Key()
	// Safe since we already made sure !iter.c.lost() by reseeking above.
	v := iter.c.valueUnchecked()
	iter.c.Next()
	return KVPair[K, V]{k, v}, true
}

func (c *cursor[K, V]) Backward() iterator.Iterator[KVPair[K, V]] {
	return &backwardIterator[K, V]{c: *c}
}

type backwardIterator[K any, V any] struct {
	c cursor[K, V]
}

func (iter *backwardIterator[K, V]) Next() (KVPair[K, V], bool) {
	if iter.c.lost() {
		iter.c.SeekLastLessOrEqual(iter.c.Key())
	}
	if iter.c.curr == nil {
		var zero KVPair[K, V]
		return zero, false
	}
	k := iter.c.Key()
	// Safe since we already made sure !iter.c.lost() by reseeking above.
	v := iter.c.valueUnchecked()
	iter.c.Prev()
	return KVPair[K, V]{k, v}, true
}
