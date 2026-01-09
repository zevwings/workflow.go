# Fallback 错误处理策略

## 统一策略

为了保持 fallback 模式的一致性，所有 fallback 函数遵循以下错误处理策略：

### 1. 读取输入失败
**场景**: `terminal.ReadLine()` 返回错误（通常是 EOF 或空输入）

**处理**: 
- 返回默认值（不返回错误）
- 显示格式化的默认结果（如果有 ResultDisplay 函数）

**原因**: Fallback 模式是为了在无法使用交互式模式时的降级方案，空输入或读取失败应该被视为用户接受默认值，而不是错误。

### 2. 输入解析失败
**场景**: 用户输入无法解析为有效值（如非数字、超出范围等）

**处理**:
- 返回默认值（不返回错误）
- 显示格式化的默认结果（如果有 ResultDisplay 函数）

**原因**: 在 fallback 模式下，用户可能不熟悉输入格式，应该宽容处理，使用默认值而不是报错。

### 3. 用户取消（Ctrl+C）
**场景**: 用户主动取消操作

**处理**:
- 返回错误（`common.HandleCancel` 返回的错误）
- 不返回默认值

**原因**: 用户明确取消操作，应该传播取消错误，让调用方处理。

### 4. 其他系统错误
**场景**: 终端操作失败、格式化错误等

**处理**:
- 返回错误
- 不返回默认值

**原因**: 这些是真正的系统错误，应该传播给调用方处理。

## 实现示例

### Select Fallback
```go
func selectFallback(cfg SelectConfig) (int, error) {
    // ...
    inputLine, err := cfg.Terminal.ReadLine()
    if err != nil {
        // 读取失败 -> 返回默认值，不返回错误
        return defaultIndex, nil
    }
    
    selectedIndex, valid := parseInput(inputLine)
    if !valid {
        // 解析失败 -> 返回默认值，不返回错误
        return defaultIndex, nil
    }
    
    return selectedIndex, nil
}
```

### Confirm Fallback
```go
func confirmFallback(cfg ConfirmConfig) (bool, error) {
    // ...
    result, err := common.ExecuteFallbackTyped(...)
    if err != nil {
        // 如果是取消错误，传播错误
        // 否则返回默认值
        return cfg.DefaultYes, err
    }
    return result, nil
}
```

## 注意事项

1. **不要混淆错误和默认值**: 只有在真正的系统错误或用户取消时才返回错误
2. **保持一致性**: 所有 fallback 函数应该遵循相同的策略
3. **文档化**: 在函数注释中明确说明错误处理策略
