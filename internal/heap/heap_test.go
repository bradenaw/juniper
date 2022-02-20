package heap

import (
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
)

func TestParentChildren(t *testing.T) {
	require2.Equal(t, 0, parent(1))
	require2.Equal(t, 0, parent(2))
	left, right := children(0)
	require2.Equal(t, 1, left)
	require2.Equal(t, 2, right)

	require2.Equal(t, 1, parent(3))
	require2.Equal(t, 1, parent(4))
	left, right = children(1)
	require2.Equal(t, 3, left)
	require2.Equal(t, 4, right)

	require2.Equal(t, 2, parent(5))
	require2.Equal(t, 2, parent(6))
	left, right = children(2)
	require2.Equal(t, 5, left)
	require2.Equal(t, 6, right)
}
