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

# Notes on Challenges with Generics So Far

## Ordering

`xsort` allows defining ordering for arbitrary types. This is used for sorting, searching, and
ordered data structures like `container/tree` and `container/xheap`. Every available option for
doing this is a little awkward.

### Option A: Less function type
```
type Less[T] func(a, b T) bool

func OrderedLess[T constraints.Ordered](a, b T) bool {
	return a < b
}

func SortByLen[T any](x [][]T) {
    xsort.Slice(x, func(a, b []T) bool {
        return len(a) < len(b)
    })
}
```

Pros:
- Easy to understand. Already the way the standard library `sort` does it.
- Allows short ad-hoc definitions for `Less`.

Cons:
- Ordering is not expressed in the type system.
- Requires passing a `Less[T]` function pointer everywhere that it's needed.
- `less` can't get inlined because its implementation isn't known at compile time.

### Option B: Ordering interface
```
type Ordering[T any] interface {
    Less(a, b T) bool
}

type NaturalOrder[T constraints.Ordered] struct{}

func (NaturalOrder[T]) Less(a, b T) bool {
	return a < b
}

func SortSlice[O Ordering[T], T any](a []T) {
    // ...
}

a := []int{5, 3, 4}
SortSlice[NaturalOrder[int]](a)
```

Pros:
- Ordering is a part of the type signature for an ordered container (e.g.
  `tree.Set[xsort.Reverse[xsort.NaturalOrder[int]]]` tells you both that the set elements are `int`
  but also they're in reverse order by `<`).
- `Less`'s concrete implementation is known at specialization time, so it should be possible for the
  compiler to inline it. `-gcflags=-m` does seem to suggest that it's trying to do this.

Cons:
- It feels odd to have this extra `struct{}` type definition to hang `Less` off of which we always
  call on the zero value.
- Defining an `Ordering` is more cumbersome than Option A. Anonymous, single-use `Ordering`s aren't
  possible.
- Type definitions feel a little verbose or clumsy, requiring some redundancy to say `[O
  Ordering[T], T any]`. Since Go is willing to infer suffixes of missing type parameters when
  calling a function, this is not as verbose as it could have been. e.g.
  `xsort.Slice[xsort.NaturalOrder[int]](x)` - `T` is inferred from `O`, so it can be left off.
  The type signatures still look confusing, and this inference is only done for functions, not
  types.
- Awkwardness abound when there are several type parameters, like for `tree.Map`. Naturally, the
  proper order is `[K, V]`. If we're adding `Ordering[K]`, then `[O Ordering[K], K any, V any]`.
  However, because `K` is inferrable from `O`, `[O, V, K]` would give us the proper shorthand. This
  allows `tree.NewMap[xsort.NaturalOrder[int], string]`, but `V` and `K` appearing in that order in
  the type parameter list is uncomfortable.


### Option C: Ordered interface
```
package xsort

type Ordered[T] interface {
    Less(other T) bool
}

type OrderedByLessOperator[T constraints.Ordered] struct {
    T
}

func (o OrderedByLessOperator[T]) Less(other T) {
    return o.T < other.T
}
```

Pros:
- This feels a little more Go-ish than option B, since there isn't this extra `struct{}` type to hold
the ordering function.

Cons:
- It has the same awkwardness as option B with having to pass both `Ordered[T]` and `T` in type
  parameters lists. Further, it requires an extra boxing/unboxing for usage.

## Chaining and Method Parameterization

This is [discussed in the proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods).

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
