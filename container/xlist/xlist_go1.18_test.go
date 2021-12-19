//go:build go1.18

package xlist

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/slices"
)

func FuzzList(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		var l List[int]
		var oracle []int

		nodeAt := func(i int) *Node[int] {
			j := 0
			curr := l.Front()
			for j < i {
				curr = curr.Next()
				j++
			}
			return curr
		}

		fuzz.Operations(
			b,
			func() { // check
				t.Logf("%v", oracle)
				require.Equal(t, len(oracle), l.Len())

				if len(oracle) == 0 {
					require.Nil(t, l.Front())
					require.Nil(t, l.Back())
					return
				}

				node := l.Front()
				for i := range oracle {
					require.NotNilf(t, node, "node nil at index %d, len(oracle)==%d", i, len(oracle))
					require.Equal(t, oracle[i], node.Value)
					if node.Next() != nil {
						require.Equal(t, node, node.Next().Prev())
					}
					node = node.Next()
				}
				require.Nil(t, node)
				require.NotNil(t, l.Back())
				require.Equal(t, oracle[len(oracle)-1], l.Back().Value)
			},
			func(value int) {
				t.Logf("PushFront(%d)", value)
				l.PushFront(value)
				oracle = append([]int{value}, oracle...)
			},
			func(value int) {
				t.Logf("PushBack(%d)", value)
				l.PushBack(value)
				oracle = append(oracle, value)
			},
			func(value int, idx int) {
				if len(oracle) == 0 || idx < 0 {
					return
				}
				idx = idx % len(oracle)
				t.Logf("InsertBefore(%d, node @ %d)", value, idx)
				l.InsertBefore(value, nodeAt(idx))
				oracle = slices.Insert(oracle, idx, value)
			},
			func(value int, idx int) {
				if len(oracle) == 0 || idx < 0 {
					return
				}
				idx = idx % len(oracle)
				t.Logf("InsertAfter(%d, node @ %d)", value, idx)
				l.InsertAfter(value, nodeAt(idx))
				oracle = slices.Insert(oracle, idx+1, value)
			},
			func(idx int) {
				if len(oracle) == 0 || idx < 0 {
					return
				}
				idx = idx % len(oracle)
				t.Logf("Remove(node @ %d)", idx)
				l.Remove(nodeAt(idx))
				oracle = slices.Remove(oracle, idx, 1)
			},
			func(src, dest int) {
				if len(oracle) == 0 || src < 0 || dest < 0 {
					return
				}
				src = src % len(oracle)
				dest = dest % len(oracle)
				t.Logf("MoveBefore(node @ %d, node @ %d)", src, dest)
				l.MoveBefore(nodeAt(src), nodeAt(dest))
				item := oracle[src]
				oracle = slices.Remove(oracle, src, 1)
				if dest > src {
					dest--
				}
				oracle = slices.Insert(oracle, dest, item)
			},
			func(src, dest int) {
				if len(oracle) == 0 || src < 0 || dest < 0 {
					return
				}
				src = src % len(oracle)
				dest = dest % len(oracle)
				t.Logf("MoveAfter(node @ %d, node @ %d)", src, dest)
				l.MoveAfter(nodeAt(src), nodeAt(dest))
				item := oracle[src]
				oracle = slices.Remove(oracle, src, 1)
				if dest >= src {
					dest--
				}
				oracle = slices.Insert(oracle, dest+1, item)
			},
			func (idx int) {
				if len(oracle) == 0 || idx < 0 {
					return
				}
				idx = idx % len(oracle)
				t.Logf("MoveToFront(node @ %d)", idx)
				l.MoveToFront(nodeAt(idx))
				item := oracle[idx]
				oracle = slices.Remove(oracle, idx, 1)
				oracle = append([]int{item}, oracle...)
			},
			func (idx int) {
				if len(oracle) == 0 || idx < 0 {
					return
				}
				idx = idx % len(oracle)
				t.Logf("MoveToBack(node @ %d)", idx)
				l.MoveToBack(nodeAt(idx))
				item := oracle[idx]
				oracle = slices.Remove(oracle, idx, 1)
				oracle = append(oracle, item)
			},
		)
	})
}
