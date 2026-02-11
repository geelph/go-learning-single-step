package main

import (
	"context"
	"log/slog"
	"os"
	"time"
)

// 定义一个新的类型作为 context 的键，以避免键冲突
type contextKey string

const loggerKey = contextKey("logger")

// main 是程序的入口点，用于演示 slog 日志库的使用方法
// 配置默认日志处理器为 JSON 格式，并自定义时间戳格式
// 创建带预设属性的日志记录器，并在不同场景下使用它记录日志
func main() {
	// 配置slog日志处理器选项
	// 设置不添加源码信息、日志级别为Debug，并自定义时间属性格式
	handlerOptions := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		// 自定义属性替换函数，将时间格式转换为Unix时间戳
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.Int64Value(time.Now().Unix()),
				}
			}
			return a
		},
	}

	// 创建JSON格式的日志处理器，输出到标准错误流
	jsonHandler := slog.NewJSONHandler(os.Stderr, handlerOptions)
	// 创建带有应用程序名称属性的logger实例
	loggerWithAppName := slog.New(jsonHandler).With(slog.String("log_app_name", "test"))
	// 将配置好的logger设置为默认全局logger
	slog.SetDefault(loggerWithAppName)

	// 创建一个带有路径属性的新日志记录器
	logger := slog.With(slog.String("path", "/ping"))

	// 使用新创建的日志记录器记录错误信息
	logger.Error("failed", slog.Int("aaa", 2))

	// 将日志记录器存储到上下文中
	ctx := context.WithValue(context.Background(), loggerKey, logger)

	// 使用默认日志记录器记录调试信息
	slog.Debug("aaaa")

	// 在测试函数中使用上下文中的日志记录器
	testCtx(ctx)
}

// testCtx 测试使用上下文中的日志记录器进行日志记录
// 该函数演示了如何从上下文中获取日志记录器并在不同场景下使用它
//
// 参数:
//   - ctx: 包含日志记录器的上下文
func testCtx(ctx context.Context) {
	// 从上下文中获取日志记录器
	logger := getLoggerByCtx(ctx)
	// 使用上下文相关的错误日志记录功能记录错误信息
	slog.ErrorContext(ctx, "testCtx")
	// 使用从上下文获取的日志记录器记录调试信息
	logger.Debug("testCtx")
}

// getLoggerByCtx 从给定的上下文中获取日志记录器
// 该函数尝试从上下文 context 中取出预先存储的日志记录器实例
//
// 参数:
//   - ctx: 包含日志记录器的上下文，通过 contextKey 存储
//
// 返回值:
//   - 如果成功从上下文获取到日志记录器则返回该记录器
//   - 如果无法从上下文获取日志记录器则返回默认日志记录器
func getLoggerByCtx(ctx context.Context) *slog.Logger {
	// 从上下文中尝试获取日志记录器
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	// 如果类型断言失败，即没有找到日志记录器或类型不匹配，则返回默认日志记录器
	if !ok {
		return slog.Default()
	}
	// 类型断言成功，返回从上下文获取到的日志记录器
	return logger
}
