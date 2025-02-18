//go:build go1.23

package parallel

import (
	"iter"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/bradenaw/juniper/container/deque"
	"github.com/bradenaw/juniper/container/xheap"
)

// MapSeq uses parallelism goroutines to call f once for each element yielded by in. The
// returned iterator returns these results in the same order that in yielded them in.
//
// If parallelism <= 0, uses GOMAXPROCS instead.
//
// bufferSize is the size of the work buffer. A larger buffer uses more memory but gives better
// throughput in the face of larger variance in the processing time for f.
func MapSeq[T any, U any](
	in iter.Seq[T],
	parallelism int,
	bufferSize int,
	f func(T) U,
) iter.Seq[U] {
	return mapSeqMutex(in, parallelism, bufferSize, f)
}

func mapSeqChan[T any, U any](
	in iter.Seq[T],
	parallelism int,
	bufferSize int,
	f func(T) U,
) iter.Seq[U] {
	if parallelism <= 0 {
		parallelism = runtime.GOMAXPROCS(-1)
	}
	if bufferSize < parallelism {
		bufferSize = parallelism
	}

	return func(yield func(U) bool) {
		var wg sync.WaitGroup
		wg.Add(parallelism + 1)
		defer wg.Wait()
		done := make(chan struct{})
		defer close(done)

		work := make(chan valueAndIndex[T])
		out := make(chan valueAndIndex[U])
		sem := make(chan struct{}, bufferSize)
		var nDone atomic.Int64
		for range parallelism {
			go func() {
				defer wg.Done()
				defer func() {
					if nDone.Add(1) == int64(parallelism) {
						close(out)
					}
				}()
				for {
					var item valueAndIndex[T]
					var ok bool
					select {
					case <-done:
						return
					case item, ok = <-work:
					}

					if !ok {
						return
					}

					select {
					case <-done:
						return
					case out <- valueAndIndex[U]{
						idx:   item.idx,
						value: f(item.value),
					}:
					}
				}
			}()
		}

		go func() {
			defer wg.Done()
			i := 0
			for x := range in {
				select {
				case <-done:
					return
				case sem <- struct{}{}:
				}

				select {
				case <-done:
					return
				case work <- valueAndIndex[T]{
					value: x,
					idx:   i,
				}:
				}
				i++
			}
			close(work)
		}()

		i := 0
		h := xheap.New(func(a, b valueAndIndex[U]) bool { return a.idx < b.idx }, nil /*initial*/)
		h.Grow(bufferSize)
		for {
			item, ok := <-out
			if !ok {
				break
			}
			h.Push(item)
			for h.Len() > 0 && h.Peek().idx == i {
				if !yield(h.Pop().value) {
					return
				}
				<-sem
				i++
			}
		}
	}
}

// MapSeqMutex uses parallelism goroutines to call f once for each element yielded by in. The
// returned iterator returns these results in the same order that in yielded them in.
//
// If parallelism <= 0, uses GOMAXPROCS instead.
//
// bufferSize is the size of the work buffer. A larger buffer uses more memory but gives better
// throughput in the face of larger variance in the processing time for f.
func mapSeqMutex[T any, U any](
	in iter.Seq[T],
	parallelism int,
	bufferSize int,
	f func(T) U,
) iter.Seq[U] {
	if parallelism <= 0 {
		parallelism = runtime.GOMAXPROCS(-1)
	}
	if bufferSize < parallelism {
		bufferSize = parallelism
	}

	return func(yield func(U) bool) {
		// TODO: The purpose of wg is to ensure that we're not leaking anything, including by the
		// caller having an f block for a very long time. Without this, we'd leave behind the worker
		// goroutine blocked in f and carry on like everyting's fine. Really, the caller needs to
		// arrange for f to return promptly if they want to wander off.
		//
		// This all seems desirable but worth explaining in the comment.
		var wg sync.WaitGroup
		wg.Add(parallelism + 1)
		defer wg.Wait()

		var mu sync.Mutex
		// True when in has stopped producing items, so all we need to do is make sure everything
		// we've alredy received from it gets processed.
		inDone := false
		// True when yield has returned false, meaning the caller has broken out of their loop. They
		// won't ask for any more items so all of the background work can exit immediately.
		outDone := false

		// Queue of work to be done by the mapper goroutines.
		var work deque.Deque[valueAndIndex[T]]
		// There can only be this many in flight at once, so we'll never need to grow larger than
		// this.
		work.Grow(bufferSize)
		// Signalled when items are added to `work`.
		workAvailable := sync.NewCond(&mu)

		// The heap of results, ordered by the index they arrived from `in` at.
		h := xheap.New(func(a, b valueAndIndex[U]) bool { return a.idx < b.idx }, nil /*initial*/)
		// There can only be this many in flight at once, so we'll never need to grow larger than
		// this.
		h.Grow(bufferSize)

		// The number of values that have been read from `in` but have not been produced to `yield`.
		// Tracked because we bound this to be `bufferSize` so that we do not use unbounded memory
		// if one invocation of `f` takes a very long time.
		inFlight := 0
		// Signalled when inFlight goes from bufferSize -> bufferSize-1, or when outDone becomes
		// true.
		outBufferAvailable := sync.NewCond(&mu)
		// Signalled when h.Peek().idx==i, meaning we're ready to yield (at least) the top of h.
		nextItemAvailable := sync.NewCond(&mu)
		// The next index that needs to be returned.
		i := 0

		// The workers that'll be calling f.
		for range parallelism {
			go func() {
				defer wg.Done()
				mu.Lock()
				for {
					// mu is always held here.

					if outDone {
						mu.Unlock()
						return
					}
					for work.Len() == 0 {
						if outDone || inDone {
							// Have to signal here because the heap might already be empty by the
							// time we realize that there's no more in.
							nextItemAvailable.Signal()
							mu.Unlock()
							return
						}
						workAvailable.Wait()
					}
					item := work.PopFront()
					mu.Unlock()

					out := valueAndIndex[U]{
						idx:   item.idx,
						value: f(item.value),
					}

					mu.Lock()
					h.Push(out)
					if out.idx == i {
						nextItemAvailable.Signal()
					}
				}
			}()
		}

		// One more goroutine to read through `in` and produce it to the workers above.
		go func() {
			defer wg.Done()
			j := 0
			for x := range in {
				mu.Lock()
				// Block if there are already too many in flight to avoid `h` becoming unreasonably
				// large if one `f` invocation takes a very long time.
				for inFlight >= bufferSize {
					if outDone {
						mu.Unlock()
						return
					}
					outBufferAvailable.Wait()
				}
				inFlight++
				work.PushBack(valueAndIndex[T]{
					value: x,
					idx:   j,
				})
				workAvailable.Signal()
				mu.Unlock()
				j++
			}

			mu.Lock()
			inDone = true
			// Wake all of the workers so that they can exit.
			workAvailable.Broadcast()
			mu.Unlock()
		}()

		for {
			mu.Lock()
			for h.Len() == 0 || h.Peek().idx != i {
				if inDone && inFlight == 0 {
					mu.Unlock()
					return
				}
				nextItemAvailable.Wait()
			}
			x := h.Pop().value
			inFlight--
			if inFlight == bufferSize-1 {
				outBufferAvailable.Signal()
			}
			mu.Unlock()
			if !yield(x) {
				mu.Lock()
				outDone = true
				workAvailable.Broadcast()
				outBufferAvailable.Broadcast()
				mu.Unlock()
				break
			}
			i++
		}
	}
}

// MapSeq2 uses parallelism goroutines to call f once for each pair yielded by in. The returned
// iterator returns these results in the same order that in yielded them in.
//
// If parallelism <= 0, uses GOMAXPROCS instead.
//
// bufferSize is the size of the work buffer. A larger buffer uses more memory but gives better
// throughput in the face of larger variance in the processing time for f.
func MapSeq2[T0 any, T1 any, U0 any, U1 any](
	in iter.Seq2[T0, T1],
	parallelism int,
	bufferSize int,
	f func(T0, T1) (U0, U1),
) iter.Seq2[U0, U1] {
	type pair[A any, B any] struct {
		a A
		b B
	}

	return func(yield func(U0, U1) bool) {
		it := MapSeq(
			func(yield func(pair[T0, T1]) bool) {
				for t0, t1 := range in {
					if !yield(pair[T0, T1]{t0, t1}) {
						return
					}
				}
			},
			parallelism,
			bufferSize,
			func(p pair[T0, T1]) pair[U0, U1] {
				u0, u1 := f(p.a, p.b)
				return pair[U0, U1]{u0, u1}
			},
		)
		for p := range it {
			if !yield(p.a, p.b) {
				return
			}
		}
	}
}
