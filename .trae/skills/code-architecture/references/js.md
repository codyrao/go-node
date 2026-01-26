# JavaScript 代码架构规范

## 概述
本文档定义了 JavaScript 的代码架构规范，适用于纯 JavaScript 项目、前端项目和通用 JavaScript 库的开发。对于特定的框架（Vue、React、Angular）和运行环境（Node.js、Electron），请参考对应的架构规范文档。

## 项目结构

### 标准项目结构
```
project-name/
├── src/                   # 源代码目录
│   ├── assets/           # 静态资源
│   │   ├── images/       # 图片资源
│   │   ├── fonts/        # 字体资源
│   │   ├── styles/       # 样式文件
│   │   └── scripts/      # 脚本文件
│   ├── components/       # 可复用组件
│   ├── utils/            # 工具函数
│   ├── services/         # 服务层
│   ├── api/              # API 接口
│   ├── constants/        # 常量定义
│   ├── types/            # 类型定义（TypeScript）
│   ├── hooks/            # 自定义 Hooks（React）
│   ├── store/            # 状态管理
│   ├── router/           # 路由配置
│   ├── views/            # 页面视图
│   ├── layouts/         # 布局组件
│   ├── config/           # 配置文件
│   └── index.js          # 入口文件
├── public/               # 公共资源
│   ├── index.html        # HTML 模板
│   ├── favicon.ico       # 网站图标
│   └── manifest.json     # PWA 配置
├── test/                  # 测试目录
│   ├── unit/             # 单元测试
│   ├── integration/      # 集成测试
│   └── e2e/              # 端到端测试
├── scripts/               # 构建脚本
│   ├── build.js          # 构建脚本
│   └── dev.js            # 开发脚本
├── docs/                  # 文档目录
│   ├── api.md            # API 文档
│   └── components.md     # 组件文档
├── dist/                  # 构建产物（自动生成）
├── node_modules/          # 依赖包（自动生成）
├── package.json           # 项目配置
├── package-lock.json      # 依赖锁定
├── tsconfig.json          # TypeScript 配置
├── .eslintrc.js           # ESLint 配置
├── .prettierrc.js         # Prettier 配置
├── .gitignore             # Git 忽略文件
├── .env                   # 环境变量
├── .env.development       # 开发环境变量
├── .env.production        # 生产环境变量
└── README.md              # 项目说明文档
```

## 目录说明

### src/
**功能**：源代码目录

**说明**：
- 所有源代码都放在 src 目录下
- 按功能模块划分子目录
- 保持目录结构清晰

#### src/assets/
**功能**：静态资源

**说明**：
- images/：图片资源（PNG、JPG、SVG 等）
- fonts/：字体资源
- styles/：全局样式文件
- scripts/：第三方脚本

**示例**：
```
src/assets/
├── images/
│   ├── logo.png
│   └── background.jpg
├── fonts/
│   ├── Roboto-Regular.ttf
│   └── Roboto-Bold.ttf
├── styles/
│   ├── main.css
│   └── variables.css
└── scripts/
    └── third-party.js
```

#### src/components/
**功能**：可复用组件

**说明**：
- 通用组件
- 业务组件
- 按功能分类

**示例**：
```
src/components/
├── common/               # 通用组件
│   ├── Button/
│   │   ├── Button.js
│   │   ├── Button.css
│   │   └── index.js
│   ├── Input/
│   │   ├── Input.js
│   │   ├── Input.css
│   │   └── index.js
│   └── Modal/
│       ├── Modal.js
│       ├── Modal.css
│       └── index.js
└── business/             # 业务组件
    ├── UserCard/
    │   ├── UserCard.js
    │   ├── UserCard.css
    │   └── index.js
    └── ProductList/
        ├── ProductList.js
        ├── ProductList.css
        └── index.js
```

**组件示例**：
```javascript
// src/components/common/Button/Button.js
import React from 'react';
import './Button.css';

const Button = ({ children, onClick, type = 'primary', disabled = false }) => {
    return (
        <button
            className={`btn btn-${type}`}
            onClick={onClick}
            disabled={disabled}
        >
            {children}
        </button>
    );
};

export default Button;
```

#### src/utils/
**功能**：工具函数

**说明**：
- 通用工具函数
- 纯函数
- 无副作用

**示例**：
```
src/utils/
├── date.js               # 日期工具
├── string.js             # 字符串工具
├── array.js              # 数组工具
├── object.js             # 对象工具
├── validation.js         # 验证工具
└── storage.js            # 存储工具
```

**工具函数示例**：
```javascript
// src/utils/date.js
export const formatDate = (date, format = 'YYYY-MM-DD') => {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    
    return format
        .replace('YYYY', year)
        .replace('MM', month)
        .replace('DD', day);
};

export const isToday = (date) => {
    const today = new Date();
    const d = new Date(date);
    return d.toDateString() === today.toDateString();
};
```

#### src/services/
**功能**：服务层

**说明**：
- 业务逻辑
- 数据处理
- 状态管理

**示例**：
```
src/services/
├── userService.js        # 用户服务
├── productService.js     # 产品服务
├── orderService.js       # 订单服务
└── authService.js        # 认证服务
```

**服务示例**：
```javascript
// src/services/userService.js
import api from '../api';

export const userService = {
    async getUsers(params) {
        const response = await api.get('/users', { params });
        return response.data;
    },

    async getUserById(id) {
        const response = await api.get(`/users/${id}`);
        return response.data;
    },

    async createUser(data) {
        const response = await api.post('/users', data);
        return response.data;
    },

    async updateUser(id, data) {
        const response = await api.put(`/users/${id}`, data);
        return response.data;
    },

    async deleteUser(id) {
        const response = await api.delete(`/users/${id}`);
        return response.data;
    }
};
```

#### src/api/
**功能**：API 接口

**说明**：
- HTTP 请求封装
- 接口定义
- 错误处理

**示例**：
```
src/api/
├── index.js              # API 实例
├── user.js               # 用户 API
├── product.js            # 产品 API
└── order.js              # 订单 API
```

**API 示例**：
```javascript
// src/api/index.js
import axios from 'axios';

const api = axios.create({
    baseURL: process.env.API_BASE_URL || 'http://localhost:3000/api',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json'
    }
});

api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

api.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    localStorage.removeItem('token');
                    window.location.href = '/login';
                    break;
                case 403:
                    console.error('Access denied');
                    break;
                case 404:
                    console.error('Resource not found');
                    break;
                case 500:
                    console.error('Server error');
                    break;
            }
        }
        return Promise.reject(error);
    }
);

export default api;
```

#### src/constants/
**功能**：常量定义

**说明**：
- 避免魔法数字
- 统一管理常量
- 便于维护

**示例**：
```
src/constants/
├── index.js              # 常量导出
├── error.js              # 错误码
├── status.js             # 状态码
└── config.js             # 配置常量
```

**常量示例**：
```javascript
// src/constants/error.js
export const ERROR_CODES = {
    SUCCESS: 0,
    INVALID_PARAMS: 1001,
    UNAUTHORIZED: 1002,
    FORBIDDEN: 1003,
    NOT_FOUND: 1004,
    SERVER_ERROR: 1005
};

export const ERROR_MESSAGES = {
    [ERROR_CODES.SUCCESS]: 'Success',
    [ERROR_CODES.INVALID_PARAMS]: 'Invalid parameters',
    [ERROR_CODES.UNAUTHORIZED]: 'Unauthorized',
    [ERROR_CODES.FORBIDDEN]: 'Access denied',
    [ERROR_CODES.NOT_FOUND]: 'Resource not found',
    [ERROR_CODES.SERVER_ERROR]: 'Server error'
};
```

#### src/types/
**功能**：类型定义（TypeScript）

**说明**：
- 接口定义
- 类型别名
- 枚举定义

**示例**：
```
src/types/
├── index.ts              # 类型导出
├── user.ts               # 用户类型
├── product.ts            # 产品类型
└── order.ts              # 订单类型
```

**类型示例**：
```typescript
// src/types/user.ts
export interface User {
    id: number;
    name: string;
    email: string;
    age: number;
    createdAt: Date;
    updatedAt: Date;
}

export interface CreateUserRequest {
    name: string;
    email: string;
    age: number;
}

export interface UpdateUserRequest {
    name?: string;
    email?: string;
    age?: number;
}

export enum UserRole {
    ADMIN = 'admin',
    USER = 'user',
    GUEST = 'guest'
}
```

#### src/store/
**功能**：状态管理

**说明**：
- 全局状态
- 状态更新逻辑
- 状态派生

**示例**：
```
src/store/
├── index.js              # Store 配置
├── modules/             # 状态模块
│   ├── user.js
│   ├── product.js
│   └── order.js
└── getters/              # 状态派生
    └── index.js
```

#### src/router/
**功能**：路由配置

**说明**：
- 路由定义
- 路由守卫
- 路由参数

**示例**：
```
src/router/
├── index.js              # 路由配置
├── routes.js             # 路由定义
└── guards.js             # 路由守卫
```

**路由示例**：
```javascript
// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import routes from './routes';
import { authGuard } from './guards';

const router = createRouter({
    history: createWebHistory(),
    routes
});

router.beforeEach(authGuard);

export default router;
```

#### src/views/
**功能**：页面视图

**说明**：
- 页面组件
- 路由对应的组件
- 业务逻辑

**示例**：
```
src/views/
├── Home/
│   ├── Home.js
│   ├── Home.css
│   └── index.js
├── User/
│   ├── UserList.js
│   ├── UserDetail.js
│   └── index.js
└── Product/
    ├── ProductList.js
    ├── ProductDetail.js
    └── index.js
```

#### src/layouts/
**功能**：布局组件

**说明**：
- 页面布局
- 通用布局
- 布局切换

**示例**：
```
src/layouts/
├── MainLayout/
│   ├── MainLayout.js
│   ├── MainLayout.css
│   └── index.js
├── AuthLayout/
│   ├── AuthLayout.js
│   ├── AuthLayout.css
│   └── index.js
└── EmptyLayout/
    ├── EmptyLayout.js
    ├── EmptyLayout.css
    └── index.js
```

#### src/config/
**功能**：配置文件

**说明**：
- 应用配置
- 环境配置
- 功能开关

**示例**：
```
src/config/
├── index.js              # 配置导出
├── app.js                # 应用配置
├── env.js                # 环境配置
└── features.js           # 功能开关
```

**配置示例**：
```javascript
// src/config/app.js
export default {
    appName: 'My App',
    version: '1.0.0',
    apiTimeout: 10000,
    pagination: {
        pageSize: 20,
        pageSizeOptions: [10, 20, 50, 100]
    },
    upload: {
        maxSize: 10 * 1024 * 1024, // 10MB
        allowedTypes: ['image/jpeg', 'image/png', 'image/gif']
    }
};
```

#### src/index.js
**功能**：入口文件

**说明**：
- 应用初始化
- 全局配置
- 挂载应用

**示例**：
```javascript
// src/index.js
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './assets/styles/main.css';
import './config';

ReactDOM.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>,
    document.getElementById('root')
);
```

### public/
**功能**：公共资源

**说明**：
- HTML 模板
- 静态资源
- PWA 配置

**示例**：
```
public/
├── index.html            # HTML 模板
├── favicon.ico           # 网站图标
├── manifest.json         # PWA 配置
└── robots.txt            # 爬虫配置
```

### test/
**功能**：测试目录

**说明**：
- 单元测试
- 集成测试
- 端到端测试

**示例**：
```
test/
├── unit/                 # 单元测试
│   ├── utils/
│   │   ├── date.test.js
│   │   └── string.test.js
│   └── components/
│       └── Button.test.js
├── integration/          # 集成测试
│   └── api.test.js
└── e2e/                  # 端到端测试
    └── user-flow.test.js
```

### scripts/
**功能**：构建脚本

**说明**：
- 构建脚本
- 开发脚本
- 部署脚本

**示例**：
```
scripts/
├── build.js              # 构建脚本
├── dev.js                # 开发脚本
└── deploy.js             # 部署脚本
```

### docs/
**功能**：文档目录

**说明**：
- API 文档
- 组件文档
- 开发文档

**示例**：
```
docs/
├── api.md                # API 文档
├── components.md         # 组件文档
└── development.md        # 开发文档
```

## 命名规范

### 文件命名
- 组件文件使用 PascalCase
- 工具文件使用 camelCase
- 常量文件使用 camelCase
- 测试文件以 `.test.js` 或 `.spec.js` 结尾

**示例**：
```
Button.js                # 组件文件
userService.js           # 服务文件
error.js                 # 常量文件
Button.test.js           # 测试文件
```

### 变量命名
- 使用 camelCase
- 常量使用 UPPER_SNAKE_CASE
- 私有变量以下划线开头

**示例**：
```javascript
const userName = 'John';           // 变量
const MAX_RETRY_COUNT = 3;         // 常量
const _privateVar = 'private';     // 私有变量
```

### 函数命名
- 使用 camelCase
- 动词开头
- 清晰表达功能

**示例**：
```javascript
function getUserById(id) {         // 函数
    return { id, name: 'John' };
}

const formatDate = (date) => {      // 箭头函数
    return date.toISOString();
};
```

### 类命名
- 使用 PascalCase
- 名词或名词短语

**示例**：
```javascript
class UserService {                 // 类
    constructor() {
        this.users = [];
    }

    getUsers() {
        return this.users;
    }
}
```

## 最佳实践

### 1. 代码风格
- 使用 ESLint 和 Prettier
- 统一代码风格
- 遵循 Airbnb JavaScript Style Guide

### 2. 错误处理
- 使用 try-catch 捕获错误
- 提供有意义的错误信息
- 记录错误日志

**示例**：
```javascript
async function fetchData() {
    try {
        const response = await api.get('/data');
        return response.data;
    } catch (error) {
        console.error('Failed to fetch data:', error);
        throw new Error('Failed to fetch data');
    }
}
```

### 3. 异步处理
- 使用 async/await
- 避免回调地狱
- 正确处理 Promise

**示例**：
```javascript
async function getUserData(userId) {
    const user = await getUserById(userId);
    const orders = await getOrdersByUserId(userId);
    return { user, orders };
}
```

### 4. 性能优化
- 避免不必要的渲染
- 使用 memoization
- 懒加载组件

**示例**：
```javascript
import { memo } from 'react';

const ExpensiveComponent = memo(({ data }) => {
    return <div>{data}</div>;
});
```

### 5. 安全考虑
- 避免使用 eval
- 防止 XSS 攻击
- 验证用户输入

**示例**：
```javascript
function sanitizeInput(input) {
    return input
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}
```

## 注意事项

1. **模块化设计**：合理划分模块，避免耦合度过高
2. **代码复用**：提取公共逻辑，提高代码复用性
3. **测试覆盖**：保持高测试覆盖率
4. **文档完整**：为公共 API 提供文档
5. **性能优化**：在必要时进行性能优化
6. **安全考虑**：注意代码安全性，避免常见漏洞
7. **版本管理**：使用语义化版本号
8. **依赖管理**：合理管理第三方依赖
