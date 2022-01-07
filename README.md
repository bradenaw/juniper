# Juniper

Juniper is a library of extensions to the Go standard library. It is **very much experimental** and
nothing is guaranteed, use at your own risk.

For example:

- `container/tree` contains equivalents to the built-in `map` that keep elements in sorted order.
- `container/deque` contains a double-ended queue.
- `container/xheap` contains a min-heap similar to `container/heap` but easier to use, although less
  powerful.
- `slices` contains some commonly-used slice operations.
- `iterator` contains an iterator interface used by the containers, along with a few functions to
  manipulate them.
- `xsort` contains mostly equivalents to `sort`.

Packages that overlap directly with a standard library package are named the same but with an `x`
prefix for "extensions", e.g. `sort` and `xsort`.

Currently it's an experiment with Go 1.18's generics and how the standard library might evolve to
take advantage of them. Go itself is not adopting changes to the standard library with Go 1.18,
instead waiting until the Go community has gotten some experience working with them.

To play with it, you'll need to use [Go 1.18 Beta 1](https://go.dev/blog/go1.18beta1) or [build Go
from source](https://go.dev/doc/install/source) as Go 1.18 is not yet released.

A few functions do not require generics, and so this library still builds with Go 1.17 and below but
with a significantly smaller API.

# Notes on Generics So Far

## `interface` can't refer to itself in method definitions

`xsort` allows defining ordering for arbitrary types.

This would've been nice to be able to do:

```
package xsort

type Ordered interface {
    Less(other X) bool
}
```

What should `X` be? It could be `Ordered`, but then implementations of `Less` need to add a
type-cast.

Here's an example. This is how we'd like to be able to implement `Less` for types that can use `<`.
```
type OrderedByLessOperator[T constraints.Ordered] struct {
    T
}

func (o OrderedByLessOperator[T]) Less(other OrderedByLess[T]) {
    return o.T < other.T
}
```

Unfortunately, the only thing we can pick for `X` that makes it so `OrderedByLessOperator[T]`
implements `xsort.Ordered` is `OrderedByLessOperator[T]`, but then we can't use `xsort.Ordered` for
anything else.

We could do this instead:

```
package xsort

type Ordered interface {
    Less(other Ordered) bool 
}

type OrderedByLessOperator[T constraints.Ordered] struct {
    T
}

func (o OrderedByLessOperator[T]) Less(other Ordered) {
    return o.T < other.(OrderedByLess[T]).T
}
```

This has a nasty typecast and thus a panic that the type system can't save us from.

We really want `X` to be "the same type" like Rust's `Self`, but Go has no such feature.

This implementation works:

```
type Less[T] func(a, b T) bool

func OrderedLess[T constraints.Ordered](a, b T) bool {
	return a < b
}
```

This requires passing a `Less[T]` function pointer everywhere that it's needed. This removes the
interface boxing and unboxing from the above solution, so the type system works nicely, but it also
means that `less` can't get inlined because it isn't known at specialization time.

There is an alternative which is a little more awkward to work with.

```
type Ordering[T any] interface {
    Less(a, b T) bool
}

type NaturalOrder[T constraints.Ordered] struct{}

func (NaturalOrder[T]) Less(a, b T) bool {
	return a < b
}

func SortSlice[T any, L Ordering[T]](a []T) {
    // ...
}

a := []int{5, 3, 4}
SortSlice[int, NaturalOrder[int]](a)
```

This works, and I like that the method of ordering becomes a part of the type signature (e.g.
`tree.Set[int, xsort.Reverse[int, xsort.NaturalOrder[int]]]` tells you both that the set elements
are `int` but also they're in reverse order by `<`).

However, it feels odd to have this extra `struct{}` type definition to hang `Less` off of which we
always call on the zero value. It also is unfortunate that this mucks up type inference - Go likes
inferring all of the types or none of them, and so moving `Ordering` into the type parameter list
means type parameters can't ever get inferred. Also note it causes an explosion in type parameters,
because in order to name `Ordering[T]` you also must separately define `T any` - see how
`tree.Set[int, xsort.Reverse[int, xsort.NaturalOrder[int]]]` had to name `int` three times.
Theoretically, `tree.Set[xsort.Reverse[xsort.NaturalOrder[int]]]` should be sufficient, since the
inner `int` would imply `T=int` to `tree.Set`.

The one advantage of this solution is that `Less`'s concrete implementation is known at
specialization time, so it should be possible for the compiler to inline it.

## Methods can't be type-parameterized

This would've made combinators on `Iterator` and `Stream` a lot more ergonomic.

Here's an example that doesn't work:

```
type SimpleIterator[T any] interface {
    Next() (T, bool)
}

type Iterator[T any] struct {
    SimpleIterator[T]
}

func (iter Iterator[T]) Filter(keep func(T) bool) Iterator[T] {
    // ...
}

func (iter Iterator[T]) Map[U any](transform func(T) U) Iterator[U] {
    // ...
}
```

This would allow more natural chaining:

```
intIterator.Filter(func(x int) bool {
    return x % 2 == 0
}).Map(func(x int) float64{
    return float64(x) / 2
})
```

Unfortunately, this is disallowed because methods cannot be parametric. This line:
```
func (iter Iterator[T]) Map[U any](transform func(T) U) Iterator[U] {
```

fails with:
```
./prog.go:23:28: methods cannot have type parameters
./prog.go:23:29: invalid AST: method must have no type parameters
```

This requires slightly awkward reordering, which is the implementation I've landed on for now.
```
type Iterator[T any] interface {
	Next() (T, bool)
}

func Filter[T any](iter Iterator[T], keep func(T) bool) Iterator[T] {
    // ...
}

func Map[T any, U any](iter Iterator[T], f func(t T) U) Iterator[U] {
    // ...
}
```

In this model, the above example looks like this, which unfortunately reads inside-out rather than
left-to-right:
```
iterator.Map(
    iterator.Filter(
        intIterator,
        func(x int) bool { return x % 2 == 0 },
    ),
    func(x int) float64{ return float64(x) / 2 },
}
```
