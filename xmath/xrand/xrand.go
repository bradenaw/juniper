package xrand

import (
	"math/rand"
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

// Sample pseudo-randomly picks k ints uniformly without replacement from [0, n).
//
// If n < k, returns all ints in [0, n).
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func Sample(n int, k int) []int {
	return rSample(defaultRand{}, n, k)
}

// RSample pseudo-randomly picks k ints uniformly without replacement from [0, n).
//
// If n < k, returns all ints in [0, n).
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func RSample(r *rand.Rand, n int, k int) []int {
	return rSample(r, n, k)
}
