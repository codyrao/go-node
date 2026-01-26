# Go2Node - 将Go代码编译为Node.js原生模块

一个用Go实现的命令行工具，可以将Go代码编译成Node.js可调用的原生模块（.node文件）。

## 功能特性

- **基于CGO和node-gyp**：使用Go的c-shared编译模式和node-gyp构建原生模块
- **自动生成桥接代码**：自动生成C++ wrapper和binding.gyp文件
- **回调支持**：支持从Go调用Node.js回调函数（同步和异步）
- **对象返回支持**：支持返回JSON格式的嵌套JavaScript对象
- **灵活的配置**：支持自定义输出目录和模块名称
- **临时文件管理**：自动管理临时文件，使用随机ID避免冲突

## 前置依赖

在运行工具之前，请确保安装以下依赖：

### 1. Go语言
```bash
# 下载地址: https://go.dev/dl/
# 版本要求: 1.23.x
```

### 2. Node.js和npm
```bash
# 下载地址: https://nodejs.org/
# 版本要求: 16.x 或更高
```

### 3. node-gyp（全局安装）
```bash
npm install -g node-gyp
```

### 4. Python 3.x
```bash
# Windows: https://www.python.org/downloads/
# 添加到PATH环境变量
```

### 5. C++编译工具（Windows）
- Visual Studio 2022 或更高版本
- 或安装 Build Tools for Visual Studio

## 安装

```bash
go install github.com/codyrao/go-node@latest
```

安装完成后，`go-node` 可执行文件将位于 `$GOPATH/bin` 目录下（Windows下为 `%USERPROFILE%\go\bin`）。确保该目录已添加到系统PATH环境变量中。

## 使用方法

### 基本用法

```bash
go-node -input=your_file.go -name=module_name
```

### 参数说明

| 参数 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `-input` | 是 | - | 输入的Go源文件路径 |
| `-name` | 是 | - | 模块名称（生成.node和.dll文件名） |
| `-package` | 否 | `main` | Go包名（仅main包下的函数会被编译） |
| `-output` | 否 | `./output` | 输出目录 |
| `-source` | 否 | `-input`所在目录 | Go源文件所在目录 |
| `-no-cleanup` | 否 | false | 编译后不清理临时文件 |

### 示例

```bash
# 基本编译
go-node -input=hello.go -name=hello

# 指定输出目录
go-node -input=hello.go -name=hello -output=./dist

# 不清理临时文件（用于调试）
go-node -input=hello.go -name=hello -no-cleanup
```

## Go代码编写规范

### 导出函数

使用 `//export` 标记导出函数：

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

### 支持的函数签名

工具支持多种函数签名模式：

#### 1. 简单同步函数

```go
//export Hello1
func Hello1(name *C.char, value *C.char) *C.char {
    // name和value都是字符串参数
    // 返回JSON字符串
    result := map[string]interface{}{
        "name":   C.GoString(name),
        "value":  C.GoString(value),
        "result": "success",
    }
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}
```

#### 2. 同步回调函数

```go
//export Hello4
func Hello4(name *C.char, callbackType *C.char) *C.char {
    if C.GoString(callbackType) == "callback" {
        // 同步触发回调
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 1"}`))
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 2"}`))
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 3"}`))
    }
    return C.CString(`{"result":42,"status":"success"}`)
}
```

#### 3. 异步回调函数

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

#### 4. 无限异步回调函数

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

### 回调支持

工具内置了回调机制，Go代码可以调用导出的`CallCallback`函数：

```go
//export ProcessWithCallback
func ProcessWithCallback() {
    // Go代码中调用回调
    response := `{"status": "success"}`
    CallCallback(C.CString("callback_type"), C.CString(response))
}
```

### 返回对象

Go函数可以返回JSON字符串，Node.js端解析为对象：

```go
//export GetUserInfo
func GetUserInfo() *C.char {
    user := map[string]interface{}{
        "name":    "John",
        "age":     30,
        "address": map[string]string{
            "city": "New York",
        },
    }
    jsonBytes, _ := json.Marshal(user)
    return C.CString(string(jsonBytes))
}
```

Node.js端使用：
```javascript
const user = hello.getUserInfo();
console.log(user.name);      // "John"
console.log(user.address.city); // "New York"
```

## Node.js端使用

### 引入模块

```javascript
const demoaddon = require('./output/hello.node')
```

### 调用同步函数

```javascript
// 简单调用
const result = demoaddon.Hello1("Test1", "10")
console.log(result)  // JSON字符串
```

### 调用回调函数

```javascript
// 同步回调
demoaddon.Hello4('{"test":"Hello4 Test"}', function(callbackType, jsonData) {
    console.log(`回调 [${callbackType}]:`, jsonData)
})

// 异步回调
demoaddon.Hello5('{"test":"Hello5 Test"}', function(callbackType, jsonData) {
    console.log(`异步回调 [${callbackType}]:`, jsonData)
})
```

## 完整示例

### Go代码 (hello.go)

```go
package main

import "C"
import (
    "encoding/json"
    "fmt"
    "time"
)

//export Hello1
func Hello1(name *C.char, value *C.char) *C.char {
    nameStr := C.GoString(name)
    valueStr := C.GoString(value)
    
    result := map[string]interface{}{
        "name":   nameStr,
        "value":  valueStr,
        "result": "Hello1 success",
    }
    jsonBytes, _ := json.Marshal(result)
    return C.CString(string(jsonBytes))
}

//export Hello4
func Hello4(name *C.char, callbackType *C.char) *C.char {
    if C.GoString(callbackType) == "callback" {
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 1"}`))
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 2"}`))
        CallCallback(C.CString("sync_callback"), C.CString(`{"result":"Callback 3"}`))
    }
    return C.CString(`{"result":42,"status":"success"}`)
}

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

func main() {}
```

### Node.js代码 (test.js)

```javascript
const demoaddon = require('./output/hello.node')

console.log('=== Go2Node 测试 ===\n')

// 测试1: 简单调用
console.log('测试1: Hello1')
const result1 = demoaddon.Hello1("Test1", "10")
console.log('   结果:', result1)
console.log()

// 测试2: 同步回调
console.log('测试2: Hello4 - 同步回调')
demoaddon.Hello4('{"test":"Hello4 Test"}', function(callbackType, jsonData) {
    console.log(`   回调 [${callbackType}]:`, jsonData)
})
console.log()

// 测试3: 异步回调
console.log('测试3: Hello5 - 异步回调')
demoaddon.Hello5('{"test":"Hello5 Test"}', function(callbackType, jsonData) {
    console.log(`   异步回调 [${callbackType}]:`, jsonData)
})

console.log('等待异步回调...')
```

### 编译和运行

```bash
# 编译
go-node -input=hello.go -name=hello

# 运行测试
node test.js
```

## 输出结构

编译完成后，输出目录结构如下：

```
output/
└── module_name.node    # Node.js原生模块（已嵌入Go动态库）
```

**注意**：Go动态库已嵌入到.node文件中，运行时会自动提取到系统临时目录使用。

## 临时文件管理

工具会自动管理临时文件：

1. 在源文件所在目录创建 `go-node-tmp_{随机ID}` 临时目录
2. 在临时目录中生成 `binding.gyp` 和 `wrapper.cc` 等文件
3. 编译完成后自动清理临时目录
4. 使用 `-no-cleanup` 参数可以保留临时文件用于调试

## 测试

查看 `test/hello.go` 了解完整示例代码，运行测试：

```bash
# 编译示例
go-node -input=test/hello.go -name=hello

# 运行测试
node test.js
```

## 已知限制

1. **包名限制**：仅支持main包下的函数编译
2. **平台相关**：当前实现主要针对Windows平台
3. **回调上下文**：异步回调使用libuv确保在正确的V8上下文中执行
4. **参数类型**：函数参数统一使用字符串传递，复杂类型使用JSON序列化

## 技术实现

- **Go编译**：使用 `go build -buildmode=c-shared` 生成动态库
- **Node.js绑定**：使用node-gyp生成C++原生模块
- **跨DLL通信**：使用Windows API动态加载和函数指针
- **异步回调**：使用libuv的uv_async_t实现线程安全的异步回调
- **回调队列**：使用互斥锁保护的队列管理回调请求

## 许可证

MIT
