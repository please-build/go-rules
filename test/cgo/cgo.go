package cgo

// #include "hello.h"
import "C"

func Hello() string {
	return C.GoString(C.hello())
}
