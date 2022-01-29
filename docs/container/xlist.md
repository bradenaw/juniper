# `package xlist`

```
import "github.com/bradenaw/juniper/container/xlist"
```

# Overview

Package xlist contains extensions to the standard library package container/list.


# Index

<pre><a href="#List">type List</a></pre>
<pre>    <a href="#Back">func (l *List[T]) Back() *Node[T]</a></pre>
<pre>    <a href="#Clear">func (l *List[T]) Clear()</a></pre>
<pre>    <a href="#Front">func (l *List[T]) Front() *Node[T]</a></pre>
<pre>    <a href="#InsertAfter">func (l *List[T]) InsertAfter(value T, mark *Node[T]) *Node[T]</a></pre>
<pre>    <a href="#InsertBefore">func (l *List[T]) InsertBefore(value T, mark *Node[T]) *Node[T]</a></pre>
<pre>    <a href="#Len">func (l *List[T]) Len() int</a></pre>
<pre>    <a href="#MoveAfter">func (l *List[T]) MoveAfter(node *Node[T], mark *Node[T])</a></pre>
<pre>    <a href="#MoveBefore">func (l *List[T]) MoveBefore(node *Node[T], mark *Node[T])</a></pre>
<pre>    <a href="#MoveToBack">func (l *List[T]) MoveToBack(node *Node[T])</a></pre>
<pre>    <a href="#MoveToFront">func (l *List[T]) MoveToFront(node *Node[T])</a></pre>
<pre>    <a href="#PushBack">func (l *List[T]) PushBack(value T) *Node[T]</a></pre>
<pre>    <a href="#PushFront">func (l *List[T]) PushFront(value T) *Node[T]</a></pre>
<pre>    <a href="#Remove">func (l *List[T]) Remove(node *Node[T])</a></pre>
<pre><a href="#Node">type Node</a></pre>
<pre>    <a href="#Next">func (n *Node[T]) Next() *Node[T]</a></pre>
<pre>    <a href="#Prev">func (n *Node[T]) Prev() *Node[T]</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

# Types

## <a id="List"></a><pre>type List</pre>
```go
type List[T any] struct {
	// contains filtered or unexported fields
}
```

List is a doubly-linked list.


## <a id="Back"></a><pre>func (l *<a href="#List">List</a>[T]) Back() *<a href="#Node">Node</a>[T]</pre>

Back returns the node at the back of the list.


## <a id="Clear"></a><pre>func (l *<a href="#List">List</a>[T]) Clear()</pre>

Clear removes all nodes from the list.


## <a id="Front"></a><pre>func (l *<a href="#List">List</a>[T]) Front() *<a href="#Node">Node</a>[T]</pre>

Front returns the node at the front of the list.


## <a id="InsertAfter"></a><pre>func (l *<a href="#List">List</a>[T]) InsertAfter(value T, mark *<a href="#Node">Node</a>[T]) *<a href="#Node">Node</a>[T]</pre>

InsertBefore adds a new node with the given value after the node mark.


## <a id="InsertBefore"></a><pre>func (l *<a href="#List">List</a>[T]) InsertBefore(value T, mark *<a href="#Node">Node</a>[T]) *<a href="#Node">Node</a>[T]</pre>

InsertBefore adds a new node with the given value before the node mark.


## <a id="Len"></a><pre>func (l *<a href="#List">List</a>[T]) Len() int</pre>

Len returns the number of items in the list.


## <a id="MoveAfter"></a><pre>func (l *<a href="#List">List</a>[T]) MoveAfter(node *<a href="#Node">Node</a>[T], mark *<a href="#Node">Node</a>[T])</pre>

MoveAfter moves node just after mark. Afterwards, mark.Next() == node && node.Prev() == mark.


## <a id="MoveBefore"></a><pre>func (l *<a href="#List">List</a>[T]) MoveBefore(node *<a href="#Node">Node</a>[T], mark *<a href="#Node">Node</a>[T])</pre>

MoveBefore moves node just before mark. Afterwards, mark.Prev() == node && node.Next() == mark.


## <a id="MoveToBack"></a><pre>func (l *<a href="#List">List</a>[T]) MoveToBack(node *<a href="#Node">Node</a>[T])</pre>

MoveToFront moves node to the back of the list.


## <a id="MoveToFront"></a><pre>func (l *<a href="#List">List</a>[T]) MoveToFront(node *<a href="#Node">Node</a>[T])</pre>

MoveToFront moves node to the front of the list.


## <a id="PushBack"></a><pre>func (l *<a href="#List">List</a>[T]) PushBack(value T) *<a href="#Node">Node</a>[T]</pre>

PushFront adds a new node with the given value to the back of the list.


## <a id="PushFront"></a><pre>func (l *<a href="#List">List</a>[T]) PushFront(value T) *<a href="#Node">Node</a>[T]</pre>

PushFront adds a new node with the given value to the front of the list.


## <a id="Remove"></a><pre>func (l *<a href="#List">List</a>[T]) Remove(node *<a href="#Node">Node</a>[T])</pre>

Remove removes node from the list.


## <a id="Node"></a><pre>type Node</pre>
```go
type Node[T any] struct {

	// Value is user-controlled, and never modified by this package.
	Value T
	// contains filtered or unexported fields
}
```

Node is a node in a linked-list.


## <a id="Next"></a><pre>func (n *<a href="#Node">Node</a>[T]) Next() *<a href="#Node">Node</a>[T]</pre>

Next returns the next node in the list that n is a part of, if there is one.


## <a id="Prev"></a><pre>func (n *<a href="#Node">Node</a>[T]) Prev() *<a href="#Node">Node</a>[T]</pre>

Prev returns the previous node in the list that n is a part of, if there is one.


