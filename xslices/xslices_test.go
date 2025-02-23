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

func TestIntersect(t *testing.T) {
	testCases := map[string]struct {
		in   [][]int
		want []int
	}{
		"nil":            {},
		"empty 1":        {in: [][]int{}},
		"empty 2":        {in: [][]int{{}}},
		"empty 3":        {in: [][]int{{}, {}, {}}},
		"empty inter 1":  {in: [][]int{{1, 2}, {1, 2}, {}}},
		"empty inter 2":  {in: [][]int{{1, 2}, {3, 4}, {5, 6}}},
		"empty inter 3":  {in: [][]int{{1, 1}, {2, 2}}},
		"single inter":   {in: [][]int{{1, 2, 3}, {1, 2}, {2, 3}}, want: []int{2}},
		"multiple inter": {in: [][]int{{1, 2, 3}, {1, 2}, {1, 2}}, want: []int{1, 2}},
		"complete inter": {in: [][]int{{1, 2, 3}, {3, 1, 2}, {2, 3, 1}}, want: []int{1, 2, 3}},
		"repeated inter": {in: [][]int{{1, 1, 2}, {1, 1, 1}, {1, 1, 3}}, want: []int{1, 1}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			require2.ElementsEqual(t, tc.want, Intersect(tc.in...))
		})
	}
}
