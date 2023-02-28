# `package deque`

```
import "github.com/bradenaw/juniper/container/deque"
```

## Overview

Package deque contains a double-ended queue.


## Index

<samp><a href="#Deque">type Deque</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Back">func (d *Deque[T]) Back() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Front">func (d *Deque[T]) Front() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Grow">func (d *Deque[T]) Grow(n int)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Item">func (d *Deque[T]) Item(i int) T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (d *Deque[T]) Iterate() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (d *Deque[T]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PopBack">func (d *Deque[T]) PopBack() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PopFront">func (d *Deque[T]) PopFront() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PushBack">func (d *Deque[T]) PushBack(item T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PushFront">func (d *Deque[T]) PushFront(item T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Shrink">func (d *Deque[T]) Shrink(n int)</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

This section is empty.

## Types

<h3><a id="Deque"></a><samp>type Deque</samp></h3>
```go
type Deque[T any] struct {
	// contains filtered or unexported fields
}
```

Deque is a double-ended queue, allowing push and pop to both the front and back of the queue.
Pushes and pops are amortized O(1). The zero-value is ready to use. Deque should not be copied
after first use.


<h3><a id="Back"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) Back() T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L165">src</a></small></sub></h3>

Back returns the item at the back of the deque. It panics if the deque is empty.


<h3><a id="Front"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) Front() T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L157">src</a></small></sub></h3>

Front returns the item at the front of the deque. It panics if the deque is empty.


<h3><a id="Grow"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) Grow(n int)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L47">src</a></small></sub></h3>

Grow allocates sufficient space to add n more items without needing to reallocate.


<h3><a id="Item"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) Item(i int) T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L170">src</a></small></sub></h3>

Item returns the ith item in the deque. 0 is the front and d.Len()-1 is the back.


<h3><a id="Iterate"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) Iterate() <a href="../iterator.html#Iterator">iterator.Iterator</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L215">src</a></small></sub></h3>

Iterate iterates over the elements of the deque.

The iterator panics if the deque has been modified since iteration started.


<h3><a id="Len"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) Len() int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L35">src</a></small></sub></h3>

Len returns the number of items in the deque.


<h3><a id="PopBack"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) PopBack() T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L137">src</a></small></sub></h3>

PopBack removes and returns the item at the back of the deque. It panics if the deque is empty.


<h3><a id="PopFront"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) PopFront() T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L117">src</a></small></sub></h3>

PopFront removes and returns the item at the front of the deque. It panics if the deque is empty.


<h3><a id="PushBack"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) PushBack(item T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L80">src</a></small></sub></h3>

PushFront adds item to the back of the deque.


<h3><a id="PushFront"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) PushFront(item T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L66">src</a></small></sub></h3>

PushFront adds item to the front of the deque.


<h3><a id="Shrink"></a><samp>func (d *<a href="#Deque">Deque</a>[T]) Shrink(n int)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/deque/deque.go#L56">src</a></small></sub></h3>

Shrink reallocates the backing buffer for d, if necessary, so that it fits only the current size
plus at most n extra items.


