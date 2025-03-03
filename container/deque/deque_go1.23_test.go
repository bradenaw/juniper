//go:build go1.23

package deque

import (
	"slices"
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
)

func TestValues(t *testing.T) {
	var deque Deque[string]

	require2.SeqsEqual(
		t,
		slices.Values([]string{}),
		deque.Values(),
	)

	deque.PushFront("c")
	deque.PushFront("b")
	deque.PushFront("a")
	deque.PushBack("d")
	deque.PushBack("e")
	deque.PushBack("f")

	require2.SeqsEqual(
		t,
		slices.Values([]string{"a", "b", "c", "d", "e", "f"}),
		deque.Values(),
	)

	deque.PopFront()

	require2.SeqsEqual(
		t,
		slices.Values([]string{"b", "c", "d", "e", "f"}),
		deque.Values(),
	)

	deque.PopBack()

	require2.SeqsEqual(
		t,
		slices.Values([]string{"b", "c", "d", "e"}),
		deque.Values(),
	)
}
