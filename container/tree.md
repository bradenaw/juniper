# `package tree`

```
import "github.com/bradenaw/juniper/container/tree"
```

## Overview

Package tree contains an implementation of a B-tree Map and Set. These are similar to Go's map
built-in, but keep elements in sorted order.


## Index

<samp><a href="#Bound">type Bound</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Excluded">func Excluded[K any](key K) Bound[K]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Included">func Included[K any](key K) Bound[K]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Unbounded">func Unbounded[K any]() Bound[K]</a></samp>

<samp><a href="#KVPair">type KVPair</a></samp>

<samp><a href="#Map">type Map</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewMap">func NewMap[K any, V any](less xsort.Less[K]) Map[K, V]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Contains">func (m Map[K, V]) Contains(k K) bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Delete">func (m Map[K, V]) Delete(k K)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#First">func (m Map[K, V]) First() (K, V)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Get">func (m Map[K, V]) Get(k K) V</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Last">func (m Map[K, V]) Last() (K, V)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (m Map[K, V]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Put">func (m Map[K, V]) Put(k K, v V)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Range">func (m Map[K, V]) Range(lower Bound[K], upper Bound[K]) iterator.Iterator[KVPair[K, V]]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#RangeReverse">func (m Map[K, V]) RangeReverse(lower Bound[K], upper Bound[K]) iterator.Iterator[KVPair[K, V]]</a></samp>

<samp><a href="#Set">type Set</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewSet">func NewSet[T any](less xsort.Less[T]) Set[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Add">func (s Set[T]) Add(item T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Contains">func (s Set[T]) Contains(item T) bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#First">func (s Set[T]) First() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (s Set[T]) Iterate() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Last">func (s Set[T]) Last() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (s Set[T]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Range">func (s Set[T]) Range(lower Bound[T], upper Bound[T]) iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#RangeReverse">func (s Set[T]) RangeReverse(lower Bound[T], upper Bound[T]) iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Remove">func (s Set[T]) Remove(item T)</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

This section is empty.

## Types

<h3><a id="Bound"></a><samp>type Bound</samp></h3>
```go
type Bound[K any] struct {
	// contains filtered or unexported fields
}
```

Bound is an endpoint for a range.


<h3><a id="Excluded"></a><samp>func Excluded[K any](key K) <a href="#Bound">Bound</a>[K]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L96">src</a></small></sub></h3>

Included returns a Bound that goes up to but not including key.


<h3><a id="Included"></a><samp>func Included[K any](key K) <a href="#Bound">Bound</a>[K]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L93">src</a></small></sub></h3>

Included returns a Bound that goes up to and including key.


<h3><a id="Unbounded"></a><samp>func Unbounded[K any]() <a href="#Bound">Bound</a>[K]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L99">src</a></small></sub></h3>

Included returns a Bound at the end of the collection.


<h3><a id="KVPair"></a><samp>type KVPair</samp></h3>
```go
type KVPair[K any, V any] struct {
	Key   K
	Value V
}
```



<h3><a id="Map"></a><samp>type Map</samp></h3>
```go
type Map[K any, V any] struct {
	// contains filtered or unexported fields
}
```

Map is a tree-structured key-value map, similar to Go's built-in map but keeps elements in sorted
order by key.

It is safe for multiple goroutines to Put concurrently with keys that are already in the map.


<h3><a id="NewMap"></a><samp>func NewMap[K any, V any](less <a href="../xsort.html#Less">xsort.Less</a>[K]) <a href="#Map">Map</a>[K, V]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L25">src</a></small></sub></h3>

NewMap returns a Map that uses less to determine the sort order of keys. If !less(a, b) &&
!less(b, a), then a and b are considered the same key. The output of less must not change for any
pair of keys while they are in the map.


<h3><a id="Contains"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Contains(k K) bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L54">src</a></small></sub></h3>

Contains returns true if the given key is present in the map.


<h3><a id="Delete"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Delete(k K)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L43">src</a></small></sub></h3>

Delete removes the given key from the map.


<h3><a id="First"></a><samp>func (m <a href="#Map">Map</a>[K, V]) First() (K, V)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L59">src</a></small></sub></h3>

First returns the lowest-keyed entry in the map according to less.


<h3><a id="Get"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Get(k K) V</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L49">src</a></small></sub></h3>

Get returns the value associated with the given key if it is present in the map. Otherwise, it
returns the zero-value of V.


<h3><a id="Iterate"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Iterate() <a href="../iterator.html#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L74">src</a></small></sub></h3>

Iterate returns an iterator that yields the elements of the map in ascending order by key.

The map may be safely modified during iteration and the iterator will continue from the
next-lowest key. Thus the iterator will see new elements that are after the current position of
the iterator according to less, but will not necessarily see a consistent snapshot of the state
of the map.


<h3><a id="Last"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Last() (K, V)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L64">src</a></small></sub></h3>

Last returns the highest-keyed entry in the map according to less.


<h3><a id="Len"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Len() int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L32">src</a></small></sub></h3>

Len returns the number of elements in the map.


<h3><a id="Put"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Put(k K, v V)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L38">src</a></small></sub></h3>

Put inserts the key-value pair into the map, overwriting the value for the key if it already
exists.


<h3><a id="Range"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Range(lower <a href="#Bound">Bound</a>[K], upper <a href="#Bound">Bound</a>[K]) <a href="../iterator.html#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L108">src</a></small></sub></h3>

Range returns an iterator that yields the elements of the map between the given bounds in
ascending order by key.

The map may be safely modified during iteration and the iterator will continue from the
next-lowest key. Thus the iterator will see new elements that are after the current position of
the iterator according to less, but will not necessarily see a consistent snapshot of the state
of the map.


<h3><a id="RangeReverse"></a><samp>func (m <a href="#Map">Map</a>[K, V]) RangeReverse(lower <a href="#Bound">Bound</a>[K], upper <a href="#Bound">Bound</a>[K]) <a href="../iterator.html#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/map.go#L119">src</a></small></sub></h3>

RangeReverse returns an iterator that yields the elements of the map between the given bounds in
descending order by key.

The map may be safely modified during iteration and the iterator will continue from the
next-lowest key. Thus the iterator will see new elements that are after the current position of
the iterator according to less, but will not necessarily see a consistent snapshot of the state
of the map.


<h3><a id="Set"></a><samp>type Set</samp></h3>
```go
type Set[T any] struct {
	// contains filtered or unexported fields
}
```

Set is a tree-structured set. Sets are a collection of unique elements. Similar to Go's built-in
map[T]struct{} but keeps elements in sorted order.


<h3><a id="NewSet"></a><samp>func NewSet[T any](less <a href="../xsort.html#Less">xsort.Less</a>[T]) <a href="#Set">Set</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L18">src</a></small></sub></h3>

NewSet returns a Set that uses less to determine the sort order of items. If !less(a, b) &&
!less(b, a), then a and b are considered the same item. The output of less must not change for
any pair of items while they are in the set.


<h3><a id="Add"></a><samp>func (s <a href="#Set">Set</a>[T]) Add(item T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L30">src</a></small></sub></h3>

Add adds item to the set if it is not already present.


<h3><a id="Contains"></a><samp>func (s <a href="#Set">Set</a>[T]) Contains(item T) bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L40">src</a></small></sub></h3>

Contains returns true if item is present in the set.


<h3><a id="First"></a><samp>func (s <a href="#Set">Set</a>[T]) First() T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L45">src</a></small></sub></h3>

First returns the lowest item in the set according to less.


<h3><a id="Iterate"></a><samp>func (s <a href="#Set">Set</a>[T]) Iterate() <a href="../iterator.html#Iterator">iterator.Iterator</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L62">src</a></small></sub></h3>

Iterate returns an iterator that yields the elements of the set in ascending order.

The set may be safely modified during iteration and the iterator will continue from the
next-lowest item. Thus the iterator will see new items that are after the current position
of the iterator according to less, but will not necessarily see a consistent snapshot of the
state of the set.


<h3><a id="Last"></a><samp>func (s <a href="#Set">Set</a>[T]) Last() T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L51">src</a></small></sub></h3>

Last returns the highest item in the set according to less.


<h3><a id="Len"></a><samp>func (s <a href="#Set">Set</a>[T]) Len() int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L25">src</a></small></sub></h3>

Len returns the number of elements in the set.


<h3><a id="Range"></a><samp>func (s <a href="#Set">Set</a>[T]) Range(lower <a href="#Bound">Bound</a>[T], upper <a href="#Bound">Bound</a>[T]) <a href="../iterator.html#Iterator">iterator.Iterator</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L73">src</a></small></sub></h3>

Range returns an iterator that yields the elements of the set between the given bounds in
ascending order.

The set may be safely modified during iteration and the iterator will continue from the
next-lowest item. Thus the iterator will see new items that are after the current position
of the iterator according to less, but will not necessarily see a consistent snapshot of the
state of the set.


<h3><a id="RangeReverse"></a><samp>func (s <a href="#Set">Set</a>[T]) RangeReverse(lower <a href="#Bound">Bound</a>[T], upper <a href="#Bound">Bound</a>[T]) <a href="../iterator.html#Iterator">iterator.Iterator</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L86">src</a></small></sub></h3>

RangeReverse returns an iterator that yields the elements of the set between the given bounds in
descending order.

The set may be safely modified during iteration and the iterator will continue from the
next-lowest item. Thus the iterator will see new items that are after the current position
of the iterator according to less, but will not necessarily see a consistent snapshot of the
state of the set.


<h3><a id="Remove"></a><samp>func (s <a href="#Set">Set</a>[T]) Remove(item T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/container/tree/set.go#L35">src</a></small></sub></h3>

Remove removes item from the set if it is present, and does nothing otherwise.


