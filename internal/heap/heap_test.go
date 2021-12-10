//go:build go1.18

package heap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParentChildren(t *testing.T) {
	require.Equal(t, 0, parent(1))
	require.Equal(t, 0, parent(2))
	left, right := children(0)
	require.Equal(t, 1, left)
	require.Equal(t, 2, right)

	require.Equal(t, 1, parent(3))
	require.Equal(t, 1, parent(4))
	left, right = children(1)
	require.Equal(t, 3, left)
	require.Equal(t, 4, right)

	require.Equal(t, 2, parent(5))
	require.Equal(t, 2, parent(6))
	left, right = children(2)
	require.Equal(t, 5, left)
	require.Equal(t, 6, right)
}
