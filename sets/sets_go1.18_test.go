//go:build go1.18

package sets_test

import (
	"fmt"

	"github.com/bradenaw/juniper/maps"
	"github.com/bradenaw/juniper/sets"
)

func ExampleUnion() {
	a := maps.Set[int]{
		1: {},
		4: {},
	}
	b := maps.Set[int]{
		3: {},
		4: {},
	}
	c := maps.Set[int]{
		1: {},
		5: {},
	}

	out := make(maps.Set[int])

	union := sets.Union[int](out, a, b, c)

	fmt.Println(union)

	// Output:
	// map[1:{} 3:{} 4:{} 5:{}]
}

func ExampleIntersection() {
	a := maps.Set[int]{
		1: {},
		2: {},
		4: {},
	}
	b := maps.Set[int]{
		1: {},
		3: {},
		4: {},
	}
	c := maps.Set[int]{
		1: {},
		4: {},
		5: {},
	}

	out := make(maps.Set[int])

	intersection := sets.Intersection[int](out, a, b, c)

	fmt.Println(intersection)

	// Output:
	// map[1:{} 4:{}]
}

func ExampleDifference() {
	a := maps.Set[int]{
		1: {},
		4: {},
		5: {},
	}
	b := maps.Set[int]{
		3: {},
		4: {},
	}

	out := make(maps.Set[int])

	difference := sets.Difference[int](out, a, b)

	fmt.Println(difference)

	// Output:
	// map[1:{} 5:{}]
}
