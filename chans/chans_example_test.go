package chans_test

import (
	"fmt"
	"sync"

	"github.com/bradenaw/juniper/chans"
)

func ExampleMerge() {
	a := make(chan int)
	go func() {
		a <- 0
		a <- 1
		a <- 2
		close(a)
	}()
	b := make(chan int)
	go func() {
		b <- 5
		b <- 6
		b <- 7
		b <- 8
		close(b)
	}()

	out := make(chan int)
	done := make(chan struct{})
	go func() {
		for i := range out {
			fmt.Println(i)
		}
		close(done)
	}()

	chans.Merge(out, a, b)
	close(out)
	<-done

	// Unordered output:
	// 0
	// 1
	// 2
	// 5
	// 6
	// 7
	// 8
}

func ExampleReplicate() {
	in := make(chan int)
	go func() {
		in <- 0
		in <- 1
		in <- 2
		in <- 3
		close(in)
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	a := make(chan int)
	go func() {
		for i := range a {
			fmt.Println(i * 2)
		}
		wg.Done()
	}()

	b := make(chan int)
	go func() {
		x := 0
		for i := range b {
			x += i
			fmt.Println(x)
		}
		wg.Done()
	}()

	chans.Replicate(in, a, b)
	close(a)
	close(b)
	wg.Wait()

	// Unordered output:
	// 0
	// 2
	// 4
	// 6
	// 0
	// 1
	// 3
	// 6
}
