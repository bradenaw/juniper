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
	rShuffle(defaultRand{}, a)
}

// RShuffle pseudo-randomizes the order of a.
func RShuffle[T any](r *rand.Rand, a []T) {
	rShuffle(r, a)
}

func rShuffle[T any](r randRand, a []T) {
	r.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}

// SampleIterator pseudo-randomly picks k items uniformly without replacement from iter.
//
// If iter yields fewer than k items, returns all of them.
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func SampleIterator[T any](iter iterator.Iterator[T], k int) []T {
	return rSampleIterator(defaultRand{}, iter, k)
}

// RSampleIterator pseudo-randomly picks k items uniformly without replacement from iter.
//
// If iter yields fewer than k items, returns all of them.
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func RSampleIterator[T any](r *rand.Rand, iter iterator.Iterator[T], k int) []T {
	return rSampleIterator(r, iter, k)
}

func rSampleIterator[T any](r randRand, iter iterator.Iterator[T], k int) []T {
	out := make([]T, k)
	i := 0
	sampler := sampleInner(r, k)
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
func SampleStream[T any](ctx context.Context, s stream.Stream[T], k int) ([]T, error) {
	return rSampleStream(ctx, defaultRand{}, s, k)
}

// RSampleStream pseudo-randomly picks k items uniformly without replacement from s.
//
// If s yields fewer than k items, returns all of them.
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func RSampleStream[T any](ctx context.Context, r *rand.Rand, s stream.Stream[T], k int) ([]T, error) {
	return rSampleStream(ctx, r, s, k)
}

func rSampleStream[T any](ctx context.Context, r randRand, s stream.Stream[T], k int) ([]T, error) {
	defer s.Close()

	out := make([]T, k)
	i := 0
	sampler := sampleInner(r, k)
Outer:
	for {
		next, replace := sampler()
		for {
			item, err := s.Next(ctx)
			if err == stream.End {
				break Outer
			} else if err != nil {
				return nil, err
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
	return out, nil
}

// SampleSlice pseudo-randomly picks k items uniformly without replacement from a.
//
// If len(a) < k, returns all items in a.
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func SampleSlice[T any](a []T, k int) []T {
	return rSampleSlice(defaultRand{}, a, k)
}

// RSampleSlice pseudo-randomly picks k items uniformly without replacement from a.
//
// If len(a) < k, returns all items in a.
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func RSampleSlice[T any](r *rand.Rand, a []T, k int) []T {
	return rSampleSlice(r, a, k)
}

func rSampleSlice[T any](r randRand, a []T, k int) []T {
	out := make([]T, k)
	sampler := sampleInner(r, k)
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
