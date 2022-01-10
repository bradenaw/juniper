//go:build go1.18

package xrand

import (
	"context"
	"math/rand"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/stream"
)

// Shuffle pseudo-randomizes the order of a.
func Shuffle[T any](a []T) {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}

// SampleIterator pseudo-randomly picks k items uniformly without replacement from iter.
//
// If iter yields fewer than k items, returns all of them.
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func SampleIterator[T any](r *rand.Rand, iter iterator.Iterator[T], k int) []T {
	out := make([]T, k)
	i := 0
	sampler := sampleInner(r.Float64, r.Intn, k)
	Outer:
	for {
		next, replace := sampler()
		for {
			item, ok := iter.Next()
			if !ok {
				break Outer
			}
			if i == next {
				out[replace] = item
				i++
				break
			}
			i++
		}
	}
	if i < k {
		out = out[:i]
	}
	return out
}

// SampleStream pseudo-randomly picks k items uniformly without replacement from s.
//
// If s yields fewer than k items, returns all of them.
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func SampleStream[T any](ctx context.Context, r *rand.Rand, s stream.Stream[T], k int) ([]T, error) {
	out := make([]T, k)
	i := 0
	sampler := sampleInner(r.Float64, r.Intn, k)
	Outer:
	for {
		next, replace := sampler()
		for {
			item, ok := s.Next(ctx)
			if !ok {
				break Outer
			}
			if i == next {
				out[replace] = item
				i++
				break
			}
			i++
		}
	}
	if i < k {
		out = out[:i]
	}
	return out, s.Close()
}

// SampleSlice pseudo-randomly picks k items uniformly without replacement from a.
//
// If len(a) < k, returns all items in a.
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func SampleSlice[T any](r *rand.Rand, a []T, k int) []T {
	out := make([]T, k)
	sampler := sampleInner(r.Float64, r.Intn, k)
	for {
		next, replace := sampler()
		if next >= len(a) {
			break
		}
		out[replace] = a[next]
	}
	if len(a) < k {
		out = out[:len(a)]
	}
	return out
}
