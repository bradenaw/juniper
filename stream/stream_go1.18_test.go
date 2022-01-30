//go:build go1.18

package stream

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/internal/require2"
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
				if sendClosed {
					// not allowed
					return
				}
				if len(oracle) == bufferSize {
					// might block
					return
				}
				t.Logf("sender.Send(ctx, %d)", x)
				err := sender.Send(context.Background(), x)
				if !recvClosed {
					require2.NoError(t, err)
				} else {
					require2.True(t, err == nil || errors.Is(err, ErrClosedPipe))
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
							require2.Equal(t, End, err)
						} else {
							require2.Equal(t, sendClosedErr, err)
						}
						return
					} else {
						// would block
						return
					}
				}

				t.Log("s.Next(ctx)")

				batch, err := s.Next(context.Background())
				require2.NoError(t, err)

				// Unfortunately we can't actually tell if the receiver has received everything that
				// we sent with Send().
				require2.Greater(t, len(batch), 0)
				require2.LessOrEqual(t, len(batch), batchSize)
				expectedBatch := oracle[:len(batch)]
				require2.SlicesEqual(t, expectedBatch, batch)

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
	require2.NoError(t, err)

	sender, receiver = Pipe[int](1)
	sender.Send(ctx, 1)

	batches = Batch(receiver, 0, 2)
	_, err = batches.Next(context.Background())
	require2.NoError(t, err)
}

func TestPipeConcurrentSend(t *testing.T) {
	ctx := context.Background()
	sender, receiver := Pipe[int](0)

	var wg sync.WaitGroup
	errs := make([]error, 4)
	for i := 0; i < 4; i++ {
		i := i
		wg.Add(1)
		go func() {
			errs[i] = sender.Send(ctx, i)
			wg.Done()
		}()
	}

	time.Sleep(2 * time.Millisecond)

	results := make([]bool, 4)

	item, err := receiver.Next(ctx)
	require2.NoError(t, err)
	results[item] = true

	item, err = receiver.Next(ctx)
	require2.NoError(t, err)
	results[item] = true

	sender.Close(intError(5))
	wg.Wait()

	for i := range results {
		require2.True(t, results[i] || errors.Is(errs[i], intError(5)))
	}
}
