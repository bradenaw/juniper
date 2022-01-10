package xrand_test

import (
	"fmt"
	"math/rand"

	"github.com/bradenaw/juniper/xmath/xrand"
)

func ExampleSample() {
	r := rand.New(rand.NewSource(0))

	sample := xrand.Sample(r, 100, 5)

	fmt.Println(sample)

}

// Output:
// [32 78 43 58 72]
