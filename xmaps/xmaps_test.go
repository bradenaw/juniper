package xmaps_test

import (
	"fmt"

	"github.com/bradenaw/juniper/xmaps"
)

func ExampleUnion() {
	a := xmaps.Set[int]{
		1: {},
		4: {},
	}
	b := xmaps.Set[int]{
		3: {},
		4: {},
	}
	c := xmaps.Set[int]{
		1: {},
		5: {},
	}

	union := xmaps.Union(a, b, c)

	fmt.Println(union)

	// Output:
	// map[1:{} 3:{} 4:{} 5:{}]
}

func ExampleIntersection() {
	a := xmaps.Set[int]{
		1: {},
		2: {},
		4: {},
	}
	b := xmaps.Set[int]{
		1: {},
		3: {},
		4: {},
	}
	c := map[int]struct{}{
		1: {},
		4: {},
		5: {},
	}

	intersection := xmaps.Intersection(a, b, c)

	fmt.Println(intersection)

	// Output:
	// map[1:{} 4:{}]
}

func ExampleIntersects() {
	a := xmaps.Set[int]{
		1: {},
		2: {},
	}
	b := xmaps.Set[int]{
		1: {},
		3: {},
	}
	c := xmaps.Set[int]{
		3: {},
		4: {},
	}

	fmt.Println(xmaps.Intersects(a, b))
	fmt.Println(xmaps.Intersects(b, c))
	fmt.Println(xmaps.Intersects(a, c))

	// Output:
	// true
	// true
	// false
}

func ExampleDifference() {
	a := xmaps.Set[int]{
		1: {},
		4: {},
		5: {},
	}
	b := xmaps.Set[int]{
		3: {},
		4: {},
	}

	difference := xmaps.Difference(a, b)

	fmt.Println(difference)

	// Output:
	// map[1:{} 5:{}]
}
