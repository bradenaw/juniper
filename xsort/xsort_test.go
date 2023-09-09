package xsort_test

import (
	"cmp"
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

func TestMergeSlices(t *testing.T) {
	check := func(in ...[]int) {
		var all []int
		for i := range in {
			require2.True(t, slices.IsSorted(in[i]))
			all = append(all, in[i]...)
		}
		merged := xsort.MergeSlices(
			cmp.Less[int],
			nil,
			in...,
		)
		require2.True(t, slices.IsSorted(merged))
		require2.ElementsMatch(t, all, merged)
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
			slices.Sort(bs[i])
		}

		expected := append([]byte{}, b...)
		xsort.Slice(expected, cmp.Less[byte])

		merged := xsort.Merge(
			cmp.Less[byte],
			iterator.Collect(
				iterator.Map(
					iterator.Slice(bs),
					func(b []byte) iterator.Iterator[byte] { return iterator.Slice(b) },
				),
			)...,
		)

		require2.SlicesEqual(t, expected, iterator.Collect(merged))
	})
}

func ExampleSearch() {
	x := []string{"a", "f", "h", "i", "p", "z"}

	fmt.Println(xsort.Search(x, cmp.Less[string], "h"))
	fmt.Println(xsort.Search(x, cmp.Less[string], "k"))

	// Output:
	// 2
	// 4
}

func ExampleMerge() {
	listOne := []string{"a", "f", "p", "x"}
	listTwo := []string{"b", "e", "o", "v"}
	listThree := []string{"s", "z"}

	merged := xsort.Merge(
		cmp.Less[string],
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

	merged := xsort.MergeSlices(
		cmp.Less[string],
		nil,
		listOne,
		listTwo,
		listThree,
	)

	fmt.Println(merged)

	// Output:
	// [a b e f o p s v x z]
}

func ExampleMinK() {
	a := []int{7, 4, 3, 8, 2, 1, 6, 9, 0, 5}

	iter := iterator.Slice(a)
	min3 := xsort.MinK(cmp.Less[int], iter, 3)
	fmt.Println(min3)

	iter = iterator.Slice(a)
	max3 := xsort.MinK(xsort.Reverse(cmp.Less[int]), iter, 3)
	fmt.Println(max3)

	// Output:
	// [0 1 2]
	// [9 8 7]
}
