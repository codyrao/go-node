package main

import "C"

import (
	"encoding/json"
	"fmt"
	"time"
)

// FunctionA emits a short callback sequence to validate callback isolation for one exported function.
//
//export FunctionA
func FunctionA(params *C.char, callbackType *C.char) *C.char {
	// Parse the callback token first so the sample can skip callback work when no callback was supplied.
	callbackID, hasCallback := parseCallbackID(callbackType)

	if hasCallback {
		// Emit a deterministic callback sequence so JavaScript can verify routing to FunctionA only.
		for i := 1; i <= 3; i++ {
			time.Sleep(100 * time.Millisecond)

			callback := nodeCallback{id: callbackID}
			callback.send(map[string]interface{}{
				"function": "FunctionA",
				"message":  fmt.Sprintf("Callback %d from FunctionA", i),
				"index":    i,
			})
		}
	}

	// Return a completion payload so the JS test can assert the direct function result.
	resultData := map[string]interface{}{
		"status":   "success",
		"function": "FunctionA",
		"result":   "FunctionA started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

// FunctionB emits a second callback sequence with different timing to validate callback bookkeeping.
//
//export FunctionB
func FunctionB(params *C.char, callbackType *C.char) *C.char {
	// Parse the callback token first so the sample can skip callback work when no callback was supplied.
	callbackID, hasCallback := parseCallbackID(callbackType)

	if hasCallback {
		// Emit a deterministic callback sequence so JavaScript can verify routing to FunctionB only.
		for i := 1; i <= 3; i++ {
			time.Sleep(150 * time.Millisecond)

			callback := nodeCallback{id: callbackID}
			callback.send(map[string]interface{}{
				"function": "FunctionB",
				"message":  fmt.Sprintf("Callback %d from FunctionB", i),
				"index":    i,
			})
		}
	}

	// Return a completion payload so the JS test can assert the direct function result.
	resultData := map[string]interface{}{
		"status":   "success",
		"function": "FunctionB",
		"result":   "FunctionB started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

// FunctionC emits a third callback sequence to validate concurrent callback isolation across exports.
//
//export FunctionC
func FunctionC(params *C.char, callbackType *C.char) *C.char {
	// Parse the callback token first so the sample can skip callback work when no callback was supplied.
	callbackID, hasCallback := parseCallbackID(callbackType)

	if hasCallback {
		// Emit a deterministic callback sequence so JavaScript can verify routing to FunctionC only.
		for i := 1; i <= 3; i++ {
			time.Sleep(200 * time.Millisecond)

			callback := nodeCallback{id: callbackID}
			callback.send(map[string]interface{}{
				"function": "FunctionC",
				"message":  fmt.Sprintf("Callback %d from FunctionC", i),
				"index":    i,
			})
		}
	}

	// Return a completion payload so the JS test can assert the direct function result.
	resultData := map[string]interface{}{
		"status":   "success",
		"function": "FunctionC",
		"result":   "FunctionC started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}
