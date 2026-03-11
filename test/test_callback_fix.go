package main

/*
#cgo CFLAGS: -I.
#include "callback.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

var gCallNodeCallback uintptr

//export RegisterGoCallback
func RegisterGoCallback(fn uintptr) {
	gCallNodeCallback = fn
}

func parseCallbackID(raw *C.char) (int32, bool) {
	if gCallNodeCallback == 0 || raw == nil {
		return 0, false
	}

	callbackID, err := strconv.Atoi(C.GoString(raw))
	if err != nil {
		return 0, false
	}

	return int32(callbackID), true
}

//export FunctionA
func FunctionA(params *C.char, callbackType *C.char) *C.char {
	callbackID, hasCallback := parseCallbackID(callbackType)

	if hasCallback {
		for i := 1; i <= 3; i++ {
			time.Sleep(100 * time.Millisecond)

			callbackData := map[string]interface{}{
				"function": "FunctionA",
				"message":  fmt.Sprintf("Callback %d from FunctionA", i),
				"index":    i,
			}
			jsonData, _ := json.Marshal(callbackData)

			cJSON := C.CString(string(jsonData))
			C.callCallbackWithId(unsafe.Pointer(gCallNodeCallback), C.int32_t(callbackID), cJSON)
			C.free(unsafe.Pointer(cJSON))
		}
	}

	resultData := map[string]interface{}{
		"status":   "success",
		"function": "FunctionA",
		"result":   "FunctionA started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export FunctionB
func FunctionB(params *C.char, callbackType *C.char) *C.char {
	callbackID, hasCallback := parseCallbackID(callbackType)

	if hasCallback {
		for i := 1; i <= 3; i++ {
			time.Sleep(150 * time.Millisecond)

			callbackData := map[string]interface{}{
				"function": "FunctionB",
				"message":  fmt.Sprintf("Callback %d from FunctionB", i),
				"index":    i,
			}
			jsonData, _ := json.Marshal(callbackData)

			cJSON := C.CString(string(jsonData))
			C.callCallbackWithId(unsafe.Pointer(gCallNodeCallback), C.int32_t(callbackID), cJSON)
			C.free(unsafe.Pointer(cJSON))
		}
	}

	resultData := map[string]interface{}{
		"status":   "success",
		"function": "FunctionB",
		"result":   "FunctionB started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export FunctionC
func FunctionC(params *C.char, callbackType *C.char) *C.char {
	callbackID, hasCallback := parseCallbackID(callbackType)

	if hasCallback {
		for i := 1; i <= 3; i++ {
			time.Sleep(200 * time.Millisecond)

			callbackData := map[string]interface{}{
				"function": "FunctionC",
				"message":  fmt.Sprintf("Callback %d from FunctionC", i),
				"index":    i,
			}
			jsonData, _ := json.Marshal(callbackData)

			cJSON := C.CString(string(jsonData))
			C.callCallbackWithId(unsafe.Pointer(gCallNodeCallback), C.int32_t(callbackID), cJSON)
			C.free(unsafe.Pointer(cJSON))
		}
	}

	resultData := map[string]interface{}{
		"status":   "success",
		"function": "FunctionC",
		"result":   "FunctionC started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

func main() {}
