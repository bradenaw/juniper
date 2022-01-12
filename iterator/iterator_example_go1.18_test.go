//go:build go1.18

package iterator_test

import (
	"fmt"
	"math"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath"
)

func ExampleIterator() {
	iter := iterator.Counter(5)

	for {
		item, ok := iter.Next()
		if !ok {
			break
		}
		fmt.Println(item)
	}

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleChunk() {
	iter := iterator.Slice([]string{"a", "b", "c", "d", "e", "f", "g", "h"})

	chunked := iterator.Chunk(iter, 3)
	item, _ := chunked.Next()
	fmt.Println(item)
	item, _ = chunked.Next()
	fmt.Println(item)
	item, _ = chunked.Next()
	fmt.Println(item)

	// Output:
	// [a b c]
	// [d e f]
	// [g h]
}

func ExampleCompact() {
	iter := iterator.Slice([]string{"a", "a", "b", "c", "c", "c", "a"})
	compacted := iterator.Compact(iter)
	fmt.Println(iterator.Collect(compacted))

	// Output:
	// [a b c a]
}

func ExampleCompactFunc() {
	iter := iterator.Slice([]string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	})
	compacted := iterator.CompactFunc(iter, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(iterator.Collect(compacted))

	// Output:
	// [bank ghost yaw]
}

func ExampleEqual() {
	fmt.Println(
		iterator.Equal(
			iterator.Slice([]string{"a", "b", "c"}),
			iterator.Slice([]string{"a", "b", "c"}),
		),
	)

	fmt.Println(
		iterator.Equal(
			iterator.Slice([]string{"a", "b", "c"}),
			iterator.Slice([]string{"a", "b", "c", "d"}),
		),
	)

	// Output:
	// true
	// false
}

func ExampleFilter() {
	iter := iterator.Slice([]int{1, 2, 3, 4, 5, 6})

	evens := iterator.Filter(iter, func(x int) bool { return x%2 == 0 })
	fmt.Println(iterator.Collect(evens))

	// Output:
	// [2 4 6]
}

func ExampleFlatten() {
	iter := iterator.Slice([]iterator.Iterator[int]{
		iterator.Slice([]int{0, 1, 2}),
		iterator.Slice([]int{3, 4, 5, 6}),
		iterator.Slice([]int{7}),
	})

	all := iterator.Flatten(iter)

	fmt.Println(iterator.Collect(all))

	// Output:
	// [0 1 2 3 4 5 6 7]
}

func ExampleFirst() {
	iter := iterator.Slice([]string{"a", "b", "c", "d", "e"})

	first3 := iterator.First(iter, 3)
	fmt.Println(iterator.Collect(first3))

	// Output:
	// [a b c]
}

func ExampleJoin() {
	iter := iterator.Join(
		iterator.Counter(3),
		iterator.Counter(5),
		iterator.Counter(2),
	)

	fmt.Println(iterator.Collect(iter))

	// Output:
	// [0 1 2 0 1 2 3 4 0 1]
}

func ExampleLast() {
	iter := iterator.Counter(10)

	last3 := iterator.Last(iter, 3)
	fmt.Println(last3)

	iter = iterator.Counter(2)
	last3 = iterator.Last(iter, 3)
	fmt.Println(last3)

	// Output:
	// [7 8 9]
	// [0 1]
}

func ExampleRuns() {
	iter := iterator.Slice([]int{2, 4, 0, 7, 1, 3, 9, 2, 8})

	parityRuns := iterator.Runs(iter, func(a, b int) bool {
		return a%2 == b%2
	})
	fmt.Println(iterator.Collect(iterator.Map(parityRuns, iterator.Collect[int])))

	// Output:
	// [[2 4 0] [7 1 3 9] [2 8]]
}

func ExampleReduce() {
	x := []int{3, 1, 2}

	iter := iterator.Slice(x)
	sum := iterator.Reduce(iter, 0, func(x, y int) int { return x + y })
	fmt.Println(sum)

	iter = iterator.Slice(x)
	min := iterator.Reduce(iter, math.MaxInt, xmath.Min[int])
	fmt.Println(min)

	// Output:
	// 6
	// 1
}

func ExampleRepeat() {
	iter := iterator.Repeat("a", 4)
	fmt.Println(iterator.Collect(iter))

	// Output:
	// [a a a a]
}
