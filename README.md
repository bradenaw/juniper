# Juniper

Juniper is a library of extensions to the Go standard library.

For example:

- `container/tree` contains equivalents to the built-in `map` that keep elements in sorted order.
- `container/xheap` contains a min-heap similar to `container/heap` but easier to use, although less
  powerful.
- `slices` contains some functions for commonly-used operations.
- `iterator` contains an iterator interface used by the containers, along with a few functions to manipulate them.
- `xsort` contains mostly equivalents to `sort`, but type-safe with the advent of generics.

Packages that overlap directly with a standard library package are named the same but with an `x`
prefix for "extensions", e.g. `sort` and `xsort`.

Currently it's an experiment with Go 1.18's generics and how the standard library might evolve to
take advantage of them. Go itself is (smartly) not adopting changes to the standard library with Go
1.18, instead waiting until the Go community has gotten some experience working with them. So here's
some experience. :)

To play with it, you'll need to [build Go from source](https://go.dev/doc/install/source) as Go 1.18
is not yet released.

A few functions do not require generics, and so this library still builds with Go 1.17 and below but
with a significantly smaller API.
