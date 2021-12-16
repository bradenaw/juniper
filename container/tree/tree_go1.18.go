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

type iterState int

const (
	iterBeforeFirst iterState = iota
	iterAt
	iterAfterLast
)

type treeIterator[K any, V any] struct {
	t *tree[K, V]
	// Always an ancestor chain, that is, stack[i] is the parent of stack[i+1], or:
	//   (stack[i].left == stack[i+1] || stack[i].right == stack[i+1])
	//
	// Should be manipulated via reset(), up(), left(), and right().
	stack []*node[K, V]
	state iterState
	gen   int
}

func (iter *treeIterator[K, V]) Next() bool {
	if iter.state == iterBeforeFirst {
		iter.SeekStart()
		return len(iter.stack) > 0
	} else if iter.state == iterAfterLast {
		return false
	} else if iter.gen != iter.t.gen && len(iter.stack) > 0 {
		// Iterator is not already done and the tree has changed structure, must re-seek to find our
		// place.
		iter.SeekFirstGreater(iter.stack[len(iter.stack)-1].key)
		return len(iter.stack) > 0
	}
	curr := iter.curr()
	if curr.right != nil {
		iter.right()
		for iter.curr().left != nil {
			iter.left()
		}
	} else {
		prev := curr
		iter.up()
		for len(iter.stack) > 0 && iter.t.less(iter.curr().key, prev.key) {
			iter.up()
		}
	}
	if len(iter.stack) == 0 {
		iter.state = iterAfterLast
		return false
	}
	return true
}
func (iter *treeIterator[K, V]) SeekStart() {
	iter.reset()
	if len(iter.stack) == 0 {
		return
	}

	for iter.curr().left != nil {
		iter.left()
	}
	iter.gen = iter.t.gen
}
func (iter *treeIterator[K, V]) SeekFirstGreater(k K) {
	iter.reset()
	if len(iter.stack) == 0 {
		return
	}

	for {
		if iter.curr().left != nil && iter.t.less(k, iter.curr().key) {
			iter.left()
		} else if iter.curr().right != nil && iter.t.less(iter.curr().key, k) {
			iter.right()
		} else {
			break
		}
	}
	iter.gen = iter.t.gen
	iter.state = iterAt

	curr := iter.curr()
	if iter.t.less(curr.key, k) || !iter.t.less(k, curr.key) {
		// If less or equal, we need to advance to the next one.
		iter.Next()
	}
}
func (iter *treeIterator[K, V]) curr() *node[K, V] {
	return iter.stack[len(iter.stack)-1]
}

func (iter *treeIterator[K, V]) reset() {
	slices.Clear(iter.stack)
	if iter.t.root == nil {
		iter.state = iterAfterLast
		iter.stack = iter.stack[:0]
		return
	}
	iter.state = iterAt
	iter.stack = append(iter.stack[:0], iter.t.root)
}
func (iter *treeIterator[K, V]) up() {
	iter.stack = iter.stack[:len(iter.stack)-1]
}
func (iter *treeIterator[K, V]) left() {
	iter.stack = append(iter.stack, iter.curr().left)
}
func (iter *treeIterator[K, V]) right() {
	iter.stack = append(iter.stack, iter.curr().right)
}
func (iter *treeIterator[K, V]) Item() KVPair[K, V] {
	curr := iter.curr()
	return KVPair[K, V]{curr.key, curr.value}
}

func (t *tree[K, V]) Iterate() iterator.Iterator[KVPair[K, V]] {
	iter := &treeIterator[K, V]{t: t}
	return iter
}
