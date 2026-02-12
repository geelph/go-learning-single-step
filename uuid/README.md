# Google UUID 库 (github.com/google/uuid) Go 语言使用

> 📌 **官方仓库**：https://github.com/google/uuid  
> ✅ **推荐理由**：Google 官方维护、线程安全、支持 UUID 多版本、社区活跃（替代已归档的 `satori/go.uuid`）

---

## 🔧 一、安装

```bash
go get github.com/google/uuid
```

---

## 🚀 二、快速入门（最常用场景）

```go
package main

import (
    "fmt"
    "github.com/google/uuid"
)

func main() {
    // 生成随机 UUID (v4) - ✅ 最推荐日常使用
    id := uuid.New() // 返回 uuid.UUID 类型（无 error，内部使用 crypto/rand）
    fmt.Println("UUID v4:", id)          // 标准格式: 550e8400-e29b-41d4-a716-446655440000
    fmt.Println("字符串:", id.String())  // 同上
    fmt.Println("URN格式:", id.URN())    // urn:uuid:550e8400-e29b-41d4-a716-446655440000
}
```

---

## 📚 三、核心功能详解

### 1️⃣ 生成不同版本 UUID

| 版本   | 函数                     | 说明                       | 适用场景                                 |
| ------ | ------------------------ | -------------------------- | ---------------------------------------- |
| **v1** | `uuid.NewUUID()`         | 基于时间戳 + MAC 地址      | 需要时间顺序（⚠️容器环境可能失败）       |
| **v3** | `uuid.NewMD5(ns, name)`  | MD5 哈希（命名空间+名称）  | 确定性生成（⚠️MD5 已不推荐用于安全场景） |
| **v4** | `uuid.New()`             | **密码学安全随机数**       | ✅ 通用首选（数据库主键、会话ID等）      |
| **v5** | `uuid.NewSHA1(ns, name)` | SHA1 哈希（命名空间+名称） | 确定性生成（比 v3 更安全）               |
| **v6** | `uuid.NewV6()`           | 有序时间戳（重排 v1）      | 需时间排序 + 避免 MAC 泄露               |
| **v7** | `uuid.NewV7()`           | **时间戳 + 随机**          | ✅ 新项目推荐（有序、安全、无硬件依赖）  |
| **v8** | `uuid.NewV8(data)`       | 自定义字节                 | 特殊需求                                 |

```go
// v5 示例：基于 DNS 命名空间生成确定性 UUID
ns := uuid.NameSpaceDNS // 预定义命名空间: DNS/URL/OID/X500
name := "example.com"
u5 := uuid.NewSHA1(ns, []byte(name))
fmt.Println("UUID v5:", u5)

// v7 示例（Go 1.20+ 推荐）
u7, err := uuid.NewV7()
if err != nil {
    panic(err)
}
fmt.Println("UUID v7 (时间有序):", u7)
```

### 2️⃣ 解析与验证

```go
// 安全解析（带错误处理）
u, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
if err != nil {
    panic("无效 UUID 格式")
}

// MustParse（测试/确定有效时使用，失败 panic）
u = uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
```

### 3️⃣ 实用操作

```go
// 比较
if u1 == u2 { /* 相等 */ }

// 转字节切片
uBytes := u1[:] // []byte 长度 16

// 无连字符格式（需自行处理）
uNoDash := strings.ReplaceAll(u1.String(), "-", "")

// 大写格式
uUpper := strings.ToUpper(u1.String())
```

### 4️⃣ 执行效果

```bash
UUID v4: e73a7988-8b02-4705-802d-4d438394a30c
字符串: e73a7988-8b02-4705-802d-4d438394a30c
URN格式: urn:uuid:e73a7988-8b02-4705-802d-4d438394a30c
UUID v5: cfbff0d1-9375-5685-968c-48ce8b15ae17
UUID v7 (时间有序): 019c4fe8-2c8b-7e45-96dc-00ba5a362c50
有效 UUID 格式: 550e8400-e29b-41d4-a716-446655440000
有效 UUID 格式: a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
转字节切片：
a0eebc999c0b4ef8bb6d6bb9bd380a11
无连字符格式: a0eebc999c0b4ef8bb6d6bb9bd380a11
UUID 大写格式: A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A11
```

---

## ⚠️ 四、关键注意事项

1. **安全性**
   - `uuid.New()` (v4) 和 `NewV7()` 使用 `crypto/rand`，**密码学安全** ✅
   - 避免自行用 `math/rand` 生成（不安全！）
2. **环境限制**
   - v1 在无网络接口环境（如 Docker）可能失败 → 改用 v4/v7
3. **版本选择建议**
   - 通用场景：**v4**（简单安全）
   - 需要时间排序（如数据库主键）：**v7**（现代首选）或 v6
   - 确定性生成：v5（SHA1）优于 v3（MD5）
4. **线程安全**：所有函数均安全，可并发调用 ✅
5. **大小写**：标准输出为小写，需大写请用 `strings.ToUpper`

---

## ❓ 五、常见问题

**Q：如何验证字符串是否为有效 UUID？**  
A：直接 `uuid.Parse(s)`，返回 error 即无效。

**Q：生成的 UUID 有连字符，能去掉吗？**  
A：库不直接提供，用 `strings.ReplaceAll(u.String(), "-", "")` 处理（注意：存储/传输时建议保留标准格式）。

**Q：与 `golang.org/x/exp/uuid` 区别？**  
A：`google/uuid` 是稳定生产级库；`x/exp/uuid` 属实验性模块，**不推荐生产使用**。

**Q：性能如何？**  
A：v4/v7 生成极快（纳秒级），实测 100 万次生成 < 1 秒，无需担忧。

---

## 📖 六、延伸学习

- [官方文档](https://pkg.go.dev/github.com/google/uuid)
- [UUID 标准 RFC 4122](https://datatracker.ietf.org/doc/html/rfc4122)
- 数据库优化：若用 UUID 作主键，**强烈推荐 v7**（减少 InnoDB 页分裂）

> 💡 **最佳实践总结**：新项目优先选 **`uuid.NewV7()`**（有序+安全+无环境依赖），通用场景用 **`uuid.New()`**（v4）。避免在安全敏感场景使用 v3/v1。
