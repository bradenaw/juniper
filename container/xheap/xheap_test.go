package xheap

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/xstd/iterator"
	"github.com/bradenaw/xstd/xsort"
)

func FuzzBasic(f *testing.F) {
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
		for {
			item, ok := h.Pop()
			if !ok {
				break
			}
			outByPop = append(outByPop, item)
		}

		expected := append(append([]byte{}, b1...), b2...)
		t.Logf("expected:        %#v", expected)
		xsort.Slice(expected, xsort.OrderedLess[byte])
		t.Logf("expected sorted: %#v", expected)

		require.Equal(t, expected, outByPop)
		require.Equal(t, expected, outByIterate)
	})
}
