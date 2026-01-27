# Go2Node - Compile Go Code to Node.js Native Module

A command-line tool implemented in Go that compiles Go code into Node.js callable native modules (.node files).

## Features

- **Based on CGO and node-gyp**: Uses Go's c-shared compilation mode and node-gyp to build native modules
- **Automatic bridge code generation**: Automatically generates C++ wrapper and binding.gyp files
- **Callback support**: Supports calling Node.js callback functions from Go (synchronous and asynchronous)
- **Multi-type return support**: Supports returning multiple Node.js types (object, int, float, bool, string)
- **Object return support**: Supports returning JSON-formatted nested JavaScript objects
- **Flexible configuration**: Supports custom output directory and module name
- **Temporary file management**: Automatically manages temporary files with random IDs to avoid conflicts

## Prerequisites

Before running the tool, ensure the following dependencies are installed:

### 1. Go Language
```bash
# Download: https://go.dev/dl/
# Version requirement: 1.23.x
```

### 2. Node.js and npm
```bash
# Download: https://nodejs.org/
# Version requirement: 16.x or higher
```

### 3. node-gyp (Global installation)
```bash
npm install -g node-gyp
```

### 4. Python 3.x
```bash
# Windows: https://www.python.org/downloads/
# Add to PATH environment variable
```

### 5. C++ Build Tools (Windows)
- Visual Studio 2022 or higher
- Or install Build Tools for Visual Studio

## Installation

```bash
go install github.com/codyrao/go-node@latest
```

After installation, the `go-node` executable will be located in the `$GOPATH/bin` directory (Windows: `%USERPROFILE%\go\bin`). Make sure this directory is added to the system PATH environment variable.

## Usage

### Basic Usage

```bash
go-node -input=your_file.go -name=module_name
```

### Parameters

| Parameter | Required | Default | Description |
|-----------|-----------|----------|-------------|
| `-input` | Yes | - | Path to the input Go source file |
| `-name` | Yes | - | Module name (generates .node file name) |
| `-package` | No | `main` | Go package name (only functions in main package will be compiled) |
| `-output` | No | `./output` | Output directory |
| `-source` | No | `-input` directory | Go source file directory |
| `-no-cleanup` | No | false | Do not clean up temporary files after compilation |

### Examples

```bash
# Basic compilation
go-node -input=hello.go -name=hello

# Specify output directory
go-node -input=hello.go -name=hello -output=./dist

# Do not clean up temporary files (for debugging)
go-node -input=hello.go -name=hello -no-cleanup
```

## Go Code Writing Guidelines

### Export Functions

Use `//export` to mark export functions:

```go
package main

import "C"

//export Hello
func Hello() {
    println("Hello from Go!")
}

//export Add
func Add(a, b int) int {
    return a + b
}

func main() {}
```

### Supported Function Signatures

The tool supports multiple function signature patterns:

#### 1. Simple Synchronous Functions

```go
//export Hello1
func Hello1(name *C.char, value *C.char) *C.char {
    // name and value are string parameters
    // Return JSON string
    result := map[string]interface{}{
        "name":   C.GoString(name),
        "value":  C.GoString(value),
        "result": "success",
    }
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}
```

#### 2. Synchronous Callback Functions

```go
//export Hello4
func Hello4(name *C.char, callbackType *C.char) *C.char {
    if C.GoString(callbackType) == "callback" {
        // Trigger synchronous callbacks
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 1"}`))
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 2"}`))
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 3"}`))
    }
    return C.CString(`{"result":42,"status":"success"}`)
}
```

#### 3. Asynchronous Callback Functions

```go
//export Hello5
func Hello5(name *C.char, callbackType *C.char) *C.char {
    if C.GoString(callbackType) == "callback" {
        go func() {
            for i := 1; i <= 5; i++ {
                time.Sleep(500 * time.Millisecond)
                CallCallback(C.CString("async_callback"), C.CString(fmt.Sprintf(`{"result":"Async callback %d"}`, i)))
            }
        }()
    }
    return C.CString(`{"result":"Async started","status":"success"}`)
}
```

#### 4. Infinite Asynchronous Callback Functions

```go
//export Hello6
func Hello6(name *C.char, callbackType *C.char) *C.char {
    if C.GoString(callbackType) == "callback" {
        go func() {
            for i := 1; ; i++ {
                time.Sleep(1 * time.Second)
                CallCallback(C.CString("infinite_callback"), C.CString(fmt.Sprintf(`{"result":"Infinite callback %d"}`, i)))
            }
        }()
    }
    return C.CString(`{"result":"Infinite async started","status":"success"}`)
}
```

### Callback Support

The tool has built-in callback mechanism, Go code can call the exported `CallCallback` function:

```go
//export ProcessWithCallback
func ProcessWithCallback() {
    // Call callback from Go code
    response := `{"status": "success"}`
    CallCallback(C.CString("callback_type"), C.CString(response))
}
```

### Multi-Type Return Support

Go functions can return multiple Node.js types, including: object, int, float, bool, string.

#### Return Type Format

Return values need to use JSON format, containing `_type` and `value` fields:

```go
result := map[string]interface{}{
    "_type": "type_name",  // object, int, float, bool, string
    "value": actual_value,
}
jsonBytes, _ := json.Marshal(result)
return C.CString(string(jsonBytes))
```

#### Supported Return Types

**1. String Type (string)**
```go
//export ReturnString
func ReturnString(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

    result := map[string]interface{}{
        "_type":  "string",
        "value":  "Hello " + nameStr + ", your value is " + valueStr,
    }
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}
```

**2. Integer Type (int)**
```go
//export ReturnInt
func ReturnInt(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

    var valueInt int
    fmt.Sscanf(valueStr, "%d", &valueInt)

    result := map[string]interface{}{
        "_type":  "int",
        "value":  valueInt * 2,
    }
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}
```

**3. Float Type (float)**
```go
//export ReturnFloat
func ReturnFloat(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

    var valueFloat float64
    fmt.Sscanf(valueStr, "%f", &valueFloat)

    result := map[string]interface{}{
        "_type":  "float",
        "value":  valueFloat * 1.5,
    }
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}
```

**4. Boolean Type (bool)**
```go
//export ReturnBool
func ReturnBool(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

    var valueFloat float64
    fmt.Sscanf(valueStr, "%f", &valueFloat)

    result := map[string]interface{}{
        "_type":  "bool",
        "value":  valueFloat > 0.0,
    }
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}
```

**5. Object Type (object)**
```go
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
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}
```

**6. Nested Object Type**
```go
//export ReturnNestedObject
func ReturnNestedObject(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

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
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}
```

#### Node.js Usage

```javascript
const types = require('./output/types.node')

// String type
const strResult = types.ReturnString("Test", "World")
console.log(strResult)  // "Hello Test, your value is World"
console.log(typeof strResult)  // "string"

// Integer type
const intResult = types.ReturnInt("Test", "10")
console.log(intResult)  // 20
console.log(typeof intResult)  // "number"

// Float type
const floatResult = types.ReturnFloat("Test", "3.14")
console.log(floatResult)  // 4.71
console.log(typeof floatResult)  // "number"

// Boolean type
const boolResult = types.ReturnBool("Test", "5.5")
console.log(boolResult)  // true
console.log(typeof boolResult)  // "boolean"

// Object type
const objResult = types.ReturnObject("Test", "30")
console.log(objResult.name)  // "Test"
console.log(objResult.age)  // 30
console.log(objResult.isActive)  // true
console.log(objResult.scores)  // [85, 90, 78]
console.log(objResult.address.city)  // "Beijing"

// Nested object type
const nestedResult = types.ReturnNestedObject("Test", "100")
console.log(nestedResult.user.name)  // "Test"
console.log(nestedResult.user.age)  // 30
console.log(nestedResult.metadata.created)  // "2024-01-01"
console.log(nestedResult.items[0].name)  // "Item 1"
```

#### Type Conversion Rules

| Go Return Type | Node.js Type | Description |
|----------------|----------------|-------------|
| `_type: "string"` | string | JavaScript string |
| `_type: "int"` | number | JavaScript number (integer) |
| `_type: "float"` | number | JavaScript number (float) |
| `_type: "bool"` | boolean | JavaScript boolean |
| `_type: "object"` | object | JavaScript object (supports nesting) |

Note: If the returned JSON does not contain the `_type` field, the JSON-parsed raw object will be returned directly.

## Node.js Usage

### Import Module

```javascript
const demoaddon = require('./output/hello.node')
```

### Call Synchronous Functions

```javascript
// Simple call
const result = demoaddon.Hello1("Test1", "10")
console.log(result)  // JSON string
```

### Call Callback Functions

```javascript
// Synchronous callback
demoaddon.Hello4('{"test":"Hello4 Test"}', function(callbackType, jsonData) {
    console.log(`Callback [${callbackType}]:`, jsonData)
})

// Asynchronous callback
demoaddon.Hello5('{"test":"Hello5 Test"}', function(callbackType, jsonData) {
    console.log(`Async callback [${callbackType}]:`, jsonData)
})
```

## Complete Example

### Go Code (hello.go)

```go
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
    jsonBytes, _ := json.Marshal(resultData)

    return C.CString(string(jsonBytes))
}

//export Hello4
func Hello4(name *C.char, callbackType *C.char) *C.char {
    nameStr := C.GoString(name)
    cbType := C.GoString(callbackType)

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
    jsonBytes, _ := json.Marshal(resultData)

    return C.CString(string(jsonBytes))
}

//export Hello5
func Hello5(name *C.char, callbackType *C.char) *C.char {
    nameStr := C.GoString(name)
    cbType := C.GoString(callbackType)

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
    jsonBytes, _ := json.Marshal(resultData)

    return C.CString(string(jsonBytes))
}

//export ReturnString
func ReturnString(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

    result := map[string]interface{}{
        "_type":  "string",
        "value":  "Hello " + nameStr + ", your value is " + valueStr,
    }
    jsonBytes, _ := json.Marshal(result)

    return C.CString(string(jsonBytes))
}

//export ReturnInt
func ReturnInt(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

    var valueInt int
    fmt.Sscanf(valueStr, "%d", &valueInt)

    result := map[string]interface{}{
        "_type":  "int",
        "value":  valueInt * 2,
    }
    jsonBytes, _ := json.Marshal(result)

    return C.CString(string(jsonBytes))
}

//export ReturnFloat
func ReturnFloat(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

    var valueFloat float64
    fmt.Sscanf(valueStr, "%f", &valueFloat)

    result := map[string]interface{}{
        "_type":  "float",
        "value":  valueFloat * 1.5,
    }
    jsonBytes, _ := json.Marshal(result)

    return C.CString(string(jsonBytes))
}

//export ReturnBool
func ReturnBool(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

    var valueFloat float64
    fmt.Sscanf(valueStr, "%f", &valueFloat)

    result := map[string]interface{}{
        "_type":  "bool",
        "value":  valueFloat > 0.0,
    }
    jsonBytes, _ := json.Marshal(result)

    return C.CString(string(jsonBytes))
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
    jsonBytes, _ := json.Marshal(result)

    return C.CString(string(jsonBytes))
}

//export ReturnNestedObject
func ReturnNestedObject(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)

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
    jsonBytes, _ := json.Marshal(result)

    return C.CString(string(jsonBytes))
}

func main() {}
```

### Node.js Code (test.js)

```javascript
const demoaddon = require('./output/hello.node')

console.log('=== Go2Node Multi-function Test ===\n')

// Test 1: Synchronous call Hello1(name string, value int) int
console.log('Test 1: Hello1 - Synchronous call (name string, value int) -> int')
const result1 = demoaddon.Hello1("Test1", "10")
console.log('   Result:', result1)
console.log()

// Test 2: Synchronous call Hello2(name string, value float) float
console.log('Test 2: Hello2 - Synchronous call (name string, value float) -> float')
const result2 = demoaddon.Hello2("Test2", "3.14")
console.log('   Result:', result2)
console.log()

// Test 3: Synchronous call Hello3(name bool, value float) bool
console.log('Test 3: Hello3 - Synchronous call (name bool, value float) -> bool')
const result3 = demoaddon.Hello3("true", "5.5")
console.log('   Result:', result3)
console.log()

// Test 4: Synchronous call - synchronous callback Hello4(name, callback(test string, res string)) int
console.log('Test 4: Hello4 - Synchronous call - synchronous callback')
const params4 = {
    test: "Hello4 Test"
}
const result4 = demoaddon.Hello4(JSON.stringify(params4), function(callbackType, jsonData) {
    console.log('   Callback [Hello4]:', callbackType, '->', jsonData)
})
console.log('   Result:', result4)
console.log()

// Test 5: Synchronous call - asynchronous callback Hello5(name, async_callback(test string, res string)) string
console.log('Test 5: Hello5 - Synchronous call - asynchronous callback')
const params5 = {
    test: "Hello5 Test"
}
const result5 = demoaddon.Hello5(JSON.stringify(params5), function(callbackType, jsonData) {
    console.log('   Async callback [Hello5]:', callbackType, '->', jsonData)
})
console.log('   Result:', result5)
console.log('   Waiting for async callbacks...\n')

// Test 6: Return string type
console.log('Test 6: ReturnString - Return string type')
const result6 = demoaddon.ReturnString("Test6", "World")
console.log('   Result:', result6)
console.log('   Type:', typeof result6)
console.log()

// Test 7: Return integer type
console.log('Test 7: ReturnInt - Return integer type')
const result7 = demoaddon.ReturnInt("Test7", "10")
console.log('   Result:', result7)
console.log('   Type:', typeof result7)
console.log()

// Test 8: Return float type
console.log('Test 8: ReturnFloat - Return float type')
const result8 = demoaddon.ReturnFloat("Test8", "3.14")
console.log('   Result:', result8)
console.log('   Type:', typeof result8)
console.log()

// Test 9: Return boolean type
console.log('Test 9: ReturnBool - Return boolean type')
const result9 = demoaddon.ReturnBool("Test9", "5.5")
console.log('   Result:', result9)
console.log('   Type:', typeof result9)
console.log()

// Test 10: Return object type
console.log('Test 10: ReturnObject - Return object type')
const result10 = demoaddon.ReturnObject("Test10", "30")
console.log('   Result:', result10)
console.log('   Type:', typeof result10)
console.log('   name:', result10.name)
console.log('   age:', result10.age)
console.log('   isActive:', result10.isActive)
console.log('   scores:', result10.scores)
console.log('   address:', result10.address)
console.log()

// Test 11: Return nested object type
console.log('Test 11: ReturnNestedObject - Return nested object type')
const result11 = demoaddon.ReturnNestedObject("Test11", "100")
console.log('   Result:', result11)
console.log('   Type:', typeof result11)
console.log('   user.name:', result11.user.name)
console.log('   user.age:', result11.user.age)
console.log('   metadata:', result11.metadata)
console.log('   items:', result11.items)
console.log()

console.log('=== Basic tests completed ===')
console.log('Waiting for async callbacks...')
```

### Compile and Run

```bash
# Compile
go-node -input=test/hello.go -name=hello

# Run test
node test.js
```

## Output Structure

After compilation, the output directory structure is as follows:

```
output/
└── module_name.node    # Node.js native module (Go dynamic library embedded)
```

**Note**: The Go dynamic library is embedded in the .node file and will be automatically extracted to the system temporary directory at runtime.

## Temporary File Management

The tool automatically manages temporary files:

1. Creates a `go-node-tmp_{randomID}` temporary directory in the source file directory
2. Generates `binding.gyp` and `wrapper.cc` files in the temporary directory
3. Automatically cleans up the temporary directory after compilation
4. Use `-no-cleanup` parameter to keep temporary files for debugging

## Testing

View `test/hello.go` for complete example code, run tests:

```bash
# Compile example
go-node -input=test/hello.go -name=hello

# Run test
node test.js
```

## Known Limitations

1. **Package name limitation**: Only supports compilation of functions in main package
2. **Platform specific**: Current implementation is mainly for Windows platform
3. **Callback context**: Asynchronous callbacks use libuv to ensure execution in the correct V8 context
4. **Parameter types**: Function parameters are uniformly passed as strings, complex types use JSON serialization

## Technical Implementation

- **Go compilation**: Uses `go build -buildmode=c-shared` to generate dynamic library
- **Node.js binding**: Uses node-gyp to generate C++ native module
- **Cross-DLL communication**: Uses Windows API dynamic loading and function pointers
- **Asynchronous callbacks**: Uses libuv's uv_async_t for thread-safe asynchronous callbacks
- **Callback queue**: Uses mutex-protected queue to manage callback requests

## License

MIT
