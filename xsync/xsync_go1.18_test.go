package xsync

import "fmt"

func ExampleLazy() {
	var (
		expensive = Lazy(func() string {
			fmt.Println("doing expensive init")
			return "foo"
		})
	)

	fmt.Println(expensive())
	fmt.Println(expensive())

	// Output:
	// doing expensive init
	// foo
	// foo
}
