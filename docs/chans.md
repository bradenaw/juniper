# `package chans`

```
import "github.com/bradenaw/juniper/chans"
```

# Overview



# Index

<samp><a href="#RecvContext">func RecvContext[T any](ctx context.Context, c &lt;-chan T) (T, bool, error)</a></samp>

<samp><a href="#SendContext">func SendContext[T any](ctx context.Context, c chan&lt;- T, item T) error</a></samp>


# Constants

This section is empty.

# Variables

This section is empty.

# Functions

<h2><a id="RecvContext"></a><samp>func <a href="#RecvContext">RecvContext</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, c &lt;-chan T) (T, bool, error)</samp></h2>

RecvContext attempts to receive from channel c. If c is closed before or during, returns (_,
false, nil). If ctx expires before or during, returns (_, _, ctx.Err()).


<h2><a id="SendContext"></a><samp>func <a href="#SendContext">SendContext</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, c chan&lt;- T, item T) error</samp></h2>

SendContext sends item on channel c and returns nil, unless ctx expires in which case it returns
ctx.Err().


# Types

