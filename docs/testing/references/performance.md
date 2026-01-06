# 性能测试指南

> 本文档介绍性能测试和基准测试的方法。

---

## 📋 目录

- [基准测试](#-基准测试)
- [性能要求](#-性能要求)
- [性能分析](#-性能分析)
- [性能优化建议](#-性能优化建议)

---

## 基准测试

Go 标准库提供了内置的基准测试功能，使用 `go test -bench` 运行。

### 基本使用

```go
import (
    "testing"
)

func BenchmarkParseTicketID(b *testing.B) {
    input := "PROJ-123"
    for i := 0; i < b.N; i++ {
        ParseTicketID(input)
    }
}
```

### 运行基准测试

```bash
# 运行所有基准测试
go test -bench=. ./...

# 运行特定包的基准测试
go test -bench=. ./internal/lib/config

# 运行匹配模式的基准测试
go test -bench=BenchmarkParse ./internal/lib/config

# 显示内存分配
go test -bench=. -benchmem ./...

# 设置基准测试时间
go test -bench=. -benchtime=5s ./...
```

### 基准测试示例

```go
func BenchmarkParseTicketID(b *testing.B) {
    input := "PROJ-123"
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ParseTicketID(input)
    }
}

func BenchmarkParseTicketID_Parallel(b *testing.B) {
    input := "PROJ-123"
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            ParseTicketID(input)
        }
    })
}
```

### 子基准测试

```go
func BenchmarkParseTicketID(b *testing.B) {
    tests := []struct {
        name  string
        input string
    }{
        {"short", "A-1"},
        {"medium", "PROJ-123"},
        {"long", "VERY-LONG-PROJECT-NAME-123"},
    }

    for _, tt := range tests {
        b.Run(tt.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                ParseTicketID(tt.input)
            }
        })
    }
}
```

---

## 性能要求

### 启动速度

- **目标**：< 200ms
- **测量方法**：使用基准测试测量启动时间

```go
func BenchmarkStartup(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 模拟启动过程
        cmd := exec.Command("workflow", "version")
        cmd.Run()
    }
}
```

### 命令执行速度

- **目标**：常用命令 < 100ms
- **测量方法**：使用基准测试测量命令执行时间

```go
func BenchmarkCommandExecution(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 执行命令
        ExecuteCommand("workflow", "config", "show")
    }
}
```

### 内存占用

- **目标**：运行时 < 50MB
- **测量方法**：使用 `go test -benchmem` 测量内存分配

```bash
go test -bench=. -benchmem ./...
```

---

## 性能分析

### CPU Profile

```bash
# 生成 CPU profile
go test -bench=. -cpuprofile=cpu.prof ./internal/lib/config

# 分析 CPU profile
go tool pprof cpu.prof

# 在 pprof 交互界面中
(pprof) top
(pprof) list ParseTicketID
(pprof) web
```

### 内存 Profile

```bash
# 生成内存 profile
go test -bench=. -memprofile=mem.prof ./internal/lib/config

# 分析内存 profile
go tool pprof mem.prof

# 在 pprof 交互界面中
(pprof) top
(pprof) list ParseTicketID
(pprof) web
```

### Trace 分析

```bash
# 生成 trace 文件
go test -bench=. -trace=trace.out ./internal/lib/config

# 分析 trace
go tool trace trace.out
```

---

## 性能优化建议

### 1. 减少内存分配

```go
// ❌ 不推荐：频繁分配内存
func ProcessData(data []byte) string {
    result := ""
    for _, b := range data {
        result += string(b) // 每次连接都分配新内存
    }
    return result
}

// ✅ 推荐：使用 strings.Builder
func ProcessData(data []byte) string {
    var builder strings.Builder
    for _, b := range data {
        builder.WriteByte(b)
    }
    return builder.String()
}
```

### 2. 使用对象池

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data []byte) {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    buf.Reset()

    // 使用 buf
    // ...
}
```

### 3. 避免不必要的复制

```go
// ❌ 不推荐：不必要的复制
func ProcessData(data []byte) {
    copy := make([]byte, len(data))
    copy(copy, data)
    // 使用 copy
}

// ✅ 推荐：直接使用原始数据
func ProcessData(data []byte) {
    // 直接使用 data
}
```

### 4. 使用并发

```go
// ✅ 推荐：使用并发处理
func ProcessDataParallel(data []byte) {
    var wg sync.WaitGroup
    chunkSize := len(data) / runtime.NumCPU()

    for i := 0; i < runtime.NumCPU(); i++ {
        wg.Add(1)
        go func(start, end int) {
            defer wg.Done()
            // 处理数据块
        }(i*chunkSize, (i+1)*chunkSize)
    }

    wg.Wait()
}
```

### 5. 延迟初始化

```go
// ✅ 推荐：延迟初始化
var configOnce sync.Once
var config *Config

func GetConfig() *Config {
    configOnce.Do(func() {
        config = loadConfig()
    })
    return config
}
```

---

## 性能测试检查清单

### 开发时

- [ ] 为新功能添加基准测试
- [ ] 运行基准测试，确保性能不下降
- [ ] 使用 `-benchmem` 检查内存分配

### 代码审查时

- [ ] 检查新代码的性能影响
- [ ] 确保关键路径的性能符合要求
- [ ] 检查是否有不必要的内存分配

### 发布前

- [ ] 运行完整的基准测试套件
- [ ] 检查性能指标是否符合要求
- [ ] 生成性能报告

---

## 相关文档

- [测试组织规范](../organization.md) - 测试组织结构
- [测试编写规范](../writing.md) - 测试编写规范
- [测试命令参考](../commands.md) - 常用测试命令

---

**最后更新**: 2025-01-28
