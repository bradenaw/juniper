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

	defer receiver.Close()
	for {
		item, err := receiver.Next(ctx)
		if err == stream.End {
			break
		} else if err != nil {
			fmt.Printf("stream ended with error: %s\n", err)
			return
		}
		fmt.Println(item)
	}

	// Output:
	// 1
	// 2
	// 3
}

func ExamplePipe_error() {
	ctx := context.Background()
	sender, receiver := stream.Pipe[int](0)

	oopsError := errors.New("oops")

	go func() {
		sender.Send(ctx, 1)
		sender.Close(oopsError)
	}()

	defer receiver.Close()
	for {
		item, err := receiver.Next(ctx)
		if err == stream.End {
			fmt.Println("stream ended normally")
			break
		} else if err != nil {
			fmt.Printf("stream ended with error: %s\n", err)
			return
		}
		fmt.Println(item)
	}

	// Output:
	// 1
	// stream ended with error: oops
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

	defer batchStream.Close()
	var batches [][]string
	for {
		batch, err := batchStream.Next(ctx)
		if err == stream.End {
			break
		} else if err != nil {
			fmt.Printf("stream ended with error: %s\n", err)
			return
		}
		batches = append(batches, batch)
		wait <- struct{}{}
	}
	fmt.Println(batches)

	// Output:
	// [[a b] [c d e] [f]]
}
