//go:build go1.18

package hashtable

import (
	"hash/maphash"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bradenaw/juniper/internal/fuzz"
	"github.com/bradenaw/juniper/xmath/xrand"
	"github.com/bradenaw/juniper/iterator"
)

func FuzzHashtable(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		h := maphash.Hash{}
		h.SetSeed(maphash.MakeSeed())
		ht := newHashtable[byte, int](
			func(key byte) uint64 {
				h.WriteByte(key)
				out := h.Sum64()
				h.Reset()
				return out
			},
			func(a, b byte) bool {
				return a == b
			},
		)
		oracle := make(map[byte]int)

		ctr := 0
		fuzz.Operations(
			b,
			func() {
				t.Log(tableToString(ht))

				require.Equal(t, len(oracle), ht.Len())
			},
			func(k byte) {
				v := ctr
				t.Logf("Put(%#v, %#v)", k, v)
				ctr++

				ht.Put(k, v)
				oracle[k] = v
			},
			func(k byte) {
				oracleV := oracle[k]
				t.Logf("Get(%#v) -> %#v", k, oracleV)

				v := ht.Get(k)
				require.Equal(t, oracleV, v)
			},
		)
	})
}

func TestMaskIter(t *testing.T) {
	iter := maskIter{m: 0b0001_1000_0110_1001}
	var ints []int
	for {
		t.Logf("m = %016b", iter.m)
		t.Logf("i = %d", iter.i)
		i, ok := iter.Next()
		if !ok {
			break
		}
		ints = append(ints, i)
	}
	require.Equal(t, []int{0, 3, 5, 6, 11, 12}, ints)
}

func TestMatchMask(t *testing.T) {
	check := func(expected uint16, a uint8, b [16]uint8) {
		actual := matchMask(a, b)
		require.Equalf(t, expected, actual, "expected %016b, got %016b", expected, actual)
	}
	check(0b0000_0000_0000_0001, 0x11, [16]uint8{
		0x11, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})
	check(0b1010_0101_1000_0001, 0x11, [16]uint8{
		0x11, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x11,
		0x11, 0x00, 0x11, 0x00,
		0x00, 0x11, 0x00, 0x11,
	})
}

func BenchmarkMatchMask(b *testing.B) {
	b.StopTimer()
	items := make([][16]uint8, 1024)
	for i := range items {
		rand.Read(items[i][:])
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		matchMask(byte(i), items[i%len(items)])
	}
}

func BenchmarkMatchMaskLoop(b *testing.B) {
	b.StopTimer()
	items := make([][16]uint8, 1024)
	for i := range items {
		rand.Read(items[i][:])
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		matchMaskLoop(byte(i), items[i%len(items)])
	}
}

func tableToString[K any, V any](ht *hashtable[K, V]) string {
	var sb strings.Builder

	for i, g := range ht.groups {
		fmt.Fprintf(&sb, "group %d\n", i)
		for j, control := range g.meta {
			fmt.Fprintf(&sb, "  slot %d  ", j)
			switch control {
			case controlEmpty:
				fmt.Fprintf(&sb, "[empty]")
			case controlDeleted:
				fmt.Fprintf(&sb, "[deleted]")
			default:
				fmt.Fprintf(&sb, "h2=%02x  %#v: %#v", control, g.pairs[j].key , g.pairs[j].value )
			}
			fmt.Fprintf(&sb, "\n")
		}
	}
	return sb.String()
}

var sizes = []int{10, 100, 1_000, 10_000, 100_000, 1_000_000}

func BenchmarkHashtableMapGet(b *testing.B) {
	for _, size := range sizes {
		m := newHashtable[int, int](defaultHash[int](), defaultEq[int])
		for i := 0; i < size; i++ {
			m.Put(i, i)
		}
		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = m.Get(i % size)
			}
		})
	}
}


func BenchmarkHashtableMapPut(b *testing.B) {
	for _, size := range sizes {
		m := newHashtable[int, int](defaultHash[int](), defaultEq[int])
		keys := iterator.Collect(iterator.Counter(size))
		xrand.Shuffle(keys)

		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m.Put(keys[i%size], i)
				if m.Len() == size {
					m = newHashtable[int, int](defaultHash[int](), defaultEq[int])
				}
			}
		})
	}
}


func BenchmarkHashtableMapPutAlreadyPresent(b *testing.B) {
	for _, size := range sizes {
		m := newHashtable[int, int](defaultHash[int](), defaultEq[int])
		keys := iterator.Collect(iterator.Counter(size))
		xrand.Shuffle(keys)
		for _, k := range keys {
			m.Put(k, 0)
		}

		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m.Put(keys[i%size], i)
			}
		})
	}
}
