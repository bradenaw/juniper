# xheap
--
    import "."

Package xheap contains extensions to the standard library package
container/heap.

## Usage

#### type Heap

```go
type Heap[T any] struct {
}
```

Heap is a min-heap (https://en.wikipedia.org/wiki/Binary_heap). Min-heaps are a
collection structure that provide constant-time access to the minimum element,
and logarithmic-time removal. They are most commonly used as a priority queue.

Push and Pop take amoritized O(log(n)) time where n is the number of items in
the heap.

Len and Peek take O(1) time.

#### func  New

```go
func New[T any](less xsort.Less[T], initial []T) Heap[T]
```
New returns a new Heap which uses less to determine the minimum element.

The elements from initial are added to the heap. initial is modified by New and
utilized by the Heap, so it should not be used after passing to New(). Passing
initial is faster (O(n)) than creating an empty heap and pushing each item (O(n
* log(n))).

#### func (*BADRECV) Grow

```go
func (h *Heap[T]) Grow(n int)
```
Grow allocates sufficient space to add n more elements without needing to
reallocate.

#### func (*BADRECV) Iterate

```go
func (h *Heap[T]) Iterate() iterator.Iterator[T]
```
Iterate iterates over the elements of the heap.

The iterator panics if the heap has been modified since iteration started.

#### func (*BADRECV) Len

```go
func (h *Heap[T]) Len() int
```
Len returns the current number of elements in the heap.

#### func (*BADRECV) Peek

```go
func (h *Heap[T]) Peek() T
```
Peek returns the minimum item in the heap. It panics if h.Len()==0.

#### func (*BADRECV) Pop

```go
func (h *Heap[T]) Pop() T
```
Pop removes and returns the minimum item in the heap. It panics if h.Len()==0.

#### func (*BADRECV) Push

```go
func (h *Heap[T]) Push(item T)
```
Push adds item to the heap.

#### type KP

```go
type KP[K any, P any] struct {
	K K
	P P
}
```

KP holds key and priority for PriorityQueue.

#### type PriorityQueue

```go
type PriorityQueue[K comparable, P any] struct {
}
```

PriorityQueue is a queue that yields items in increasing order of priority.

#### func  NewPriorityQueue

```go
func NewPriorityQueue[K comparable, P any](
	less xsort.Less[P],
	initial []KP[K, P],
) PriorityQueue[K, P]
```
NewPriorityQueue returns a new PriorityQueue which uses less to determine the
minimum element.

The elements from initial are added to the priority queue. initial is modified
by NewPriorityQueue and utilized by the PriorityQueue, so it should not be used
after passing to NewPriorityQueue. Passing initial is faster (O(n)) than
creating an empty priority queue and pushing each item (O(n * log(n))).

Pop, Remove, and Update all take amoritized O(log(n)) time where n is the number
of items in the queue.

Len, Peek, Contains, and Priority take O(1) time.

#### func (*BADRECV) Contains

```go
func (h *PriorityQueue[K, P]) Contains(k K) bool
```
Contains returns true if the given key is present in the priority queue.

#### func (*BADRECV) Grow

```go
func (h *PriorityQueue[K, P]) Grow(n int)
```
Grow allocates sufficient space to add n more elements without needing to
reallocate.

#### func (*BADRECV) Iterate

```go
func (h *PriorityQueue[K, P]) Iterate() iterator.Iterator[K]
```
Iterate iterates over the elements of the priority queue.

The iterator panics if the priority queue has been modified since iteration
started.

#### func (*BADRECV) Len

```go
func (h *PriorityQueue[K, P]) Len() int
```
Len returns the current number of elements in the priority queue.

#### func (*BADRECV) Peek

```go
func (h *PriorityQueue[K, P]) Peek() K
```
Peek returns the key of the lowest-P item in the priority queue. It panics if
h.Len()==0.

#### func (*BADRECV) Pop

```go
func (h *PriorityQueue[K, P]) Pop() K
```
Pop removes and returns the lowest-P item in the priority queue. It panics if
h.Len()==0.

#### func (*BADRECV) Priority

```go
func (h *PriorityQueue[K, P]) Priority(k K) P
```
Priority returns the priority of k, or the zero value of P if k is not present.

#### func (*BADRECV) Remove

```go
func (h *PriorityQueue[K, P]) Remove(k K)
```
Remove removes the item with the given key if present.

#### func (*BADRECV) Update

```go
func (h *PriorityQueue[K, P]) Update(k K, p P)
```
Update updates the priority of k to p, or adds it to the priority queue if not
present.
