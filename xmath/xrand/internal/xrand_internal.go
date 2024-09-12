package internal

import (
	"math"
)

type Rand interface {
	IntN(n int) int
	Float64() float64
	Shuffle(n int, swap func(int, int))
}

func RShuffleSlice[S ~[]E, R Rand, E any](r R, s S) {
	r.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

func RSample[R Rand](r R, n int, k int) []int {
	out := make([]int, k)
	samp := NewSampler(r, k)
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
	RShuffleSlice(r, out)
	return out
}

func RSampleSlice[T any, R Rand](r R, a []T, k int) []T {
	out := make([]T, k)
	samp := NewSampler(r, k)
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
	RShuffleSlice(r, out)
	return out
}

type Sampler[R Rand] struct {
	i     int
	first bool
	w     float64
	k     int
	r     R
}

func NewSampler[R Rand](r R, k int) Sampler[R] {
	return Sampler[R]{
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
func (s *Sampler[R]) Next() (int, int) {
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
	return s.i, s.r.IntN(s.k)
}
