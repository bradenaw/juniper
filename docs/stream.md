# `package stream`

```
import "github.com/bradenaw/juniper/stream"
```

# Overview

package stream allows iterating over sequences of values where iteration may fail.


# Index

<samp><a href="#Collect">func Collect[T any](ctx context.Context, s Stream[T]) ([]T, error)</a></samp>

<samp><a href="#Last">func Last[T any](ctx context.Context, s Stream[T], n int) ([]T, error)</a></samp>

<samp><a href="#Pipe">func Pipe[T any](bufferSize int) (*PipeSender[T], Stream[T])</a></samp>

<samp><a href="#Reduce">func Reduce[T any, U any](
	ctx context.Context,
	s Stream[T],
	initial U,
	f func(U, T) (U, error),
) (U, error)</a></samp>

<samp><a href="#Peekable">type Peekable</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#WithPeek">func WithPeek[T any](s Stream[T]) Peekable[T]</a></samp>

<samp><a href="#PipeSender">type PipeSender</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Close">func (s *PipeSender[T]) Close(err error)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Send">func (s *PipeSender[T]) Send(ctx context.Context, x T) error</a></samp>

<samp><a href="#Stream">type Stream</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Batch">func Batch[T any](s Stream[T], maxWait time.Duration, batchSize int) Stream[[]T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#BatchFunc">func BatchFunc[T any](
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	s Stream[T],
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	maxWait time.Duration,
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	full func(batch []T) bool,
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;) Stream[[]T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Chan">func Chan[T any](c &lt;-chan T) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Chunk">func Chunk[T any](s Stream[T], chunkSize int) Stream[[]T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Compact">func Compact[T comparable](s Stream[T]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#CompactFunc">func CompactFunc[T comparable](s Stream[T], eq func(T, T) bool) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Filter">func Filter[T any](s Stream[T], keep func(T) (bool, error)) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#First">func First[T any](s Stream[T], n int) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Flatten">func Flatten[T any](s Stream[Stream[T]]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#FromIterator">func FromIterator[T any](iter iterator.Iterator[T]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Join">func Join[T any](streams ...Stream[T]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Map">func Map[T any, U any](s Stream[T], f func(t T) (U, error)) Stream[U]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Runs">func Runs[T any](s Stream[T], same func(a, b T) bool) Stream[Stream[T]]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#While">func While[T any](s Stream[T], f func(T) (bool, error)) Stream[T]</a></samp>


# Constants

This section is empty.

# Variables

<pre>
<a id="ErrClosedPipe"></a><a id="End"></a>var (
    ErrClosedPipe = errors.New("closed pipe")
    End = errors.New("end of stream")
)
</pre>

# Functions

<h2><a id="Collect"></a><samp>func <a href="#Collect">Collect</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="#Stream">Stream</a>[T]) ([]T, error)</samp></h2>

Collect advances s to the end and returns all of the items seen as a slice.


### Example 
```go
{
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{"a", "b", "c"}))

	x, err := stream.Collect(ctx, s)
	fmt.Println(err)
	fmt.Println(x)

}
```

Output:
```text
<nil>
[a b c]
```
<h2><a id="Last"></a><samp>func <a href="#Last">Last</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="#Stream">Stream</a>[T], n int) ([]T, error)</samp></h2>

Last consumes s and returns the last n items. If s yields fewer than n items, Last returns
all of them.


<h2><a id="Pipe"></a><samp>func <a href="#Pipe">Pipe</a>[T any](bufferSize int) (*<a href="#PipeSender">PipeSender</a>[T], <a href="#Stream">Stream</a>[T])</samp></h2>

Pipe returns a linked sender and receiver pair. Values sent using sender.Send will be delivered
to the given Stream. The Stream will terminate when the sender is closed.

bufferSize is the number of elements in the buffer between the sender and the receiver. 0 has the
same meaning as for the built-in make(chan).


### Example 
```go
{
	ctx := context.Background()
	sender, receiver := stream.Pipe[int](0)

	go func() {
		sender.Send(ctx, 1)
		sender.Send(ctx, 2)
		sender.Send(ctx, 3)
		sender.Close(nil)
	}()

	defer receiver.Close()
	for {
		item, err := receiver.Next(ctx)
		if err == stream.End {
			break
		} else if err != nil {
			fmt.Printf("stream ended with error: %s\n", err)
			return
		}
		fmt.Println(item)
	}

}
```

Output:
```text
1
2
3
```
### Example error
```go
{
	ctx := context.Background()
	sender, receiver := stream.Pipe[int](0)

	oopsError := errors.New("oops")

	go func() {
		sender.Send(ctx, 1)
		sender.Close(oopsError)
	}()

	defer receiver.Close()
	for {
		item, err := receiver.Next(ctx)
		if err == stream.End {
			fmt.Println("stream ended normally")
			break
		} else if err != nil {
			fmt.Printf("stream ended with error: %s\n", err)
			return
		}
		fmt.Println(item)
	}

}
```

Output:
```text
1
stream ended with error: oops
```
<h2><a id="Reduce"></a><samp>func <a href="#Reduce">Reduce</a>[T any, U any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="#Stream">Stream</a>[T], initial U, f (U, T) (U, error)) (U, error)</samp></h2>

Reduce reduces s to a single value using the reduction function f.


# Types

<h2><a id="Peekable"></a><samp>type Peekable</samp></h2>
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


<h2><a id="WithPeek"></a><samp>func WithPeek[T any](s <a href="#Stream">Stream</a>[T]) <a href="#Peekable">Peekable</a>[T]</samp></h2>

WithPeek returns iter with a Peek() method attached.


<h2><a id="PipeSender"></a><samp>type PipeSender</samp></h2>
```go
type PipeSender[T any] struct {
	// contains filtered or unexported fields
}
```

PipeSender is the send half of a pipe returned by Pipe.


<h2><a id="Close"></a><samp>func (s *<a href="#PipeSender">PipeSender</a>[T]) Close(err error)</samp></h2>

Close closes the PipeSender, signalling to the receiver that no more values will be sent. If an
error is provided, it will surface when closing the receiver.


<h2><a id="Send"></a><samp>func (s *<a href="#PipeSender">PipeSender</a>[T]) Send(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, x T) error</samp></h2>

Send attempts to send x to the receiver. If the receiver closes before x can be sent, returns
ErrClosedPipe immediately. If ctx expires before x can be sent, returns ctx.Err().

A nil return does not necessarily mean that the receiver will see x, since the receiver may close
early.


<h2><a id="Stream"></a><samp>type Stream</samp></h2>
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

Stream is used to iterate over a sequence of values. It is similar to Iterator, except intended
for use when iteration may fail for some reason, usually because the sequence requires I/O to
produce.

Streams and the combinator functions are lazy, meaning they do no work until a call to Next().

Streams do not need to be fully consumed, but streams must be closed. Functions in this package
that are passed streams expect to be the sole user of that stream going forward, and so will
handle closing on your behalf so long as all streams they return are closed appropriately.


<h2><a id="Batch"></a><samp>func Batch[T any](s <a href="#Stream">Stream</a>[T], maxWait <a href="https://pkg.go.dev/time#Duration">time.Duration</a>, batchSize int) <a href="#Stream">Stream</a>[[]T]</samp></h2>

Batch returns a stream of non-overlapping batches from s of size batchSize. Batch is similar to
Chunk with the added feature that an underfilled batch will be delivered to the output stream if
any item has been in the batch for more than maxWait.


### Example 
```go
{
	ctx := context.Background()

	sender, receiver := stream.Pipe[string](0)
	batchStream := stream.Batch(receiver, 50*time.Millisecond, 3)

	wait := make(chan struct{}, 3)
	go func() {
		_ = sender.Send(ctx, "a")
		_ = sender.Send(ctx, "b")

		<-wait
		_ = sender.Send(ctx, "c")
		_ = sender.Send(ctx, "d")
		_ = sender.Send(ctx, "e")
		_ = sender.Send(ctx, "f")
		sender.Close(nil)
	}()

	defer batchStream.Close()
	var batches [][]string
	for {
		batch, err := batchStream.Next(ctx)
		if err == stream.End {
			break
		} else if err != nil {
			fmt.Printf("stream ended with error: %s\n", err)
			return
		}
		batches = append(batches, batch)
		wait <- struct{}{}
	}
	fmt.Println(batches)

}
```

Output:
```text
[[a b] [c d e] [f]]
```
<h2><a id="BatchFunc"></a><samp>func BatchFunc[T any](s <a href="#Stream">Stream</a>[T], maxWait <a href="https://pkg.go.dev/time#Duration">time.Duration</a>, full (batch []T) bool) <a href="#Stream">Stream</a>[[]T]</samp></h2>

BatchFunc returns a stream of non-overlapping batches from s, using full to determine when a
batch is full. BatchFunc is similar to Chunk with the added feature that an underfilled batch
will be delivered to the output stream if any item has been in the batch for more than maxWait.


<h2><a id="Chan"></a><samp>func Chan[T any](c &lt;-chan T) <a href="#Stream">Stream</a>[T]</samp></h2>

Chan returns a Stream that receives values from c.


<h2><a id="Chunk"></a><samp>func Chunk[T any](s <a href="#Stream">Stream</a>[T], chunkSize int) <a href="#Stream">Stream</a>[[]T]</samp></h2>

Chunk returns a stream of non-overlapping chunks from s of size chunkSize. The last chunk will be
smaller than chunkSize if the stream does not contain an even multiple.


<h2><a id="Compact"></a><samp>func Compact[T comparable](s <a href="#Stream">Stream</a>[T]) <a href="#Stream">Stream</a>[T]</samp></h2>

Compact elides adjacent duplicates from s.


<h2><a id="CompactFunc"></a><samp>func CompactFunc[T comparable](s <a href="#Stream">Stream</a>[T], eq (T, T) bool) <a href="#Stream">Stream</a>[T]</samp></h2>

CompactFunc elides adjacent duplicates from s, using eq to determine duplicates.


<h2><a id="Filter"></a><samp>func Filter[T any](s <a href="#Stream">Stream</a>[T], keep (T) (bool, error)) <a href="#Stream">Stream</a>[T]</samp></h2>

Filter returns a Stream that yields only the items from s for which keep returns true. If keep
returns an error, terminates the stream early.


<h2><a id="First"></a><samp>func First[T any](s <a href="#Stream">Stream</a>[T], n int) <a href="#Stream">Stream</a>[T]</samp></h2>

First returns a Stream that yields the first n items from s.


<h2><a id="Flatten"></a><samp>func Flatten[T any](s <a href="#Stream">Stream</a>[<a href="#Stream">Stream</a>[T]]) <a href="#Stream">Stream</a>[T]</samp></h2>

Flatten returns a stream that yields all items from all streams yielded by s.


<h2><a id="FromIterator"></a><samp>func FromIterator[T any](iter <a href="./iterator.md#Iterator">iterator.Iterator</a>[T]) <a href="#Stream">Stream</a>[T]</samp></h2>

FromIterator returns a Stream that yields the values from iter. This stream ignores the context
passed to Next during the call to iter.Next.


<h2><a id="Join"></a><samp>func Join[T any](streams ...) <a href="#Stream">Stream</a>[T]</samp></h2>

Join returns a Stream that yields all elements from streams[0], then all elements from
streams[1], and so on.


<h2><a id="Map"></a><samp>func Map[T any, U any](s <a href="#Stream">Stream</a>[T], f (t T) (U, error)) <a href="#Stream">Stream</a>[U]</samp></h2>

Map transforms the values of s using the conversion f. If f returns an error, terminates the
stream early.


<h2><a id="Runs"></a><samp>func Runs[T any](s <a href="#Stream">Stream</a>[T], same (a, b T) bool) <a href="#Stream">Stream</a>[<a href="#Stream">Stream</a>[T]]</samp></h2>

Runs returns a stream of streams. The inner streams yield contiguous elements from s such that
same(a, b) returns true for any a and b in the run.

The inner stream should be drained before calling Next on the outer stream.

same(a, a) must return true. If same(a, b) and same(b, c) both return true, then same(a, c) must
also.


<h2><a id="While"></a><samp>func While[T any](s <a href="#Stream">Stream</a>[T], f (T) (bool, error)) <a href="#Stream">Stream</a>[T]</samp></h2>

While returns a Stream that terminates before the first item from s for which f returns false.
If f returns an error, terminates the stream early.


