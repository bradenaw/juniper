# sets
--
    import "."

package sets contains set operations like union, intersection, and difference.

## Usage

#### type Set

```go
type Set[T any] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Len() int
	Iterate() iterator.Iterator[T]
}
```


#### func  Difference

```go
func Difference[T comparable](out, a, b Set[T]) Set[T]
```
Difference adds to out all items that appear in a but not in b and returns out.

#### func  Intersection

```go
func Intersection[T comparable](out Set[T], sets ...Set[T]) Set[T]
```
Intersection adds to out all items that appear in all sets and returns out.

#### func  Union

```go
func Union[T any](out Set[T], sets ...Set[T]) Set[T]
```
Union adds to out out all items from sets and returns out.
