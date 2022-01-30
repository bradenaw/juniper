# `package xtime`

```
import "github.com/bradenaw/juniper/xtime"
```

# Overview



# Index

<samp><a href="#SleepContext">func SleepContext(ctx context.Context, d time.Duration) error</a></samp>
<samp><a href="#DeadlineTooSoonError">type DeadlineTooSoonError</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Error">func (err DeadlineTooSoonError) Error() string</a></samp>


# Constants

This section is empty.

# Variables

This section is empty.

# Functions

<h2><a id="SleepContext"></a><samp>func <a href="#SleepContext">SleepContext</a>(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, d <a href="https://pkg.go.dev/time#Duration">time.Duration</a>) error</samp></h2>

SleepContext pauses the current goroutine for at least the duration d and returns nil, unless ctx
expires in the mean time in which case it returns ctx.Err().

A negative or zero duration causes SleepContext to return nil immediately.

If there is less than d left until ctx's deadline, returns DeadlineTooSoonError immediately.


# Types

<h2><a id="DeadlineTooSoonError"></a><samp>type DeadlineTooSoonError</samp></h2>
```go
type DeadlineTooSoonError struct {
	// contains filtered or unexported fields
}
```



<h2><a id="Error"></a><samp>func (err <a href="#DeadlineTooSoonError">DeadlineTooSoonError</a>) Error() string</samp></h2>



