# `package xheap`

```
import "github.com/bradenaw/juniper/container/xheap"
```

# Overview

Package xheap contains extensions to the standard library package container/heap.


# Index

<pre><a href="#Heap">type Heap</a></pre>
<pre>    <a href="#New">func New[T any](less xsort.Less[T], initial []T) Heap[T]</a></pre>
<pre>    <a href="#Grow">func (h *Heap[T]) Grow(n int)</a></pre>
<pre>    <a href="#Iterate">func (h *Heap[T]) Iterate() iterator.Iterator[T]</a></pre>
<pre>    <a href="#Len">func (h *Heap[T]) Len() int</a></pre>
<pre>    <a href="#Peek">func (h *Heap[T]) Peek() T</a></pre>
<pre>    <a href="#Pop">func (h *Heap[T]) Pop() T</a></pre>
<pre>    <a href="#Push">func (h *Heap[T]) Push(item T)</a></pre>
<pre><a href="#KP">type KP</a></pre>
<pre><a href="#PriorityQueue">type PriorityQueue</a></pre>
<pre>    <a href="#NewPriorityQueue">func NewPriorityQueue[K comparable, P any](
    	less xsort.Less[P],
    	initial []KP[K, P],
    ) PriorityQueue[K, P]</a></pre>
<pre>    <a href="#Contains">func (h *PriorityQueue[K, P]) Contains(k K) bool</a></pre>
<pre>    <a href="#Grow">func (h *PriorityQueue[K, P]) Grow(n int)</a></pre>
<pre>    <a href="#Iterate">func (h *PriorityQueue[K, P]) Iterate() iterator.Iterator[K]</a></pre>
<pre>    <a href="#Len">func (h *PriorityQueue[K, P]) Len() int</a></pre>
<pre>    <a href="#Peek">func (h *PriorityQueue[K, P]) Peek() K</a></pre>
<pre>    <a href="#Pop">func (h *PriorityQueue[K, P]) Pop() K</a></pre>
<pre>    <a href="#Priority">func (h *PriorityQueue[K, P]) Priority(k K) P</a></pre>
<pre>    <a href="#Remove">func (h *PriorityQueue[K, P]) Remove(k K)</a></pre>
<pre>    <a href="#Update">func (h *PriorityQueue[K, P]) Update(k K, p P)</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

# Types

## <a id="Heap"></a><pre>type Heap</pre>
```go
type Heap[T any] struct {
	// contains filtered or unexported fields
}
```

Heap is a min-heap (https://en.wikipedia.org/wiki/Binary_heap). Min-heaps are a collection
structure that provide constant-time access to the minimum element, and logarithmic-time removal.
They are most commonly used as a priority queue.

Push and Pop take amoritized O(log(n)) time where n is the number of items in the heap.

Len and Peek take O(1) time.


## <a id="New"></a><pre>func New[T any](less <a href="../xsort.md#Less">xsort.Less</a>[T], initial []T) <a href="#Heap">Heap</a>[T]</pre>

New returns a new Heap which uses less to determine the minimum element.

The elements from initial are added to the heap. initial is modified by New and utilized by the
Heap, so it should not be used after passing to New(). Passing initial is faster (O(n)) than
creating an empty heap and pushing each item (O(n * log(n))).


## <a id="Grow"></a><pre>func (h *<a href="#Heap">Heap</a>[T]) Grow(n int)</pre>

Grow allocates sufficient space to add n more elements without needing to reallocate.


## <a id="Iterate"></a><pre>func (h *<a href="#Heap">Heap</a>[T]) Iterate() <a href="../iterator.md#Iterator">iterator.Iterator</a>[T]</pre>

Iterate iterates over the elements of the heap.

The iterator panics if the heap has been modified since iteration started.


## <a id="Len"></a><pre>func (h *<a href="#Heap">Heap</a>[T]) Len() int</pre>

Len returns the current number of elements in the heap.


## <a id="Peek"></a><pre>func (h *<a href="#Heap">Heap</a>[T]) Peek() T</pre>

Peek returns the minimum item in the heap. It panics if h.Len()==0.


## <a id="Pop"></a><pre>func (h *<a href="#Heap">Heap</a>[T]) Pop() T</pre>

Pop removes and returns the minimum item in the heap. It panics if h.Len()==0.


## <a id="Push"></a><pre>func (h *<a href="#Heap">Heap</a>[T]) Push(item T)</pre>

Push adds item to the heap.


## <a id="KP"></a><pre>type KP</pre>
```go
type KP[K any, P any] struct {
	K K
	P P
}
```

KP holds key and priority for PriorityQueue.


## <a id="PriorityQueue"></a><pre>type PriorityQueue</pre>
```go
type PriorityQueue[K comparable, P any] struct {
	// contains filtered or unexported fields
}
```

PriorityQueue is a queue that yields items in increasing order of priority.


## <a id="NewPriorityQueue"></a><pre>func NewPriorityQueue[K comparable, P any](less <a href="../xsort.md#Less">xsort.Less</a>[P], initial []<a href="#KP">KP</a>[K, P]) <a href="#PriorityQueue">PriorityQueue</a>[K, P]</pre>

NewPriorityQueue returns a new PriorityQueue which uses less to determine the minimum element.

The elements from initial are added to the priority queue. initial is modified by
NewPriorityQueue and utilized by the PriorityQueue, so it should not be used after passing to
NewPriorityQueue. Passing initial is faster (O(n)) than creating an empty priority queue and
pushing each item (O(n * log(n))).

Pop, Remove, and Update all take amoritized O(log(n)) time where n is the number of items in the
queue.

Len, Peek, Contains, and Priority take O(1) time.


## <a id="Contains"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Contains(k K) bool</pre>

Contains returns true if the given key is present in the priority queue.


## <a id="Grow"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Grow(n int)</pre>

Grow allocates sufficient space to add n more elements without needing to reallocate.


## <a id="Iterate"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Iterate() <a href="../iterator.md#Iterator">iterator.Iterator</a>[K]</pre>

Iterate iterates over the elements of the priority queue.

The iterator panics if the priority queue has been modified since iteration started.


## <a id="Len"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Len() int</pre>

Len returns the current number of elements in the priority queue.


## <a id="Peek"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Peek() K</pre>

Peek returns the key of the lowest-P item in the priority queue. It panics if h.Len()==0.


## <a id="Pop"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Pop() K</pre>

Pop removes and returns the lowest-P item in the priority queue. It panics if h.Len()==0.


## <a id="Priority"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Priority(k K) P</pre>

Priority returns the priority of k, or the zero value of P if k is not present.


## <a id="Remove"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Remove(k K)</pre>

Remove removes the item with the given key if present.


## <a id="Update"></a><pre>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Update(k K, p P)</pre>

Update updates the priority of k to p, or adds it to the priority queue if not present.


