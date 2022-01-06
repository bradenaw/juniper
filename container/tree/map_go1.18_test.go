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

// Results as of 9e9150ea8f67082486a73c45ee0b8846246c7e96
// Run under WSL since I don't have a native Linux machine handy at the moment.
//
// TL;DR:
// tree.Map is slower. It's within about a factor of 2 for Get(). It's much slower to write to, and
// the gap is (expectedly) larger for larger maps. It only makes any sense to use if you need
// ordered iteration while the map is also evolving. If you only need ordered iteration
// occasionally compared to how often the map changes, you might still be better off using the
// builtin map and sorting the keys when necessary.
//
// tree.Map is for now not very optimized, however. Writes are recursive because they're easiest to
// implement that way - they can be done non-recursively by adding a parent pointer (would also
// improve Iterator speed since it wouldn't need to keep its own stack, but increase memory usage by
// one more pointer per item) or by self-allocating a stack of visited nodes (may or may not be
// faster because of the heap allocation, but rebalancing can terminate earlier).
//
//
//
// results="$(mktemp)"
// echo "$results"
// for i in {0..5}; do
//   go test --bench . github.com/bradenaw/juniper/container/tree | tee --append "$results"
// done
// go run golang.org/x/perf/cmd/benchstat@latest \
//   <(grep -E "^BenchmarkBuiltin" "$results" | sed -E "s/^BenchmarkBuiltin/Benchmark/g") \
//   <(grep -E "^BenchmarkTree" "$results" | sed -E "s/^BenchmarkTree/Benchmark/g")
//
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz
//
//
// ("old" means builtin, and "new" means tree.Map)
//
// name                                  old time/op    new time/op    delta
// MapGet/Size=10-16                       11.0ns ± 2%    12.6ns ± 0%   +14.93%  (p=0.002 n=6+6)
// MapGet/Size=100-16                      11.2ns ± 1%    18.6ns ± 1%   +65.86%  (p=0.002 n=6+6)
// MapGet/Size=1000-16                     16.5ns ± 3%    42.9ns ± 0%  +160.52%  (p=0.004 n=6+5)
// MapGet/Size=10000-16                    23.0ns ± 1%    50.2ns ± 1%  +118.59%  (p=0.004 n=5+6)
// MapGet/Size=100000-16                   28.2ns ± 1%    58.3ns ± 1%  +106.75%  (p=0.004 n=5+6)
// MapGet/Size=1000000-16                  51.4ns ± 1%    71.1ns ± 1%   +38.20%  (p=0.008 n=5+5)
// MapPut/Size=10-16                       42.0ns ± 1%    56.4ns ± 1%   +34.27%  (p=0.002 n=6+6)
// MapPut/Size=100-16                      55.3ns ± 0%    82.9ns ± 2%   +49.90%  (p=0.002 n=6+6)
// MapPut/Size=1000-16                     64.1ns ± 1%   138.5ns ± 3%  +116.21%  (p=0.002 n=6+6)
// MapPut/Size=10000-16                    59.3ns ± 1%   194.4ns ± 1%  +227.99%  (p=0.002 n=6+6)
// MapPut/Size=100000-16                   60.5ns ± 1%   306.9ns ± 1%  +407.65%  (p=0.002 n=6+6)
// MapPut/Size=1000000-16                   111ns ± 1%     598ns ± 3%  +436.55%  (p=0.002 n=6+6)
// MapPutAlreadyPresent/Size=10-16         13.4ns ± 1%    23.5ns ± 1%   +75.29%  (p=0.004 n=5+6)
// MapPutAlreadyPresent/Size=100-16        15.1ns ± 4%    44.6ns ± 1%  +196.38%  (p=0.002 n=6+6)
// MapPutAlreadyPresent/Size=1000-16       20.2ns ± 3%    88.8ns ± 0%  +340.59%  (p=0.002 n=6+6)
// MapPutAlreadyPresent/Size=10000-16      26.4ns ± 1%   136.4ns ± 1%  +417.21%  (p=0.004 n=5+6)
// MapPutAlreadyPresent/Size=100000-16     32.2ns ± 0%   220.8ns ± 1%  +584.94%  (p=0.004 n=5+6)
// MapPutAlreadyPresent/Size=1000000-16    62.8ns ± 1%   583.0ns ± 3%  +828.83%  (p=0.002 n=6+6)
// MapIterate/Size=10-16                   81.4ns ± 2%   302.7ns ± 3%  +271.62%  (p=0.002 n=6+6)
// MapIterate/Size=100-16                   723ns ± 2%    1146ns ± 2%   +58.56%  (p=0.002 n=6+6)
// MapIterate/Size=1000-16                 8.67µs ± 4%   10.45µs ± 1%   +20.56%  (p=0.002 n=6+6)
// MapIterate/Size=10000-16                78.6µs ± 1%   161.0µs ± 1%  +104.80%  (p=0.004 n=6+5)
// MapIterate/Size=100000-16                742µs ± 1%    1967µs ± 0%  +165.06%  (p=0.002 n=6+6)
// MapIterate/Size=1000000-16              9.23ms ± 1%   83.15ms ± 2%  +800.98%  (p=0.002 n=6+6)
//
//
// tree.Map appears to use a little less memory than the builtin map:
//
// name                                  old alloc/op   new alloc/op   delta
// MapPut/Size=10-16                        48.0B ± 0%     51.0B ± 0%    +6.25%  (p=0.002 n=6+6)
// MapPut/Size=100-16                       55.0B ± 0%     48.0B ± 0%   -12.73%  (p=0.002 n=6+6)
// MapPut/Size=1000-16                      86.0B ± 0%     48.0B ± 0%   -44.19%  (p=0.002 n=6+6)
// MapPut/Size=10000-16                     68.0B ± 0%     48.0B ± 0%   -29.41%  (p=0.002 n=6+6)
// MapPut/Size=100000-16                    57.0B ± 0%     48.0B ± 0%   -15.79%  (p=0.002 n=6+6)
// MapPut/Size=1000000-16                   87.0B ± 0%     48.0B ± 0%   -44.83%  (p=0.026 n=5+6)
//
// name                      old alloc/op   new alloc/op   delta
// MapBuild/Size=10-16           292B ± 0%      480B ± 0%    +64.38%  (p=0.004 n=6+5)
// MapBuild/Size=100-16        5.36kB ± 0%    4.80kB ± 0%    -10.53%  (p=0.002 n=6+6)
// MapBuild/Size=1000-16       86.5kB ± 0%    48.0kB ± 0%    -44.55%  (p=0.002 n=6+6)
// MapBuild/Size=10000-16       687kB ± 0%     479kB ± 0%    -30.22%  (p=0.002 n=6+6)
// MapBuild/Size=100000-16     5.74MB ± 0%    4.75MB ± 0%    -17.16%  (p=0.002 n=6+6)
// MapBuild/Size=1000000-16    86.7MB ± 0%    47.5MB ± 0%    -45.22%  (p=0.010 n=4+6)
//
//
// Although tree.Map requires a _lot_ more allocations, so is likely more costly for the GC to keep
// track of:
//
// name                      old allocs/op  new allocs/op  delta
// MapBuild/Size=10-16           1.00 ± 0%      9.50 ± 5%   +850.00%  (p=0.002 n=6+6)
// MapBuild/Size=100-16          16.0 ± 0%      99.0 ± 0%   +518.75%  (p=0.002 n=6+6)
// MapBuild/Size=1000-16         64.0 ± 0%     999.0 ± 0%  +1460.94%  (p=0.002 n=6+6)
// MapBuild/Size=10000-16         275 ± 0%      9984 ± 0%  +3526.15%  (p=0.002 n=6+6)
// MapBuild/Size=100000-16      3.99k ± 0%    99.00k ± 0%  +2379.66%  (p=0.002 n=6+6)
// MapBuild/Size=1000000-16     38.0k ± 0%    990.0k ± 0%  +2503.11%  (p=0.004 n=5+6)

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

func BenchmarkTreeMapPutAlreadyPresent(b *testing.B) {
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

func BenchmarkBuiltinMapPutAlreadyPresent(b *testing.B) {
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

func BenchmarkTreeMapBuild(b *testing.B) {
	for _, size := range sizes {
		keys := iterator.Collect(iterator.Count(size))

		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 1; i < b.N; i++ {
				b.StopTimer()
				xrand.Shuffle(keys)
				m := NewMap[int, int](xsort.OrderedLess[int])
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
		keys := iterator.Collect(iterator.Count(size))

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
