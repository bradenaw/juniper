# `package xerrors`

```
import "github.com/bradenaw/juniper/xerrors"
```

# Overview



# Index

<pre><a href="#WithStack">func WithStack(err error) error</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

## <a id="WithStack"></a><pre>func <a href="#WithStack">WithStack</a>(err error) error</pre>

WithStack returns an error that wraps err and adds the call stack of the call to WithStack to
Error().


### Example 
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

# Types

