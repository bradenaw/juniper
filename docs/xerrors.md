# `package xerrors`

```
import "github.com/bradenaw/juniper/xerrors"
```

## Overview



## Index

<samp><a href="#WithStack">func WithStack(err error) error</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="WithStack"></a><samp>func <a href="#WithStack">WithStack</a>(err error) error</samp></h3>

WithStack returns an error that wraps err and adds the call stack of the call to WithStack to
Error().


#### Example 
```go
{
	err := WithStack(errors.New("foo"))

	fmt.Println(strings.Join(strings.Split(err.Error(), "\n")[:3], "\n"))

}
```

Output:
```text

foo

github.com/bradenaw/juniper/xerrors.ExampleWithStack(...)
```
## Types

