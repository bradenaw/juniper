//go:build go1.18

package stream_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/stream"
	"github.com/bradenaw/juniper/xmath"
)

func ExamplePipe() {
	ctx := context.Background()
	sender, receiver := stream.Pipe[int](0)

	go func() {
		sender.Send(ctx, 1)
		sender.Send(ctx, 2)
		sender.Send(ctx, 3)
		sender.Close(nil)
	}()

	for {
		item, ok := receiver.Next(ctx)
		if !ok {
			break
		}
		fmt.Println(item)
	}
	err := receiver.Close()
	fmt.Println(err)

	// Output:
	// 1
	// 2
	// 3
	// <nil>
}

func ExamplePipe_error() {
	ctx := context.Background()
	sender, receiver := stream.Pipe[int](0)

	oopsError := errors.New("oops")

	go func() {
		sender.Send(ctx, 1)
		sender.Close(oopsError)
	}()

	for {
		item, ok := receiver.Next(ctx)
		if !ok {
			break
		}
		fmt.Println(item)
	}
	err := receiver.Close()
	fmt.Println(err)

	// Output:
	// 1
	// oops
}

func ExampleCollect() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{"a", "b", "c"}))

	x, err := stream.Collect(ctx, s)
	fmt.Println(err)
	fmt.Println(x)

	// Output:
	// <nil>
	// [a b c]
}

func ExampleBatch() {
	ctx := context.Background()

	sender, receiver := stream.Pipe[string](0)

	batchStream := stream.Batch(receiver, 3, 50*time.Millisecond)

	wait := make(chan struct{}, 3)
	go func() {
		_ = sender.Send(ctx, "a")
		_ = sender.Send(ctx, "b")
		// Wait here before sending any more to show that the first batch will flush early because
		// of maxTime=50*time.Millisecond.
		<-wait
		_ = sender.Send(ctx, "c")
		_ = sender.Send(ctx, "d")
		_ = sender.Send(ctx, "e")
		_ = sender.Send(ctx, "f")
		sender.Close(nil)
	}()

	var batches [][]string
	for {
		batch, ok := batchStream.Next(ctx)
		if !ok {
			break
		}
		batches = append(batches, batch)
		wait <- struct{}{}
	}
	err := batchStream.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println(batches)

	// Output:
	// [[a b] [c d e] [f]]
}

type intError int

func (err intError) Error() string {
	return strconv.Itoa(int(err))
}

func FuzzBatch(f *testing.F) {
	f.Fuzz(func(t *testing.T, bufferSize int, batchSize int, b []byte) {
		bufferSize = xmath.Clamp(bufferSize, 0, 1000)
		batchSize = xmath.Clamp(bufferSize, 0, bufferSize)

		t.Logf("bufferSize = %#v", bufferSize)
		t.Logf("batchSize = %#v", batchSize)

		sender, receiver := stream.Pipe[int](bufferSize)
		s := stream.Batch(receiver, batchSize, 10*time.Millisecond)

		var oracle []int
		sendClosed := false
		var sendClosedErr error
		recvClosed := false

		x := 0

		fuzz.Operations(
			b,
			func() { // check
				t.Logf("  oracle        = %#v", oracle)
				t.Logf("  sendClosed    = %#v", sendClosed)
				t.Logf("  sendClosedErr = %#v", sendClosedErr)
				t.Logf("  recvClosed    = %#v", recvClosed)
			},
			func() {
				if len(oracle) == bufferSize {
					// would block
					return
				} else if sendClosed {
					// not allowed
					return
				}
				t.Logf("sender.Send(ctx, %d)", x)
				err := sender.Send(context.Background(), x)
				if recvClosed {
					require.Error(t, stream.ErrClosedPipe)
				} else {
					require.NoError(t, err)
				}
				oracle = append(oracle, x)
				x++
			},
			func(withErr bool) {
				if sendClosed {
					// not allowed
					return
				}
				if withErr {
					t.Logf("sender.Close(intError(%d))", x)
					sendClosedErr = intError(x)
					sender.Close(sendClosedErr)
					x++
				} else {
					t.Log("sender.Close(nil)")
					sender.Close(nil)
				}
				sendClosed = true
			},
			func() {
				if recvClosed {
					// not allowed
					return
				}
				if len(oracle) == 0 {
					if sendClosed {
						t.Log("s.Next(ctx) at end")
						_, ok := s.Next(context.Background())
						require.False(t, ok)
						return
					} else {
						// would block
						return
					}
				}

				t.Log("s.Next(ctx)")

				batch, ok := s.Next(context.Background())
				require.True(t, ok)

				expectedSize := xmath.Min(len(oracle), batchSize)
				expectedBatch := oracle[:expectedSize]
				require.Equal(t, expectedBatch, batch)

				t.Logf(" -> %#v", batch)

				oracle = oracle[expectedSize:]
			},
			func() {
				if recvClosed {
					return
				}
				t.Log("s.Close()")
				err := s.Close()
				require.Equalf(t, sendClosedErr, err, "%s", err)
				recvClosed = true
			},
		)
	})
}
