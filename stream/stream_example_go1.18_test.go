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

func ExampleBatch() {
	ctx := context.Background()

	sender, receiver := stream.Pipe[string](0)
	batchStream := stream.Batch(receiver, 50*time.Millisecond, 3)

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

func ExampleChan() {
	ctx := context.Background()

	c := make(chan string, 3)
	c <- "a"
	c <- "b"
	c <- "c"
	close(c)
	s := stream.Chan(c)

	x, err := stream.Collect(ctx, s)
	fmt.Println(err)
	fmt.Println(x)

	// Output:
	// <nil>
	// [a b c]
}

func ExampleChunk() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{"a", "b", "c", "d", "e", "f", "g", "h"}))

	chunked := stream.Chunk(s, 3)
	item, _ := chunked.Next(ctx)
	fmt.Println(item)
	item, _ = chunked.Next(ctx)
	fmt.Println(item)
	item, _ = chunked.Next(ctx)
	fmt.Println(item)

	// Output:
	// [a b c]
	// [d e f]
	// [g h]
}

func ExampleCompact() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{"a", "a", "b", "c", "c", "c", "a"}))
	compactStream := stream.Compact(s)
	compacted, _ := stream.Collect(ctx, compactStream)
	fmt.Println(compacted)

	// Output:
	// [a b c a]
}

func ExampleCompactFunc() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}))
	compactStream := stream.CompactFunc(s, func(a, b string) bool {
		return a[0] == b[0]
	})
	compacted, _ := stream.Collect(ctx, compactStream)
	fmt.Println(compacted)

	// Output:
	// [bank ghost yaw]
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

func ExampleError() {
	ctx := context.Background()

	s := stream.Error[int](errors.New("foo"))

	_, err := s.Next(ctx)
	fmt.Println(err)

	// Output:
	// foo
}

func ExampleFilter() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]int{1, 2, 3, 4, 5, 6}))

	evensStream := stream.Filter(s, func(ctx context.Context, x int) (bool, error) {
		return x%2 == 0, nil
	})
	evens, _ := stream.Collect(ctx, evensStream)
	fmt.Println(evens)

	// Output:
	// [2 4 6]
}

func ExampleFirst() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]string{"a", "b", "c", "d", "e"}))

	first3Stream := stream.First(s, 3)
	first3, _ := stream.Collect(ctx, first3Stream)
	fmt.Println(first3)

	// Output:
	// [a b c]
}

func ExampleFlatten() {
	ctx := context.Background()

	s := stream.FromIterator(iterator.Slice([]stream.Stream[int]{
		stream.FromIterator(iterator.Slice([]int{0, 1, 2})),
		stream.FromIterator(iterator.Slice([]int{3, 4, 5, 6})),
		stream.FromIterator(iterator.Slice([]int{7})),
	}))

	allStream := stream.Flatten(s)
	all, _ := stream.Collect(ctx, allStream)

	fmt.Println(all)

	// Output:
	// [0 1 2 3 4 5 6 7]
}

func ExampleJoin() {
	ctx := context.Background()

	s := stream.Join(
		stream.FromIterator(iterator.Counter(3)),
		stream.FromIterator(iterator.Counter(5)),
		stream.FromIterator(iterator.Counter(2)),
	)

	all, _ := stream.Collect(ctx, s)

	fmt.Println(all)

	// Output:
	// [0 1 2 0 1 2 3 4 0 1]
}

func ExampleLast() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Counter(10))
	last5, _ := stream.Last(ctx, s, 5)
	fmt.Println(last5)

	s = stream.FromIterator(iterator.Counter(3))
	last5, _ = stream.Last(ctx, s, 5)
	fmt.Println(last5)

	// Output:
	// [5 6 7 8 9]
	// [0 1 2]
}

func ExampleMap() {
	ctx := context.Background()

	s := stream.FromIterator(iterator.Counter(5))
	halfStream := stream.Map(s, func(ctx context.Context, x int) (float64, error) {
		return float64(x) / 2, nil
	})
	all, _ := stream.Collect(ctx, halfStream)
	fmt.Println(all)

	// Output:
	// [0 0.5 1 1.5 2]
}

func ExamplePeekable() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]int{1, 2, 3}))

	p := stream.WithPeek(s)
	x, _ := p.Peek(ctx)
	fmt.Println(x)
	x, _ = p.Next(ctx)
	fmt.Println(x)
	x, _ = p.Next(ctx)
	fmt.Println(x)
	x, _ = p.Peek(ctx)
	fmt.Println(x)

	// Output:
	// 1
	// 1
	// 2
	// 3
}

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

func ExampleReduce() {
	ctx := context.Background()
	s := stream.FromIterator(iterator.Slice([]int{1, 2, 3, 4, 5}))

	sum, _ := stream.Reduce(ctx, s, 0, func(x, y int) (int, error) {
		return x + y, nil
	})
	fmt.Println(sum)

	s = stream.FromIterator(iterator.Slice([]int{1, 3, 2, 3}))
	// Computes the exponentially-weighted moving average of the values of s.
	first := true
	ewma, _ := stream.Reduce(ctx, s, 0, func(running float64, item int) (float64, error) {
		if first {
			first = false
			return float64(item), nil
		}
		return running*0.5 + float64(item)*0.5, nil
	})
	// Should end as 1/8 + 3/8 + 2/4 + 3/2
	//             = 1/8 + 3/8 + 4/8 + 12/8
	//             = 20/8
	//             = 2.5
	fmt.Println(ewma)

	// Output:
	// 15
	// 2.5
}

func ExampleRuns() {
	ctx := context.Background()

	s := stream.FromIterator(iterator.Slice([]int{2, 4, 0, 7, 1, 3, 9, 2, 8}))

	// Contiguous runs of evens/odds.
	parityRuns := stream.Runs(s, func(a, b int) bool {
		return a%2 == b%2
	})

	one, _ := parityRuns.Next(ctx)
	allOne, _ := stream.Collect(ctx, one)
	fmt.Println(allOne)
	two, _ := parityRuns.Next(ctx)
	allTwo, _ := stream.Collect(ctx, two)
	fmt.Println(allTwo)
	three, _ := parityRuns.Next(ctx)
	allThree, _ := stream.Collect(ctx, three)
	fmt.Println(allThree)

	// Output:
	// [2 4 0]
	// [7 1 3 9]
	// [2 8]
}

func ExampleWhile() {
	ctx := context.Background()

	s := stream.FromIterator(iterator.Slice([]string{
		"aardvark",
		"badger",
		"cheetah",
		"dinosaur",
		"egret",
	}))

	beforeD := stream.While(s, func(ctx context.Context, s string) (bool, error) {
		return s < "d", nil
	})

	out, _ := stream.Collect(ctx, beforeD)
	fmt.Println(out)

	// Output:
	// [aardvark badger cheetah]
}
