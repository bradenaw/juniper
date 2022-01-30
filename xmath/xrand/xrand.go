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

func rSample(r randRand, n int, k int) []int {
	out := make([]int, k)
	sampler := sampleInner(r, k)
	for {
		next, replace := sampler()
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

// Returns a function that, when called, does one iteration for a reservoir sample.
//
// Returns (next, replace) such that next is always increasing, and that the input item at index
// next (if there is one) should replace the reservoir item at index replace.
//
// As such, for the first k iterations, always returns (i, i) to build the reservoir.
func sampleInner(
	r randRand,
	k int,
) func() (int, int) {
	i := 0
	first := true
	w := math.Exp(math.Log(r.Float64()) / float64(k))
	return func() (int, int) {
		if i < k {
			j := i
			i++
			return j, j
		}
		if first && i == k {
			i--
			first = false
		}
		skip := math.Floor(math.Log(r.Float64()) / math.Log(1-w))
		if math.IsInf(skip, 0) || math.IsNaN(skip) {
			return math.MaxInt, 0
		}
		i += int(skip) + 1
		w *= math.Exp(math.Log(r.Float64()) / float64(k))
		return i, r.Intn(k)
	}
}
