# xmath
--
    import "."

Package xmath contains extensions to the standard library package math.

## Usage

#### func  Abs

```go
func Abs[T constraints.Signed](x T) T
```
Abs returns the absolute value of x. It panics if this value is not
representable, for example because -math.MinInt32 requires more than 32 bits to
represent and so does not fit in an int32.

#### func  Clamp

```go
func Clamp[T constraints.Ordered](x, min, max T) T
```
Clamp clamps the value of x to within min and max.

#### func  Max

```go
func Max[T constraints.Ordered](a, b T) T
```
Max returns the maximum of a and b based on the > operator.

#### func  Min

```go
func Min[T constraints.Ordered](a, b T) T
```
Min returns the minimum of a and b based on the < operator.
