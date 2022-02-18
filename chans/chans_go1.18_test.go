//go:build go1.18

package chans

import (
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bradenaw/juniper/internal/require2"
	"github.com/bradenaw/juniper/slices"
)

func FuzzMerge(f *testing.F) {
	f.Fuzz(func(t *testing.T, n int, b []byte) {
		if n > 5 || n <= 0 {
			return
		}

		t.Logf("n = %d", n)

		out := make(chan byte)
		ins := make([]chan byte, n)
		for i := range ins {
			ins[i] = make(chan byte)
		}
		ins2 := slices.Map(ins, func(c chan byte) <-chan byte { return c })

		go func() {
			Merge(out, ins2...)
			close(out)
		}()

		var inSlice []byte
		var outSlice []byte
		done := make(chan struct{})
		go func() {
			for item := range out {
				outSlice = append(outSlice, item)
			}
			close(done)
		}()

	Loop:
		for {
			if len(b) < 3 {
				break
			}
			idx := int(b[0])
			if idx >= len(ins) {
				break
			}
			switch b[1] {
			case 0:
				inSlice = append(inSlice, b[2])
				ins[idx] <- b[2]
			case 1:
				close(ins[idx])
				ins = slices.RemoveUnordered(ins, idx, 1)
			default:
				break Loop
			}
			b = b[3:]
		}
		for _, in := range ins {
			close(in)
		}
		<-done

		require2.SlicesEqual(t, inSlice, outSlice)
	})
}

func TestStressMerge(t *testing.T) {
	t.Skip()
	count := uint64(0)
	start := time.Now()

	go func() {
		for {
			t.Logf("%s %d", time.Since(start).Round(time.Second), count)
			time.Sleep(3 * time.Second)
		}
	}()

	for i := 0; i < runtime.GOMAXPROCS(-1); i++ {
		go func() {
			r := rand.New(rand.NewSource(time.Now().Unix()))

			for {
				n := r.Intn(4) + 1

				atomic.AddUint64(&count, 1)

				out := make(chan byte)
				ins := make([]chan byte, n)
				for i := range ins {
					ins[i] = make(chan byte)
				}

				ins2 := slices.Map(ins, func(c chan byte) <-chan byte { return c })
				go func() {
					Merge(out, ins2...)
					close(out)
				}()

				var inS []byte
				var outS []byte
				done := make(chan struct{})
				go func() {
					for item := range out {
						outS = append(outS, item)
					}
					close(done)
				}()

				for {
					if len(ins) == 0 {
						break
					}
					idx := r.Intn(len(ins))
					switch r.Intn(2) {
					case 0:
						v := byte(r.Intn(256))
						inS = append(inS, v)
						ins[idx] <- v
					case 1:
						close(ins[idx])
						nBefore := len(ins)
						ins = slices.RemoveUnordered(ins, idx, 1)
						require2.Equal(t, len(ins), nBefore-1)
					}
				}
				<-done

				require2.SlicesEqual(t, inS, outS)
			}
		}()
	}

	c := make(chan struct{})
	<-c
}
