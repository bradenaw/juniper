# `package tree`

```
import "github.com/bradenaw/juniper/container/tree"
```

# Overview

Package tree contains an implementation of a B-tree Map and Set. These are similar to Go's map
built-in, but keep elements in sorted order.


# Index

<pre><a href="#KVPair">type KVPair</a></pre>
<pre><a href="#Map">type Map</a></pre>
<pre>    <a href="#NewMap">func NewMap[K any, V any](less xsort.Less[K]) Map[K, V]</a></pre>
<pre>    <a href="#Contains">func (m Map[K, V]) Contains(k K) bool</a></pre>
<pre>    <a href="#Cursor">func (m Map[K, V]) Cursor() *MapCursor[K, V]</a></pre>
<pre>    <a href="#Delete">func (m Map[K, V]) Delete(k K)</a></pre>
<pre>    <a href="#First">func (m Map[K, V]) First() (K, V)</a></pre>
<pre>    <a href="#Get">func (m Map[K, V]) Get(k K) V</a></pre>
<pre>    <a href="#Iterate">func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]]</a></pre>
<pre>    <a href="#Last">func (m Map[K, V]) Last() (K, V)</a></pre>
<pre>    <a href="#Len">func (m Map[K, V]) Len() int</a></pre>
<pre>    <a href="#Put">func (m Map[K, V]) Put(k K, v V)</a></pre>
<pre><a href="#MapCursor">type MapCursor</a></pre>
<pre>    <a href="#Backward">func (c *MapCursor[K, V]) Backward() iterator.Iterator[KVPair[K, V]]</a></pre>
<pre>    <a href="#Forward">func (c *MapCursor[K, V]) Forward() iterator.Iterator[KVPair[K, V]]</a></pre>
<pre>    <a href="#Key">func (c *MapCursor[K, V]) Key() K</a></pre>
<pre>    <a href="#Next">func (c *MapCursor[K, V]) Next()</a></pre>
<pre>    <a href="#Ok">func (c *MapCursor[K, V]) Ok() bool</a></pre>
<pre>    <a href="#Prev">func (c *MapCursor[K, V]) Prev()</a></pre>
<pre>    <a href="#SeekFirst">func (c *MapCursor[K, V]) SeekFirst()</a></pre>
<pre>    <a href="#SeekFirstGreater">func (c *MapCursor[K, V]) SeekFirstGreater(k K)</a></pre>
<pre>    <a href="#SeekFirstGreaterOrEqual">func (c *MapCursor[K, V]) SeekFirstGreaterOrEqual(k K)</a></pre>
<pre>    <a href="#SeekLast">func (c *MapCursor[K, V]) SeekLast()</a></pre>
<pre>    <a href="#SeekLastLess">func (c *MapCursor[K, V]) SeekLastLess(k K)</a></pre>
<pre>    <a href="#SeekLastLessOrEqual">func (c *MapCursor[K, V]) SeekLastLessOrEqual(k K)</a></pre>
<pre>    <a href="#Value">func (c *MapCursor[K, V]) Value() V</a></pre>
<pre><a href="#Set">type Set</a></pre>
<pre>    <a href="#NewSet">func NewSet[T any](less xsort.Less[T]) Set[T]</a></pre>
<pre>    <a href="#Add">func (s Set[T]) Add(item T)</a></pre>
<pre>    <a href="#Contains">func (s Set[T]) Contains(item T) bool</a></pre>
<pre>    <a href="#Cursor">func (s Set[T]) Cursor() *SetCursor[T]</a></pre>
<pre>    <a href="#First">func (s Set[T]) First() T</a></pre>
<pre>    <a href="#Iterate">func (s Set[T]) Iterate() iterator.Iterator[T]</a></pre>
<pre>    <a href="#Last">func (s Set[T]) Last() T</a></pre>
<pre>    <a href="#Len">func (s Set[T]) Len() int</a></pre>
<pre>    <a href="#Remove">func (s Set[T]) Remove(item T)</a></pre>
<pre><a href="#SetCursor">type SetCursor</a></pre>
<pre>    <a href="#Backward">func (c *SetCursor[T]) Backward() iterator.Iterator[T]</a></pre>
<pre>    <a href="#Forward">func (c *SetCursor[T]) Forward() iterator.Iterator[T]</a></pre>
<pre>    <a href="#Item">func (c *SetCursor[T]) Item() T</a></pre>
<pre>    <a href="#Next">func (c *SetCursor[T]) Next()</a></pre>
<pre>    <a href="#Ok">func (c *SetCursor[T]) Ok() bool</a></pre>
<pre>    <a href="#Prev">func (c *SetCursor[T]) Prev()</a></pre>
<pre>    <a href="#SeekFirst">func (c *SetCursor[T]) SeekFirst()</a></pre>
<pre>    <a href="#SeekFirstGreater">func (c *SetCursor[T]) SeekFirstGreater(x T)</a></pre>
<pre>    <a href="#SeekFirstGreaterOrEqual">func (c *SetCursor[T]) SeekFirstGreaterOrEqual(x T)</a></pre>
<pre>    <a href="#SeekLast">func (c *SetCursor[T]) SeekLast()</a></pre>
<pre>    <a href="#SeekLastLess">func (c *SetCursor[T]) SeekLastLess(x T)</a></pre>
<pre>    <a href="#SeekLastLessOrEqual">func (c *SetCursor[T]) SeekLastLessOrEqual(x T)</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

# Types

## <a id="KVPair"></a><pre>type KVPair</pre>
```go
type KVPair[K any, V any] struct {
	Key   K
	Value V
}
```



## <a id="Map"></a><pre>type Map</pre>
```go
type Map[K any, V any] struct {
	// contains filtered or unexported fields
}
```

Map is a tree-structured key-value map, similar to Go's built-in map but keeps elements in sorted
order by key.


## <a id="NewMap"></a><pre>func NewMap[K any, V any](less <a href="../xsort.md#Less">xsort.Less</a>[K]) <a href="#Map">Map</a>[K, V]</pre>

NewMap returns a Map that uses less to determine the sort order of keys. If !less(a, b) &&
!less(b, a), then a and b are considered the same key. The output of less must not change for any
pair of keys while they are in the map.


## <a id="Contains"></a><pre>func (m <a href="#Map">Map</a>[K, V]) Contains(k K) bool</pre>

Contains returns true if the given key is present in the map.


## <a id="Cursor"></a><pre>func (m <a href="#Map">Map</a>[K, V]) Cursor() *<a href="#MapCursor">MapCursor</a>[K, V]</pre>

Cursor returns a cursor into the map placed at the first element.


## <a id="Delete"></a><pre>func (m <a href="#Map">Map</a>[K, V]) Delete(k K)</pre>

Delete removes the given key from the map.


## <a id="First"></a><pre>func (m <a href="#Map">Map</a>[K, V]) First() (K, V)</pre>

First returns the lowest-keyed entry in the map according to less.


## <a id="Get"></a><pre>func (m <a href="#Map">Map</a>[K, V]) Get(k K) V</pre>

Get returns the value associated with the given key if it is present in the map. Otherwise, it
returns the zero-value of V.


## <a id="Iterate"></a><pre>func (m <a href="#Map">Map</a>[K, V]) Iterate() <a href="../iterator.md#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</pre>

Iterate returns an iterator that yields the elements of the map in sorted order by key.

The map may be safely modified during iteration and the iterator will continue from the
next-lowest key. Thus if the map is modified, the iterator will not necessarily return all of
the keys present in the map.


## <a id="Last"></a><pre>func (m <a href="#Map">Map</a>[K, V]) Last() (K, V)</pre>

Last returns the highest-keyed entry in the map according to less.


## <a id="Len"></a><pre>func (m <a href="#Map">Map</a>[K, V]) Len() int</pre>

Len returns the number of elements in the map.


## <a id="Put"></a><pre>func (m <a href="#Map">Map</a>[K, V]) Put(k K, v V)</pre>

Put inserts the key-value pair into the map, overwriting the value for the key if it already
exists.


## <a id="MapCursor"></a><pre>type MapCursor</pre>
```go
type MapCursor[K any, V any] struct {
	// contains filtered or unexported fields
}
```

MapCursor is a cursor into a Map.

A cursor is usable while a map is being modified. If the element the cursor is at is deleted, the
cursor will still return the old value.


## <a id="Backward"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Backward() <a href="../iterator.md#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</pre>

Backward returns an iterator that starts from the cursor's position and yields all of the
elements less than or equal to the cursor in descending order.

This iterator's Next method is amoritized O(1), unless the map changes in which case the
following Next is O(log(n)) where n is the number of elements in the map.


## <a id="Forward"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Forward() <a href="../iterator.md#Iterator">iterator.Iterator</a>[<a href="#KVPair">KVPair</a>[K, V]]</pre>

Forward returns an iterator that starts from the cursor's position and yields all of the elements
greater than or equal to the cursor in ascending order.

This iterator's Next method is amoritized O(1), unless the map changes in which case the
following Next is O(log(n)) where n is the number of elements in the map.


## <a id="Key"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Key() K</pre>

Key returns the key of the element that the cursor is at. Panics if Ok is false.


## <a id="Next"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Next()</pre>

Next moves the cursor to the next element in the map.

Next is amoritized O(1) unless the map has been modified since the last cursor move, in which
case it's O(log(n)).


## <a id="Ok"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Ok() bool</pre>

Ok returns false if the cursor is not currently placed at an element, for example if Next
advances past the last element.


## <a id="Prev"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Prev()</pre>

Prev moves the cursor to the previous element in the map.

Prev is amoritized O(1) unless the map has been modified since the last cursor move, in which
case it's O(log(n)).


## <a id="SeekFirst"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekFirst()</pre>

SeekFirst moves the cursor to the first element in the map.

SeekFirst is O(log(n)).


## <a id="SeekFirstGreater"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekFirstGreater(k K)</pre>

SeekFirstGreater moves the cursor to the element in the map just after k.

SeekFirstGreater is O(log(n)).


## <a id="SeekFirstGreaterOrEqual"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekFirstGreaterOrEqual(k K)</pre>

SeekFirstGreaterOrEqual moves the cursor to the element in the map with the least key that is
greater than or equal to k.

SeetFirstGreaterOrEqual is O(log(n)).


## <a id="SeekLast"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekLast()</pre>

SeekLast moves the cursor to the last element in the map.

SeekLast is O(log(n)).


## <a id="SeekLastLess"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekLastLess(k K)</pre>

SeekLastLess moves the cursor to the element in the map just before k.

SeekLastLess is O(log(n)).


## <a id="SeekLastLessOrEqual"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) SeekLastLessOrEqual(k K)</pre>

SeekLastLessOrEqual moves the cursor to the element in the map with the greatest key that is less
than or equal to k.

SeekLastLessOrEqual is O(log(n)).


## <a id="Value"></a><pre>func (c *<a href="#MapCursor">MapCursor</a>[K, V]) Value() V</pre>

Value returns the value of the element that the cursor is at. Panics if Ok is false.


## <a id="Set"></a><pre>type Set</pre>
```go
type Set[T any] struct {
	// contains filtered or unexported fields
}
```

Set is a tree-structured set. Sets are a collection of unique elements. Similar to Go's built-in
map[T]struct{} but keeps elements in sorted order.


## <a id="NewSet"></a><pre>func NewSet[T any](less <a href="../xsort.md#Less">xsort.Less</a>[T]) <a href="#Set">Set</a>[T]</pre>

NewSet returns a Set that uses less to determine the sort order of items. If !less(a, b) &&
!less(b, a), then a and b are considered the same item. The output of less must not change for
any pair of items while they are in the set.


## <a id="Add"></a><pre>func (s <a href="#Set">Set</a>[T]) Add(item T)</pre>

Add adds item to the set if it is not already present.


## <a id="Contains"></a><pre>func (s <a href="#Set">Set</a>[T]) Contains(item T) bool</pre>

Contains returns true if item is present in the set.


## <a id="Cursor"></a><pre>func (s <a href="#Set">Set</a>[T]) Cursor() *<a href="#SetCursor">SetCursor</a>[T]</pre>

Cursor returns a cursor into the set placed at the first item.


## <a id="First"></a><pre>func (s <a href="#Set">Set</a>[T]) First() T</pre>

First returns the lowest item in the set according to less.


## <a id="Iterate"></a><pre>func (s <a href="#Set">Set</a>[T]) Iterate() <a href="../iterator.md#Iterator">iterator.Iterator</a>[T]</pre>

Iterate returns an iterator that yields the elements of the set in sorted order.

The set may be safely modified during iteration and the iterator will continue from the
next-lowest item. Thus if the set is modified, the iterator will not necessarily return all of
the items present in the set.


## <a id="Last"></a><pre>func (s <a href="#Set">Set</a>[T]) Last() T</pre>

Last returns the highest item in the set according to less.


## <a id="Len"></a><pre>func (s <a href="#Set">Set</a>[T]) Len() int</pre>

Len returns the number of elements in the set.


## <a id="Remove"></a><pre>func (s <a href="#Set">Set</a>[T]) Remove(item T)</pre>

Remove removes item from the set if it is present, and does nothing otherwise.


## <a id="SetCursor"></a><pre>type SetCursor</pre>
```go
type SetCursor[T any] struct {
	// contains filtered or unexported fields
}
```

SetCursor is a cursor into a Set.

A cursor is usable while a set is being modified. If the item the cursor is at is deleted, the
cursor will still return the old item.


## <a id="Backward"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) Backward() <a href="../iterator.md#Iterator">iterator.Iterator</a>[T]</pre>

Backward returns an iterator that starts from the cursor's position and yields all of the
elements less than or equal to the cursor in descending order.

This iterator's Next method is amoritized O(1), unless the map changes in which case the
following Next is O(log(n)) where n is the number of elements in the map.


## <a id="Forward"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) Forward() <a href="../iterator.md#Iterator">iterator.Iterator</a>[T]</pre>

Forward returns an iterator that starts from the cursor's position and yields all of the elements
greater than or equal to the cursor in ascending order.

This iterator's Next method is amoritized O(1), unless the map changes in which case the
following Next is O(log(n)) where n is the number of elements in the map.


## <a id="Item"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) Item() T</pre>

Item returns the item that the cursor is at. Panics if Ok is false.


## <a id="Next"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) Next()</pre>

Next moves the cursor to the next item in the set.

Next is amoritized O(1) unless the map has been modified since the last cursor move, in which
case it's O(log(n)).


## <a id="Ok"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) Ok() bool</pre>

Ok returns false if the cursor is not currently placed at an item, for example if Next advances
past the last item.


## <a id="Prev"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) Prev()</pre>

Prev moves the cursor to the previous item in the set.

Prev is amoritized O(1) unless the map has been modified since the last cursor move, in which
case it's O(log(n)).


## <a id="SeekFirst"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekFirst()</pre>

SeekFirst moves the cursor to the first item in the set.

SeekFirst is O(log(n)).


## <a id="SeekFirstGreater"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekFirstGreater(x T)</pre>

SeekFirstGreater moves the cursor to the item in the set just after x.

SeekFirstGreater is O(log(n)).


## <a id="SeekFirstGreaterOrEqual"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekFirstGreaterOrEqual(x T)</pre>

SeekFirstGreaterOrEqual moves the cursor to the least item in the set that is greater than or
equal to x.

SeetFirstGreaterOrEqual is O(log(n)).


## <a id="SeekLast"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekLast()</pre>

SeekLast moves the cursor to the last item in the set.

SeekLast is O(log(n)).


## <a id="SeekLastLess"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekLastLess(x T)</pre>

SeekLastLess moves the cursor to the item in the set just before x.

SeekLastLess is O(log(n)).


## <a id="SeekLastLessOrEqual"></a><pre>func (c *<a href="#SetCursor">SetCursor</a>[T]) SeekLastLessOrEqual(x T)</pre>

SeekLastLessOrEqual moves the cursor to the greatest item in the set that is less than or equal
to x.

SeekLastLessOrEqual is O(log(n)).


