# `package xmath`

```
import "github.com/bradenaw/juniper/xmath"
```

## Overview

Package xmath contains extensions to the standard library package math.


## Index

<samp><a href="#Abs">func Abs[T constraints.Signed](x T) T</a></samp>

<samp><a href="#Clamp">func Clamp[T constraints.Ordered](x, min, max T) T</a></samp>

<samp><a href="#Max">func Max[T constraints.Ordered](a, b T) T</a></samp>

<samp><a href="#Min">func Min[T constraints.Ordered](a, b T) T</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="Abs"></a><samp>func <a href="#Abs">Abs</a>[T <a href="https://pkg.go.dev/golang.org/x/exp/constraints#Signed">constraints.Signed</a>](x T) T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xmath.go#L24">src</a></small></sub></h3>

Abs returns the absolute value of x. It panics if this value is not representable, for example
because -math.MinInt32 requires more than 32 bits to represent and so does not fit in an int32.


<h3><a id="Clamp"></a><samp>func <a href="#Clamp">Clamp</a>[T <a href="https://pkg.go.dev/golang.org/x/exp/constraints#Ordered">constraints.Ordered</a>](x, min, max T) T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xmath.go#L35">src</a></small></sub></h3>

Clamp clamps the value of x to within min and max.


<h3><a id="Max"></a><samp>func <a href="#Max">Max</a>[T <a href="https://pkg.go.dev/golang.org/x/exp/constraints#Ordered">constraints.Ordered</a>](a, b T) T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xmath.go#L15">src</a></small></sub></h3>

Max returns the maximum of a and b based on the > operator.


<h3><a id="Min"></a><samp>func <a href="#Min">Min</a>[T <a href="https://pkg.go.dev/golang.org/x/exp/constraints#Ordered">constraints.Ordered</a>](a, b T) T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xmath.go#L7">src</a></small></sub></h3>

Min returns the minimum of a and b based on the < operator.


## Types

This section is empty.

