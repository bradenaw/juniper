# `package xatomic`

```
import "github.com/bradenaw/juniper/xsync/xatomic"
```

## Overview

Package xatomic contains extensions to the standard library package sync/atomic.


## Index

<samp><a href="#Value">type Value</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#CompareAndSwap">func (v *Value[T]) CompareAndSwap(old, new T) (swapped bool)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Load">func (v *Value[T]) Load() T</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Store">func (v *Value[T]) Store(t T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Swap">func (v *Value[T]) Swap(new T) (old T)</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

This section is empty.

## Types

<h3><a id="Value"></a><samp>type Value</samp></h3>
```go
type Value[T any] struct {
	// contains filtered or unexported fields
}
```

Value is equivalent to sync/atomic.Value, except strongly typed.


<h3><a id="CompareAndSwap"></a><samp>func (v *<a href="#Value">Value</a>[T]) CompareAndSwap(old, new T) swapped bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xsync/xatomic/xatomic.go#L13">src</a></small></sub></h3>



<h3><a id="Load"></a><samp>func (v *<a href="#Value">Value</a>[T]) Load() T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xsync/xatomic/xatomic.go#L17">src</a></small></sub></h3>



<h3><a id="Store"></a><samp>func (v *<a href="#Value">Value</a>[T]) Store(t T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xsync/xatomic/xatomic.go#L21">src</a></small></sub></h3>



<h3><a id="Swap"></a><samp>func (v *<a href="#Value">Value</a>[T]) Swap(new T) old T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xsync/xatomic/xatomic.go#L25">src</a></small></sub></h3>



