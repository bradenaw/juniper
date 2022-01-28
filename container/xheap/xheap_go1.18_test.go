//go:build go1.18

package xheap

import (
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

func FuzzHeap(f *testing.F) {
	f.Fuzz(func(t *testing.T, b1 []byte, b2 []byte) {
		t.Logf("initial: %#v", b1)
		t.Logf("pushed:  %#v", b2)
		h := New(xsort.OrderedLess[byte], append([]byte{}, b1...))
		for i := range b2 {
			h.Push(b2[i])
		}

		outByIterate := iterator.Collect(h.Iterate())
		xsort.Slice(outByIterate, xsort.OrderedLess[byte])
		if outByIterate == nil {
			outByIterate = []byte{}
		}

		outByPop := []byte{}
		for h.Len() > 0 {
			item := h.Pop()
			outByPop = append(outByPop, item)
		}

		expected := append(append([]byte{}, b1...), b2...)
		t.Logf("expected:        %#v", expected)
		xsort.Slice(expected, xsort.OrderedLess[byte])
		t.Logf("expected sorted: %#v", expected)

		require2.SlicesEqual(t, expected, outByPop)
		require2.SlicesEqual(t, expected, outByIterate)
	})
}

func FuzzPriorityQueue(f *testing.F) {
	const (
		Update = iota
		Pop
		Peek
		Contains
		Priority
		Remove
		Iterate
	)
	f.Fuzz(func(t *testing.T, b1 []byte, b2 []byte) {
		initial := make([]KP[int, float32], 0, len(b1))
		oracle := make(map[int]float32)
		for i := range b1 {
			k := int((b1[i] & 0b00011100) >> 2)
			p := float32((b1[i] & 0b00000011))

			_, ok := oracle[k]
			if ok {
				continue
			}

			initial = append(initial, KP[int, float32]{k, p})
			oracle[k] = p
		}
		t.Logf("initial:        %#v", initial)
		t.Logf("initial oracle: %#v", oracle)

		h := NewPriorityQueue(xsort.OrderedLess[float32], initial)

		oracleLowestP := func() float32 {
			first := true
			lowest := float32(0)
			for _, p := range oracle {
				if first || p < lowest {
					lowest = p
				}
				first = false
			}
			return lowest
		}

		for _, b := range b2 {
			op := (b & 0b11100000) >> 5
			k := int((b & 0b00011100) >> 2)
			p := float32(b & 0b00000011)

			switch op {
			case Update:
				t.Logf("Update(%d, %f)", k, p)
				oracle[k] = p
				h.Update(k, p)
			case Pop:
				t.Logf("Pop()")
				if len(oracle) == 0 {
					require2.Equal(t, 0, h.Len())
					continue
				}
				lowestP := oracleLowestP()
				hPopped := h.Pop()
				require2.Equal(t, lowestP, oracle[hPopped])
				delete(oracle, hPopped)
			case Peek:
				t.Logf("Peek()")
				if len(oracle) == 0 {
					require2.Equal(t, 0, h.Len())
					continue
				}
				lowestP := oracleLowestP()
				hPeeked := h.Peek()
				require2.Equal(t, lowestP, oracle[hPeeked])
			case Contains:
				t.Logf("Contains(%d)", k)
				_, oracleContains := oracle[k]
				require2.Equal(t, oracleContains, h.Contains(k))
			case Priority:
				t.Logf("Priority(%d)", k)
				require2.Equal(t, oracle[k], h.Priority(k))
			case Remove:
				t.Logf("Remove(%d)", k)
				delete(oracle, k)
				h.Remove(k)
			case Iterate:
				t.Logf("Iterate()")
				oracleItems := make([]int, 0, len(oracle))
				for k := range oracle {
					oracleItems = append(oracleItems, k)
				}
				items := iterator.Collect(h.Iterate())
				require2.ElementsMatch(t, oracleItems, items)
			}

			require2.Equal(t, len(oracle), h.Len())
		}
	})
}
