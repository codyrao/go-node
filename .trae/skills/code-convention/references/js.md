# JavaScript 语言代码规范

## 命名规范

### 变量和函数
- 使用驼峰命名法（camelCase）
- 变量名使用有意义的名词
- 函数名使用动词或动词短语
- 布尔变量以 `is`、`has`、`can`、`should` 等前缀开头
- 私有变量以 `_` 前缀开头

### 常量
- 使用全大写，单词间用下划线分隔
- 常量应使用 `const` 声明
- 常量名应表达其值的含义

### 类和构造函数
- 使用帕斯卡命名法（PascalCase）
- 类名应为名词
- 构造函数名与类名相同

### 文件和目录
- 文件名使用小写字母，单词间用连字符分隔
- 目录名使用小写字母，单词间用连字符分隔
- 组件文件名应与组件名保持一致

## 代码格式

### 使用 ESLint
- 所有代码必须通过 ESLint 检查
- 使用 Prettier 进行代码格式化
- 在项目中配置统一的 ESLint 和 Prettier 规则

### 缩进和空格
- 使用 2 个空格缩进
- 运算符两侧添加空格
- 逗号后添加空格
- 对象和数组的括号内添加空格

### 行长度
- 单行代码不超过 100 字符
- 长字符串使用模板字符串或拼接

### 分号
- 始终使用分号结束语句
- 避免依赖自动分号插入（ASI）

### 引号
- 优先使用单引号
- 字符串中包含单引号时使用双引号
- 模板字符串使用反引号

### 大括号
- 左大括号不换行
- 右大括号独占一行

### 注释
- 单行注释使用 `//`
- 多行注释使用 `/* */`
- JSDoc 注释用于函数、类、模块的文档
- 注释应解释"为什么"而不是"是什么"
- **每个函数都需要加上中文注释，注释说明函数作用**
- **函数中特别复杂的逻辑也需要加上中文注释，说明为什么这么写**

## 变量声明

### 使用 const 和 let
- 优先使用 `const` 声明不会重新赋值的变量
- 使用 `let` 声明需要重新赋值的变量
- 不使用 `var` 声明变量

### 作用域
- 使用块级作用域
- 避免在块外使用块内声明的变量
- 变量应在使用前声明

### 变量提升
- 避免依赖变量提升
- 函数声明会提升，函数表达式不会

## 函数规范

### 函数声明
- 优先使用函数声明（function 声明）
- 函数声明会被提升
- 函数表达式用于回调和高阶函数

### 参数
- 参数数量不超过 5 个
- 参数过多时使用对象参数
- 为参数提供默认值
- 使用解构赋值提取参数

### 返回值
- 函数应尽早返回
- 避免使用多个返回值
- 返回值类型应一致

### 箭头函数
- 箭头函数用于回调和高阶函数
- 箭头函数不绑定 `this`
- 避免在对象方法中使用箭头函数

## 对象和数组

### 对象
- 使用对象字面量创建对象
- 使用简写属性和方法
- 使用计算属性名
- 使用展开运算符复制和合并对象

### 数组
- 使用数组字面量创建数组
- 使用展开运算符复制和合并数组
- 使用数组方法（map、filter、reduce）处理数组
- 避免使用 `for` 循环遍历数组

### 解构赋值
- 使用解构赋值提取对象和数组的值
- 为解构的值提供默认值
- 使用解构交换变量值

## 异步编程

### Promise
- 使用 Promise 处理异步操作
- 使用 `Promise.all` 并行执行多个 Promise
- 使用 `Promise.race` 获取最快完成的 Promise

### async/await
- 优先使用 async/await 处理异步操作
- async 函数总是返回 Promise
- 使用 try-catch 捕获 async 函数中的错误

### 错误处理
- 始终处理 Promise 的错误
- 使用 catch 捕获 Promise 的错误
- 使用 try-catch 捕获 async/await 的错误
- 错误信息应包含足够的上下文

## 类和模块

### 类
- 使用 class 语法定义类
- 使用 constructor 定义构造函数
- 使用 extends 继承类
- 使用 super 调用父类方法
- 使用 static 定义静态方法

### 模块
- 使用 ES6 模块语法
- 使用 export 导出模块
- 使用 import 导入模块
- 避免使用默认导出
- 导入时使用命名导入

## DOM 操作

### 选择元素
- 使用 `querySelector` 和 `querySelectorAll` 选择元素
- 避免使用 `getElementById` 和 `getElementsByClassName`

### 事件处理
- 使用 `addEventListener` 添加事件监听器
- 避免使用 HTML 属性添加事件处理
- 使用事件委托处理动态元素的事件

### 性能优化
- 避免频繁的 DOM 操作
- 使用文档片段批量操作 DOM
- 使用 `requestAnimationFrame` 优化动画

## 性能优化

### 避免全局变量
- 避免在全局作用域声明变量
- 使用 IIFE 或模块隔离作用域

### 事件节流和防抖
- 对频繁触发的事件使用节流（throttle）
- 对需要延迟执行的事件使用防抖（debounce）

### 懒加载
- 对图片和组件使用懒加载
- 使用动态导入（dynamic import）懒加载模块

### 内存管理
- 及时清理不再使用的引用
- 避免内存泄漏
- 使用 WeakMap 和 WeakSet 管理弱引用

## 安全规范

### XSS 防护
- 避免使用 `innerHTML` 插入用户输入
- 使用 `textContent` 或 `innerText` 插入文本
- 使用 DOMPurify 清理 HTML

### CSRF 防护
- 使用 CSRF Token
- 验证请求的来源

### 敏感信息
- 不在代码中硬编码密钥、密码等敏感信息
- 使用环境变量管理敏感信息
- 不在前端代码中暴露敏感信息

## 测试规范

### 测试框架
- 使用 Jest 或 Mocha 进行单元测试
- 使用 Cypress 或 Playwright 进行端到端测试

### 测试覆盖
- 核心业务逻辑测试覆盖率应不低于 80%
- 边界条件和错误情况必须有测试

### 测试结构
- 使用 describe 分组测试
- 使用 it/test 定义测试用例
- 使用 beforeAll、beforeEach、afterEach、afterAll 管理测试生命周期

### Mock 和桩
- 使用 mock 函数和模块进行单元测试
- 避免在测试中使用真实的外部依赖

## 最佳实践

### 代码复用
- 提取公共逻辑为函数或模块
- 使用高阶函数复用逻辑
- 使用组合优于继承

### 错误处理
- 尽早捕获和处理错误
- 提供有用的错误信息
- 使用自定义错误类型

### 日志记录
- 使用结构化日志
- 日志级别合理使用
- 避免记录敏感信息

### 配置管理
- 使用配置文件管理配置
- 配置项有合理的默认值
- 支持环境变量覆盖配置

## 禁止事项

- 禁止使用 `eval`
- 禁止使用 `with` 语句
- 禁止使用 `==` 和 `!=`，使用 `===` 和 `!==`
- 禁止在循环中创建函数
- 禁止在条件语句中使用赋值
- 禁止使用 `arguments` 对象，使用剩余参数
- 禁止使用 `new Function`
- 禁止在模板字符串中嵌套模板字符串

## 代码复用规范

- 如果发现相似的代码超过2个，就要封装成函数调用
- 提取公共逻辑到独立的函数中
- 避免代码重复，提高代码复用性
- 函数命名应清晰表达其功能

**示例**：
```javascript
// 错误示例：相似的代码重复
function processUserA(user) {
    if (!user.name) {
        throw new Error('name is empty');
    }
    if (!user.email) {
        throw new Error('email is empty');
    }
    // 处理用户A
}

function processUserB(user) {
    if (!user.name) {
        throw new Error('name is empty');
    }
    if (!user.email) {
        throw new Error('email is empty');
    }
    // 处理用户B
}

// 正确示例：封装成函数调用
// validateUser 验证用户信息
function validateUser(user) {
    if (!user.name) {
        throw new Error('name is empty');
    }
    if (!user.email) {
        throw new Error('email is empty');
    }
}

function processUserA(user) {
    validateUser(user);
    // 处理用户A
}

function processUserB(user) {
    validateUser(user);
    // 处理用户B
}
```

## 函数长度规范

- 函数的长度不宜过长，不易超过1千行
- 如果超过一千行，需要封装内部子函数缩短单个函数的长度
- 函数应保持单一职责，只做一件事
- 复杂的逻辑应拆分为多个小函数

**示例**：
```javascript
// 错误示例：函数过长
function processOrder(order) {
    // 1000多行代码
    // 验证订单
    // 检查库存
    // 计算价格
    // 创建支付
    // 发送通知
    // 更新库存
    // 记录日志
    // ... 更多代码
}

// 正确示例：封装内部子函数
// validateOrder 验证订单
function validateOrder(order) {
    // 验证逻辑
}

// checkStock 检查库存
function checkStock(order) {
    // 检查库存逻辑
}

// calculatePrice 计算价格
function calculatePrice(order) {
    // 计算价格逻辑
}

// createPayment 创建支付
function createPayment(order) {
    // 创建支付逻辑
}

// sendNotification 发送通知
function sendNotification(order) {
    // 发送通知逻辑
}

// processOrder 处理订单
function processOrder(order) {
    validateOrder(order);
    checkStock(order);
    const price = calculatePrice(order);
    createPayment(order);
    sendNotification(order);
}
```
