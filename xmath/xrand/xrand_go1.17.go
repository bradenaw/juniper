//go:build !go1.18

package xrand

import (
	"math"
)

func rSample(r randRand, n int, k int) []int {
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

type sampler struct {
	i     int
	first bool
	w     float64
	k     int
	r     randRand
}

func newSampler(r randRand, k int) sampler {
	return sampler{
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
func (s *sampler) Next() (int, int) {
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
