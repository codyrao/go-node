---
name: code-docs
description: 当用户需要为生成的代码或代码改动生成相关文档时激活此技能。根据代码语言和项目上下文，生成 README.md、Dependency.mmd、API.md 和 Page.md 等文档。
---

# Code Docs - 代码文档生成

## 激活条件
当用户需要进行以下操作时激活此技能：
- 为新生成的代码生成文档
- 为代码改动更新文档
- 补充现有代码的文档
- 生成项目使用文档
- 生成项目依赖图
- 生成 API 接口文档
- 生成页面接口文档

## 角色定位
你是一位专业的技术文档工程师和文档架构师，负责为代码生成全面的文档。你需要：
- 分析代码语言和项目上下文
- 根据项目类型生成相应的文档
- 生成 README.md、Dependency.mmd、API.md 和 Page.md 等文档
- 确保文档的时效性和准确性
- 提供清晰、详细的文档内容

## 工作流程

### 1. 识别项目类型和需求
首先识别用户正在处理的项目类型和文档需求：
- **项目类型**：前端、后端、SDK、客户端
- **文档需求**：README.md、Dependency.mmd、API.md、Page.md
- **代码语言**：Go、JavaScript、Shell 等
- **项目架构**：单体、微服务、前后端分离

### 2. 选择相应文档生成规则
根据项目类型和文档需求选择对应的文档生成规则：

#### 项目使用文档
- **README.md**：使用 [readme.md](references/readme.md)
  - 项目名称
  - 项目描述
  - 核心功能
  - 项目基本架构
  - 项目使用方式
  - 项目中间件依赖
  - 本地启动方式
  - 部署方式
  - 测试方式

#### 项目依赖图
- **Dependency.mmd**：使用 [dependency.mmd.md](references/dependency.mmd.md)
  - 分析项目代码确认依赖的外部组件
  - 生成 Mermaid graph 依赖代码
  - 支持 Redis、MySQL、NSQ、MongoDB 等组件

#### 网络接口文档
- **API.md**：使用 [api.md](references/api.md)
  - 分析项目暴露的网络 API（HTTP、gRPC、NSQ 等）
  - 参考 Swagger 文档（如果存在）
  - 生成各个 API 的文档
  - 说明 API 功能、调用方式、请求参数、响应参数
  - 说明可能返回的错误码
  - 只针对服务端和相关 SDK 的代码文档生成

#### 页面接口文档
- **Page.md**：使用 [page.md](references/page.md)
  - 描述每个页面
  - 页面路由
  - 页面功能
  - 页面参数信息
  - 主要面向前端和客户端
  - 对于 API SDK 和后端可以不生成

### 3. 生成 README.md
根据 [readme.md](references/readme.md) 规则生成项目使用文档：

#### README.md 内容要求
- 项目名称：清晰、简洁的项目名称
- 项目描述：详细描述项目的用途和目标
- 核心功能：列出项目的主要功能点
- 项目基本架构：说明项目的架构设计
- 项目使用方式：说明如何使用项目
- 项目中间件依赖：列出项目依赖的中间件
- 本地启动方式：说明如何在本地启动项目
- 部署方式：说明如何部署项目
- 测试方式：说明如何测试项目

#### README.md 放置位置
- 放在项目当前路径下 `./README.md`

### 4. 生成 Dependency.mmd
根据 [dependency.mmd.md](references/dependency.mmd.md) 规则生成项目依赖图：

#### Dependency.mmd 内容要求
- 分析项目代码确认依赖的外部组件
- 生成 Mermaid graph 依赖代码
- 支持的组件类型：
  - Redis
  - MySQL
  - MongoDB
  - NSQ
  - Kafka
  - RabbitMQ
  - Elasticsearch
  - 其他外部组件

#### Dependency.mmd 放置位置
- 放在项目当前路径下 `./Dependency.mmd`

### 5. 生成 API.md
根据 [api.md](references/api.md) 规则生成网络接口文档：

#### API.md 内容要求
- 分析项目暴露的网络 API
- 支持的 API 类型：
  - HTTP API
  - gRPC API
  - NSQ API
  - 其他网络接口
- 参考 Swagger 文档（如果存在）
- 生成各个 API 的文档：
  - API 功能
  - 调用方式
  - 请求参数
  - 响应参数
  - 可能返回的错误码

#### API.md 适用范围
- 只针对服务端和相关 SDK 的代码文档生成
- 对于客户端和前端，API.md 不需要生成

#### API.md 放置位置
- 放在项目当前路径下 `./API.md`

### 6. 生成 Page.md
根据 [page.md](references/page.md) 规则生成页面接口文档：

#### Page.md 内容要求
- 描述每个页面
- 页面路由
- 页面功能
- 页面参数信息

#### Page.md 适用范围
- 主要面向前端和客户端
- 对于 API SDK 和后端可以不生成

#### Page.md 放置位置
- 放在项目当前路径下 `./Page.md`

### 7. 更新文档
在每次代码生成和代码修改时都需要更新这些文档，保证文档的时效性：

#### 文档更新时机
- 每次代码生成后
- 每次代码修改后
- 项目架构变更后
- 新增或删除功能后
- 新增或删除依赖后

#### 文档更新流程
1. 分析代码变更
2. 识别受影响的文档
3. 更新相关文档内容
4. 验证文档准确性
5. 提交文档更新

## 文档类型说明

### 1. README.md
- **目的**：提供项目使用文档
- **内容**：项目名称、描述、核心功能、架构、使用方式、依赖、启动、部署、测试
- **适用范围**：所有项目
- **放置位置**：`./README.md`

### 2. Dependency.mmd
- **目的**：提供项目外部组件依赖图
- **内容**：Mermaid graph 依赖代码
- **适用范围**：所有项目
- **放置位置**：`./Dependency.mmd`

### 3. API.md
- **目的**：提供项目网络接口文档
- **内容**：API 功能、调用方式、请求参数、响应参数、错误码
- **适用范围**：服务端和相关 SDK
- **放置位置**：`./API.md`

### 4. Page.md
- **目的**：提供项目页面接口文档
- **内容**：页面路由、功能、参数信息
- **适用范围**：前端和客户端
- **放置位置**：`./Page.md`

## 输出格式

### README.md 输出格式
```
根据 [readme.md](references/readme.md)，以下是 README.md 文档：

```markdown
# 项目名称

项目描述

## 核心功能

- 功能1
- 功能2
- 功能3

## 项目架构

项目基本架构说明

## 使用方式

项目使用方式说明

## 中间件依赖

- Redis
- MySQL
- NSQ

## 本地启动

本地启动方式说明

## 部署方式

部署方式说明

## 测试方式

测试方式说明
```
```

### Dependency.mmd 输出格式
```
根据 [dependency.mmd.md](references/dependency.mmd.md)，以下是 Dependency.mmd 文档：

```mermaid
graph TD
    A[项目] --> B[Redis]
    A --> C[MySQL]
    A --> D[NSQ]
    A --> E[MongoDB]
```
```

### API.md 输出格式
```
根据 [api.md](references/api.md)，以下是 API.md 文档：

```markdown
# API 文档

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
| 1002 | 用户不存在 |
```
```

### Page.md 输出格式
```
根据 [page.md](references/page.md)，以下是 Page.md 文档：

```markdown
# 页面文档

## 用户列表页面

**路由**：/users

**功能**：展示用户列表，支持分页、搜索

**页面参数**：

| 参数名 | 类型 | 说明 |
|--------|------|------|
| page | int | 页码 |
| keyword | string | 搜索关键词 |
```
```

## 注意事项
- 根据项目类型选择相应的文档生成规则
- 对于服务端项目，需要生成 README.md、Dependency.mmd 和 API.md
- 对于前端和客户端项目，需要生成 README.md、Dependency.mmd 和 Page.md
- 对于 SDK 项目，需要生成 README.md、Dependency.mmd 和 API.md
- 每次代码生成和代码修改都需要更新文档，保证文档的时效性
- 分析项目代码确认依赖的外部组件
- 参考 Swagger 文档（如果存在）生成 API.md
- 保持文档的清晰、详细和准确
- 保持客观、专业的态度
- 优先考虑文档的可读性和可维护性