# `package deque`

```
import "github.com/bradenaw/juniper/container/deque"
```

# Overview

Package deque contains a double-ended queue.


# Index

<samp><a href="#Deque">type Deque</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Back">func (r *Deque[T]) Back() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Front">func (r *Deque[T]) Front() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Grow">func (r *Deque[T]) Grow(n int)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Item">func (r *Deque[T]) Item(i int) T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (r *Deque[T]) Iterate() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (r *Deque[T]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PopBack">func (r *Deque[T]) PopBack() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PopFront">func (r *Deque[T]) PopFront() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PushBack">func (r *Deque[T]) PushBack(item T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PushFront">func (r *Deque[T]) PushFront(item T)</a></samp>


# Constants

This section is empty.

# Variables

This section is empty.

# Functions

# Types

<h2><a id="Deque"></a><samp>type Deque</samp></h2>
```go
type Deque[T any] struct {
	// contains filtered or unexported fields
}
```

Deque is a double-ended queue, allowing push and pop to both the front and back of the queue.
Pushes and pops are amoritized O(1). The zero-value is ready to use. Deque should not be copied
after first use.


<h2><a id="Back"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) Back() T</samp></h2>

Back returns the item at the back of the deque. It panics if the deque is empty.


<h2><a id="Front"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) Front() T</samp></h2>

Front returns the item at the front of the deque. It panics if the deque is empty.


<h2><a id="Grow"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) Grow(n int)</samp></h2>

Grow allocates sufficient space to add n more items without needing to reallocate.


<h2><a id="Item"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) Item(i int) T</samp></h2>

Item returns the ith item in the deque. 0 is the front and r.Len()-1 is the back.


<h2><a id="Iterate"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) Iterate() <a href="../iterator.md#Iterator">iterator.Iterator</a>[T]</samp></h2>

Iterate iterates over the elements of the deque.

The iterator panics if the deque has been modified since iteration started.


<h2><a id="Len"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) Len() int</samp></h2>

Len returns the number of items in the deque.


<h2><a id="PopBack"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) PopBack() T</samp></h2>

PopBack removes and returns the item at the back of the deque. It panics if the deque is empty.


<h2><a id="PopFront"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) PopFront() T</samp></h2>

PopFront removes and returns the item at the front of the deque. It panics if the deque is empty.


<h2><a id="PushBack"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) PushBack(item T)</samp></h2>

PushFront adds item to the back of the deque.


<h2><a id="PushFront"></a><samp>func (r *<a href="#Deque">Deque</a>[T]) PushFront(item T)</samp></h2>

PushFront adds item to the front of the deque.


