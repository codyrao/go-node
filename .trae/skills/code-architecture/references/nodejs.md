# Node.js 代码架构规范

## 概述
本文档定义了 Node.js 项目的代码架构规范，遵循 Node.js 官方的最佳实践和社区约定。适用于 Express、Koa、NestJS 等 Node.js 框架。

## 项目结构

### 标准项目结构
```
project-name/
├── src/                   # 源代码目录
│   ├── app/               # 应用程序代码
│   │   ├── controllers/   # 控制器
│   │   ├── services/      # 服务层
│   │   ├── repositories/  # 数据访问层
│   │   ├── models/        # 数据模型
│   │   ├── dto/           # 数据传输对象
│   │   ├── middleware/    # 中间件
│   │   ├── guards/        # 守卫
│   │   ├── interceptors/  # 拦截器
│   │   ├── filters/       # 过滤器
│   │   ├── pipes/         # 管道
│   │   ├── decorators/    # 装饰器
│   │   ├── utils/         # 工具函数
│   │   ├── constants/     # 常量定义
│   │   ├── config/        # 配置文件
│   │   ├── modules/       # 模块
│   │   ├── common/        # 公共模块
│   │   └── main.ts       # 应用入口
│   ├── assets/            # 静态资源
│   │   ├── images/
│   │   ├── fonts/
│   │   └── styles/
│   ├── database/          # 数据库相关
│   │   ├── migrations/   # 数据库迁移
│   │   ├── seeds/        # 数据填充
│   │   └── schema/       # 数据库模式
│   ├── jobs/              # 定时任务
│   ├── queues/            # 消息队列
│   ├── events/            # 事件处理
│   ├── interfaces/        # 接口定义
│   ├── types/             # 类型定义
│   └── index.ts          # 入口文件
├── test/                  # 测试目录
│   ├── unit/             # 单元测试
│   ├── integration/      # 集成测试
│   └── e2e/              # 端到端测试
├── scripts/               # 构建脚本
│   ├── build.js          # 构建脚本
│   ├── dev.js            # 开发脚本
│   └── seed.js           # 数据填充脚本
├── docs/                  # 文档目录
├── dist/                  # 构建产物（自动生成）
├── logs/                  # 日志目录
├── uploads/               # 上传文件目录
├── node_modules/          # 依赖包（自动生成）
├── package.json           # 项目配置
├── package-lock.json      # 依赖锁定
├── tsconfig.json          # TypeScript 配置
├── nest-cli.json          # NestJS CLI 配置
├── .eslintrc.js           # ESLint 配置
├── .prettierrc.js         # Prettier 配置
├── .gitignore             # Git 忽略文件
├── .env                   # 环境变量
├── .env.development       # 开发环境变量
├── .env.production        # 生产环境变量
├── .env.test             # 测试环境变量
└── README.md              # 项目说明文档
```

## 目录说明

### src/app/
**功能**：应用程序代码

**说明**：
- controllers/：控制器层，处理 HTTP 请求
- services/：服务层，业务逻辑
- repositories/：数据访问层，数据库操作
- models/：数据模型
- dto/：数据传输对象
- middleware/：中间件
- guards/：守卫，权限控制
- interceptors/：拦截器
- filters/：过滤器，异常处理
- pipes/：管道，数据验证
- decorators/：装饰器
- utils/：工具函数
- constants/：常量定义
- config/：配置文件
- modules/：模块
- common/：公共模块

**示例**：
```
src/app/
├── controllers/
│   ├── user.controller.ts
│   ├── product.controller.ts
│   └── order.controller.ts
├── services/
│   ├── user.service.ts
│   ├── product.service.ts
│   └── order.service.ts
├── repositories/
│   ├── user.repository.ts
│   ├── product.repository.ts
│   └── order.repository.ts
├── models/
│   ├── user.model.ts
│   ├── product.model.ts
│   └── order.model.ts
├── dto/
│   ├── user.dto.ts
│   ├── product.dto.ts
│   └── order.dto.ts
├── middleware/
│   ├── logger.middleware.ts
│   ├── auth.middleware.ts
│   └── validation.middleware.ts
├── guards/
│   ├── auth.guard.ts
│   └── roles.guard.ts
├── interceptors/
│   ├── logging.interceptor.ts
│   └── transform.interceptor.ts
├── filters/
│   ├── http-exception.filter.ts
│   └── all-exceptions.filter.ts
├── pipes/
│   ├── validation.pipe.ts
│   └── parse-int.pipe.ts
├── decorators/
│   ├── roles.decorator.ts
│   └── public.decorator.ts
├── utils/
│   ├── date.util.ts
│   ├── string.util.ts
│   └── crypto.util.ts
├── constants/
│   ├── error.constants.ts
│   └── status.constants.ts
├── config/
│   ├── app.config.ts
│   ├── database.config.ts
│   └── redis.config.ts
├── modules/
│   ├── user.module.ts
│   ├── product.module.ts
│   └── order.module.ts
├── common/
│   ├── common.module.ts
│   └── common.constants.ts
└── main.ts
```

**控制器示例**：
```typescript
// src/app/controllers/user.controller.ts
import { Controller, Get, Post, Put, Delete, Body, Param, UseGuards } from '@nestjs/common';
import { UserService } from '../services/user.service';
import { CreateUserDto, UpdateUserDto } from '../dto/user.dto';
import { AuthGuard } from '../guards/auth.guard';

@Controller('users')
@UseGuards(AuthGuard)
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Get()
  findAll() {
    return this.userService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.userService.findOne(+id);
  }

  @Post()
  create(@Body() createUserDto: CreateUserDto) {
    return this.userService.create(createUserDto);
  }

  @Put(':id')
  update(@Param('id') id: string, @Body() updateUserDto: UpdateUserDto) {
    return this.userService.update(+id, updateUserDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.userService.remove(+id);
  }
}
```

**服务示例**：
```typescript
// src/app/services/user.service.ts
import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from '../models/user.model';
import { CreateUserDto, UpdateUserDto } from '../dto/user.dto';

@Injectable()
export class UserService {
  constructor(
    @InjectRepository(User)
    private readonly userRepository: Repository<User>,
  ) {}

  async findAll(): Promise<User[]> {
    return this.userRepository.find();
  }

  async findOne(id: number): Promise<User> {
    const user = await this.userRepository.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException(`User with ID ${id} not found`);
    }
    return user;
  }

  async create(createUserDto: CreateUserDto): Promise<User> {
    const user = this.userRepository.create(createUserDto);
    return this.userRepository.save(user);
  }

  async update(id: number, updateUserDto: UpdateUserDto): Promise<User> {
    const user = await this.findOne(id);
    this.userRepository.merge(user, updateUserDto);
    return this.userRepository.save(user);
  }

  async remove(id: number): Promise<void> {
    const user = await this.findOne(id);
    await this.userRepository.remove(user);
  }
}
```

**仓库示例**：
```typescript
// src/app/repositories/user.repository.ts
import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from '../models/user.model';

@Injectable()
export class UserRepository {
  constructor(
    @InjectRepository(User)
    private readonly repository: Repository<User>,
  ) {}

  async findById(id: number): Promise<User | null> {
    return this.repository.findOne({ where: { id } });
  }

  async findByEmail(email: string): Promise<User | null> {
    return this.repository.findOne({ where: { email } });
  }

  async create(userData: Partial<User>): Promise<User> {
    const user = this.repository.create(userData);
    return this.repository.save(user);
  }

  async update(id: number, userData: Partial<User>): Promise<User> {
    await this.repository.update(id, userData);
    return this.findById(id);
  }

  async delete(id: number): Promise<void> {
    await this.repository.delete(id);
  }
}
```

**DTO 示例**：
```typescript
// src/app/dto/user.dto.ts
import { IsEmail, IsNotEmpty, IsOptional, IsString, Min, Max } from 'class-validator';

export class CreateUserDto {
  @IsEmail()
  @IsNotEmpty()
  email: string;

  @IsString()
  @IsNotEmpty()
  name: string;

  @IsOptional()
  @Min(0)
  @Max(120)
  age?: number;
}

export class UpdateUserDto {
  @IsOptional()
  @IsEmail()
  email?: string;

  @IsOptional()
  @IsString()
  name?: string;

  @IsOptional()
  @Min(0)
  @Max(120)
  age?: number;
}
```

**中间件示例**：
```typescript
// src/app/middleware/logger.middleware.ts
import { Injectable, NestMiddleware, Logger } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';

@Injectable()
export class LoggerMiddleware implements NestMiddleware {
  private readonly logger = new Logger('HTTP');

  use(req: Request, res: Response, next: NextFunction) {
    const { method, url, ip } = req;
    const userAgent = req.get('user-agent') || '';

    this.logger.log(`${method} ${url} - ${ip} - ${userAgent}`);

    res.on('finish', () => {
      const { statusCode } = res;
      this.logger.log(`${method} ${url} - ${statusCode}`);
    });

    next();
  }
}
```

**守卫示例**：
```typescript
// src/app/guards/auth.guard.ts
import { Injectable, CanActivate, ExecutionContext, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private readonly jwtService: JwtService) {}

  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest();
    const token = this.extractTokenFromHeader(request);

    if (!token) {
      throw new UnauthorizedException();
    }

    try {
      const payload = this.jwtService.verify(token);
      request.user = payload;
    } catch {
      throw new UnauthorizedException();
    }

    return true;
  }

  private extractTokenFromHeader(request: Request): string | undefined {
    const [type, token] = request.headers.authorization?.split(' ') ?? [];
    return type === 'Bearer' ? token : undefined;
  }
}
```

**过滤器示例**：
```typescript
// src/app/filters/http-exception.filter.ts
import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import { Request, Response } from 'express';

@Catch(HttpException)
export class HttpExceptionFilter implements ExceptionFilter {
  catch(exception: HttpException, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();
    const status = exception.getStatus();

    response.status(status).json({
      statusCode: status,
      timestamp: new Date().toISOString(),
      path: request.url,
      message: exception.message,
    });
  }
}
```

**管道示例**：
```typescript
// src/app/pipes/validation.pipe.ts
import {
  PipeTransform,
  Injectable,
  ArgumentMetadata,
  BadRequestException,
} from '@nestjs/common';
import { validate } from 'class-validator';
import { plainToClass } from 'class-transformer';

@Injectable()
export class ValidationPipe implements PipeTransform<any> {
  async transform(value: any, { metatype }: ArgumentMetadata) {
    if (!metatype || !this.toValidate(metatype)) {
      return value;
    }

    const object = plainToClass(metatype, value);
    const errors = await validate(object);

    if (errors.length > 0) {
      throw new BadRequestException(this.formatErrors(errors));
    }

    return value;
  }

  private toValidate(metatype: Function): boolean {
    const types: Function[] = [String, Boolean, Number, Array, Object];
    return !types.includes(metatype);
  }

  private formatErrors(errors: any[]): string {
    return errors
      .map((error) => Object.values(error.constraints).join(', '))
      .join(', ');
  }
}
```

### src/database/
**功能**：数据库相关

**说明**：
- migrations/：数据库迁移
- seeds/：数据填充
- schema/：数据库模式

**示例**：
```
src/database/
├── migrations/
│   ├── 20240101000001-create-users-table.ts
│   ├── 20240101000002-create-products-table.ts
│   └── 20240101000003-create-orders-table.ts
├── seeds/
│   ├── user.seed.ts
│   ├── product.seed.ts
│   └── order.seed.ts
└── schema/
    ├── user.schema.ts
    ├── product.schema.ts
    └── order.schema.ts
```

### src/jobs/
**功能**：定时任务

**说明**：
- 定时任务
- 后台任务
- 调度任务

**示例**：
```
src/jobs/
├── cleanup.job.ts
├── backup.job.ts
└── notification.job.ts
```

### src/queues/
**功能**：消息队列

**说明**：
- 消息队列
- 任务队列
- 事件队列

**示例**：
```
src/queues/
├── email.queue.ts
├── notification.queue.ts
└── report.queue.ts
```

### src/events/
**功能**：事件处理

**说明**：
- 事件监听
- 事件发布
- 事件订阅

**示例**：
```
src/events/
├── user-created.event.ts
├── user-updated.event.ts
└── user-deleted.event.ts
```

## 命名规范

### 文件命名
- 使用 kebab-case
- 测试文件以 `.spec.ts` 或 `.test.ts` 结尾

**示例**：
```
user.controller.ts       # 控制器
user.service.ts          # 服务
user.repository.ts       # 仓库
user.model.ts            # 模型
user.dto.ts              # DTO
user.controller.spec.ts  # 测试文件
```

### 类命名
- 使用 PascalCase
- 控制器类名以 Controller 结尾
- 服务类名以 Service 结尾
- 仓库类名以 Repository 结尾
- 模型类名以 Model 或 Entity 结尾
- DTO 类名以 Dto 结尾

**示例**：
```typescript
export class UserController {}       # 控制器
export class UserService {}          # 服务
export class UserRepository {}       # 仓库
export class User {}                # 模型
export class CreateUserDto {}       # DTO
```

### 接口命名
- 使用 PascalCase
- 接口名以 I 开头（可选）

**示例**：
```typescript
export interface IUser {}            # 接口
export interface User {}            # 不带前缀的接口
```

### 变量命名
- 使用 camelCase
- 常量使用 UPPER_SNAKE_CASE

**示例**：
```typescript
const userName = 'John';           # 变量
const MAX_RETRY_COUNT = 3;         # 常量
```

### 方法命名
- 使用 camelCase
- 动词开头
- 清晰表达功能

**示例**：
```typescript
findAll(): Promise<User[]> {        # 方法
  return this.userRepository.find();
}

async createUser(dto: CreateUserDto): Promise<User> {  # 异步方法
  return this.userRepository.create(dto);
}
```

## 最佳实践

### 1. 依赖注入
- 使用依赖注入
- 构造函数注入
- 避免循环依赖

**示例**：
```typescript
@Injectable()
export class UserService {
  constructor(
    private readonly userRepository: UserRepository,
    private readonly emailService: EmailService,
  ) {}
}
```

### 2. 错误处理
- 使用异常过滤器
- 自定义异常
- 统一错误响应

**示例**：
```typescript
// 自定义异常
export class UserNotFoundException extends NotFoundException {
  constructor(id: number) {
    super(`User with ID ${id} not found`);
  }
}

// 使用自定义异常
throw new UserNotFoundException(userId);
```

### 3. 数据验证
- 使用 class-validator
- 使用验证管道
- 自定义验证器

**示例**：
```typescript
export class CreateUserDto {
  @IsEmail()
  @IsNotEmpty()
  email: string;

  @IsString()
  @IsNotEmpty()
  @MinLength(2)
  name: string;
}
```

### 4. 日志记录
- 使用日志服务
- 记录关键操作
- 结构化日志

**示例**：
```typescript
import { Logger } from '@nestjs/common';

@Injectable()
export class UserService {
  private readonly logger = new Logger(UserService.name);

  async createUser(dto: CreateUserDto) {
    this.logger.log(`Creating user with email: ${dto.email}`);
    const user = await this.userRepository.create(dto);
    this.logger.log(`User created with ID: ${user.id}`);
    return user;
  }
}
```

### 5. 配置管理
- 使用配置模块
- 环境变量
- 配置验证

**示例**：
```typescript
// config/app.config.ts
export default () => ({
  port: parseInt(process.env.PORT, 10) || 3000,
  nodeEnv: process.env.NODE_ENV || 'development',
  jwt: {
    secret: process.env.JWT_SECRET || 'secret',
    expiresIn: process.env.JWT_EXPIRES_IN || '7d',
  },
});
```

### 6. 数据库操作
- 使用事务
- 优化查询
- 避免N+1查询

**示例**：
```typescript
async transferFunds(fromId: number, toId: number, amount: number) {
  await this.connection.transaction(async (manager) => {
    await manager.decrement(User, fromId, 'balance', amount);
    await manager.increment(User, toId, 'balance', amount);
  });
}
```

### 7. 安全性
- 输入验证
- 输出编码
- 认证授权
- HTTPS

**示例**：
```typescript
// 使用 helmet 中间件
import helmet from 'helmet';

app.use(helmet());

// 使用 rate limiting
import rateLimit from 'express-rate-limit';

app.use(rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100 // limit each IP to 100 requests per windowMs
}));
```

## 注意事项

1. **模块化设计**：合理划分模块，避免模块过大
2. **依赖注入**：正确使用依赖注入
3. **测试覆盖**：保持高测试覆盖率
4. **文档完整**：为公共 API 提供文档
5. **安全考虑**：注意输入验证、输出编码、认证授权
6. **性能优化**：优化数据库查询、缓存、异步处理
7. **错误处理**：统一错误处理，提供友好的错误信息
8. **日志记录**：记录关键操作，便于调试和监控
