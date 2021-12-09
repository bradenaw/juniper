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
