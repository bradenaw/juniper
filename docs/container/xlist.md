# `package xlist`

```
import "github.com/bradenaw/juniper/container/xlist"
```

## Overview

Package xlist contains extensions to the standard library package container/list.


## Index

<samp><a href="#List">type List</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Back">func (l *List[T]) Back() *Node[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Clear">func (l *List[T]) Clear()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Front">func (l *List[T]) Front() *Node[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#InsertAfter">func (l *List[T]) InsertAfter(value T, mark *Node[T]) *Node[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#InsertBefore">func (l *List[T]) InsertBefore(value T, mark *Node[T]) *Node[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (l *List[T]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#MoveAfter">func (l *List[T]) MoveAfter(node *Node[T], mark *Node[T])</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#MoveBefore">func (l *List[T]) MoveBefore(node *Node[T], mark *Node[T])</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#MoveToBack">func (l *List[T]) MoveToBack(node *Node[T])</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#MoveToFront">func (l *List[T]) MoveToFront(node *Node[T])</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PushBack">func (l *List[T]) PushBack(value T) *Node[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PushFront">func (l *List[T]) PushFront(value T) *Node[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Remove">func (l *List[T]) Remove(node *Node[T])</a></samp>

<samp><a href="#Node">type Node</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Next">func (n *Node[T]) Next() *Node[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Prev">func (n *Node[T]) Prev() *Node[T]</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

This section is empty.
## Types

<h3><a id="List"></a><samp>type List</samp></h3>
```go
type List[T any] struct {
	// contains filtered or unexported fields
}
```

List is a doubly-linked list.


<h3><a id="Back"></a><samp>func (l *<a href="#List">List</a>[T]) Back() *<a href="#Node">Node</a>[T]</samp></h3>

Back returns the node at the back of the list.


<h3><a id="Clear"></a><samp>func (l *<a href="#List">List</a>[T]) Clear()</samp></h3>

Clear removes all nodes from the list.


<h3><a id="Front"></a><samp>func (l *<a href="#List">List</a>[T]) Front() *<a href="#Node">Node</a>[T]</samp></h3>

Front returns the node at the front of the list.


<h3><a id="InsertAfter"></a><samp>func (l *<a href="#List">List</a>[T]) InsertAfter(value T, mark *<a href="#Node">Node</a>[T]) *<a href="#Node">Node</a>[T]</samp></h3>

InsertBefore adds a new node with the given value after the node mark.


<h3><a id="InsertBefore"></a><samp>func (l *<a href="#List">List</a>[T]) InsertBefore(value T, mark *<a href="#Node">Node</a>[T]) *<a href="#Node">Node</a>[T]</samp></h3>

InsertBefore adds a new node with the given value before the node mark.


<h3><a id="Len"></a><samp>func (l *<a href="#List">List</a>[T]) Len() int</samp></h3>

Len returns the number of items in the list.


<h3><a id="MoveAfter"></a><samp>func (l *<a href="#List">List</a>[T]) MoveAfter(node *<a href="#Node">Node</a>[T], mark *<a href="#Node">Node</a>[T])</samp></h3>

MoveAfter moves node just after mark. Afterwards, mark.Next() == node && node.Prev() == mark.


<h3><a id="MoveBefore"></a><samp>func (l *<a href="#List">List</a>[T]) MoveBefore(node *<a href="#Node">Node</a>[T], mark *<a href="#Node">Node</a>[T])</samp></h3>

MoveBefore moves node just before mark. Afterwards, mark.Prev() == node && node.Next() == mark.


<h3><a id="MoveToBack"></a><samp>func (l *<a href="#List">List</a>[T]) MoveToBack(node *<a href="#Node">Node</a>[T])</samp></h3>

MoveToFront moves node to the back of the list.


<h3><a id="MoveToFront"></a><samp>func (l *<a href="#List">List</a>[T]) MoveToFront(node *<a href="#Node">Node</a>[T])</samp></h3>

MoveToFront moves node to the front of the list.


<h3><a id="PushBack"></a><samp>func (l *<a href="#List">List</a>[T]) PushBack(value T) *<a href="#Node">Node</a>[T]</samp></h3>

PushFront adds a new node with the given value to the back of the list.


<h3><a id="PushFront"></a><samp>func (l *<a href="#List">List</a>[T]) PushFront(value T) *<a href="#Node">Node</a>[T]</samp></h3>

PushFront adds a new node with the given value to the front of the list.


<h3><a id="Remove"></a><samp>func (l *<a href="#List">List</a>[T]) Remove(node *<a href="#Node">Node</a>[T])</samp></h3>

Remove removes node from the list.


<h3><a id="Node"></a><samp>type Node</samp></h3>
```go
type Node[T any] struct {

	// Value is user-controlled, and never modified by this package.
	Value T
	// contains filtered or unexported fields
}
```

Node is a node in a linked-list.


<h3><a id="Next"></a><samp>func (n *<a href="#Node">Node</a>[T]) Next() *<a href="#Node">Node</a>[T]</samp></h3>

Next returns the next node in the list that n is a part of, if there is one.


<h3><a id="Prev"></a><samp>func (n *<a href="#Node">Node</a>[T]) Prev() *<a href="#Node">Node</a>[T]</samp></h3>

Prev returns the previous node in the list that n is a part of, if there is one.


