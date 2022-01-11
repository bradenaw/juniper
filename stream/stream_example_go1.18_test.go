//go:build go1.18

package stream_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/stream"
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
