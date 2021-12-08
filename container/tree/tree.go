package tree

import (
	"github.com/bradenaw/xstd/xsort"
)

type tree[T any] struct {
	root *node[T]
	less xsort.Less[T]
	size int
}

type node[T any] struct {
	left *node[T]
	right *node[T]
	value T
}

func newTree[T any](less xsort.Less[T]) tree[T] {
	return tree[T]{
		root: nil,
		less: less,
		size: 0,
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
				curr.left = &node[T]{value:item}
				t.size++
				return
			}
			curr = curr.left
		} else if t.less(curr.value, item) {
			if curr.right == nil {
				curr.right = &node[T]{value:item}
				t.size++
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

// func (t *tree) Contains(item T) bool {
// }
