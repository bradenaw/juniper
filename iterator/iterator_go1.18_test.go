//go:build go1.18

package iterator

import (
	"fmt"
)

func ExampleIterator() {
	i := 0
	iter := New(func() (int, bool) {
		if i >= 5 {
			return 0, false
		}
		item := i
		i++
		return item, true
	})

	for iter.Next() {
		fmt.Println(iter.Item())
	}

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleChunk() {
	iter := Slice([]string{"a", "b", "c", "d", "e", "f", "g", "h"})

	chunked := Chunk(iter, 3)
	chunked.Next()
	fmt.Println(chunked.Item())
	chunked.Next()
	fmt.Println(chunked.Item())
	chunked.Next()
	fmt.Println(chunked.Item())

	// Output:
	// [a b c]
	// [d e f]
	// [g h]
}

func ExampleEqual() {
	fmt.Println(
		Equal(
			Slice([]string{"a", "b", "c"}),
			Slice([]string{"a", "b", "c"}),
		),
	)

	fmt.Println(
		Equal(
			Slice([]string{"a", "b", "c"}),
			Slice([]string{"a", "b", "c", "d"}),
		),
	)

	// Output:
	// true
	// false
}
