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
	fmt.Println("Go: RegisterGoCallback called, fn=", fn)
}

//export Hello1
func Hello1(name *C.char, value *C.char) *C.char {
	nameStr := C.GoString(name)
	valueStr := C.GoString(value)
	fmt.Printf("Go: Hello1 called with name=%s, value=%s\n", nameStr, valueStr)

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
	fmt.Printf("Go: Hello2 called with name=%s, value=%s\n", nameStr, valueStr)

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
	fmt.Printf("Go: Hello3 called with name=%s, value=%s\n", nameStr, valueStr)

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
	cbType := C.GoString(callbackType)
	fmt.Printf("Go: Hello4 called with name=%s, callbackType=%s\n", nameStr, cbType)

	var inputData map[string]interface{}
	json.Unmarshal([]byte(nameStr), &inputData)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

	if gCallNodeCallback != 0 {
		for i := 1; i <= 3; i++ {
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("Go: Hello4 Triggering callback #%d\n", i)

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
	cbType := C.GoString(callbackType)
	fmt.Printf("Go: Hello5 called with name=%s, callbackType=%s\n", nameStr, cbType)

	var inputData map[string]interface{}
	json.Unmarshal([]byte(nameStr), &inputData)

	testMsg := "default"
	if msg, ok := inputData["test"].(string); ok {
		testMsg = msg
	}

	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Go: Hello5 Triggering async callback #%d\n", i)

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
	cbType := C.GoString(callbackType)
	fmt.Printf("Go: Hello6 called with name=%s, callbackType=%s\n", nameStr, cbType)

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
			fmt.Printf("Go: Hello6 Triggering infinite callback #%d\n", i)

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

func main() {}
