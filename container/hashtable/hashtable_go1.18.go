//go:build go1.18

package hashtable

import (
	"hash/maphash"
	"math/bits"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/bradenaw/juniper/xmath"
)

// hashtable is a Swiss table. It's a regular probing hashtable, with just a little rearranging from
// the usual in order to increase cache locality and to use SIMD intrinsics to dramatically speed
// probing.
//
// https://abseil.io/docs/cpp/guides/container
// https://www.youtube.com/watch?v=ncHmEUmJZf4
type hashtable[K any, V any] struct {
	hash func(K) uint64
	eq   func(K, K) bool

	groups []group[K, V]
	size   int
}

func defaultEq[K comparable](a, b K) bool { return a == b }
func defaultHash[K any]() func(k K) uint64 {
	s := maphash.MakeSeed()
	var zero K
	l := int(unsafe.Sizeof(zero))
	return func(k K) uint64 {
		var h maphash.Hash
		h.SetSeed(s)
		ak := &k
		var b []byte
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
		hdr.Data = uintptr(unsafe.Pointer(ak))
		hdr.Len = l
		hdr.Cap = l
		h.Write(b)
		runtime.KeepAlive(ak)
		return h.Sum64()
	}
}

func newHashtable[K any, V any](
	hash func(K) uint64,
	eq func(K, K) bool,
) *hashtable[K, V] {
	return &hashtable[K, V]{
		hash: hash,
		eq:   eq,
	}
}

const (
	groupSize  = 16
	maxLFNum   = 7
	maxLFDenom = 10
)

const (
	controlEmpty   uint8 = 0b10000000
	controlDeleted uint8 = 0b11111110
	// controlSentinel uint8 = 0b11111111
	// controlFull  uint8 = 0b0xxxxxxx
)

type pair[K any, V any] struct {
	key   K
	value V
}

type group[K any, V any] struct {
	meta  [groupSize]uint8
	pairs [groupSize]pair[K, V]
}

func (ht *hashtable[K, V]) Len() int {
	return ht.size
}

func (ht *hashtable[K, V]) Put(key K, value V) {
	if ht.shouldGrow() {
		ht.grow()
	}
	h2, control, p, found := ht.find(key, ht.groups)
	if !found {
		*control = h2
		*p = pair[K, V]{key, value}
		ht.size++
	} else {
		p.value = value
	}
}

func (ht *hashtable[K, V]) Get(key K) V {
	var zero V
	if ht.size == 0 {
		return zero
	}
	_, _, p, found := ht.find(key, ht.groups)
	if !found {
		return zero
	}
	return p.value
}

// returns the h2 of key, and pointers into a control byte and pair. If found is true, then this is
// where the key is located. If found is false, then this is a suitable location to insert it.
func (ht *hashtable[K, V]) find(
	key K,
	groups []group[K, V],
) (h2_ uint8, control_ *uint8, pair_ *pair[K, V], found_ bool) {
	groupIdx, h2 := ht.splitHash(key, len(groups))
	j := groupIdx
	for {
		g := &groups[j]
		m := matchMask(h2, g.meta)
		iter := maskIter{m: m}
		for {
			i, ok := iter.Next()
			if !ok {
				break
			}
			if ht.eq(key, g.pairs[i].key) {
				return h2, &g.meta[i], &g.pairs[i], true
			}
		}
		// If there's an empty slot in this group, that means key can't have gotten into any other
		// group by probing on put, so we're done.
		m = matchMask(controlEmpty, g.meta)
		if m != 0 {
			// should actually prefer to overwrite a tombstone here
			i := bits.TrailingZeros16(m)
			return h2, &g.meta[i], &g.pairs[i], false
		}
		// Move on to the next group.
		j = (j + 1) % len(groups)
		// If we looped all the way around, we're done.
		if j == groupIdx {
			break
		}
	}
	return h2, nil, nil, false
}

func (ht *hashtable[K, V]) grow() {
	newGroups := make([]group[K, V], xmath.Max(len(ht.groups)*2, 1))
	for i := range newGroups {
		newGroups[i].meta = [groupSize]uint8{
			controlEmpty, controlEmpty, controlEmpty, controlEmpty,
			controlEmpty, controlEmpty, controlEmpty, controlEmpty,
			controlEmpty, controlEmpty, controlEmpty, controlEmpty,
			controlEmpty, controlEmpty, controlEmpty, controlEmpty,
		}
	}

	for _, g := range ht.groups {
		for i := range g.meta {
			if g.meta[i]&0x80 != 0 {
				continue
			}

			// found is always false unless the hash function is lying to us
			h2, control, p, _ := ht.find(g.pairs[i].key, newGroups)
			*control = h2
			*p = g.pairs[i]
		}
	}

	ht.groups = newGroups
}

func (ht *hashtable[K, V]) shouldGrow() bool {
	// size / nSlots > maxLFNum / maxLFDenom
	// with a little rearranging to use int math
	return ht.size == 0 || (ht.size+1)*maxLFDenom/(len(ht.groups)*groupSize) > maxLFNum
}

// bit i (counting from lsb) in the output is 1 if b[i] == a
func matchMask(a uint8, b [16]uint8) uint16

func matchMaskLoop(a uint8, b [16]uint8) uint16 {
	out := uint16(0)
	for i := range b {
		if a == b[i] {
			out |= 1 << i
		}
	}
	return out
}

// iterator over the indexes of the set bits in m
type maskIter struct {
	m uint16
	i int
}

func (iter *maskIter) Next() (int, bool) {
	tz := bits.TrailingZeros16(iter.m)
	if tz == 16 {
		return 0, false
	}
	iter.i += tz
	out := iter.i
	iter.m = iter.m >> (tz + 1)
	iter.i += 1
	return out, true
}

func (ht *hashtable[K, V]) splitHash(key K, nGroups int) (groupIdx_ int, h2_ uint8) {
	h := ht.hash(key)
	h1 := (h >> 7) & 0xFFFFFFFF
	h2 := uint8(h & 0x7F)
	return int((h1 * uint64(nGroups)) >> 32), h2
}
