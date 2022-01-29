# xerrors
--
    import "."


## Usage

#### func  WithStack

```go
func WithStack(err error) error
```
WithStack returns an error that wraps err and adds the call stack of the call to
WithStack to Error().
