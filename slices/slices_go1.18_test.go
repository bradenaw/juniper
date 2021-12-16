//go:build go1.18

package slices

import (
	"fmt"
)

func ExampleGrow() {
	x := make([]int, 0, 1)
	x = Grow(x, 4)
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
	x = Filter(x, func(value int) bool {
		return value > 0
	})
	fmt.Println(x)

	// Output:
	// [5 1 8 3]
}

func ExampleReverse() {
	x := []string{"a", "b", "c", "d", "e"}
	Reverse(x)
	fmt.Println(x)

	// Output:
	// [e d c b a]
}

func ExampleInsert() {
	x := []string{"a", "b", "c", "d", "e"}
	x = Insert(x, 3, "f", "g")
	fmt.Println(x)

	// Output:
	// [a b c f g d e]
}

func ExampleRemove() {
	x := []int{1, 2, 3, 4, 5}
	x = Remove(x, 1, 2)
	fmt.Println(x)

	// Output:
	// [1 4 5]
}

func ExampleClear() {
	x := []int{1, 2, 3}
	Clear(x)
	fmt.Println(x)

	// Output:
	// [0 0 0]
}

func ExampleClone() {
	x := []int{1, 2, 3}
	cloned := Clone(x)
	fmt.Println(cloned)

	// Output:
	// [1 2 3]
}

func ExampleCompact() {
	x := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := Compact(x)
	fmt.Println(compacted)

	// Output:
	// [a b c a]
}

func ExampleEqual() {
	x := []string{"a", "b", "c"}
	y := []string{"a", "b", "c"}
	z := []string{"a", "b", "d"}

	fmt.Println(Equal(x, y))
	fmt.Println(Equal(x[:2], y))
	fmt.Println(Equal(z, y))

	// Output:
	// true
	// false
	// false
}

func ExampleCount() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(Count(x, "a"))

	// Output:
	// 3
}

func ExampleIndex() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(Index(x, "b"))
	fmt.Println(Index(x, "c"))

	// Output:
	// 1
	// -1
}

func ExampleLastIndex() {
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(LastIndex(x, "a"))
	fmt.Println(LastIndex(x, "c"))

	// Output:
	// 3
	// -1
}

func ExampleJoin() {
	joined := Join(
		[]string{"a", "b", "c"},
		[]string{"x", "y"},
		[]string{"l", "m", "n", "o"},
	)

	fmt.Println(joined)

	// Output:
	// [a b c x y l m n o]
}
