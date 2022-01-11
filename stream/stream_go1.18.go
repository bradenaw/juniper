//go:build go1.18

// package stream allows iterating over sequences of values where iteration may fail.
package stream

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/bradenaw/juniper/iterator"
)

// Returned from Sender.Send() when the associated stream has already been closed.
var ErrClosedPipe = errors.New("closed pipe")

// Stream is used to iterate over a sequence of values. It is similar to Iterator, except intended
// for use when iteration may fail for some reason, usually because the sequence requires I/O to
// produce.
//
// Streams and the combinator functions are lazy, meaning they do no work until a call to Next().
type Stream[T any] interface {
	// Next advances the stream and returns the next item. Once the stream is finished, the first
	// return is meaningless and the second return is false. The final value of the stream will have
	// true in the second return.
	//
	// If an error is encountered during Next(), it should return false immediately and surface the
	// error from Close().
	Next(ctx context.Context) (T, bool)
	// Close ends receiving from the stream and returns any error encountered either in the course
	// of Next or while cleaning up afterward.
	Close() error
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Converters + Constructors                                                                      //
// Functions that produce a Stream.                                                               //
////////////////////////////////////////////////////////////////////////////////////////////////////

// Chan returns a Stream that receives values from c.
func Chan[T any](c <-chan T) Stream[T] {
	return &chanStream[T]{c: c}
}

type chanStream[T any] struct {
	c   <-chan T
	err error
}

func (s *chanStream[T]) Next(ctx context.Context) (T, bool) {
	var zero T
	select {
	case item, ok := <-s.c:
		return item, ok
	case <-ctx.Done():
		s.err = ctx.Err()
		return zero, false
	}
}

func (s *chanStream[T]) Close() error {
	return s.err
}

// FromIterator returns a Stream that yields the values from iter. This stream ignores the context
// passed to Next during the call to iter.Next.
func FromIterator[T any](iter iterator.Iterator[T]) Stream[T] {
	return &iteratorStream[T]{iter: iter}
}

type iteratorStream[T any] struct {
	iter iterator.Iterator[T]
	err  error
}

func (s *iteratorStream[T]) Next(ctx context.Context) (T, bool) {
	if ctx.Err() != nil {
		s.err = ctx.Err()
		var zero T
		return zero, false
	}
	return s.iter.Next()
}

func (s *iteratorStream[T]) Close() error {
	return s.err
}

// Pipe returns a linked sender and receiver pair. Values sent using sender.Send will be delivered
// to the given Stream. The Stream will terminate when the sender is closed.
//
// bufferSize is the number of elements in the buffer between the sender and the receiver. 0 has the
// same meaning as for the built-in make(chan).
func Pipe[T any](bufferSize int) (*Sender[T], Stream[T]) {
	c := make(chan T, bufferSize)
	senderDone := make(chan error, 1)
	streamDone := make(chan struct{})

	sender := &Sender[T]{
		c:          c,
		senderDone: senderDone,
		streamDone: streamDone,
	}
	receiver := &pipeStream[T]{
		c:          c,
		senderDone: senderDone,
		streamDone: streamDone,
	}

	return sender, receiver
}

// Sender is the send half of a pipe returned by Pipe.
type Sender[T any] struct {
	c          chan<- T
	senderDone chan<- error
	streamDone <-chan struct{}
}

// Send attemps to send x to the receiver. If the receiver closes before x can be sent, returns
// ErrClosedPipe immediately. If ctx expires before x can be sent, returns ctx.Err().
//
// A nil return does not necessarily mean that the receiver will see x, since the receiver may close
// early.
func (s *Sender[T]) Send(ctx context.Context, x T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.streamDone:
		return ErrClosedPipe
	case s.c <- x:
		return nil
	}
}

// Close closes the Sender, signalling to the receiver that no more values will be sent. If an error
// is provided, it will surface when closing the receiver.
func (s *Sender[T]) Close(err error) {
	s.senderDone <- err
	close(s.c)
}

type pipeStream[T any] struct {
	c          <-chan T
	senderDone <-chan error
	streamDone chan<- struct{}
	err        error
}

// sentinel to early-exit from pipeStream.Next.
var errPipeDone = errors.New("pipe done")

func (s *pipeStream[T]) Next(ctx context.Context) (T, bool) {
	var zero T
	if s.err != nil {
		return zero, false
	}
	select {
	case <-ctx.Done():
		s.err = ctx.Err()
		return zero, false
	case item, ok := <-s.c:
		if !ok {
			s.err = errPipeDone
			return zero, false
		}
		return item, true
	}
}

func (s *pipeStream[T]) Close() error {
	close(s.streamDone)
	select {
	case err := <-s.senderDone:
		return err
	default:
		if s.err == errPipeDone {
			return nil
		}
		return s.err
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Reducers                                                                                       //
// Functions that consume a stream and produce some kind of final value.                          //
////////////////////////////////////////////////////////////////////////////////////////////////////

// Collect advances s to the end and returns all of the items seen as a slice.
func Collect[T any](ctx context.Context, s Stream[T]) ([]T, error) {
	var out []T
	for {
		item, ok := s.Next(ctx)
		if !ok {
			break
		}
		out = append(out, item)
	}
	err := s.Close()
	if err != nil {
		return nil, err
	}
	return out, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Combinators                                                                                    //
// Functions that take and return iterators, transforming the output somehow.                     //
////////////////////////////////////////////////////////////////////////////////////////////////////

// Batch returns a stream of non-overlapping batches from s of size batchSize. Batch is similar to
// Chunk with the added feature that an underfilled batch will be delivered to the output stream if
// any item has been in the batch for more than maxWait.
func Batch[T any](s Stream[T], batchSize int, maxWait time.Duration) Stream[[]T] {
	bgCtx, bgCancel := context.WithCancel(context.Background())

	out := &batchStream[T]{
		batchC:   make(chan []T),
		waiting:  make(chan struct{}),
		bgCancel: bgCancel,
	}

	c := make(chan T)

	out.eg.Go(func() error {
		for {
			item, ok := s.Next(bgCtx)
			if !ok {
				break
			}
			c <- item
		}
		err := s.Close()
		close(c)
		if err == context.Canceled && bgCtx.Err() == context.Canceled {
			// Implies the caller Close()d without finishing. That's not meant to be an error.
			return nil
		}
		return err
	})

	// Build up batches and flush them when either:
	// A) The batch is full.
	// B) It's been at least maxWait since the first item arrived _and_ there is somebody waiting.
	// No sense in underfilling a batch if nobody's actually asking for it yet.
	// C) There aren't any more items.
	out.eg.Go(func() error {
		batch := make([]T, 0, batchSize)
		var batchStart time.Time
		var timer *time.Timer
		// Starts off as nil so that the timerC select arm isn't chosen until populated.  Also set
		// to nil when we've already stopped or received from timer to know when it needs to be
		// drained.
		var timerC <-chan time.Time
		waitingAtEmpty := false

		defer func() {
			if timer != nil {
				timer.Stop()
			}
			close(out.batchC)
		}()

		flush := func() error {
			select {
			case <-bgCtx.Done():
				return nil
			case out.batchC <- batch:
			}
			batch = make([]T, 0, batchSize)
			waitingAtEmpty = false
			return nil
		}

		stopTimer := func() {
			if timer == nil {
				return
			}
			stopped := timer.Stop()
			if !stopped && timerC != nil {
				<-timerC
			}
			timerC = nil
		}

		startTimer := func() {
			stopTimer()
			if timer == nil {
				timer = time.NewTimer(maxWait - time.Since(batchStart))
			} else {
				timer.Reset(maxWait - time.Since(batchStart))
			}
			timerC = timer.C
		}

		for {
			select {
			case item, ok := <-c:
				if !ok { // Case (C): we're done.
					// Flush what we have so far, if any.
					if len(batch) > 0 {
						return flush()
					}
					return nil
				}
				batch = append(batch, item)
				if len(batch) == batchSize { // Case (A): the batch is full.
					stopTimer()
					err := flush()
					if err != nil {
						return err
					}
				} else if len(batch) == 1 { // Bookkeeping for case (B).
					batchStart = time.Now()
					if waitingAtEmpty {
						startTimer()
					}
				}
			case <-timerC: // Case (B).
				timerC = nil
				// Being here already implies the conditions are true, since the timer is only
				// running while the batch is non-empty and there's somebody waiting.
				err := flush()
				if err != nil {
					return err
				}
			case <-out.waiting: // Bookkeeping for case (B).
				if len(batch) > 0 {
					// Time already elapsed, just deliver the batch now.
					if time.Since(batchStart) > maxWait {
						err := flush()
						if err != nil {
							return err
						}
					} else {
						startTimer()
					}
				} else {
					// Timer will start when the first item shows up.
					waitingAtEmpty = true
				}
			}
		}
	})

	return out
}

type batchStream[T any] struct {
	bgCancel context.CancelFunc
	batchC   chan []T
	waiting  chan struct{}
	err      error
	eg       errgroup.Group
}

func (iter *batchStream[T]) Next(ctx context.Context) ([]T, bool) {
	select {
	// There might be a batch already ready because it filled before we even asked.
	case batch, ok := <-iter.batchC:
		return batch, ok
	// Otherwise, we need to let the sender know we're waiting so that they can flush an underfilled
	// batch at interval.
	case iter.waiting <- struct{}{}:
		select {
		case batch, ok := <-iter.batchC:
			return batch, ok
		case <-ctx.Done():
			iter.err = ctx.Err()
			return nil, false
		}
	case <-ctx.Done():
		iter.err = ctx.Err()
		return nil, false
	}
}

func (iter *batchStream[T]) Close() error {
	iter.bgCancel()
	err := iter.eg.Wait()
	if iter.err != nil {
		return iter.err
	}
	return err
}

// Chain returns a Stream that yields all elements from streams[0], then all elements from
// streams[1], and so on.
func Chain[T any](streams ...Stream[T]) Stream[T] {
	return &chainStream[T]{remaining: streams}
}

type chainStream[T any] struct {
	remaining []Stream[T]
	err       error
}

func (s *chainStream[T]) Next(ctx context.Context) (T, bool) {
	var zero T
	if s.err != nil {
		return zero, false
	}
	for len(s.remaining) > 0 {
		item, ok := s.remaining[0].Next(ctx)
		if !ok {
			err := s.remaining[0].Close()
			s.remaining = s.remaining[1:]
			if err != nil {
				s.err = err
				return zero, false
			}
		}
		return item, true
	}
	return zero, false
}

func (s *chainStream[T]) Close() error {
	for i := range s.remaining {
		err := s.remaining[i].Close()
		if err != nil && s.err == nil {
			s.err = err
		}
	}
	return s.err
}

// Chunk returns a stream of non-overlapping chunks from s of size chunkSize. The last chunk will be
// smaller than chunkSize if the stream does not contain an even multiple.
func Chunk[T any](s Stream[T], chunkSize int) Stream[[]T] {
	return &chunkStream[T]{
		inner:     s,
		chunkSize: chunkSize,
	}
}

type chunkStream[T any] struct {
	inner     Stream[T]
	chunkSize int
}

func (s *chunkStream[T]) Next(ctx context.Context) ([]T, bool) {
	chunk := make([]T, 0, s.chunkSize)
	for {
		item, ok := s.inner.Next(ctx)
		if !ok {
			break
		}
		chunk = append(chunk, item)
		if len(chunk) == s.chunkSize {
			return chunk, true
		}
	}
	if len(chunk) > 0 {
		return chunk, true
	}
	return nil, false
}

func (s *chunkStream[T]) Close() error {
	return s.inner.Close()
}

// Filter returns a Stream that yields only the items from s for which keep returns true. If keep
// returns an error, terminates the stream early.
func Filter[T any](s Stream[T], keep func(T) (bool, error)) Stream[T] {
	return &filterStream[T]{inner: s, keep: keep}
}

type filterStream[T any] struct {
	inner Stream[T]
	keep  func(T) (bool, error)
	err   error
}

func (s *filterStream[T]) Next(ctx context.Context) (T, bool) {
	var zero T
	if s.err != nil {
		return zero, false
	}
	for {
		item, ok := s.inner.Next(ctx)
		if !ok {
			return zero, false
		}
		ok, err := s.keep(item)
		if err != nil {
			s.err = err
			return zero, false
		}
		if ok {
			return item, true
		}
	}
}

func (s *filterStream[T]) Close() error {
	closeErr := s.inner.Close()
	if s.err != nil {
		return s.err
	}
	return closeErr
}

// First returns a Stream that yields the first n items from s.
func First[T any](s Stream[T], n int) Stream[T] {
	return &firstStream[T]{inner: s, x: n}
}

type firstStream[T any] struct {
	inner Stream[T]
	x     int
}

func (s *firstStream[T]) Next(ctx context.Context) (T, bool) {
	if s.x <= 0 {
		var zero T
		return zero, false
	}
	s.x--
	return s.inner.Next(ctx)
}

func (s *firstStream[T]) Close() error {
	return s.inner.Close()
}

// Map transforms the values of s using the conversion f. If f returns an error, terminates the
// stream early.
func Map[T any, U any](s Stream[T], f func(t T) (U, error)) Stream[U] {
	return &mapStream[T, U]{inner: s, f: f}
}

type mapStream[T any, U any] struct {
	inner Stream[T]
	f     func(t T) (U, error)
	err   error
}

func (s *mapStream[T, U]) Next(ctx context.Context) (U, bool) {
	var zero U
	if s.err != nil {
		return zero, false
	}
	item, ok := s.inner.Next(ctx)
	if !ok {
		return zero, false
	}
	mapped, err := s.f(item)
	if err != nil {
		s.err = err
		return zero, false
	}
	return mapped, true
}

func (s *mapStream[T, U]) Close() error {
	closeErr := s.inner.Close()
	if s.err != nil {
		return s.err
	}
	return closeErr
}

// While returns a Stream that terminates before the first item from s for which f returns false.
// If f returns an error, terminates the stream early.
func While[T any](s Stream[T], f func(T) (bool, error)) Stream[T] {
	return &whileStream[T]{
		inner: s,
		f:     f,
	}
}

type whileStream[T any] struct {
	inner Stream[T]
	f     func(T) (bool, error)
	done  bool
	err   error
}

func (s *whileStream[T]) Next(ctx context.Context) (T, bool) {
	var zero T
	if s.done {
		return zero, false
	}
	item, ok := s.inner.Next(ctx)
	if !ok {
		return zero, false
	}
	ok, err := s.f(item)
	if !ok || err != nil {
		s.done = true
		s.err = err
		return zero, false
	}
	return item, true
}

func (s *whileStream[T]) Close() error {
	closeErr := s.inner.Close()
	if s.err != nil {
		return s.err
	}
	return closeErr
}
