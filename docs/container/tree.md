# `package tree`

```
import "github.com/bradenaw/juniper/container/tree"
```

## Overview

Package tree contains an implementation of a B-tree Map and Set. These are similar to Go's map
built-in, but keep elements in sorted order.


## Index

<samp><a href="#KVPair">type KVPair</a></samp>

<samp><a href="#Map">type Map</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewMap">func NewMap[K any, V any](less xsort.Less[K]) Map[K, V]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Contains">func (m Map[K, V]) Contains(k K) bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Cursor">func (m Map[K, V]) Cursor() *MapCursor[K, V]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Delete">func (m Map[K, V]) Delete(k K)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#First">func (m Map[K, V]) First() (K, V)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Get">func (m Map[K, V]) Get(k K) V</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Last">func (m Map[K, V]) Last() (K, V)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (m Map[K, V]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Put">func (m Map[K, V]) Put(k K, v V)</a></samp>

<samp><a href="#MapCursor">type MapCursor</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Backward">func (c *MapCursor[K, V]) Backward() iterator.Iterator[KVPair[K, V]]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Forward">func (c *MapCursor[K, V]) Forward() iterator.Iterator[KVPair[K, V]]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Key">func (c *MapCursor[K, V]) Key() K</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Next">func (c *MapCursor[K, V]) Next()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Ok">func (c *MapCursor[K, V]) Ok() bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Prev">func (c *MapCursor[K, V]) Prev()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekFirst">func (c *MapCursor[K, V]) SeekFirst()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekFirstGreater">func (c *MapCursor[K, V]) SeekFirstGreater(k K)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekFirstGreaterOrEqual">func (c *MapCursor[K, V]) SeekFirstGreaterOrEqual(k K)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekLast">func (c *MapCursor[K, V]) SeekLast()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekLastLess">func (c *MapCursor[K, V]) SeekLastLess(k K)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekLastLessOrEqual">func (c *MapCursor[K, V]) SeekLastLessOrEqual(k K)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Value">func (c *MapCursor[K, V]) Value() V</a></samp>

<samp><a href="#Set">type Set</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewSet">func NewSet[T any](less xsort.Less[T]) Set[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Add">func (s Set[T]) Add(item T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Contains">func (s Set[T]) Contains(item T) bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Cursor">func (s Set[T]) Cursor() *SetCursor[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#First">func (s Set[T]) First() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (s Set[T]) Iterate() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Last">func (s Set[T]) Last() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (s Set[T]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Remove">func (s Set[T]) Remove(item T)</a></samp>

<samp><a href="#SetCursor">type SetCursor</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Backward">func (c *SetCursor[T]) Backward() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Forward">func (c *SetCursor[T]) Forward() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Item">func (c *SetCursor[T]) Item() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Next">func (c *SetCursor[T]) Next()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Ok">func (c *SetCursor[T]) Ok() bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Prev">func (c *SetCursor[T]) Prev()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekFirst">func (c *SetCursor[T]) SeekFirst()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekFirstGreater">func (c *SetCursor[T]) SeekFirstGreater(x T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekFirstGreaterOrEqual">func (c *SetCursor[T]) SeekFirstGreaterOrEqual(x T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekLast">func (c *SetCursor[T]) SeekLast()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekLastLess">func (c *SetCursor[T]) SeekLastLess(x T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#SeekLastLessOrEqual">func (c *SetCursor[T]) SeekLastLessOrEqual(x T)</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

## Types

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


<h3><a id="NewMap"></a><samp>func NewMap[K any, V any](less <a href="../xsort.html#Less">xsort.Less</a>[K]) <a href="#Map">Map</a>[K, V]</samp></h3>

NewMap returns a Map that uses less to determine the sort order of keys. If !less(a, b) &&
!less(b, a), then a and b are considered the same key. The output of less must not change for any
pair of keys while they are in the map.


<h3><a id="Contains"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Contains(k K) bool</samp></h3>

Contains returns true if the given key is present in the map.


<h3><a id="Cursor"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Cursor() *<a href="#MapCursor">MapCursor</a>[K, V]</samp></h3>

Cursor returns a cursor into the map placed at the first element.


<h3><a id="Delete"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Delete(k K)</samp></h3>

Delete removes the given key from the map.


<h3><a id="First"></a><samp>func (m <a href="#Map">Map</a>[K, V]) First() (K, V)</samp></h3>

First returns the lowest-keyed entry in the map according to less.


<h3><a id="Get"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Get(k K) V</samp></h3>

Get returns the value associated with the given key if it is present in the map. Otherwise, it
returns the zero-value of V.


<h3><a id="Iterate"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Iterate() <a href="../iterator.html#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</samp></h3>

Iterate returns an iterator that yields the elements of the map in sorted order by key.

The map may be safely modified during iteration and the iterator will continue from the
next-lowest key. Thus if the map is modified, the iterator will not necessarily return all of
the keys present in the map.


<h3><a id="Last"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Last() (K, V)</samp></h3>

Last returns the highest-keyed entry in the map according to less.


<h3><a id="Len"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Len() int</samp></h3>

Len returns the number of elements in the map.


<h3><a id="Put"></a><samp>func (m <a href="#Map">Map</a>[K, V]) Put(k K, v V)</samp></h3>

Put inserts the key-value pair into the map, overwriting the value for the key if it already
exists.


<h3><a id="MapCursor"></a><samp>type MapCursor</samp></h3>
```go
type MapCursor[K any, V any] struct {
	// contains filtered or unexported fields
}
```

MapCursor is a cursor into a Map.

A cursor is usable while a map is being modified. If the element the cursor is at is deleted, the
cursor will still return the old value.


<h3><a id="Backward"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Backward() <a href="../iterator.html#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</samp></h3>

Backward returns an iterator that starts from the cursor's position and yields all of the
elements less than or equal to the cursor in descending order.

This iterator's Next method is amoritized O(1), unless the map changes in which case the
following Next is O(log(n)) where n is the number of elements in the map.


<h3><a id="Forward"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Forward() <a href="../iterator.html#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</samp></h3>

Forward returns an iterator that starts from the cursor's position and yields all of the elements
greater than or equal to the cursor in ascending order.

This iterator's Next method is amoritized O(1), unless the map changes in which case the
following Next is O(log(n)) where n is the number of elements in the map.


<h3><a id="Key"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Key() K</samp></h3>

Key returns the key of the element that the cursor is at. Panics if Ok is false.


<h3><a id="Next"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Next()</samp></h3>

Next moves the cursor to the next element in the map.

Next is amoritized O(1) unless the map has been modified since the last cursor move, in which
case it's O(log(n)).


<h3><a id="Ok"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Ok() bool</samp></h3>

Ok returns false if the cursor is not currently placed at an element, for example if Next
advances past the last element.


<h3><a id="Prev"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Prev()</samp></h3>

Prev moves the cursor to the previous element in the map.

Prev is amoritized O(1) unless the map has been modified since the last cursor move, in which
case it's O(log(n)).


<h3><a id="SeekFirst"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekFirst()</samp></h3>

SeekFirst moves the cursor to the first element in the map.

SeekFirst is O(log(n)).


<h3><a id="SeekFirstGreater"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekFirstGreater(k K)</samp></h3>

SeekFirstGreater moves the cursor to the element in the map just after k.

SeekFirstGreater is O(log(n)).


<h3><a id="SeekFirstGreaterOrEqual"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekFirstGreaterOrEqual(k K)</samp></h3>

SeekFirstGreaterOrEqual moves the cursor to the element in the map with the least key that is
greater than or equal to k.

SeetFirstGreaterOrEqual is O(log(n)).


<h3><a id="SeekLast"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekLast()</samp></h3>

SeekLast moves the cursor to the last element in the map.

SeekLast is O(log(n)).


<h3><a id="SeekLastLess"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekLastLess(k K)</samp></h3>

SeekLastLess moves the cursor to the element in the map just before k.

SeekLastLess is O(log(n)).


<h3><a id="SeekLastLessOrEqual"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekLastLessOrEqual(k K)</samp></h3>

SeekLastLessOrEqual moves the cursor to the element in the map with the greatest key that is less
than or equal to k.

SeekLastLessOrEqual is O(log(n)).


<h3><a id="Value"></a><samp>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Value() V</samp></h3>

Value returns the value of the element that the cursor is at. Panics if Ok is false.


<h3><a id="Set"></a><samp>type Set</samp></h3>
```go
type Set[T any] struct {
	// contains filtered or unexported fields
}
```

Set is a tree-structured set. Sets are a collection of unique elements. Similar to Go's built-in
map[T]struct{} but keeps elements in sorted order.


<h3><a id="NewSet"></a><samp>func NewSet[T any](less <a href="../xsort.html#Less">xsort.Less</a>[T]) <a href="#Set">Set</a>[T]</samp></h3>

NewSet returns a Set that uses less to determine the sort order of items. If !less(a, b) &&
!less(b, a), then a and b are considered the same item. The output of less must not change for
any pair of items while they are in the set.


<h3><a id="Add"></a><samp>func (s <a href="#Set">Set</a>[T]) Add(item T)</samp></h3>

Add adds item to the set if it is not already present.


<h3><a id="Contains"></a><samp>func (s <a href="#Set">Set</a>[T]) Contains(item T) bool</samp></h3>

Contains returns true if item is present in the set.


<h3><a id="Cursor"></a><samp>func (s <a href="#Set">Set</a>[T]) Cursor() *<a href="#SetCursor">SetCursor</a>[T]</samp></h3>

Cursor returns a cursor into the set placed at the first item.


<h3><a id="First"></a><samp>func (s <a href="#Set">Set</a>[T]) First() T</samp></h3>

First returns the lowest item in the set according to less.


<h3><a id="Iterate"></a><samp>func (s <a href="#Set">Set</a>[T]) Iterate() <a href="../iterator.html#Iterator">iterator.Iterator</a>[T]</samp></h3>

Iterate returns an iterator that yields the elements of the set in sorted order.

The set may be safely modified during iteration and the iterator will continue from the
next-lowest item. Thus if the set is modified, the iterator will not necessarily return all of
the items present in the set.


<h3><a id="Last"></a><samp>func (s <a href="#Set">Set</a>[T]) Last() T</samp></h3>

Last returns the highest item in the set according to less.


<h3><a id="Len"></a><samp>func (s <a href="#Set">Set</a>[T]) Len() int</samp></h3>

Len returns the number of elements in the set.


<h3><a id="Remove"></a><samp>func (s <a href="#Set">Set</a>[T]) Remove(item T)</samp></h3>

Remove removes item from the set if it is present, and does nothing otherwise.


<h3><a id="SetCursor"></a><samp>type SetCursor</samp></h3>
```go
type SetCursor[T any] struct {
	// contains filtered or unexported fields
}
```

SetCursor is a cursor into a Set.

A cursor is usable while a set is being modified. If the item the cursor is at is deleted, the
cursor will still return the old item.


<h3><a id="Backward"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) Backward() <a href="../iterator.html#Iterator">iterator.Iterator</a>[T]</samp></h3>

Backward returns an iterator that starts from the cursor's position and yields all of the
elements less than or equal to the cursor in descending order.

This iterator's Next method is amoritized O(1), unless the map changes in which case the
following Next is O(log(n)) where n is the number of elements in the map.


<h3><a id="Forward"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) Forward() <a href="../iterator.html#Iterator">iterator.Iterator</a>[T]</samp></h3>

Forward returns an iterator that starts from the cursor's position and yields all of the elements
greater than or equal to the cursor in ascending order.

This iterator's Next method is amoritized O(1), unless the map changes in which case the
following Next is O(log(n)) where n is the number of elements in the map.


<h3><a id="Item"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) Item() T</samp></h3>

Item returns the item that the cursor is at. Panics if Ok is false.


<h3><a id="Next"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) Next()</samp></h3>

Next moves the cursor to the next item in the set.

Next is amoritized O(1) unless the map has been modified since the last cursor move, in which
case it's O(log(n)).


<h3><a id="Ok"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) Ok() bool</samp></h3>

Ok returns false if the cursor is not currently placed at an item, for example if Next advances
past the last item.


<h3><a id="Prev"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) Prev()</samp></h3>

Prev moves the cursor to the previous item in the set.

Prev is amoritized O(1) unless the map has been modified since the last cursor move, in which
case it's O(log(n)).


<h3><a id="SeekFirst"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekFirst()</samp></h3>

SeekFirst moves the cursor to the first item in the set.

SeekFirst is O(log(n)).


<h3><a id="SeekFirstGreater"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekFirstGreater(x T)</samp></h3>

SeekFirstGreater moves the cursor to the item in the set just after x.

SeekFirstGreater is O(log(n)).


<h3><a id="SeekFirstGreaterOrEqual"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekFirstGreaterOrEqual(x T)</samp></h3>

SeekFirstGreaterOrEqual moves the cursor to the least item in the set that is greater than or
equal to x.

SeetFirstGreaterOrEqual is O(log(n)).


<h3><a id="SeekLast"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekLast()</samp></h3>

SeekLast moves the cursor to the last item in the set.

SeekLast is O(log(n)).


<h3><a id="SeekLastLess"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekLastLess(x T)</samp></h3>

SeekLastLess moves the cursor to the item in the set just before x.

SeekLastLess is O(log(n)).


<h3><a id="SeekLastLessOrEqual"></a><samp>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekLastLessOrEqual(x T)</samp></h3>

SeekLastLessOrEqual moves the cursor to the greatest item in the set that is less than or equal
to x.

SeekLastLessOrEqual is O(log(n)).


