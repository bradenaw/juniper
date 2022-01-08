//go:build go1.18

package tree

import (
	"fmt"
	"testing"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath/xrand"
	"github.com/bradenaw/juniper/xsort"
)

var sizes = []int{10, 100, 1_000, 10_000, 100_000, 1_000_000}

// Run under WSL since I don't have a native Linux machine handy at the moment.
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
//
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz
//
//
//
// tree.Map is slower. It's within about a factor of 2 for Get(). It's much slower to write to, and
// the gap is (expectedly) larger for larger maps. It only makes any sense to use if you need
// ordered iteration while the map is also evolving. If you only need ordered iteration
// occasionally compared to how often the map changes, you might still be better off using the
// builtin map and sorting the keys when necessary.
//
// ("old" means builtin, and "new" means tree.Map)
//
//   name                                  old time/op    new time/op    delta
//   MapGet/Size=10-16                       10.7ns ± 0%    12.4ns ± 1%    +16.00%  (p=0.008 n=5+5)
//   MapGet/Size=100-16                      11.2ns ± 2%    18.7ns ± 2%    +67.01%  (p=0.008 n=5+5)
//   MapGet/Size=1000-16                     16.2ns ± 2%    41.3ns ± 1%   +154.99%  (p=0.008 n=5+5)
//   MapGet/Size=10000-16                    22.9ns ± 2%    50.7ns ± 1%   +122.01%  (p=0.008 n=5+5)
//   MapGet/Size=100000-16                   28.1ns ± 2%    58.6ns ± 1%   +108.66%  (p=0.008 n=5+5)
//   MapGet/Size=1000000-16                  51.5ns ± 1%    71.7ns ± 2%    +39.28%  (p=0.008 n=5+5)
//   MapPut/Size=10-16                       42.6ns ± 2%    50.7ns ± 0%    +18.79%  (p=0.008 n=5+5)
//   MapPut/Size=100-16                      55.8ns ± 2%    79.6ns ± 1%    +42.63%  (p=0.008 n=5+5)
//   MapPut/Size=1000-16                     65.0ns ± 2%   132.0ns ± 3%   +103.21%  (p=0.008 n=5+5)
//   MapPut/Size=10000-16                    60.0ns ± 2%   180.2ns ± 4%   +200.35%  (p=0.008 n=5+5)
//   MapPut/Size=100000-16                   61.9ns ± 3%   274.4ns ± 5%   +343.41%  (p=0.008 n=5+5)
//   MapPut/Size=1000000-16                   112ns ± 2%     541ns ± 0%   +381.02%  (p=0.016 n=5+4)
//   MapPutAlreadyPresent/Size=10-16         13.7ns ± 6%    12.3ns ± 4%    -10.62%  (p=0.008 n=5+5)
//   MapPutAlreadyPresent/Size=100-16        15.2ns ± 4%    19.5ns ± 1%    +27.86%  (p=0.008 n=5+5)
//   MapPutAlreadyPresent/Size=1000-16       20.3ns ± 2%    54.4ns ± 1%   +168.08%  (p=0.008 n=5+5)
//   MapPutAlreadyPresent/Size=10000-16      25.8ns ± 2%    88.0ns ± 2%   +241.53%  (p=0.008 n=5+5)
//   MapPutAlreadyPresent/Size=100000-16     32.1ns ± 0%   144.6ns ± 1%   +350.76%  (p=0.008 n=5+5)
//   MapPutAlreadyPresent/Size=1000000-16    62.4ns ± 1%   382.9ns ±10%   +514.06%  (p=0.008 n=5+5)
//   MapIterate/Size=10-16                   79.2ns ± 1%    98.6ns ± 1%    +24.49%  (p=0.008 n=5+5)
//   MapIterate/Size=100-16                   703ns ± 1%     697ns ± 1%     -0.94%  (p=0.032 n=5+5)
//   MapIterate/Size=1000-16                 8.66µs ± 1%    8.65µs ± 1%       ~     (p=1.000 n=5+5)
//   MapIterate/Size=10000-16                77.5µs ± 1%   147.3µs ± 0%    +90.15%  (p=0.008 n=5+5)
//   MapIterate/Size=100000-16                740µs ± 3%    1844µs ± 0%   +149.23%  (p=0.008 n=5+5)
//   MapIterate/Size=1000000-16              9.20ms ± 1%   82.40ms ± 1%   +795.51%  (p=0.008 n=5+5)
//   MapBuild/Size=10-16                     1.36µs ± 3%    1.84µs ± 3%    +34.77%  (p=0.008 n=5+5)
//   MapBuild/Size=100-16                    7.46µs ± 3%   11.56µs ± 2%    +55.10%  (p=0.008 n=5+5)
//   MapBuild/Size=1000-16                   58.5µs ± 4%   126.2µs ± 1%   +115.80%  (p=0.008 n=5+5)
//   MapBuild/Size=10000-16                   519µs ± 4%    1715µs ± 3%   +230.28%  (p=0.008 n=5+5)
//   MapBuild/Size=100000-16                 5.20ms ± 4%   25.87ms ± 1%   +397.17%  (p=0.008 n=5+5)
//   MapBuild/Size=1000000-16                98.1ms ± 2%   545.1ms ± 1%   +455.71%  (p=0.016 n=5+4)
//
//
// tree.Map appears to use a little less memory than the builtin map. Your mileage may vary, since
// this is only tested with int keys and values and for this set of sizes.
//
//   name                                  old alloc/op   new alloc/op   delta
//   MapPut/Size=10-16                        48.0B ± 0%     51.0B ± 0%     +6.25%  (p=0.008 n=5+5)
//   MapPut/Size=100-16                       55.0B ± 0%     48.0B ± 0%    -12.73%  (p=0.008 n=5+5)
//   MapPut/Size=1000-16                      86.0B ± 0%     48.0B ± 0%    -44.19%  (p=0.008 n=5+5)
//   MapPut/Size=10000-16                     68.0B ± 0%     48.0B ± 0%    -29.41%  (p=0.008 n=5+5)
//   MapPut/Size=100000-16                    57.0B ± 0%     48.0B ± 0%    -15.79%  (p=0.008 n=5+5)
//   MapPut/Size=1000000-16                   86.8B ± 2%     48.0B ± 0%    -44.70%  (p=0.008 n=5+5)
//   MapBuild/Size=10-16                       292B ± 0%      480B ± 0%    +64.25%  (p=0.008 n=5+5)
//   MapBuild/Size=100-16                    5.36kB ± 0%    4.80kB ± 0%    -10.54%  (p=0.008 n=5+5)
//   MapBuild/Size=1000-16                   86.5kB ± 0%    48.0kB ± 0%    -44.55%  (p=0.008 n=5+5)
//   MapBuild/Size=10000-16                   687kB ± 0%     479kB ± 0%    -30.21%  (p=0.008 n=5+5)
//   MapBuild/Size=100000-16                 5.74MB ± 0%    4.75MB ± 0%    -17.19%  (p=0.008 n=5+5)
//   MapBuild/Size=1000000-16                86.7MB ± 0%    47.5MB ± 0%    -45.22%  (p=0.008 n=5+5)
//
//
// Although tree.Map requires a lot more allocations, since the builtin map allocates buckets of a
// few elements at a time and tree.Map allocates once per element. This may mean that it's more
// costly for the GC to keep track of.
//
//   name                                  old allocs/op  new allocs/op  delta
//   MapBuild/Size=10-16                       1.00 ± 0%      9.00 ± 0%   +800.00%  (p=0.008 n=5+5)
//   MapBuild/Size=100-16                      16.0 ± 0%      99.0 ± 0%   +518.75%  (p=0.008 n=5+5)
//   MapBuild/Size=1000-16                     64.0 ± 0%     999.0 ± 0%  +1460.94%  (p=0.008 n=5+5)
//   MapBuild/Size=10000-16                     276 ± 0%      9985 ± 0%  +3523.00%  (p=0.008 n=5+5)
//   MapBuild/Size=100000-16                  4.00k ± 0%    99.00k ± 0%  +2377.73%  (p=0.008 n=5+5)
//   MapBuild/Size=1000000-16                 38.0k ± 0%    990.0k ± 0%  +2503.74%  (p=0.008 n=5+5)

func BenchmarkTreeMapGet(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int, xsort.NaturalOrder[int]]()
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
				_ = m[i%size]
			}
		})
	}
}

func BenchmarkTreeMapPut(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int, xsort.NaturalOrder[int]]()
		keys := iterator.Collect(iterator.Count(size))
		xrand.Shuffle(keys)

		b.Run(fmt.Sprintf("Size=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m.Put(keys[i%size], i)
				if m.Len() == size {
					m = NewMap[int, int, xsort.NaturalOrder[int]]()
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
				m[keys[i%size]] = i
				if len(m) == size {
					m = make(map[int]int)
				}
			}
		})
	}
}

func BenchmarkTreeMapPutAlreadyPresent(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int, xsort.NaturalOrder[int]]()
		keys := iterator.Collect(iterator.Count(size))
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
				m[keys[i%size]] = i
			}
		})
	}
}

func BenchmarkTreeMapIterate(b *testing.B) {
	for _, size := range sizes {
		m := NewMap[int, int, xsort.NaturalOrder[int]]()
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
				m := NewMap[int, int, xsort.NaturalOrder[int]]()
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
