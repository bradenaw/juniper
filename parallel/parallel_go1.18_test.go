//go:build go1.18
// +build go1.18

package parallel

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}
	strs, err := Map(
		context.Background(),
		2, // parallelism
		ints,
		func(ctx context.Context, i int) (string, error) {
			return strconv.Itoa(i), nil
		},
	)
	require.NoError(t, err)
	require.Equal(t, strs, []string{"1", "2", "3", "4", "5"})
}
