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
