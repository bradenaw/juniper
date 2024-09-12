package xrand

import (
	"context"
	"math"
	"math/rand"
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/stream"
)

func stddev(a []int) float64 {
	m := mean(a)
	sumSquaredDeviation := float64(0)
	for i := range a {
		deviation := m - float64(a[i])
		sumSquaredDeviation += (deviation * deviation)
	}
	return math.Sqrt(sumSquaredDeviation / float64(len(a)))
}

func mean(a []int) float64 {
	sum := 0
	for i := range a {
		sum += a[i]
	}
	return float64(sum) / float64(len(a))
}

// f must return the same as Sample(r, 20, 5).
func testSample(t *testing.T, f func(r *rand.Rand) []int) {
	r := rand.New(rand.NewSource(0))

	counts := make([]int, 20)

	for i := 0; i < 10000; i++ {
		sample := f(r)
		for _, item := range sample {
			counts[item]++
		}
	}
	m := mean(counts)

	t.Logf("counts        %#v", counts)
	t.Logf("stddev        %#v", stddev(counts))
	t.Logf("stddev / mean %#v", stddev(counts)/m)

	// There's certainly a better statistical test than this, but I haven't bothered to break out
	// the stats book yet.
	require2.InDelta(t, 0.02, stddev(counts)/m, 0.01)

}
func TestSample(t *testing.T) {
	testSample(t, func(r *rand.Rand) []int {
		return RSample(r, 20, 5)
	})
}

func TestSampleSlice(t *testing.T) {
	a := iterator.Collect(iterator.Counter(20))
	testSample(t, func(r *rand.Rand) []int {
		return RSampleSlice(r, a, 5)
	})
}

func TestSampleIterator(t *testing.T) {
	testSample(t, func(r *rand.Rand) []int {
		return RSampleIterator(r, iterator.Counter(20), 5)
	})
}

func TestSampleStream(t *testing.T) {
	testSample(t, func(r *rand.Rand) []int {
		out, err := RSampleStream(
			context.Background(),
			r,
			stream.FromIterator(iterator.Counter(20)),
			5,
		)
		require2.NoError(t, err)
		return out
	})
}
