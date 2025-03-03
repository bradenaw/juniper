//go:build go1.23

package require2

import (
	"iter"
	"testing"
)

func SeqsEqual[T comparable](t *testing.T, a iter.Seq[T], b iter.Seq[T]) {
	aNext, aStop := iter.Pull(a)
	defer aStop()
	bNext, bStop := iter.Pull(b)
	defer bStop()

	i := 0
	for {
		aItem, aOK := aNext()
		bItem, bOK := bNext()

		if !aOK && !bOK {
			break
		}

		if !aOK {
			t.Fatalf("left sequence had fewer items, ended after %d", i)
		}
		if !bOK {
			t.Fatalf("right sequence had fewer items, ended after %d", i)
		}

		if aItem != bItem {
			t.Fatalf("items at index %d differed: %#v != %#v", i, a, b)
		}

		i++
	}
}
