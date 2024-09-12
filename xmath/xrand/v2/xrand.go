//go:build go1.23

package xrand

import (
	"iter"
	"math/rand/v2"

	"github.com/bradenaw/juniper/xmath/xrand/internal"
)

type defaultRand struct{}

var _ internal.Rand = defaultRand{}

func (defaultRand) Float64() float64                   { return rand.Float64() }
func (defaultRand) IntN(n int) int                     { return rand.IntN(n) }
func (defaultRand) Shuffle(n int, swap func(int, int)) { rand.Shuffle(n, swap) }

// ShuffleSlice pseudo-randomizes the order of s.
func ShuffleSlice[S ~[]E, E any](s S) {
	internal.RShuffleSlice(defaultRand{}, s)
}

// RShuffleSlice pseudo-randomizes the order of s.
func RShuffleSlice[S ~[]E, E any](r *rand.Rand, s S) {
	internal.RShuffleSlice(r, s)
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
	return internal.RSample(r, n, k)
}

// SampleSlice pseudo-randomly picks k items uniformly without replacement from s.
//
// If len(a) < k, returns all items in s.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses O(k) time
// and space.
func SampleSlice[S ~[]E, E any](s S, k int) S {
	return internal.RSampleSlice(defaultRand{}, s, k)
}

// RSampleSlice pseudo-randomly picks k items uniformly without replacement from s.
//
// If len(a) < k, returns all items in s.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses O(k) time
// and space.
func RSampleSlice[S ~[]E, E any](r *rand.Rand, s S, k int) S {
	return internal.RSampleSlice(r, s, k)
}

// SampleIterator pseudo-randomly picks k items uniformly without replacement from seq.
//
// If seq yields fewer than k items, returns all of them.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
// linear in the length of seq but only O(k) extra space.
func SampleSeq[V any](seq iter.Seq[V], k int) []V {
	return rSampleSeq(defaultRand{}, seq, k)
}

// RSampleIterator pseudo-randomly picks k items uniformly without replacement from seq.
//
// If seq yields fewer than k items, returns all of them.
//
// Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
// linear in the length of seq but only O(k) extra space.
func RSampleSeq[V any](r *rand.Rand, seq iter.Seq[V], k int) []V {
	return rSampleSeq(r, seq, k)
}

func rSampleSeq[V any, R internal.Rand](r R, seq iter.Seq[V], k int) []V {
	out := make([]V, k)
	i := 0
	samp := internal.NewSampler(r, k)
	next, replace := samp.Next()
	for item := range seq {
		if i == next {
			out[replace] = item
			next, replace = samp.Next()
		}
		i++
	}
	if i < k {
		out = out[:i]
	}
	internal.RShuffleSlice(r, out)
	return out
}
