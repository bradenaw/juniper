//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/slices"
	"github.com/bradenaw/juniper/xsort"
)

type tree[T any] struct {
	// TODO: rebalancing.
	root *node[T]
	less xsort.Less[T]
	size int
	// Incremented whenever the tree structure changes, so that iterators know to reset.
	gen int
}

type node[T any] struct {
	left  *node[T]
	right *node[T]
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
	curr := t.root
	for {
		if t.less(item, curr.value) {
			if curr.left == nil {
				curr.left = &node[T]{value: item}
				t.size++
				t.gen++
				return
			}
			curr = curr.left
		} else if t.less(curr.value, item) {
			if curr.right == nil {
				curr.right = &node[T]{value: item}
				t.size++
				t.gen++
				return
			}
			curr = curr.right
		} else {
			curr.value = item
			return
		}
	}
}

func (t *tree[T]) Delete(item T) {
	in := &t.root
	curr := t.root
	for curr != nil {
		if t.less(item, curr.value) {
			in = &curr.left
			curr = curr.left
		} else if t.less(curr.value, item) {
			in = &curr.right
			curr = curr.right
		} else {
			if curr.left != nil && curr.right != nil {
				*in = curr.right
				curr2 := &curr.right.left
				for *curr2 != nil {
					curr2 = &(*curr2).left
				}
				*curr2 = curr.left
			} else if curr.left != nil {
				*in = curr.left
			} else {
				*in = curr.right
			}
			t.size--
			t.gen++
			return
		}
	}
}

func (t *tree[T]) Get(item T) (T, bool) {
	curr := t.root
	for curr != nil {
		if t.less(item, curr.value) {
			curr = curr.left
		} else if t.less(curr.value, item) {
			curr = curr.right
		} else {
			return curr.value, true
		}
	}
	var zero T
	return zero, false
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
