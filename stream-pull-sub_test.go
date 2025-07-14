package nats_client

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestStreamPullSub(t *testing.T) {
	// 1. 连接到 NATS 服务器
	nc, err := NewNATSConnect()
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer nc.Close()

	// 2. 获取 JetStream 上下文
	js, err := nc.JetStream(nats.Domain("hub"))
	if err != nil {
		log.Fatalf("获取JetStream失败: %v", err)
	}

	// 5. 订阅并拉取持久化消息
	sub, err := js.PullSubscribe("events.*", "my_consumer")

	if err != nil {
		log.Fatalf("订阅失败: %v", err)
	}

	fmt.Println("拉取消息:")
	for {
		msgs, err := sub.Fetch(3, nats.MaxWait(2*time.Second))
		if err != nil && err != nats.ErrTimeout {
			log.Fatalf("拉取消息失败: %v", err)
		}
		if len(msgs) == 0 {
			fmt.Println("没有更多消息，退出。")
			break
		}
		for _, msg := range msgs {
			fmt.Printf("收到消息:%s %s\n", msg.Subject, string(msg.Data))
			msg.Ack()
		}
	}
}
