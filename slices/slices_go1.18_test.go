//go:build go1.18

package slices_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/bradenaw/juniper/slices"
	"github.com/bradenaw/juniper/xmath"
)

func ExampleGrow() {
	x := make([]int, 0, 1)
	x = slices.Grow(x, 4)
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

func ExampleFilter() {
	x := []int{5, -9, -2, 1, -4, 8, 3}
	x = slices.Filter(x, func(value int) bool {
		return value > 0
	})
	fmt.Println(x)

	// Output:
	// [5 1 8 3]
}

func ExampleReverse() {
	x := []string{"a", "b", "c", "d", "e"}
	slices.Reverse(x)
	fmt.Println(x)

	// Output:
	// [e d c b a]
}

func ExampleInsert() {
	x := []string{"a", "b", "c", "d", "e"}
	x = slices.Insert(x, 3, "f", "g")
	fmt.Println(x)

	// Output:
	// [a b c f g d e]
}

func ExampleRemove() {
	x := []int{1, 2, 3, 4, 5}
	x = slices.Remove(x, 1, 2)
	fmt.Println(x)

	// Output:
	// [1 4 5]
}

func ExampleClear() {
	x := []int{1, 2, 3}
	slices.Clear(x)
	fmt.Println(x)

	// Output:
	// [0 0 0]
}

func ExampleClone() {
	x := []int{1, 2, 3}
	cloned := slices.Clone(x)
	fmt.Println(cloned)

	// Output:
	// [1 2 3]
}

func ExampleCompact() {
	x := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := slices.Compact(x)
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
	compacted := slices.CompactFunc(x, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(compacted)

	// Output:
	// [bank ghost yaw]
}

func ExampleEqual() {
	x := []string{"a", "b", "c"}
	y := []string{"a", "b", "c"}
	z := []string{"a", "b", "d"}

	fmt.Println(slices.Equal(x, y))
	fmt.Println(slices.Equal(x[:2], y))
	fmt.Println(slices.Equal(z, y))

	// Output:
	// true
	// false
	// false
}

func ExampleCount() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(slices.Count(x, "a"))

	// Output:
	// 3
}

func ExampleIndex() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(slices.Index(x, "b"))
	fmt.Println(slices.Index(x, "c"))

	// Output:
	// 1
	// -1
}

func ExampleLastIndex() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(slices.LastIndex(x, "a"))
	fmt.Println(slices.LastIndex(x, "c"))

	// Output:
	// 3
	// -1
}

func ExampleJoin() {
	joined := slices.Join(
		[]string{"a", "b", "c"},
		[]string{"x", "y"},
		[]string{"l", "m", "n", "o"},
	)

	fmt.Println(joined)

	// Output:
	// [a b c x y l m n o]
}

func ExamplePartition() {
	a := []int{11, 3, 4, 2, 7, 8, 0, 1, 14}

	slices.Partition(a, func(x int) bool { return x%2 == 0 })

	fmt.Println(a)

	// Output:
	// [11 3 1 7 2 8 0 4 14]
}

func FuzzPartition(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		test := func(x byte) bool { return x%2 == 0 }

		slices.Partition(b, test)

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

func ExampleUnique() {
	a := []string{"a", "b", "b", "c", "a", "b", "b", "c"}
	unique := slices.Unique(a)
	fmt.Println(unique)

	// Output:
	// [a b c]
}

func ExampleChunk() {
	a := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	chunks := slices.Chunk(a, 3)
	fmt.Println(chunks)

	// Output:
	// [[a b c] [d e f] [g h]]
}

func ExampleAny() {
	isOdd := func(x int) bool {
		return x%2 != 0
	}

	anyOdd := slices.Any([]int{2, 3, 4}, isOdd)
	fmt.Println(anyOdd)

	anyOdd = slices.Any([]int{2, 4, 6}, isOdd)
	fmt.Println(anyOdd)

	// Output:
	// true
	// false
}

func ExampleAll() {
	isOdd := func(x int) bool {
		return x%2 != 0
	}

	allOdd := slices.All([]int{1, 3, 5}, isOdd)
	fmt.Println(allOdd)

	allOdd = slices.All([]int{1, 3, 6}, isOdd)
	fmt.Println(allOdd)

	// Output:
	// true
	// false
}

func ExampleMap() {
	toHalfFloat := func(x int) float32 {
		return float32(x) / 2
	}

	a := []int{1, 2, 3}
	floats := slices.Map(a, toHalfFloat)
	fmt.Println(floats)

	// Output:
	// [0.5 1 1.5]
}

func ExampleRuns() {
	x := []int{2, 4, 0, 7, 1, 3, 9, 2, 8}

	parityRuns := slices.Runs(x, func(a, b int) bool {
		return a%2 == b%2
	})

	fmt.Println(parityRuns)

	// Output:
	// [[2 4 0] [7 1 3 9] [2 8]]
}

func ExampleReduce() {
	x := []int{3, 1, 2}

	sum := slices.Reduce(x, 0, func(x, y int) int { return x + y })
	fmt.Println(sum)

	min := slices.Reduce(x, math.MaxInt, xmath.Min[int])
	fmt.Println(min)

	// Output:
	// 6
	// 1
}
