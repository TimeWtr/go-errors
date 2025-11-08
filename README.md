# go-errors
go-errors 是一个功能丰富的 Go 错误处理库，专为生产环境设计。它提供结构化错误、堆栈跟踪、监控集成和框架支持，帮助您构建更可靠、更易维护的 Go 应用程序。

## 🚀 特性概览
### 核心功能
🏗️ 结构化错误 - 统一的错误接口，支持错误码、类型分类和元数据

🔍 完整堆栈跟踪 - 自动捕获调用堆栈，支持开发和生产环境优化

📊 监控集成 - 内置 Prometheus 指标，支持错误统计和告警

🔄 错误包装 - 保持错误链，支持 errors.Is 和 errors.As

🎯 类型安全 - 预定义错误码，编译期检查

### 生产环境就绪
🚀 高性能 - 零内存分配的错误创建，配置化的堆栈跟踪

🔧 可配置 - 支持环境特定的错误行为和显示级别

📈 可观测性 - 结构化日志、分布式追踪、指标收集

🛡️ 安全 - 生产环境敏感信息过滤

### 框架集成
🌐 HTTP 支持 - 开箱即用的 Gin 中间件

🔌 可扩展 - 支持自定义错误处理器和监控后端

📋 标准化 - 统一的 API 错误响应格式

## 📦 快速开始
### 安装
```bash
go get github.com/your-org/go-errors
```
### 基础使用
```go
package main

import (
    "fmt"

    "github.com/TimeWtr/go-errors"
)

func main() {
    // 创建预定义错误
    err := errors.New(errors.ErrUserNotFound)
    fmt.Printf("Error: %s, Code: %s\n", err.Error(), err.Code())

    // 带上下文的错误
    err = errors.UserNotFound("user-123").
        WithMetadata("operation", "get_user").
        WithMetadata("attempt", 3)
    
    // 错误包装
    if _, err := someOperation(); err != nil {
        wrapped := errors.Wrapf(err, errors.ErrInternal, "operation failed")
        fmt.Printf("Wrapped: %s\n", wrapped.Error())
    }
}
```