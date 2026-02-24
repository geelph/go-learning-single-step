# Elo（埃罗等级分系统）算法实现

## 简介

Elo评分系统是一种用于计算对弈双方相对技能水平的方法，最初用于国际象棋，现在也被广泛应用于各种竞技游戏中。本项目实现了基本的Elo评分算法，可以根据对战结果更新两位玩家的评分。

## 基本概念

### 核心公式

Elo系统的核心在于预测胜率和更新评分：

1. 计算预期胜率：
   ```
   Ea = 1 / (1 + 10^((Rb - Ra) / 400))
   Eb = 1 - Ea
   ```

2. 更新评分：
   ```
   新Ra = 旧Ra + K × (实际结果 - 预期结果)
   新Rb = 旧Rb + K × (实际结果 - 预期结果)
   ```

### 参数说明

- **K**: K因子，控制评分变化幅度，默认值为32
- **D**: 偏差系数，默认值为400
- **Ra**: A玩家的当前评分
- **Rb**: B玩家的当前评分
- **Sa**: A玩家的实际比赛结果（胜=1，平=0.5，负=0）
- **Sb**: B玩家的实际比赛结果（胜=1，平=0.5，负=0）

## 代码结构

### 常量定义

```go
const (
    // K is the default K-Factor
    K = 32
    // D is the default deviation
    D = 400
)
```

### Elo结构体

```go
type Elo struct {
    A  uint32    // A玩家当前的Rating
    B  uint32    // B玩家当前的Rating
    Sa float64   // 实际胜负值，胜=1，平=0.5，负=0
                 // 传入值默认A的胜负 / 1 A胜利 B失败 / 0 B胜利 A失败
}
```

### 主要函数

#### `EloRating(elo Elo) (a uint32, b uint32)`

根据输入的Elo结构体计算新的评分，返回更新后的A、B玩家评分。

#### `Decimal(value float64, f string) float64`

格式化浮点数，保留指定精度的小数位。

### 函数测试

```go
package elo

import "testing"

func Test_EloRating(t *testing.T) {
	a, b := EloRating(Elo{
		A:  1500,
		B:  1600,
		Sa: 1,
	})
	t.Log("a", a)
	t.Log("b", b)
}

func Test_Decimal(t *testing.T) {
	t.Log(Decimal(22.222222222, "%.2f"))
	t.Log(Decimal(22.222222222, "%.0f"))
	t.Log(Decimal(22.6666666666, "%.2f"))
	t.Log(Decimal(22.66666666666, "%.0f"))
}

```

运行结果：

```bash
# go test -v elo_test.go elo.go
=== RUN   Test_EloRating
    elo_test.go:11: a 1520
    elo_test.go:12: b 1580
--- PASS: Test_EloRating (0.00s)
=== RUN   Test_Decimal
    elo_test.go:16: 22.22
    elo_test.go:17: 22
    elo_test.go:18: 22.67
    elo_test.go:19: 23
--- PASS: Test_Decimal (0.00s)
PASS
ok      command-line-arguments  0.751s
```

## 使用示例

```go
package main

import (
    "fmt"
    "your-project/elo"
)

func main() {
    // 创建Elo对象，A玩家评分为1500，B玩家评分为1400
    // 假设A获胜（Sa=1.0）
    game := elo.Elo{
        A:  1500,      // A玩家当前评分
        B:  1400,      // B玩家当前评分
        Sa: 1.0,       // A获胜，实际结果为1.0
    }

    // 计算新评分
    newA, newB := elo.EloRating(game)

    fmt.Printf("A玩家新评分: %d\n", newA)  // A玩家新评分
    fmt.Printf("B玩家新评分: %d\n", newB)  // B玩家新评分
}
```

## 算法流程

1. **计算预期胜率**：
   - 根据两玩家的评分差异计算A玩家的预期胜率Ea
   - B玩家的预期胜率为Eb = 1 - Ea

2. **确定实际结果**：
   - 如果A获胜，则Sa=1.0，Sb=0.0
   - 如果平局，则Sa=0.5，Sb=0.5
   - 如果B获胜，则Sa=0.0，Sb=1.0

3. **更新评分**：
   - 根据K因子和实际结果与预期结果的差异更新评分
   - 实际结果优于预期时，评分上升；反之则下降

## 应用场景

- 游戏中的排位系统
- 竞技比赛的等级评定
- 在线对战平台的匹配机制
- 任何需要评估两个实体相对实力的场景

## 注意事项

- K因子的选择会影响评分变动的速度，K值越大，评分变化越快
- 初始评分设定对后续评分有一定影响
- 该实现中使用了四舍五入处理小数部分

## 参考资料

- [Elo Rating System - Wikipedia](https://en.wikipedia.org/wiki/Elo_rating_system)
- [How Chess Ratings Work](https://www.chess.com/terms/elo-rating-chess)