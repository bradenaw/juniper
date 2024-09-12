// Package xrand contains extensions to the standard library package math/rand.
package xrand

import (
	"context"
	"math/rand"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/stream"
	"github.com/bradenaw/juniper/xmath/xrand/internal"
)

type defaultRand struct{}

var _ internal.Rand = defaultRand{}

func (defaultRand) Float64() float64                   { return rand.Float64() }
func (defaultRand) IntN(n int) int                     { return rand.Intn(n) }
func (defaultRand) Shuffle(n int, swap func(int, int)) { rand.Shuffle(n, swap) }

type v1Rand struct{ *rand.Rand }

var _ internal.Rand = v1Rand{}

func (r v1Rand) IntN(n int) int                     { return r.Rand.Intn(n) }
func (r v1Rand) Float64() float64                   { return r.Rand.Float64() }
func (r v1Rand) Shuffle(n int, swap func(int, int)) { r.Rand.Shuffle(n, swap) }

// Shuffle pseudo-randomizes the order of a.
func Shuffle[T any](a []T) {
	internal.RShuffleSlice(defaultRand{}, a)
}

// RShuffle pseudo-randomizes the order of a.
func RShuffle[T any](r *rand.Rand, a []T) {
	internal.RShuffleSlice(v1Rand{Rand: r}, a)
}

// Sample pseudo-randomly picks k ints uniformly without replacement from [0, n).
//
// If n < k, returns all ints in [0, n).
//
// Requires O(k) time and space.
func Sample(n int, k int) []int {
	return internal.RSample(defaultRand{}, n, k)
}

// RSample pseudo-randomly picks k ints uniformly without replacement from [0, n).
//
// If n < k, returns all ints in [0, n).
//
// Requires O(k) time and space.
func RSample(r *rand.Rand, n int, k int) []int {
	return internal.RSample(v1Rand{r}, n, k)
}

// SampleIterator pseudo-randomly picks k items uniformly without replacement from iter.
//
// If iter yields fewer than k items, returns all of them.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
// linear in the length of iter but only O(k) extra space.
func SampleIterator[T any](iter iterator.Iterator[T], k int) []T {
	return rSampleIterator(defaultRand{}, iter, k)
}

// RSampleIterator pseudo-randomly picks k items uniformly without replacement from iter.
//
// If iter yields fewer than k items, returns all of them.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
// linear in the length of iter but only O(k) extra space.
func RSampleIterator[T any](r *rand.Rand, iter iterator.Iterator[T], k int) []T {
	return rSampleIterator(v1Rand{r}, iter, k)
}

func rSampleIterator[T any, R internal.Rand](r R, iter iterator.Iterator[T], k int) []T {
	out := make([]T, k)
	i := 0
	samp := internal.NewSampler(r, k)
Outer:
	for {
		next, replace := samp.Next()
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
	internal.RShuffleSlice(r, out)
	return out
}

// SampleStream pseudo-randomly picks k items uniformly without replacement from s.
//
// If s yields fewer than k items, returns all of them.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
// linear in the length of s but only O(k) extra space.
func SampleStream[T any](ctx context.Context, s stream.Stream[T], k int) ([]T, error) {
	return rSampleStream(ctx, defaultRand{}, s, k)
}

// RSampleStream pseudo-randomly picks k items uniformly without replacement from s.
//
// If s yields fewer than k items, returns all of them.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
// linear in the length of s but only O(k) extra space.
func RSampleStream[T any](
	ctx context.Context,
	r *rand.Rand,
	s stream.Stream[T],
	k int,
) ([]T, error) {
	return rSampleStream(ctx, v1Rand{Rand: r}, s, k)
}

func rSampleStream[T any, R internal.Rand](
	ctx context.Context,
	r R,
	s stream.Stream[T],
	k int,
) ([]T, error) {
	defer s.Close()

	out := make([]T, k)
	i := 0
	samp := internal.NewSampler(r, k)
Outer:
	for {
		next, replace := samp.Next()
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
	internal.RShuffleSlice(r, out)
	return out, nil
}

// SampleSlice pseudo-randomly picks k items uniformly without replacement from a.
//
// If len(a) < k, returns all items in a.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses O(k) time
// and space.
func SampleSlice[T any](a []T, k int) []T {
	return internal.RSampleSlice(defaultRand{}, a, k)
}

// RSampleSlice pseudo-randomly picks k items uniformly without replacement from a.
//
// If len(a) < k, returns all items in a.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses O(k) time
// and space.
func RSampleSlice[T any](r *rand.Rand, a []T, k int) []T {
	return internal.RSampleSlice(v1Rand{Rand: r}, a, k)
}
