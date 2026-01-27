package main

/*
#cgo CFLAGS: -I.
#include <stdlib.h>
#include <stdint.h>
#include <windows.h>
#include <string.h>

typedef void (*CallbackFunc)(const char*, const char*);

static void callCallback(void* ptr, const char* callbackType, const char* data) {
    if (ptr != NULL) {
        ((CallbackFunc)ptr)(callbackType, data);
    }
}
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

//export Hello1
func Hello1(name *C.char, value *C.char) *C.char {
	nameStr := C.GoString(name)
	valueStr := C.GoString(value)

	valueInt, _ := strconv.Atoi(valueStr)
	result := valueInt * 2

	resultData := map[string]interface{}{
		"name":   nameStr,
		"value":  valueInt,
		"result": result,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export Hello2
func Hello2(name *C.char, value *C.char) *C.char {
	nameStr := C.GoString(name)
	valueStr := C.GoString(value)

	valueFloat, _ := strconv.ParseFloat(valueStr, 64)
	result := valueFloat * 1.5

	resultData := map[string]interface{}{
		"name":   nameStr,
		"value":  valueFloat,
		"result": result,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export Hello3
func Hello3(name *C.char, value *C.char) *C.char {
	nameStr := C.GoString(name)
	valueStr := C.GoString(value)

	valueFloat, _ := strconv.ParseFloat(valueStr, 64)
	result := valueFloat > 0.0

	resultData := map[string]interface{}{
		"name":   nameStr,
		"value":  valueFloat,
		"result": result,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export Hello4
func Hello4(name *C.char, callbackType *C.char) *C.char {
	nameStr := C.GoString(name)

	var inputData map[string]interface{}
	json.Unmarshal([]byte(nameStr), &inputData)

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

			C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString("sync_callback"), C.CString(string(jsonData)))
		}
	}

	resultData := map[string]interface{}{
		"status": "success",
		"result": 42,
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export Hello5
func Hello5(name *C.char, callbackType *C.char) *C.char {
	nameStr := C.GoString(name)

	var inputData map[string]interface{}
	json.Unmarshal([]byte(nameStr), &inputData)

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

			C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString("async_callback"), C.CString(string(jsonData)))
		}
	}()

	resultData := map[string]interface{}{
		"status": "success",
		"result": "Async started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export Hello6
func Hello6(name *C.char, callbackType *C.char) *C.char {
	nameStr := C.GoString(name)

	var inputData map[string]interface{}
	json.Unmarshal([]byte(nameStr), &inputData)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

	go func() {
		i := 1
		for {
			time.Sleep(1000 * time.Millisecond)

			callbackData := map[string]interface{}{
				"test":   testMsg,
				"result": fmt.Sprintf("Infinite callback %d", i),
			}
			jsonData, _ := json.Marshal(callbackData)

			C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString("infinite_callback"), C.CString(string(jsonData)))
			i++
		}
	}()

	resultData := map[string]interface{}{
		"status": "success",
		"result": "Infinite async started",
	}
	resultJson, _ := json.Marshal(resultData)

	return C.CString(string(resultJson))
}

//export ReturnString
func ReturnString(name *C.char, value *C.char) *C.char {
	nameStr := C.GoString(name)
	valueStr := C.GoString(value)

	result := map[string]interface{}{
		"_type": "string",
		"value": "Hello " + nameStr + ", your value is " + valueStr,
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

//export ReturnInt
func ReturnInt(name *C.char, value *C.char) *C.char {
	valueStr := C.GoString(value)

	var valueInt int
	fmt.Sscanf(valueStr, "%d", &valueInt)

	result := map[string]interface{}{
		"_type": "int",
		"value": valueInt * 2,
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

//export ReturnFloat
func ReturnFloat(name *C.char, value *C.char) *C.char {
	valueStr := C.GoString(value)

	var valueFloat float64
	fmt.Sscanf(valueStr, "%f", &valueFloat)

	result := map[string]interface{}{
		"_type": "float",
		"value": valueFloat * 1.5,
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

//export ReturnBool
func ReturnBool(name *C.char, value *C.char) *C.char {
	valueStr := C.GoString(value)

	var valueFloat float64
	fmt.Sscanf(valueStr, "%f", &valueFloat)

	result := map[string]interface{}{
		"_type": "bool",
		"value": valueFloat > 0.0,
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

//export ReturnObject
func ReturnObject(name *C.char, value *C.char) *C.char {
	nameStr := C.GoString(name)
	valueStr := C.GoString(value)

	var valueInt int
	fmt.Sscanf(valueStr, "%d", &valueInt)

	result := map[string]interface{}{
		"_type": "object",
		"value": map[string]interface{}{
			"name":     nameStr,
			"age":      valueInt,
			"isActive": true,
			"scores":   []int{85, 90, 78},
			"address": map[string]string{
				"city":    "Beijing",
				"country": "China",
			},
		},
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

//export ReturnNestedObject
func ReturnNestedObject(name *C.char, value *C.char) *C.char {
	nameStr := C.GoString(name)

	result := map[string]interface{}{
		"_type": "object",
		"value": map[string]interface{}{
			"user": map[string]interface{}{
				"name": nameStr,
				"age":  30,
			},
			"metadata": map[string]interface{}{
				"created": "2024-01-01",
				"tags":    []string{"tag1", "tag2"},
			},
			"items": []map[string]interface{}{
				{"id": 1, "name": "Item 1"},
				{"id": 2, "name": "Item 2"},
				{"id": 3, "name": "Item 3"},
			},
		},
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

//export ReturnWithCallback
func ReturnWithCallback(name *C.char, callbackType *C.char) *C.char {
	if gCallNodeCallback != 0 {
		C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString("test_callback"), C.CString(`{"message":"Callback from Go"}`))
	}

	result := map[string]interface{}{
		"_type": "string",
		"value": "Callback triggered",
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

func main() {}
