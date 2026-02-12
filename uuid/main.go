package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func main() {
	// 生成随机 UUID (v4) - ✅ 最推荐日常使用
	id := uuid.New()                 // 返回 uuid.UUID 类型（无 error，内部使用 crypto/rand）
	fmt.Println("UUID v4:", id)      // 标准格式: 550e8400-e29b-41d4-a716-446655440000
	fmt.Println("字符串:", id.String()) // 同上
	fmt.Println("URN格式:", id.URN())  // urn:uuid:550e8400-e29b-41d4-a716-446655440000

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

	// 解析与验证
	// 安全解析（带错误处理）
	u, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		panic("无效 UUID 格式")
	}
	fmt.Println("有效 UUID 格式:", u)

	// MustParse（测试/确定有效时使用，失败 panic）
	u = uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	fmt.Println("有效 UUID 格式:", u)

	// 转字节切片
	uBytes := u[:] // []byte 长度 16
	fmt.Println("转字节切片：")
	for _, b := range uBytes {
		fmt.Printf("%02x", b)
	}
	fmt.Println()

	// 无连字符格式（需自行处理）
	uNoDash := strings.ReplaceAll(u.String(), "-", "")
	fmt.Println("无连字符格式:", uNoDash)

	// 大写格式
	uUpper := strings.ToUpper(u.String())
	fmt.Println("UUID 大写格式:", uUpper)
}
