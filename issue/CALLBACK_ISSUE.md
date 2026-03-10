# Go-Node 回调混淆问题修复指南

## 问题描述

当前 SDK 的所有导出函数（`ScanGames`, `StartCloudBackup`, `EnableGameCloudSave` 等）都使用同一个全局回调函数 `gCallNodeCallback` 来返回结果。这导致当多个函数并发调用时，回调消息会混淆，无法区分是哪个函数的回调。

### 具体表现

在 `test.js` 中同时调用多个函数：
```javascript
game.scanGames({...}, callback);
game.startCloudBackup({...}, callback);
game.enableGameCloudSave({...}, callback);
```

当 `startCloudBackup` 返回 JWT Token 验证失败错误时：
```json
{
  "error": "download backup file failed: token validation failed: token illegal",
  "success": false
}
```

这个错误会被错误地归属到 `scanGames` 的回调输出中，因为所有函数共享同一个回调通道。

## 根本原因

1. **单一全局回调**: 所有函数使用同一个 `gCallNodeCallback` 变量
2. **无标识机制**: 返回的 JSON 数据中没有字段标识是哪个函数的回调
3. **并发冲突**: 多个函数同时执行时，回调消息无法区分归属

## 解决方案

### 方案 1: 在返回数据中添加 `callbackType` 字段（推荐）

修改 `main.go` 中所有导出函数，在返回的 JSON 数据中添加 `callbackType` 字段来标识回调来源。

#### 修改示例

**修改前:**
```go
//export ScanGames
func ScanGames(req *C.char, callbackType *C.char) {
    // ...
    resultData := map[string]interface{}{
        "success": false,
        "error":   err.Error(),
    }
    resultJson, _ := json.Marshal(resultData)
    C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(resultJson)))
}
```

**修改后:**
```go
//export ScanGames
func ScanGames(req *C.char, callbackType *C.char) {
    // ...
    resultData := map[string]interface{}{
        "callbackType": "scanGames",  // 添加标识字段
        "success":      false,
        "error":        err.Error(),
    }
    resultJson, _ := json.Marshal(resultData)
    C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(resultJson)))
}
```

### 需要修改的函数列表

| 函数名 | callbackType 值 |
|--------|-----------------|
| `ScanGames` | `"scanGames"` |
| `Init` | `"init"` |
| `EnableGameCloudSave` | `"enableGameCloudSave"` |
| `DisableGameCloudSave` | `"disableGameCloudSave"` |
| `ListCloudSaves` | `"listCloudSaves"` |
| `RenameCloudSaveName` | `"renameCloudSaveName"` |
| `DeleteOneCloudSave` | `"deleteOneCloudSave"` |
| `SyncCloudSave` | `"syncCloudSave"` |
| `StartCloudBackup` | `"startCloudBackup"` |

### 修改要点

1. **所有返回点都要添加**: 包括成功返回、错误返回、panic 恢复等所有调用 `C.callCallback` 的地方
2. **保持一致性**: 所有函数的 `callbackType` 值使用驼峰命名，与函数名对应
3. **不要修改 C 接口**: 保持 `callbackType *C.char` 参数不变（虽然当前未使用，但可用于未来扩展）

### 修改示例代码

```go
// ScanGames 中的修改示例
for game := range gamesCh {
    firstScanFinished := game.FirstScanFinished
    callbackData := map[string]interface{}{}
    if firstScanFinished {
        callbackData = map[string]interface{}{
            "callbackType": "scanGames",  // 添加标识
            "success":      true,
            "error":        "ok",
            "data": map[string]interface{}{
                "firstScanFinished": firstScanFinished,
            },
        }
    } else {
        callbackData = map[string]interface{}{
            "callbackType": "scanGames",  // 添加标识
            "success":      true,
            "error":        "ok",
            "data": map[string]interface{}{
                "firstScanFinished": firstScanFinished,
                "game":              game,
            },
        }
    }
    jsonData, _ := json.Marshal(callbackData)
    C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(jsonData)))
}
```

```go
// StartCloudBackup 中的修改示例
if err != nil {
    log.Printf("StartCloudBackup 失败: %v\n", err)
    resultData := map[string]interface{}{
        "callbackType": "startCloudBackup",  // 添加标识
        "success":      false,
        "error":        err.Error(),
    }
    resultJson, _ := json.Marshal(resultData)
    C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(resultJson)))
    return
}

resultData := map[string]interface{}{
    "callbackType": "startCloudBackup",  // 添加标识
    "success":      true,
    "error":        "ok",
}
resultJson, _ := json.Marshal(resultData)
C.callCallback(unsafe.Pointer(gCallNodeCallback), C.CString(string(resultJson)))
```

## Node.js 层使用方式

修改后，Node.js 层可以通过检查 `callbackType` 字段来区分回调：

```javascript
const game = require('./output/game.node');

// 统一回调处理函数
function handleCallback(res) {
    switch(res.callbackType) {
        case 'scanGames':
            console.log('扫描游戏回调:', res);
            break;
        case 'startCloudBackup':
            console.log('云备份回调:', res);
            break;
        case 'enableGameCloudSave':
            console.log('启用云存档回调:', res);
            break;
        // ... 其他 case
        default:
            console.log('未知回调:', res);
    }
}

// 调用多个函数
game.scanGames({...}, handleCallback);
game.startCloudBackup({...}, handleCallback);
game.enableGameCloudSave({...}, handleCallback);
```

## 文件位置

- **main.go**: `c:\Users\admin\Documents\goproject\xlasgame\game-asst-sdk\main.go`
- **所有导出函数**: 位于 `main.go` 中，以 `//export` 开头的函数

## 注意事项

1. 不要修改函数的 C 接口签名（保持 `callbackType *C.char` 参数）
2. 确保所有返回路径都添加了 `callbackType` 字段
3. `ScanGames` 会多次回调，每次都需要包含 `callbackType`
4. 保持向后兼容：如果 Node.js 层不检查 `callbackType`，代码仍然可以工作
