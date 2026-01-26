# 服务端高并发规范

## 并发安全

### 共享数据保护
- 所有共享数据必须使用锁或其他同步机制保护
- 使用互斥锁（Mutex）保护临界区
- 使用读写锁（RWMutex）保护读多写少的场景
- 使用原子操作（Atomic）保护简单数据类型
- 避免在锁内执行耗时操作

### 锁的使用
- 锁的粒度应尽可能小
- 锁的持有时间应尽可能短
- 使用 defer 确保锁的释放
- 避免死锁（按固定顺序获取锁）
- 避免锁嵌套
- 避免在锁内调用外部函数

### 无锁编程
- 优先使用无锁数据结构
- 使用 CAS（Compare-And-Swap）操作
- 使用 channel 进行 goroutine 间通信
- 避免使用全局变量

### 线程局部存储
- 使用线程局部存储避免共享数据
- 使用 context 传递请求上下文
- 避免在 goroutine 间传递大对象

## 并发数量控制

### Goroutine 池
- 使用 goroutine 池限制并发数量
- 池大小根据 CPU 核心数和任务类型设置
- CPU 密集型任务：goroutine 数量 = CPU 核心数
- IO 密集型任务：goroutine 数量 = CPU 核心数 * 2
- 使用 sync.Pool 复用对象

### 信号量
- 使用信号量控制并发数量
- 使用带缓冲的 channel 作为信号量
- 设置合理的并发上限
- 动态调整并发数量

### 限流
- 实现接口级别的限流
- 使用令牌桶算法限流
- 使用漏桶算法限流
- 限制单个用户的并发请求数
- 限制单个 IP 的并发请求数

### 队列
- 使用队列缓冲请求
- 使用消息队列处理异步任务
- 队列长度应有限制
- 队列满时拒绝请求

## 并发编程场景

### 优先使用并行的场景
- 多个独立的 IO 操作
- 多个独立的计算任务
- 批量数据处理
- 并发请求外部服务
- 并发读写多个数据源

### 并发编程规范
- 多个可以并行的 IO 操作优先使用多线程多协程并发编程
- 不可预知并发数量的场景，应该限制并发 goroutine 数量
- CPU 密集型任务：goroutine 数量 = CPU 核心数
- IO 密集型任务：goroutine 数量 = CPU 核心数 * 2

### IO 密集型任务
- 使用 goroutine 并发执行 IO 操作
- 使用 channel 收集结果
- 使用 context 控制超时
- 使用 errgroup 管理多个 goroutine

### CPU 密集型任务
- 使用 worker pool 模式
- 根据 CPU 核心数设置 worker 数量
- 使用 channel 分发任务
- 使用 channel 收集结果

### 批量处理
- 使用 goroutine 并发处理批量数据
- 每个批次大小适中（100-1000）
- 使用 channel 分发批次
- 使用 channel 收集结果

### 并发请求外部服务
- 使用 goroutine 并发请求多个外部服务
- 使用 context 控制超时
- 使用 channel 收集结果
- 实现熔断和降级

## 并发模式

### Worker Pool 模式
- 使用固定数量的 worker 处理任务
- 使用 channel 分发任务
- 使用 channel 收集结果
- 使用 sync.WaitGroup 等待 worker 完成

### Pipeline 模式
- 将任务分解为多个阶段
- 每个阶段使用独立的 goroutine
- 使用 channel 连接各个阶段
- 使用 context 控制取消

### Fan-out/Fan-in 模式
- 使用多个 goroutine 并发处理任务
- 使用 channel 分发任务（fan-out）
- 使用 channel 收集结果（fan-in）
- 使用 errgroup 管理 goroutine

### Future/Promise 模式
- 使用 goroutine 异步执行任务
- 使用 channel 返回结果
- 使用 select 等待多个结果
- 使用 context 控制超时

## Context 使用

### Context 传递
- context 应作为函数的第一个参数
- context 应贯穿整个调用链
- 不将 context 存储在结构体中
- 不将 context 作为结构体字段

### Context 取消
- 使用 context.WithCancel 创建可取消的 context
- 使用 context.WithTimeout 创建带超时的 context
- 使用 context.WithDeadline 创建带截止时间的 context
- 使用 context.WithValue 传递请求上下文

### Context 使用场景
- HTTP 请求使用 context
- 数据库查询使用 context
- 外部服务调用使用 context
- goroutine 间通信使用 context

### Context 注意事项
- 不要传递 nil context
- 不要将 context 放入结构体
- 不要将 context 作为函数的可选参数
- context 应该是可变的，传递时应传递其副本

## 超时控制

### 超时设置
- 所有外部调用必须设置超时
- HTTP 请求超时：5-30 秒
- 数据库查询超时：1-10 秒
- 缓存操作超时：100ms-1 秒
- 消息队列操作超时：1-5 秒

### 超时实现
- 使用 context.WithTimeout 设置超时
- 使用 http.Client.Timeout 设置 HTTP 超时
- 使用数据库驱动的超时参数
- 使用 select + time.After 实现超时

### 超时处理
- 超时后取消正在进行的操作
- 超时后返回明确的错误信息
- 超时后释放资源
- 超时后记录日志

### 超时优化
- 根据业务场景设置合理的超时时间
- 根据网络状况动态调整超时时间
- 使用指数退避重试
- 实现熔断机制

## 错误处理

### 错误传播
- 使用 context 传递取消信号
- 使用 channel 传递错误
- 使用 errgroup 管理 goroutine 错误
- 使用 defer recover 捕获 panic

### 错误重试
- 实现自动重试机制
- 重试次数不超过 3 次
- 重试间隔指数退避
- 只重试可重试的错误

### 错误降级
- 实现服务降级
- 降级策略清晰
- 提供降级通知
- 记录降级日志

### 错误监控
- 监控错误率
- 监控错误类型
- 设置错误告警
- 分析错误原因

## 资源管理

### 连接池
- 使用连接池管理数据库连接
- 使用连接池管理 Redis 连接
- 使用连接池管理 HTTP 连接
- 配置合理的连接池大小
- 设置合理的连接超时时间

### 内存管理
- 避免内存泄漏
- 及时释放不再使用的资源
- 使用 sync.Pool 复用对象
- 避免频繁的内存分配

### 文件描述符
- 及时关闭文件
- 及时关闭网络连接
- 及时关闭数据库连接
- 避免文件描述符泄漏

### Goroutine 管理
- 避免创建过多的 goroutine
- 使用 goroutine 池管理 goroutine
- 及时退出不再需要的 goroutine
- 避免 goroutine 泄漏

## 性能优化

### 减少锁竞争
- 减小锁的粒度
- 使用读写锁
- 使用无锁数据结构
- 使用原子操作

### 减少上下文切换
- 减少 goroutine 数量
- 使用 goroutine 池
- 减少锁的使用
- 使用无锁编程

### 减少内存分配
- 使用 sync.Pool 复用对象
- 预分配切片和 map 的容量
- 避免频繁的内存分配
- 使用值类型而非指针类型

### 减少系统调用
- 使用批量操作
- 使用缓存
- 使用连接池
- 使用异步 IO

## 监控和调试

### 性能监控
- 监控并发数
- 监控响应时间
- 监控错误率
- 监控资源使用

### 日志记录
- 记录并发操作
- 记录锁竞争
- 记录超时
- 记录错误

### 性能分析
- 使用 pprof 分析性能
- 使用 trace 分析并发
- 使用 race detector 检测竞态
- 使用 benchstat 比较性能

### 调试技巧
- 使用调试日志
- 使用断点调试
- 使用性能分析工具
- 使用并发检测工具

## 最佳实践

### 并发设计
- 优先使用 channel 进行通信
- 避免共享内存
- 使用 context 控制生命周期
- 使用 errgroup 管理 goroutine

### 错误处理
- 及时处理错误
- 提供有用的错误信息
- 实现重试机制
- 实现降级机制

### 资源管理
- 及时释放资源
- 使用连接池
- 使用对象池
- 避免资源泄漏

### 性能优化
- 减少锁竞争
- 减少上下文切换
- 减少内存分配
- 减少系统调用

## 禁止事项

- 禁止在锁内执行耗时操作
- 禁止在锁内调用外部函数
- 禁止在锁内创建 goroutine
- 禁止在锁内使用 channel
- 禁止在锁内使用 defer recover
- 禁止在锁内使用 panic
- 禁止在锁内使用 os.Exit
- 禁止在锁内使用 log.Fatal
- 禁止在锁内使用 log.Panic
- 禁止在锁内使用 runtime.Goexit
- 禁止在锁内使用 runtime.Gosched
- 禁止在锁内使用 runtime.LockOSThread
- 禁止在锁内使用 runtime.UnlockOSThread
- 禁止在锁内使用 runtime.GC
- 禁止在锁内使用 runtime.SetFinalizer
- 禁止在锁内使用 runtime.SetCPUProfileRate
- 禁止在锁内使用 runtime.SetBlockProfileRate
- 禁止在锁内使用 runtime.SetMutexProfileFraction
- 禁止在锁内使用 runtime.MemProfileRate
