package main

/*
#cgo CFLAGS: -I.
#include "callback.h"
*/
import "C"

import (
	"encoding/json"
	"strconv"
	"unsafe"
)

var gCallNodeCallback uintptr

type nodeCallback struct {
	id    int32
	token string
}

// RegisterGoCallback stores the wrapper callback entrypoint so sample exports can send data back into Node.
//
//export RegisterGoCallback
func RegisterGoCallback(fn uintptr) {
	// Capture the callback entrypoint once so all sample exports share the same bridge state.
	gCallNodeCallback = fn
}

// newNodeCallback converts the callback token string received from Node into a reusable callback descriptor.
func newNodeCallback(raw *C.char) nodeCallback {
	// Reject missing runtime state first so callers can no-op safely when no callback was provided.
	if gCallNodeCallback == 0 || raw == nil {
		return nodeCallback{id: -1}
	}

	// Parse the callback token into the integer identifier expected by the native bridge.
	token := C.GoString(raw)
	callbackID, err := strconv.Atoi(token)
	if err != nil {
		return nodeCallback{id: -1, token: token}
	}

	return nodeCallback{
		id:    int32(callbackID),
		token: token,
	}
}

// parseCallbackID exposes the callback identifier form used by the callback-fix sample exports.
func parseCallbackID(raw *C.char) (int32, bool) {
	// Reuse shared token parsing so callback validation stays consistent across samples.
	callback := newNodeCallback(raw)
	if callback.id < 0 {
		return 0, false
	}

	return callback.id, true
}

// send serializes callback payload data and forwards it through the Node bridge.
func (cb nodeCallback) send(data map[string]interface{}) {
	// Skip work when no callback bridge is active or the callback token was invalid.
	if gCallNodeCallback == 0 || cb.id < 0 {
		return
	}

	// Encode the payload to JSON because the wrapper callback protocol transports strings.
	jsonData, _ := json.Marshal(data)
	cJSON := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(cJSON))

	// Invoke the native callback bridge with the parsed callback identifier and payload bytes.
	C.callCallbackWithId(unsafe.Pointer(gCallNodeCallback), C.int32_t(cb.id), cJSON)
}

// main keeps the sample package buildable as a c-shared package while exposing no standalone CLI behavior.
func main() {
	// Leave main empty because these files are loaded through the generated DLL entrypoints.
}
