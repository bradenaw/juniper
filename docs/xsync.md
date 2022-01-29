# xsync
--
    import "."

Package xsync contains extensions to the standard library package sync.

## Usage

#### func  Lazy

```go
func Lazy[T any](f func() T) func() T
```
Lazy makes a lazily-initialized value. On first access, it uses f to create the
value. Later accesses all receive the same value.

#### type ContextCond

```go
type ContextCond struct {
	L sync.Locker
}
```

ContextCond is equivalent to sync.Cond, except its Wait function accepts a
context.Context.

ContextConds should not be copied after first use.

#### func  NewContextCond

```go
func NewContextCond(l sync.Locker) *ContextCond
```
NewContextCond returns a new ContextCond with l as its Locker.

#### func (*ContextCond) Broadcast

```go
func (c *ContextCond) Broadcast()
```
Broadcast wakes all goroutines blocked in Wait(), if there are any.

It is allowed but not required for the caller to hold c.L during the call.

#### func (*ContextCond) Signal

```go
func (c *ContextCond) Signal()
```
Signal wakes one goroutine blocked in Wait(), if there is any. No guarantee is
made as to which goroutine will wake.

It is allowed but not required for the caller to hold c.L during the call.

#### func (*ContextCond) Wait

```go
func (c *ContextCond) Wait(ctx context.Context) error
```
Wait is equivalent to sync.Cond.Wait, except it accepts a context.Context. If
the context expires before this goroutine is woken by Broadcast or Signal, it
returns ctx.Err() immediately. If an error is returned, does not reaquire c.L
before returning.

#### type Future

```go
type Future[T any] struct {
}
```

Future can be filled with a value exactly once. Many goroutines can concurrently
wait for it to be filled. After filling, Wait() immediately returns the value it
was filled with.

Futures must be created by NewFuture and should not be copied after first use.

#### func  NewFuture

```go
func NewFuture[T any]() *Future[T]
```
NewFuture returns a ready-to-use Future.

#### func (*BADRECV) Fill

```go
func (f *Future[T]) Fill(x T)
```
Fill fills f with value x. All active calls to Wait return x, and all future
calls to Wait return x immediately.

Panics if f has already been filled.

#### func (*BADRECV) Wait

```go
func (f *Future[T]) Wait() T
```
Wait waits for f to be filled with a value and returns it. Returns immediately
if f is already filled.

#### func (*BADRECV) WaitContext

```go
func (f *Future[T]) WaitContext(ctx context.Context) (T, error)
```
Wait waits for f to be filled with a value and returns it, or returns ctx.Err()
if ctx expires before this happens. Returns immediately if f is already filled.

#### type Group

```go
type Group struct {
}
```

Group manages a group of goroutines.

#### func  NewGroup

```go
func NewGroup(ctx context.Context) *Group
```
NewGroup returns a Group ready for use. The context passed to any of the f
functions will be a descendant of ctx.

#### func (*Group) Once

```go
func (g *Group) Once(f func(ctx context.Context))
```
Once calls f once from another goroutine.

#### func (*Group) Periodic

```go
func (g *Group) Periodic(
	interval time.Duration,
	jitter time.Duration,
	f func(ctx context.Context),
)
```
Periodic spawns a goroutine that calls f once per interval +/- jitter.

#### func (*Group) PeriodicOrTrigger

```go
func (g *Group) PeriodicOrTrigger(
	interval time.Duration,
	jitter time.Duration,
	f func(ctx context.Context),
) func()
```
PeriodicOrTrigger spawns a goroutine which calls f whenever the returned
function is called. If f is already running when triggered, f will run again
immediately when it finishes. Also calls f when it has been interval+/-jitter
since the last trigger.

#### func (*Group) Stop

```go
func (g *Group) Stop()
```
Stop cancels the context passed to spawned goroutines.

#### func (*Group) Trigger

```go
func (g *Group) Trigger(f func(ctx context.Context)) func()
```
Trigger spawns a goroutine which calls f whenever the returned function is
called. If f is already running when triggered, f will run again immediately
when it finishes.

#### func (*Group) Wait

```go
func (g *Group) Wait()
```
Wait cancels the context passed to any of the spawned goroutines and waits for
all spawned goroutines to exit.

It is not safe to call Wait concurrently with any other method on g.

#### type Map

```go
type Map[K comparable, V any] struct {
}
```

Map is a typesafe wrapper over sync.Map.

#### func (*BADRECV) Delete

```go
func (m *Map[K, V]) Delete(key K)
```

#### func (*BADRECV) Load

```go
func (m *Map[K, V]) Load(key K) (value V, ok bool)
```

#### func (*BADRECV) LoadAndDelete

```go
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool)
```

#### func (*BADRECV) LoadOrStore

```go
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool)
```

#### func (*BADRECV) Range

```go
func (m *Map[K, V]) Range(f func(key K, value V) bool)
```

#### func (*BADRECV) Store

```go
func (m *Map[K, V]) Store(key K, value V)
```

#### type Pool

```go
type Pool[T any] struct {
}
```

Pool is a typesafe wrapper over sync.Pool.

#### func  NewPool

```go
func NewPool[T any](new_ func() T) Pool[T]
```

#### func (*BADRECV) Get

```go
func (p *Pool[T]) Get() T
```

#### func (*BADRECV) Put

```go
func (p *Pool[T]) Put(x T)
```
