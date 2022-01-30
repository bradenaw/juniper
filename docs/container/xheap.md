# `package xheap`

```
import "github.com/bradenaw/juniper/container/xheap"
```

## Overview

Package xheap contains extensions to the standard library package container/heap.


## Index

<samp><a href="#Heap">type Heap</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#New">func New[T any](less xsort.Less[T], initial []T) Heap[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Grow">func (h *Heap[T]) Grow(n int)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (h *Heap[T]) Iterate() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (h *Heap[T]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Peek">func (h *Heap[T]) Peek() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Pop">func (h *Heap[T]) Pop() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Push">func (h *Heap[T]) Push(item T)</a></samp>

<samp><a href="#KP">type KP</a></samp>

<samp><a href="#PriorityQueue">type PriorityQueue</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewPriorityQueue">func NewPriorityQueue[K comparable, P any](
&nbsp;&nbsp;&nbsp;&nbsp;	less xsort.Less[P],
&nbsp;&nbsp;&nbsp;&nbsp;	initial []KP[K, P],
&nbsp;&nbsp;&nbsp;&nbsp;) PriorityQueue[K, P]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Contains">func (h *PriorityQueue[K, P]) Contains(k K) bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Grow">func (h *PriorityQueue[K, P]) Grow(n int)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (h *PriorityQueue[K, P]) Iterate() iterator.Iterator[K]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (h *PriorityQueue[K, P]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Peek">func (h *PriorityQueue[K, P]) Peek() K</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Pop">func (h *PriorityQueue[K, P]) Pop() K</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Priority">func (h *PriorityQueue[K, P]) Priority(k K) P</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Remove">func (h *PriorityQueue[K, P]) Remove(k K)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Update">func (h *PriorityQueue[K, P]) Update(k K, p P)</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

## Types

<h3><a id="Heap"></a><samp>type Heap</samp></h3>
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


<h3><a id="New"></a><samp>func New[T any](less <a href="../xsort.html#Less">xsort.Less</a>[T], initial []T) <a href="#Heap">Heap</a>[T]</samp></h3>

New returns a new Heap which uses less to determine the minimum element.

The elements from initial are added to the heap. initial is modified by New and utilized by the
Heap, so it should not be used after passing to New(). Passing initial is faster (O(n)) than
creating an empty heap and pushing each item (O(n * log(n))).


<h3><a id="Grow"></a><samp>func (h *<a href="#Heap">Heap</a>[T]) Grow(n int)</samp></h3>

Grow allocates sufficient space to add n more elements without needing to reallocate.


<h3><a id="Iterate"></a><samp>func (h *<a href="#Heap">Heap</a>[T]) Iterate() <a href="../iterator.html#Iterator">iterator.Iterator</a>[T]</samp></h3>

Iterate iterates over the elements of the heap.

The iterator panics if the heap has been modified since iteration started.


<h3><a id="Len"></a><samp>func (h *<a href="#Heap">Heap</a>[T]) Len() int</samp></h3>

Len returns the current number of elements in the heap.


<h3><a id="Peek"></a><samp>func (h *<a href="#Heap">Heap</a>[T]) Peek() T</samp></h3>

Peek returns the minimum item in the heap. It panics if h.Len()==0.


<h3><a id="Pop"></a><samp>func (h *<a href="#Heap">Heap</a>[T]) Pop() T</samp></h3>

Pop removes and returns the minimum item in the heap. It panics if h.Len()==0.


<h3><a id="Push"></a><samp>func (h *<a href="#Heap">Heap</a>[T]) Push(item T)</samp></h3>

Push adds item to the heap.


<h3><a id="KP"></a><samp>type KP</samp></h3>
```go
type KP[K any, P any] struct {
	K K
	P P
}
```

KP holds key and priority for PriorityQueue.


<h3><a id="PriorityQueue"></a><samp>type PriorityQueue</samp></h3>
```go
type PriorityQueue[K comparable, P any] struct {
	// contains filtered or unexported fields
}
```

PriorityQueue is a queue that yields items in increasing order of priority.


<h3><a id="NewPriorityQueue"></a><samp>func NewPriorityQueue[K comparable, P any](less <a href="../xsort.html#Less">xsort.Less</a>[P], initial []<a href="#KP">KP</a>[K, P]) <a href="#PriorityQueue">PriorityQueue</a>[K, P]</samp></h3>

NewPriorityQueue returns a new PriorityQueue which uses less to determine the minimum element.

The elements from initial are added to the priority queue. initial is modified by
NewPriorityQueue and utilized by the PriorityQueue, so it should not be used after passing to
NewPriorityQueue. Passing initial is faster (O(n)) than creating an empty priority queue and
pushing each item (O(n * log(n))).

Pop, Remove, and Update all take amoritized O(log(n)) time where n is the number of items in the
queue.

Len, Peek, Contains, and Priority take O(1) time.


<h3><a id="Contains"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Contains(k K) bool</samp></h3>

Contains returns true if the given key is present in the priority queue.


<h3><a id="Grow"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Grow(n int)</samp></h3>

Grow allocates sufficient space to add n more elements without needing to reallocate.


<h3><a id="Iterate"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Iterate() <a href="../iterator.html#Iterator">iterator.Iterator</a>[K]</samp></h3>

Iterate iterates over the elements of the priority queue.

The iterator panics if the priority queue has been modified since iteration started.


<h3><a id="Len"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Len() int</samp></h3>

Len returns the current number of elements in the priority queue.


<h3><a id="Peek"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Peek() K</samp></h3>

Peek returns the key of the lowest-P item in the priority queue. It panics if h.Len()==0.


<h3><a id="Pop"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Pop() K</samp></h3>

Pop removes and returns the lowest-P item in the priority queue. It panics if h.Len()==0.


<h3><a id="Priority"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Priority(k K) P</samp></h3>

Priority returns the priority of k, or the zero value of P if k is not present.


<h3><a id="Remove"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Remove(k K)</samp></h3>

Remove removes the item with the given key if present.


<h3><a id="Update"></a><samp>func (h *<a href="#PriorityQueue">PriorityQueue</a>[K, P]) Update(k K, p P)</samp></h3>

Update updates the priority of k to p, or adds it to the priority queue if not present.


