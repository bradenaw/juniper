//go:build go1.18

package iterator_test

import (
	"fmt"

	"github.com/bradenaw/juniper/iterator"
)

func ExampleIterator() {
	i := 0
	iter := iterator.FromNext(func() (int, bool) {
		if i >= 5 {
			return 0, false
		}
		item := i
		i++
		return item, true
	})

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
