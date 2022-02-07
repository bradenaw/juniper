package chans

import (
	"math/rand"
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

		go func() {
			Merge(out, slices.Map(ins, func(c chan byte) <-chan byte { return c })...)
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
	r := rand.New(rand.NewSource(time.Now().Unix()))

	count := 0
	start := time.Now()
	lastReport := time.Now()

	for {
		n := r.Intn(4) + 1

		if count%10000 == 0 || time.Since(lastReport) > 3*time.Second {
			t.Logf("%s %d", time.Since(start).Round(time.Second), count)
			lastReport = time.Now()
		}
		count++

		out := make(chan byte)
		ins := make([]chan byte, n)
		for i := range ins {
			ins[i] = make(chan byte)
		}
		inDone := make([]bool, n)

		go func() {
			Merge(out, slices.Map(ins, func(c chan byte) <-chan byte { return c })...)
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
			notDone := make([]int, 0, n)
			for i := range ins {
				if !inDone[i] {
					notDone = append(notDone, i)
				}
			}
			if len(notDone) == 0 {
				break
			}
			idx := notDone[r.Intn(len(notDone))]
			switch r.Intn(2) {
			case 0:
				v := byte(r.Intn(256))
				inS = append(inS, v)
				ins[idx] <- v
			case 1:
				close(ins[idx])
				inDone[idx] = true
			}
		}
		<-done

		require2.SlicesEqual(t, inS, outS)
	}
}
