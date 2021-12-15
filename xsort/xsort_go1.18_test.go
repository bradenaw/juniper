//go:build go1.18

package xsort

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/iterator"
)

func TestMergeSlices(t *testing.T) {
	check := func(in ...[]int) {
		var all []int
		for i := range in {
			require.True(t, SliceIsSorted(in[i], OrderedLess[int]))
			all = append(all, in[i]...)
		}
		merged := MergeSlices(
			OrderedLess[int],
			nil,
			in...,
		)
		require.True(t, SliceIsSorted(merged, OrderedLess[int]))
		require.ElementsMatch(t, all, merged)
	}

	check([]int{1, 2, 3})
	check(
		[]int{1, 2, 3},
		[]int{4, 5, 6},
	)
	check(
		[]int{1, 3, 5},
		[]int{2, 4, 6},
	)
	check(
		[]int{1, 12, 19, 27},
		[]int{2, 7, 13},
		[]int{},
		[]int{5},
	)
}

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

		merged := Merge(
			OrderedLess[byte],
			iterator.Collect(
				iterator.Map(
					iterator.Slice(bs),
					func(b []byte) iterator.Iterator[byte] { return iterator.Slice(b) },
				),
			)...,
		)

		require.Equal(t, expected, iterator.Collect(merged))
	})
}

func ExampleMerge() {
	listOne := []string{"a", "f", "p", "x"}
	listTwo := []string{"b", "e", "o", "v"}
	listThree := []string{"s", "z"}

	merged := Merge(
		OrderedLess[string],
		iterator.Slice(listOne),
		iterator.Slice(listTwo),
		iterator.Slice(listThree),
	)

	fmt.Println(iterator.Collect(merged))

	// Output:
	// [a b e f o p s v x z]
}

func ExampleMergeSlices() {
	listOne := []string{"a", "f", "p", "x"}
	listTwo := []string{"b", "e", "o", "v"}
	listThree := []string{"s", "z"}

	merged := MergeSlices(
		OrderedLess[string],
		nil,
		listOne,
		listTwo,
		listThree,
	)

	fmt.Println(merged)

	// Output:
	// [a b e f o p s v x z]
}
