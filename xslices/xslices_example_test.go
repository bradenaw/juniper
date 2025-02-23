package xslices_test

import (
	"bytes"
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
	s := []int{1, 2, 3}
	xslices.Clear(s)
	fmt.Println(s)

	// Output:
	// [0 0 0]
}

func ExampleChunk() {
	s := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	chunks := xslices.Chunk(s, 3)
	fmt.Println(chunks)

	// Output:
	// [[a b c] [d e f] [g h]]
}

func ExampleClone() {
	s := []int{1, 2, 3}
	cloned := xslices.Clone(s)
	fmt.Println(cloned)

	// Output:
	// [1 2 3]
}

func ExampleCompact() {
	s := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := xslices.Compact(s)
	fmt.Println(compacted)

	// Output:
	// [a b c a]
}

func ExampleCompactFunc() {
	s := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}
	compacted := xslices.CompactFunc(s, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(compacted)

	// Output:
	// [bank ghost yaw]
}

func ExampleCompactInPlace() {
	s := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := xslices.CompactInPlace(s)
	fmt.Println(compacted)

	// Output:
	// [a b c a]
}

func ExampleCompactInPlaceFunc() {
	s := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}
	compacted := xslices.CompactInPlaceFunc(s, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(compacted)

	// Output:
	// [bank ghost yaw]
}

func ExampleCount() {
	s := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.Count(s, "a"))

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

func ExampleEqualFunc() {
	x := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	y := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	z := [][]byte{[]byte("a"), []byte("b"), []byte("d")}

	fmt.Println(xslices.EqualFunc(x, y, bytes.Equal))
	fmt.Println(xslices.EqualFunc(x[:2], y, bytes.Equal))
	fmt.Println(xslices.EqualFunc(z, y, bytes.Equal))

	// Output:
	// true
	// false
	// false
}

func ExampleFill() {
	s := []int{1, 2, 3}
	xslices.Fill(s, 5)
	fmt.Println(s)

	// Output:
	// [5 5 5]
}

func ExampleFilter() {
	s := []int{5, -9, -2, 1, -4, 8, 3}
	s = xslices.Filter(s, func(value int) bool {
		return value > 0
	})
	fmt.Println(s)

	// Output:
	// [5 1 8 3]
}

func ExampleFilterInPlace() {
	s := []int{5, -9, -2, 1, -4, 8, 3}
	s = xslices.FilterInPlace(s, func(value int) bool {
		return value > 0
	})
	fmt.Println(s)

	// Output:
	// [5 1 8 3]
}

func ExampleGrow() {
	s := make([]int, 0, 1)
	s = xslices.Grow(s, 4)
	fmt.Println(len(s))
	fmt.Println(cap(s))
	s = append(s, 1)
	addr := &s[0]
	s = append(s, 2)
	fmt.Println(addr == &s[0])
	s = append(s, 3)
	fmt.Println(addr == &s[0])
	s = append(s, 4)
	fmt.Println(addr == &s[0])

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
	s := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.Index(s, "b"))
	fmt.Println(xslices.Index(s, "c"))

	// Output:
	// 1
	// -1
}

func ExampleIndexFunc() {
	s := []string{
		"blue",
		"green",
		"yellow",
		"gold",
		"red",
	}

	fmt.Println(xslices.IndexFunc(s, func(s string) bool {
		return strings.HasPrefix(s, "g")
	}))
	fmt.Println(xslices.IndexFunc(s, func(s string) bool {
		return strings.HasPrefix(s, "p")
	}))

	// Output:
	// 1
	// -1
}

func ExampleInsert() {
	s := []string{"a", "b", "c", "d", "e"}
	s = xslices.Insert(s, 3, "f", "g")
	fmt.Println(s)

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
	s := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.LastIndex(s, "a"))
	fmt.Println(xslices.LastIndex(s, "c"))

	// Output:
	// 3
	// -1
}

func ExampleLastIndexFunc() {
	s := []string{
		"blue",
		"green",
		"yellow",
		"gold",
		"red",
	}

	fmt.Println(xslices.LastIndexFunc(s, func(s string) bool {
		return strings.HasPrefix(s, "g")
	}))
	fmt.Println(xslices.LastIndexFunc(s, func(s string) bool {
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

	s := []int{1, 2, 3}
	floats := xslices.Map(s, toHalfFloat)
	fmt.Println(floats)

	// Output:
	// [0.5 1 1.5]
}

func ExamplePartition() {
	s := []int{11, 3, 4, 2, 7, 8, 0, 1, 14}

	xslices.Partition(s, func(x int) bool { return x%2 == 0 })

	fmt.Println(s)

	// Output:
	// [11 3 1 7 2 8 0 4 14]
}

func ExampleReduce() {
	s := []int{3, 1, 2}

	sum := xslices.Reduce(s, 0, func(x, y int) int { return x + y })
	fmt.Println(sum)

	min := xslices.Reduce(s, math.MaxInt, xmath.Min[int])
	fmt.Println(min)

	// Output:
	// 6
	// 1
}

func ExampleRemove() {
	s := []int{1, 2, 3, 4, 5}
	s = xslices.Remove(s, 1, 2)
	fmt.Println(s)

	// Output:
	// [1 4 5]
}

func ExampleRemoveUnordered() {
	s := []int{1, 2, 3, 4, 5}
	s = xslices.RemoveUnordered(s, 1, 1)
	fmt.Println(s)

	s = xslices.RemoveUnordered(s, 1, 2)
	fmt.Println(s)

	// Output:
	// [1 5 3 4]
	// [1 4]
}

func ExampleRepeat() {
	s := xslices.Repeat("a", 4)
	fmt.Println(s)

	// Output:
	// [a a a a]
}

func ExampleReverse() {
	s := []string{"a", "b", "c", "d", "e"}
	xslices.Reverse(s)
	fmt.Println(s)

	// Output:
	// [e d c b a]
}

func ExampleRuns() {
	s := []int{2, 4, 0, 7, 1, 3, 9, 2, 8}

	parityRuns := xslices.Runs(s, func(a, b int) bool {
		return a%2 == b%2
	})

	fmt.Println(parityRuns)

	// Output:
	// [[2 4 0] [7 1 3 9] [2 8]]
}

func ExampleShrink() {
	s := make([]int, 3, 15)
	s[0] = 0
	s[1] = 1
	s[2] = 2

	fmt.Println(s)
	fmt.Println(cap(s))

	s = xslices.Shrink(s, 0)

	fmt.Println(s)
	fmt.Println(cap(s))

	// Output:
	// [0 1 2]
	// 15
	// [0 1 2]
	// 3
}

func ExampleUnique() {
	s := []string{"a", "b", "b", "c", "a", "b", "b", "c"}
	unique := xslices.Unique(s)
	fmt.Println(unique)

	// Output:
	// [a b c]
}

func ExampleUniqueInPlace() {
	s := []string{"a", "b", "b", "c", "a", "b", "b", "c"}
	unique := xslices.UniqueInPlace(s)
	fmt.Println(unique)

	// Output:
	// [a b c]
}

func ExampleIntersect() {
	intersection := xslices.Intersect(
		[]string{"a", "b", "b", "c", "d"},
		[]string{"a", "b", "b", "c"},
		[]string{"a", "b", "b"},
	)

	fmt.Println(intersection)

	// Output:
	// [a b b]
}
