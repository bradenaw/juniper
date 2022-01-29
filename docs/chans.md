# chans
--
    import "."


## Usage

#### func  RecvContext

```go
func RecvContext[T any](ctx context.Context, c <-chan T) (T, bool, error)
```
RecvContext attempts to receive from channel c. If c is closed before or during,
returns (_, false, nil). If ctx expires before or during, returns (_, _,
ctx.Err()).

#### func  SendContext

```go
func SendContext[T any](ctx context.Context, c chan<- T, item T) error
```
SendContext sends item on channel c and returns nil, unless ctx expires in which
case it returns ctx.Err().
