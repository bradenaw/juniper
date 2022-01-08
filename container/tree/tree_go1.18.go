//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xsort"
)

// tree is an AVL tree: https://en.wikipedia.org/wiki/AVL_tree
type tree[K any, V any, O xsort.Ordering[K]] struct {
	root     *node[K, V]
	size     int
	ordering O
	// Incremented whenever the tree structure changes, so that iterators know to reset.
	gen int
}

type node[K any, V any] struct {
	left   *node[K, V]
	right  *node[K, V]
	parent *node[K, V]
	// The height of this node. Leaves have height 0, internal nodes are one higher than their
	// highest child.
	height int
	// iterator depends on key not changing with respect to tree.less.
	key   K
	value V
}

func newTree[K any, V any, O xsort.Ordering[K]]() *tree[K, V, O] {
	return &tree[K, V, O]{
		root: nil,
		size: 0,
		gen:  1,
	}
}

func (t *tree[K, V, O]) Put(k K, v V) {
	if t.root == nil {
		t.root = &node[K, V]{
			key:   k,
			value: v,
		}
		t.size++
		t.gen++
		return
	}

	curr := t.root
	for {
		if t.ordering.Less(k, curr.key) {
			if curr.left == nil {
				curr.left = &node[K, V]{
					parent: curr,
					key:    k,
					value:  v,
					height: 0,
				}
				break
			}
			curr = curr.left
		} else if t.ordering.Less(curr.key, k) {
			if curr.right == nil {
				curr.right = &node[K, V]{
					parent: curr,
					key:    k,
					value:  v,
					height: 0,
				}
				break
			}
			curr = curr.right
		} else {
			curr.key = k
			curr.value = v
			return
		}
	}

	for {
		curr = t.rebalance(curr)
		if curr.parent == nil {
			break
		}
		curr = curr.parent
	}
	t.root = curr
	t.size++
	t.gen++
}

func (t *tree[K, V, O]) Delete(k K) {
	in := &t.root
	var last *node[K, V]
	for {
		if *in == nil {
			// k isn't in the map
			return
		}
		curr := *in
		if t.ordering.Less(k, curr.key) {
			in = &curr.left
		} else if t.ordering.Less(curr.key, k) {
			in = &curr.right
		} else {
			// curr contains k
			if curr.left != nil && curr.right != nil {
				// curr has both children, so replace it with its successor.
				successor := t.removeLeftmost(curr, &curr.right)
				successor.left = curr.left
				if successor.left != nil {
					successor.left.parent = successor
				}
				successor.right = curr.right
				if successor.right != nil {
					successor.right.parent = successor
				}
				successor.parent = curr.parent
				*in = successor
				last = successor
			} else if curr.left != nil {
				// curr has only a left child, just hoist it upwards
				curr.left.parent = curr.parent
				*in = curr.left
				last = curr.left
			} else if curr.right != nil {
				// curr has only a right child, just hoist it upwards
				curr.right.parent = curr.parent
				*in = curr.right
				last = curr.right
			} else {
				// curr has no children, remove it
				*in = nil
				last = curr.parent
			}
			break
		}
	}

	curr := last
	for curr != nil {
		curr = t.rebalance(curr)
		if curr.parent == nil {
			break
		}
		curr = curr.parent
	}
	t.root = curr
	t.size--
	t.gen++
}

func (t *tree[K, V, O]) removeLeftmost(
	parent *node[K, V],
	in **node[K, V],
) *node[K, V] {

	for {
		curr := *in
		if curr.left == nil {
			break
		}
		in = &curr.left
	}

	leftmost := *in
	x := leftmost.parent
	*in = leftmost.right
	if leftmost.right != nil {
		leftmost.right.parent = leftmost.parent
	}

	curr := x
	if curr == nil {
	}
	for curr != parent {
		t.setHeight(curr)
		curr = t.rebalance(curr)
		curr = curr.parent
	}

	return leftmost
}

func (t *tree[K, V, O]) Get(k K) V {
	curr := t.root
	for curr != nil {
		if t.ordering.Less(k, curr.key) {
			curr = curr.left
		} else if t.ordering.Less(curr.key, k) {
			curr = curr.right
		} else {
			return curr.value
		}
	}
	var zero V
	return zero
}

func (t *tree[K, V, O]) Contains(k K) bool {
	curr := t.root
	for curr != nil {
		if t.ordering.Less(k, curr.key) {
			curr = curr.left
		} else if t.ordering.Less(curr.key, k) {
			curr = curr.right
		} else {
			return true
		}
	}
	return false
}

func (t *tree[K, V, O]) First() (K, V) {
	curr := t.root
	for curr != nil {
		if curr.left == nil {
			return curr.key, curr.value
		}
		curr = curr.left
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV
}

func (t *tree[K, V, O]) Last() (K, V) {
	curr := t.root
	for curr != nil {
		if curr.right == nil {
			return curr.key, curr.value
		}
		curr = curr.right
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV
}

func (t *tree[K, V, O]) rebalance(curr *node[K, V]) *node[K, V] {
	imbalance := t.imbalance(curr)
	newSubtreeRoot := curr
	if imbalance > 1 {
		if t.imbalance(curr.left) < 0 {
			t.rotateLeft(curr.left)
		}
		newSubtreeRoot = t.rotateRight(curr)
		t.gen++
	} else if imbalance < -1 {
		if t.imbalance(curr.right) > 0 {
			t.rotateRight(curr.right)
		}
		newSubtreeRoot = t.rotateLeft(curr)
		t.gen++
	} else {
		t.setHeight(curr)
	}
	return newSubtreeRoot
}

//      parent                           parent
//        │                                │
//        b                                d
//   ┌────┴────┐                     ┌─────┴─────┐
//   a         d        ╶──>         b           e
//          ┌──┴──┐               ┌──┴──┐
//          c     e               a     c
func (t *tree[K, V, O]) rotateLeft(b *node[K, V]) *node[K, V] {
	parent := b.parent
	d := b.right
	c := d.left

	d.left = b
	b.parent = d

	b.right = c
	if c != nil {
		c.parent = b
	}

	d.parent = parent

	if parent != nil {
		if parent.left == b {
			parent.left = d
		} else {
			parent.right = d
		}
	}

	t.setHeight(b)
	t.setHeight(d)
	return d
}

//          parent                      parent
//            │                           │
//            d                           b
//      ┌─────┴─────┐                ┌────┴────┐
//      b           e      ╶──>      a         d
//   ┌──┴──┐                                ┌──┴──┐
//   a     c                                c     e
func (t *tree[K, V, O]) rotateRight(d *node[K, V]) *node[K, V] {
	parent := d.parent
	b := d.left
	c := b.right

	d.left = c
	if c != nil {
		c.parent = d
	}

	b.right = d
	d.parent = b

	b.parent = parent

	if parent != nil {
		if parent.left == d {
			parent.left = b
		} else {
			parent.right = b
		}
	}

	t.setHeight(d)
	t.setHeight(b)
	return b
}

// The height of x's left node, or -1 if no child.
func (t *tree[K, V, O]) leftHeight(x *node[K, V]) int {
	if x.left != nil {
		return x.left.height
	}
	return -1
}

// The height of x's right node, or -1 if no child.
func (t *tree[K, V, O]) rightHeight(x *node[K, V]) int {
	if x.right != nil {
		return x.right.height
	}
	return -1
}

// imbalance is the difference in height between x's children.
// 0 means perfectly balanced.
// >0 means the left tree is higher.
// <0 means the right tree is higher.
func (t *tree[K, V, O]) imbalance(x *node[K, V]) int {
	return t.leftHeight(x) - t.rightHeight(x)
}

func (t *tree[K, V, O]) setHeight(x *node[K, V]) {
	x.height = xmath.Max(t.leftHeight(x), t.rightHeight(x)) + 1
}

type cursor[K any, V any, O xsort.Ordering[K]] struct {
	t *tree[K, V, O]
	// Should be manipulated via reset(), up(), left(), and right().
	curr *node[K, V]
	gen  int
}

func (c *cursor[K, V, O]) clone() cursor[K, V, O] {
	return cursor[K, V, O]{
		t:    c.t,
		curr: c.curr,
		gen:  c.gen,
	}
}

func (c *cursor[K, V, O]) Next() {
	if !c.Ok() {
		return
	}
	if c.gen != c.t.gen {
		// Tree has changed structure, must re-seek to find our place.
		c.SeekFirstGreater(c.curr.key)
		return
	}
	if c.curr.right != nil {
		c.right()
		for c.curr.left != nil {
			c.left()
		}
	} else {
		prev := c.curr
		c.up()
		for c.curr != nil && c.t.ordering.Less(c.curr.key, prev.key) {
			c.up()
		}
	}
}
func (c *cursor[K, V, O]) Prev() {
	if !c.Ok() {
		return
	}
	if c.gen != c.t.gen && c.curr != nil {
		c.SeekLastLess(c.curr.key)
		return
	}
	if c.curr.left != nil {
		c.left()
		for c.curr.right != nil {
			c.right()
		}
	} else {
		prev := c.curr
		c.up()
		for c.curr != nil && c.t.ordering.Less(prev.key, c.curr.key) {
			c.up()
		}
	}
}

func (c *cursor[K, V, O]) Ok() bool {
	return c.curr != nil
}

func (c *cursor[K, V, O]) Key() K {
	return c.curr.key
}

func (c *cursor[K, V, O]) Value() V {
	return c.curr.value
}

func (c *cursor[K, V, O]) seek(k K) bool {
	if !c.reset() {
		return false
	}
	for {
		if c.curr.left != nil && c.t.ordering.Less(k, c.curr.key) {
			c.left()
		} else if c.curr.right != nil && c.t.ordering.Less(c.curr.key, k) {
			c.right()
		} else {
			break
		}
	}
	c.gen = c.t.gen
	return true
}

func (c *cursor[K, V, O]) SeekFirst() {
	if !c.reset() {
		return
	}
	for c.curr.left != nil {
		c.left()
	}
	c.gen = c.t.gen
}

func (c *cursor[K, V, O]) SeekLast() {
	if !c.reset() {
		return
	}
	for c.curr.right != nil {
		c.right()
	}
	c.gen = c.t.gen
}

func (c *cursor[K, V, O]) SeekLastLess(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.LessOrEqual[K, O](k, c.curr.key) {
		c.Prev()
	}
}

func (c *cursor[K, V, O]) SeekLastLessOrEqual(k K) {
	if !c.seek(k) {
		return
	}
	if c.t.ordering.Less(k, c.curr.key) {
		c.Prev()
	}
}

func (c *cursor[K, V, O]) SeekFirstGreaterOrEqual(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.Greater[K, O](k, c.curr.key) {
		c.Next()
	}
}

func (c *cursor[K, V, O]) SeekFirstGreater(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.GreaterOrEqual[K, O](k, c.curr.key) {
		c.Next()
	}
}

type forwardIterator[K any, V any, O xsort.Ordering[K]] struct {
	c cursor[K, V, O]
}

func (iter *forwardIterator[K, V, O]) Next() (KVPair[K, V], bool) {
	if !iter.c.Ok() {
		var zero KVPair[K, V]
		return zero, false
	}
	k := iter.c.Key()
	v := iter.c.Value()
	iter.c.Next()
	return KVPair[K, V]{k, v}, true
}

func (c *cursor[K, V, O]) Forward() iterator.Iterator[KVPair[K, V]] {
	c2 := c.clone()
	if c2.gen != c2.t.gen && c2.Ok() {
		c2.SeekFirstGreaterOrEqual(c2.Key())
	}
	return &forwardIterator[K, V, O]{c: c2}
}

type backwardIterator[K any, V any, O xsort.Ordering[K]] struct {
	c cursor[K, V, O]
}

func (iter *backwardIterator[K, V, O]) Next() (KVPair[K, V], bool) {
	if !iter.c.Ok() {
		var zero KVPair[K, V]
		return zero, false
	}
	k := iter.c.Key()
	v := iter.c.Value()
	iter.c.Prev()
	return KVPair[K, V]{k, v}, true
}

func (c *cursor[K, V, O]) Backward() iterator.Iterator[KVPair[K, V]] {
	c2 := c.clone()
	if c2.gen != c2.t.gen && c2.Ok() {
		c2.SeekLastLessOrEqual(c2.Key())
	}
	return &backwardIterator[K, V, O]{c: c2}
}

func (c *cursor[K, V, O]) reset() bool {
	c.curr = c.t.root
	return c.curr != nil
}
func (c *cursor[K, V, O]) up() {
	c.curr = c.curr.parent
}
func (c *cursor[K, V, O]) left() {
	c.curr = c.curr.left
}
func (c *cursor[K, V, O]) right() {
	c.curr = c.curr.right
}

func (t *tree[K, V, O]) Cursor() cursor[K, V, O] {
	c := cursor[K, V, O]{t: t}
	c.SeekFirst()
	return c
}
