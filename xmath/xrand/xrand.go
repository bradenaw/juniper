package xrand

import (
	"math"
	"math/rand"
)

// Sample pseudo-randomly picks k ints uniformly without replacement from [0, n).
//
// If n < k, returns all ints in [0, n).
//
// The output is not in any particular order. If a pseudo-random order is desired, the output should
// be passed to Shuffle.
func Sample(r *rand.Rand, n int, k int) []int {
	out := make([]int, k)
	sampler := sampleInner(r.Float64, r.Intn, k)
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
	randFloat64 func() float64,
	randIntn func(int) int,
	k int,
) func() (int, int) {
	i := 0
	first := true
	w := math.Exp(math.Log(randFloat64()) / float64(k))
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
		skip := math.Floor(math.Log(randFloat64()) / math.Log(1-w))
		if math.IsInf(skip, 0) || math.IsNaN(skip) {
			return math.MaxInt, 0
		}
		i += int(skip) + 1
		w *= math.Exp(math.Log(randFloat64()) / float64(k))
		return i, randIntn(k)
	}
}
