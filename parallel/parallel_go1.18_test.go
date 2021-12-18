//go:build go1.18
// +build go1.18

package parallel

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}
	strs := Map(
		2, // parallelism
		ints,
		func(i int) string {
			return strconv.Itoa(i)
		},
	)
	require.Equal(t, strs, []string{"1", "2", "3", "4", "5"})
}
