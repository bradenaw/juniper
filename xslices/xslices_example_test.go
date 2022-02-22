package xslices_test

import (
	"fmt"
	"math"
	"strings"

	"github.com/bradenaw/juniper/xmath"
	"github.com/bradenaw/juniper/xslices"
)

func ExampleAll() {
	isOdd := func(x int) bool {
		return x%2 != 0
	}

	allOdd := xslices.All([]int{1, 3, 5}, isOdd)
	fmt.Println(allOdd)

	allOdd = xslices.All([]int{1, 3, 6}, isOdd)
	fmt.Println(allOdd)

	// Output:
	// true
	// false
}

func ExampleAny() {
	isOdd := func(x int) bool {
		return x%2 != 0
	}

	anyOdd := xslices.Any([]int{2, 3, 4}, isOdd)
	fmt.Println(anyOdd)

	anyOdd = xslices.Any([]int{2, 4, 6}, isOdd)
	fmt.Println(anyOdd)

	// Output:
	// true
	// false
}

func ExampleClear() {
	x := []int{1, 2, 3}
	xslices.Clear(x)
	fmt.Println(x)

	// Output:
	// [0 0 0]
}

func ExampleChunk() {
	a := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	chunks := xslices.Chunk(a, 3)
	fmt.Println(chunks)

	// Output:
	// [[a b c] [d e f] [g h]]
}

func ExampleClone() {
	x := []int{1, 2, 3}
	cloned := xslices.Clone(x)
	fmt.Println(cloned)

	// Output:
	// [1 2 3]
}

func ExampleCompact() {
	x := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := xslices.Compact(x)
	fmt.Println(compacted)

	// Output:
	// [a b c a]
}

func ExampleCompactFunc() {
	x := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}
	compacted := xslices.CompactFunc(x, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(compacted)

	// Output:
	// [bank ghost yaw]
}

func ExampleCompactInPlace() {
	x := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := xslices.CompactInPlace(x)
	fmt.Println(compacted)

	// Output:
	// [a b c a]
}

func ExampleCompactInPlaceFunc() {
	x := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}
	compacted := xslices.CompactInPlaceFunc(x, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(compacted)

	// Output:
	// [bank ghost yaw]
}

func ExampleCount() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.Count(x, "a"))

	// Output:
	// 3
}

func ExampleEqual() {
	x := []string{"a", "b", "c"}
	y := []string{"a", "b", "c"}
	z := []string{"a", "b", "d"}

	fmt.Println(xslices.Equal(x, y))
	fmt.Println(xslices.Equal(x[:2], y))
	fmt.Println(xslices.Equal(z, y))

	// Output:
	// true
	// false
	// false
}

func ExampleFill() {
	x := []int{1, 2, 3}
	xslices.Fill(x, 5)
	fmt.Println(x)

	// Output:
	// [5 5 5]
}

func ExampleFilter() {
	x := []int{5, -9, -2, 1, -4, 8, 3}
	x = xslices.Filter(x, func(value int) bool {
		return value > 0
	})
	fmt.Println(x)

	// Output:
	// [5 1 8 3]
}

func ExampleFilterInPlace() {
	x := []int{5, -9, -2, 1, -4, 8, 3}
	x = xslices.FilterInPlace(x, func(value int) bool {
		return value > 0
	})
	fmt.Println(x)

	// Output:
	// [5 1 8 3]
}

func ExampleGrow() {
	x := make([]int, 0, 1)
	x = xslices.Grow(x, 4)
	fmt.Println(len(x))
	fmt.Println(cap(x))
	x = append(x, 1)
	addr := &x[0]
	x = append(x, 2)
	fmt.Println(addr == &x[0])
	x = append(x, 3)
	fmt.Println(addr == &x[0])
	x = append(x, 4)
	fmt.Println(addr == &x[0])

	// Output:
	// 0
	// 4
	// true
	// true
	// true
}

func ExampleGroup() {
	words := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}

	groups := xslices.Group(words, func(s string) rune {
		return ([]rune(s))[0]
	})

	for firstChar, group := range groups {
		fmt.Printf("%c: %v\n", firstChar, group)
	}

	// Unordered output:
	// b: [bank beach]
	// g: [ghost goat group]
	// y: [yaw yew]
}

func ExampleIndex() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.Index(x, "b"))
	fmt.Println(xslices.Index(x, "c"))

	// Output:
	// 1
	// -1
}

func ExampleIndexFunc() {
	x := []string{
		"blue",
		"green",
		"yellow",
		"gold",
		"red",
	}

	fmt.Println(xslices.IndexFunc(x, func(s string) bool {
		return strings.HasPrefix(s, "g")
	}))
	fmt.Println(xslices.IndexFunc(x, func(s string) bool {
		return strings.HasPrefix(s, "p")
	}))

	// Output:
	// 1
	// -1
}

func ExampleInsert() {
	x := []string{"a", "b", "c", "d", "e"}
	x = xslices.Insert(x, 3, "f", "g")
	fmt.Println(x)

	// Output:
	// [a b c f g d e]
}

func ExampleJoin() {
	joined := xslices.Join(
		[]string{"a", "b", "c"},
		[]string{"x", "y"},
		[]string{"l", "m", "n", "o"},
	)

	fmt.Println(joined)

	// Output:
	// [a b c x y l m n o]
}

func ExampleLastIndex() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.LastIndex(x, "a"))
	fmt.Println(xslices.LastIndex(x, "c"))

	// Output:
	// 3
	// -1
}

func ExampleLastIndexFunc() {
	x := []string{
		"blue",
		"green",
		"yellow",
		"gold",
		"red",
	}

	fmt.Println(xslices.LastIndexFunc(x, func(s string) bool {
		return strings.HasPrefix(s, "g")
	}))
	fmt.Println(xslices.LastIndexFunc(x, func(s string) bool {
		return strings.HasPrefix(s, "p")
	}))

	// Output:
	// 3
	// -1
}

func ExampleMap() {
	toHalfFloat := func(x int) float32 {
		return float32(x) / 2
	}

	a := []int{1, 2, 3}
	floats := xslices.Map(a, toHalfFloat)
	fmt.Println(floats)

	// Output:
	// [0.5 1 1.5]
}

func ExamplePartition() {
	a := []int{11, 3, 4, 2, 7, 8, 0, 1, 14}

	xslices.Partition(a, func(x int) bool { return x%2 == 0 })

	fmt.Println(a)

	// Output:
	// [11 3 1 7 2 8 0 4 14]
}

func ExampleReduce() {
	x := []int{3, 1, 2}

	sum := xslices.Reduce(x, 0, func(x, y int) int { return x + y })
	fmt.Println(sum)

	min := xslices.Reduce(x, math.MaxInt, xmath.Min[int])
	fmt.Println(min)

	// Output:
	// 6
	// 1
}

func ExampleRemove() {
	x := []int{1, 2, 3, 4, 5}
	x = xslices.Remove(x, 1, 2)
	fmt.Println(x)

	// Output:
	// [1 4 5]
}

func ExampleRemoveUnordered() {
	x := []int{1, 2, 3, 4, 5}
	x = xslices.RemoveUnordered(x, 1, 1)
	fmt.Println(x)

	x = xslices.RemoveUnordered(x, 1, 2)
	fmt.Println(x)

	// Output:
	// [1 5 3 4]
	// [1 4]
}

func ExampleRepeat() {
	x := xslices.Repeat("a", 4)
	fmt.Println(x)

	// Output:
	// [a a a a]
}

func ExampleReverse() {
	x := []string{"a", "b", "c", "d", "e"}
	xslices.Reverse(x)
	fmt.Println(x)

	// Output:
	// [e d c b a]
}

func ExampleRuns() {
	x := []int{2, 4, 0, 7, 1, 3, 9, 2, 8}

	parityRuns := xslices.Runs(x, func(a, b int) bool {
		return a%2 == b%2
	})

	fmt.Println(parityRuns)

	// Output:
	// [[2 4 0] [7 1 3 9] [2 8]]
}

func ExampleShrink() {
	x := make([]int, 3, 15)
	x[0] = 0
	x[1] = 1
	x[2] = 2

	fmt.Println(x)
	fmt.Println(cap(x))

	x = xslices.Shrink(x, 0)

	fmt.Println(x)
	fmt.Println(cap(x))

	// Output:
	// [0 1 2]
	// 15
	// [0 1 2]
	// 3
}

func ExampleUnique() {
	a := []string{"a", "b", "b", "c", "a", "b", "b", "c"}
	unique := xslices.Unique(a)
	fmt.Println(unique)

	// Output:
	// [a b c]
}

func ExampleUniqueInPlace() {
	a := []string{"a", "b", "b", "c", "a", "b", "b", "c"}
	unique := xslices.UniqueInPlace(a)
	fmt.Println(unique)

	// Output:
	// [a b c]
}
