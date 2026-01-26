# Vue 代码架构规范

## 概述
本文档定义了 Vue.js 项目的代码架构规范，遵循 Vue.js 官方的最佳实践和社区约定。适用于 Vue 2 和 Vue 3 项目。

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
│   ├── views/            # 页面视图
│   ├── router/           # 路由配置
│   ├── store/            # 状态管理（Vuex/Pinia）
│   ├── api/              # API 接口
│   ├── utils/            # 工具函数
│   ├── composables/      # 组合式函数（Vue 3）
│   ├── directives/       # 自定义指令
│   ├── filters/          # 过滤器（Vue 2）
│   ├── mixins/           # 混入（Vue 2）
│   ├── plugins/          # 插件
│   ├── config/           # 配置文件
│   ├── types/            # 类型定义（TypeScript）
│   ├── App.vue           # 根组件
│   └── main.js           # 入口文件
├── public/               # 公共资源
│   ├── index.html        # HTML 模板
│   ├── favicon.ico       # 网站图标
│   └── robots.txt        # 爬虫配置
├── test/                  # 测试目录
│   ├── unit/             # 单元测试
│   ├── e2e/              # 端到端测试
│   └── snapshots/        # 快照测试
├── docs/                  # 文档目录
├── dist/                  # 构建产物（自动生成）
├── node_modules/          # 依赖包（自动生成）
├── package.json           # 项目配置
├── package-lock.json      # 依赖锁定
├── vue.config.js          # Vue CLI 配置
├── vite.config.js         # Vite 配置
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
│   ├── BaseButton/
│   │   ├── BaseButton.vue
│   │   └── index.js
│   ├── BaseInput/
│   │   ├── BaseInput.vue
│   │   └── index.js
│   ├── BaseModal/
│   │   ├── BaseModal.vue
│   │   └── index.js
│   └── BaseTable/
│       ├── BaseTable.vue
│       └── index.js
└── business/             # 业务组件
    ├── UserCard/
    │   ├── UserCard.vue
    │   └── index.js
    ├── ProductList/
    │   ├── ProductList.vue
    │   └── index.js
    └── OrderForm/
        ├── OrderForm.vue
        └── index.js
```

**组件示例**：
```vue
<!-- src/components/common/BaseButton/BaseButton.vue -->
<template>
  <button
    :class="['base-button', `base-button--${type}`, { 'base-button--disabled': disabled }]"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot></slot>
  </button>
</template>

<script>
export default {
  name: 'BaseButton',
  props: {
    type: {
      type: String,
      default: 'primary',
      validator: (value) => ['primary', 'secondary', 'danger'].includes(value)
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    handleClick(event) {
      if (!this.disabled) {
        this.$emit('click', event);
      }
    }
  }
}
</script>

<style scoped>
.base-button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.base-button--primary {
  background-color: #1890ff;
  color: white;
}

.base-button--secondary {
  background-color: #f0f0f0;
  color: #333;
}

.base-button--danger {
  background-color: #ff4d4f;
  color: white;
}

.base-button--disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
```

### src/views/
**功能**：页面视图

**说明**：
- 路由对应的页面组件
- 页面级别的业务逻辑
- 页面布局

**示例**：
```
src/views/
├── Home/
│   ├── Home.vue
│   └── components/
│       └── HomeBanner.vue
├── User/
│   ├── UserList.vue
│   ├── UserDetail.vue
│   └── components/
│       ├── UserFilter.vue
│       └── UserTable.vue
├── Product/
│   ├── ProductList.vue
│   ├── ProductDetail.vue
│   └── components/
│       ├── ProductCard.vue
│       └── ProductFilter.vue
└── Order/
    ├── OrderList.vue
    ├── OrderDetail.vue
    └── components/
        └── OrderForm.vue
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
├── index.js              # 路由配置
├── routes.js             # 路由定义
└── guards.js             # 路由守卫
```

**路由示例**：
```javascript
// src/router/index.js
import Vue from 'vue';
import VueRouter from 'vue-router';
import routes from './routes';
import { beforeEach } from './guards';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
});

router.beforeEach(beforeEach);

export default router;
```

```javascript
// src/router/routes.js
export default [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home/Home.vue'),
    meta: {
      title: '首页',
      requiresAuth: false
    }
  },
  {
    path: '/users',
    name: 'UserList',
    component: () => import('@/views/User/UserList.vue'),
    meta: {
      title: '用户列表',
      requiresAuth: true
    }
  },
  {
    path: '/users/:id',
    name: 'UserDetail',
    component: () => import('@/views/User/UserDetail.vue'),
    meta: {
      title: '用户详情',
      requiresAuth: true
    }
  }
];
```

### src/store/
**功能**：状态管理（Vuex/Pinia）

**说明**：
- 全局状态管理
- 状态模块化
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

**Vuex 示例**：
```javascript
// src/store/modules/user.js
const state = {
  users: [],
  currentUser: null,
  loading: false,
  error: null
};

const mutations = {
  SET_USERS(state, users) {
    state.users = users;
  },
  SET_CURRENT_USER(state, user) {
    state.currentUser = user;
  },
  SET_LOADING(state, loading) {
    state.loading = loading;
  },
  SET_ERROR(state, error) {
    state.error = error;
  }
};

const actions = {
  async fetchUsers({ commit }) {
    commit('SET_LOADING', true);
    try {
      const response = await api.getUsers();
      commit('SET_USERS', response.data);
    } catch (error) {
      commit('SET_ERROR', error.message);
    } finally {
      commit('SET_LOADING', false);
    }
  },
  async fetchUserById({ commit }, id) {
    commit('SET_LOADING', true);
    try {
      const response = await api.getUserById(id);
      commit('SET_CURRENT_USER', response.data);
    } catch (error) {
      commit('SET_ERROR', error.message);
    } finally {
      commit('SET_LOADING', false);
    }
  }
};

const getters = {
  users: (state) => state.users,
  currentUser: (state) => state.currentUser,
  loading: (state) => state.loading,
  error: (state) => state.error
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
};
```

**Pinia 示例**：
```javascript
// src/store/user.js
import { defineStore } from 'pinia';
import api from '@/api';

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [],
    currentUser: null,
    loading: false,
    error: null
  }),
  getters: {
    users: (state) => state.users,
    currentUser: (state) => state.currentUser,
    loading: (state) => state.loading,
    error: (state) => state.error
  },
  actions: {
    async fetchUsers() {
      this.loading = true;
      try {
        const response = await api.getUsers();
        this.users = response.data;
      } catch (error) {
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    },
    async fetchUserById(id) {
      this.loading = true;
      try {
        const response = await api.getUserById(id);
        this.currentUser = response.data;
      } catch (error) {
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    }
  }
});
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
  baseURL: process.env.VUE_APP_API_BASE_URL || 'http://localhost:3000/api',
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

### src/utils/
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

### src/composables/
**功能**：组合式函数（Vue 3）

**说明**：
- 可复用的逻辑
- 响应式状态
- 生命周期钩子

**示例**：
```
src/composables/
├── useUser.js            # 用户逻辑
├── useProduct.js         # 产品逻辑
└── usePagination.js      # 分页逻辑
```

**组合式函数示例**：
```javascript
// src/composables/useUser.js
import { ref, computed } from 'vue';
import api from '@/api';

export function useUser() {
  const users = ref([]);
  const loading = ref(false);
  const error = ref(null);

  const fetchUsers = async () => {
    loading.value = true;
    try {
      const response = await api.getUsers();
      users.value = response.data;
    } catch (err) {
      error.value = err.message;
    } finally {
      loading.value = false;
    }
  };

  const userCount = computed(() => users.value.length);

  return {
    users,
    loading,
    error,
    fetchUsers,
    userCount
  };
}
```

### src/directives/
**功能**：自定义指令

**说明**：
- 自定义指令
- 指令钩子
- 指令参数

**示例**：
```
src/directives/
├── index.js              # 指令注册
├── focus.js              # 焦点指令
├── lazy.js               # 懒加载指令
└── permission.js         # 权限指令
```

**指令示例**：
```javascript
// src/directives/permission.js
import store from '@/store';

export default {
  inserted(el, binding) {
    const { value } = binding;
    const permissions = store.getters['user/permissions'];

    if (value && !permissions.includes(value)) {
      el.parentNode && el.parentNode.removeChild(el);
    }
  }
};
```

### src/plugins/
**功能**：插件

**说明**：
- Vue 插件
- 第三方库集成
- 全局配置

**示例**：
```
src/plugins/
├── index.js              # 插件注册
├── element.js            # Element UI
├── ant-design.js         # Ant Design Vue
└── i18n.js               # 国际化
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
├── index.js              # 配置导出
├── app.js                # 应用配置
├── env.js                # 环境配置
└── features.js           # 功能开关
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

## 命名规范

### 组件命名
- 组件名使用 PascalCase
- 多个单词组合
- 避免使用保留字

**示例**：
```
UserCard.vue             # 正确
user-card.vue            # 错误
UserCardComponent.vue    # 冗余
```

### 文件命名
- 组件文件使用 PascalCase
- 工具文件使用 camelCase
- 常量文件使用 camelCase

**示例**：
```
UserCard.vue             # 组件文件
userService.js           # 服务文件
error.js                 # 常量文件
```

### 变量命名
- 使用 camelCase
- 常量使用 UPPER_SNAKE_CASE

**示例**：
```javascript
const userName = 'John';           // 变量
const MAX_RETRY_COUNT = 3;         // 常量
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
```

## 最佳实践

### 1. 组件设计
- 单一职责原则
- 组件复用性
- Props 验证
- 事件命名

**示例**：
```vue
<script>
export default {
  name: 'UserCard',
  props: {
    user: {
      type: Object,
      required: true,
      validator: (value) => {
        return value.id && value.name;
      }
    }
  },
  methods: {
    handleDelete() {
      this.$emit('delete', this.user.id);
    }
  }
}
</script>
```

### 2. 状态管理
- 合理划分模块
- 避免过度使用
- 使用 getters 派生状态

### 3. 路由管理
- 懒加载路由
- 路由守卫
- 路由元信息

### 4. 性能优化
- 使用 v-show 替代 v-if（频繁切换）
- 使用 key 优化列表渲染
- 使用计算属性缓存
- 使用 keep-alive 缓存组件

**示例**：
```vue
<template>
  <div>
    <div v-for="item in items" :key="item.id">
      {{ item.name }}
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    filteredItems() {
      return this.items.filter(item => item.active);
    }
  }
}
</script>
```

### 5. 代码风格
- 使用 ESLint 和 Prettier
- 统一代码风格
- 遵循 Vue 官方风格指南

## 注意事项

1. **组件复用**：提取公共逻辑，提高组件复用性
2. **性能优化**：避免不必要的渲染和计算
3. **测试覆盖**：保持高测试覆盖率
4. **文档完整**：为公共组件提供文档
5. **安全考虑**：注意 XSS 攻击，使用 v-html 时要谨慎
6. **版本管理**：使用语义化版本号
7. **依赖管理**：合理管理第三方依赖
8. **响应式设计**：确保组件在不同设备上正常显示
