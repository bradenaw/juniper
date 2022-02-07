//go:build go1.18

package xrand

import (
	"context"
	"encoding/binary"
	"math"
	"math/rand"
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/stream"
)

type fuzzRand struct {
	t *testing.T
	b []byte
}

func (r *fuzzRand) Intn(n int) int {
	if len(r.b) < 4 {
		return 0
	}
	x := binary.BigEndian.Uint32(r.b[:4])
	r.b = r.b[4:]
	return int(x) % n
}
func (r *fuzzRand) Float64() float64 {
	if len(r.b) < 8 {
		return 0
	}
	x := binary.BigEndian.Uint64(r.b[:8])
	r.b = r.b[8:]
	out := float64(x) / math.MaxUint64
	if out == 1 {
		out = math.Nextafter(out, 0)
	}
	require2.GreaterOrEqual(r.t, out, float64(0))
	require2.Less(r.t, out, float64(1))
	r.t.Logf("%f", out)
	return out
}
func (r *fuzzRand) Shuffle(int, func(int, int)) {
	panic("unimplemented")
}

func FuzzSampleInner(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte, k int) {
		if k <= 0 {
			return
		}

		t.Logf("k %d", k)

		r := &fuzzRand{t, b}

		samp := newSampler(r, k)
		prev := 0
		for i := 0; i < 100; i++ {
			next, replace := samp.Next()
			t.Logf("%d: next %d replace %d", i, next, replace)
			if next == math.MaxInt {
				break
			}
			if i < k {
				require2.Equal(t, next, i)
				require2.Equal(t, replace, i)
			} else {
				require2.Greater(t, next, prev)
				require2.GreaterOrEqual(t, replace, 0)
				require2.Less(t, replace, k)
			}

			prev = next
		}
	})
}

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
