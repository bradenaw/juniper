package xsort

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzMerge(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte, n int, seed int64) {
		if len(b) == 0 {
			return
		}
		if n <= 0 {
			return
		}
		r := rand.New(rand.NewSource(seed))
		bs := make([][]byte, (n%len(b))+1)
		for i := range b {
			j := r.Intn(len(bs))
			bs[j] = append(bs[j], b[i])
		}
		for i := range bs {
			Slice(bs[i], OrderedLess[byte])
		}

		expected := append([]byte{}, b...)
		Slice(expected, OrderedLess[byte])

		merged := Merge(OrderedLess[byte], make([]byte, r.Intn(len(b))), bs...)

		require.Equal(t, expected, merged)
	})
}
