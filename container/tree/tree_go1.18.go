//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xsort"
)

// tree is an AVL tree: https://en.wikipedia.org/wiki/AVL_tree
type tree[K any, V any] struct {
	root *node[K, V]
	less xsort.Less[K]
	size int
}

type node[K any, V any] struct {
	left  *node[K, V]
	right *node[K, V]
	// nil if root or unlinked
	parent *node[K, V]
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

	curr := t.root
	for {
		if t.less(k, curr.key) {
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
		} else if t.less(curr.key, k) {
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
		height := curr.height
		curr = t.rebalance(curr)
		if curr.parent == nil {
			t.root = curr
			break
		}
		if height == curr.height {
			break
		}
		curr = curr.parent
	}
	t.size++
}

func (t *tree[K, V]) Delete(k K) {
	in := &t.root
	var last *node[K, V]
	for {
		if *in == nil {
			// k isn't in the map
			return
		}
		curr := *in
		if t.less(k, curr.key) {
			in = &curr.left
		} else if t.less(curr.key, k) {
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
			// unlink so cursors know to re-seek
			curr.parent = nil
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
}

func (t *tree[K, V]) removeLeftmost(
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
	imbalance := t.imbalance(curr)
	newSubtreeRoot := curr
	if imbalance > 1 {
		if t.imbalance(curr.left) < 0 {
			t.rotateLeft(curr.left)
		}
		newSubtreeRoot = t.rotateRight(curr)
	} else if imbalance < -1 {
		if t.imbalance(curr.right) > 0 {
			t.rotateRight(curr.right)
		}
		newSubtreeRoot = t.rotateLeft(curr)
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
func (t *tree[K, V]) rotateLeft(b *node[K, V]) *node[K, V] {
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
func (t *tree[K, V]) rotateRight(d *node[K, V]) *node[K, V] {
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
	t    *tree[K, V]
	curr *node[K, V]
}

func (c *cursor[K, V]) clone() cursor[K, V] {
	return cursor[K, V]{
		t:    c.t,
		curr: c.curr,
	}
}

func (c *cursor[K, V]) Next() {
	if !c.Ok() {
		return
	}
	if c.currDeleted() {
		c.SeekFirstGreater(c.curr.key)
		return
	}
	if c.curr.right != nil {
		c.curr = c.curr.right
		for c.curr.left != nil {
			c.curr = c.curr.left
		}
	} else {
		prev := c.curr
		c.curr = c.curr.parent
		for c.curr != nil && c.t.less(c.curr.key, prev.key) {
			c.curr = c.curr.parent
		}
	}
}
func (c *cursor[K, V]) Prev() {
	if !c.Ok() {
		return
	}
	if c.currDeleted() {
		c.SeekLastLess(c.curr.key)
		return
	}
	if c.curr.left != nil {
		c.curr = c.curr.left
		for c.curr.right != nil {
			c.curr = c.curr.right
		}
	} else {
		prev := c.curr
		c.curr = c.curr.parent
		for c.curr != nil && c.t.less(prev.key, c.curr.key) {
			c.curr = c.curr.parent
		}
	}
}

func (c *cursor[K, V]) Ok() bool {
	return c.curr != nil
}

func (c *cursor[K, V]) Key() K {
	return c.curr.key
}

func (c *cursor[K, V]) Value() V {
	return c.curr.value
}

func (c *cursor[K, V]) seek(k K) bool {
	c.curr = c.t.root
	if c.curr == nil {
		return false
	}
	for {
		if c.curr.left != nil && c.t.less(k, c.curr.key) {
			c.curr = c.curr.left
		} else if c.curr.right != nil && c.t.less(c.curr.key, k) {
			c.curr = c.curr.right
		} else {
			break
		}
	}
	return true
}

func (c *cursor[K, V]) SeekFirst() {
	c.curr = c.t.root
	if c.curr == nil {
		return
	}
	for c.curr.left != nil {
		c.curr = c.curr.left
	}
}

func (c *cursor[K, V]) SeekLast() {
	c.curr = c.t.root
	if c.curr == nil {
		return
	}
	for c.curr.right != nil {
		c.curr = c.curr.right
	}
}

func (c *cursor[K, V]) SeekLastLess(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.LessOrEqual(c.t.less, k, c.curr.key) {
		c.Prev()
	}
}

func (c *cursor[K, V]) SeekLastLessOrEqual(k K) {
	if !c.seek(k) {
		return
	}
	if c.t.less(k, c.curr.key) {
		c.Prev()
	}
}

func (c *cursor[K, V]) SeekFirstGreaterOrEqual(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.Greater(c.t.less, k, c.curr.key) {
		c.Next()
	}
}

func (c *cursor[K, V]) SeekFirstGreater(k K) {
	if !c.seek(k) {
		return
	}
	if xsort.GreaterOrEqual(c.t.less, k, c.curr.key) {
		c.Next()
	}
}

type forwardIterator[K any, V any] struct {
	c cursor[K, V]
}

func (iter *forwardIterator[K, V]) Next() (KVPair[K, V], bool) {
	if !iter.c.Ok() {
		var zero KVPair[K, V]
		return zero, false
	}
	k := iter.c.Key()
	v := iter.c.Value()
	iter.c.Next()
	return KVPair[K, V]{k, v}, true
}

func (c *cursor[K, V]) Forward() iterator.Iterator[KVPair[K, V]] {
	c2 := c.clone()
	if c2.currDeleted() {
		c2.SeekFirstGreaterOrEqual(c2.Key())
	}
	return &forwardIterator[K, V]{c: c2}
}

type backwardIterator[K any, V any] struct {
	c cursor[K, V]
}

func (iter *backwardIterator[K, V]) Next() (KVPair[K, V], bool) {
	if !iter.c.Ok() {
		var zero KVPair[K, V]
		return zero, false
	}
	k := iter.c.Key()
	v := iter.c.Value()
	iter.c.Prev()
	return KVPair[K, V]{k, v}, true
}

func (c *cursor[K, V]) Backward() iterator.Iterator[KVPair[K, V]] {
	c2 := c.clone()
	if c2.currDeleted() {
		c2.SeekLastLessOrEqual(c2.Key())
	}
	return &backwardIterator[K, V]{c: c2}
}

func (c *cursor[K, V]) currDeleted() bool {
	return c.curr != nil && c.curr.parent == nil && c.t.root != c.curr
}

func (t *tree[K, V]) Cursor() cursor[K, V] {
	c := cursor[K, V]{t: t}
	c.SeekFirst()
	return c
}
