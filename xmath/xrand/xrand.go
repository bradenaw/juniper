package xrand

import (
	"math"
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
