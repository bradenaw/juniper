//go:build go1.18

package stream_test

import (
	"context"
	"errors"
	"fmt"

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

func ExamplePipe_Error() {
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
