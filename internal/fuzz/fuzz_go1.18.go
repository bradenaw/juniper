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
	takeByte := func() (byte, bool) {
		if len(b) < 1 {
			return 0, false
		}
		x := b[0]
		b = b[1:]
		return x, true
	}
	takeInt := func() (int, bool) {
		if len(b) < 8 {
			return 0, false
		}
		x := int(binary.BigEndian.Uint64(b[:8]))
		b = b[8:]
		return x, true
	}
	takeBool := func() (bool, bool) {
		b, ok := takeByte()
		return b != 0, ok
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
			case reflect.Uint8:
				x, ok := takeByte()
				if !ok {
					break Loop
				}
				args[j] = reflect.ValueOf(x)
			case reflect.Bool:
				x, ok := takeBool()
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
