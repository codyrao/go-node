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
	"strings"
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
		for i := 1; i <= 5; i++ {
			time.Sleep(1000 * time.Millisecond)

			callbackData := map[string]interface{}{
				"test":   testMsg,
				"result": fmt.Sprintf("Limited callback %d", i),
			}
			jsonData, _ := json.Marshal(callbackData)

			C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString("infinite_callback"), C.CString(string(jsonData)))
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

	result := "Hello " + nameStr + ", your value is " + valueStr
	return C.CString(result)
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
		callbackData := map[string]interface{}{
			"message": "Callback from Go",
			"status":  "success",
			"data": map[string]interface{}{
				"id":        123,
				"name":      C.GoString(name),
				"timestamp": time.Now().Unix(),
			},
		}
		jsonData, _ := json.Marshal(callbackData)
		C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString("test_callback"), C.CString(string(jsonData)))
	}

	result := map[string]interface{}{
		"_type": "string",
		"value": "Callback triggered",
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

//export ReturnWithObjectCallback
func ReturnWithObjectCallback(name *C.char, callbackType *C.char) *C.char {
	if gCallNodeCallback != 0 {
		callbackData := map[string]interface{}{
			"message": "Object callback from Go",
			"status":  "success",
			"user": map[string]interface{}{
				"id":       456,
				"username": C.GoString(name),
				"email":    "test@example.com",
				"roles":    []string{"admin", "user"},
			},
			"timestamp": time.Now().Unix(),
		}
		jsonData, _ := json.Marshal(callbackData)
		C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString("object_callback"), C.CString(string(jsonData)))
	}

	result := map[string]interface{}{
		"_type": "string",
		"value": "Object callback triggered",
	}
	resultJson, _ := json.Marshal(result)

	return C.CString(string(resultJson))
}

//export ProcessArray
func ProcessArray(jsonArray *C.char) *C.char {
	arrayStr := C.GoString(jsonArray)

	var arrayData []interface{}
	if err := json.Unmarshal([]byte(arrayStr), &arrayData); err != nil {
		result := map[string]interface{}{
			"_type": "error",
			"value": "Failed to parse array: " + err.Error(),
		}
		resultJson, _ := json.Marshal(result)
		return C.CString(string(resultJson))
	}

	// Process array - sum numbers, count strings, etc.
	var sum float64
	var count int
	var strings []string

	for _, item := range arrayData {
		switch v := item.(type) {
		case float64:
			sum += v
			count++
		case string:
			strings = append(strings, v)
			count++
		case int:
			sum += float64(v)
			count++
		}
	}

	result := map[string]interface{}{
		"_type": "object",
		"value": map[string]interface{}{
			"originalArray": arrayData,
			"itemCount":     count,
			"sum":           sum,
			"strings":       strings,
			"processed":     true,
		},
	}

	resultJson, _ := json.Marshal(result)
	return C.CString(string(resultJson))
}

//export FilterArray
func FilterArray(jsonArray *C.char) *C.char {
	arrayStr := C.GoString(jsonArray)

	var arrayData []map[string]interface{}
	if err := json.Unmarshal([]byte(arrayStr), &arrayData); err != nil {
		result := map[string]interface{}{
			"_type": "error",
			"value": "Failed to parse array: " + err.Error(),
		}
		resultJson, _ := json.Marshal(result)
		return C.CString(string(resultJson))
	}

	// Filter array - return only items with specific conditions
	var filtered []map[string]interface{}
	for _, item := range arrayData {
		if name, ok := item["name"].(string); ok && len(name) > 0 {
			if age, ok := item["age"].(float64); ok && age >= 18 {
				filtered = append(filtered, item)
			}
		}
	}

	result := map[string]interface{}{
		"_type": "array",
		"value": filtered,
	}

	resultJson, _ := json.Marshal(result)
	return C.CString(string(resultJson))
}

//export ProcessObject
func ProcessObject(jsonObject *C.char) *C.char {
	objectStr := C.GoString(jsonObject)

	var objectData map[string]interface{}
	if err := json.Unmarshal([]byte(objectStr), &objectData); err != nil {
		result := map[string]interface{}{
			"_type": "error",
			"value": "Failed to parse object: " + err.Error(),
		}
		resultJson, _ := json.Marshal(result)
		return C.CString(string(resultJson))
	}

	// Process object - extract and transform data
	processed := map[string]interface{}{
		"processed": true,
		"timestamp": time.Now().Unix(),
	}

	// Copy all original fields
	for key, value := range objectData {
		processed[key] = value
	}

	// Add some computed fields
	if name, ok := objectData["name"].(string); ok {
		processed["nameLength"] = len(name)
		processed["nameUpperCase"] = strings.ToUpper(name)
	}

	if age, ok := objectData["age"].(float64); ok {
		processed["isAdult"] = age >= 18
		processed["ageInDays"] = int(age * 365)
	}

	if items, ok := objectData["items"].([]interface{}); ok {
		processed["itemCount"] = len(items)
	}

	result := map[string]interface{}{
		"_type": "object",
		"value": processed,
	}

	resultJson, _ := json.Marshal(result)
	return C.CString(string(resultJson))
}

func main() {}
