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

func MapSeqChan[T any, U any](
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

func MapSeqMutex[T any, U any](
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
		var mu sync.Mutex
		inDone := false
		outDone := false

		var work deque.Deque[valueAndIndex[T]]
		work.Grow(bufferSize)
		workAvailable := sync.NewCond(&mu)

		h := xheap.New(func(a, b valueAndIndex[U]) bool { return a.idx < b.idx }, nil /*initial*/)
		h.Grow(bufferSize)

		inFlight := 0
		// signalled when inFlight goes from bufferSize -> bufferSize-1
		outBufferAvailable := sync.NewCond(&mu)
		// signalled when h.Peek().idx==i
		nextItemAvailable := sync.NewCond(&mu)
		i := 0

		for range parallelism {
			go func() {
				mu.Lock()
				for {
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

		go func() {
			j := 0
			for x := range in {
				mu.Lock()
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
