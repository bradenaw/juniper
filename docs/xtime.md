# `package xtime`

```
import "github.com/bradenaw/juniper/xtime"
```

# Overview



# Index

<pre><a href="#SleepContext">func SleepContext(ctx context.Context, d time.Duration) error</a></pre>
<pre><a href="#DeadlineTooSoonError">type DeadlineTooSoonError</a></pre>
<pre>    <a href="#Error">func (err DeadlineTooSoonError) Error() string</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

<h2><a id="SleepContext"></a><pre>func <a href="#SleepContext">SleepContext</a>(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, d <a href="https://pkg.go.dev/time#Duration">time.Duration</a>) error</pre></h2>

SleepContext pauses the current goroutine for at least the duration d and returns nil, unless ctx
expires in the mean time in which case it returns ctx.Err().

A negative or zero duration causes SleepContext to return nil immediately.

If there is less than d left until ctx's deadline, returns DeadlineTooSoonError immediately.


# Types

## <a id="DeadlineTooSoonError"></a><pre>type DeadlineTooSoonError</pre>
```go
type DeadlineTooSoonError struct {
	// contains filtered or unexported fields
}
```



<h2><a id="Error"></a><pre>func (err <a href="#DeadlineTooSoonError">DeadlineTooSoonError</a>) Error() string</pre></h2>



