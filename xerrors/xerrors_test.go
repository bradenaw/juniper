package xerrors

import (
	"errors"
	"fmt"
	"strings"
)

func ExampleWithStack() {
	err := WithStack(errors.New("foo"))

	fmt.Println(strings.Join(strings.Split(err.Error(), "\n")[:3], "\n"))

	// Output:
	//
	// foo
	//
	// github.com/bradenaw/juniper/xerrors.ExampleWithStack(...)
}
