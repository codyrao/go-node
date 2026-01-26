# Angular 代码架构规范

## 概述
本文档定义了 Angular 项目的代码架构规范，遵循 Angular 官方的最佳实践和风格指南。适用于 Angular 12+、Angular 14+、Angular 16+ 项目。

## 项目结构

### 标准项目结构
```
project-name/
├── src/                   # 源代码目录
│   ├── app/               # 应用程序代码
│   │   ├── core/          # 核心模块
│   │   │   ├── guards/    # 路由守卫
│   │   │   ├── interceptors/ # HTTP 拦截器
│   │   │   ├── services/  # 核心服务
│   │   │   └── models/    # 核心模型
│   │   ├── shared/        # 共享模块
│   │   │   ├── components/ # 共享组件
│   │   │   ├── directives/ # 共享指令
│   │   │   ├── pipes/     # 共享管道
│   │   │   ├── services/  # 共享服务
│   │   │   └── modules/   # 共享模块
│   │   ├── features/      # 功能模块
│   │   │   ├── user/      # 用户功能
│   │   │   │   ├── components/
│   │   │   │   ├── services/
│   │   │   │   ├── models/
│   │   │   │   ├── pages/
│   │   │   │   └── user.module.ts
│   │   │   ├── product/   # 产品功能
│   │   │   └── order/     # 订单功能
│   │   ├── layouts/       # 布局组件
│   │   ├── pages/         # 页面组件
│   │   ├── assets/        # 静态资源
│   │   │   ├── images/
│   │   │   ├── fonts/
│   │   │   ├── styles/
│   │   │   └── icons/
│   │   ├── environments/  # 环境配置
│   │   │   ├── environment.ts
│   │   │   ├── environment.prod.ts
│   │   │   └── environment.dev.ts
│   │   ├── app.component.ts
│   │   ├── app.component.html
│   │   ├── app.component.css
│   │   ├── app.module.ts
│   │   └── app-routing.module.ts
│   ├── styles/            # 全局样式
│   │   ├── styles.scss
│   │   └── variables.scss
│   ├── index.html         # HTML 模板
│   ├── main.ts            # 入口文件
│   └── polyfills.ts       # Polyfills
├── e2e/                   # 端到端测试
├── docs/                  # 文档目录
├── dist/                  # 构建产物（自动生成）
├── node_modules/          # 依赖包（自动生成）
├── angular.json           # Angular CLI 配置
├── tsconfig.json          # TypeScript 配置
├── tsconfig.app.json      # 应用 TypeScript 配置
├── tsconfig.spec.json     # 测试 TypeScript 配置
├── package.json           # 项目配置
├── package-lock.json      # 依赖锁定
├── .eslintrc.js           # ESLint 配置
├── .prettierrc.js         # Prettier 配置
├── .gitignore             # Git 忽略文件
└── README.md              # 项目说明文档
```

## 目录说明

### src/app/
**功能**：应用程序代码

**说明**：
- core/：核心模块，单例服务、全局配置
- shared/：共享模块，可复用的组件和功能
- features/：功能模块，按业务功能划分
- layouts/：布局组件
- pages/：页面组件

### src/app/core/
**功能**：核心模块

**说明**：
- 路由守卫
- HTTP 拦截器
- 核心服务
- 核心模型

**示例**：
```
src/app/core/
├── guards/
│   ├── auth.guard.ts      # 认证守卫
│   └── role.guard.ts      # 角色守卫
├── interceptors/
│   ├── auth.interceptor.ts # 认证拦截器
│   └── error.interceptor.ts # 错误拦截器
├── services/
│   ├── api.service.ts    # API 服务
│   └── storage.service.ts # 存储服务
├── models/
│   └── response.model.ts # 响应模型
└── core.module.ts        # 核心模块
```

**守卫示例**：
```typescript
// src/app/core/guards/auth.guard.ts
import { Injectable } from '@angular/core';
import { CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { AuthService } from '../services/auth.service';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate {
  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): boolean {
    if (this.authService.isLoggedIn()) {
      return true;
    }
    this.router.navigate(['/login'], { queryParams: { returnUrl: state.url } });
    return false;
  }
}
```

**拦截器示例**：
```typescript
// src/app/core/interceptors/auth.interceptor.ts
import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from '../services/auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private authService: AuthService) {}

  intercept(
    request: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    const token = this.authService.getToken();
    if (token) {
      request = request.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`
        }
      });
    }
    return next.handle(request);
  }
}
```

### src/app/shared/
**功能**：共享模块

**说明**：
- 共享组件
- 共享指令
- 共享管道
- 共享服务

**示例**：
```
src/app/shared/
├── components/
│   ├── button/
│   │   ├── button.component.ts
│   │   ├── button.component.html
│   │   ├── button.component.scss
│   │   └── button.module.ts
│   ├── input/
│   │   ├── input.component.ts
│   │   ├── input.component.html
│   │   ├── input.component.scss
│   │   └── input.module.ts
│   └── modal/
│       ├── modal.component.ts
│       ├── modal.component.html
│       ├── modal.component.scss
│       └── modal.module.ts
├── directives/
│   ├── click-outside/
│   │   ├── click-outside.directive.ts
│   │   └── click-outside.module.ts
│   └── auto-focus/
│       ├── auto-focus.directive.ts
│       └── auto-focus.module.ts
├── pipes/
│   ├── safe-html/
│   │   ├── safe-html.pipe.ts
│   │   └── safe-html.module.ts
│   └── date-format/
│       ├── date-format.pipe.ts
│       └── date-format.module.ts
├── services/
│   ├── notification.service.ts
│   └── confirmation.service.ts
└── shared.module.ts
```

**组件示例**：
```typescript
// src/app/shared/components/button/button.component.ts
import { Component, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-button',
  templateUrl: './button.component.html',
  styleUrls: ['./button.component.scss']
})
export class ButtonComponent {
  @Input() type: 'primary' | 'secondary' | 'danger' = 'primary';
  @Input() disabled = false;
  @Input() loading = false;
  @Output() click = new EventEmitter<void>();

  handleClick(): void {
    if (!this.disabled && !this.loading) {
      this.click.emit();
    }
  }
}
```

```html
<!-- src/app/shared/components/button/button.component.html -->
<button
  [class]="'btn btn-' + type"
  [disabled]="disabled || loading"
  (click)="handleClick()"
>
  <ng-content></ng-content>
</button>
```

### src/app/features/
**功能**：功能模块

**说明**：
- 按业务功能划分
- 每个功能模块独立
- 包含组件、服务、模型等

**示例**：
```
src/app/features/
├── user/
│   ├── components/
│   │   ├── user-card/
│   │   │   ├── user-card.component.ts
│   │   │   ├── user-card.component.html
│   │   │   └── user-card.component.scss
│   │   └── user-list/
│   │       ├── user-list.component.ts
│   │       ├── user-list.component.html
│   │       └── user-list.component.scss
│   ├── services/
│   │   ├── user.service.ts
│   │   └── user.service.spec.ts
│   ├── models/
│   │   ├── user.model.ts
│   │   └── user-create.model.ts
│   ├── pages/
│   │   ├── user-list-page/
│   │   │   ├── user-list-page.component.ts
│   │   │   ├── user-list-page.component.html
│   │   │   └── user-list-page.component.scss
│   │   └── user-detail-page/
│   │       ├── user-detail-page.component.ts
│   │       ├── user-detail-page.component.html
│   │       └── user-detail-page.component.scss
│   ├── user-routing.module.ts
│   └── user.module.ts
├── product/
│   └── ...
└── order/
    └── ...
```

**功能模块示例**：
```typescript
// src/app/features/user/user.module.ts
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

import { UserRoutingModule } from './user-routing.module';
import { UserCardComponent } from './components/user-card/user-card.component';
import { UserListComponent } from './components/user-list/user-list.component';
import { UserListPageComponent } from './pages/user-list-page/user-list-page.component';
import { UserDetailPageComponent } from './pages/user-detail-page/user-detail-page.component';

@NgModule({
  declarations: [
    UserCardComponent,
    UserListComponent,
    UserListPageComponent,
    UserDetailPageComponent
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    UserRoutingModule
  ]
})
export class UserModule {}
```

**服务示例**：
```typescript
// src/app/features/user/services/user.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User, CreateUserRequest, UpdateUserRequest } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = 'api/users';

  constructor(private http: HttpClient) {}

  getUsers(): Observable<User[]> {
    return this.http.get<User[]>(this.apiUrl);
  }

  getUserById(id: number): Observable<User> {
    return this.http.get<User>(`${this.apiUrl}/${id}`);
  }

  createUser(data: CreateUserRequest): Observable<User> {
    return this.http.post<User>(this.apiUrl, data);
  }

  updateUser(id: number, data: UpdateUserRequest): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}/${id}`, data);
  }

  deleteUser(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }
}
```

### src/app/layouts/
**功能**：布局组件

**说明**：
- 主布局
- 认证布局
- 空布局

**示例**：
```
src/app/layouts/
├── main-layout/
│   ├── main-layout.component.ts
│   ├── main-layout.component.html
│   └── main-layout.component.scss
├── auth-layout/
│   ├── auth-layout.component.ts
│   ├── auth-layout.component.html
│   └── auth-layout.component.scss
└── empty-layout/
    ├── empty-layout.component.ts
    ├── empty-layout.component.html
    └── empty-layout.component.scss
```

### src/app/pages/
**功能**：页面组件

**说明**：
- 路由对应的页面组件
- 页面级别的业务逻辑

**示例**：
```
src/app/pages/
├── home/
│   ├── home.component.ts
│   ├── home.component.html
│   └── home.component.scss
├── login/
│   ├── login.component.ts
│   ├── login.component.html
│   └── login.component.scss
└── not-found/
    ├── not-found.component.ts
    ├── not-found.component.html
    └── not-found.component.scss
```

### src/app/assets/
**功能**：静态资源

**说明**：
- images/：图片资源
- fonts/：字体资源
- styles/：样式文件
- icons/：图标资源

**示例**：
```
src/app/assets/
├── images/
│   ├── logo.png
│   └── background.jpg
├── fonts/
│   ├── Roboto-Regular.ttf
│   └── Roboto-Bold.ttf
├── styles/
│   ├── main.scss
│   └── variables.scss
└── icons/
    ├── user.svg
    └── settings.svg
```

### src/styles/
**功能**：全局样式

**说明**：
- 全局样式文件
- 样式变量
- 样式重置

**示例**：
```scss
// src/styles/variables.scss
$primary-color: #1890ff;
$secondary-color: #f0f0f0;
$danger-color: #ff4d4f;
$success-color: #52c41a;
$warning-color: #faad14;

$font-family-base: 'Roboto', sans-serif;
$font-size-base: 14px;
$line-height-base: 1.5;

$border-radius-base: 4px;
$box-shadow-base: 0 2px 8px rgba(0, 0, 0, 0.15);
```

```scss
// src/styles/styles.scss
@import './variables.scss';

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: $font-family-base;
  font-size: $font-size-base;
  line-height: $line-height-base;
}

a {
  text-decoration: none;
  color: $primary-color;

  &:hover {
    text-decoration: underline;
  }
}
```

## 命名规范

### 文件命名
- 组件文件使用 kebab-case
- 服务文件使用 kebab-case
- 模型文件使用 kebab-case
- 测试文件以 `.spec.ts` 结尾

**示例**：
```
user-card.component.ts     # 组件文件
user.service.ts            # 服务文件
user.model.ts              # 模型文件
user.service.spec.ts       # 测试文件
```

### 类命名
- 使用 PascalCase
- 组件类名以 Component 结尾
- 服务类名以 Service 结尾
- 指令类名以 Directive 结尾
- 管道类名以 Pipe 结尾

**示例**：
```typescript
export class UserCardComponent {}       // 组件
export class UserService {}              // 服务
export class ClickOutsideDirective {}   // 指令
export class DateFormatPipe {}          // 管道
```

### 接口命名
- 使用 PascalCase
- 接口名以 I 开头（可选）
- 模型接口可以不加 I 前缀

**示例**：
```typescript
export interface User {}                 // 接口
export interface IUser {}               # 带前缀的接口
export interface CreateUserRequest {}   # 请求接口
export interface UpdateUserRequest {}   # 请求接口
```

### 选择器命名
- 使用 kebab-case
- 组件选择器以 app- 前缀
- 指令选择器以 app- 前缀

**示例**：
```typescript
@Component({
  selector: 'app-user-card',             # 组件选择器
  templateUrl: './user-card.component.html',
  styleUrls: ['./user-card.component.scss']
})
export class UserCardComponent {}

@Directive({
  selector: '[appClickOutside]'        # 指令选择器
})
export class ClickOutsideDirective {}
```

### 变量命名
- 使用 camelCase
- 常量使用 UPPER_SNAKE_CASE
- 私有成员以下划线开头

**示例**：
```typescript
const userName = 'John';               // 变量
const MAX_RETRY_COUNT = 3;             // 常量
private _userService: UserService;      // 私有成员
```

### 方法命名
- 使用 camelCase
- 动词开头
- 清晰表达功能

**示例**：
```typescript
getUserById(id: number): Observable<User> {    // 方法
  return this.http.get<User>(`${this.apiUrl}/${id}`);
}

handleClick(): void {                          // 事件处理
  this.click.emit();
}
```

## 最佳实践

### 1. 组件设计
- 单一职责原则
- 组件复用性
- 输入输出清晰
- 使用 OnPush 变更检测策略

**示例**：
```typescript
@Component({
  selector: 'app-user-card',
  templateUrl: './user-card.component.html',
  styleUrls: ['./user-card.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class UserCardComponent {
  @Input() user: User;
  @Output() delete = new EventEmitter<number>();

  handleDelete(): void {
    this.delete.emit(this.user.id);
  }
}
```

### 2. 服务设计
- 使用 Injectable 装饰器
- providedIn: 'root' 提供单例服务
- 使用 Observable 处理异步
- 使用 RxJS 操作符

**示例**：
```typescript
@Injectable({
  providedIn: 'root'
})
export class UserService {
  constructor(private http: HttpClient) {}

  getUsers(): Observable<User[]> {
    return this.http.get<User[]>(this.apiUrl).pipe(
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    console.error('An error occurred:', error);
    return throwError(() => error);
  }
}
```

### 3. 路由管理
- 懒加载模块
- 路由守卫
- 路由参数
- 路由数据

**示例**：
```typescript
const routes: Routes = [
  {
    path: '',
    component: MainLayoutComponent,
    children: [
      {
        path: '',
        redirectTo: '/home',
        pathMatch: 'full'
      },
      {
        path: 'home',
        component: HomeComponent
      },
      {
        path: 'users',
        loadChildren: () => import('./features/user/user.module').then(m => m.UserModule),
        canActivate: [AuthGuard]
      }
    ]
  },
  {
    path: 'login',
    component: AuthLayoutComponent,
    children: [
      {
        path: '',
        component: LoginComponent
      }
    ]
  }
];
```

### 4. 表单处理
- 使用响应式表单
- 表单验证
- 自定义验证器
- 表单状态管理

**示例**：
```typescript
@Component({
  selector: 'app-user-form',
  templateUrl: './user-form.component.html'
})
export class UserFormComponent implements OnInit {
  userForm: FormGroup;

  constructor(private fb: FormBuilder) {}

  ngOnInit(): void {
    this.userForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      age: ['', [Validators.required, Validators.min(18)]]
    });
  }

  onSubmit(): void {
    if (this.userForm.valid) {
      console.log(this.userForm.value);
    }
  }
}
```

### 5. 性能优化
- 使用 OnPush 变更检测
- 懒加载模块
- 使用 trackBy 优化列表
- 使用纯管道

**示例**：
```typescript
@Component({
  selector: 'app-user-list',
  template: `
    <div *ngFor="let user of users; trackBy: trackById">
      {{ user.name }}
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class UserListComponent {
  @Input() users: User[];

  trackById(index: number, user: User): number {
    return user.id;
  }
}
```

### 6. 代码风格
- 使用 ESLint 和 Prettier
- 统一代码风格
- 遵循 Angular 官方风格指南

## 注意事项

1. **模块化设计**：合理划分模块，避免模块过大
2. **依赖注入**：正确使用 Angular 依赖注入
3. **测试覆盖**：保持高测试覆盖率
4. **文档完整**：为公共组件和服务提供文档
5. **安全考虑**：注意 XSS 攻击，使用 DomSanitizer
6. **版本管理**：使用语义化版本号
7. **依赖管理**：合理管理第三方依赖
8. **响应式设计**：确保组件在不同设备上正常显示
