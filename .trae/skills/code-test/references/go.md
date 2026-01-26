# Go 语言测试规范

## 概述
本文档定义了 Go 语言的测试规范，包括单元测试、系统功能测试、系统压力测试和系统破坏性测试。

## 测试框架
使用 **GoConvey** 测试框架，它提供了 BDD 风格的测试语法和 Web 可视化界面。

### 安装 GoConvey
```bash
go get github.com/smartystreets/goconvey
```

## 单元测试

### 单元测试规范
- 单元测试文件名必须以 `_test.go` 结尾
- 单元测试函数名必须以 `Test` 开头
- 使用 GoConvey 的 BDD 风格语法
- 测试用例要覆盖正常功能、边界条件和异常输入

### 单元测试模板
```go
package main

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

// TestFunctionName 测试函数功能
func TestFunctionName(t *testing.T) {
    Convey("测试函数功能", t, func() {
        // 正常功能测试
        Convey("正常输入", func() {
            result, err := FunctionName("正常参数")
            So(err, ShouldBeNil)
            So(result, ShouldNotBeNil)
        })

        // 边界条件测试
        Convey("边界条件", func() {
            result, err := FunctionName("边界值")
            So(err, ShouldBeNil)
            So(result, ShouldEqual, "预期值")
        })

        // 异常输入测试（根据项目上下文判断）
        Convey("异常输入", func() {
            result, err := FunctionName("")
            So(err, ShouldNotBeNil)
            So(result, ShouldBeNil)
        })
    })
}
```

### 单元测试示例
```go
package user

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

// TestGetUserByID 测试根据 ID 获取用户
func TestGetUserByID(t *testing.T) {
    Convey("测试根据 ID 获取用户", t, func() {
        // 正常功能测试
        Convey("正常用户 ID", func() {
            user, err := GetUserByID(1)
            So(err, ShouldBeNil)
            So(user, ShouldNotBeNil)
            So(user.ID, ShouldEqual, 1)
        })

        // 边界条件测试
        Convey("边界值 ID", func() {
            user, err := GetUserByID(0)
            So(err, ShouldNotBeNil)
            So(user, ShouldBeNil)
        })

        // 异常输入测试（根据项目上下文，如果项目中 ID 不可能为负数，则不需要生成该测试）
        Convey("负数 ID", func() {
            user, err := GetUserByID(-1)
            So(err, ShouldNotBeNil)
            So(user, ShouldBeNil)
        })
    })
}
```

## 系统功能测试

### 系统功能测试规范
- 系统功能测试放在 `./test/system/functional/` 目录下
- 使用 GoConvey 框架
- 测试要覆盖所有单元测试
- 对于网络 API、消息队列等网络中间件，要实现接口的功能测试

### 系统功能测试模板
```go
package functional

import (
    "net/http"
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

// TestAPIFunction 测试 API 功能
func TestAPIFunction(t *testing.T) {
    Convey("测试 API 功能", t, func() {
        // API 接口功能测试
        Convey("正常请求", func() {
            resp, err := http.Get("http://localhost:8080/api/users")
            So(err, ShouldBeNil)
            So(resp.StatusCode, ShouldEqual, 200)
        })

        // 数据库操作测试
        Convey("数据库操作", func() {
            // 测试数据库操作
        })

        // 消息队列功能测试
        Convey("消息队列功能", func() {
            // 测试消息队列
        })
    })
}
```

### 系统功能测试示例
```go
package functional

import (
    "encoding/json"
    "net/http"
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

// TestUserAPI 测试用户 API 功能
func TestUserAPI(t *testing.T) {
    Convey("测试用户 API 功能", t, func() {
        baseURL := "http://localhost:8080/api"

        // 获取用户列表
        Convey("获取用户列表", func() {
            resp, err := http.Get(baseURL + "/users")
            So(err, ShouldBeNil)
            So(resp.StatusCode, ShouldEqual, 200)

            var result struct {
                Code int         `json:"code"`
                Data interface{} `json:"data"`
            }
            err = json.NewDecoder(resp.Body).Decode(&result)
            So(err, ShouldBeNil)
            So(result.Code, ShouldEqual, 0)
        })

        // 创建用户
        Convey("创建用户", func() {
            // 测试创建用户
        })

        // 更新用户
        Convey("更新用户", func() {
            // 测试更新用户
        })

        // 删除用户
        Convey("删除用户", func() {
            // 测试删除用户
        })
    })
}
```

## 系统压力测试

### 系统压力测试规范
- 系统压力测试放在 `./test/system/stress/` 目录下
- 压力测试只针对外部网络 API 生成，如 HTTP、gRPC、NSQ 等外部接口
- 压力测试要支持多梯度，从而最大程度了解接口的性能
- 测试不同并发级别下的性能表现

### 压力测试梯度
- **低并发**：1-10 并发
- **中并发**：10-100 并发
- **高并发**：100-1000 并发
- **极高并发**：1000+ 并发

### 系统压力测试模板
```go
package stress

import (
    "fmt"
    "net/http"
    "sync"
    "testing"
    "time"
    . "github.com/smartystreets/goconvey/convey"
)

// TestAPIStress 测试 API 压力
func TestAPIStress(t *testing.T) {
    Convey("测试 API 压力", t, func() {
        baseURL := "http://localhost:8080/api"

        // 低并发测试
        Convey("低并发测试 (1-10)", func() {
            testConcurrency(baseURL+"/users", 5)
        })

        // 中并发测试
        Convey("中并发测试 (10-100)", func() {
            testConcurrency(baseURL+"/users", 50)
        })

        // 高并发测试
        Convey("高并发测试 (100-1000)", func() {
            testConcurrency(baseURL+"/users", 500)
        })

        // 极高并发测试
        Convey("极高并发测试 (1000+)", func() {
            testConcurrency(baseURL+"/users", 2000)
        })
    })
}

// testConcurrency 测试并发性能
func testConcurrency(url string, concurrency int) {
    var wg sync.WaitGroup
    successCount := 0
    failCount := 0
    startTime := time.Now()

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            resp, err := http.Get(url)
            if err != nil {
                failCount++
                return
            }
            if resp.StatusCode == 200 {
                successCount++
            } else {
                failCount++
            }
            resp.Body.Close()
        }()
    }

    wg.Wait()
    duration := time.Since(startTime)

    fmt.Printf("并发数: %d, 成功: %d, 失败: %d, 耗时: %v\n",
        concurrency, successCount, failCount, duration)

    So(failCount, ShouldBeLessThan, concurrency/10)
}
```

### 系统压力测试示例
```go
package stress

import (
    "fmt"
    "net/http"
    "sync"
    "testing"
    "time"
    . "github.com/smartystreets/goconvey/convey"
)

// TestUserAPIStress 测试用户 API 压力
func TestUserAPIStress(t *testing.T) {
    Convey("测试用户 API 压力", t, func() {
        baseURL := "http://localhost:8080/api"

        // 低并发测试
        Convey("低并发测试 (1-10)", func() {
            testAPIStress(baseURL+"/users", 5)
        })

        // 中并发测试
        Convey("中并发测试 (10-100)", func() {
            testAPIStress(baseURL+"/users", 50)
        })

        // 高并发测试
        Convey("高并发测试 (100-1000)", func() {
            testAPIStress(baseURL+"/users", 500)
        })

        // 极高并发测试
        Convey("极高并发测试 (1000+)", func() {
            testAPIStress(baseURL+"/users", 2000)
        })
    })
}

// testAPIStress 测试 API 压力
func testAPIStress(url string, concurrency int) {
    var wg sync.WaitGroup
    successCount := 0
    failCount := 0
    startTime := time.Now()

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            resp, err := http.Get(url)
            if err != nil {
                failCount++
                return
            }
            if resp.StatusCode == 200 {
                successCount++
            } else {
                failCount++
            }
            resp.Body.Close()
        }()
    }

    wg.Wait()
    duration := time.Since(startTime)

    fmt.Printf("URL: %s, 并发数: %d, 成功: %d, 失败: %d, 耗时: %v\n",
        url, concurrency, successCount, failCount, duration)

    So(failCount, ShouldBeLessThan, concurrency/10)
}
```

## 系统破坏性测试

### 系统破坏性测试规范
- 系统破坏性测试放在 `./test/system/destructive/` 目录下
- 破坏性测试只针对外部网络 API 生成，如 HTTP、gRPC、NSQ 等外部接口
- 传入空异常值和 SQL 注入等异常值
- 对异常情况，如果没有出现错误则认为测试失败
- 对于错误，不是用简单的协议状态码，而是根据项目功能场景判断
- 比如：返回 HTTP 错误码是 200，但是返回的 JSON 里面 code 不为 0，也认为是返回了错误

### 系统破坏性测试模板
```go
package destructive

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

// TestAPIDestructive 测试 API 破坏性
func TestAPIDestructive(t *testing.T) {
    Convey("测试 API 破坏性", t, func() {
        baseURL := "http://localhost:8080/api"

        // 空值测试
        Convey("空值测试", func() {
            resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewReader([]byte("{}")))
            So(err, ShouldBeNil)
            So(resp.StatusCode, ShouldEqual, 200)

            var result struct {
                Code int         `json:"code"`
                Msg  string      `json:"msg"`
                Data interface{} `json:"data"`
            }
            err = json.NewDecoder(resp.Body).Decode(&result)
            So(err, ShouldBeNil)
            So(result.Code, ShouldNotEqual, 0)
        })

        // 异常类型测试
        Convey("异常类型测试", func() {
            // 测试异常类型
        })

        // SQL 注入测试
        Convey("SQL 注入测试", func() {
            // 测试 SQL 注入
        })

        // XSS 攻击测试
        Convey("XSS 攻击测试", func() {
            // 测试 XSS 攻击
        })
    })
}
```

### 系统破坏性测试示例
```go
package destructive

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

// TestUserAPIDestructive 测试用户 API 破坏性
func TestUserAPIDestructive(t *testing.T) {
    Convey("测试用户 API 破坏性", t, func() {
        baseURL := "http://localhost:8080/api"

        // 空值测试
        Convey("空值测试", func() {
            resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewReader([]byte("{}")))
            So(err, ShouldBeNil)
            So(resp.StatusCode, ShouldEqual, 200)

            var result struct {
                Code int         `json:"code"`
                Msg  string      `json:"msg"`
                Data interface{} `json:"data"`
            }
            err = json.NewDecoder(resp.Body).Decode(&result)
            So(err, ShouldBeNil)
            So(result.Code, ShouldNotEqual, 0)
        })

        // 异常类型测试
        Convey("异常类型测试", func() {
            payload := map[string]interface{}{
                "name": 123,
                "age":  "invalid",
            }
            body, _ := json.Marshal(payload)
            resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewReader(body))
            So(err, ShouldBeNil)
            So(resp.StatusCode, ShouldEqual, 200)

            var result struct {
                Code int         `json:"code"`
                Msg  string      `json:"msg"`
                Data interface{} `json:"data"`
            }
            err = json.NewDecoder(resp.Body).Decode(&result)
            So(err, ShouldBeNil)
            So(result.Code, ShouldNotEqual, 0)
        })

        // SQL 注入测试
        Convey("SQL 注入测试", func() {
            payload := map[string]interface{}{
                "name": "admin' OR '1'='1",
                "age":  25,
            }
            body, _ := json.Marshal(payload)
            resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewReader(body))
            So(err, ShouldBeNil)
            So(resp.StatusCode, ShouldEqual, 200)

            var result struct {
                Code int         `json:"code"`
                Msg  string      `json:"msg"`
                Data interface{} `json:"data"`
            }
            err = json.NewDecoder(resp.Body).Decode(&result)
            So(err, ShouldBeNil)
            So(result.Code, ShouldNotEqual, 0)
        })

        // XSS 攻击测试
        Convey("XSS 攻击测试", func() {
            payload := map[string]interface{}{
                "name": "<script>alert('xss')</script>",
                "age":  25,
            }
            body, _ := json.Marshal(payload)
            resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewReader(body))
            So(err, ShouldBeNil)
            So(resp.StatusCode, ShouldEqual, 200)

            var result struct {
                Code int         `json:"code"`
                Msg  string      `json:"msg"`
                Data interface{} `json:"data"`
            }
            err = json.NewDecoder(resp.Body).Decode(&result)
            So(err, ShouldBeNil)
            So(result.Code, ShouldNotEqual, 0)
        })
    })
}
```

## 测试可视化界面

### GoConvey Web UI
GoConvey 提供了一个强大的 Web 可视化界面，可以实时查看测试结果。

### 启动 Web UI
```bash
cd ./test
goconvey
```

### 访问 Web UI
在浏览器中打开 `http://localhost:8080`

### Web UI 功能
- 显示测试用例列表
- 显示测试执行进度
- 显示测试结果
- 显示测试覆盖率
- 支持选择运行特定测试
- 支持实时查看测试日志

## 测试目录结构
```
./test/
├── README.md              # 测试使用说明
├── run.sh                 # 测试执行脚本
├── unit/                  # 单元测试
│   └── user_test.go
├── system/                # 系统测试
│   ├── functional/        # 功能测试
│   │   └── user_api_test.go
│   ├── stress/            # 压力测试
│   │   └── user_api_stress_test.go
│   └── destructive/       # 破坏性测试
│       └── user_api_destructive_test.go
└── web/                   # 测试可视化界面
    └── server.go
```

## 测试执行脚本

### run.sh 脚本内容
```bash
#!/bin/bash

# 测试执行脚本

case "$1" in
    all)
        echo "运行所有测试..."
        go test -v ./...
        ;;
    unit)
        echo "运行单元测试..."
        go test -v ./unit/...
        ;;
    functional)
        echo "运行系统功能测试..."
        go test -v ./system/functional/...
        ;;
    stress)
        echo "运行系统压力测试..."
        go test -v ./system/stress/...
        ;;
    destructive)
        echo "运行系统破坏性测试..."
        go test -v ./system/destructive/...
        ;;
    web)
        echo "启动测试可视化界面..."
        goconvey
        ;;
    *)
        echo "使用方法: $0 {all|unit|functional|stress|destructive|web}"
        exit 1
        ;;
esac
```

## 测试使用说明

### 测试环境要求
- Go 1.16 或更高版本
- GoConvey 测试框架
- 项目依赖已安装

### 测试依赖安装
```bash
go get github.com/smartystreets/goconvey
```

### 测试命令说明
```bash
# 运行所有测试
./test/run.sh all

# 运行单元测试
./test/run.sh unit

# 运行系统功能测试
./test/run.sh functional

# 运行系统压力测试
./test/run.sh stress

# 运行系统破坏性测试
./test/run.sh destructive

# 启动测试可视化界面
./test/run.sh web
```

### 测试结果查看
- 命令行查看测试结果
- Web UI 查看测试结果和覆盖率

### 测试报告生成
```bash
# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 注意事项
- 单元测试要结合项目函数上下文，避免不必要的测试用例
- 系统压力测试要支持多梯度，全面了解接口性能
- 系统破坏性测试要根据项目功能场景判断错误，而不是简单的协议状态码
- 测试可视化界面使用 GoConvey Web UI
- 测试命令要简单易用，支持选择性执行
- 保持测试代码的简洁和可维护性
