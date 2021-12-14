//go:build go1.18

package deque

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/iterator"
)

func FuzzDeque(f *testing.F) {
	const (
		PushFront = 0b00000000
		PushBack  = 0b00100000
		PopFront  = 0b01000000
		PopBack   = 0b01100000
		PeekFront = 0b10000000
		PeekBack  = 0b10100000
		Iterate   = 0b11000000
	)

	f.Add([]byte{
		PushFront | 0x01,
		PopBack,
	})

	f.Fuzz(func(t *testing.T, b []byte) {
		var oracle []byte
		var deque Deque[byte]

		for i := range b {
			operand := b[i] & 0b00011111
			switch b[i] & 0b11100000 {
			case PushFront:
				t.Logf("PushFront(%#v)", operand)
				deque.PushFront(operand)
				oracle = append([]byte{operand}, oracle...)
			case PushBack:
				t.Logf("PushBack(%#v)", operand)
				deque.PushBack(operand)
				oracle = append(oracle, operand)
			case PopFront:
				if len(oracle) == 0 {
					continue
				}
				oracleItem := oracle[0]
				t.Logf("PopFront() -> %#v", oracleItem)
				oracle = oracle[1:]
				dequeItem := deque.PopFront()
				require.Equal(t, oracleItem, dequeItem)
			case PopBack:
				if len(oracle) == 0 {
					continue
				}
				oracleItem := oracle[len(oracle)-1]
				t.Logf("PopBack() -> %#v", oracleItem)
				oracle = oracle[:len(oracle)-1]
				dequeItem := deque.PopBack()
				require.Equal(t, oracleItem, dequeItem)
			case PeekFront:
				if len(oracle) == 0 {
					continue
				}
				oracleItem := oracle[0]
				t.Logf("PeekFront() -> %#v", oracleItem)
				dequeItem := deque.PeekFront()
				require.Equal(t, oracleItem, dequeItem)
			case PeekBack:
				if len(oracle) == 0 {
					continue
				}
				oracleItem := oracle[len(oracle)-1]
				t.Logf("PeekBack() -> %#v", oracleItem)
				dequeItem := deque.PeekBack()
				require.Equal(t, oracleItem, dequeItem)
			case Iterate:
				oracleAll := oracle
				if len(oracleAll) == 0 {
					oracleAll = nil
				}
				dequeAll := iterator.Collect(deque.Iterate())
				if len(dequeAll) == 0 {
					dequeAll = nil
				}
				require.Equal(t, oracleAll, dequeAll)
			}

			t.Logf("oracle state: %#v", oracle)
			t.Logf("deque state:  %#v", deque)
		}
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
