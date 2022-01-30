# `package deque`

```
import "github.com/bradenaw/juniper/container/deque"
```

# Overview

Package deque contains a double-ended queue.


# Index

<pre><a href="#Deque">type Deque</a></pre>
<pre>    <a href="#Back">func (r *Deque[T]) Back() T</a></pre>
<pre>    <a href="#Front">func (r *Deque[T]) Front() T</a></pre>
<pre>    <a href="#Grow">func (r *Deque[T]) Grow(n int)</a></pre>
<pre>    <a href="#Item">func (r *Deque[T]) Item(i int) T</a></pre>
<pre>    <a href="#Iterate">func (r *Deque[T]) Iterate() iterator.Iterator[T]</a></pre>
<pre>    <a href="#Len">func (r *Deque[T]) Len() int</a></pre>
<pre>    <a href="#PopBack">func (r *Deque[T]) PopBack() T</a></pre>
<pre>    <a href="#PopFront">func (r *Deque[T]) PopFront() T</a></pre>
<pre>    <a href="#PushBack">func (r *Deque[T]) PushBack(item T)</a></pre>
<pre>    <a href="#PushFront">func (r *Deque[T]) PushFront(item T)</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

# Types

## <a id="Deque"></a><pre>type Deque</pre>
```go
type Deque[T any] struct {
	// contains filtered or unexported fields
}
```

Deque is a double-ended queue, allowing push and pop to both the front and back of the queue.
Pushes and pops are amoritized O(1). The zero-value is ready to use. Deque should not be copied
after first use.


<h2><a id="Back"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) Back() T</pre></h2>

Back returns the item at the back of the deque. It panics if the deque is empty.


<h2><a id="Front"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) Front() T</pre></h2>

Front returns the item at the front of the deque. It panics if the deque is empty.


<h2><a id="Grow"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) Grow(n int)</pre></h2>

Grow allocates sufficient space to add n more items without needing to reallocate.


<h2><a id="Item"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) Item(i int) T</pre></h2>

Item returns the ith item in the deque. 0 is the front and r.Len()-1 is the back.


<h2><a id="Iterate"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) Iterate() <a href="../iterator.md#Iterator">iterator.Iterator</a>[T]</pre></h2>

Iterate iterates over the elements of the deque.

The iterator panics if the deque has been modified since iteration started.


<h2><a id="Len"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) Len() int</pre></h2>

Len returns the number of items in the deque.


<h2><a id="PopBack"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) PopBack() T</pre></h2>

PopBack removes and returns the item at the back of the deque. It panics if the deque is empty.


<h2><a id="PopFront"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) PopFront() T</pre></h2>

PopFront removes and returns the item at the front of the deque. It panics if the deque is empty.


<h2><a id="PushBack"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) PushBack(item T)</pre></h2>

PushFront adds item to the back of the deque.


<h2><a id="PushFront"></a><pre>func (r *<a href="#Deque">Deque</a>[T]) PushFront(item T)</pre></h2>

PushFront adds item to the front of the deque.


