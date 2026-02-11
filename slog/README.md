# 描述

`log/slog` 是 Go 语言标准库中自 Go 1.21 版本起引入的一个结构化日志记录包。

它旨在提供一个现代化、高性能、结构化（structured logging）的日志记录解决方案，作为对传统 `log` 包的补充和升级，以满足更复杂应用的需求。

# 核心组件

1. ** Logger (**`*slog.Logger`** )** : 日志记录器实例。你可以创建多个具有不同配置（如级别、处理器、属性）的 Logger。
2. ** Handler (**`slog.Handler`** )** : 处理日志记录的实际后端。它决定日志的格式（如文本、JSON）和输出目的地（如标准输出、文件、网络等）。标准库提供了 `TextHandler` 和 `JSONHandler` 。
3. ** Level (**`slog.Level`** )** : 日志级别，用于控制日志的严重程度和过滤。常见的级别有：
   - `LevelDebug` (-4)
   - `LevelInfo` (0) - 默认级别
   - `LevelWarn` (4)
   - `LevelError` (8)
4. ** Record (**`slog.Record`** )** : 代表一条日志记录的结构体，包含时间、级别、消息、属性等信息。通常由 `Handler` 处理，用户较少直接操作。

# 基本用法

## 使用默认 Logger

最简单的用法是直接使用包级别的函数，它们使用一个默认配置的 `Logger`（通常是 `TextHandler` 输出到 `os.Stderr`，级别为 `Info`）。

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // 使用包级别函数
    slog.Info("Application started")
    slog.Warn("This is a warning", "user_id", 12345)
    slog.Error("Something went wrong", "error_code", "E1001", "retry", true)
}
```

输出：

```plain
2025/09/19 16:12:26 INFO Application started
2025/09/19 16:12:26 WARN This is a warning user_id=12345
2025/09/19 16:12:26 ERROR Something went wrong error_code=E1001 retry=true
```

## 设置日志级别

自定义一个 Handler 设置输出和日志等级并将其设置为默认日志

```go
package main

import (
	"log/slog"
	"os"
)

// main 函数是程序的入口点，演示了 slog 日志库的基本使用方法
// 该函数不接受任何参数，也没有返回值
func main() {
	// 创建一个新的文本处理程序，并设置日志级别为 Info
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	// 设置全局的日志处理程序
	slog.SetDefault(slog.New(handler))

	// 记录不同级别的日志信息
	// Debug 级别日志由于低于设置的 Info 级别，因此不会被输出
	slog.Debug("这是一条调试日志", "debugKey", "debugValue")
	// Info 级别日志符合设置的日志级别要求，会被正常输出
	slog.Info("这是一条信息日志", "infoKey", "infoValue")
}


```

输出：

```plain
time=2025-09-19T17:00:31.190+08:00 level=INFO msg=这是一条信息日志 infoKey=infoValue
```

## 创建自定义 Logger

你可以创建带有特定 `Handler`、级别和默认属性的 `Logger`。

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	// 创建一个输出到标准输出、使用 JSON 格式的 Handler
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)

	// 使用自定义 Logger
	logger.Info("User logged in", "username", "alice", "ip", "192.168.1.1")
	// 输出类似: {"time":"2023-10-27T10:00:00.000Z","level":"INFO","msg":"User logged in","username":"alice","ip":"192.168.1.1"}

	// 设置默认属性（会附加到每条日志）
	opts := &slog.HandlerOptions{
		AddSource: true,            // 添加源码位置
		Level:     slog.LevelDebug, // 设置最低记录级别为 Debug
	}
	textHandler := slog.NewTextHandler(os.Stderr, opts)
	appLogger := slog.New(textHandler).With("service", "auth-service", "version", "1.0")

	appLogger.Info("Service initialized")
	// 输出类似: time=... level=INFO source="main.go:25" msg="Service initialized" service=auth-service version=1.0
}

```

输出：

```plain
{"time":"2025-09-19T16:12:57.9142786+08:00","level":"INFO","msg":"User logged in","username":"alice","ip":"192.168.1.1"}
time=2025-09-19T16:12:57.934+08:00 level=INFO source=E:/Source/TestScript/GoProjects/test/main.go:25 msg="Service initialized" service=auth-service version=1.0
```

## 使用 With 添加上下文

`Logger.With` 方法可以创建一个带有预设属性的新 `Logger`，这对于为特定组件或请求添加上下文非常有用。

```go
func handleRequest(logger *slog.Logger, reqID string) {
    // 为这个请求创建一个带有 request_id 上下文的 Logger
    reqLogger := logger.With("request_id", reqID)

    reqLogger.Info("Processing request")
    // ... 处理逻辑 ...
    if err := doSomething(); err != nil {
        reqLogger.Error("Failed to process request", "error", err)
        return
    }
    reqLogger.Info("Request processed successfully")
}
```

## 记录错误

虽然可以用 `slog.Error("message", "error", err)`，但更推荐使用 `slog.Any("key", value)` 或直接利用 `slog` 对 `error` 类型的特殊处理（如果 `error` 是第一个参数，它会自动作为 `msg` 的补充或单独的 `err `字段，具体取决于 `Handler`）。

```go
if err := someOperation(); err != nil {
    slog.Error("Operation failed", slog.Any("error", err))
    // 或者更简洁地（slog 会自动处理 error 类型）
    slog.Error("Operation failed", "error", err)
    // 对于 JSONHandler，可能会输出 err 字段
}
```

# 优势

- **结构化**: 日志是键值对形式，便于机器解析、搜索、过滤和分析。
- **高性能**: 设计时考虑了性能，尤其是在高吞吐量场景下。
- **灵活**: 易于扩展 Handler 以支持不同的格式和输出目标。
- **标准化**: 作为标准库的一部分，减少了对外部依赖的需求，促进了生态系统的统一。
- **上下文**: 通过 With 方法轻松添加和传递上下文信息。

# 总结

`log/slog` 是 Go 语言在日志记录方面的一个重要进步。它鼓励使用结构化日志，提供了良好的性能和灵活性。对于新项目，推荐优先考虑使用 `slog`。对于现有项目，可以根据需要逐步迁移到 `slog` 或将其与现有日志库结合使用。

# 参考链接

- [Go 语言中 slog 日志库的全面使用指南\_go slog-CSDN博客](https://blog.csdn.net/m0_57836225/article/details/145715828)
