//go:build go1.18

package xatomic

import (
	"sync/atomic"
)

type Value[T any] struct {
	v atomic.Value
}

func (v *Value[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.v.CompareAndSwap(old, new)
}

func (v *Value[T]) Load() T {
	return v.v.Load().(T)
}

func (v *Value[T]) Store(t T) {
	v.v.Store(t)
}

func (v *Value[T]) Swap(new T) (old T) {
	var zero T
	out := v.v.Swap(new)
	if out == nil {
		return zero
	}
	return out.(T)
}
