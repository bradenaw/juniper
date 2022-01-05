//go:build go1.18

package tree

import (
	"fmt"
	"testing"

	"github.com/bradenaw/juniper/xmath/xrand"
	"github.com/bradenaw/juniper/xsort"
	"github.com/bradenaw/juniper/iterator"
)

var sizes = []int{10, 100, 1_000, 10_000, 100_000, 1_000_000}

func BenchmarkTreeMapGet(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int](xsort.OrderedLess[int])
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

func BenchmarkBuiltinMapGet(b *testing.B) {
	for _, size := range sizes {
		m := make(map[int]int, size)
		for i := 0; i < size; i++ {
			m[i] = i
		}
		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = m[i % size]
			}
		})
	}
}

func BenchmarkTreeMapPut(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int](xsort.OrderedLess[int])
		keys := iterator.Collect(iterator.Count(size))
		xrand.Shuffle(keys)
		
		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m.Put(keys[i % size], i)
				if m.Len() == size {
					m = NewMap[int, int](xsort.OrderedLess[int])
				}
			}
		})
	}
}

func BenchmarkBuiltinMapPut(b *testing.B) {
	for _, size := range sizes {
		m := make(map[int]int)
		keys := iterator.Collect(iterator.Count(size))
		xrand.Shuffle(keys)
		
		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m[keys[i % size]] = i
				if len(m) == size {
					m = make(map[int]int)
				}
			}
		})
	}
}

func BenchmarkTreeMapPutFull(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int](xsort.OrderedLess[int])
		keys := iterator.Collect(iterator.Count(size))
		xrand.Shuffle(keys)
		for _, k := range keys {
			m.Put(k, 0)
		}
		
		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m.Put(keys[i % size], i)
			}
		})
	}
}

func BenchmarkBuiltinMapPutFull(b *testing.B) {
	for _, size := range sizes {
		m := make(map[int]int)
		keys := iterator.Collect(iterator.Count(size))
		xrand.Shuffle(keys)
		for _, k := range keys {
			m[k] = 0
		}
		
		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m[keys[i % size]] = i
			}
		})
	}
}

func BenchmarkTreeMapIterate(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int](xsort.OrderedLess[int])
		keys := iterator.Collect(iterator.Count(size))
		xrand.Shuffle(keys)
		for i, k := range keys {
			m.Put(k, i)
		}
		
		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				iter := m.Iterate()
				for {
					_, ok := iter.Next()
					if !ok {
						break
					}
				}
			}
		})
	}
}

func BenchmarkBuiltinMapIterate(b *testing.B) {
	for _, size := range sizes {
		m := make(map[int]int)
		keys := iterator.Collect(iterator.Count(size))
		xrand.Shuffle(keys)
		for i, k := range keys {
			m[k] = i
		}
		
		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, _ = range m {
				}
			}
		})
	}
}
