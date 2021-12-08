package slices

func Grow[T any](x []T, n int) []T {
	if cap(x) - len(x) < n {
		x2 := make([]T, len(x) + n)
		copy(x2, x)
		return x2[:len(x)]
	}
	return x
}

func Filter[T any](x []T, keep func(t T) bool) []T {
	filtered := x[:0]
	for i := range x {
		if keep(x[i]) {
			filtered = append(filtered, x[i])
		}
	}
	// Zero out the rest in case they contain pointers, so that filtered doesn't retain references.
	var zero T
	for i := range x[len(filtered):] {
		x[i] = zero
	}
	return filtered
}

func Reverse[T any](x []T) {
	for i := 0; i < len(x) / 2 - 1; i++ {
		x[i], x[len(x)-i-1] = x[len(x)-i-1], x[i]
	}
}

func Remove[T any](x []T, i int) []T {
	var zero T
	copy(x[i:], x[i+1:])
	x[len(x)-1] = zero
	return x[:len(x)-1]
}
