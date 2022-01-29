# stream
--
    import "."

package stream allows iterating over sequences of values where iteration may
fail.

## Usage

```go
var (
	// ErrClosedPipe is returned from PipeSender.Send() when the associated stream has already been
	// closed.
	ErrClosedPipe = errors.New("closed pipe")
	// End is returned from Stream.Next when iteration ends successfully.
	End = errors.New("end of stream")
)
```

#### func  Collect

```go
func Collect[T any](ctx context.Context, s Stream[T]) ([]T, error)
```
Collect advances s to the end and returns all of the items seen as a slice.

#### func  Last

```go
func Last[T any](ctx context.Context, s Stream[T], n int) ([]T, error)
```
Last consumes s and returns the last n items. If s yields fewer than n items,
Last returns all of them.

#### func  Pipe

```go
func Pipe[T any](bufferSize int) (*PipeSender[T], Stream[T])
```
Pipe returns a linked sender and receiver pair. Values sent using sender.Send
will be delivered to the given Stream. The Stream will terminate when the sender
is closed.

bufferSize is the number of elements in the buffer between the sender and the
receiver. 0 has the same meaning as for the built-in make(chan).

#### func  Reduce

```go
func Reduce[T any, U any](
	ctx context.Context,
	s Stream[T],
	initial U,
	f func(U, T) (U, error),
) (U, error)
```
Reduce reduces s to a single value using the reduction function f.

#### type Peekable

```go
type Peekable[T any] interface {
	Stream[T]
	// Peek returns the next item of the stream if there is one without consuming it.
	//
	// If Peek returns a value, the next call to Next will return the same value.
	Peek(ctx context.Context) (T, error)
}
```

Peekable allows viewing the next item from a stream without consuming it.

#### func  WithPeek

```go
func WithPeek[T any](s Stream[T]) Peekable[T]
```
WithPeek returns iter with a Peek() method attached.

#### type PipeSender

```go
type PipeSender[T any] struct {
}
```

PipeSender is the send half of a pipe returned by Pipe.

#### func (*BADRECV) Close

```go
func (s *PipeSender[T]) Close(err error)
```
Close closes the PipeSender, signalling to the receiver that no more values will
be sent. If an error is provided, it will surface when closing the receiver.

#### func (*BADRECV) Send

```go
func (s *PipeSender[T]) Send(ctx context.Context, x T) error
```
Send attemps to send x to the receiver. If the receiver closes before x can be
sent, returns ErrClosedPipe immediately. If ctx expires before x can be sent,
returns ctx.Err().

A nil return does not necessarily mean that the receiver will see x, since the
receiver may close early.

#### type Stream

```go
type Stream[T any] interface {
	// Next advances the stream and returns the next item. If the stream is already over, Next
	// returns stream.End in the second return. Note that the final item of the stream has nil in
	// the second return, and it's the following call that returns stream.End.
	Next(ctx context.Context) (T, error)
	// Close ends receiving from the stream. It is invalid to call Next after calling Close.
	Close()
}
```

Stream is used to iterate over a sequence of values. It is similar to Iterator,
except intended for use when iteration may fail for some reason, usually because
the sequence requires I/O to produce.

Streams and the combinator functions are lazy, meaning they do no work until a
call to Next().

Streams do not need to be fully consumed, but streams must be closed. Functions
in this package that are passed streams expect to be the sole user of that
stream going forward, and so will handle closing on your behalf so long as all
streams they return are closed appropriately.

#### func  Batch

```go
func Batch[T any](s Stream[T], maxWait time.Duration, batchSize int) Stream[[]T]
```
Batch returns a stream of non-overlapping batches from s of size batchSize.
Batch is similar to Chunk with the added feature that an underfilled batch will
be delivered to the output stream if any item has been in the batch for more
than maxWait.

#### func  BatchFunc

```go
func BatchFunc[T any](
	s Stream[T],
	maxWait time.Duration,
	full func(batch []T) bool,
) Stream[[]T]
```
BatchFunc returns a stream of non-overlapping batches from s, using full to
determine when a batch is full. BatchFunc is similar to Chunk with the added
feature that an underfilled batch will be delivered to the output stream if any
item has been in the batch for more than maxWait.

#### func  Chan

```go
func Chan[T any](c <-chan T) Stream[T]
```
Chan returns a Stream that receives values from c.

#### func  Chunk

```go
func Chunk[T any](s Stream[T], chunkSize int) Stream[[]T]
```
Chunk returns a stream of non-overlapping chunks from s of size chunkSize. The
last chunk will be smaller than chunkSize if the stream does not contain an even
multiple.

#### func  Compact

```go
func Compact[T comparable](s Stream[T]) Stream[T]
```
Compact elides adjacent duplicates from s.

#### func  CompactFunc

```go
func CompactFunc[T comparable](s Stream[T], eq func(T, T) bool) Stream[T]
```
CompactFunc elides adjacent duplicates from s, using eq to determine duplicates.

#### func  Filter

```go
func Filter[T any](s Stream[T], keep func(T) (bool, error)) Stream[T]
```
Filter returns a Stream that yields only the items from s for which keep returns
true. If keep returns an error, terminates the stream early.

#### func  First

```go
func First[T any](s Stream[T], n int) Stream[T]
```
First returns a Stream that yields the first n items from s.

#### func  Flatten

```go
func Flatten[T any](s Stream[Stream[T]]) Stream[T]
```
Flatten returns a stream that yields all items from all streams yielded by s.

#### func  FromIterator

```go
func FromIterator[T any](iter iterator.Iterator[T]) Stream[T]
```
FromIterator returns a Stream that yields the values from iter. This stream
ignores the context passed to Next during the call to iter.Next.

#### func  Join

```go
func Join[T any](streams ...Stream[T]) Stream[T]
```
Join returns a Stream that yields all elements from streams[0], then all
elements from streams[1], and so on.

#### func  Map

```go
func Map[T any, U any](s Stream[T], f func(t T) (U, error)) Stream[U]
```
Map transforms the values of s using the conversion f. If f returns an error,
terminates the stream early.

#### func  Runs

```go
func Runs[T any](s Stream[T], same func(a, b T) bool) Stream[Stream[T]]
```
Runs returns a stream of streams. The inner streams yield contiguous elements
from s such that same(a, b) returns true for any a and b in the run.

The inner stream should be drained before calling Next on the outer stream.

same(a, a) must return true. If same(a, b) and same(b, c) both return true, then
same(a, c) must also.

#### func  While

```go
func While[T any](s Stream[T], f func(T) (bool, error)) Stream[T]
```
While returns a Stream that terminates before the first item from s for which f
returns false. If f returns an error, terminates the stream early.
