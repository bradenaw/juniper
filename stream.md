# `package stream`

```
import "github.com/bradenaw/juniper/stream"
```

## Overview

Package stream allows iterating over sequences of values where iteration may fail, for example
when it involves I/O.


## Index

<samp><a href="#Collect">func Collect[T any](ctx context.Context, s Stream[T]) ([]T, error)</a></samp>

<samp><a href="#Last">func Last[T any](ctx context.Context, s Stream[T], n int) ([]T, error)</a></samp>

<samp><a href="#One">func One[T any](ctx context.Context, s Stream[T]) (T, error)</a></samp>

<samp><a href="#Pipe">func Pipe[T any](bufferSize int) (*PipeSender[T], Stream[T])</a></samp>

<samp><a href="#Reduce">func Reduce[T any, U any](
	ctx context.Context,
	s Stream[T],
	initial U,
	f func(U, T) (U, error),
) (U, error)</a></samp>

<samp><a href="#Peekable">type Peekable</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#WithPeek">func WithPeek[T any](s Stream[T]) Peekable[T]</a></samp>

<samp><a href="#PipeSender">type PipeSender</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Close">func (s *PipeSender[T]) Close(err error)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Send">func (s *PipeSender[T]) Send(ctx context.Context, x T) error</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#TrySend">func (s *PipeSender[T]) TrySend(ctx context.Context, x T) (bool, error)</a></samp>

<samp><a href="#Stream">type Stream</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Batch">func Batch[T any](s Stream[T], maxWait time.Duration, batchSize int) Stream[[]T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#BatchFunc">func BatchFunc[T any](
&nbsp;&nbsp;&nbsp;&nbsp;	s Stream[T],
&nbsp;&nbsp;&nbsp;&nbsp;	maxWait time.Duration,
&nbsp;&nbsp;&nbsp;&nbsp;	full func(batch []T) bool,
&nbsp;&nbsp;&nbsp;&nbsp;) Stream[[]T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Chan">func Chan[T any](c &lt;-chan T) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Chunk">func Chunk[T any](s Stream[T], chunkSize int) Stream[[]T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Compact">func Compact[T comparable](s Stream[T]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#CompactFunc">func CompactFunc[T any](s Stream[T], eq func(T, T) bool) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Empty">func Empty[T any]() Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Error">func Error[T any](err error) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Filter">func Filter[T any](s Stream[T], keep func(context.Context, T) (bool, error)) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#First">func First[T any](s Stream[T], n int) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Flatten">func Flatten[T any](s Stream[Stream[T]]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#FromIterator">func FromIterator[T any](iter iterator.Iterator[T]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Join">func Join[T any](streams ...Stream[T]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Map">func Map[T any, U any](s Stream[T], f func(context.Context, T) (U, error)) Stream[U]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Merge">func Merge[T any](in ...Stream[T]) Stream[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Runs">func Runs[T any](s Stream[T], same func(a, b T) bool) Stream[Stream[T]]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#While">func While[T any](s Stream[T], f func(context.Context, T) (bool, error)) Stream[T]</a></samp>


## Constants

This section is empty.

## Variables

<pre>
<a id="End"></a><a id="ErrClosedPipe"></a><a id="ErrMoreThanOne"></a><a id="ErrEmpty"></a>var (
    End = errors.New("end of stream")
    ErrClosedPipe = errors.New("closed pipe")
    ErrMoreThanOne = errors.New("stream had more than one item")
    ErrEmpty = errors.New("stream empty")
)
</pre>

## Functions

<h3><a id="Collect"></a><samp>func <a href="#Collect">Collect</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="#Stream">Stream</a>[T]) ([]T, error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L302">src</a></small></sub></h3>

Collect advances s to the end and returns all of the items seen as a slice.


#### Example 
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
<h3><a id="Last"></a><samp>func <a href="#Last">Last</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="#Stream">Stream</a>[T], n int) ([]T, error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L319">src</a></small></sub></h3>

Last consumes s and returns the last n items. If s yields fewer than n items, Last returns
all of them.


#### Example 
```go
{
	ctx := context.Background()
	s := stream.FromIterator(iterator.Counter(10))
	last5, _ := stream.Last(ctx, s, 5)
	fmt.Println(last5)

	s = stream.FromIterator(iterator.Counter(3))
	last5, _ = stream.Last(ctx, s, 5)
	fmt.Println(last5)

}
```

Output:
```text
[5 6 7 8 9]
[0 1 2]
```
<h3><a id="One"></a><samp>func <a href="#One">One</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="#Stream">Stream</a>[T]) (T, error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L345">src</a></small></sub></h3>

One returns the only item that s yields. Returns an error if encountered, or if s yields zero or
more than one item.


#### Example 
```go
{
	ctx := context.Background()

	s := stream.FromIterator(iterator.Slice([]string{"a"}))
	item, err := stream.One(ctx, s)
	fmt.Println(err == nil)
	fmt.Println(item)

	s = stream.FromIterator(iterator.Slice([]string{"a", "b"}))
	_, err = stream.One(ctx, s)
	fmt.Println(err == stream.ErrMoreThanOne)

	s = stream.FromIterator(iterator.Slice([]string{}))
	_, err = stream.One(ctx, s)
	fmt.Println(err == stream.ErrEmpty)

}
```

Output:
```text
true
a
true
true
```
<h3><a id="Pipe"></a><samp>func <a href="#Pipe">Pipe</a>[T any](bufferSize int) (*<a href="#PipeSender">PipeSender</a>[T], <a href="#Stream">Stream</a>[T])</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L137">src</a></small></sub></h3>

Pipe returns a linked sender and receiver pair. Values sent using sender.Send will be delivered
to the given Stream. The Stream will terminate when the sender is closed.

bufferSize is the number of elements in the buffer between the sender and the receiver. 0 has the
same meaning as for the built-in make(chan).


#### Example 
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
#### Example error
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
<h3><a id="Reduce"></a><samp>func <a href="#Reduce">Reduce</a>[T any, U any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="#Stream">Stream</a>[T], initial U, f func(U, T) (U, error)) (U, error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L363">src</a></small></sub></h3>

Reduce reduces s to a single value using the reduction function f.


#### Example 
```go
{
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]int{1, 2, 3, 4, 5}))

	sum, _ := stream.Reduce(ctx, s, 0, func(x, y int) (int, error) {
		return x + y, nil
	})
	fmt.Println(sum)

	s = stream.FromIterator(iterator.Slice([]int{1, 3, 2, 3}))

	first := true
	ewma, _ := stream.Reduce(ctx, s, 0, func(running float64, item int) (float64, error) {
		if first {
			first = false
			return float64(item), nil
		}
		return running*0.5 + float64(item)*0.5, nil
	})

	fmt.Println(ewma)

}
```

Output:
```text
15
2.5
```
## Types

<h3><a id="Peekable"></a><samp>type Peekable</samp></h3>
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


<h3><a id="WithPeek"></a><samp>func WithPeek[T any](s <a href="#Stream">Stream</a>[T]) <a href="#Peekable">Peekable</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L257">src</a></small></sub></h3>

WithPeek returns iter with a Peek() method attached.


<h3><a id="PipeSender"></a><samp>type PipeSender</samp></h3>
```go
type PipeSender[T any] struct {
	// contains filtered or unexported fields
}
```

PipeSender is the send half of a pipe returned by Pipe.


<h3><a id="Close"></a><samp>func (s *<a href="#PipeSender">PipeSender</a>[T]) Close(err error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L217">src</a></small></sub></h3>

Close closes the PipeSender, signalling to the receiver that no more values will be sent. If an
error is provided, it will surface to the receiver's Next and to any concurrent Sends.

Close may only be called once.


<h3><a id="Send"></a><samp>func (s *<a href="#PipeSender">PipeSender</a>[T]) Send(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, x T) error</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L174">src</a></small></sub></h3>

Send attempts to send x to the receiver. If the receiver closes before x can be sent, returns
ErrClosedPipe immediately. If ctx expires before x can be sent, returns ctx.Err().

A nil return does not necessarily mean that the receiver will see x, since the receiver may close
early.

Send may be called concurrently with other Sends and with Close.


<h3><a id="TrySend"></a><samp>func (s *<a href="#PipeSender">PipeSender</a>[T]) TrySend(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, x T) (bool, error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L195">src</a></small></sub></h3>

TrySend attempts to send x to the receiver, but returns (false, nil) if the pipe's buffer is
already full instead of blocking. If the receiver is already closed, returns ErrClosedPipe. If
ctx expires before x can be sent, returns ctx.Err().

A (true, nil) return does not necessarily mean that the receiver will see x, since the receiver
may close early.

TrySend may be called concurrently with other Sends and with Close.


<h3><a id="Stream"></a><samp>type Stream</samp></h3>
```go
type Stream[T any] interface {
	// Next advances the stream and returns the next item. If the stream is already over, Next
	// returns stream.End in the second return. Note that the final item of the stream has nil in
	// the second return, and it's the following call that returns stream.End.
	//
	// Once a Next call returns stream.End, it is expected that the Stream will return stream.End to
	// every Next call afterwards.
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


<h3><a id="Batch"></a><samp>func Batch[T any](s <a href="#Stream">Stream</a>[T], maxWait <a href="https://pkg.go.dev/time#Duration">time.Duration</a>, batchSize int) <a href="#Stream">Stream</a>[[]T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L394">src</a></small></sub></h3>

Batch returns a stream of non-overlapping batches from s of size batchSize. Batch is similar to
Chunk with the added feature that an underfilled batch will be delivered to the output stream if
any item has been in the batch for more than maxWait.


#### Example 
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
<h3><a id="BatchFunc"></a><samp>func BatchFunc[T any](s <a href="#Stream">Stream</a>[T], maxWait <a href="https://pkg.go.dev/time#Duration">time.Duration</a>, full func(batch []T) bool) <a href="#Stream">Stream</a>[[]T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L403">src</a></small></sub></h3>

BatchFunc returns a stream of non-overlapping batches from s, using full to determine when a
batch is full. BatchFunc is similar to Chunk with the added feature that an underfilled batch
will be delivered to the output stream if any item has been in the batch for more than maxWait.


<h3><a id="Chan"></a><samp>func Chan[T any](c &lt;-chan T) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L55">src</a></small></sub></h3>

Chan returns a Stream that receives values from c.


#### Example 
```go
{
	ctx := context.Background()

	c := make(chan string, 3)
	c <- "a"
	c <- "b"
	c <- "c"
	close(c)
	s := stream.Chan(c)

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
<h3><a id="Chunk"></a><samp>func Chunk[T any](s <a href="#Stream">Stream</a>[T], chunkSize int) <a href="#Stream">Stream</a>[[]T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L595">src</a></small></sub></h3>

Chunk returns a stream of non-overlapping chunks from s of size chunkSize. The last chunk will be
smaller than chunkSize if the stream does not contain an even multiple.


#### Example 
```go
{
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{"a", "b", "c", "d", "e", "f", "g", "h"}))

	chunked := stream.Chunk(s, 3)
	item, _ := chunked.Next(ctx)
	fmt.Println(item)
	item, _ = chunked.Next(ctx)
	fmt.Println(item)
	item, _ = chunked.Next(ctx)
	fmt.Println(item)

}
```

Output:
```text
[a b c]
[d e f]
[g h]
```
<h3><a id="Compact"></a><samp>func Compact[T comparable](s <a href="#Stream">Stream</a>[T]) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L636">src</a></small></sub></h3>

Compact elides adjacent duplicates from s.


#### Example 
```go
{
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{"a", "a", "b", "c", "c", "c", "a"}))
	compactStream := stream.Compact(s)
	compacted, _ := stream.Collect(ctx, compactStream)
	fmt.Println(compacted)

}
```

Output:
```text
[a b c a]
```
<h3><a id="CompactFunc"></a><samp>func CompactFunc[T any](s <a href="#Stream">Stream</a>[T], eq func(T, T) bool) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L643">src</a></small></sub></h3>

CompactFunc elides adjacent duplicates from s, using eq to determine duplicates.


#### Example 
```go
{
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}))
	compactStream := stream.CompactFunc(s, func(a, b string) bool {
		return a[0] == b[0]
	})
	compacted, _ := stream.Collect(ctx, compactStream)
	fmt.Println(compacted)

}
```

Output:
```text
[bank ghost yaw]
```
<h3><a id="Empty"></a><samp>func Empty[T any]() <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L79">src</a></small></sub></h3>

Empty returns a Stream that yields stream.End immediately.


<h3><a id="Error"></a><samp>func Error[T any](err error) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L93">src</a></small></sub></h3>

Error returns a Stream that immediately produces err from Next.


#### Example 
```go
{
	ctx := context.Background()

	s := stream.Error[int](errors.New("foo"))

	_, err := s.Next(ctx)
	fmt.Println(err)

}
```

Output:
```text
foo
```
<h3><a id="Filter"></a><samp>func Filter[T any](s <a href="#Stream">Stream</a>[T], keep func(<a href="https://pkg.go.dev/context#Context">context.Context</a>, T) (bool, error)) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L682">src</a></small></sub></h3>

Filter returns a Stream that yields only the items from s for which keep returns true. If keep
returns an error, terminates the stream early.


#### Example 
```go
{
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]int{1, 2, 3, 4, 5, 6}))

	evensStream := stream.Filter(s, func(ctx context.Context, x int) (bool, error) {
		return x%2 == 0, nil
	})
	evens, _ := stream.Collect(ctx, evensStream)
	fmt.Println(evens)

}
```

Output:
```text
[2 4 6]
```
<h3><a id="First"></a><samp>func First[T any](s <a href="#Stream">Stream</a>[T], n int) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L713">src</a></small></sub></h3>

First returns a Stream that yields the first n items from s.


#### Example 
```go
{
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{"a", "b", "c", "d", "e"}))

	first3Stream := stream.First(s, 3)
	first3, _ := stream.Collect(ctx, first3Stream)
	fmt.Println(first3)

}
```

Output:
```text
[a b c]
```
<h3><a id="Flatten"></a><samp>func Flatten[T any](s <a href="#Stream">Stream</a>[<a href="#Stream">Stream</a>[T]]) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L740">src</a></small></sub></h3>

Flatten returns a stream that yields all items from all streams yielded by s.


#### Example 
```go
{
	ctx := context.Background()

	s := stream.FromIterator(iterator.Slice([]stream.Stream[int]{
		stream.FromIterator(iterator.Slice([]int{0, 1, 2})),
		stream.FromIterator(iterator.Slice([]int{3, 4, 5, 6})),
		stream.FromIterator(iterator.Slice([]int{7})),
	}))

	allStream := stream.Flatten(s)
	all, _ := stream.Collect(ctx, allStream)

	fmt.Println(all)

}
```

Output:
```text
[0 1 2 3 4 5 6 7]
```
<h3><a id="FromIterator"></a><samp>func FromIterator[T any](iter <a href="./iterator.html#Iterator">iterator.Iterator</a>[T]) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L110">src</a></small></sub></h3>

FromIterator returns a Stream that yields the values from iter. This stream ignores the context
passed to Next during the call to iter.Next.


<h3><a id="Join"></a><samp>func Join[T any](streams ...<a href="#Stream">Stream</a>[T]) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L782">src</a></small></sub></h3>

Join returns a Stream that yields all elements from streams[0], then all elements from
streams[1], and so on.


#### Example 
```go
{
	ctx := context.Background()

	s := stream.Join(
		stream.FromIterator(iterator.Counter(3)),
		stream.FromIterator(iterator.Counter(5)),
		stream.FromIterator(iterator.Counter(2)),
	)

	all, _ := stream.Collect(ctx, s)

	fmt.Println(all)

}
```

Output:
```text
[0 1 2 0 1 2 3 4 0 1]
```
<h3><a id="Map"></a><samp>func Map[T any, U any](s <a href="#Stream">Stream</a>[T], f func(<a href="https://pkg.go.dev/context#Context">context.Context</a>, T) (U, error)) <a href="#Stream">Stream</a>[U]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L814">src</a></small></sub></h3>

Map transforms the values of s using the conversion f. If f returns an error, terminates the
stream early.


#### Example 
```go
{
	ctx := context.Background()

	s := stream.FromIterator(iterator.Counter(5))
	halfStream := stream.Map(s, func(ctx context.Context, x int) (float64, error) {
		return float64(x) / 2, nil
	})
	all, _ := stream.Collect(ctx, halfStream)
	fmt.Println(all)

}
```

Output:
```text
[0 0.5 1 1.5 2]
```
<h3><a id="Merge"></a><samp>func Merge[T any](in ...<a href="#Stream">Stream</a>[T]) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L842">src</a></small></sub></h3>

Merge merges the in streams, returning a stream that yields all elements from all of them as they
arrive.


#### Example 
```go
{
	ctx := context.Background()

	a := stream.FromIterator(iterator.Slice([]string{"a", "b", "c"}))
	b := stream.FromIterator(iterator.Slice([]string{"x", "y", "z"}))
	c := stream.FromIterator(iterator.Slice([]string{"m", "n"}))

	s := stream.Merge(a, b, c)

	for {
		item, err := s.Next(ctx)
		if err == stream.End {
			break
		} else if err != nil {
			panic(err)
		}

		fmt.Println(item)
	}

}
```

Unordered output:
```text
m
b
a
n
x
c
z
y
```
<h3><a id="Runs"></a><samp>func Runs[T any](s <a href="#Stream">Stream</a>[T], same func(a, b T) bool) <a href="#Stream">Stream</a>[<a href="#Stream">Stream</a>[T]]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L899">src</a></small></sub></h3>

Runs returns a stream of streams. The inner streams yield contiguous elements from s such that
same(a, b) returns true for any a and b in the run.

The inner stream should be drained before calling Next on the outer stream.

same(a, a) must return true. If same(a, b) and same(b, c) both return true, then same(a, c) must
also.


#### Example 
```go
{
	ctx := context.Background()

	s := stream.FromIterator(iterator.Slice([]int{2, 4, 0, 7, 1, 3, 9, 2, 8}))

	parityRuns := stream.Runs(s, func(a, b int) bool {
		return a%2 == b%2
	})

	one, _ := parityRuns.Next(ctx)
	allOne, _ := stream.Collect(ctx, one)
	fmt.Println(allOne)
	two, _ := parityRuns.Next(ctx)
	allTwo, _ := stream.Collect(ctx, two)
	fmt.Println(allTwo)
	three, _ := parityRuns.Next(ctx)
	allThree, _ := stream.Collect(ctx, three)
	fmt.Println(allThree)

}
```

Output:
```text
[2 4 0]
[7 1 3 9]
[2 8]
```
<h3><a id="While"></a><samp>func While[T any](s <a href="#Stream">Stream</a>[T], f func(<a href="https://pkg.go.dev/context#Context">context.Context</a>, T) (bool, error)) <a href="#Stream">Stream</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/stream/stream.go#L963">src</a></small></sub></h3>

While returns a Stream that terminates before the first item from s for which f returns false.
If f returns an error, terminates the stream early.


#### Example 
```go
{
	ctx := context.Background()

	s := stream.FromIterator(iterator.Slice([]string{
		"aardvark",
		"badger",
		"cheetah",
		"dinosaur",
		"egret",
	}))

	beforeD := stream.While(s, func(ctx context.Context, s string) (bool, error) {
		return s < "d", nil
	})

	out, _ := stream.Collect(ctx, beforeD)
	fmt.Println(out)

}
```

Output:
```text
[aardvark badger cheetah]
```
