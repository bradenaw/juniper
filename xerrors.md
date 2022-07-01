# `package xerrors`

```
import "github.com/bradenaw/juniper/xerrors"
```

## Overview

Package xerrors contains extensions to the standard library package errors.


## Index

<samp><a href="#WithStack">func WithStack(err error) error</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="WithStack"></a><samp>func <a href="#WithStack">WithStack</a>(err error) error</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xerrors/xerrors.go#L44">src</a></small></sub></h3>

WithStack returns an error that wraps err and adds the call stack of the call to WithStack to
Error(). If err is nil or already has a stack attached, returns err.


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

This section is empty.

