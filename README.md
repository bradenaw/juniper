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
