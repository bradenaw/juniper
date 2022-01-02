//go:build go1.18

// package stream allows iterating over sequences of values where iteration may fail.
package stream

import (
	"context"
	"errors"

	"github.com/bradenaw/juniper/iterator"
)

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

// FromIterator returns a Stream that yields the values from iter. This stream ignores the context
// passed to Next during the call to iter.Next.
func FromIterator[T any](iter iterator.Iterator[T]) Stream[T] {
	return &iteratorStream[T]{iter: iter}
}

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

// Chan returns a Stream that receives values from c.
func Chan[T any](c <-chan T) Stream[T] {
	return &chanStream[T]{c: c}
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

// Map transforms the values of s using the conversion f. If f returns an error, terminates the
// stream early.
func Map[T any, U any](s Stream[T], f func(t T) (U, error)) Stream[U] {
	return &mapStream[T, U]{inner: s, f: f}
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

// Chain returns an Iterator that yields all elements of streams[0], then all elements of
// streams[1], and so on.
func Chain[T any](streams ...Stream[T]) Stream[T] {
	return &chainStream[T]{remaining: streams}
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

// Filter returns a Stream that yields only the items from s for which keep returns true. If keep
// returns an error, terminates the stream early.
func Filter[T any](s Stream[T], keep func(T) (bool, error)) Stream[T] {
	return &filterStream[T]{inner: s, keep: keep}
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

// First returns an iterator that yields the first n items from s.
func First[T any](s Stream[T], n int) Stream[T] {
	return &firstStream[T]{inner: s, x: n}
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

// While returns an iterator that terminates before the first item from iter for which f returns
// false. If f returns an error, terminates the stream early.
func While[T any](s Stream[T], f func(T) (bool, error)) Stream[T] {
	return &whileStream[T]{
		inner: s,
		f:     f,
	}
}

// Returned from Sender.Send() when the associated stream has already been closed.
var ErrClosedPipe = errors.New("closed pipe")

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
