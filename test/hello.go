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

type nodeCallback struct {
	id    int32
	token string
}

func newNodeCallback(raw *C.char) nodeCallback {
	if gCallNodeCallback == 0 || raw == nil {
		return nodeCallback{id: -1}
	}

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

func (cb nodeCallback) send(data map[string]interface{}) {
	if gCallNodeCallback == 0 || cb.id < 0 {
		return
	}

	jsonData, _ := json.Marshal(data)
	cJSON := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(cJSON))

	C.callCallbackWithId(unsafe.Pointer(gCallNodeCallback), C.int32_t(cb.id), cJSON)
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

//export HelloWithCallback
func HelloWithCallback(params *C.char, callbackType *C.char) *C.char {
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	callback := newNodeCallback(callbackType)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

	for i := 1; i <= 3; i++ {
		time.Sleep(300 * time.Millisecond)

		callback.send(map[string]interface{}{
			"callbackType": callback.token,
			"test":         testMsg,
			"result":       fmt.Sprintf("Callback %d", i),
		})
	}

	resultData := map[string]interface{}{
		"callbackType": callback.token,
		"status":       "success",
		"result":       42,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export AsyncHello
func AsyncHello(params *C.char, callbackType *C.char) *C.char {
	var inputData map[string]interface{}
	json.Unmarshal([]byte(C.GoString(params)), &inputData)

	callback := newNodeCallback(callbackType)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

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

	resultData := map[string]interface{}{
		"callbackType": callback.token,
		"status":       "success",
		"result":       "Async started",
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

func main() {}
