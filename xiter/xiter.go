//go:build go1.23

package xiter

import (
	"iter"
)

func Chunks[V comparable](seq iter.Seq[V], n int) iter.Seq[[]V] {
	chunk := make([]V, 0, n)
	return func(yield func([]V) bool) {
		for x := range seq {
			chunk = append(chunk, x)
			if len(chunk) == n {
				if !yield(chunk) {
					return
				}
				chunk = chunk[:0]
			}
		}
		if len(chunk) > 0 {
			yield(chunk)
		}
	}
}

func Compact[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var prev V
		hasPrev := false
		for v := range seq {
			if !hasPrev || v != prev {
				if !yield(v) {
					break
				}
			}
			prev = v
			hasPrev = true
		}
	}
}

func CompactFunc[V comparable](seq iter.Seq[V], eq func(V, V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		var prev V
		hasPrev := false
		for v := range seq {
			if !hasPrev || !eq(v, prev) {
				if !yield(v) {
					break
				}
			}
			prev = v
			hasPrev = true
		}
	}
}

func Counter(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			yield(i)
		}
	}
}

func Empty[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {}
}

func Equal[V comparable](seqs ...iter.Seq[V]) bool {
	nexts := make([]func() (V, bool), len(seqs))
	for i, seq := range seqs {
		var stop func()
		nexts[i], stop = iter.Pull(seq)
		defer stop()
	}
	for {
		v, ok := nexts[0]()

		for _, next := range nexts[1:] {
			v2, ok2 := next()
			if ok != ok2 || v != v2 {
				return false
			}
		}
		if !ok {
			break
		}
	}
	return true
}

func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !keep(v) {
				continue
			}
			if !yield(v) {
				break
			}
		}
	}
}

func First[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if !yield(v) {
				break
			}
			i++
			if i >= n {
				break
			}
		}
	}
}

func Flatten[V any](seq iter.Seq[iter.Seq[V]]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for seq2 := range seq {
			for v := range seq2 {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Join[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Map[V any, U any](seq iter.Seq[V], f func(V) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

func Last[V any](seq iter.Seq[V], n int) []V {
	buf := make([]V, n)
	i := 0
	for v := range seq {
		buf[i%n] = v
		i++
	}
	if i < n {
		return buf[:i]
	}
	for j := i % n; j < n; j++ {
		buf[j], buf[j-i%n] = buf[j-i%n], buf[j]
	}
	return buf
}

func One[V any](seq iter.Seq[V]) (V, bool) {
	var zero V
	var out V
	seen := false
	for v := range seq {
		if seen {
			return zero, false
		}
		out = v
		seen = true
	}
	return out, seen

}

func Repeat[V any](v V, n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _ = range n {
			if !yield(v) {
				return
			}
		}
	}
}

func Runs[V any](seq iter.Seq[V], same func(a, b V) bool) iter.Seq[iter.Seq[V]] {
	return func(yield func(iter.Seq[V]) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		prev, ok := next()
		if !ok {
			return
		}
		hasMore := true
		for hasMore {
			if !yield(func(yield2 func(V) bool) {
				// We always have one left over from either the start or from the previous run,
				// since we had to look at it to decide that either we should start at all or that
				// we need to move on to the next run.
				if !yield2(prev) {
					return
				}
				broken := false
				for {
					var curr V
					curr, hasMore = next()
					prev = curr
					if !hasMore {
						break
					}
					if !same(prev, curr) {
						break
					}
					if !broken {
						// We can't actually break here because we still have to find the start of
						// the next run, but we can stop emitting like they asked.
						broken = !yield2(curr)
					}
				}
			}) {
				break
			}
		}
	}
}

func While[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !f(v) {
				break
			}
			if !yield(v) {
				break
			}
		}
	}
}
