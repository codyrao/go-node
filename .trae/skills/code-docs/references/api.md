# API.md 文档生成规则

## 概述
本文档定义了 API.md 的生成规则，API.md 是项目的网络接口文档，分析项目暴露的网络 API，如 HTTP、gRPC、NSQ 等，生成各个 API 的文档，说明 API 的功能、调用方式、请求参数、响应参数以及可能返回的错误码。

## 文档结构

### 1. 文档标题
- **要求**：清晰、简洁的文档标题
- **格式**：一级标题 `# API 文档`

- **示例**：
```markdown
# API 文档
```

### 2. API 分组
- **要求**：按功能模块分组 API
- **格式**：二级标题 `## 模块名`

- **示例**：
```markdown
## 用户 API

## 订单 API

## 商品 API
```

### 3. API 接口文档

#### HTTP API 文档格式
```markdown
### 接口名称

**功能**：接口功能描述

**调用方式**：HTTP 方法 URL

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| param1 | string | 是 | 参数说明 |
| param2 | int | 否 | 参数说明 |

**请求示例**：

```json
{
  "param1": "value1",
  "param2": 123
}
```

**响应参数**：

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码 |
| msg | string | 消息 |
| data | object | 数据 |

**响应示例**：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "test"
  }
}
```

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 用户不存在 |
```

#### gRPC API 文档格式
```markdown
### 接口名称

**功能**：接口功能描述

**调用方式**：gRPC 方法

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| param1 | string | 是 | 参数说明 |
| param2 | int32 | 否 | 参数说明 |

**请求示例**：

```protobuf
message Request {
  string param1 = 1;
  int32 param2 = 2;
}
```

**响应参数**：

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int32 | 状态码 |
| msg | string | 消息 |
| data | object | 数据 |

**响应示例**：

```protobuf
message Response {
  int32 code = 1;
  string msg = 2;
  Data data = 3;
}
```

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 用户不存在 |
```

#### NSQ API 文档格式
```markdown
### 接口名称

**功能**：接口功能描述

**调用方式**：NSQ Topic

**消息格式**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| param1 | string | 是 | 参数说明 |
| param2 | int | 否 | 参数说明 |

**消息示例**：

```json
{
  "param1": "value1",
  "param2": 123
}
```

**处理说明**：

说明消息的处理逻辑和注意事项。

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 处理失败 |
```

### 4. 分析项目暴露的网络 API

#### HTTP API
分析以下内容：
- 路由定义（如 Go 的 `http.HandleFunc`、`gin.Engine` 等）
- 控制器方法（如 Spring 的 `@RequestMapping` 等）
- API 注解（如 Swagger 的 `@ApiOperation` 等）

**示例分析（Go + Gin）**：
```go
r.GET("/users", GetUsers)
r.POST("/users", CreateUser)
r.PUT("/users/:id", UpdateUser)
r.DELETE("/users/:id", DeleteUser)
```

**生成的 API 文档**：
```markdown
## 用户 API

### 获取用户列表

**功能**：获取用户列表

**调用方式**：GET /api/users

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**响应参数**：

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码 |
| msg | string | 消息 |
| data | object | 数据 |

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
```

#### gRPC API
分析以下内容：
- Proto 文件定义
- gRPC 服务定义
- gRPC 方法定义

**示例分析（Proto）**：
```protobuf
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}
```

**生成的 API 文档**：
```markdown
## 用户 API

### 获取用户信息

**功能**：获取用户信息

**调用方式**：gRPC GetUser

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int64 | 是 | 用户 ID |

**响应参数**：

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int32 | 状态码 |
| msg | string | 消息 |
| data | object | 数据 |

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1002 | 用户不存在 |
```

#### NSQ API
分析以下内容：
- NSQ Topic 定义
- NSQ 消息处理函数
- NSQ 消息格式

**示例分析（Go）**：
```go
consumer, _ := nsq.NewConsumer("users", "channel", nsq.NewConfig())
consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
    // 处理消息
    return nil
}))
```

**生成的 API 文档**：
```markdown
## 用户 API

### 用户消息

**功能**：处理用户相关消息

**调用方式**：NSQ Topic: users

**消息格式**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| action | string | 是 | 操作类型：create/update/delete |
| id | int64 | 否 | 用户 ID |
| data | object | 否 | 用户数据 |

**处理说明**：

根据 action 类型执行相应的操作：
- create：创建用户
- update：更新用户
- delete：删除用户

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 处理失败 |
```

### 5. 参考 Swagger 文档

#### 如果项目中有 Swagger 文档
- 参考 Swagger 文档生成 API.md
- 提取 Swagger 中的 API 定义
- 转换为 Markdown 格式

#### Swagger 注解示例
```go
// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} Response
// @Router /users [get]
func GetUsers(c *gin.Context) {
    // 实现
}
```

**生成的 API 文档**：
```markdown
### 获取用户列表

**功能**：获取用户列表

**调用方式**：GET /api/users

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**响应参数**：

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码 |
| msg | string | 消息 |
| data | object | 数据 |

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
```

### 6. 错误码说明

#### 通用错误码
```markdown
## 通用错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 系统错误 |
| 1003 | 未授权 |
| 1004 | 权限不足 |
```

#### 业务错误码
```markdown
## 业务错误码

| 错误码 | 说明 |
|--------|------|
| 2001 | 用户不存在 |
| 2002 | 用户已存在 |
| 2003 | 密码错误 |
```

### 7. 文档模板

#### 完整模板
```markdown
# API 文档

## 通用错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 系统错误 |
| 1003 | 未授权 |
| 1004 | 权限不足 |

## 用户 API

### 获取用户列表

**功能**：获取用户列表

**调用方式**：GET /api/users

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**请求示例**：

```json
{
  "page": 1,
  "size": 10
}
```

**响应参数**：

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码 |
| msg | string | 消息 |
| data | object | 数据 |

**响应示例**：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 100,
    "list": [
      {
        "id": 1,
        "name": "test"
      }
    ]
  }
}
```

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
```

## 订单 API

### 创建订单

**功能**：创建订单

**调用方式**：POST /api/orders

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | int64 | 是 | 用户 ID |
| productId | int64 | 是 | 商品 ID |
| quantity | int | 是 | 数量 |

**请求示例**：

```json
{
  "userId": 1,
  "productId": 100,
  "quantity": 2
}
```

**响应参数**：

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码 |
| msg | string | 消息 |
| data | object | 数据 |

**响应示例**：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "orderId": 1001,
    "status": "pending"
  }
}
```

**错误码**：

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 3001 | 商品不存在 |
| 3002 | 库存不足 |
```
```

## 文档放置位置
- 放在项目当前路径下 `./API.md`

## 适用范围
- 只针对服务端和相关 SDK 的代码文档生成
- 对于客户端和前端，API.md 不需要生成

## 注意事项
- 分析项目暴露的网络 API（HTTP、gRPC、NSQ 等）
- 参考 Swagger 文档（如果存在）
- 生成各个 API 的文档，说明 API 功能、调用方式、请求参数、响应参数
- 说明可能返回的错误码
- 保持 API 文档的清晰和详细
- 及时更新 API 文档，反映最新的 API 变更
- 使用表格格式展示参数和错误码，提高可读性
- 提供请求和响应示例，方便开发者理解
