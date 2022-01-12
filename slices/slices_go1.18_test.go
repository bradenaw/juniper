//go:build go1.18

package slices

import "testing"

func FuzzPartition(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		test := func(x byte) bool { return x%2 == 0 }

		Partition(b, test)

		for i := range b {
			if test(b[i]) {
				for j := i + 1; j < len(b); j++ {
					if !test(b[j]) {
						t.FailNow()
					}
				}
				break
			}
		}
	})
}
