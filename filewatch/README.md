# 文件系统监控 - fsnotify 使用教程

本项目演示了如何使用 `fsnotify` 库来监控文件系统的变化。`fsnotify` 是一个跨平台的文件系统监控库，可以监听文件或目录的各种事件。

## 安装

```bash
go get github.com/fsnotify/fsnotify
```

## 基本用法

### 1. 创建监控器

```go
watcher, err := fsnotify.NewWatcher()
if err != nil {
    log.Fatal(err)
}
defer watcher.Close()
```

### 2. 添加要监控的目录或文件

```go
err = watcher.Add("/path/to/directory")
if err != nil {
    log.Fatal(err)
}
```

### 3. 监听事件

```go
go func() {
    for {
        select {
        case event := <-watcher.Events:
            log.Println("Event:", event)
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Println("Modified file:", event.Name)
            }
        case err := <-watcher.Errors:
            log.Println("Error:", err)
        }
    }
}()
```

## 事件类型

`fsnotify` 支持以下几种事件类型：

- `fsnotify.Create`: 文件或目录被创建
- `fsnotify.Write`: 文件被写入（内容更改）
- `fsnotify.Remove`: 文件或目录被删除
- `fsnotify.Rename`: 文件或目录被重命名
- `fsnotify.Chmod`: 文件权限被修改

可以通过位运算符组合检测特定事件：

```go
if event.Op&fsnotify.Create == fsnotify.Create {
    // 处理创建事件
}
```

## 完整示例

以下是一个完整的示例，展示了如何监控一个目录及其子目录的所有变化：

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// 监控服务类型
type Watch struct {
	watch *fsnotify.Watcher
}

func main() {
	// 创建一个文件监控服务
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("创建监控失败: ", err)
	}
	defer watch.Close()

	// 创建监控服务
	w := Watch{
		watch: watch,
	}

	// 检查并创建监控目录
	watchDir := "./watch" // 修改为相对简单的路径
	if _, err := os.Stat(watchDir); os.IsNotExist(err) {
		log.Printf("%s 目录不存在，创建中...", watchDir)
		if err := os.MkdirAll(watchDir, 0755); err != nil {
			log.Fatal("监控目录创建失败: ", err)
		}
	}

	// 指定监控目录
	w.watchDir(watchDir)
	select {}
}

// 监控目录
func (w *Watch) watchDir(dir string) {
	log.Println("开始监控目录 : ", dir)

	// 使用filepath.Walk遍历目录并添加监控
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("访问路径 %s 时发生错误: %v", path, err)
			return nil // 继续遍历其他目录
		}

		if info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				log.Printf("获取绝对路径失败 %s: %v", path, err)
				return err
			}

			err = w.watch.Add(absPath)
			if err != nil {
				log.Printf("添加监控失败 %s: %v", absPath, err)
				return err
			}
			log.Printf("已添加监控: %s", absPath)
		}
		return nil
	})

	if err != nil {
		log.Printf("遍历目录时出错: %v", err)
		return
	}

	log.Println("监控服务已经启动")
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("创建文件 : ", ev.Name)
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							err = w.watch.Add(ev.Name)
							if err != nil {
								log.Printf("添加新目录监控失败: %v", err)
							} else {
								fmt.Println("添加监控 : ", ev.Name)
							}
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("写入文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Println("删除文件 : ", ev.Name)
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							err = w.watch.Remove(ev.Name)
							if err != nil {
								log.Printf("移除目录监控失败: %v", err)
							} else {
								fmt.Println("删除监控 : ", ev.Name)
							}
						}
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						fmt.Println("重命名文件 : ", ev.Name)
						// 移除监控前判断文件是否在监控列表中
						for _, path := range w.watch.WatchList() {
							if path == ev.Name {
								// 移除监控中被重命名的文件
								err := w.watch.Remove(ev.Name)
								if err != nil {
									log.Printf("移除重命名文件监控失败: %v", err)
								}
								break
							}
						}
					}
				}
			case err := <-w.watch.Errors:
				{
					fmt.Println("error : ", err)
					return
				}
			}
		}
	}()
}

```

执行效果：

```bash
# go run main.go
2026/02/12 15:31:26 开始监控目录 :  ./watch
2026/02/12 15:31:26 已添加监控: E:\Source\TestScript\GoProjects\go-learning-single-step\filewatch\watch
2026/02/12 15:31:26 监控服务已经启动
创建文件 :  E:\Source\TestScript\GoProjects\go-learning-single-step\filewatch\watch\新建 文本文档.txt
重命名文件 :  E:\Source\TestScript\GoProjects\go-learning-single-step\filewatch\watch\新建 文本文档.txt
创建文件 :  E:\Source\TestScript\GoProjects\go-learning-single-step\filewatch\watch\test.txt
写入文件 :  E:\Source\TestScript\GoProjects\go-learning-single-step\filewatch\watch\test.txt
删除文件 :  E:\Source\TestScript\GoProjects\go-learning-single-step\filewatch\watch\test.txt
```

## 常见应用场景

1. **实时同步**: 当文件发生变化时，立即同步到其他位置
2. **热重载**: 在开发过程中，当代码改变时自动重新加载应用程序
3. **日志监控**: 监控日志文件的变化并实时处理
4. **备份系统**: 实时监控文件变化并备份

## 注意事项

- 监控目录需要相应的读取权限
- 不同操作系统对监控数量有限制
- 需要注意处理监控器关闭和资源释放
- 对于递归监控，需要手动遍历子目录并添加监控
- 监控事件是异步的，需要使用 goroutine 处理

## 参考资料

- [fsnotify GitHub 仓库](https://github.com/fsnotify/fsnotify)
- [fsnotify 文档](https://pkg.go.dev/github.com/fsnotify/fsnotify)
