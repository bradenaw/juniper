# xlist
--
    import "."

Package xlist contains extensions to the standard library package
container/list.

## Usage

#### type List

```go
type List[T any] struct {
}
```

List is a doubly-linked list.

#### func (*BADRECV) Back

```go
func (l *List[T]) Back() *Node[T]
```
Back returns the node at the back of the list.

#### func (*BADRECV) Clear

```go
func (l *List[T]) Clear()
```
Clear removes all nodes from the list.

#### func (*BADRECV) Front

```go
func (l *List[T]) Front() *Node[T]
```
Front returns the node at the front of the list.

#### func (*BADRECV) InsertAfter

```go
func (l *List[T]) InsertAfter(value T, mark *Node[T]) *Node[T]
```
InsertBefore adds a new node with the given value after the node mark.

#### func (*BADRECV) InsertBefore

```go
func (l *List[T]) InsertBefore(value T, mark *Node[T]) *Node[T]
```
InsertBefore adds a new node with the given value before the node mark.

#### func (*BADRECV) Len

```go
func (l *List[T]) Len() int
```
Len returns the number of items in the list.

#### func (*BADRECV) MoveAfter

```go
func (l *List[T]) MoveAfter(node *Node[T], mark *Node[T])
```
MoveAfter moves node just after mark. Afterwards, mark.Next() == node &&
node.Prev() == mark.

#### func (*BADRECV) MoveBefore

```go
func (l *List[T]) MoveBefore(node *Node[T], mark *Node[T])
```
MoveBefore moves node just before mark. Afterwards, mark.Prev() == node &&
node.Next() == mark.

#### func (*BADRECV) MoveToBack

```go
func (l *List[T]) MoveToBack(node *Node[T])
```
MoveToFront moves node to the back of the list.

#### func (*BADRECV) MoveToFront

```go
func (l *List[T]) MoveToFront(node *Node[T])
```
MoveToFront moves node to the front of the list.

#### func (*BADRECV) PushBack

```go
func (l *List[T]) PushBack(value T) *Node[T]
```
PushFront adds a new node with the given value to the back of the list.

#### func (*BADRECV) PushFront

```go
func (l *List[T]) PushFront(value T) *Node[T]
```
PushFront adds a new node with the given value to the front of the list.

#### func (*BADRECV) Remove

```go
func (l *List[T]) Remove(node *Node[T])
```
Remove removes node from the list.

#### type Node

```go
type Node[T any] struct {

	// Value is user-controlled, and never modified by this package.
	Value T
}
```

Node is a node in a linked-list.

#### func (*BADRECV) Next

```go
func (n *Node[T]) Next() *Node[T]
```
Next returns the next node in the list that n is a part of, if there is one.

#### func (*BADRECV) Prev

```go
func (n *Node[T]) Prev() *Node[T]
```
Prev returns the previous node in the list that n is a part of, if there is one.
