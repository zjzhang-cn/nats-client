package nats_client

import (
	"log"
	"testing"

	"github.com/nats-io/nats.go"
)

func TestStreamQueueSub(t *testing.T) {
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
	js.QueueSubscribe("events.user.*", "my_consumer", func(msg *nats.Msg) {
		log.Printf("收到消息: %s %s\n", msg.Subject, string(msg.Data))
		msg.Ack()
	}, nats.Durable("worker-group"))
	select {}

}
