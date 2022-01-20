//go:build go1.18

package stream

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/xmath"
)

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

		sender, receiver := Pipe[int](bufferSize)
		s := Batch(receiver, 10*time.Millisecond, batchSize)

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
					require.Error(t, ErrClosedPipe)
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
						_, err := s.Next(context.Background())
						if sendClosedErr == nil {
							require.Equal(t, End, err)
						} else {
							require.Equal(t, sendClosedErr, err)
						}
						return
					} else {
						// would block
						return
					}
				}

				t.Log("s.Next(ctx)")

				batch, err := s.Next(context.Background())
				require.NoError(t, err)

				// Unfortunately we can't actually tell if the receiver has received everything that
				// we sent with Send().
				require.Greater(t, len(batch), 0)
				require.LessOrEqual(t, len(batch), batchSize)
				expectedBatch := oracle[:len(batch)]
				require.Equal(t, expectedBatch, batch)

				t.Logf(" -> %#v", batch)

				oracle = oracle[len(expectedBatch):]
			},
			func() {
				if recvClosed {
					return
				}
				t.Log("s.Close()")
				s.Close()
				recvClosed = true
			},
		)
	})
}

func TestBatch(t *testing.T) {
	ctx := context.Background()
	sender, receiver := Pipe[int](1)
	sender.Send(ctx, 1)

	batches := Batch(receiver, 365*24*time.Hour, 1)
	_, err := batches.Next(ctx)
	require.NoError(t, err)

	sender, receiver = Pipe[int](1)
	sender.Send(ctx, 1)

	batches = Batch(receiver, 0, 2)
	_, err = batches.Next(context.Background())
	require.NoError(t, err)
}
