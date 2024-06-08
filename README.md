# Juniper

[![Go Reference](https://pkg.go.dev/badge/github.com/bradenaw/juniper.svg)](https://pkg.go.dev/github.com/bradenaw/juniper)
[![Go 1.19](https://github.com/bradenaw/juniper/actions/workflows/go1.19.yml/badge.svg)](https://github.com/bradenaw/juniper/actions/workflows/go1.19.yml)
[![Go 1.20](https://github.com/bradenaw/juniper/actions/workflows/go1.20.yml/badge.svg)](https://github.com/bradenaw/juniper/actions/workflows/go1.20.yml)
[![Go 1.21](https://github.com/bradenaw/juniper/actions/workflows/go1.21.yml/badge.svg)](https://github.com/bradenaw/juniper/actions/workflows/go1.21.yml)
[![Go 1.22](https://github.com/bradenaw/juniper/actions/workflows/go1.22.yml/badge.svg)](https://github.com/bradenaw/juniper/actions/workflows/go1.22.yml)
[![Fuzz](https://github.com/bradenaw/juniper/actions/workflows/fuzz.yml/badge.svg)](https://github.com/bradenaw/juniper/actions/workflows/fuzz.yml)

Juniper is a library of extensions to the Go standard library using generics, including containers,
iterators, and streams.

- `container/tree` contains a `Map` and `Set` that keep elements in sorted order. They are
  implemented using a B-tree, which performs better than a binary search tree.
- `container/deque` contains a double-ended queue implemented with a ring buffer.
- `container/xheap` contains a min-heap similar to the standard library's `container/heap` but
  more ergonomic, along with a `PriorityQueue` that allows setting priorities by key.
- `container/xlist` contains a linked-list similar to the standard library's `container/list`, but
  type-safe.
- `xslices` contains some commonly-used slice operations, like `Chunk`, `Reverse`, `Clear`, and
  `Join`.
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
- You can probably guess what's in the packages `xerrors`, `xmath`, `xmath/xrand`, `xsync`, and
  `xtime`.

Packages that overlap directly with a standard library package are named the same but with an `x`
prefix for "extensions", e.g. `sort` and `xsort`.

See the [docs](https://pkg.go.dev/github.com/bradenaw/juniper) for more.
