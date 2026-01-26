# README.md 文档生成规则

## 概述
本文档定义了 README.md 的生成规则，README.md 是项目的使用文档，需要说明项目的基本信息、架构、使用方式、依赖、启动、部署和测试等。

## 文档结构

### 1. 项目名称
- **要求**：清晰、简洁的项目名称
- **格式**：一级标题 `# 项目名称`
- **示例**：
```markdown
# 用户管理系统
```

### 2. 项目描述
- **要求**：详细描述项目的用途和目标
- **格式**：二级标题 `## 项目描述`
- **内容**：
  - 项目的主要用途
  - 项目解决的问题
  - 项目的目标用户
  - 项目的应用场景

- **示例**：
```markdown
## 项目描述

用户管理系统是一个用于管理用户信息的后端服务，提供用户的增删改查、权限管理、数据统计等功能。该系统适用于企业内部用户管理、SaaS 平台用户管理等场景。
```

### 3. 核心功能
- **要求**：列出项目的主要功能点
- **格式**：二级标题 `## 核心功能`
- **内容**：使用无序列表列出核心功能

- **示例**：
```markdown
## 核心功能

- 用户管理：支持用户的增删改查
- 权限管理：支持角色和权限的管理
- 数据统计：支持用户数据的统计分析
- 日志记录：支持操作日志的记录和查询
```

### 4. 项目架构
- **要求**：说明项目的架构设计
- **格式**：二级标题 `## 项目架构`
- **内容**：
  - 项目的基本架构
  - 主要模块和组件
  - 模块之间的关系
  - 技术栈

- **示例**：
```markdown
## 项目架构

项目采用微服务架构，主要包含以下模块：

- 用户服务：负责用户的基本管理
- 权限服务：负责权限的管理和验证
- 统计服务：负责数据的统计分析
- 网关服务：负责请求的路由和转发

技术栈：
- 后端：Go 1.16+
- 数据库：MySQL 8.0
- 缓存：Redis 6.0
- 消息队列：NSQ
```

### 5. 使用方式
- **要求**：说明如何使用项目
- **格式**：二级标题 `## 使用方式`
- **内容**：
  - SDK 使用方式（如果是 SDK 项目）
  - API 调用方式（如果是服务端项目）
  - 前端使用方式（如果是前端项目）

- **示例（SDK 项目）**：
```markdown
## 使用方式

### 安装

```bash
go get github.com/example/user-sdk
```

### 初始化

```go
import "github.com/example/user-sdk"

client := user.NewClient("http://api.example.com", "your-api-key")
```

### 使用示例

```go
user, err := client.GetUser(1)
if err != nil {
    log.Fatal(err)
}
fmt.Println(user)
```
```

- **示例（服务端项目）**：
```markdown
## 使用方式

### API 调用

所有 API 都通过 HTTP 协议调用，基础 URL 为：`http://api.example.com`

详细 API 文档请参考 [API.md](./API.md)
```

### 6. 中间件依赖
- **要求**：列出项目依赖的中间件
- **格式**：二级标题 `## 中间件依赖`
- **内容**：使用无序列表列出依赖的中间件

- **示例**：
```markdown
## 中间件依赖

- Redis 6.0+：用于缓存用户信息
- MySQL 8.0+：用于持久化存储用户数据
- NSQ 1.2+：用于异步消息处理
- Elasticsearch 7.0+：用于用户搜索
```

### 7. 本地启动
- **要求**：说明如何在本地启动项目
- **格式**：二级标题 `## 本地启动`
- **内容**：
  - 环境要求
  - 配置说明
  - 启动步骤

- **示例**：
```markdown
## 本地启动

### 环境要求

- Go 1.16+
- MySQL 8.0+
- Redis 6.0+
- NSQ 1.2+

### 配置说明

复制配置文件模板：

```bash
cp config/config.example.yaml config/config.yaml
```

修改 `config/config.yaml` 中的数据库、Redis 等配置。

### 启动步骤

1. 启动依赖服务

```bash
docker-compose up -d
```

2. 安装依赖

```bash
go mod download
```

3. 运行项目

```bash
go run main.go
```

4. 访问服务

打开浏览器访问：`http://localhost:8080`
```

### 8. 部署方式
- **要求**：说明如何部署项目
- **格式**：二级标题 `## 部署方式`
- **内容**：
  - 部署环境要求
  - 部署步骤
  - 配置说明

- **示例**：
```markdown
## 部署方式

### 环境要求

- Linux 服务器
- Go 1.16+
- MySQL 8.0+
- Redis 6.0+
- NSQ 1.2+

### 部署步骤

1. 编译项目

```bash
go build -o user-service main.go
```

2. 上传到服务器

```bash
scp user-service user@server:/opt/user-service/
```

3. 配置服务

修改服务器上的配置文件 `/opt/user-service/config/config.yaml`

4. 启动服务

```bash
cd /opt/user-service
./user-service
```

5. 使用 systemd 管理

创建 `/etc/systemd/system/user-service.service` 文件：

```ini
[Unit]
Description=User Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/user-service
ExecStart=/opt/user-service/user-service
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
systemctl daemon-reload
systemctl start user-service
systemctl enable user-service
```
```

### 9. 测试方式
- **要求**：说明如何测试项目
- **格式**：二级标题 `## 测试方式`
- **内容**：
  - 单元测试
  - 集成测试
  - 系统测试

- **示例**：
```markdown
## 测试方式

### 单元测试

```bash
go test ./...
```

### 集成测试

```bash
go test -tags=integration ./...
```

### 系统测试

```bash
cd test
./run.sh all
```

详细测试说明请参考 [测试文档](./test/README.md)
```

### 10. 其他章节（可选）
根据项目需要，可以添加以下章节：

#### 贡献指南
```markdown
## 贡献指南

欢迎贡献代码，请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request
```

#### 许可证
```markdown
## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件
```

#### 联系方式
```markdown
## 联系方式

如有问题或建议，请通过以下方式联系：

- Email: support@example.com
- Issue: https://github.com/example/user-service/issues
```

## 文档模板

### 完整模板
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

环境要求

配置说明

启动步骤

## 部署方式

环境要求

部署步骤

配置说明

## 测试方式

单元测试

集成测试

系统测试

## 贡献指南

贡献指南说明

## 许可证

许可证说明

## 联系方式

联系方式说明
```

## 注意事项
- README.md 应该简洁明了，易于阅读
- 使用 Markdown 格式，确保格式正确
- 代码示例应该清晰、可运行
- 配置说明应该详细，避免遗漏
- 部署方式应该包含完整的步骤
- 测试方式应该包含所有测试类型
- 保持文档的时效性，及时更新
