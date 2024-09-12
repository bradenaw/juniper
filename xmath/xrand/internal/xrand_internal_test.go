package internal

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
)

type fuzzRand struct {
	t *testing.T
	b []byte
}

func (r *fuzzRand) IntN(n int) int {
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

		samp := NewSampler(r, k)
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
