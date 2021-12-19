package fuzz

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

func Operations(b []byte, check func(), fns ...interface{}) {
	choose := func(n int) (int, bool) {
		if len(b) < 1 {
			return 0, false
		}
		if n > 255 {
			panic("")
		}
		choice := int(b[0]) % n
		b = b[1:]
		return choice, true
	}
	takeInt := func() (int, bool) {
		if len(b) < 8 {
			return 0, false
		}
		x := int(binary.BigEndian.Uint64(b[:8]))
		b = b[8:]
		return x, true
	}

Loop:
	for {
		check()
		i, ok := choose(len(fns))
		if !ok {
			break
		}
		fnV := reflect.ValueOf(fns[i])

		args := make([]reflect.Value, fnV.Type().NumIn())
		for j := range args {
			argType := fnV.Type().In(j)

			switch argType.Kind() {
			case reflect.Int:
				x, ok := takeInt()
				if !ok {
					break Loop
				}
				args[j] = reflect.ValueOf(x)
			default:
				panic(fmt.Sprintf("arg type %s not supported", argType.Kind()))
			}
		}
		fnV.Call(args)
	}
}
