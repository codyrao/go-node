# Go2Node - 将 Go 代码编译为 Node.js 原生模块

一个用 Go 实现的命令行工具，可以将 Go 代码编译成 Node.js 可调用的原生模块（.node 文件）。

## 特性

- **基于 CGO 和 node-gyp**：使用 Go 的 c-shared 编译模式和 node-gyp 构建原生模块
- **自动桥接代码生成**：自动生成 C++ 包装器和 binding.gyp 文件
- **回调支持**：支持从 Go 调用 Node.js 回调函数（同步和异步）
- **Electron 版本支持**：通过 `-ev` 参数指定 Electron 版本
- **固定函数签名**：所有 Go 导出函数使用一致的双参数结构
- **仅返回对象**：函数返回 0 或 1 个对象类型值
- **对象回调参数**：回调函数接收单个对象参数
- **自动回调 JSON 解析**：在回调中自动将 JSON 字符串解析为 JavaScript 对象
- **灵活的 DLL 加载**：支持从本地目录或嵌入资源加载 DLL
- **灵活配置**：支持自定义输出目录和模块名
- **临时文件管理**：使用随机 ID 自动管理临时文件以避免冲突

## 前置要求

运行工具前，确保已安装以下依赖：

### 1. Go 语言
```bash
# 下载地址：https://go.dev/dl/
# 版本要求：1.23.x
```

### 2. Node.js 和 npm
```bash
# 下载地址：https://nodejs.org/
# 版本要求：16.x 或更高
```

### 3. node-gyp（全局安装）
```bash
npm install -g node-gyp
```

### 4. Python 3.x
```bash
# Windows：https://www.python.org/downloads/
# 添加到 PATH 环境变量
```

### 5. C++ 构建工具（Windows）
- Visual Studio 2022 或更高版本
- 或安装 Build Tools for Visual Studio

## 安装

```bash
go install github.com/codyrao/go-node@latest
```

安装后，`go-node` 可执行文件将位于 `$GOPATH/bin` 目录（Windows：`%USERPROFILE%\go\bin`）。确保将此目录添加到系统 PATH 环境变量中。

## 使用方法

### 基本用法

```bash
go-node -input=your_file.go -name=module_name
```

### 参数说明

| 参数 | 必需 | 默认值 | 说明 |
|-----------|-----------|----------|-------------|
| `-input` | 是 | - | 输入 Go 源文件路径 |
| `-name` | 是 | - | 模块名称（生成 .node 文件名） |
| `-package` | 否 | `main` | Go 包名（仅编译 main 包中的函数） |
| `-output` | 否 | `./output` | 输出目录 |
| `-source` | 否 | `-input` 目录 | Go 源文件目录 |
| `-no-cleanup` | 否 | false | 编译后不清理临时文件 |
| `-ev` | 否 | - | Electron 版本（如 28.0.0）。如未指定，使用 node-gyp 的默认 Node.js 版本 |

### 使用示例

```bash
# 基本编译
go-node -input=hello.go -name=hello

# 指定输出目录
go-node -input=hello.go -name=hello -output=./dist

# 为 Electron 编译
go-node -input=hello.go -name=hello -ev=28.0.0

# 不清理临时文件（用于调试）
go-node -input=hello.go -name=hello -no-cleanup
```

### 输出文件

编译成功后，输出目录中将生成以下文件：

| 文件 | 说明 |
|------|-------------|
| `module_name.node` | Node.js 原生模块（可直接在 Node.js 中 require） |
| `module_name.dll` | Go DLL（用于独立使用或调试） |

`.node` 文件包含嵌入的 DLL，可以独立使用。`.dll` 文件也提供用于调试或直接使用。

### DLL 加载机制

工具使用灵活的 DLL 加载策略：

1. **主要方式**：尝试从与 `.node` 文件相同的目录加载 DLL
2. **备用方式**：如果未找到，尝试从 `.node` 文件的嵌入资源中提取 DLL

这确保了在各种部署场景中的可靠加载。

## Go 代码编写指南

### 函数签名要求

所有 Go 导出函数必须遵循此固定签名：

```go
func FunctionName(params *C.char, callbackType *C.char) *C.char
```

**参数结构：**
- **第一个参数（`params`）**：对象类型参数（JSON 字符串）
- **第二个参数（`callbackType`）**：回调类型指示器（"callback" 或空字符串）

**返回值：**
- **返回类型**：`*C.char`（JSON 字符串或 nil）
- **返回格式**：JSON 对象或 nil

### 导出函数

使用 `//export` 标记导出函数：

```go
package main

import "C"

//export Hello
func Hello(params *C.char, callbackType *C.char) *C.char {
    // params 是表示对象的 JSON 字符串
    // callbackType 如果提供了回调则为 "callback"
    // 返回 JSON 字符串或 nil
}

func main() {}
```

### 基本函数示例

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

### 回调函数

#### 同步回调

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

#### 异步回调

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

### 无返回值

```go
//export NoReturn
func NoReturn(params *C.char, callbackType *C.char) *C.char {
    return nil
}
```

### 对象处理示例

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

## Node.js 使用方法

### 导入模块

```javascript
const hello = require('./output/hello.node')
```

### 调用基本函数

```javascript
// 使用对象参数调用函数
const result1 = hello.hello1({ name: 'Alice', value: 21 })
console.log(result1)
// 输出：{ name: 'Alice', value: 21, result: 42 }
```

### 调用回调函数

```javascript
// 同步回调
const result2 = hello.helloWithCallback({ test: 'Hello from Node' }, (data) => {
    console.log('Callback received:', data)
    // 输出：{ test: 'Hello from Node', result: 'Callback 1' }
    // 输出：{ test: 'Hello from Node', result: 'Callback 2' }
    // 输出：{ test: 'Hello from Node', result: 'Callback 3' }
})
console.log('Result:', result2)
// 输出：{ status: 'success', result: 42 }

// 异步回调
const result3 = hello.asyncHello({ test: 'Async test' }, (data) => {
    console.log('Async callback received:', data)
    // 输出：{ test: 'Async test', result: 'Async callback 1' }
    // 输出：{ test: 'Async test', result: 'Async callback 2' }
    // ...（共 5 个回调）
})
console.log('Result:', result3)
// 输出：{ status: 'success', result: 'Async started' }
```

### 调用对象处理函数

```javascript
const result4 = hello.processObject({
    name: 'Charlie',
    age: 25,
    items: ['item1', 'item2', 'item3']
})
console.log(result4)
// 输出：{
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

### 无返回值

```javascript
const result5 = hello.noReturn({})
console.log(result5)
// 输出：undefined
```

## 完整示例

### Go 代码（hello.go）

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

### Node.js 代码（test.js）

```javascript
const hello = require('./output/hello.node')

console.log('=== 测试 go-node ===\n')

console.log('1. 测试 Hello1 - 基本对象参数：')
const result1 = hello.hello1({ name: 'Alice', value: 21 })
console.log('Result:', result1)
console.log()

console.log('2. 测试 ProcessObject - 对象处理：')
const result2 = hello.processObject({
    name: 'Charlie',
    age: 25,
    items: ['item1', 'item2', 'item3']
})
console.log('Result:', result2)
console.log()

console.log('3. 测试 HelloWithCallback - 同步回调：')
const result3 = hello.helloWithCallback({ test: 'Hello from Node' }, (data) => {
    console.log('Callback received:', data)
})
console.log('Result:', result3)
console.log()

console.log('=== 所有测试完成 ===')
```

## 函数签名规则

### Go 函数签名

所有导出的 Go 函数必须遵循此签名：

```go
func FunctionName(params *C.char, callbackType *C.char) *C.char
```

**参数：**
- `params`：表示对象的 JSON 字符串
- `callbackType`：如果提供了回调则为 "callback"，否则为空

**返回：**
- `*C.char`：表示对象的 JSON 字符串，或 nil

### Node.js 函数调用

```javascript
const result = module.functionName(objectParam, callbackFunction)
```

**参数：**
- `objectParam`：JavaScript 对象（必需）
- `callbackFunction`：JavaScript 函数（可选）

**返回：**
- JavaScript 对象或 undefined

## 回调函数规则

### 回调参数

回调函数接收单个对象参数：

```javascript
module.functionName({ key: 'value' }, (callbackData) => {
    // callbackData 是一个 JavaScript 对象
    console.log(callbackData.key)
})
```

### Go 回调调用

```go
callbackData := map[string]interface{}{
    "key":   "value",
    "result": "success",
}
jsonData, _ := json.Marshal(callbackData)
C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(jsonData)))
```

## Electron 支持

要为 Electron 编译，请使用 `-ev` 参数：

```bash
go-node -input=hello.go -name=hello -ev=28.0.0
```

这将使用指定 Electron 版本的头文件，而不是默认的 Node.js 版本。

### Electron 版本兼容性

不同的 Electron 版本使用不同的 Node.js 版本和 ABI。您必须为正在使用的特定 Electron 版本编译模块：

| Electron 版本 | Node.js 版本 | ABI |
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

### Electron 示例

#### Go 代码（electron-hello.go）

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

#### 为 Electron 编译

```bash
go-node -input=electron-hello.go -name=hello -output=output -ev=39.0.0
```

#### Electron 代码（electron-test.js）

```javascript
const { app, BrowserWindow } = require('electron');
const hello = require('./output/hello.node');

console.log('Electron:', process.versions.electron);
console.log('Node:', process.versions.node);
console.log('ABI:', process.versions.modules);

// 测试基本函数调用
console.log('Hello1:', hello.hello1({name: 'Alice', value: 21}));
// 输出：{ name: 'Alice', value: 21, result: 42 }

// 测试带回调
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

#### 运行 Electron 应用

```bash
npx electron@v39.0.0 electron-test.js
```

### Electron 头文件缓存

当您为 Electron 编译时，工具会自动从 `https://electronjs.org/headers` 下载 Electron 头文件并缓存到：

```
C:\Users\<username>\App\Local\node-gyp\Cache\<electron-version>
```

这意味着对同一 Electron 版本的后续编译会更快。

## 许可证

MIT License

## 贡献

欢迎贡献！请随时提交 Pull Request。
