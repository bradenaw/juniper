# `package xmath`

```
import "github.com/bradenaw/juniper/xmath"
```

# Overview

Package xmath contains extensions to the standard library package math.


# Index

<pre><a href="#Abs">func Abs[T constraints.Signed](x T) T</a></pre>
<pre><a href="#Clamp">func Clamp[T constraints.Ordered](x, min, max T) T</a></pre>
<pre><a href="#Max">func Max[T constraints.Ordered](a, b T) T</a></pre>
<pre><a href="#Min">func Min[T constraints.Ordered](a, b T) T</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

## <a id="Abs"></a><pre>func <a href="#Abs">Abs</a>[T <a href="https://pkg.go.dev/constraints#Signed">constraints.Signed</a>](x T) T</pre>

Abs returns the absolute value of x. It panics if this value is not representable, for example
because -math.MinInt32 requires more than 32 bits to represent and so does not fit in an int32.


## <a id="Clamp"></a><pre>func <a href="#Clamp">Clamp</a>[T <a href="https://pkg.go.dev/constraints#Ordered">constraints.Ordered</a>](x, min, max T) T</pre>

Clamp clamps the value of x to within min and max.


## <a id="Max"></a><pre>func <a href="#Max">Max</a>[T <a href="https://pkg.go.dev/constraints#Ordered">constraints.Ordered</a>](a, b T) T</pre>

Max returns the maximum of a and b based on the > operator.


## <a id="Min"></a><pre>func <a href="#Min">Min</a>[T <a href="https://pkg.go.dev/constraints#Ordered">constraints.Ordered</a>](a, b T) T</pre>

Min returns the minimum of a and b based on the < operator.


# Types

