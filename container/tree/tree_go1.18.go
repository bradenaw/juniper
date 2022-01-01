//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/slices"
	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xsort"
)

// tree is an AVL tree: https://en.wikipedia.org/wiki/AVL_tree
type tree[K any, V any] struct {
	root *node[K, V]
	less xsort.Less[K]
	size int
	// Incremented whenever the tree structure changes, so that iterators know to reset.
	gen int
}

type node[K any, V any] struct {
	left  *node[K, V]
	right *node[K, V]
	// The height of this node. Leaves have height 0, internal nodes are one higher than their
	// highest child.
	height int
	// iterator depends on key not changing with respect to tree.less.
	key   K
	value V
}

func newTree[K any, V any](less xsort.Less[K]) *tree[K, V] {
	return &tree[K, V]{
		root: nil,
		less: less,
		size: 0,
		gen:  1,
	}
}

func (t *tree[K, V]) Put(k K, v V) {
	if t.root == nil {
		t.root = &node[K, V]{
			key:   k,
			value: v,
		}
		t.size++
		return
	}
	var added bool
	t.root, added = t.putTraverse(k, v, t.root)
	if added {
		t.size++
		t.gen++
	}
}

func (t *tree[K, V]) putTraverse(k K, v V, curr *node[K, V]) (_newCurr *node[K, V], _added bool) {
	if curr == nil {
		return &node[K, V]{key: k, value: v}, true
	}

	var added bool
	if t.less(k, curr.key) {
		curr.left, added = t.putTraverse(k, v, curr.left)
	} else if t.less(curr.key, k) {
		curr.right, added = t.putTraverse(k, v, curr.right)
	} else {
		curr.key = k
		curr.value = v
		return curr, false
	}

	return t.rebalance(curr), added
}

func (t *tree[K, V]) Delete(k K) {
	var deleted bool
	t.root, deleted = t.deleteTraverse(k, t.root)
	if deleted {
		t.size--
		t.gen++
	}
}

func (t *tree[K, V]) deleteTraverse(k K, curr *node[K, V]) (_newCurr *node[K, V], _deleted bool) {
	if curr == nil {
		// k isn't in the tree
		return nil, false
	}
	var deleted bool
	if t.less(k, curr.key) {
		curr.left, deleted = t.deleteTraverse(k, curr.left)
	} else if t.less(curr.key, k) {
		curr.right, deleted = t.deleteTraverse(k, curr.right)
	} else {
		// curr contains k
		if curr.left != nil && curr.right != nil {
			// curr has both children, so replace it with its successor.
			right, successor := t.removeLeftmost(curr.right)
			successor.left = curr.left
			successor.right = right
			return t.rebalance(successor), true
		} else if curr.left != nil {
			// curr has only a left child, just hoist it upwards
			return curr.left, true
		} else {
			// curr has only a right child, just hoist it upwards
			return curr.right, true
		}
	}
	return t.rebalance(curr), deleted
}

func (t *tree[K, V]) removeLeftmost(
	curr *node[K, V],
) (_newCurr *node[K, V], _leftmost *node[K, V]) {
	if curr.left == nil {
		return curr.right, curr
	}
	var leftmost *node[K, V]
	curr.left, leftmost = t.removeLeftmost(curr.left)
	t.setHeight(curr)
	return t.rebalance(curr), leftmost
}

func (t *tree[K, V]) Get(k K) V {
	curr := t.root
	for curr != nil {
		if t.less(k, curr.key) {
			curr = curr.left
		} else if t.less(curr.key, k) {
			curr = curr.right
		} else {
			return curr.value
		}
	}
	var zero V
	return zero
}

func (t *tree[K, V]) Contains(k K) bool {
	curr := t.root
	for curr != nil {
		if t.less(k, curr.key) {
			curr = curr.left
		} else if t.less(curr.key, k) {
			curr = curr.right
		} else {
			return true
		}
	}
	return false
}

func (t *tree[K, V]) First() (K, V) {
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

func (t *tree[K, V]) Last() (K, V) {
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

func (t *tree[K, V]) rebalance(curr *node[K, V]) *node[K, V] {
	if curr == nil {
		return nil
	}
	imbalance := t.imbalance(curr)
	newSubtreeRoot := curr
	if imbalance > 1 {
		if t.imbalance(curr.left) < 0 {
			curr.left = t.rotateLeft(curr.left)
		}
		newSubtreeRoot = t.rotateRight(curr)
		t.gen++
	} else if imbalance < -1 {
		if t.imbalance(curr.right) > 0 {
			curr.right = t.rotateRight(curr.right)
		}
		newSubtreeRoot = t.rotateLeft(curr)
		t.gen++
	} else {
		t.setHeight(curr)
	}
	return newSubtreeRoot
}

//        b                                d
//   ┌────┴────┐                     ┌─────┴─────┐
//   a         d        ╶──>         b           e
//          ┌──┴──┐               ┌──┴──┐
//          c     e               a     c
func (t *tree[K, V]) rotateLeft(b *node[K, V]) *node[K, V] {
	d := b.right
	c := d.left
	d.left = b
	b.right = c
	t.setHeight(b)
	t.setHeight(d)
	return d
}

//            d                           b
//      ┌─────┴─────┐                ┌────┴────┐
//      b           e      ╶──>      a         d
//   ┌──┴──┐                                ┌──┴──┐
//   a     c                                c     e
func (t *tree[K, V]) rotateRight(d *node[K, V]) *node[K, V] {
	b := d.left
	c := b.right
	d.left = c
	b.right = d
	t.setHeight(d)
	t.setHeight(b)
	return b
}

// The height of x's left node, or -1 if no child.
func (t *tree[K, V]) leftHeight(x *node[K, V]) int {
	if x.left != nil {
		return x.left.height
	}
	return -1
}

// The height of x's right node, or -1 if no child.
func (t *tree[K, V]) rightHeight(x *node[K, V]) int {
	if x.right != nil {
		return x.right.height
	}
	return -1
}

// imbalance is the difference in height between x's children.
// 0 means perfectly balanced.
// >0 means the left tree is higher.
// <0 means the right tree is higher.
func (t *tree[K, V]) imbalance(x *node[K, V]) int {
	return t.leftHeight(x) - t.rightHeight(x)
}

func (t *tree[K, V]) setHeight(x *node[K, V]) {
	x.height = xmath.Max(t.leftHeight(x), t.rightHeight(x)) + 1
}

type cursor[K any, V any] struct {
	t *tree[K, V]
	// Always an ancestor chain, that is, stack[i] is the parent of stack[i+1], or:
	//   (stack[i].left == stack[i+1] || stack[i].right == stack[i+1])
	//
	// Should be manipulated via reset(), up(), left(), and right().
	stack []*node[K, V]
	gen   int
}

func (c *cursor[K, V]) clone() *cursor[K, V] {
	return &cursor[K, V]{
		t:     c.t,
		stack: slices.Clone(c.stack),
		gen:   c.gen,
	}
}

func (c *cursor[K, V]) Next() {
	if !c.Ok() {
		return
	}
	if c.gen != c.t.gen {
		// Tree has changed structure, must re-seek to find our place.
		c.SeekFirstGreater(c.curr().key)
		return
	}
	if c.curr().right != nil {
		c.right()
		for c.curr().left != nil {
			c.left()
		}
	} else {
		prev := c.curr()
		c.up()
		for len(c.stack) > 0 && c.t.less(c.curr().key, prev.key) {
			c.up()
		}
	}
}
func (c *cursor[K, V]) Prev() {
	if !c.Ok() {
		return
	}
	if c.gen != c.t.gen && len(c.stack) > 0 {
		c.SeekLastLess(c.curr().key)
		return
	}
	if c.curr().left != nil {
		c.left()
		for c.curr().right != nil {
			c.right()
		}
	} else {
		prev := c.curr()
		c.up()
		for len(c.stack) > 0 && c.t.less(prev.key, c.curr().key) {
			c.up()
		}
	}
}

// This isn't very sane behavior and probably needs some rethinking. The goal here is to get it to
// do something reasonable when the key the cursor is sitting at gets deleted. We could immediately
// axe the cursor, but that makes it so you can't filter the contents of a map by looping and
// removing items. We could just keep the cursor at the removed node and continue returning its
// contents, but then you get weirdness like:
//
//   t.Put(1, 1)
//   c := t.Cursor()
//   t.Delete(1)
//   t.Put(1, 2)
//   c.Value() // 1!
//
// Advancing inside Ok/Key/Value seems _extremely_ weird, though. Maybe the right thing to do is
// make Ok() return false, and thus make it illegal to call Key/Value, but legal to call
// Next/Prev/Seek?
func (c *cursor[K, V]) maybeRecover() {
	if len(c.stack) > 0 && c.gen != c.t.gen {
		c.SeekFirstGreaterOrEqual(c.curr().key)
	}
}

func (c *cursor[K, V]) Ok() bool {
	c.maybeRecover()
	return len(c.stack) > 0
}

func (c *cursor[K, V]) Key() K {
	c.maybeRecover()
	return c.curr().key
}

func (c *cursor[K, V]) Value() V {
	c.maybeRecover()
	return c.curr().value
}

func (c *cursor[K, V]) seek(k K) bool {
	if !c.reset() {
		return false
	}
	for {
		if c.curr().left != nil && c.t.less(k, c.curr().key) {
			c.left()
		} else if c.curr().right != nil && c.t.less(c.curr().key, k) {
			c.right()
		} else {
			break
		}
	}
	c.gen = c.t.gen
	return true
}

func (c *cursor[K, V]) SeekFirst() {
	if !c.reset() {
		return
	}
	for c.curr().left != nil {
		c.left()
	}
	c.gen = c.t.gen
}

func (c *cursor[K, V]) SeekLast() {
	if !c.reset() {
		return
	}
	for c.curr().right != nil {
		c.right()
	}
	c.gen = c.t.gen
}

func (c *cursor[K, V]) SeekLastLess(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.LessOrEqual(c.t.less, k, c.curr().key) {
		c.Prev()
	}
}

func (c *cursor[K, V]) SeekLastLessOrEqual(k K) {
	if !c.seek(k) {
		return
	}
	if c.t.less(k, c.curr().key) {
		c.Prev()
	}
}

func (c *cursor[K, V]) SeekFirstGreaterOrEqual(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.Greater(c.t.less, k, c.curr().key) {
		c.Next()
	}
}

func (c *cursor[K, V]) SeekFirstGreater(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.GreaterOrEqual(c.t.less, k, c.curr().key) {
		c.Next()
	}
}

func (c *cursor[K, V]) Forward() iterator.Iterator[KVPair[K, V]] {
	c2 := c.clone()
	if c2.gen != c2.t.gen && c2.Ok() {
		c2.SeekFirstGreaterOrEqual(c2.Key())
	}
	return iterator.FromNext(func() (KVPair[K, V], bool) {
		if !c2.Ok() {
			var zero KVPair[K, V]
			return zero, false
		}
		k := c2.Key()
		v := c2.Value()
		c2.Next()
		return KVPair[K, V]{k, v}, true
	})
}

func (c *cursor[K, V]) Backward() iterator.Iterator[KVPair[K, V]] {
	c2 := c.clone()
	if c2.gen != c2.t.gen && c2.Ok() {
		c2.SeekLastLessOrEqual(c2.Key())
	}
	return iterator.FromNext(func() (KVPair[K, V], bool) {
		if !c2.Ok() {
			var zero KVPair[K, V]
			return zero, false
		}
		k := c2.Key()
		v := c2.Value()
		c2.Prev()
		return KVPair[K, V]{k, v}, true
	})
}

func (c *cursor[K, V]) curr() *node[K, V] {
	return c.stack[len(c.stack)-1]
}
func (c *cursor[K, V]) reset() bool {
	slices.Clear(c.stack)
	if c.t.root == nil {
		c.stack = nil
		return false
	}
	c.stack = append(c.stack[:0], c.t.root)
	return true
}
func (c *cursor[K, V]) up() {
	c.stack = c.stack[:len(c.stack)-1]
}
func (c *cursor[K, V]) left() {
	c.stack = append(c.stack, c.curr().left)
}
func (c *cursor[K, V]) right() {
	c.stack = append(c.stack, c.curr().right)
}

func (t *tree[K, V]) Cursor() cursor[K, V] {
	c := cursor[K, V]{t: t}
	c.SeekFirst()
	return c
}
