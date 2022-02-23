package xslices

import (
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
)

func FuzzPartition(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		test := func(x byte) bool { return x%2 == 0 }

		t.Logf("in: %#v", b)
		t.Logf("in test: %#v", Map(b, test))
		idx := Partition(b, test)
		t.Logf("out: %#v", b)
		t.Logf("out test: %#v", Map(b, test))
		t.Logf("out idx: %d", idx)

		for i := 0; i < idx; i++ {
			require2.True(t, !test(b[i]))
		}
		for i := idx; i < len(b); i++ {
			require2.True(t, test(b[i]))
		}
	})
}

func FuzzRemoveUnordered(f *testing.F) {
	f.Fuzz(func(t *testing.T, l int, idx int, n int) {
		if l < 0 || l > 255 || idx < 0 || idx > l-1 || n < 0 || n > l-idx {
			return
		}

		t.Logf("l   = %d", l)
		t.Logf("idx = %d", idx)
		t.Logf("n   = %d", n)

		x := make([]int, l)
		expected := make([]int, 0, l)
		for i := range x {
			x[i] = i

			if !(i >= idx && i < idx+n) {
				expected = append(expected, i)
			}
		}

		actual := RemoveUnordered(Clone(x), idx, n)

		require2.ElementsMatch(t, expected, actual)
	})
}
