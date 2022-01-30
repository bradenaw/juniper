# `package xsync`

```
import "github.com/bradenaw/juniper/xsync"
```

# Overview

Package xsync contains extensions to the standard library package sync.


# Index

<samp><a href="#Lazy">func Lazy[T any](f func() T) func() T</a></samp>

<samp><a href="#ContextCond">type ContextCond</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewContextCond">func NewContextCond(l sync.Locker) *ContextCond</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Broadcast">func (c *ContextCond) Broadcast()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Signal">func (c *ContextCond) Signal()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Wait">func (c *ContextCond) Wait(ctx context.Context) error</a></samp>

<samp><a href="#Future">type Future</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewFuture">func NewFuture[T any]() *Future[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Fill">func (f *Future[T]) Fill(x T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Wait">func (f *Future[T]) Wait() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#WaitContext">func (f *Future[T]) WaitContext(ctx context.Context) (T, error)</a></samp>

<samp><a href="#Group">type Group</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewGroup">func NewGroup(ctx context.Context) *Group</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Once">func (g *Group) Once(f func(ctx context.Context))</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Periodic">func (g *Group) Periodic(
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	interval time.Duration,
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	jitter time.Duration,
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	f func(ctx context.Context),
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#PeriodicOrTrigger">func (g *Group) PeriodicOrTrigger(
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	interval time.Duration,
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	jitter time.Duration,
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	f func(ctx context.Context),
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;) func()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Stop">func (g *Group) Stop()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Trigger">func (g *Group) Trigger(f func(ctx context.Context)) func()</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Wait">func (g *Group) Wait()</a></samp>

<samp><a href="#Map">type Map</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Delete">func (m *Map[K, V]) Delete(key K)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Load">func (m *Map[K, V]) Load(key K) (value V, ok bool)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#LoadAndDelete">func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#LoadOrStore">func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Range">func (m *Map[K, V]) Range(f func(key K, value V) bool)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Store">func (m *Map[K, V]) Store(key K, value V)</a></samp>

<samp><a href="#Pool">type Pool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#NewPool">func NewPool[T any](new_ func() T) Pool[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Get">func (p *Pool[T]) Get() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Put">func (p *Pool[T]) Put(x T)</a></samp>


# Constants

This section is empty.

# Variables

This section is empty.

# Functions

<h2><a id="Lazy"></a><samp>func <a href="#Lazy">Lazy</a>[T any](f () T) () T</samp></h2>

Lazy makes a lazily-initialized value. On first access, it uses f to create the value. Later
accesses all receive the same value.


### Example 
```go
{
	var (
		expensive = Lazy(func() string {
			fmt.Println("doing expensive init")
			return "foo"
		})
	)

	fmt.Println(expensive())
	fmt.Println(expensive())

}
```

Output:
```text
doing expensive init
foo
foo
```
# Types

<h2><a id="ContextCond"></a><samp>type ContextCond</samp></h2>
```go
type ContextCond struct {
	L sync.Locker
	// contains filtered or unexported fields
}
```

ContextCond is equivalent to sync.Cond, except its Wait function accepts a context.Context.

ContextConds should not be copied after first use.


<h2><a id="NewContextCond"></a><samp>func NewContextCond(l <a href="https://pkg.go.dev/sync#Locker">sync.Locker</a>) *<a href="#ContextCond">ContextCond</a></samp></h2>

NewContextCond returns a new ContextCond with l as its Locker.


<h2><a id="Broadcast"></a><samp>func (c *<a href="#ContextCond">ContextCond</a>) Broadcast()</samp></h2>

Broadcast wakes all goroutines blocked in Wait(), if there are any.

It is allowed but not required for the caller to hold c.L during the call.


<h2><a id="Signal"></a><samp>func (c *<a href="#ContextCond">ContextCond</a>) Signal()</samp></h2>

Signal wakes one goroutine blocked in Wait(), if there is any. No guarantee is made as to which
goroutine will wake.

It is allowed but not required for the caller to hold c.L during the call.


<h2><a id="Wait"></a><samp>func (c *<a href="#ContextCond">ContextCond</a>) Wait(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>) error</samp></h2>

Wait is equivalent to sync.Cond.Wait, except it accepts a context.Context. If the context expires
before this goroutine is woken by Broadcast or Signal, it returns ctx.Err() immediately. If an
error is returned, does not reaquire c.L before returning.


<h2><a id="Future"></a><samp>type Future</samp></h2>
```go
type Future[T any] struct {
	// contains filtered or unexported fields
}
```

Future can be filled with a value exactly once. Many goroutines can concurrently wait for it to
be filled. After filling, Wait() immediately returns the value it was filled with.

Futures must be created by NewFuture and should not be copied after first use.


<h2><a id="NewFuture"></a><samp>func NewFuture[T any]() *<a href="#Future">Future</a>[T]</samp></h2>

NewFuture returns a ready-to-use Future.


<h2><a id="Fill"></a><samp>func (f *<a href="#Future">Future</a>[T]) Fill(x T)</samp></h2>

Fill fills f with value x. All active calls to Wait return x, and all future calls to Wait return
x immediately.

Panics if f has already been filled.


<h2><a id="Wait"></a><samp>func (f *<a href="#Future">Future</a>[T]) Wait() T</samp></h2>

Wait waits for f to be filled with a value and returns it. Returns immediately if f is already
filled.


<h2><a id="WaitContext"></a><samp>func (f *<a href="#Future">Future</a>[T]) WaitContext(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>) (T, error)</samp></h2>

Wait waits for f to be filled with a value and returns it, or returns ctx.Err() if ctx expires
before this happens. Returns immediately if f is already filled.


<h2><a id="Group"></a><samp>type Group</samp></h2>
```go
type Group struct {
	// contains filtered or unexported fields
}
```

Group manages a group of goroutines.


<h2><a id="NewGroup"></a><samp>func NewGroup(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>) *<a href="#Group">Group</a></samp></h2>

NewGroup returns a Group ready for use. The context passed to any of the f functions will be a
descendant of ctx.


<h2><a id="Once"></a><samp>func (g *<a href="#Group">Group</a>) Once(f (ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>))</samp></h2>

Once calls f once from another goroutine.


<h2><a id="Periodic"></a><samp>func (g *<a href="#Group">Group</a>) Periodic(interval <a href="https://pkg.go.dev/time#Duration">time.Duration</a>, jitter <a href="https://pkg.go.dev/time#Duration">time.Duration</a>, f (ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>))</samp></h2>

Periodic spawns a goroutine that calls f once per interval +/- jitter.


<h2><a id="PeriodicOrTrigger"></a><samp>func (g *<a href="#Group">Group</a>) PeriodicOrTrigger(interval <a href="https://pkg.go.dev/time#Duration">time.Duration</a>, jitter <a href="https://pkg.go.dev/time#Duration">time.Duration</a>, f (ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>)) ()</samp></h2>

PeriodicOrTrigger spawns a goroutine which calls f whenever the returned function is called.  If
f is already running when triggered, f will run again immediately when it finishes. Also calls f
when it has been interval+/-jitter since the last trigger.


<h2><a id="Stop"></a><samp>func (g *<a href="#Group">Group</a>) Stop()</samp></h2>

Stop cancels the context passed to spawned goroutines.


<h2><a id="Trigger"></a><samp>func (g *<a href="#Group">Group</a>) Trigger(f (ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>)) ()</samp></h2>

Trigger spawns a goroutine which calls f whenever the returned function is called. If f is
already running when triggered, f will run again immediately when it finishes.


<h2><a id="Wait"></a><samp>func (g *<a href="#Group">Group</a>) Wait()</samp></h2>

Wait cancels the context passed to any of the spawned goroutines and waits for all spawned
goroutines to exit.

It is not safe to call Wait concurrently with any other method on g.


<h2><a id="Map"></a><samp>type Map</samp></h2>
```go
type Map[K comparable, V any] struct {
	// contains filtered or unexported fields
}
```

Map is a typesafe wrapper over sync.Map.


<h2><a id="Delete"></a><samp>func (m *<a href="#Map">Map</a>[K, V]) Delete(key K)</samp></h2>



<h2><a id="Load"></a><samp>func (m *<a href="#Map">Map</a>[K, V]) Load(key K) (value V, ok bool)</samp></h2>



<h2><a id="LoadAndDelete"></a><samp>func (m *<a href="#Map">Map</a>[K, V]) LoadAndDelete(key K) (value V, loaded bool)</samp></h2>



<h2><a id="LoadOrStore"></a><samp>func (m *<a href="#Map">Map</a>[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool)</samp></h2>



<h2><a id="Range"></a><samp>func (m *<a href="#Map">Map</a>[K, V]) Range(f (key K, value V) bool)</samp></h2>



<h2><a id="Store"></a><samp>func (m *<a href="#Map">Map</a>[K, V]) Store(key K, value V)</samp></h2>



<h2><a id="Pool"></a><samp>type Pool</samp></h2>
```go
type Pool[T any] struct {
	// contains filtered or unexported fields
}
```

Pool is a typesafe wrapper over sync.Pool.


<h2><a id="NewPool"></a><samp>func NewPool[T any](new_ () T) <a href="#Pool">Pool</a>[T]</samp></h2>



<h2><a id="Get"></a><samp>func (p *<a href="#Pool">Pool</a>[T]) Get() T</samp></h2>



<h2><a id="Put"></a><samp>func (p *<a href="#Pool">Pool</a>[T]) Put(x T)</samp></h2>



