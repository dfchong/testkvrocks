package main

import (
	"context"
	"fmt"

	"github.com/redis/rueidis"
)

func main() {
	// 创建 rueidis 客户端配置
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"kvrocks.kvrocks-prod.svc.cluster.local:6666"},
		Password:    "kz501",
		SelectDB:    0,

		// --- 关键修复点 ---
		DisableCache: true, // Kvrocks 不支持客户端缓存，必须设置为 true
		// 建议显式指定协议为 RESP2，因为 Kvrocks 对 RESP3 的支持可能不完整

		// -----------------
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx := context.Background()

	// 写入测试 (Kvrocks 完美支持 SET 命令)
	err = client.Do(ctx, client.B().Set().Key("k3s_test").Value("rueidis_v1.0.73").Build()).Error()
	if err != nil {
		fmt.Printf("写入失败: %v\n", err)
		return
	}

	// 读取测试
	val, err := client.Do(ctx, client.B().Get().Key("k3s_test").Build()).ToString()
	if err != nil {
		fmt.Printf("读取失败: %v\n", err)
		return
	}

	fmt.Printf("读取成功: %s\n", val)
}
