package main

import "C"

import (
	"encoding/json"
	"fmt"
	"time"
)

// Hello1 doubles the provided value and echoes the request fields for basic wrapper verification.
//
//export Hello1
func Hello1(params *C.char, callbackType *C.char) *C.char {
	// Decode the input payload so the sample can validate JSON object argument handling.
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	// Extract typed fields with permissive defaults because sample inputs come from dynamic JavaScript values.
	name := ""
	if n, ok := inputData["name"].(string); ok {
		name = n
	}

	value := 0
	if v, ok := inputData["value"].(float64); ok {
		value = int(v)
	}

	// Build the response payload that demonstrates successful synchronous execution.
	result := value * 2

	resultData := map[string]interface{}{
		"name":   name,
		"value":  value,
		"result": result,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

// HelloWithCallback emits three synchronous callbacks before returning a final success payload.
//
//export HelloWithCallback
func HelloWithCallback(params *C.char, callbackType *C.char) *C.char {
	// Decode input first so callback payloads can echo a caller-provided test marker.
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	// Resolve the callback token into the shared callback runtime representation.
	callback := newNodeCallback(callbackType)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

	// Send a small deterministic callback sequence so the JavaScript test can assert ordering and payload shape.
	for i := 1; i <= 3; i++ {
		time.Sleep(300 * time.Millisecond)

		callback.send(map[string]interface{}{
			"callbackType": callback.token,
			"test":         testMsg,
			"result":       fmt.Sprintf("Callback %d", i),
		})
	}

	// Return a final JSON payload after the callback sequence completes.
	resultData := map[string]interface{}{
		"callbackType": callback.token,
		"status":       "success",
		"result":       42,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

// AsyncHello starts a background callback sequence and returns immediately to exercise async callback handling.
//
//export AsyncHello
func AsyncHello(params *C.char, callbackType *C.char) *C.char {
	// Decode input so the asynchronous payloads can include the caller-provided marker string.
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	// Resolve the callback token before starting the goroutine so later callback sends use stable data.
	callback := newNodeCallback(callbackType)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

	// Emit callbacks from a goroutine to verify the wrapper's async callback queueing behavior.
	go func(cb nodeCallback, msg string) {
		for i := 1; i <= 5; i++ {
			time.Sleep(500 * time.Millisecond)

			cb.send(map[string]interface{}{
				"callbackType": cb.token,
				"test":         msg,
				"result":       fmt.Sprintf("Async callback %d", i),
			})
		}
	}(callback, testMsg)

	// Return immediately so JavaScript can observe that callbacks continue after the function returns.
	resultData := map[string]interface{}{
		"callbackType": callback.token,
		"status":       "success",
		"result":       "Async started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

// ProcessObject transforms an input object and returns derived fields that exercise object marshalling.
//
//export ProcessObject
func ProcessObject(params *C.char, callbackType *C.char) *C.char {
	// Decode the incoming object so the sample can enrich it with derived metadata.
	var objectData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &objectData)

	// Seed the response with generic processing metadata before copying user-provided fields.
	processed := map[string]interface{}{
		"processed": true,
		"timestamp": time.Now().Unix(),
	}

	for key, value := range objectData {
		processed[key] = value
	}

	if name, ok := objectData["name"].(string); ok {
		processed["nameLength"] = len(name)
		processed["nameUpperCase"] = name
	}

	if age, ok := objectData["age"].(float64); ok {
		processed["isAdult"] = age >= 18
		processed["ageInDays"] = int(age * 365)
	}

	if items, ok := objectData["items"].([]interface{}); ok {
		processed["itemCount"] = len(items)
	}

	// Serialize the enriched object back to JSON for the wrapper to return into JavaScript.
	resultJson, _ := json.Marshal(processed)
	return C.CString(string(resultJson))
}
