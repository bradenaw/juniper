# `package xatomic`

```
import "github.com/bradenaw/juniper/xsync/xatomic"
```

# Overview

Package xatomic contains extensions to the standard library package sync/atomic.


# Index

<pre><a href="#Value">type Value</a></pre>
<pre>    <a href="#CompareAndSwap">func (v *Value[T]) CompareAndSwap(old, new T) (swapped bool)</a></pre>
<pre>    <a href="#Load">func (v *Value[T]) Load() T</a></pre>
<pre>    <a href="#Store">func (v *Value[T]) Store(t T)</a></pre>
<pre>    <a href="#Swap">func (v *Value[T]) Swap(new T) (old T)</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

# Types

## <a id="Value"></a><pre>type Value</pre>
```go
type Value[T any] struct {
	// contains filtered or unexported fields
}
```

Value is equivalent to sync/atomic.Value, except strongly typed.


## <a id="CompareAndSwap"></a><pre>func (v *<a href="#Value">Value</a>[T]) CompareAndSwap(old, new T) swapped bool</pre>



## <a id="Load"></a><pre>func (v *<a href="#Value">Value</a>[T]) Load() T</pre>



## <a id="Store"></a><pre>func (v *<a href="#Value">Value</a>[T]) Store(t T)</pre>



## <a id="Swap"></a><pre>func (v *<a href="#Value">Value</a>[T]) Swap(new T) old T</pre>



