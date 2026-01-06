# 性能优化规范

> 本文档定义了 Workflow CLI 项目的性能优化规范和最佳实践，所有贡献者都应遵循这些规范。

---

## 📋 目录

- [概述](#-概述)
- [性能测试要求](#-性能测试要求)
- [内存使用优化规则](#-内存使用优化规则)
- [异步操作使用规则](#-异步操作使用规则)
- [相关文档](#-相关文档)

---

## 📋 概述

本文档定义了性能优化规范，包括性能测试要求、内存使用优化规则和异步操作使用规则。

### 核心原则

- **关键路径测试**：关键路径必须进行性能测试
- **内存优化**：避免不必要的内存分配
- **异步操作**：网络请求应使用异步操作
- **流式处理**：大文件处理应使用流式处理

### 使用场景

- 编写性能关键代码时参考
- 性能优化时使用
- 性能测试时参考

---

## 性能测试要求

### 关键路径性能测试

1. **识别关键路径**：
   - 识别应用中的性能关键路径（如频繁调用的函数、主循环、数据处理流程）
   - 识别用户感知明显的操作（如命令执行、文件处理、网络请求）

2. **性能测试工具**：
   - 使用 Go 标准库 `testing` 包进行基准测试（推荐）
   - 使用 `go test -bench` 运行基准测试
   - 使用 `go test -benchmem` 运行基准测试并显示内存分配

3. **性能测试实现**：

```go
// internal/lib/module/benchmark_test.go
package module

import (
    "testing"
)

func BenchmarkMyFunction(b *testing.B) {
    input := "test input"
    for i := 0; i < b.N; i++ {
        MyFunction(input)
    }
}
```

4. **性能测试要求**：
   - 关键路径必须进行性能测试
   - 性能测试应作为 CI/CD 的一部分
   - 性能测试结果应记录和跟踪

### 性能回归测试要求

1. **回归测试时机**：
   - 性能关键代码变更后必须运行性能测试
   - 重构性能关键代码后必须运行性能测试
   - 添加新依赖后评估性能影响

2. **性能阈值**：
   - 建立性能基准线
   - 设置性能回归阈值（如性能下降不超过 5%）
   - 如果性能下降超过阈值，需要优化或回滚

3. **性能监控**：
   - 使用 `criterion` 的统计功能跟踪性能趋势
   - 记录性能测试历史数据
   - 识别性能回归趋势

### 性能基准测试要求

1. **建立基准线**：
   - 为关键路径建立性能基准线
   - 记录基准测试结果（平均值、中位数、P95、P99）
   - 在文档中记录性能基准

2. **定期运行基准测试**：
   - 定期运行基准测试（如每次发布前）
   - 记录性能趋势
   - 识别性能退化

3. **基准测试工具**：

```bash
# 运行所有基准测试
go test -bench=. ./...

# 运行特定基准测试
go test -bench=BenchmarkMyFunction ./internal/lib/module

# 显示详细输出和内存分配
go test -bench=. -benchmem ./...
```

4. **基准测试最佳实践**：
   - 使用 `black_box` 防止编译器过度优化
   - 多次运行取平均值
   - 在稳定的环境中运行（避免 CPU 频率变化影响）

---

## 内存使用优化规则

### 避免不必要的内存分配

1. **优先使用栈分配**：
   - 优先使用栈分配，避免堆分配
   - 小数据结构应使用栈分配
   - 大数据结构才使用堆分配

2. **使用指针避免复制大结构体**：
   - 使用指针 `*Struct` 或接口避免大结构体复制
   - 使用切片 `[]byte` 或 `[]string` 而不是数组
   - 使用引用传递而不是值传递

```go
// ✅ 好的做法：使用指针传递
func ProcessData(data *LargeStruct) {
    // 使用指针，只复制指针，不复制数据
}

// ❌ 不好的做法：值传递大结构体
func ProcessData(data LargeStruct) {
    // 会复制整个结构体，性能较差
}
```

3. **优化内存分配**：
   - 复用缓冲区避免重复分配
   - 使用对象池（如 `sync.Pool`）复用对象
   - 避免不必要的字符串拼接

```go
// ✅ 好的做法：复用缓冲区
var buf bytes.Buffer
for i := 0; i < 1000; i++ {
    buf.Reset()  // 复用缓冲区
    buf.WriteString("data")
}

// ❌ 不好的做法：每次都创建新缓冲区
for i := 0; i < 1000; i++ {
    buf := bytes.NewBufferString("data")  // 每次都分配
}
```

### 预分配内存

1. **预分配 Vec**：
   - 使用 `Vec::with_capacity` 预分配内存
   - 如果知道大小，预分配可以避免多次重新分配

```go
// ✅ 好的做法：预分配容量
vec := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    vec = append(vec, i)
}

// ❌ 不好的做法：多次重新分配
var vec []int
for i := 0; i < 1000; i++ {
    vec = append(vec, i)  // 可能多次重新分配
}
```

2. **预分配字符串**：
   - 使用 `bytes.Buffer` 或 `strings.Builder` 预分配字符串
   - 如果知道字符串长度，预分配可以提高性能

```go
// ✅ 好的做法：预分配字符串
var builder strings.Builder
builder.Grow(100)  // 预分配容量
for i := 0; i < 100; i++ {
    builder.WriteString(strconv.Itoa(i))
}
s := builder.String()
```

3. **预分配 Map**：
   - 使用 `make(map[K]V, capacity)` 预分配 Map
   - 如果知道元素数量，预分配可以避免多次重新哈希

```go
// ✅ 好的做法：预分配 Map
m := make(map[int]int, 100)
for i := 0; i < 100; i++ {
    m[i] = i * 2
}
```

### 大文件处理

1. **使用流式处理**：
   - 大文件处理时使用 `BufReader`、`BufWriter`
   - 避免一次性将整个文件加载到内存

```go
import (
    "bufio"
    "os"
)

// ✅ 好的做法：流式处理
file, err := os.Open("large_file.txt")
if err != nil {
    return err
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    processLine(line)
}

// ❌ 不好的做法：一次性加载整个文件
content, err := os.ReadFile("large_file.txt")  // 可能内存不足
```

---

## 异步操作使用规则

### 网络请求应使用异步操作

1. **使用异步 HTTP 客户端**：
   - 网络请求应使用 goroutine 进行并发处理（如 `net/http` 配合 goroutine）
   - 避免阻塞主线程

```go
// ✅ 好的做法：使用 goroutine 并发处理
func fetchData(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(body), nil
}

// ✅ 好的做法：并发处理多个请求
func fetchMultiple(urls []string) ([]string, error) {
    var wg sync.WaitGroup
    results := make([]string, len(urls))
    errs := make([]error, len(urls))

    for i, url := range urls {
        wg.Add(1)
        go func(idx int, u string) {
            defer wg.Done()
            data, err := fetchData(u)
            results[idx] = data
            errs[idx] = err
        }(i, url)
    }

    wg.Wait()

    // 检查错误
    for _, err := range errs {
        if err != nil {
            return nil, err
        }
    }

    return results, nil
}
```

---

## 🔍 故障排除

### 问题 1：性能测试结果不稳定

**症状**：性能测试结果波动较大

**解决方案**：

1. 使用 `black_box` 防止编译器过度优化
2. 多次运行取平均值
3. 在稳定的环境中运行（避免 CPU 频率变化影响）

### 问题 2：内存使用过高

**症状**：程序内存使用过高

**解决方案**：

1. 检查是否有不必要的内存分配
2. 使用引用而非拥有所有权
3. 使用流式处理大文件

---

## 📚 相关文档

### 开发规范

- [代码风格规范](../code-style.md) - 代码风格规范
- [模块组织规范](../module-organization.md) - 模块组织规范

### 工具文档

- [Go 基准测试](https://pkg.go.dev/testing#hdr-Benchmarks) - Go 标准库基准测试文档
- [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) - 基准测试结果分析工具

---

## ✅ 检查清单

使用本规范时，请确保：

- [ ] 关键路径已进行性能测试
- [ ] 避免不必要的内存分配
- [ ] 网络请求使用异步操作
- [ ] 大文件处理使用流式处理

---

**最后更新**: 2025-01-27

