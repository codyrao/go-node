# Go 语言代码架构规范

## 概述
本文档定义了 Go 语言的代码架构规范，遵循 Go 语言官方的最佳实践和社区约定。

## 项目结构

### 标准项目结构
```
project-name/
├── api/                    # 对外保留的接口
│   ├── http/              # HTTP API 定义
│   ├── grpc/              # gRPC API 定义
│   └── openapi/           # OpenAPI/Swagger 定义
├── bin/                   # 编译产物
│   ├── project-name       # 可执行文件
│   └── project-name.exe   # Windows 可执行文件
├── cmd/                   # 项目代码的执行入口
│   ├── server/            # 服务器入口
│   │   └── main.go
│   ├── client/            # 客户端入口
│   │   └── main.go
│   └── worker/            # 工作进程入口
│       └── main.go
├── config/                # 项目代码和相关配置
│   ├── config.go          # 配置结构定义
│   ├── config.yaml        # 配置文件
│   ├── config.dev.yaml    # 开发环境配置
│   ├── config.test.yaml   # 测试环境配置
│   └── config.prod.yaml   # 生产环境配置
├── example/               # 项目代码的测试 Demo
│   ├── basic/             # 基础示例
│   ├── advanced/          # 高级示例
│   └── README.md
├── internal/              # 项目代码的核心逻辑
│   ├── handler/           # HTTP 处理器
│   ├── service/           # 业务逻辑服务
│   ├── repository/        # 数据访问层
│   ├── model/             # 数据模型
│   ├── middleware/        # 中间件
│   ├── pkg/               # 对第三方的库的相关引用
│   │   ├── redis/         # Redis 封装
│   │   ├── mysql/         # MySQL 封装
│   │   ├── nsq/           # NSQ 封装
│   │   └── utils/         # 工具函数
│   └── const/             # 常量定义
├── scripts/               # 项目代码的辅助脚本
│   ├── build.sh           # 构建脚本
│   ├── deploy.sh          # 部署脚本
│   ├── test.sh            # 测试脚本
│   └── lint.sh            # 代码检查脚本
├── test/                  # 项目代码的系统测试入口
│   ├── integration/       # 集成测试
│   ├── e2e/               # 端到端测试
│   ├── stress/            # 压力测试
│   └── destructive/       # 破坏性测试
├── go.mod                 # Go 模块定义
├── go.sum                 # Go 模块校验和
├── Makefile               # Makefile 构建文件，统筹项目所有可执行的命令
├── Dockerfile             # Docker 构建文件
├── docker-compose.yml     # Docker Compose 配置
├── .gitignore             # Git 忽略文件
└── README.md              # 项目说明文档
```

## 目录说明

### api/
**功能**：存放对外保留的接口定义

**子目录**：
- `http/`：HTTP API 定义，包括路由、请求/响应结构
- `grpc/`：gRPC API 定义，包括 proto 文件
- `openapi/`：OpenAPI/Swagger 定义

**示例**：
```
api/
├── http/
│   ├── user.go            # 用户 API 定义
│   └── product.go         # 产品 API 定义
├── grpc/
│   ├── user.proto         # 用户 gRPC 定义
│   └── product.proto      # 产品 gRPC 定义
└── openapi/
    └── swagger.yaml       # OpenAPI 规范
```

### bin/
**功能**：存放编译产物

**说明**：
- 可执行文件直接放在 bin 目录下
- 文件名通常与项目名相同
- Windows 平台添加 .exe 后缀

**示例**：
```
bin/
├── myproject              # Linux/Mac 可执行文件
└── myproject.exe          # Windows 可执行文件
```

### cmd/
**功能**：存放项目代码的执行入口

**说明**：
- 每个子目录代表一个可执行程序
- main.go 文件作为入口点
- 保持简洁，主要逻辑放在 internal 中

**示例**：
```
cmd/
├── server/
│   └── main.go            # 服务器入口
├── client/
│   └── main.go            # 客户端入口
└── worker/
    └── main.go            # 工作进程入口
```

**main.go 示例**：
```go
package main

import (
    "log"
    "myproject/internal/server"
)

func main() {
    if err := server.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### config/
**功能**：存放项目代码和相关配置

**说明**：
- config.go 定义配置结构
- config.yaml 为默认配置
- 不同环境使用不同的配置文件

**示例**：
```
config/
├── config.go              # 配置结构定义
├── config.yaml            # 默认配置
├── config.dev.yaml        # 开发环境配置
├── config.test.yaml       # 测试环境配置
└── config.prod.yaml       # 生产环境配置
```

**config.go 示例**：
```go
package config

import (
    "os"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
    Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Database string `mapstructure:"database"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

func Load(configPath string) (*Config, error) {
    viper.SetConfigFile(configPath)
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

### example/
**功能**：存放项目代码的测试 Demo

**说明**：
- 提供使用示例
- 按难度或功能分类
- 包含 README 说明

**示例**：
```
example/
├── basic/
│   ├── user_create.go    # 创建用户示例
│   └── user_query.go     # 查询用户示例
├── advanced/
│   ├── batch_process.go  # 批量处理示例
│   └── concurrent.go     # 并发处理示例
└── README.md             # 示例说明
```

### internal/
**功能**：存放项目代码的核心逻辑

**说明**：
- internal 目录下的代码不能被外部项目导入
- 包含业务逻辑、数据访问、中间件等

#### internal/handler/
**功能**：HTTP 处理器

**说明**：
- 处理 HTTP 请求
- 调用 service 层
- 返回 HTTP 响应

**示例**：
```
internal/handler/
├── user.go               # 用户处理器
├── product.go            # 产品处理器
└── middleware.go        # 中间件处理器
```

#### internal/service/
**功能**：业务逻辑服务

**说明**：
- 实现业务逻辑
- 调用 repository 层
- 处理业务规则

**示例**：
```
internal/service/
├── user.go               # 用户服务
├── product.go            # 产品服务
└── order.go              # 订单服务
```

#### internal/repository/
**功能**：数据访问层

**说明**：
- 封装数据库操作
- 提供 CRUD 接口
- 处理数据持久化

**示例**：
```
internal/repository/
├── user.go               # 用户数据访问
├── product.go            # 产品数据访问
└── order.go              # 订单数据访问
```

#### internal/model/
**功能**：数据模型

**说明**：
- 定义数据结构
- 数据库表映射
- API 请求/响应结构

**示例**：
```
internal/model/
├── user.go               # 用户模型
├── product.go            # 产品模型
└── order.go              # 订单模型
```

#### internal/middleware/
**功能**：中间件

**说明**：
- 认证中间件
- 日志中间件
- 限流中间件

**示例**：
```
internal/middleware/
├── auth.go               # 认证中间件
├── logging.go            # 日志中间件
└── rate_limit.go         # 限流中间件
```

#### internal/pkg/
**功能**：对第三方的库的相关引用

**说明**：
- 封装第三方库
- 提供统一的接口
- 便于替换和测试

**示例**：
```
internal/pkg/
├── redis/                # Redis 封装
│   ├── redis.go
│   └── redis_test.go
├── mysql/                # MySQL 封装
│   ├── mysql.go
│   └── mysql_test.go
├── nsq/                  # NSQ 封装
│   ├── nsq.go
│   └── nsq_test.go
└── utils/                # 工具函数
    ├── string.go
    ├── time.go
    └── crypto.go
```

#### internal/const/
**功能**：常量定义

**说明**：
- 定义项目常量
- 避免魔法数字
- 统一管理常量

**示例**：
```
internal/const/
├── error.go              # 错误码定义
├── status.go             # 状态码定义
└── config.go             # 配置常量
```

### scripts/
**功能**：存放项目代码的辅助脚本

**说明**：
- 构建脚本
- 部署脚本
- 测试脚本
- 代码检查脚本

**示例**：
```
scripts/
├── build.sh              # 构建脚本
├── deploy.sh             # 部署脚本
├── test.sh               # 测试脚本
└── lint.sh               # 代码检查脚本
```

**build.sh 示例**：
```bash
#!/bin/bash

set -e

echo "Building project..."

# 设置构建参数
VERSION=${VERSION:-"1.0.0"}
BUILD_TIME=$(date +%Y-%m-%d\ %H:%M:%S)
GIT_COMMIT=$(git rev-parse --short HEAD)

# 构建
go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" -o bin/myproject cmd/server/main.go

echo "Build completed: bin/myproject"
```

### test/
**功能**：存放项目代码的系统测试入口

**说明**：
- 集成测试
- 端到端测试
- 压力测试
- 破坏性测试

**示例**：
```
test/
├── integration/          # 集成测试
│   ├── user_test.go
│   └── product_test.go
├── e2e/                  # 端到端测试
│   ├── user_flow_test.go
│   └── order_flow_test.go
├── stress/               # 压力测试
│   ├── api_stress_test.go
│   └── grpc_stress_test.go
├── destructive/          # 破坏性测试
│   ├── api_destructive_test.go
│   └── grpc_destructive_test.go
└── README.md             # 测试说明
```

### Makefile
**功能**：统筹项目所有可执行的命令

**说明**：
- Makefile 应放在项目根目录
- 统筹项目中所有目录下的可执行命令或脚本
- 执行 `make {命令}` 就能执行项目中所有目录下的可执行命令或脚本
- 提供统一的命令入口，简化项目操作

**常用命令**：
- `make build`：构建项目
- `make run`：运行项目
- `make test`：运行测试
- `make lint`：代码检查
- `make clean`：清理构建产物
- `make deploy`：部署项目
- `make docker-build`：构建 Docker 镜像
- `make docker-run`：运行 Docker 容器

**Makefile 示例**：
```makefile
# Makefile 统筹项目所有可执行的命令

.PHONY: build run test lint clean deploy docker-build docker-run help

# 项目名称
PROJECT_NAME := myproject

# Go 编译参数
GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# 版本信息
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# 构建目录
BUILD_DIR := bin
CMD_DIR := cmd

# 默认目标
all: build

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make build        - 构建项目"
	@echo "  make run          - 运行项目"
	@echo "  make test         - 运行测试"
	@echo "  make lint         - 代码检查"
	@echo "  make clean        - 清理构建产物"
	@echo "  make deploy       - 部署项目"
	@echo "  make docker-build  - 构建 Docker 镜像"
	@echo "  make docker-run   - 运行 Docker 容器"
	@echo "  make help         - 显示帮助信息"

# 构建项目
build:
	@echo "构建项目..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) $(CMD_DIR)/server/main.go
	@echo "构建完成: $(BUILD_DIR)/$(PROJECT_NAME)"

# 运行项目
run:
	@echo "运行项目..."
	$(GO) run $(CMD_DIR)/server/main.go

# 运行测试
test:
	@echo "运行测试..."
	$(GO) test -v ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 清理构建产物
clean:
	@echo "清理构建产物..."
	rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 部署项目
deploy:
	@echo "部署项目..."
	./scripts/deploy.sh

# 构建 Docker 镜像
docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t $(PROJECT_NAME):$(VERSION) .
	docker tag $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest
	@echo "Docker 镜像构建完成"

# 运行 Docker 容器
docker-run:
	@echo "运行 Docker 容器..."
	docker-compose up -d
	@echo "Docker 容器运行完成"
```

**Makefile 规范**：
- 使用 `.PHONY` 声明伪目标
- 使用变量定义可配置的参数
- 提供帮助信息（`make help`）
- 命令名称使用小写，单词间用连字符分隔
- 每个命令前添加注释说明
- 使用 `@echo` 输出执行信息
- 使用 `@` 静默执行命令

**使用示例**：
```bash
# 构建项目
make build

# 运行项目
make run

# 运行测试
make test

# 代码检查
make lint

# 清理构建产物
make clean

# 部署项目
make deploy

# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run

# 查看帮助
make help
```

## 命名规范

### 包命名
- 包名使用小写单词
- 不使用下划线或驼峰
- 包名应简洁明了
- 包名与目录名一致

**示例**：
```
user/          # 正确
userService/   # 错误
user_service/  # 错误
```

### 文件命名
- 文件名使用小写单词
- 多个单词使用下划线分隔
- 测试文件以 `_test.go` 结尾

**示例**：
```
user.go              # 正确
user_service.go      # 正确
userService.go       # 错误
user_test.go         # 正确
```

### 接口命名
- 接口名以 `er` 结尾
- 或使用 `I` 前缀

**示例**：
```go
type Userer interface {     # 正确
    GetUser(id int) (*User, error)
}

type IUser interface {     # 正确
    GetUser(id int) (*User, error)
}
```

### 函数命名
- 函数名使用驼峰命名
- 导出函数首字母大写
- 私有函数首字母小写

**示例**：
```go
func GetUser(id int) (*User, error) {    # 导出函数
    return &User{ID: id}, nil
}

func getUser(id int) (*User, error) {    # 私有函数
    return &User{ID: id}, nil
}
```

### 变量命名
- 变量名使用驼峰命名
- 导出变量首字母大写
- 私有变量首字母小写
- 常量使用全大写，单词间用下划线分隔

**示例**：
```go
var UserName string           # 导出变量
var userName string           # 私有变量
const MAX_RETRY_COUNT = 3     # 常量
```

## 最佳实践

### 1. 错误处理
- 总是检查错误
- 使用 errors.Wrap 包装错误
- 提供有意义的错误信息

**示例**：
```go
func GetUser(id int) (*User, error) {
    user, err := repository.GetUser(id)
    if err != nil {
        return nil, errors.Wrap(err, "failed to get user")
    }
    return user, nil
}
```

### 2. 依赖注入
- 使用接口定义依赖
- 通过构造函数注入依赖
- 便于测试和替换

**示例**：
```go
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

### 3. 并发安全
- 使用 sync 包处理并发
- 使用 channel 进行通信
- 避免共享状态

**示例**：
```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}
```

### 4. 日志记录
- 使用结构化日志
- 记录关键操作
- 避免记录敏感信息

**示例**：
```go
log.WithFields(log.Fields{
    "user_id": userID,
    "action":  "create",
}).Info("User created successfully")
```

### 5. 配置管理
- 使用配置文件
- 支持多环境配置
- 敏感信息使用环境变量

**示例**：
```go
config, err := config.Load("config/config.yaml")
if err != nil {
    log.Fatal(err)
}

// 使用环境变量覆盖配置
if os.Getenv("DB_PASSWORD") != "" {
    config.Database.Password = os.Getenv("DB_PASSWORD")
}
```

## 注意事项

1. **保持简洁**：代码应该简洁明了，避免过度设计
2. **遵循规范**：遵循 Go 语言官方规范和社区约定
3. **测试覆盖**：保持高测试覆盖率
4. **文档完整**：为公共 API 提供文档
5. **性能优化**：在必要时进行性能优化
6. **安全考虑**：注意代码安全性，避免常见漏洞
7. **版本管理**：使用语义化版本号
8. **依赖管理**：合理管理第三方依赖
