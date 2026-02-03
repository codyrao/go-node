# Go2Node - Compile Go Code to Node.js Native Module

A command-line tool implemented in Go that compiles Go code into Node.js callable native modules (.node files).

## Features

- **Based on CGO and node-gyp**: Uses Go's c-shared compilation mode and node-gyp to build native modules
- **Automatic bridge code generation**: Automatically generates C++ wrapper and binding.gyp files
- **Callback support**: Supports calling Node.js callback functions from Go (synchronous and asynchronous)
- **Electron version support**: Supports specifying Electron version via `-ev` parameter
- **Fixed function signature**: All Go export functions use a consistent two-parameter structure
- **Object return only**: Functions return either 0 or 1 object type value
- **Object callback parameter**: Callback functions receive a single object parameter
- **Automatic callback JSON parsing**: Automatically parses JSON strings to JavaScript objects in callbacks
- **Flexible DLL loading**: Supports loading DLL from local directory or embedded resources
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
| `-ev` | No | - | Electron version (e.g., 28.0.0). If not specified, uses node-gyp's default Node.js version |

### Examples

```bash
# Basic compilation
go-node -input=hello.go -name=hello

# Specify output directory
go-node -input=hello.go -name=hello -output=./dist

# Compile for Electron
go-node -input=hello.go -name=hello -ev=28.0.0

# Do not clean up temporary files (for debugging)
go-node -input=hello.go -name=hello -no-cleanup
```

### Output Files

After successful compilation, the following files will be generated in the output directory:

| File | Description |
|------|-------------|
| `module_name.node` | Node.js native module (can be directly required in Node.js) |
| `module_name.dll` | Go DLL (for standalone use or debugging) |

The `.node` file contains an embedded DLL and can be used independently. The `.dll` file is also provided for debugging or direct use.

### DLL Loading Mechanism

The tool uses a flexible DLL loading strategy:

1. **Primary**: Attempts to load DLL from the same directory as the `.node` file
2. **Fallback**: If not found, attempts to extract DLL from embedded resources in the `.node` file

This ensures reliable loading in various deployment scenarios.

## Go Code Writing Guidelines

### Function Signature Requirements

All Go export functions must follow this fixed signature:

```go
func FunctionName(params *C.char, callbackType *C.char) *C.char
```

**Parameter Structure:**
- **First parameter (`params`)**: Object type parameter (JSON string)
- **Second parameter (`callbackType`)**: Callback type indicator ("callback" or empty string)

**Return Value:**
- **Return type**: `*C.char` (JSON string or nil)
- **Return format**: JSON object or nil

### Export Functions

Use `//export` to mark export functions:

```go
package main

import "C"

//export Hello
func Hello(params *C.char, callbackType *C.char) *C.char {
    // params is a JSON string representing an object
    // callbackType is "callback" if a callback is provided
    // Return a JSON string or nil
}

func main() {}
```

### Basic Function Example

```go
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
```

### Callback Functions

#### Synchronous Callback

```go
//export HelloWithCallback
func HelloWithCallback(params *C.char, callbackType *C.char) *C.char {
    var inputData map[string]interface{}
    json.Unmarshal([]byte(C.GoString(params)), &inputData)

    testMsg := "default"
    if msg, ok := inputData["test"].(string); ok {
        testMsg = msg
    }

    if C.GoString(callbackType) == "callback" && gCallNodeCallback != 0 {
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
```

#### Asynchronous Callback

```go
//export AsyncHello
func AsyncHello(params *C.char, callbackType *C.char) *C.char {
    var inputData map[string]interface{}
    json.Unmarshal([]byte(C.GoString(params)), &inputData)

    testMsg := "default"
    if msg, ok := inputData["test"].(string); ok {
        testMsg = msg
    }

    if C.GoString(callbackType) == "callback" && gCallNodeCallback != 0 {
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
    }

    resultData := map[string]interface{}{
        "status": "success",
        "result": "Async started",
    }
    resultJson, _ := json.Marshal(resultData)

    return C.CString(string(resultJson))
}
```

### No Return Value

```go
//export NoReturn
func NoReturn(params *C.char, callbackType *C.char) *C.char {
    return nil
}
```

### Object Processing Example

```go
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
```

## Node.js Usage

### Import Module

```javascript
const hello = require('./output/hello.node')
```

### Call Basic Functions

```javascript
// Call function with object parameter
const result1 = hello.hello1({ name: 'Alice', value: 21 })
console.log(result1)
// Output: { name: 'Alice', value: 21, result: 42 }
```

### Call Callback Functions

```javascript
// Synchronous callback
const result2 = hello.helloWithCallback({ test: 'Hello from Node' }, (data) => {
    console.log('Callback received:', data)
    // Output: { test: 'Hello from Node', result: 'Callback 1' }
    // Output: { test: 'Hello from Node', result: 'Callback 2' }
    // Output: { test: 'Hello from Node', result: 'Callback 3' }
})
console.log('Result:', result2)
// Output: { status: 'success', result: 42 }

// Asynchronous callback
const result3 = hello.asyncHello({ test: 'Async test' }, (data) => {
    console.log('Async callback received:', data)
    // Output: { test: 'Async test', result: 'Async callback 1' }
    // Output: { test: 'Async test', result: 'Async callback 2' }
    // ... (5 callbacks total)
})
console.log('Result:', result3)
// Output: { status: 'success', result: 'Async started' }
```

### Call Object Processing Functions

```javascript
const result4 = hello.processObject({
    name: 'Charlie',
    age: 25,
    items: ['item1', 'item2', 'item3']
})
console.log(result4)
// Output: {
//   name: 'Charlie',
//   age: 25,
//   items: ['item1', 'item2', 'item3'],
//   nameLength: 7,
//   nameUpperCase: 'Charlie',
//   isAdult: true,
//   ageInDays: 9125,
//   itemCount: 3,
//   processed: true,
//   timestamp: 1234567890
// }
```

### No Return Value

```javascript
const result5 = hello.noReturn({})
console.log(result5)
// Output: undefined
```

## Complete Example

### Go Code (hello.go)

```go
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

//export HelloWithCallback
func HelloWithCallback(params *C.char, callbackType *C.char) *C.char {
    var inputData map[string]interface{}
    json.Unmarshal([]byte(C.GoString(params)), &inputData)

    testMsg := "default"
    if msg, ok := inputData["test"].(string); ok {
        testMsg = msg
    }

    if C.GoString(callbackType) == "callback" && gCallNodeCallback != 0 {
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
```

### Node.js Code (test.js)

```javascript
const hello = require('./output/hello.node')

console.log('=== Testing go-node ===\n')

console.log('1. Test Hello1 - basic object parameter:')
const result1 = hello.hello1({ name: 'Alice', value: 21 })
console.log('Result:', result1)
console.log()

console.log('2. Test ProcessObject - object processing:')
const result2 = hello.processObject({
    name: 'Charlie',
    age: 25,
    items: ['item1', 'item2', 'item3']
})
console.log('Result:', result2)
console.log()

console.log('3. Test HelloWithCallback - synchronous callback:')
const result3 = hello.helloWithCallback({ test: 'Hello from Node' }, (data) => {
    console.log('Callback received:', data)
})
console.log('Result:', result3)
console.log()

console.log('=== All tests completed ===')
```

## Function Signature Rules

### Go Function Signature

All exported Go functions must follow this signature:

```go
func FunctionName(params *C.char, callbackType *C.char) *C.char
```

**Parameters:**
- `params`: JSON string representing an object
- `callbackType`: "callback" if a callback is provided, otherwise empty

**Return:**
- `*C.char`: JSON string representing an object, or nil

### Node.js Function Call

```javascript
const result = module.functionName(objectParam, callbackFunction)
```

**Parameters:**
- `objectParam`: JavaScript object (required)
- `callbackFunction`: JavaScript function (optional)

**Return:**
- JavaScript object or undefined

## Callback Function Rules

### Callback Parameter

Callback functions receive a single object parameter:

```javascript
module.functionName({ key: 'value' }, (callbackData) => {
    // callbackData is a JavaScript object
    console.log(callbackData.key)
})
```

### Go Callback Invocation

```go
callbackData := map[string]interface{}{
    "key":   "value",
    "result": "success",
}
jsonData, _ := json.Marshal(callbackData)
C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(jsonData)))
```

## Electron Support

To compile for Electron, use the `-ev` parameter:

```bash
go-node -input=hello.go -name=hello -ev=28.0.0
```

This will use the specified Electron version's headers instead of the default Node.js version.

### Electron Version Compatibility

Different Electron versions use different Node.js versions and ABIs. You must compile your module for the specific Electron version you're using:

| Electron Version | Node.js Version | ABI |
|-----------------|-----------------|-----|
| 39.0.0 | 22.20.0 | 140 |
| 38.0.0 | 22.12.0 | 138 |
| 37.0.0 | 22.8.0 | 136 |
| 36.0.0 | 22.6.0 | 134 |
| 35.0.0 | 22.4.0 | 132 |
| 34.0.0 | 22.2.0 | 130 |
| 33.0.0 | 22.0.0 | 128 |
| 32.0.0 | 20.16.0 | 127 |
| 31.0.0 | 20.14.0 | 125 |
| 30.0.0 | 20.12.0 | 123 |
| 29.0.0 | 20.11.0 | 121 |
| 28.0.0 | 20.9.0 | 119 |
| 27.0.0 | 18.18.0 | 118 |
| 26.0.0 | 18.16.0 | 116 |

### Electron Example

#### Go Code (electron-hello.go)

```go
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

func main() {}
```

#### Compile for Electron

```bash
go-node -input=electron-hello.go -name=hello -output=output -ev=39.0.0
```

#### Electron Code (electron-test.js)

```javascript
const { app, BrowserWindow } = require('electron');
const hello = require('./output/hello.node');

console.log('Electron:', process.versions.electron);
console.log('Node:', process.versions.node);
console.log('ABI:', process.versions.modules);

// Test basic function call
console.log('Hello1:', hello.hello1({name: 'Alice', value: 21}));
// Output: { name: 'Alice', value: 21, result: 42 }

// Test with callback
console.log('HelloWithCallback:', hello.helloWithCallback(
    {test: 'Hello from Electron'}, 
    (data) => {
        console.log('Callback received:', data);
    }
));

app.whenReady().then(() => {
    const win = new BrowserWindow({ width: 800, height: 600 });
    win.loadFile('index.html');
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});
```

#### Run Electron App

```bash
npx electron@v39.0.0 electron-test.js
```

### Electron Headers Cache

When you compile for Electron, the tool automatically downloads the Electron headers from `https://electronjs.org/headers` and caches them in:

```
C:\Users\<username>\App\Local\node-gyp\Cache\<electron-version>
```

This means subsequent compilations for the same Electron version will be faster.

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
