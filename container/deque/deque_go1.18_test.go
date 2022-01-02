//go:build go1.18

package deque

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/iterator"
)

func FuzzDeque(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		var oracle []byte
		var deque Deque[byte]

		fuzz.Operations(
			b,
			func() {
				require.Equal(t, len(oracle), deque.Len())
				t.Logf("  len = %d", len(oracle))
				t.Logf("  oracle state: %#v", oracle)
				t.Logf("  deque state:  (len(r.a) = %d) %#v", len(deque.a), deque)
			}, // check
			func(x byte) {
				t.Logf("PushFront(%#v)", x)
				deque.PushFront(x)
				oracle = append([]byte{x}, oracle...)
			},
			func(x byte) {
				t.Logf("PushBack(%#v)", x)
				deque.PushBack(x)
				oracle = append(oracle, x)
			},
			func() {
				if len(oracle) == 0 {
					return
				}
				oracleItem := oracle[0]
				t.Logf("PopFront() -> %#v", oracleItem)
				oracle = oracle[1:]
				dequeItem := deque.PopFront()
				require.Equal(t, oracleItem, dequeItem)
			},
			func() {
				if len(oracle) == 0 {
					return
				}
				oracleItem := oracle[len(oracle)-1]
				t.Logf("PopBack() -> %#v", oracleItem)
				oracle = oracle[:len(oracle)-1]
				dequeItem := deque.PopBack()
				require.Equal(t, oracleItem, dequeItem)
			},
			func() {
				if len(oracle) == 0 {
					return
				}
				oracleItem := oracle[0]
				t.Logf("PeekFront() -> %#v", oracleItem)
				dequeItem := deque.PeekFront()
				require.Equal(t, oracleItem, dequeItem)
			},
			func() {
				if len(oracle) == 0 {
					return
				}
				oracleItem := oracle[len(oracle)-1]
				t.Logf("PeekBack() -> %#v", oracleItem)
				dequeItem := deque.PeekBack()
				require.Equal(t, oracleItem, dequeItem)
			},
			func() {
				t.Log("Iterate()")
				oracleAll := oracle
				if len(oracleAll) == 0 {
					oracleAll = nil
				}
				dequeAll := iterator.Collect(deque.Iterate())
				if len(dequeAll) == 0 {
					dequeAll = nil
				}
				require.Equal(t, oracleAll, dequeAll)
			},
			func(n byte) {
				t.Logf("Grow(%d)", n)
				deque.Grow(int(n))
			},
		)
	})
}

func Example() {
	var deque Deque[string]

	deque.PushFront("a")
	deque.PushFront("b")
	fmt.Println(deque.PopFront())
	deque.PushBack("c")
	deque.PushBack("d")
	fmt.Println(deque.PopBack())
	fmt.Println(deque.PopFront())

	// Output:
	// b
	// d
	// a
}
