# xtime
--
    import "."


## Usage

#### func  SleepContext

```go
func SleepContext(ctx context.Context, d time.Duration) error
```
SleepContext pauses the current goroutine for at least the duration d and
returns nil, unless ctx expires in the mean time in which case it returns
ctx.Err().

A negative or zero duration causes SleepContext to return nil immediately.

If there is less than d left until ctx's deadline, returns DeadlineTooSoonError
immediately.

#### type DeadlineTooSoonError

```go
type DeadlineTooSoonError struct {
}
```


#### func (DeadlineTooSoonError) Error

```go
func (err DeadlineTooSoonError) Error() string
```
