# Juniper

[![Go Reference](https://pkg.go.dev/badge/github.com/bradenaw/juniper.svg)](https://pkg.go.dev/github.com/bradenaw/juniper) (Note: go 1.17 features only until pkg.go.dev builds with 1.18)

Juniper is a library of extensions to the Go standard library using generics, including containers.

- `container/tree` contains a `Map` and `Set` that keep elements in sorted order. They are
  implemented using a B-tree, which performs better than a binary search tree.
- `container/deque` contains a double-ended queue implemented with a ring buffer.
- `container/xheap` contains a min-heap similar to the standard library's `container/heap` but
  more ergonomic, along with a `PriorityQueue` that allows setting priorities by key.
- `container/xlist` contains a linked-list similar to the standard library's `container/list`, but
  type-safe.
- `slices` contains some commonly-used slice operations, like `Insert`, `Remove`, `Chunk`, `Filter`,
  and `Compact`.
- `iterator` contains an iterator interface used by the containers, along with functions to
  manipulate them, like `Map`, `While`, and `Reduce`.
- `stream` contains a stream interface, which is an iterator that can fail. Useful for iterating
  over collections that require I/O. It has most of the same combinators as `iterator`, plus some
  extras like `Pipe` and `Batch`.
- `parallel` contains some shorthand for common uses of goroutines to process slices, iterators, and
  streams in parallel, like `parallel.MapStream`.
- `xsort` contains extensions to the standard library package `sort`. Notably, it also has the
  definition for `xsort.Less`, which is how custom orderings can be defined for sorting and also for
  ordered collections like from `container/tree`.
- You can probably guess what's in the packages `maps`, `sets`, `xerrors`, `xmath`, `xmath/xrand`,
  `xsync`, and `xtime`.

Packages that overlap directly with a standard library package are named the same but with an `x`
prefix for "extensions", e.g. `sort` and `xsort`.

A few functions do not require generics (e.g. `parallel.Do` and `xmath/xrand.Sample`), and so this
library still builds with Go 1.17 and below but with a significantly smaller API. Note that
pkg.go.dev still builds with 1.17, so the documentation is incomplete.

## Status

Things should basically work. The container packages have been tested decently well using the [new
built-in coverage-based fuzzer](https://go.dev/doc/fuzz/) (it's a pleasure, by the way, other than
having to translate from the built-in fuzz argument types). `container/tree` has been benchmarked
and tweaked for some extra performance. It's far from hyper-optimized, but should be efficient
enough. Most of the simpler functions are tested only with their examples.

Since I no longer work at a megacorp running a huge global deployment of Go, I no longer have that
at my disposal to certify any of this as battle-hardened. However, the quality of code here is high
enough that I would've been comfortable using anything here in the systems that I worked on.
