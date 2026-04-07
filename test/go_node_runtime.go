package main

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// FreeCString releases a C string allocated inside the Go shared library.
//
//export FreeCString
func FreeCString(ptr *C.char) {
	// Guard against nil pointers to keep release operation idempotent.
	if ptr == nil {
		return
	}

	// Free inside the same DLL allocator boundary to avoid cross-runtime deallocation issues.
	C.free(unsafe.Pointer(ptr))
}
