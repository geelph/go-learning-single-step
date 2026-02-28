# urfave/cli/v3 使用指南

`urfave/cli` 是一个简单、快速且有趣的包，用于构建命令行应用程序。

## 安装

确保你已经安装了 Go (版本不低于 1.18)，然后创建一个新的项目：

```bash
mkdir myproject
cd myproject
go mod init myproject
go get github.com/urfave/cli/v3
```

## 命令

### 一行代码

```go
package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	(&cli.Command{}).Run(context.Background(), os.Args)
}

```

运行结果：

```bash
# go run .
NAME:
   cli.exe - A new cli application

USAGE:
   cli.exe [global options]

GLOBAL OPTIONS:
   --help, -h  show help
```

### 简单示例-不获取参数

最简单的 CLI 应用程序如下所示：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

```

运行结果：

```bash
# go run .
boom! I say!

# go run . -h
NAME:
   boom - make an explosive entrance

USAGE:
   boom [global options]

GLOBAL OPTIONS:
   --help, -h  show help
```

### 简单示例-获取参数

```go
cmd := &cli.Command{
  Name:    "print",
  Aliases: []string{"p"},
  Usage:   "打印参数",
  Action: func(ctx context.Context, c *cli.Command) error {
    fmt.Printf("打印参数: %v\n", c.Args().Slice())
    return nil
  },
}
```

运行结果：

```bash
# go run . -h
NAME:
   print - Printing parameters

USAGE:
   print [global options]

GLOBAL OPTIONS:
   --help, -h  show help
   

```



## 子命令

子命令是嵌套在主命令下的命令，用于构建更复杂的命令行接口。

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app := &cli.Command{
		Name:  "host",
		Usage: "manage hosts",
		Commands: []*cli.Command{
			{
				Name:    "configure",
				Aliases: []string{"c"},
				Usage:   "configure host",
				Commands: []*cli.Command{
					{
						Name:  "ssh",
						Usage: "configure SSH access",
						Action: func(ctx context.Context, c *cli.Command) error {
							fmt.Println("Configuring SSH access...")
							return nil
						},
					},
					{
						Name:  "ftp",
						Usage: "configure FTP access",
						Action: func(ctx context.Context, c *cli.Command) error {
							fmt.Println("Configuring FTP access...")
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

```

运行结果：

```bash
# go run .
NAME:
   host - manage hosts

USAGE:
   host [global options] [command [command options]]

COMMANDS:
   configure, c  configure host
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help


# go run . configure
NAME:
   host configure - configure host

USAGE:
   host configure [command [command options]]

COMMANDS:
   ssh  configure SSH access
   ftp  configure FTP access

OPTIONS:
   --help, -h  show help


# go run . configure ssh
Configuring SSH access...
```



## 标志

标志用于向命令传递选项。

### 字符串标志

```go
&cli.StringFlag{
  Name:     "config",
  Aliases:  []string{"c"},
  Value:    "config.json",
  Usage:    "配置文件路径",
  Required: false,
}
```

### 整数标志

```go
&cli.IntFlag{
  Name:  "port",
  Value: 8080,
  Usage: "监听端口",
}
```

### 布尔标志

```go
&cli.BoolFlag{
  Name:  "verbose",
  Usage: "启用详细输出",
}
```

### 获取标志值

```go
Action: func(ctx context.Context, cmd *cli.Command) error {
  port := cmd.Int("port")              // 获取整数值
  configFile := cmd.String("config")   // 获取字符串值
  verbose := cmd.Bool("verbose")       // 获取布尔值

  fmt.Printf("端口: %d, 配置文件: %s, 详细模式: %t\n", port, configFile, verbose)
  return nil
}
```

## 上下文

urfave/cli 使用 Go 的 context 包来处理取消和超时。

```go
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app := &cli.Command{
		Name:  "timeout",
		Usage: "Presentation context timeout",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			log.Println("Starting task...")
			// 模拟长时间运行的任务
			select {
			case <-time.After(15 * time.Second):
				log.Println("Task completed.")
			case <-ctx.Done():
				log.Println("The task was cancelled or timed out.")
			}
			return nil
		},
	}

	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
```

运行结果：

```bash
E:\Source\TestScript\GoProjects\go-learning-single-step\cli (master -> origin)
# go run .
2026/02/27 17:19:40 Task completed.

E:\Source\TestScript\GoProjects\go-learning-single-step\cli (master -> origin)
# go run .
2026/02/27 17:19:56 Starting task...
2026/02/27 17:20:01 Task completed.

E:\Source\TestScript\GoProjects\go-learning-single-step\cli (master -> origin)
# go run .
2026/02/27 17:20:24 Starting task...
2026/02/27 17:20:34 The task was cancelled or timed out.
```



## 生命周期回调

urfave/cli 支持在命令执行前后运行函数。

### Before 回调

在命令执行前运行，可用于设置、验证等：

```go
&cli.Command{
  Name:  "template",
  Usage: "模板命令",
  Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
    fmt.Println("在命令执行前做一些准备工作...")
    // 这里可以进行验证、初始化资源等
    return ctx, nil
  },
  Action: func(ctx context.Context, cmd *cli.Command) error {
    fmt.Println("执行核心命令...")
    return nil
  },
  After: func(ctx context.Context, cmd *cli.Command) error {
    fmt.Println("清理资源...")
    return nil
  },
}
```

### After 回调

在命令执行后运行，可用于清理资源：

```go
After: func(ctx context.Context, cmd *cli.Command) error {
  fmt.Println("执行后的清理工作...")
  return nil
},
```

## 实际应用示例

以下是我们项目中使用的实际示例，展示了微服务管理工具的功能：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := &cli.Command{
		Name:  "microctl",
		Usage: "微服务集群管理工具",
		Commands: []*cli.Command{
			{
				Name:    "deploy",
				Aliases: []string{"d"},
				Usage:   "部署微服务",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "image", Required: true, Usage: "Docker镜像地址"},
					&cli.IntFlag{Name: "replicas", Value: 3, Usage: "副本数量"},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					image := c.String("image")
					replicas := c.Int("replicas")
					fmt.Printf("🚀 部署镜像 %s，副本数: %d\n", image, replicas)
					return nil
				},
			},
			{
				Name:  "scale",
				Usage: "动态扩缩容",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "service", Required: true},
					&cli.IntFlag{Name: "count", Required: true},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Printf("🔄 调整服务 %s 至 %d 实例\n",
						c.String("service"), c.Int("count"))
					return nil
				},
			},
			{
				Name:  "secure",
				Usage: "安全操作",
				Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
					fmt.Println("🔒 执行前：验证权限...")
					// 可在此处添加认证/日志逻辑
					return ctx, nil // 修复：返回正确的 context
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Println("✅ 执行核心逻辑")
					return nil
				},
				After: func(ctx context.Context, c *cli.Command) error {
					fmt.Println("📝 执行后：记录审计日志")
					return nil
				},
			},
		},
	}

	// 传递带超时的 context
	if err := cmd.Run(ctx, os.Args); err != nil {
		if err == context.DeadlineExceeded {
			log.Fatal("⏰ 命令执行超时！")
		}
		log.Fatal(err)
	}
}

```

## 运行应用

编译并运行您的应用程序：

```bash
go run main.go
```

尝试不同的命令和选项：

```bash
# 显示帮助
go run main.go --help

# 运行 deploy 命令
go run main.go deploy --image nginx --replicas 5

# 使用别名
go run main.go d --image redis

# 运行 scale 命令
go run main.go scale --service my-service --count 10

# 运行 secure 命令
go run main.go secure
```

## 高级特性

### 自定义帮助文本

```go
&cli.Command{
  Name:  "serve",
  Usage: "启动服务",
  Description: `启动服务器并监听指定端口。

  示例:
    {{.HelpName}} -p 8080
    {{.HelpName}} --port 9000`,
  Flags: []cli.Flag{
    &cli.IntFlag{Name: "port", Value: 8080, Usage: "监听端口"},
  },
}
```

### 全局标志

全局标志可用于任何命令：

```go
cmd := &cli.Command{
    Flags: []cli.Flag{
        &cli.StringFlag{Name: "host", Value: "localhost", Usage: "服务器主机地址"},
    },
    Commands: []*cli.Command{
        {
            Name:  "serve",
            Usage: "启动服务器",
            Action: func(ctx context.Context, c *cli.Command) error {
                host := c.String("host") // 使用全局标志
                fmt.Printf("服务器将在 %s 上启动\n", host)
                return nil
            },
        },
    },
}
```

## 总结

urfave/cli/v3 提供了一套简洁而强大的 API 来构建命令行应用程序。它支持:

- 嵌套命令
- 标志 (选项)
- 上下文管理
- 生命周期回调 (Before/After)
- 自定义帮助文本
- 全局标志
- 别名

## 参考链接

-  [Getting Started - urfave/cli](https://cli.urfave.org/v3/getting-started/)
- 📦 [GitHub 仓库](https://github.com/urfave/cli)
