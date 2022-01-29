# deque
--
    import "."

Package deque contains a double-ended queue.

## Usage

#### type Deque

```go
type Deque[T any] struct {
}
```

Deque is a double-ended queue, allowing push and pop to both the front and back
of the queue. Pushes and pops are amoritized O(1). The zero-value is ready to
use. Deque should not be copied after first use.

#### func (*BADRECV) Back

```go
func (r *Deque[T]) Back() T
```
Back returns the item at the back of the deque. It panics if the deque is empty.

#### func (*BADRECV) Front

```go
func (r *Deque[T]) Front() T
```
Front returns the item at the front of the deque. It panics if the deque is
empty.

#### func (*BADRECV) Grow

```go
func (r *Deque[T]) Grow(n int)
```
Grow allocates sufficient space to add n more items without needing to
reallocate.

#### func (*BADRECV) Item

```go
func (r *Deque[T]) Item(i int) T
```
Item returns the ith item in the deque. 0 is the front and r.Len()-1 is the
back.

#### func (*BADRECV) Iterate

```go
func (r *Deque[T]) Iterate() iterator.Iterator[T]
```
Iterate iterates over the elements of the deque.

The iterator panics if the deque has been modified since iteration started.

#### func (*BADRECV) Len

```go
func (r *Deque[T]) Len() int
```
Len returns the number of items in the deque.

#### func (*BADRECV) PopBack

```go
func (r *Deque[T]) PopBack() T
```
PopBack removes and returns the item at the back of the deque. It panics if the
deque is empty.

#### func (*BADRECV) PopFront

```go
func (r *Deque[T]) PopFront() T
```
PopFront removes and returns the item at the front of the deque. It panics if
the deque is empty.

#### func (*BADRECV) PushBack

```go
func (r *Deque[T]) PushBack(item T)
```
PushFront adds item to the back of the deque.

#### func (*BADRECV) PushFront

```go
func (r *Deque[T]) PushFront(item T)
```
PushFront adds item to the front of the deque.
