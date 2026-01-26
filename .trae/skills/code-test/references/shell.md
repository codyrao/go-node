# Shell 脚本测试规范

## 概述
本文档定义了 Shell 脚本的测试规范，包括单元测试和系统功能测试。

## 测试框架
使用 **Bats (Bash Automated Testing System)** 测试框架，它提供了简洁的 Shell 脚本测试语法。

### 安装 Bats
```bash
git clone https://github.com/sstephenson/bats.git
cd bats
sudo ./install.sh /usr/local
```

## 单元测试

### 单元测试规范
- 单元测试文件名必须以 `.bats` 结尾
- 使用 Bats 的测试语法
- 测试用例要覆盖正常功能、边界条件和异常输入

### 单元测试模板
```bash
#!/usr/bin/env bats

# 加载被测试的脚本
source ./script.sh

@test "测试函数正常功能" {
    # 正常功能测试
    result=$(function_name "正常参数")
    [ "$result" = "预期值" ]
}

@test "测试函数边界条件" {
    # 边界条件测试
    result=$(function_name "边界值")
    [ "$result" = "预期值" ]
}

@test "测试函数异常输入" {
    # 异常输入测试（根据项目上下文判断）
    run function_name ""
    [ "$status" -ne 0 ]
}
```

### 单元测试示例
```bash
#!/usr/bin/env bats

# 加载被测试的脚本
source ./user.sh

@test "测试获取用户信息" {
    # 正常功能测试
    result=$(get_user_info 1)
    [ "$result" = "用户信息" ]
}

@test "测试获取用户信息边界条件" {
    # 边界条件测试
    run get_user_info 0
    [ "$status" -ne 0 ]
}

@test "测试获取用户信息异常输入" {
    # 异常输入测试（根据项目上下文，如果项目中 ID 不可能为负数，则不需要测试）
    # run get_user_info -1
    # [ "$status" -ne 0 ]
}
```

## 系统功能测试

### 系统功能测试规范
- 系统功能测试放在 `./test/system/functional/` 目录下
- 使用 Bats 测试框架
- 测试要覆盖所有单元测试
- 对于网络 API、消息队列等网络中间件，要实现接口的功能测试

### 系统功能测试模板
```bash
#!/usr/bin/env bats

@test "测试 API 功能" {
    # API 接口功能测试
    result=$(curl -s http://localhost:8080/api/users)
    echo "$result" | grep -q '"code":0'
}

@test "测试数据库操作" {
    # 数据库操作测试
    result=$(mysql -u root -p123456 -e "SELECT * FROM users WHERE id=1;")
    echo "$result" | grep -q "用户信息"
}

@test "测试消息队列功能" {
    # 消息队列功能测试
    result=$(nsq_stat --lookupd-http-address=127.0.0.1:4161)
    echo "$result" | grep -q "topic"
}
```

### 系统功能测试示例
```bash
#!/usr/bin/env bats

@test "测试用户 API 功能" {
    baseURL="http://localhost:8080/api"

    # 获取用户列表
    @test "获取用户列表" {
        result=$(curl -s "$baseURL/users")
        echo "$result" | grep -q '"code":0'
    }

    # 创建用户
    @test "创建用户" {
        result=$(curl -s -X POST -H "Content-Type: application/json" \
            -d '{"name":"test","age":25}' \
            "$baseURL/users")
        echo "$result" | grep -q '"code":0'
    }

    # 更新用户
    @test "更新用户" {
        result=$(curl -s -X PUT -H "Content-Type: application/json" \
            -d '{"name":"test2","age":26}' \
            "$baseURL/users/1")
        echo "$result" | grep -q '"code":0'
    }

    # 删除用户
    @test "删除用户" {
        result=$(curl -s -X DELETE "$baseURL/users/1")
        echo "$result" | grep -q '"code":0'
    }
}
```

## 测试可视化界面

### 使用 Bats-Reporter
Bats-Reporter 是一个用于生成 HTML 测试报告的工具。

### 安装 Bats-Reporter
```bash
npm install -g bats-reporter
```

### 生成测试报告
```bash
bats --reporter bats-reporter test/*.bats > report.html
```

### 访问测试报告
在浏览器中打开 `report.html`

## 测试目录结构
```
./test/
├── README.md              # 测试使用说明
├── run.sh                 # 测试执行脚本
├── unit/                  # 单元测试
│   └── user_test.bats
├── system/                # 系统测试
│   └── functional/        # 功能测试
│       └── user_api_test.bats
└── report/                # 测试报告
    └── report.html
```

## 测试执行脚本

### run.sh 脚本内容
```bash
#!/bin/bash

# 测试执行脚本

case "$1" in
    all)
        echo "运行所有测试..."
        bats --reporter bats-reporter ./unit/*.bats ./system/functional/*.bats > ./report/report.html
        ;;
    unit)
        echo "运行单元测试..."
        bats --reporter bats-reporter ./unit/*.bats > ./report/report.html
        ;;
    functional)
        echo "运行系统功能测试..."
        bats --reporter bats-reporter ./system/functional/*.bats > ./report/report.html
        ;;
    web)
        echo "启动测试可视化界面..."
        if [ -f ./report/report.html ]; then
            echo "测试报告已生成: ./report/report.html"
            echo "请在浏览器中打开该文件查看测试结果"
        else
            echo "测试报告不存在，请先运行测试"
        fi
        ;;
    *)
        echo "使用方法: $0 {all|unit|functional|web}"
        exit 1
        ;;
esac
```

## 测试使用说明

### 测试环境要求
- Bash 4.0 或更高版本
- Bats 测试框架
- curl（用于 API 测试）

### 测试依赖安装
```bash
# 安装 Bats
git clone https://github.com/sstephenson/bats.git
cd bats
sudo ./install.sh /usr/local

# 安装 Bats-Reporter
npm install -g bats-reporter
```

### 测试命令说明
```bash
# 运行所有测试
./test/run.sh all

# 运行单元测试
./test/run.sh unit

# 运行系统功能测试
./test/run.sh functional

# 启动测试可视化界面
./test/run.sh web
```

### 测试结果查看
- 命令行查看测试结果
- HTML 报告查看测试结果

### 测试报告生成
```bash
# 生成测试报告
bats --reporter bats-reporter test/*.bats > report.html
```

## 注意事项
- 单元测试要结合项目函数上下文，避免不必要的测试用例
- 测试可视化界面使用 Bats-Reporter 生成 HTML 报告
- 测试命令要简单易用，支持选择性执行
- 保持测试代码的简洁和可维护性
