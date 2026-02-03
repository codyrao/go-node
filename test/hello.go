package main

/*
#cgo CFLAGS: -I.
#include "callback.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"time"
	"unsafe"
)

var gCallNodeCallback uintptr

//export RegisterGoCallback
func RegisterGoCallback(fn uintptr) {
	gCallNodeCallback = fn
}

//export Hello1
func Hello1(params *C.char, callbackType *C.char) *C.char {
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	name := ""
	if n, ok := inputData["name"].(string); ok {
		name = n
	}

	value := 0
	if v, ok := inputData["value"].(float64); ok {
		value = int(v)
	}

	result := value * 2

	resultData := map[string]interface{}{
		"name":   name,
		"value":  value,
		"result": result,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export Hello2
func Hello2(params *C.char, callbackType *C.char) *C.char {
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	name := ""
	if n, ok := inputData["name"].(string); ok {
		name = n
	}

	value := 0.0
	if v, ok := inputData["value"].(float64); ok {
		value = v
	}

	result := value * 1.5

	resultData := map[string]interface{}{
		"name":   name,
		"value":  value,
		"result": result,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export HelloWithCallback
func HelloWithCallback(params *C.char, callbackType *C.char) *C.char {
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

	if gCallNodeCallback != 0 {
		for i := 1; i <= 3; i++ {
			time.Sleep(300 * time.Millisecond)

			callbackData := map[string]interface{}{
				"test":   testMsg,
				"result": fmt.Sprintf("Callback %d", i),
			}
			jsonData, _ := json.Marshal(callbackData)

			C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(jsonData)))
		}
	}

	resultData := map[string]interface{}{
		"status": "success",
		"result": 42,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export AsyncHello
func AsyncHello(params *C.char, callbackType *C.char) *C.char {
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(500 * time.Millisecond)

			callbackData := map[string]interface{}{
				"test":   testMsg,
				"result": fmt.Sprintf("Async callback %d", i),
			}
			jsonData, _ := json.Marshal(callbackData)

			C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(jsonData)))
		}
	}()

	resultData := map[string]interface{}{
		"status": "success",
		"result": "Async started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export ProcessObject
func ProcessObject(params *C.char, callbackType *C.char) *C.char {
	var objectData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &objectData)

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

	resultJson, _ := json.Marshal(processed)
	return C.CString(string(resultJson))
}

//export Calculate
func Calculate(params *C.char, callbackType *C.char) *C.char {
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	a := 0.0
	if v, ok := inputData["a"].(float64); ok {
		a = v
	}

	b := 0.0
	if v, ok := inputData["b"].(float64); ok {
		b = v
	}

	operation := "add"
	if op, ok := inputData["operation"].(string); ok {
		operation = op
	}

	var result float64
	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b != 0 {
			result = a / b
		}
	}

	resultData := map[string]interface{}{
		"a":         a,
		"b":         b,
		"operation": operation,
		"result":    result,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export NoReturn
func NoReturn(params *C.char, callbackType *C.char) *C.char {
	return nil
}

func main() {}
