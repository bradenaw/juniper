# xatomic
--
    import "."

Package xatomic contains extensions to the standard library package sync/atomic.

## Usage

#### type Value

```go
type Value[T any] struct {
}
```

Value is equivalent to sync/atomic.Value, except strongly typed.

#### func (*BADRECV) CompareAndSwap

```go
func (v *Value[T]) CompareAndSwap(old, new T) (swapped bool)
```

#### func (*BADRECV) Load

```go
func (v *Value[T]) Load() T
```

#### func (*BADRECV) Store

```go
func (v *Value[T]) Store(t T)
```

#### func (*BADRECV) Swap

```go
func (v *Value[T]) Swap(new T) (old T)
```
