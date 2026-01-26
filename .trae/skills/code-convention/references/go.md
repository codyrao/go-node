# Go 语言代码规范

## 命名规范

### 包名
- 包名使用小写单词，不使用下划线或驼峰命名
- 包名应简短、有意义，通常使用单个单词
- 包名与所在目录名保持一致
- 避免使用 `util`、`common`、`base` 等过于泛化的包名

### 变量和常量
- 使用驼峰命名法（camelCase）
- 导出的变量和常量首字母大写
- 局部变量首字母小写
- 常量使用全大写，单词间用下划线分隔（仅在跨包共享时）
- 布尔变量应以 `Is`、`Has`、`Can`、`Should` 等前缀开头

### 函数和方法
- 使用驼峰命名法
- 导出的函数和方法首字母大写
- 私有函数和方法首字母小写
- 函数名应清晰表达其功能，避免使用缩写
- 接口方法名应简洁，通常以动词开头

### 接口
- 接口名通常以 `-er` 后缀结尾（如 `Reader`、`Writer`）
- 单方法接口命名以方法名加 `er` 后缀
- 接口应保持小而专一，避免臃肿

### 结构体
- 使用驼峰命名法
- 导出的结构体首字母大写
- 私有结构体首字母小写
- 结构体字段遵循相同的命名规则

## 代码格式

### 使用 gofmt
- 所有代码必须通过 `gofmt` 格式化
- 使用 `goimports` 管理导入包的顺序和分组

### 缩进和空格
- 使用 Tab 缩进，不使用空格
- 运算符两侧添加空格
- 逗号后添加空格
- 函数参数列表中逗号后添加空格

### 行长度
- 单行代码不超过 120 字符
- 长字符串使用反引号或拼接

### 大括号
- 左大括号不换行
- 右大括号独占一行

### 注释
- 包注释应放在包声明之前，描述包的功能
- 导出的函数、方法、类型、常量、变量必须有注释
- 注释以被注释对象的名称开头
- 使用完整的句子，首字母大写，以句号结尾
- 注释应解释"为什么"而不是"是什么"
- **每个函数都需要加上中文注释，注释说明函数作用**
- **函数中特别复杂的逻辑也需要加上中文注释，说明为什么这么写**

## 错误处理

### 错误检查
- 所有可能返回错误的函数调用必须检查错误
- 不要忽略错误，使用 `_` 丢弃错误必须有明确理由
- 错误信息应包含足够的上下文信息

### 错误包装
- 使用 `fmt.Errorf` 和 `%w` 包装错误
- 使用 `errors.Is` 和 `errors.As` 判断错误类型
- 避免重复包装错误

### 错误定义
- 使用 `errors.New` 创建简单错误
- 实现 `error` 接口定义自定义错误类型
- 错误变量以 `Err` 前缀开头

### panic 使用
- 仅在不可恢复的错误情况下使用 panic
- 不应在普通业务逻辑中使用 panic
- 库代码应返回错误而不是 panic

## 并发安全

### 互斥锁
- 使用 `sync.Mutex` 保护共享数据
- 使用 `defer` 确保锁的释放
- 避免在持有锁时调用外部函数

### 读写锁
- 读多写少的场景使用 `sync.RWMutex`
- 注意读写锁的性能开销

### WaitGroup
- 使用 `sync.WaitGroup` 等待一组 goroutine 完成
- 确保 `Add` 在 `Wait` 之前调用

### Channel
- 优先使用 channel 进行 goroutine 间通信
- 避免在 channel 上发送 nil 值
- 关闭 channel 应由发送方执行
- 使用 `range` 接收 channel 数据

### Context
- 使用 `context.Context` 传递取消信号和超时
- 不要将 context 存储在结构体中
- context 应作为函数的第一个参数

## 性能优化

### 内存分配
- 避免不必要的内存分配
- 使用 `sync.Pool` 复用对象
- 预分配切片和 map 的容量

### 字符串处理
- 使用 `strings.Builder` 拼接字符串
- 避免频繁的字符串转换

### 循环优化
- 将不变的计算移出循环
- 使用 `for range` 遍历切片和 map
- 避免在循环中创建闭包

### 反射
- 避免在性能敏感的代码中使用反射
- 反射代码应做好性能测试

## 依赖管理

### 导入规范
- 按标准库、第三方库、本地包的顺序分组
- 使用空行分隔不同组
- 删除未使用的导入

### 版本管理
- 使用 Go Modules 管理依赖
- 在 `go.mod` 中明确指定依赖版本
- 定期更新依赖版本

### 避免循环依赖
- 包之间不应存在循环依赖
- 通过接口解耦依赖关系

## 测试规范

### 测试文件
- 测试文件以 `_test.go` 结尾
- 测试函数以 `Test` 开头
- 基准测试以 `Benchmark` 开头
- 示例测试以 `Example` 开头

### 测试覆盖
- 核心业务逻辑测试覆盖率应不低于 80%
- 边界条件和错误情况必须有测试

### 表驱动测试
- 使用表驱动测试覆盖多种情况
- 测试用例命名清晰描述测试场景

### Mock 和桩
- 使用接口和 mock 进行单元测试
- 避免在测试中使用真实的外部依赖

## 安全规范

### 输入验证
- 所有外部输入必须进行验证
- 验证失败应返回明确的错误信息

### SQL 注入
- 使用参数化查询
- 避免拼接 SQL 语句

### 敏感信息
- 不在代码中硬编码密钥、密码等敏感信息
- 使用环境变量或配置管理敏感信息

### 并发安全
- 确保共享数据的并发访问安全
- 使用原子操作或锁保护临界区

## 最佳实践

### 接口设计
- 接口应小而专一
- 接口应在使用者端定义
- 避免过早抽象

### 错误处理
- 尽早返回错误
- 错误处理逻辑清晰
- 提供有用的错误信息

### 资源管理
- 使用 `defer` 确保资源释放
- 及时关闭文件、数据库连接等资源
- 使用 `io.Close` 接口统一管理资源

### 日志记录
- 使用结构化日志
- 日志级别合理使用
- 避免记录敏感信息

### 配置管理
- 使用配置文件管理配置
- 配置项有合理的默认值
- 支持环境变量覆盖配置

## 禁止事项

- 禁止使用 `goto` 语句
- 禁止在循环中创建不必要的 goroutine
- 禁止忽略错误
- 禁止在公共 API 中使用可变参数
- 禁止在循环中捕获循环变量
- 禁止使用 `time.Sleep` 等待 goroutine
- 禁止在 defer 中修改返回值

## 代码复用规范

- 如果发现相似的代码超过2个，就要封装成函数调用
- 提取公共逻辑到独立的函数中
- 避免代码重复，提高代码复用性
- 函数命名应清晰表达其功能

**示例**：
```go
// 错误示例：相似的代码重复
func ProcessUserA(user User) error {
    if user.Name == "" {
        return errors.New("name is empty")
    }
    if user.Email == "" {
        return errors.New("email is empty")
    }
    // 处理用户A
    return nil
}

func ProcessUserB(user User) error {
    if user.Name == "" {
        return errors.New("name is empty")
    }
    if user.Email == "" {
        return errors.New("email is empty")
    }
    // 处理用户B
    return nil
}

// 正确示例：封装成函数调用
// ValidateUser 验证用户信息
func ValidateUser(user User) error {
    if user.Name == "" {
        return errors.New("name is empty")
    }
    if user.Email == "" {
        return errors.New("email is empty")
    }
    return nil
}

func ProcessUserA(user User) error {
    if err := ValidateUser(user); err != nil {
        return err
    }
    // 处理用户A
    return nil
}

func ProcessUserB(user User) error {
    if err := ValidateUser(user); err != nil {
        return err
    }
    // 处理用户B
    return nil
}
```

## 函数长度规范

- 函数的长度不宜过长，不易超过1千行
- 如果超过一千行，需要封装内部子函数缩短单个函数的长度
- 函数应保持单一职责，只做一件事
- 复杂的逻辑应拆分为多个小函数

**示例**：
```go
// 错误示例：函数过长
func ProcessOrder(order Order) error {
    // 1000多行代码
    // 验证订单
    // 检查库存
    // 计算价格
    // 创建支付
    // 发送通知
    // 更新库存
    // 记录日志
    // ... 更多代码
    return nil
}

// 正确示例：封装内部子函数
// ValidateOrder 验证订单
func ValidateOrder(order Order) error {
    // 验证逻辑
    return nil
}

// CheckStock 检查库存
func CheckStock(order Order) error {
    // 检查库存逻辑
    return nil
}

// CalculatePrice 计算价格
func CalculatePrice(order Order) float64 {
    // 计算价格逻辑
    return 0
}

// CreatePayment 创建支付
func CreatePayment(order Order) error {
    // 创建支付逻辑
    return nil
}

// SendNotification 发送通知
func SendNotification(order Order) error {
    // 发送通知逻辑
    return nil
}

// ProcessOrder 处理订单
func ProcessOrder(order Order) error {
    if err := ValidateOrder(order); err != nil {
        return err
    }
    
    if err := CheckStock(order); err != nil {
        return err
    }
    
    price := CalculatePrice(order)
    
    if err := CreatePayment(order); err != nil {
        return err
    }
    
    if err := SendNotification(order); err != nil {
        return err
    }
    
    return nil
}
```
