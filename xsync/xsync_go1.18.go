//go:build go1.18

package xsync

import "sync"

// Lazy makes a lazily-initialized value. On first access, it uses f to create the value. Later
// accesses all receive the same value.
func Lazy[T any](f func() T) func() T {
	var once sync.Once
	var val T
	return func() T {
		once.Do(func() {
			val = f()
		})
		return val
	}
}
