package utils

import "unsafe"

// IsNilInterface check that interface is really nil
func IsNilInterface(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}
