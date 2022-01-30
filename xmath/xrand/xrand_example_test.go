package xrand_test

import (
	"fmt"
	"math/rand"

	"github.com/bradenaw/juniper/xmath/xrand"
)

func ExampleSample() {
	r := rand.New(rand.NewSource(0))

	sample := xrand.RSample(r, 100, 5)

	fmt.Println(sample)

	// Output:
	// [45 71 88 93 60]
}
