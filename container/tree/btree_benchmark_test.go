package tree

import (
	"cmp"
	"fmt"
	"testing"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath/xrand"
)

// Run under WSL since I don't have a native Linux machine handy at the moment.
//
// Obviously not 100% scientific since there are more dimensions than represented here. The builtin
// map reallocates at different points, so its alloc/op is sort of random numbers depending on how
// well the benchmark size fits. These are all for int keys and values using the builtin < to
// compare, and things may shift for different sized types and differently complex comparisons.
//
// However, in this small study, the btree map is about half as fast as the builtin hashmap for gets
// and puts, and will usually use less memory. The builtin is (expectedly) about 4.5x faster at
// replacing keys that are already present.
//
// The purpose was to pick branchFactor. Too small and we waste more space and, with more
// allocations, put more strain on GC. Too large we spend more time searching and shifting inside
// nodes. branchFactor=32 requires a similar number of objects to the builtin map, which may have
// visible advantages for GC in a real program. Unfortunately, we'll have to wait until there's a
// real program using it to find out which does better. branchFactor=16 seems like a decent balance
// for now, it's nearly as small memory-wise and is a little faster than branchFactor=32 at nearly
// everything.
//
// In addition to drastically reducing allocations, this B-tree implementation drastically
// outperformed a now-removed AVL tree implementation on everything except Get, which branchFactor=4
// is about 10% slower than and branchFactor=16 is about 20% slower than. Outperformance on writes
// is unsurprising, since allocation is a significant cost. Reads get better cache locality with the
// B-tree, but have to call less more.
//
//
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz
//
// benchmark          size           builtin map       btree                                                                 //
//                                                     branchFactor=4    branchFactor=8    branchFactor=16   branchFactor=32 //
// time ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────── //
// Get                10              10.7ns            15.03 ns/op       17.00 ns/op       21.27 ns/op       21.25 ns/op    //
//                    100             11.2ns            25.79 ns/op       25.93 ns/op       36.56 ns/op       39.52 ns/op    //
//                    1000            16.2ns            48.51 ns/op       50.57 ns/op       58.23 ns/op       67.26 ns/op    //
//                    10000           22.9ns            63.11 ns/op       63.58 ns/op       71.59 ns/op       84.33 ns/op    //
//                    100000          28.1ns            70.84 ns/op       74.67 ns/op       83.78 ns/op      102.1 ns/op     //
//                    1000000         51.5ns            81.46 ns/op       89.05 ns/op       94.68 ns/op      116.0 ns/op     //
// Put                10              42.6ns            54.55 ns/op       49.36 ns/op       41.41 ns/op       45.82 ns/op    //
//                    100             55.8ns            72.56 ns/op       63.41 ns/op       65.07 ns/op       88.99 ns/op    //
//                    1000            65.0ns           127.0 ns/op       107.4 ns/op       108.3 ns/op       123.6 ns/op     //
//                    10000           60.0ns           158.3 ns/op       131.2 ns/op       135.7 ns/op       154.2 ns/op     //
//                    100000          61.9ns           217.6 ns/op       183.6 ns/op       181.5 ns/op       203.9 ns/op     //
//                    1000000        112ns             386.5 ns/op       308.8 ns/op       264.2 ns/op       297.2 ns/op     //
// PutAlreadyPresent  10              13.7ns            18.19 ns/op       17.68 ns/op       21.92 ns/op       21.40 ns/op    //
//                    100             15.2ns            27.27 ns/op       28.76 ns/op       37.43 ns/op       44.87 ns/op    //
//                    1000            20.3ns            65.64 ns/op       56.90 ns/op       62.54 ns/op       75.38 ns/op    //
//                    10000           25.8ns           104.3 ns/op        86.67 ns/op       86.52 ns/op       107.1 ns/op    //
//                    100000          32.1ns           155.1 ns/op       126.1 ns/op       127.8 ns/op       142.7 ns/op     //
//                    1000000         62.4ns           385.3 ns/op       293.5 ns/op       281.0 ns/op       270.7 ns/op     //
// Iterate            10             125.9 ns/op       171.8 ns/op       160.1 ns/op       150.3 ns/op       150.3 ns/op     //
//                    100           1068 ns/op        1611 ns/op        1257 ns/op        1058 ns/op         947.4 ns/op     //
//                    1000         12572 ns/op       14581 ns/op       12062 ns/op       10375 ns/op        9348 ns/op       //
//                    10000       109337 ns/op      189691 ns/op      142145 ns/op      112312 ns/op       96263 ns/op       //
//                    100000     1028813 ns/op     2060605 ns/op     1447266 ns/op     1145753 ns/op      984982 ns/op       //
//                    1000000   13181522 ns/op    66242895 ns/op    33671273 ns/op    19829590 ns/op    14673612 ns/op       //
//                                                                                                                           //
// alloc bytes ───────────────────────────────────────────────────────────────────────────────────────────────────────────── //
// Put                10              48 B/op           50 B/op           60 B/op           40 B/op           79 B/op        //
//                    100             55 B/op           47 B/op           40 B/op           34 B/op           46 B/op        //
//                    1000            86 B/op           50 B/op           40 B/op           36 B/op           37 B/op        //
//                    10000           68 B/op           49 B/op           40 B/op           37 B/op           35 B/op        //
//                    100000          57 B/op           50 B/op           40 B/op           37 B/op           36 B/op        //
//                    1000000         86 B/op           49 B/op           40 B/op           37 B/op           35 B/op        //
//                                                                                                                           //
// alloc objects ─────────────────────────────────────────────────────────────────────────────────────────────────────────── //
// Build              10                1                3                 1                 1                 1             //
//                    100              16               50                17                 8                 3             //
//                    1000             64              518               170                96                37             //
//                    10000           276                5.1K              1.7K            971               375             //
//                    100000            4.0K            51K               16.9K              9.6K              3.7K          //
//                    1000000          38.0K           514K              169K               96K               37K            //

var sizes = []int{10, 100, 1_000, 10_000, 100_000, 1_000_000}

func BenchmarkBtreeMapGet(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int](cmp.Less[int])
		for i := 0; i < size; i++ {
			m.Put(i, i)
		}
		b.Run(fmt.Sprintf("Size=%d,BranchFactor=%d", size, branchFactor), func(b *testing.B) {
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
				_ = m[i%size]
			}
		})
	}
}

func BenchmarkBtreeMapPut(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int](cmp.Less[int])
		keys := iterator.Collect(iterator.Counter(size))
		xrand.Shuffle(keys)

		b.Run(fmt.Sprintf("Size=%d,BranchFactor=%d", size, branchFactor), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m.Put(keys[i%size], i)
				if m.Len() == size {
					m = NewMap[int, int](cmp.Less[int])
				}
			}
		})
	}
}

func BenchmarkBuiltinMapPut(b *testing.B) {
	for _, size := range sizes {
		m := make(map[int]int)
		keys := iterator.Collect(iterator.Counter(size))
		xrand.Shuffle(keys)

		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m[keys[i%size]] = i
				if len(m) == size {
					m = make(map[int]int)
				}
			}
		})
	}
}

func BenchmarkBtreeMapPutAlreadyPresent(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int](cmp.Less[int])
		keys := iterator.Collect(iterator.Counter(size))
		xrand.Shuffle(keys)
		for _, k := range keys {
			m.Put(k, 0)
		}

		b.Run(fmt.Sprintf("Size=%d,BranchFactor=%d", size, branchFactor), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m.Put(keys[i%size], i)
			}
		})
	}
}

func BenchmarkBuiltinMapPutAlreadyPresent(b *testing.B) {
	for _, size := range sizes {
		m := make(map[int]int)
		keys := iterator.Collect(iterator.Counter(size))
		xrand.Shuffle(keys)
		for _, k := range keys {
			m[k] = 0
		}

		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m[keys[i%size]] = i
			}
		})
	}
}

func BenchmarkBtreeMapIterate(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int](cmp.Less[int])
		keys := iterator.Collect(iterator.Counter(size))
		xrand.Shuffle(keys)
		for i, k := range keys {
			m.Put(k, i)
		}

		b.Run(fmt.Sprintf("Size=%d,BranchFactor=%d", size, branchFactor), func(b *testing.B) {
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
		keys := iterator.Collect(iterator.Counter(size))
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

func BenchmarkBtreeMapBuild(b *testing.B) {
	for _, size := range sizes {
		keys := iterator.Collect(iterator.Counter(size))

		b.Run(fmt.Sprintf("Size=%d,BranchFactor=%d", size, branchFactor), func(b *testing.B) {
			b.ReportAllocs()
			for i := 1; i < b.N; i++ {
				b.StopTimer()
				xrand.Shuffle(keys)
				m := NewMap[int, int](cmp.Less[int])
				b.StartTimer()

				for j := 0; j < size; j++ {
					m.Put(keys[j], j)
				}
			}
		})
	}
}

func BenchmarkBuiltinMapBuild(b *testing.B) {
	for _, size := range sizes {
		keys := iterator.Collect(iterator.Counter(size))

		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 1; i < b.N; i++ {
				b.StopTimer()
				xrand.Shuffle(keys)
				m := make(map[int]int)
				b.StartTimer()

				for j := 0; j < size; j++ {
					m[keys[j]] = j
				}
			}
		})
	}
}
