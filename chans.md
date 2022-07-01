# `package chans`

```
import "github.com/bradenaw/juniper/chans"
```

## Overview

Package chans contains functions for manipulating channels.


## Index

<samp><a href="#Merge">func Merge[T any](out chan&lt;- T, in ...&lt;-chan T)</a></samp>

<samp><a href="#RecvContext">func RecvContext[T any](ctx context.Context, c &lt;-chan T) (T, bool, error)</a></samp>

<samp><a href="#Replicate">func Replicate[T any](src &lt;-chan T, dsts ...chan&lt;- T)</a></samp>

<samp><a href="#SendContext">func SendContext[T any](ctx context.Context, c chan&lt;- T, item T) error</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="Merge"></a><samp>func <a href="#Merge">Merge</a>[T any](out chan&lt;- T, in ...&lt;-chan T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/chans/chans.go#L37">src</a></small></sub></h3>

Merge sends all values from all in channels to out.

Merge blocks until all ins have closed and all values have been sent. It does not close out.


#### Example 
```go
{
	a := make(chan int)
	go func() {
		a <- 0
		a <- 1
		a <- 2
		close(a)
	}()
	b := make(chan int)
	go func() {
		b <- 5
		b <- 6
		b <- 7
		b <- 8
		close(b)
	}()

	out := make(chan int)
	done := make(chan struct{})
	go func() {
		for i := range out {
			fmt.Println(i)
		}
		close(done)
	}()

	chans.Merge(out, a, b)
	close(out)
	<-done

}
```

Unordered output:
```text
0
1
2
5
6
7
8
```
<h3><a id="RecvContext"></a><samp>func <a href="#RecvContext">RecvContext</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, c &lt;-chan T) (T, bool, error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/chans/chans.go#L24">src</a></small></sub></h3>

RecvContext attempts to receive from channel c. If c is closed before or during, returns (_,
false, nil). If ctx expires before or during, returns (_, _, ctx.Err()).


<h3><a id="Replicate"></a><samp>func <a href="#Replicate">Replicate</a>[T any](src &lt;-chan T, dsts ...chan&lt;- T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/chans/chans.go#L142">src</a></small></sub></h3>

Replicate sends all values sent to src to every channel in dsts.

Replicate blocks until src is closed and all values have been sent to all dsts. It does not close
dsts.


#### Example 
```go
{
	in := make(chan int)
	go func() {
		in <- 0
		in <- 1
		in <- 2
		in <- 3
		close(in)
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	a := make(chan int)
	go func() {
		for i := range a {
			fmt.Println(i * 2)
		}
		wg.Done()
	}()

	b := make(chan int)
	go func() {
		x := 0
		for i := range b {
			x += i
			fmt.Println(x)
		}
		wg.Done()
	}()

	chans.Replicate(in, a, b)
	close(a)
	close(b)
	wg.Wait()

}
```

Unordered output:
```text
0
2
4
6
0
1
3
6
```
<h3><a id="SendContext"></a><samp>func <a href="#SendContext">SendContext</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, c chan&lt;- T, item T) error</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/chans/chans.go#L13">src</a></small></sub></h3>

SendContext sends item on channel c and returns nil, unless ctx expires in which case it returns
ctx.Err().


## Types

This section is empty.

