# React 代码架构规范

## 概述
本文档定义了 React 项目的代码架构规范，遵循 React 官方的最佳实践和社区约定。适用于 React 16+、React 18+ 项目。

## 项目结构

### 标准项目结构
```
project-name/
├── src/                   # 源代码目录
│   ├── assets/           # 静态资源
│   │   ├── images/       # 图片资源
│   │   ├── fonts/        # 字体资源
│   │   ├── styles/       # 样式文件
│   │   └── icons/        # 图标资源
│   ├── components/       # 可复用组件
│   │   ├── common/       # 通用组件
│   │   └── business/     # 业务组件
│   ├── pages/            # 页面组件
│   ├── hooks/            # 自定义 Hooks
│   ├── contexts/         # Context API
│   ├── services/         # 服务层
│   ├── api/              # API 接口
│   ├── utils/            # 工具函数
│   ├── constants/        # 常量定义
│   ├── types/            # 类型定义（TypeScript）
│   ├── store/            # 状态管理（Redux/Zustand）
│   ├── router/           # 路由配置
│   ├── layouts/          # 布局组件
│   ├── config/           # 配置文件
│   ├── App.tsx           # 根组件
│   ├── App.css           # 根样式
│   ├── index.tsx         # 入口文件
│   └── index.css         # 全局样式
├── public/               # 公共资源
│   ├── index.html        # HTML 模板
│   ├── favicon.ico       # 网站图标
│   ├── manifest.json     # PWA 配置
│   └── robots.txt        # 爬虫配置
├── test/                  # 测试目录
│   ├── unit/             # 单元测试
│   ├── integration/      # 集成测试
│   └── e2e/              # 端到端测试
├── docs/                  # 文档目录
├── build/                 # 构建产物（自动生成）
├── node_modules/          # 依赖包（自动生成）
├── package.json           # 项目配置
├── package-lock.json      # 依赖锁定
├── tsconfig.json          # TypeScript 配置
├── vite.config.ts         # Vite 配置
├── webpack.config.js      # Webpack 配置
├── .eslintrc.js           # ESLint 配置
├── .prettierrc.js         # Prettier 配置
├── .gitignore             # Git 忽略文件
├── .env                   # 环境变量
├── .env.development       # 开发环境变量
├── .env.production        # 生产环境变量
└── README.md              # 项目说明文档
```

## 目录说明

### src/assets/
**功能**：静态资源

**说明**：
- images/：图片资源（PNG、JPG、SVG 等）
- fonts/：字体资源
- styles/：全局样式文件
- icons/：图标资源

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
│   ├── variables.css
│   └── reset.css
└── icons/
    ├── user.svg
    └── settings.svg
```

### src/components/
**功能**：可复用组件

**说明**：
- common/：通用组件（按钮、输入框、模态框等）
- business/：业务组件（用户卡片、产品列表等）

**示例**：
```
src/components/
├── common/               # 通用组件
│   ├── Button/
│   │   ├── Button.tsx
│   │   ├── Button.test.tsx
│   │   ├── Button.module.css
│   │   └── index.ts
│   ├── Input/
│   │   ├── Input.tsx
│   │   ├── Input.test.tsx
│   │   ├── Input.module.css
│   │   └── index.ts
│   ├── Modal/
│   │   ├── Modal.tsx
│   │   ├── Modal.test.tsx
│   │   ├── Modal.module.css
│   │   └── index.ts
│   └── Table/
│       ├── Table.tsx
│       ├── Table.test.tsx
│       ├── Table.module.css
│       └── index.ts
└── business/             # 业务组件
    ├── UserCard/
    │   ├── UserCard.tsx
    │   ├── UserCard.test.tsx
    │   ├── UserCard.module.css
    │   └── index.ts
    ├── ProductList/
    │   ├── ProductList.tsx
    │   ├── ProductList.test.tsx
    │   ├── ProductList.module.css
    │   └── index.ts
    └── OrderForm/
        ├── OrderForm.tsx
        ├── OrderForm.test.tsx
        ├── OrderForm.module.css
        └── index.ts
```

**组件示例**：
```tsx
// src/components/common/Button/Button.tsx
import React from 'react';
import styles from './Button.module.css';

export interface ButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  type?: 'primary' | 'secondary' | 'danger';
  disabled?: false;
}

const Button: React.FC<ButtonProps> = ({
  children,
  onClick,
  type = 'primary',
  disabled = false
}) => {
  return (
    <button
      className={`${styles.button} ${styles[`button--${type}`]} ${disabled ? styles['button--disabled'] : ''}`}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  );
};

export default Button;
```

```css
/* src/components/common/Button/Button.module.css */
.button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.button--primary {
  background-color: #1890ff;
  color: white;
}

.button--primary:hover {
  background-color: #40a9ff;
}

.button--secondary {
  background-color: #f0f0f0;
  color: #333;
}

.button--secondary:hover {
  background-color: #d9d9d9;
}

.button--danger {
  background-color: #ff4d4f;
  color: white;
}

.button--danger:hover {
  background-color: #ff7875;
}

.button--disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
```

### src/pages/
**功能**：页面组件

**说明**：
- 路由对应的页面组件
- 页面级别的业务逻辑
- 页面布局

**示例**：
```
src/pages/
├── Home/
│   ├── Home.tsx
│   ├── Home.test.tsx
│   ├── Home.module.css
│   └── components/
│       └── HomeBanner.tsx
├── User/
│   ├── UserList.tsx
│   ├── UserList.test.tsx
│   ├── UserList.module.css
│   ├── UserDetail.tsx
│   ├── UserDetail.test.tsx
│   ├── UserDetail.module.css
│   └── components/
│       ├── UserFilter.tsx
│       └── UserTable.tsx
├── Product/
│   ├── ProductList.tsx
│   ├── ProductList.test.tsx
│   ├── ProductList.module.css
│   ├── ProductDetail.tsx
│   ├── ProductDetail.test.tsx
│   ├── ProductDetail.module.css
│   └── components/
│       ├── ProductCard.tsx
│       └── ProductFilter.tsx
└── Order/
    ├── OrderList.tsx
    ├── OrderList.test.tsx
    ├── OrderList.module.css
    ├── OrderDetail.tsx
    ├── OrderDetail.test.tsx
    ├── OrderDetail.module.css
    └── components/
        └── OrderForm.tsx
```

### src/hooks/
**功能**：自定义 Hooks

**说明**：
- 可复用的逻辑
- 响应式状态
- 副作用处理

**示例**：
```
src/hooks/
├── useUser.ts            # 用户逻辑
├── useProduct.ts         # 产品逻辑
├── usePagination.ts      # 分页逻辑
├── useDebounce.ts        # 防抖逻辑
└── useLocalStorage.ts    # 本地存储逻辑
```

**Hook 示例**：
```tsx
// src/hooks/useUser.ts
import { useState, useEffect } from 'react';
import api from '@/api';

export interface User {
  id: number;
  name: string;
  email: string;
}

export function useUser() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const response = await api.getUsers();
      setUsers(response.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  return {
    users,
    loading,
    error,
    fetchUsers,
    userCount: users.length
  };
}
```

### src/contexts/
**功能**：Context API

**说明**：
- 全局状态管理
- 主题管理
- 用户认证

**示例**：
```
src/contexts/
├── AuthContext.tsx       # 认证上下文
├── ThemeContext.tsx      # 主题上下文
└── AppContext.tsx        # 应用上下文
```

**Context 示例**：
```tsx
// src/contexts/AuthContext.tsx
import React, { createContext, useContext, useState, ReactNode } from 'react';

export interface User {
  id: number;
  name: string;
  email: string;
}

interface AuthContextType {
  user: User | null;
  login: (user: User) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);

  const login = (userData: User) => {
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData));
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('user');
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
```

### src/services/
**功能**：服务层

**说明**：
- 业务逻辑
- 数据处理
- 状态管理

**示例**：
```
src/services/
├── userService.ts        # 用户服务
├── productService.ts     # 产品服务
├── orderService.ts       # 订单服务
└── authService.ts        # 认证服务
```

**服务示例**：
```tsx
// src/services/userService.ts
import api from '@/api';

export interface User {
  id: number;
  name: string;
  email: string;
  age: number;
}

export const userService = {
  async getUsers(params?: any) {
    const response = await api.get('/users', { params });
    return response.data;
  },

  async getUserById(id: number) {
    const response = await api.get(`/users/${id}`);
    return response.data;
  },

  async createUser(data: Partial<User>) {
    const response = await api.post('/users', data);
    return response.data;
  },

  async updateUser(id: number, data: Partial<User>) {
    const response = await api.put(`/users/${id}`, data);
    return response.data;
  },

  async deleteUser(id: number) {
    const response = await api.delete(`/users/${id}`);
    return response.data;
  }
};
```

### src/api/
**功能**：API 接口

**说明**：
- HTTP 请求封装
- 接口定义
- 错误处理

**示例**：
```
src/api/
├── index.ts              # API 实例
├── user.ts               # 用户 API
├── product.ts            # 产品 API
└── order.ts              # 订单 API
```

**API 示例**：
```tsx
// src/api/index.ts
import axios, { AxiosInstance, AxiosError } from 'axios';

const api: AxiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL || 'http://localhost:3000/api',
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
  (error: AxiosError) => {
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

### src/utils/
**功能**：工具函数

**说明**：
- 通用工具函数
- 纯函数
- 无副作用

**示例**：
```
src/utils/
├── date.ts               # 日期工具
├── string.ts             # 字符串工具
├── array.ts              # 数组工具
├── object.ts             # 对象工具
├── validation.ts         # 验证工具
└── storage.ts            # 存储工具
```

### src/constants/
**功能**：常量定义

**说明**：
- 避免魔法数字
- 统一管理常量
- 便于维护

**示例**：
```
src/constants/
├── index.ts              # 常量导出
├── error.ts              # 错误码
├── status.ts             # 状态码
└── config.ts             # 配置常量
```

**常量示例**：
```tsx
// src/constants/error.ts
export const ERROR_CODES = {
  SUCCESS: 0,
  INVALID_PARAMS: 1001,
  UNAUTHORIZED: 1002,
  FORBIDDEN: 1003,
  NOT_FOUND: 1004,
  SERVER_ERROR: 1005
} as const;

export const ERROR_MESSAGES: Record<number, string> = {
  [ERROR_CODES.SUCCESS]: 'Success',
  [ERROR_CODES.INVALID_PARAMS]: 'Invalid parameters',
  [ERROR_CODES.UNAUTHORIZED]: 'Unauthorized',
  [ERROR_CODES.FORBIDDEN]: 'Access denied',
  [ERROR_CODES.NOT_FOUND]: 'Resource not found',
  [ERROR_CODES.SERVER_ERROR]: 'Server error'
};
```

### src/types/
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
```tsx
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

### src/store/
**功能**：状态管理（Redux/Zustand）

**说明**：
- 全局状态
- 状态更新逻辑
- 状态派生

**示例**：
```
src/store/
├── index.ts              # Store 配置
├── slices/              # Redux Slices
│   ├── userSlice.ts
│   ├── productSlice.ts
│   └── orderSlice.ts
└── hooks.ts              # Store Hooks
```

**Redux Toolkit 示例**：
```tsx
// src/store/slices/userSlice.ts
import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import api from '@/api';

export interface User {
  id: number;
  name: string;
  email: string;
}

interface UserState {
  users: User[];
  currentUser: User | null;
  loading: boolean;
  error: string | null;
}

const initialState: UserState = {
  users: [],
  currentUser: null,
  loading: false,
  error: null
};

export const fetchUsers = createAsyncThunk(
  'user/fetchUsers',
  async () => {
    const response = await api.getUsers();
    return response.data;
  }
);

export const fetchUserById = createAsyncThunk(
  'user/fetchUserById',
  async (id: number) => {
    const response = await api.getUserById(id);
    return response.data;
  }
);

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setCurrentUser: (state, action: PayloadAction<User>) => {
      state.currentUser = action.payload;
    },
    clearCurrentUser: (state) => {
      state.currentUser = null;
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchUsers.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchUsers.fulfilled, (state, action) => {
        state.loading = false;
        state.users = action.payload;
      })
      .addCase(fetchUsers.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch users';
      })
      .addCase(fetchUserById.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchUserById.fulfilled, (state, action) => {
        state.loading = false;
        state.currentUser = action.payload;
      })
      .addCase(fetchUserById.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch user';
      });
  }
});

export const { setCurrentUser, clearCurrentUser } = userSlice.actions;
export default userSlice.reducer;
```

### src/router/
**功能**：路由配置

**说明**：
- 路由定义
- 路由守卫
- 路由参数

**示例**：
```
src/router/
├── index.tsx             # 路由配置
├── routes.tsx            # 路由定义
└── guards.tsx            # 路由守卫
```

**路由示例**：
```tsx
// src/router/index.tsx
import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ProtectedRoute } from './guards';
import Home from '@/pages/Home';
import UserList from '@/pages/User/UserList';
import UserDetail from '@/pages/User/UserDetail';
import Login from '@/pages/Login';

const Router: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route
          path="/users"
          element={
            <ProtectedRoute>
              <UserList />
            </ProtectedRoute>
          }
        />
        <Route
          path="/users/:id"
          element={
            <ProtectedRoute>
              <UserDetail />
            </ProtectedRoute>
          }
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
```

### src/layouts/
**功能**：布局组件

**说明**：
- 页面布局
- 通用布局
- 布局切换

**示例**：
```
src/layouts/
├── MainLayout/
│   ├── MainLayout.tsx
│   ├── MainLayout.module.css
│   └── index.ts
├── AuthLayout/
│   ├── AuthLayout.tsx
│   ├── AuthLayout.module.css
│   └── index.ts
└── EmptyLayout/
    ├── EmptyLayout.tsx
    ├── EmptyLayout.module.css
    └── index.ts
```

### src/config/
**功能**：配置文件

**说明**：
- 应用配置
- 环境配置
- 功能开关

**示例**：
```
src/config/
├── index.ts              # 配置导出
├── app.ts                # 应用配置
├── env.ts                # 环境配置
└── features.ts           # 功能开关
```

## 命名规范

### 组件命名
- 组件名使用 PascalCase
- 多个单词组合
- 避免使用保留字

**示例**：
```
UserCard.tsx             # 正确
user-card.tsx            # 错误
UserCardComponent.tsx    # 冗余
```

### 文件命名
- 组件文件使用 PascalCase
- 工具文件使用 camelCase
- 常量文件使用 camelCase
- 测试文件以 `.test.tsx` 或 `.spec.tsx` 结尾

**示例**：
```
UserCard.tsx             # 组件文件
userService.ts           # 服务文件
error.ts                 # 常量文件
UserCard.test.tsx        # 测试文件
```

### 变量命名
- 使用 camelCase
- 常量使用 UPPER_SNAKE_CASE

**示例**：
```tsx
const userName = 'John';           // 变量
const MAX_RETRY_COUNT = 3;         // 常量
```

### 函数命名
- 使用 camelCase
- 动词开头
- 清晰表达功能

**示例**：
```tsx
function getUserById(id: number): User {         // 函数
    return { id, name: 'John' };
}

const formatDate = (date: Date): string => {      // 箭头函数
    return date.toISOString();
};
```

## 最佳实践

### 1. 组件设计
- 单一职责原则
- 组件复用性
- Props 类型定义
- 事件处理

**示例**：
```tsx
interface UserCardProps {
  user: User;
  onDelete?: (id: number) => void;
}

const UserCard: React.FC<UserCardProps> = ({ user, onDelete }) => {
  const handleDelete = () => {
    if (onDelete) {
      onDelete(user.id);
    }
  };

  return (
    <div className="user-card">
      <h3>{user.name}</h3>
      <p>{user.email}</p>
      <button onClick={handleDelete}>Delete</button>
    </div>
  );
};
```

### 2. Hooks 使用
- 自定义 Hooks 提取逻辑
- 遵循 Hooks 规则
- 合理使用 useEffect

**示例**：
```tsx
// 自定义 Hook
function useFetch<T>(url: string) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const response = await fetch(url);
        const json = await response.json();
        setData(json);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [url]);

  return { data, loading, error };
}
```

### 3. 状态管理
- 合理选择状态管理方案
- 避免过度使用 Context
- 使用 Redux Toolkit 简化 Redux

### 4. 性能优化
- 使用 React.memo 避免不必要的渲染
- 使用 useMemo 缓存计算结果
- 使用 useCallback 缓存函数
- 懒加载组件和路由

**示例**：
```tsx
import React, { memo, useMemo, useCallback } from 'react';

const ExpensiveComponent = memo(({ data }: { data: any[] }) => {
  const processedData = useMemo(() => {
    return data.map(item => ({ ...item, processed: true }));
  }, [data]);

  const handleClick = useCallback((id: number) => {
    console.log('Clicked:', id);
  }, []);

  return (
    <div>
      {processedData.map(item => (
        <div key={item.id} onClick={() => handleClick(item.id)}>
          {item.name}
        </div>
      ))}
    </div>
  );
});
```

### 5. 错误处理
- 使用 Error Boundary 捕获错误
- 提供友好的错误提示
- 记录错误日志

**示例**：
```tsx
class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean; error: Error | null }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return <div>Something went wrong.</div>;
    }

    return this.props.children;
  }
}
```

### 6. 代码风格
- 使用 ESLint 和 Prettier
- 统一代码风格
- 遵循 Airbnb React/JSX Style Guide

## 注意事项

1. **组件复用**：提取公共逻辑，提高组件复用性
2. **性能优化**：避免不必要的渲染和计算
3. **测试覆盖**：保持高测试覆盖率
4. **文档完整**：为公共组件提供文档
5. **安全考虑**：注意 XSS 攻击，使用 dangerouslySetInnerHTML 时要谨慎
6. **版本管理**：使用语义化版本号
7. **依赖管理**：合理管理第三方依赖
8. **响应式设计**：确保组件在不同设备上正常显示
