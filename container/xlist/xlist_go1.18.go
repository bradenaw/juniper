//go:build go1.18

// Package xlist contains extensions to the standard library package container/list.
package xlist

// List is a doubly-linked list.
type List[T any] struct {
	front *Node[T]
	back  *Node[T]
	size  int
}

// Len returns the number of items in the list.
func (l *List[T]) Len() int        { return l.size }

// Front returns the node at the front of the list.
func (l *List[T]) Front() *Node[T] { return l.front }

// Back returns the node at the back of the list.
func (l *List[T]) Back() *Node[T]  { return l.back }

// Clear removes all nodes from the list.
func (l *List[T]) Clear() { l.front = nil; l.back = nil; l.size = 0 }

// PushFront adds a new node with the given value to the front of the list.
func (l *List[T]) PushFront(value T) *Node[T] {
	node := &Node[T]{
		next:  l.front,
		Value: value,
	}
	if l.front != nil {
		l.front.prev = node
	}
	l.front = node
	if l.back == nil {
		l.back = node
	}
	l.size++
	return node
}

// PushFront adds a new node with the given value to the back of the list.
func (l *List[T]) PushBack(value T) *Node[T] {
	node := &Node[T]{
		prev:  l.back,
		Value: value,
	}
	if l.back != nil {
		l.back.next = node
	}
	l.back = node
	if l.front == nil {
		l.front = node
	}
	l.size++
	return node
}

// InsertBefore adds a new node with the given value before the node mark.
func (l *List[T]) InsertBefore(value T, mark *Node[T]) *Node[T] {
	node := &Node[T]{
		Value: value,
		prev: mark.prev,
		next: mark,
	}
	mark.prev = node
	if node.prev != nil {
		node.prev.next = node
	}
	if l.front == mark {
		l.front = node
	}
	l.size++
	return node
}

// InsertBefore adds a new node with the given value after the node mark.
func (l *List[T]) InsertAfter(value T, mark *Node[T]) *Node[T] {
	node := &Node[T]{
		Value: value,
		prev: mark,
		next: mark.next,
	}
	mark.next = node
	if node.next != nil {
		node.next.prev = node
	}
	if l.back == mark {
		l.back = node
	}
	l.size++
	return node
}

// Remove removes node from the list.
func (l *List[T]) Remove(node *Node[T]) {
	l.remove(node)
	node.prev = nil
	node.next = nil
	l.size--
}

func (l *List[T]) remove(node *Node[T]) {
	if l.front == node {
		l.front = l.front.next
	} else {
		node.prev.next = node.next
	}
	if l.back == node {
		l.back = l.back.prev
	} else {
		node.next.prev = node.prev
	}
}

// MoveBefore moves node just before mark. Afterwards, mark.Prev() == node && node.Next() == mark.
func (l *List[T]) MoveBefore(node *Node[T], mark *Node[T]) {
	if node == mark {
		return
	}
	l.remove(node)
	node.prev = mark.prev
	mark.prev = node
	node.next = mark
	if node.prev != nil {
		node.prev.next = node
	}
	if l.front == mark {
		l.front = node
	}
}

// MoveAfter moves node just after mark. Afterwards, mark.Next() == node && node.Prev() == mark.
func (l *List[T]) MoveAfter(node *Node[T], mark *Node[T]) {
	if node == mark {
		return
	}
	l.remove(node)
	node.next = mark.next
	mark.next = node
	node.prev = mark
	if node.next != nil {
		node.next.prev = node
	}
	if l.back == mark {
		l.back = node
	}
}

// MoveToFront moves node to the front of the list.
func (l *List[T]) MoveToFront(node *Node[T]) {
	l.MoveBefore(node, l.Front())
}

// MoveToFront moves node to the back of the list.
func (l *List[T]) MoveToBack(node *Node[T]) {
	l.MoveAfter(node, l.Back())
}

// Node is a node in a linked-list.
type Node[T any] struct {
	prev  *Node[T]
	next  *Node[T]
	// Value is user-controlled, and never modified by this package.
	Value T
}

// Next returns the next node in the list that n is a part of, if there is one.
func (n *Node[T]) Next() *Node[T] {
	return n.next
}

// Prev returns the previous node in the list that n is a part of, if there is one.
func (n *Node[T]) Prev() *Node[T] {
	return n.prev
}
