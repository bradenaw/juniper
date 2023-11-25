// Package xmath contains extensions to the standard library package math.
package xmath

// Abs returns the absolute value of x. It panics if this value is not representable, for example
// because -math.MinInt32 requires more than 32 bits to represent and so does not fit in an int32.
func Abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64](x T) T {
	if x < 0 {
		if -x == x {
			panic("can't xmath.Abs minimum value: positive equivalent not representable")
		}
		return -x
	}
	return x
}
