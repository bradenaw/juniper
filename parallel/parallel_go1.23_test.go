//go:build go1.23

package parallel

import (
	"fmt"
	"iter"
	"sync"
	"testing"
	"time"
)

type mapSeqInts func(
	in iter.Seq[int],
	parallelism int,
	bufferSize int,
	f func(int) int,
) iter.Seq[int]

func forEachParallelism(b *testing.B, f func(b *testing.B, parallelism int)) {
	for _, parallelism := range []int{1, 4, 16, 64} {
		b.Run(fmt.Sprintf("parallelism=%d", parallelism), func(b *testing.B) {
			f(b, parallelism)
		})
	}
}

func BenchmarkMapSeq1000VariableMapTime(b *testing.B) {
	run := func(b *testing.B, f mapSeqInts) {
		forEachParallelism(b, func(b *testing.B, parallelism int) {
			b.ReportAllocs()
			for range b.N {
				it := f(count(1000), parallelism, -1, func(x int) int {
					time.Sleep(time.Duration(x%5) * time.Millisecond)

					return x * x
				})
				sum := 0
				for x := range it {
					x += sum
				}
			}
		})
	}

	b.Run("Chan", func(b *testing.B) {
		run(b, MapSeqChan)
	})
	b.Run("Mutex", func(b *testing.B) {
		run(b, MapSeqMutex)
	})
}

// Use b.N as the length of the in/out iterators and have m not take any time, so that we're
// mostly measuring the overhead of MapSeq.
func BenchmarkMapSeqSingleItemTime(b *testing.B) {
	b.ReportAllocs()
	run := func(b *testing.B, f mapSeqInts) {
		b.ReportAllocs()
		it := f(count(b.N), -1, -1, func(x int) int {
			return x * x
		})

		sum := 0
		for x := range it {
			x += sum
		}
	}

	b.Run("Chan", func(b *testing.B) {
		run(b, MapSeqChan)
	})
	b.Run("Mutex", func(b *testing.B) {
		run(b, MapSeqMutex)
	})
}

func TestMapSeq(t *testing.T) {
	check := func(t *testing.T, f mapSeqInts) {
		measurer := &parallelismMeasurer{}
		actual := f(count(1000), 64, -1, func(x int) int {
			measurer.Start()
			// Ensure that even if the map function finishes at different times, we still put the
			// items back in the right order for output.
			time.Sleep(time.Duration(10) * time.Millisecond)
			measurer.Stop()
			return x
		})

		err := checkEqual(count(1000), actual)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("effective parallelism: %f", measurer.Average())
	}

	t.Run("Chan", func(t *testing.T) {
		check(t, MapSeqChan)
	})

	t.Run("Mutex", func(t *testing.T) {
		check(t, MapSeqMutex)
	})
}

func checkEqual[T comparable](a iter.Seq[T], b iter.Seq[T]) error {
	aNext, aStop := iter.Pull(a)
	defer aStop()
	bNext, bStop := iter.Pull(b)
	defer bStop()

	i := 0
	for {
		aItem, aOK := aNext()
		bItem, bOK := bNext()

		if !aOK && !bOK {
			return nil
		}
		if !aOK {
			return fmt.Errorf("left seq ended early after %d items", i+1)
		}
		if !bOK {
			return fmt.Errorf("right seq ended early after %d items", i+1)
		}
		if aItem != bItem {
			return fmt.Errorf("item %d didn't match: %v != %v", i, aItem, bItem)
		}
		i++
	}
}

func count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				break
			}
		}
	}
}

type parallelismMeasurer struct {
	mu    sync.Mutex
	start time.Time
	last  time.Time
	n     int

	running float64
}

func (m *parallelismMeasurer) Start() {
	m.mu.Lock()
	m.advance()
	m.n++
	m.mu.Unlock()
}

func (m *parallelismMeasurer) Stop() {
	m.mu.Lock()
	m.advance()
	m.n--
	m.mu.Unlock()
}

func (m *parallelismMeasurer) advance() {
	now := time.Now()
	if m.start.IsZero() {
		m.start = now
	} else {
		m.running += now.Sub(m.last).Seconds() * float64(m.n)
	}
	m.last = now
}

func (m *parallelismMeasurer) Average() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running / time.Since(m.start).Seconds()
}
