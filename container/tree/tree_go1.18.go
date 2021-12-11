//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/slices"
	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xsort"
)

// tree is an AVL tree: https://en.wikipedia.org/wiki/AVL_tree
type tree[T any] struct {
	root *node[T]
	less xsort.Less[T]
	size int
	// Incremented whenever the tree structure changes, so that iterators know to reset.
	gen int
}

type node[T any] struct {
	left  *node[T]
	right *node[T]
	// The height of this node. Leaves have height 0, internal nodes are one higher than their
	// highest child.
	height int
	// iterator depends on the value of this not changing with respect to tree.less.
	value T
}

func newTree[T any](less xsort.Less[T]) *tree[T] {
	return &tree[T]{
		root: nil,
		less: less,
		size: 0,
		gen:  1,
	}
}

func (t *tree[T]) Put(item T) {
	if t.root == nil {
		t.root = &node[T]{
			value: item,
		}
		t.size++
		return
	}
	var added bool
	t.root, added = t.putTraverse(item, t.root)
	if added {
		t.size++
		t.gen++
	}
}

// Returns true if a new node got added.
func (t *tree[T]) putTraverse(item T, curr *node[T]) (*node[T], bool) {
	if curr == nil {
		return &node[T]{value: item}, true
	}

	var added bool
	if t.less(item, curr.value) {
		curr.left, added = t.putTraverse(item, curr.left)
	} else if t.less(curr.value, item) {
		curr.right, added = t.putTraverse(item, curr.right)
	} else {
		curr.value = item
		return curr, false
	}

	return t.rebalance(curr), added
}

func (t *tree[T]) Delete(item T) {
	var deleted bool
	t.root, deleted = t.deleteTraverse(item, t.root)
	if deleted {
		t.size--
		t.gen++
	}
}

func (t *tree[T]) deleteTraverse(item T, curr *node[T]) (*node[T], bool) {
	if curr == nil {
		// item isn't in the tree
		return nil, false
	}
	var deleted bool
	if t.less(item, curr.value) {
		curr.left, deleted = t.deleteTraverse(item, curr.left)
	} else if t.less(curr.value, item) {
		curr.right, deleted = t.deleteTraverse(item, curr.right)
	} else {
		// curr contains item
		if curr.left != nil && curr.right != nil {
			// curr has both children, so replace it with its successor - one move right and then
			// all the way left. Since we're removing the successor from this subtree, correct the
			// heihts all the way down
			in := &curr.right
			successor := curr.right
			successor.height = t.rightHeight(successor) + 1
			for successor.left != nil {
				in = &successor.left
				successor = successor.left
				successor.height = t.rightHeight(successor) + 1
			}
			*in = successor.right
			successor.left = curr.left
			successor.right = t.rebalance(curr.right)
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

func (t *tree[T]) Get(item T) T {
	curr := t.root
	for curr != nil {
		if t.less(item, curr.value) {
			curr = curr.left
		} else if t.less(curr.value, item) {
			curr = curr.right
		} else {
			return curr.value
		}
	}
	var zero T
	return zero
}

func (t *tree[T]) Contains(item T) bool {
	curr := t.root
	for curr != nil {
		if t.less(item, curr.value) {
			curr = curr.left
		} else if t.less(curr.value, item) {
			curr = curr.right
		} else {
			return true
		}
	}
	return false
}

func (t *tree[T]) rebalance(curr *node[T]) *node[T] {
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
func (t *tree[T]) rotateLeft(b *node[T]) *node[T] {
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
func (t *tree[T]) rotateRight(d *node[T]) *node[T] {
	b := d.left
	c := b.right
	d.left = c
	b.right = d
	t.setHeight(d)
	t.setHeight(b)
	return b
}

// The height of x's left node, or -1 if no child.
func (t *tree[T]) leftHeight(x *node[T]) int {
	if x.left != nil {
		return x.left.height
	}
	return -1
}

// The height of x's right node, or -1 if no child.
func (t *tree[T]) rightHeight(x *node[T]) int {
	if x.right != nil {
		return x.right.height
	}
	return -1
}

// imbalance is the difference in height between x's children.
// 0 means perfectly balanced.
// >0 means the left tree is higher.
// <0 means the right tree is higher.
func (t *tree[T]) imbalance(x *node[T]) int {
	return t.leftHeight(x) - t.rightHeight(x)
}

func (t *tree[T]) setHeight(x *node[T]) {
	x.height = xmath.Max(t.leftHeight(x), t.rightHeight(x)) + 1
}

type iterState int

const (
	iterBeforeFirst iterState = iota
	iterAt
	iterAfterLast
)

type treeIterator[T any] struct {
	t *tree[T]
	// Always an ancestor chain, that is, stack[i] is the parent of stack[i+1], or:
	//   (stack[i].left == stack[i+1] || stack[i].right == stack[i+1])
	//
	// Should be manipulated via reset(), up(), left(), and right().
	stack []*node[T]
	state iterState
	gen   int
}

func (iter *treeIterator[T]) Next() bool {
	if iter.state == iterBeforeFirst {
		iter.SeekStart()
		return len(iter.stack) > 0
	} else if iter.state == iterAfterLast {
		return false
	} else if iter.gen != iter.t.gen && len(iter.stack) > 0 {
		// Iterator is not already done and the tree has changed structure, must re-seek to find our
		// place.
		iter.SeekFirstGreater(iter.stack[len(iter.stack)-1].value)
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
		for len(iter.stack) > 0 && iter.t.less(iter.curr().value, prev.value) {
			iter.up()
		}
	}
	if len(iter.stack) == 0 {
		iter.state = iterAfterLast
		return false
	}
	return true
}
func (iter *treeIterator[T]) SeekStart() {
	iter.reset()
	if len(iter.stack) == 0 {
		return
	}

	for iter.curr().left != nil {
		iter.left()
	}
	iter.gen = iter.t.gen
}
func (iter *treeIterator[T]) SeekFirstGreater(item T) {
	iter.reset()
	if len(iter.stack) == 0 {
		return
	}

	for {
		if iter.curr().left != nil && iter.t.less(item, iter.curr().value) {
			iter.left()
		} else if iter.curr().right != nil && iter.t.less(iter.curr().value, item) {
			iter.right()
		} else {
			break
		}
	}
	iter.gen = iter.t.gen
	iter.state = iterAt

	curr := iter.curr().value
	if iter.t.less(curr, item) || !iter.t.less(item, curr) {
		// If less or equal, we need to advance to the next one.
		iter.Next()
	}
}
func (iter *treeIterator[T]) curr() *node[T] {
	return iter.stack[len(iter.stack)-1]
}

func (iter *treeIterator[T]) reset() {
	slices.Clear(iter.stack)
	if iter.t.root == nil {
		iter.state = iterAfterLast
		iter.stack = iter.stack[:0]
		return
	}
	iter.state = iterAt
	iter.stack = append(iter.stack[:0], iter.t.root)
}
func (iter *treeIterator[T]) up() {
	iter.stack = iter.stack[:len(iter.stack)-1]
}
func (iter *treeIterator[T]) left() {
	iter.stack = append(iter.stack, iter.curr().left)
}
func (iter *treeIterator[T]) right() {
	iter.stack = append(iter.stack, iter.curr().right)
}
func (iter *treeIterator[T]) Item() T {
	return iter.curr().value
}

func (t *tree[T]) Iterate() iterator.Iterator[T] {
	iter := &treeIterator[T]{
		stack: []*node[T]{},
		t:     t,
	}
	return iterator.New(func() (T, bool) {
		ok := iter.Next()
		if !ok {
			var zero T
			return zero, false
		}
		return iter.Item(), true
	})
}
