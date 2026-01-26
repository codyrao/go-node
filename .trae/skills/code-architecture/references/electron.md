# Electron 代码架构规范

## 概述
本文档定义了 Electron 项目的代码架构规范，遵循 Electron 官方的最佳实践和安全指南。Electron 是一个使用 JavaScript、HTML 和 CSS 构建跨平台桌面应用程序的框架。

## 项目结构

### 标准项目结构
```
project-name/
├── src/                   # 源代码目录
│   ├── main/              # 主进程代码
│   │   ├── index.ts       # 主进程入口
│   │   ├── window.ts      # 窗口管理
│   │   ├── menu.ts       # 菜单配置
│   │   ├── tray.ts       # 系统托盘
│   │   ├── ipc/          # IPC 通信
│   │   │   ├── handlers/ # IPC 处理器
│   │   │   └── channels/ # IPC 通道定义
│   │   ├── services/     # 主进程服务
│   │   ├── utils/        # 工具函数
│   │   └── config/       # 配置文件
│   ├── renderer/          # 渲染进程代码
│   │   ├── index.html    # 渲染进程入口
│   │   ├── index.ts      # 渲染进程入口
│   │   ├── assets/       # 静态资源
│   │   │   ├── images/
│   │   │   ├── fonts/
│   │   │   ├── styles/
│   │   │   └── icons/
│   │   ├── components/   # 组件
│   │   ├── pages/        # 页面
│   │   ├── store/        # 状态管理
│   │   ├── hooks/        # 自定义 Hooks
│   │   ├── utils/        # 工具函数
│   │   ├── api/          # API 接口
│   │   ├── constants/    # 常量定义
│   │   ├── types/        # 类型定义
│   │   └── preload/      # 预加载脚本
│   │       └── index.ts
├── resources/             # 资源文件
│   ├── icons/            # 应用图标
│   │   ├── icon.png
│   │   ├── icon.ico
│   │   └── icon.icns
│   └── updater/          # 更新资源
├── build/                 # 构建配置
│   ├── builder.config.ts  # Electron Builder 配置
│   └── webpack.config.ts  # Webpack 配置
├── scripts/               # 构建脚本
│   ├── build.js          # 构建脚本
│   ├── dev.js            # 开发脚本
│   └── release.js        # 发布脚本
├── test/                  # 测试目录
│   ├── unit/             # 单元测试
│   ├── integration/      # 集成测试
│   └── e2e/              # 端到端测试
├── docs/                  # 文档目录
├── dist/                  # 构建产物（自动生成）
├── out/                   # 打包产物（自动生成）
├── node_modules/          # 依赖包（自动生成）
├── package.json           # 项目配置
├── package-lock.json      # 依赖锁定
├── tsconfig.json          # TypeScript 配置
├── tsconfig.main.json     # 主进程 TypeScript 配置
├── tsconfig.renderer.json # 渲染进程 TypeScript 配置
├── .eslintrc.js           # ESLint 配置
├── .prettierrc.js         # Prettier 配置
├── .gitignore             # Git 忽略文件
├── .env                   # 环境变量
└── README.md              # 项目说明文档
```

## 目录说明

### src/main/
**功能**：主进程代码

**说明**：
- 主进程入口
- 窗口管理
- 菜单配置
- 系统托盘
- IPC 通信
- 主进程服务

**示例**：
```
src/main/
├── index.ts              # 主进程入口
├── window.ts             # 窗口管理
├── menu.ts               # 菜单配置
├── tray.ts               # 系统托盘
├── ipc/
│   ├── handlers/         # IPC 处理器
│   │   ├── app.handler.ts
│   │   ├── window.handler.ts
│   │   └── file.handler.ts
│   └── channels/         # IPC 通道定义
│       └── index.ts
├── services/             # 主进程服务
│   ├── app.service.ts
│   ├── window.service.ts
│   └── file.service.ts
├── utils/                # 工具函数
│   ├── logger.ts
│   └── path.ts
└── config/               # 配置文件
    └── app.config.ts
```

**主进程入口示例**：
```typescript
// src/main/index.ts
import { app, BrowserWindow } from 'electron';
import { createWindow } from './window';
import { createMenu } from './menu';
import { createTray } from './tray';
import { registerIpcHandlers } from './ipc/handlers';
import { logger } from './utils/logger';

let mainWindow: BrowserWindow | null = null;

app.whenReady().then(() => {
  mainWindow = createWindow();
  createMenu();
  createTray();
  registerIpcHandlers();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      mainWindow = createWindow();
    }
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', () => {
  logger.info('Application is quitting');
});
```

**窗口管理示例**：
```typescript
// src/main/window.ts
import { BrowserWindow, app } from 'electron';
import * as path from 'path';
import { isDev } from './utils/env';

export function createWindow(): BrowserWindow {
  const mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: path.join(__dirname, '../preload/index.js'),
      nodeIntegration: false,
      contextIsolation: true,
      sandbox: true
    },
    icon: path.join(__dirname, '../../resources/icons/icon.png')
  });

  if (isDev) {
    mainWindow.loadURL('http://localhost:3000');
    mainWindow.webContents.openDevTools();
  } else {
    mainWindow.loadFile(path.join(__dirname, '../renderer/index.html'));
  }

  return mainWindow;
}

export function getMainWindow(): BrowserWindow | null {
  return BrowserWindow.getAllWindows()[0] || null;
}
```

**IPC 处理器示例**：
```typescript
// src/main/ipc/handlers/app.handler.ts
import { ipcMain } from 'electron';
import { IPC_CHANNELS } from '../channels';
import { logger } from '../../utils/logger';

export function registerAppHandlers() {
  ipcMain.handle(IPC_CHANNELS.APP.GET_VERSION, () => {
    return app.getVersion();
  });

  ipcMain.handle(IPC_CHANNELS.APP.GET_PATH, (event, name) => {
    return app.getPath(name);
  });

  ipcMain.on(IPC_CHANNELS.APP.QUIT, () => {
    logger.info('Quit requested from renderer');
    app.quit();
  });
}
```

### src/renderer/
**功能**：渲染进程代码

**说明**：
- 渲染进程入口
- 静态资源
- 组件
- 页面
- 状态管理
- API 接口

**示例**：
```
src/renderer/
├── index.html            # 渲染进程入口
├── index.ts             # 渲染进程入口
├── assets/              # 静态资源
│   ├── images/
│   ├── fonts/
│   ├── styles/
│   └── icons/
├── components/          # 组件
│   ├── common/
│   └── business/
├── pages/               # 页面
│   ├── Home/
│   ├── Settings/
│   └── About/
├── store/               # 状态管理
├── hooks/               # 自定义 Hooks
├── utils/               # 工具函数
├── api/                 # API 接口
├── constants/           # 常量定义
├── types/               # 类型定义
└── preload/             # 预加载脚本
    └── index.ts
```

**预加载脚本示例**：
```typescript
// src/renderer/preload/index.ts
import { contextBridge, ipcRenderer } from 'electron';
import { IPC_CHANNELS } from '../../main/ipc/channels';

const electronAPI = {
  app: {
    getVersion: () => ipcRenderer.invoke(IPC_CHANNELS.APP.GET_VERSION),
    getPath: (name: string) => ipcRenderer.invoke(IPC_CHANNELS.APP.GET_PATH, name),
    quit: () => ipcRenderer.send(IPC_CHANNELS.APP.QUIT)
  },
  window: {
    minimize: () => ipcRenderer.send(IPC_CHANNELS.WINDOW.MINIMIZE),
    maximize: () => ipcRenderer.send(IPC_CHANNELS.WINDOW.MAXIMIZE),
    close: () => ipcRenderer.send(IPC_CHANNELS.WINDOW.CLOSE)
  },
  file: {
    selectFile: () => ipcRenderer.invoke(IPC_CHANNELS.FILE.SELECT),
    readFile: (path: string) => ipcRenderer.invoke(IPC_CHANNELS.FILE.READ, path),
    writeFile: (path: string, content: string) => ipcRenderer.invoke(IPC_CHANNELS.FILE.WRITE, path, content)
  }
};

contextBridge.exposeInMainWorld('electronAPI', electronAPI);
```

**渲染进程入口示例**：
```typescript
// src/renderer/index.ts
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './assets/styles/main.css';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
```

### resources/
**功能**：资源文件

**说明**：
- 应用图标
- 更新资源
- 其他资源文件

**示例**：
```
resources/
├── icons/               # 应用图标
│   ├── icon.png         # PNG 图标
│   ├── icon.ico         # Windows 图标
│   └── icon.icns        # macOS 图标
├── updater/             # 更新资源
│   ├── latest.yml       # 更新配置
│   └── update.exe       # Windows 更新程序
└── splash.png           # 启动画面
```

### build/
**功能**：构建配置

**说明**：
- Electron Builder 配置
- Webpack 配置
- 其他构建配置

**示例**：
```
build/
├── builder.config.ts     # Electron Builder 配置
├── webpack.config.ts    # Webpack 配置
└── webpack.main.config.ts # 主进程 Webpack 配置
```

**Electron Builder 配置示例**：
```typescript
// build/builder.config.ts
import { Configuration } from 'electron-builder';

const config: Configuration = {
  appId: 'com.example.app',
  productName: 'My App',
  copyright: 'Copyright © 2024',
  directories: {
    output: 'dist',
    buildResources: 'resources'
  },
  files: [
    'dist/**/*',
    'package.json'
  ],
  win: {
    target: ['nsis', 'portable'],
    icon: 'resources/icons/icon.ico'
  },
  mac: {
    target: ['dmg', 'zip'],
    icon: 'resources/icons/icon.icns',
    category: 'public.app-category.productivity'
  },
  linux: {
    target: ['AppImage', 'deb', 'rpm'],
    icon: 'resources/icons/icon.png',
    category: 'Utility'
  },
  nsis: {
    oneClick: false,
    allowToChangeInstallationDirectory: true,
    createDesktopShortcut: true,
    createStartMenuShortcut: true
  },
  publish: {
    provider: 'github',
    owner: 'your-username',
    repo: 'your-repo'
  }
};

export default config;
```

### scripts/
**功能**：构建脚本

**说明**：
- 构建脚本
- 开发脚本
- 发布脚本

**示例**：
```
scripts/
├── build.js             # 构建脚本
├── dev.js               # 开发脚本
└── release.js           # 发布脚本
```

## 命名规范

### 文件命名
- 主进程文件使用 kebab-case
- 渲染进程组件使用 PascalCase
- 工具文件使用 camelCase
- 测试文件以 `.spec.ts` 或 `.test.ts` 结尾

**示例**：
```
window.service.ts        # 主进程服务
UserCard.tsx            # 渲染进程组件
dateUtils.ts            # 工具文件
window.service.spec.ts   # 测试文件
```

### 类命名
- 使用 PascalCase
- 服务类名以 Service 结尾
- 窗口类名以 Window 结尾

**示例**：
```typescript
export class WindowService {}     # 服务
export class MainWindow {}       # 窗口
export class AppHandler {}       # 处理器
```

### 接口命名
- 使用 PascalCase
- 接口名以 I 开头（可选）

**示例**：
```typescript
export interface WindowConfig {}  # 接口
export interface IWindowConfig {} # 带前缀的接口
```

### 变量命名
- 使用 camelCase
- 常量使用 UPPER_SNAKE_CASE

**示例**：
```typescript
const mainWindow = null;        // 变量
const APP_VERSION = '1.0.0';    // 常量
```

### 方法命名
- 使用 camelCase
- 动词开头
- 清晰表达功能

**示例**：
```typescript
createWindow(): BrowserWindow {   # 方法
  return new BrowserWindow();
}

handleQuit(): void {              # 事件处理
  app.quit();
}
```

## 最佳实践

### 1. 安全性
- 禁用 nodeIntegration
- 启用 contextIsolation
- 使用预加载脚本
- 验证 IPC 通信

**示例**：
```typescript
// 安全的窗口配置
const mainWindow = new BrowserWindow({
  webPreferences: {
    preload: path.join(__dirname, '../preload/index.js'),
    nodeIntegration: false,
    contextIsolation: true,
    sandbox: true,
    webSecurity: true
  }
});
```

### 2. IPC 通信
- 使用 contextBridge 暴露 API
- 验证 IPC 通信
- 处理错误
- 使用类型定义

**示例**：
```typescript
// 主进程
ipcMain.handle('file:read', async (event, filePath) => {
  try {
    const content = await fs.promises.readFile(filePath, 'utf-8');
    return { success: true, content };
  } catch (error) {
    return { success: false, error: error.message };
  }
});

// 渲染进程
const result = await window.electronAPI.file.readFile('/path/to/file');
if (result.success) {
  console.log(result.content);
} else {
  console.error(result.error);
}
```

### 3. 窗口管理
- 单例窗口管理
- 窗口状态保存
- 窗口恢复
- 多窗口支持

**示例**：
```typescript
class WindowManager {
  private windows: Map<string, BrowserWindow> = new Map();

  createWindow(id: string, options: BrowserWindowConstructorOptions): BrowserWindow {
    if (this.windows.has(id)) {
      return this.windows.get(id)!;
    }

    const window = new BrowserWindow(options);
    this.windows.set(id, window);

    window.on('closed', () => {
      this.windows.delete(id);
    });

    return window;
  }

  getWindow(id: string): BrowserWindow | undefined {
    return this.windows.get(id);
  }

  closeWindow(id: string): void {
    const window = this.getWindow(id);
    if (window) {
      window.close();
    }
  }
}
```

### 4. 自动更新
- 使用 electron-updater
- 检查更新
- 下载更新
- 安装更新

**示例**：
```typescript
import { autoUpdater } from 'electron-updater';

export function setupAutoUpdater() {
  autoUpdater.checkForUpdatesAndNotify();

  autoUpdater.on('update-available', () => {
    logger.info('Update available');
  });

  autoUpdater.on('update-downloaded', () => {
    logger.info('Update downloaded');
    autoUpdater.quitAndInstall();
  });
}
```

### 5. 打包发布
- 使用 electron-builder
- 配置签名
- 配置更新服务器
- 多平台支持

**示例**：
```json
{
  "scripts": {
    "build": "electron-builder",
    "build:win": "electron-builder --win",
    "build:mac": "electron-builder --mac",
    "build:linux": "electron-builder --linux",
    "release": "electron-builder --publish always"
  }
}
```

### 6. 性能优化
- 懒加载模块
- 优化渲染性能
- 减少内存占用
- 优化启动时间

**示例**：
```typescript
// 懒加载模块
app.whenReady().then(async () => {
  const { createWindow } = await import('./window');
  const { createMenu } = await import('./menu');
  
  createWindow();
  createMenu();
});
```

## 注意事项

1. **安全性**：始终禁用 nodeIntegration，启用 contextIsolation
2. **跨平台**：考虑不同平台的差异
3. **性能**：优化渲染性能和内存占用
4. **测试**：进行充分的测试
5. **文档**：为公共 API 提供文档
6. **更新**：实现自动更新功能
7. **依赖管理**：合理管理第三方依赖
8. **打包**：正确配置打包和发布流程
