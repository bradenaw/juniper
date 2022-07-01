// Package xrand contains extensions to the standard library package math/rand.
package xrand

import (
	"context"
	"math"
	"math/rand"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/stream"
)

type randRand interface {
	Float64() float64
	Intn(int) int
	Shuffle(int, func(int, int))
}

type defaultRand struct{}

func (defaultRand) Float64() float64                   { return rand.Float64() }
func (defaultRand) Intn(n int) int                     { return rand.Intn(n) }
func (defaultRand) Shuffle(n int, swap func(int, int)) { rand.Shuffle(n, swap) }

// Shuffle pseudo-randomizes the order of a.
func Shuffle[T any](a []T) {
	rShuffle(defaultRand{}, a)
}

// RShuffle pseudo-randomizes the order of a.
func RShuffle[T any](r *rand.Rand, a []T) {
	rShuffle(r, a)
}

func rShuffle[T any, R randRand](r R, a []T) {
	r.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}

// Sample pseudo-randomly picks k ints uniformly without replacement from [0, n).
//
// If n < k, returns all ints in [0, n).
//
// Requires O(k) time and space.
func Sample(n int, k int) []int {
	return rSample(defaultRand{}, n, k)
}

// RSample pseudo-randomly picks k ints uniformly without replacement from [0, n).
//
// If n < k, returns all ints in [0, n).
//
// Requires O(k) time and space.
func RSample(r *rand.Rand, n int, k int) []int {
	return rSample(r, n, k)
}

func rSample[R randRand](r R, n int, k int) []int {
	out := make([]int, k)
	samp := newSampler(r, k)
	for {
		next, replace := samp.Next()
		if next >= n {
			break
		}
		out[replace] = next
	}
	if n < k {
		out = out[:n]
	}
	rShuffle(r, out)
	return out
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
	return rSampleIterator(r, iter, k)
}

func rSampleIterator[T any, R randRand](r R, iter iterator.Iterator[T], k int) []T {
	out := make([]T, k)
	i := 0
	samp := newSampler(r, k)
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
	rShuffle(r, out)
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
	return rSampleStream(ctx, r, s, k)
}

func rSampleStream[T any, R randRand](
	ctx context.Context,
	r R,
	s stream.Stream[T],
	k int,
) ([]T, error) {
	defer s.Close()

	out := make([]T, k)
	i := 0
	samp := newSampler(r, k)
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
	rShuffle(r, out)
	return out, nil
}

// SampleSlice pseudo-randomly picks k items uniformly without replacement from a.
//
// If len(a) < k, returns all items in a.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses O(k) time
// and space.
func SampleSlice[T any](a []T, k int) []T {
	return rSampleSlice(defaultRand{}, a, k)
}

// RSampleSlice pseudo-randomly picks k items uniformly without replacement from a.
//
// If len(a) < k, returns all items in a.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses O(k) time
// and space.
func RSampleSlice[T any](r *rand.Rand, a []T, k int) []T {
	return rSampleSlice(r, a, k)
}

func rSampleSlice[T any, R randRand](r R, a []T, k int) []T {
	out := make([]T, k)
	samp := newSampler(r, k)
	for {
		next, replace := samp.Next()
		if next >= len(a) {
			break
		}
		out[replace] = a[next]
	}
	if len(a) < k {
		out = out[:len(a)]
	}
	rShuffle(r, out)
	return out
}

type sampler[R randRand] struct {
	i     int
	first bool
	w     float64
	k     int
	r     R
}

func newSampler[R randRand](r R, k int) sampler[R] {
	return sampler[R]{
		i:     0,
		first: true,
		w:     math.Exp(math.Log(r.Float64()) / float64(k)),
		k:     k,
		r:     r,
	}
}

// Returns (next, replace) such that next is always increasing, and that the input item at index
// next (if there is one) should replace the reservoir item at index replace.
//
// As such, for the first k iterations, always returns (i, i) to build the reservoir.
func (s *sampler[R]) Next() (int, int) {
	if s.i < s.k {
		j := s.i
		s.i++
		return j, j
	}
	if s.first && s.i == s.k {
		s.i--
		s.first = false
	}
	skip := math.Floor(math.Log(s.r.Float64()) / math.Log(1-s.w))
	if math.IsInf(skip, 0) || math.IsNaN(skip) {
		return math.MaxInt, 0
	}
	s.i += int(skip) + 1
	s.w *= math.Exp(math.Log(s.r.Float64()) / float64(s.k))
	return s.i, s.r.Intn(s.k)
}
