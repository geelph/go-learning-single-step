package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

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
