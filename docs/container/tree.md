# tree
--
    import "."

Package tree contains an implementation of a B-tree Map and Set. These are
similiar to Go's map built-in, but keep elements in sorted order.

## Usage

#### type KVPair

```go
type KVPair[K any, V any] struct {
	Key   K
	Value V
}
```


#### type Map

```go
type Map[K any, V any] struct {
}
```

Map is a tree-structured key-value map, similar to Go's built-in map but keeps
elements in sorted order by key.

#### func  NewMap

```go
func NewMap[K any, V any](less xsort.Less[K]) Map[K, V]
```
NewMap returns a Map that uses less to determine the sort order of keys. If
!less(a, b) && !less(b, a), then a and b are considered the same key. The output
of less must not change for any pair of keys while they are in the map.

#### func (BADRECV) Contains

```go
func (m Map[K, V]) Contains(k K) bool
```
Contains returns true if the given key is present in the map.

#### func (BADRECV) Cursor

```go
func (m Map[K, V]) Cursor() *MapCursor[K, V]
```
Cursor returns a cursor into the map placed at the first element.

#### func (BADRECV) Delete

```go
func (m Map[K, V]) Delete(k K)
```
Delete removes the given key from the map.

#### func (BADRECV) First

```go
func (m Map[K, V]) First() (K, V)
```
First returns the lowest-keyed entry in the map according to less.

#### func (BADRECV) Get

```go
func (m Map[K, V]) Get(k K) V
```
Get returns the value associated with the given key if it is present in the map.
Otherwise, it returns the zero-value of V.

#### func (BADRECV) Iterate

```go
func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]]
```
Iterate returns an iterator that yields the elements of the map in sorted order
by key.

The map may be safely modified during iteration and the iterator will continue
from the next-lowest key. Thus if the map is modified, the iterator will not
necessarily return all of the keys present in the map.

#### func (BADRECV) Last

```go
func (m Map[K, V]) Last() (K, V)
```
Last returns the highest-keyed entry in the map according to less.

#### func (BADRECV) Len

```go
func (m Map[K, V]) Len() int
```
Len returns the number of elements in the map.

#### func (BADRECV) Put

```go
func (m Map[K, V]) Put(k K, v V)
```
Put inserts the key-value pair into the map, overwriting the value for the key
if it already exists.

#### type MapCursor

```go
type MapCursor[K any, V any] struct {
}
```

MapCursor is a cursor into a Map.

A cursor is usable while a map is being modified. If the element the cursor is
at is deleted, the cursor will still return the old value.

#### func (*BADRECV) Backward

```go
func (c *MapCursor[K, V]) Backward() iterator.Iterator[KVPair[K, V]]
```
Backward returns an iterator that starts from the cursor's position and yields
all of the elements less than or equal to the cursor in descending order.

This iterator's Next method is amoritized O(1), unless the map changes in which
case the following Next is O(log(n)) where n is the number of elements in the
map.

#### func (*BADRECV) Forward

```go
func (c *MapCursor[K, V]) Forward() iterator.Iterator[KVPair[K, V]]
```
Forward returns an iterator that starts from the cursor's position and yields
all of the elements greater than or equal to the cursor in ascending order.

This iterator's Next method is amoritized O(1), unless the map changes in which
case the following Next is O(log(n)) where n is the number of elements in the
map.

#### func (*BADRECV) Key

```go
func (c *MapCursor[K, V]) Key() K
```
Key returns the key of the element that the cursor is at. Panics if Ok is false.

#### func (*BADRECV) Next

```go
func (c *MapCursor[K, V]) Next()
```
Next moves the cursor to the next element in the map.

Next is amoritized O(1) unless the map has been modified since the last cursor
move, in which case it's O(log(n)).

#### func (*BADRECV) Ok

```go
func (c *MapCursor[K, V]) Ok() bool
```
Ok returns false if the cursor is not currently placed at an element, for
example if Next advances past the last element.

#### func (*BADRECV) Prev

```go
func (c *MapCursor[K, V]) Prev()
```
Prev moves the cursor to the previous element in the map.

Prev is amoritized O(1) unless the map has been modified since the last cursor
move, in which case it's O(log(n)).

#### func (*BADRECV) SeekFirst

```go
func (c *MapCursor[K, V]) SeekFirst()
```
SeekFirst moves the cursor to the first element in the map.

SeekFirst is O(log(n)).

#### func (*BADRECV) SeekFirstGreater

```go
func (c *MapCursor[K, V]) SeekFirstGreater(k K)
```
SeekFirstGreater moves the cursor to the element in the map just after k.

SeekFirstGreater is O(log(n)).

#### func (*BADRECV) SeekFirstGreaterOrEqual

```go
func (c *MapCursor[K, V]) SeekFirstGreaterOrEqual(k K)
```
SeekFirstGreaterOrEqual moves the cursor to the element in the map with the
least key that is greater than or equal to k.

SeetFirstGreaterOrEqual is O(log(n)).

#### func (*BADRECV) SeekLast

```go
func (c *MapCursor[K, V]) SeekLast()
```
SeekLast moves the cursor to the last element in the map.

SeekLast is O(log(n)).

#### func (*BADRECV) SeekLastLess

```go
func (c *MapCursor[K, V]) SeekLastLess(k K)
```
SeekLastLess moves the cursor to the element in the map just before k.

SeekLastLess is O(log(n)).

#### func (*BADRECV) SeekLastLessOrEqual

```go
func (c *MapCursor[K, V]) SeekLastLessOrEqual(k K)
```
SeekLastLessOrEqual moves the cursor to the element in the map with the greatest
key that is less than or equal to k.

SeekLastLessOrEqual is O(log(n)).

#### func (*BADRECV) Value

```go
func (c *MapCursor[K, V]) Value() V
```
Value returns the value of the element that the cursor is at. Panics if Ok is
false.

#### type Set

```go
type Set[T any] struct {
}
```

Set is a tree-structured set. Sets are a collection of unique elements. Similar
to Go's built-in map[T]struct{} but keeps elements in sorted order.

#### func  NewSet

```go
func NewSet[T any](less xsort.Less[T]) Set[T]
```
NewSet returns a Set that uses less to determine the sort order of items. If
!less(a, b) && !less(b, a), then a and b are considered the same item. The
output of less must not change for any pair of items while they are in the set.

#### func (BADRECV) Add

```go
func (s Set[T]) Add(item T)
```
Add adds item to the set if it is not already present.

#### func (BADRECV) Contains

```go
func (s Set[T]) Contains(item T) bool
```
Contains returns true if item is present in the set.

#### func (BADRECV) Cursor

```go
func (s Set[T]) Cursor() *SetCursor[T]
```
Cursor returns a cursor into the set placed at the first item.

#### func (BADRECV) First

```go
func (s Set[T]) First() T
```
First returns the lowest item in the set according to less.

#### func (BADRECV) Iterate

```go
func (s Set[T]) Iterate() iterator.Iterator[T]
```
Iterate returns an iterator that yields the elements of the set in sorted order.

The set may be safely modified during iteration and the iterator will continue
from the next-lowest item. Thus if the set is modified, the iterator will not
necessarily return all of the items present in the set.

#### func (BADRECV) Last

```go
func (s Set[T]) Last() T
```
Last returns the highest item in the set according to less.

#### func (BADRECV) Len

```go
func (s Set[T]) Len() int
```
Len returns the number of elements in the set.

#### func (BADRECV) Remove

```go
func (s Set[T]) Remove(item T)
```
Remove removes item from the set if it is present, and does nothing otherwise.

#### type SetCursor

```go
type SetCursor[T any] struct {
}
```

SetCursor is a cursor into a Set.

A cursor is usable while a set is being modified. If the item the cursor is at
is deleted, the cursor will still return the old item.

#### func (*BADRECV) Backward

```go
func (c *SetCursor[T]) Backward() iterator.Iterator[T]
```
Backward returns an iterator that starts from the cursor's position and yields
all of the elements less than or equal to the cursor in descending order.

This iterator's Next method is amoritized O(1), unless the map changes in which
case the following Next is O(log(n)) where n is the number of elements in the
map.

#### func (*BADRECV) Forward

```go
func (c *SetCursor[T]) Forward() iterator.Iterator[T]
```
Forward returns an iterator that starts from the cursor's position and yields
all of the elements greater than or equal to the cursor in ascending order.

This iterator's Next method is amoritized O(1), unless the map changes in which
case the following Next is O(log(n)) where n is the number of elements in the
map.

#### func (*BADRECV) Item

```go
func (c *SetCursor[T]) Item() T
```
Item returns the item that the cursor is at. Panics if Ok is false.

#### func (*BADRECV) Next

```go
func (c *SetCursor[T]) Next()
```
Next moves the cursor to the next item in the set.

Next is amoritized O(1) unless the map has been modified since the last cursor
move, in which case it's O(log(n)).

#### func (*BADRECV) Ok

```go
func (c *SetCursor[T]) Ok() bool
```
Ok returns false if the cursor is not currently placed at an item, for example
if Next advances past the last item.

#### func (*BADRECV) Prev

```go
func (c *SetCursor[T]) Prev()
```
Prev moves the cursor to the previous item in the set.

Prev is amoritized O(1) unless the map has been modified since the last cursor
move, in which case it's O(log(n)).

#### func (*BADRECV) SeekFirst

```go
func (c *SetCursor[T]) SeekFirst()
```
SeekFirst moves the cursor to the first item in the set.

SeekFirst is O(log(n)).

#### func (*BADRECV) SeekFirstGreater

```go
func (c *SetCursor[T]) SeekFirstGreater(x T)
```
SeekFirstGreater moves the cursor to the item in the set just after x.

SeekFirstGreater is O(log(n)).

#### func (*BADRECV) SeekFirstGreaterOrEqual

```go
func (c *SetCursor[T]) SeekFirstGreaterOrEqual(x T)
```
SeekFirstGreaterOrEqual moves the cursor to the least item in the set that is
greater than or equal to x.

SeetFirstGreaterOrEqual is O(log(n)).

#### func (*BADRECV) SeekLast

```go
func (c *SetCursor[T]) SeekLast()
```
SeekLast moves the cursor to the last item in the set.

SeekLast is O(log(n)).

#### func (*BADRECV) SeekLastLess

```go
func (c *SetCursor[T]) SeekLastLess(x T)
```
SeekLastLess moves the cursor to the item in the set just before x.

SeekLastLess is O(log(n)).

#### func (*BADRECV) SeekLastLessOrEqual

```go
func (c *SetCursor[T]) SeekLastLessOrEqual(x T)
```
SeekLastLessOrEqual moves the cursor to the greatest item in the set that is
less than or equal to x.

SeekLastLessOrEqual is O(log(n)).
