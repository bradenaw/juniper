//go:build go1.18
// +build go1.18

package parallel

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/stream"
)

func TestMap(t *testing.T) {
	for _, parallelism := range []int{1, 2} {
		t.Run(fmt.Sprintf("parallelism=%d", parallelism), func(t *testing.T) {
			ints := []int{0, 1, 2, 3, 4}
			strs := Map(
				parallelism,
				ints,
				func(i int) string {
					return strconv.Itoa(i)
				},
			)
			require2.SlicesEqual(t, []string{"0", "1", "2", "3", "4"}, strs)
		})
	}
}

func TestMapContext(t *testing.T) {
	for _, parallelism := range []int{1, 2} {
		t.Run(fmt.Sprintf("parallelism=%d", parallelism), func(t *testing.T) {
			ctx := context.Background()
			ints := []int{0, 1, 2, 3, 4}
			strs, err := MapContext(
				ctx,
				parallelism,
				ints,
				func(ctx context.Context, i int) (string, error) {
					return strconv.Itoa(i), nil
				},
			)
			require2.NoError(t, err)
			require2.SlicesEqual(t, []string{"0", "1", "2", "3", "4"}, strs)
		})
	}
}

func TestMapIterator(t *testing.T) {
	strs := MapIterator(
		iterator.Counter(5),
		2, // parallelism
		0, // bufferSize
		func(i int) string {
			return strconv.Itoa(i)
		},
	)
	require2.SlicesEqual(t, []string{"0", "1", "2", "3", "4"}, iterator.Collect(strs))
}

func TestMapStream(t *testing.T) {
	strsStream := MapStream(
		stream.FromIterator(iterator.Counter(5)),
		2, // parallelism
		0, // bufferSize
		func(ctx context.Context, i int) (string, error) {
			return strconv.Itoa(i), nil
		},
	)
	strs, err := stream.Collect(context.Background(), strsStream)
	require2.NoError(t, err)
	require2.SlicesEqual(t, []string{"0", "1", "2", "3", "4"}, strs)
}

func TestMapStreamError(t *testing.T) {
	sender, receiver := stream.Pipe[int](0)
	strsStream := MapStream(
		receiver,
		2, // parallelism
		0, // bufferSize
		func(ctx context.Context, i int) (string, error) {
			return strconv.Itoa(i), nil
		},
	)

	oopsError := errors.New("oops")

	err := sender.Send(context.Background(), 0)
	require2.NoError(t, err)
	err = sender.Send(context.Background(), 1)
	require2.NoError(t, err)
	sender.Close(oopsError)

	for {
		_, err := strsStream.Next(context.Background())
		if err == nil {
			continue
		}
		require2.Equal(t, oopsError, err)
		break
	}
}
