# maps
--
    import "."


## Usage

#### func  Keys

```go
func Keys[K comparable, V any](m map[K]V) []K
```
Keys returns the keys of m as a slice.

#### func  Values

```go
func Values[K comparable, V any](m map[K]V) []V
```
Values returns the values of m as a slice.

#### type Set

```go
type Set[T comparable] map[T]struct{}
```

Set implements sets.Set for map[T]struct{}.

#### func (BADRECV) Add

```go
func (s Set[T]) Add(item T)
```

#### func (BADRECV) Contains

```go
func (s Set[T]) Contains(item T) bool
```

#### func (BADRECV) Iterate

```go
func (s Set[T]) Iterate() iterator.Iterator[T]
```

#### func (BADRECV) Len

```go
func (s Set[T]) Len() int
```

#### func (BADRECV) Remove

```go
func (s Set[T]) Remove(item T)
```
