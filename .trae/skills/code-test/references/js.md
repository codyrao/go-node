# JavaScript 测试规范

## 概述
本文档定义了 JavaScript 的测试规范，包括单元测试、系统功能测试、系统压力测试（仅针对页面）和系统破坏性测试（仅针对页面）。

## 测试框架
使用 **Jest** 测试框架，它提供了简洁的测试语法和强大的功能。

### 安装 Jest
```bash
npm install --save-dev jest
```

## 单元测试

### 单元测试规范
- 单元测试文件名必须以 `.test.js` 或 `.spec.js` 结尾
- 使用 Jest 的测试语法
- 测试用例要覆盖正常功能、边界条件和异常输入

### 单元测试模板
```javascript
const { functionName } = require('./script');

describe('测试函数功能', () => {
    // 正常功能测试
    test('正常输入', () => {
        const result = functionName('正常参数');
        expect(result).toBe('预期值');
    });

    // 边界条件测试
    test('边界条件', () => {
        const result = functionName('边界值');
        expect(result).toBe('预期值');
    });

    // 异常输入测试（根据项目上下文判断）
    test('异常输入', () => {
        expect(() => {
            functionName('');
        }).toThrow();
    });
});
```

### 单元测试示例
```javascript
const { getUserById } = require('./user');

describe('测试根据 ID 获取用户', () => {
    // 正常功能测试
    test('正常用户 ID', () => {
        const user = getUserById(1);
        expect(user).toBeDefined();
        expect(user.id).toBe(1);
    });

    // 边界条件测试
    test('边界值 ID', () => {
        expect(() => {
            getUserById(0);
        }).toThrow();
    });

    // 异常输入测试（根据项目上下文，如果项目中 ID 不可能为负数，则不需要测试）
    // test('负数 ID', () => {
    //     expect(() => {
    //         getUserById(-1);
    //     }).toThrow();
    // });
});
```

## 系统功能测试

### 系统功能测试规范
- 系统功能测试放在 `./test/system/functional/` 目录下
- 使用 Jest 测试框架
- 测试要覆盖所有单元测试
- 对于网络 API、消息队列等网络中间件，要实现接口的功能测试

### 系统功能测试模板
```javascript
const axios = require('axios');

describe('测试 API 功能', () => {
    const baseURL = 'http://localhost:8080/api';

    // API 接口功能测试
    test('正常请求', async () => {
        const response = await axios.get(`${baseURL}/users`);
        expect(response.status).toBe(200);
        expect(response.data.code).toBe(0);
    });

    // 数据库操作测试
    test('数据库操作', async () => {
        // 测试数据库操作
    });

    // 消息队列功能测试
    test('消息队列功能', async () => {
        // 测试消息队列
    });
});
```

### 系统功能测试示例
```javascript
const axios = require('axios');

describe('测试用户 API 功能', () => {
    const baseURL = 'http://localhost:8080/api';

    // 获取用户列表
    test('获取用户列表', async () => {
        const response = await axios.get(`${baseURL}/users`);
        expect(response.status).toBe(200);
        expect(response.data.code).toBe(0);
        expect(response.data.data).toBeDefined();
    });

    // 创建用户
    test('创建用户', async () => {
        const userData = {
            name: 'test',
            age: 25
        };
        const response = await axios.post(`${baseURL}/users`, userData);
        expect(response.status).toBe(200);
        expect(response.data.code).toBe(0);
    });

    // 更新用户
    test('更新用户', async () => {
        const userData = {
            name: 'test2',
            age: 26
        };
        const response = await axios.put(`${baseURL}/users/1`, userData);
        expect(response.status).toBe(200);
        expect(response.data.code).toBe(0);
    });

    // 删除用户
    test('删除用户', async () => {
        const response = await axios.delete(`${baseURL}/users/1`);
        expect(response.status).toBe(200);
        expect(response.data.code).toBe(0);
    });
});
```

## 系统压力测试

### 系统压力测试规范
- 系统压力测试放在 `./test/system/stress/` 目录下
- 压力测试只针对页面
- 压力测试要支持多梯度，从而最大程度了解页面的性能
- 使用 Puppeteer 或 Playwright 进行页面压力测试

### 压力测试梯度
- **低并发**：1-10 并发
- **中并发**：10-100 并发
- **高并发**：100-1000 并发
- **极高并发**：1000+ 并发

### 系统压力测试模板
```javascript
const puppeteer = require('puppeteer');

describe('测试页面压力', () => {
    const baseURL = 'http://localhost:3000';

    // 低并发测试
    test('低并发测试 (1-10)', async () => {
        await testPageStress(`${baseURL}/users`, 5);
    });

    // 中并发测试
    test('中并发测试 (10-100)', async () => {
        await testPageStress(`${baseURL}/users`, 50);
    });

    // 高并发测试
    test('高并发测试 (100-1000)', async () => {
        await testPageStress(`${baseURL}/users`, 500);
    });

    // 极高并发测试
    test('极高并发测试 (1000+)', async () => {
        await testPageStress(`${baseURL}/users`, 2000);
    });
});

async function testPageStress(url, concurrency) {
    const startTime = Date.now();
    const promises = [];

    for (let i = 0; i < concurrency; i++) {
        promises.push(
            (async () => {
                const browser = await puppeteer.launch();
                const page = await browser.newPage();
                try {
                    await page.goto(url, { waitUntil: 'networkidle2' });
                    return { success: true };
                } catch (error) {
                    return { success: false, error: error.message };
                } finally {
                    await browser.close();
                }
            })()
        );
    }

    const results = await Promise.all(promises);
    const duration = Date.now() - startTime;
    const successCount = results.filter(r => r.success).length;
    const failCount = results.filter(r => !r.success).length;

    console.log(`URL: ${url}, 并发数: ${concurrency}, 成功: ${successCount}, 失败: ${failCount}, 耗时: ${duration}ms`);

    expect(failCount).toBeLessThan(concurrency / 10);
}
```

### 系统压力测试示例
```javascript
const puppeteer = require('puppeteer');

describe('测试用户页面压力', () => {
    const baseURL = 'http://localhost:3000';

    // 低并发测试
    test('低并发测试 (1-10)', async () => {
        await testPageLoad(`${baseURL}/users`, 5);
    });

    // 中并发测试
    test('中并发测试 (10-100)', async () => {
        await testPageLoad(`${baseURL}/users`, 50);
    });

    // 高并发测试
    test('高并发测试 (100-1000)', async () => {
        await testPageLoad(`${baseURL}/users`, 500);
    });

    // 极高并发测试
    test('极高并发测试 (1000+)', async () => {
        await testPageLoad(`${baseURL}/users`, 2000);
    });
});

async function testPageLoad(url, concurrency) {
    const startTime = Date.now();
    const promises = [];

    for (let i = 0; i < concurrency; i++) {
        promises.push(
            (async () => {
                const browser = await puppeteer.launch();
                const page = await browser.newPage();
                try {
                    await page.goto(url, { waitUntil: 'networkidle2' });
                    return { success: true };
                } catch (error) {
                    return { success: false, error: error.message };
                } finally {
                    await browser.close();
                }
            })()
        );
    }

    const results = await Promise.all(promises);
    const duration = Date.now() - startTime;
    const successCount = results.filter(r => r.success).length;
    const failCount = results.filter(r => !r.success).length;

    console.log(`URL: ${url}, 并发数: ${concurrency}, 成功: ${successCount}, 失败: ${failCount}, 耗时: ${duration}ms`);

    expect(failCount).toBeLessThan(concurrency / 10);
}
```

## 系统破坏性测试

### 系统破坏性测试规范
- 系统破坏性测试放在 `./test/system/destructive/` 目录下
- 破坏性测试只针对页面
- 传入空异常值和 SQL 注入等异常值
- 对异常情况，如果没有出现错误则认为测试失败
- 对于错误，不是用简单的协议状态码，而是根据项目功能场景判断
- 比如：返回 HTTP 错误码是 200，但是页面显示错误信息，也认为是返回了错误
- 使用 Puppeteer 或 Playwright 进行页面破坏性测试

### 系统破坏性测试模板
```javascript
const puppeteer = require('puppeteer');

describe('测试页面破坏性', () => {
    const baseURL = 'http://localhost:3000';

    // 空值测试
    test('空值测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', '');
            await page.type('#age', '');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });

    // 异常类型测试
    test('异常类型测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', '123');
            await page.type('#age', 'invalid');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });

    // SQL 注入测试
    test('SQL 注入测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', "admin' OR '1'='1");
            await page.type('#age', '25');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });

    // XSS 攻击测试
    test('XSS 攻击测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', "<script>alert('xss')</script>");
            await page.type('#age', '25');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });
});
```

### 系统破坏性测试示例
```javascript
const puppeteer = require('puppeteer');

describe('测试用户页面破坏性', () => {
    const baseURL = 'http://localhost:3000';

    // 空值测试
    test('空值测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users/create`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', '');
            await page.type('#age', '');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });

    // 异常类型测试
    test('异常类型测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users/create`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', '123');
            await page.type('#age', 'invalid');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });

    // SQL 注入测试
    test('SQL 注入测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users/create`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', "admin' OR '1'='1");
            await page.type('#age', '25');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });

    // XSS 攻击测试
    test('XSS 攻击测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users/create`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', "<script>alert('xss')</script>");
            await page.type('#age', '25');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });

    // 超长字符串测试
    test('超长字符串测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users/create`, { waitUntil: 'networkidle2' });
            
            const longString = 'a'.repeat(10000);
            await page.type('#name', longString);
            await page.type('#age', '25');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });

    // 特殊字符测试
    test('特殊字符测试', async () => {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        try {
            await page.goto(`${baseURL}/users/create`, { waitUntil: 'networkidle2' });
            
            await page.type('#name', "!@#$%^&*()_+-=[]{}|;:,.<>?");
            await page.type('#age', '25');
            await page.click('#submit');
            
            await page.waitForSelector('.error', { timeout: 3000 });
            const errorText = await page.$eval('.error', el => el.textContent);
            expect(errorText).toBeTruthy();
        } finally {
            await browser.close();
        }
    });
});
```

## 测试可视化界面

### 使用 Jest HTML Reporters
Jest HTML Reporters 是一个用于生成 HTML 测试报告的工具。

### 安装 Jest HTML Reporters
```bash
npm install --save-dev jest-html-reporters
```

### 配置 Jest
在 `package.json` 中添加配置：
```json
{
  "scripts": {
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage"
  },
  "jest": {
    "reporters": [
      "default",
      [
        "jest-html-reporters",
        {
          "publicPath": "./test/report",
          "filename": "report.html",
          "expand": true
        }
      ]
    ]
  }
}
```

### 访问测试报告
在浏览器中打开 `./test/report/report.html`

## 测试目录结构
```
./test/
├── README.md              # 测试使用说明
├── run.sh                 # 测试执行脚本
├── unit/                  # 单元测试
│   └── user.test.js
├── system/                # 系统测试
│   ├── functional/        # 功能测试
│   │   └── user_api.test.js
│   ├── stress/            # 压力测试
│   │   └── user_api_stress.test.js
│   └── destructive/       # 破坏性测试
│       └── user_api_destructive.test.js
└── report/                # 测试报告
    └── report.html
```

## 测试执行脚本

### run.sh 脚本内容
```bash
#!/bin/bash

# 测试执行脚本

case "$1" in
    all)
        echo "运行所有测试..."
        npm test
        ;;
    unit)
        echo "运行单元测试..."
        npm test ./unit/
        ;;
    functional)
        echo "运行系统功能测试..."
        npm test ./system/functional/
        ;;
    stress)
        echo "运行系统压力测试..."
        npm test ./system/stress/
        ;;
    destructive)
        echo "运行系统破坏性测试..."
        npm test ./system/destructive/
        ;;
    web)
        echo "启动测试可视化界面..."
        if [ -f ./test/report/report.html ]; then
            echo "测试报告已生成: ./test/report/report.html"
            echo "请在浏览器中打开该文件查看测试结果"
        else
            echo "测试报告不存在，请先运行测试"
        fi
        ;;
    *)
        echo "使用方法: $0 {all|unit|functional|stress|destructive|web}"
        exit 1
        ;;
esac
```

## 测试使用说明

### 测试环境要求
- Node.js 12 或更高版本
- Jest 测试框架
- axios（用于 API 测试）
- Puppeteer 或 Playwright（用于页面压力测试和破坏性测试）

### 测试依赖安装
```bash
npm install --save-dev jest
npm install --save-dev jest-html-reporters
npm install axios
npm install --save-dev puppeteer
# 或者
npm install --save-dev playwright
```

### 测试命令说明
```bash
# 运行所有测试
./test/run.sh all

# 运行单元测试
./test/run.sh unit

# 运行系统功能测试
./test/run.sh functional

# 运行系统压力测试
./test/run.sh stress

# 运行系统破坏性测试
./test/run.sh destructive

# 启动测试可视化界面
./test/run.sh web
```

### 测试结果查看
- 命令行查看测试结果
- HTML 报告查看测试结果和覆盖率

### 测试报告生成
```bash
# 生成测试覆盖率报告
npm run test:coverage
```

## 注意事项
- 单元测试要结合项目函数上下文，避免不必要的测试用例
- 系统压力测试只针对页面，使用 Puppeteer 或 Playwright 进行页面压力测试
- 系统破坏性测试只针对页面，根据项目功能场景判断错误，而不是简单的协议状态码
- 测试可视化界面使用 Jest HTML Reporters 生成 HTML 报告
- 测试命令要简单易用，支持选择性执行
- 保持测试代码的简洁和可维护性
