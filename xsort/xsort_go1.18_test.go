//go:build go1.18

package xsort_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

func TestMergeSlices(t *testing.T) {
	check := func(in ...[]int) {
		var all []int
		for i := range in {
			require.True(t, xsort.SliceIsSorted[int, xsort.NaturalOrder[int]](in[i]))
			all = append(all, in[i]...)
		}
		merged := xsort.MergeSlices[int, xsort.NaturalOrder[int]](
			nil,
			in...,
		)
		require.True(t, xsort.SliceIsSorted[int, xsort.NaturalOrder[int]](merged))
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
			xsort.Slice[byte, xsort.NaturalOrder[byte]](bs[i])
		}

		expected := append([]byte{}, b...)
		xsort.Slice[byte, xsort.NaturalOrder[byte]](expected)

		merged := xsort.Merge[byte, xsort.NaturalOrder[byte]](
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

func ExampleSearch() {
	x := []string{"a", "f", "h", "i", "p", "z"}

	fmt.Println(xsort.Search[string, xsort.NaturalOrder[string]](x, "h"))
	fmt.Println(xsort.Search[string, xsort.NaturalOrder[string]](x, "k"))

	// Output:
	// 2
	// 4
}

func ExampleSlice() {
	x := []int{3, 5, 1, 4, 2}
	xsort.Slice[int, xsort.NaturalOrder[int]](x)
	fmt.Println(x)

	// Output:
	// [1 2 3 4 5]
}

func ExampleMerge() {
	listOne := []string{"a", "f", "p", "x"}
	listTwo := []string{"b", "e", "o", "v"}
	listThree := []string{"s", "z"}

	merged := xsort.Merge[string, xsort.NaturalOrder[string]](
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

	merged := xsort.MergeSlices[string, xsort.NaturalOrder[string]](
		nil,
		listOne,
		listTwo,
		listThree,
	)

	fmt.Println(merged)

	// Output:
	// [a b e f o p s v x z]
}
