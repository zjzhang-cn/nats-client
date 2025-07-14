package nats_client

import (
	"fmt"
	"log"
	"testing"

	"github.com/nats-io/nats.go"
)

func TestStreamPub(t *testing.T) {
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

	// 3. 创建/更新一个持久化流
	streamName := "EVENTS"
	subject := "events.>"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.FileStorage, // 持久化到磁盘
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatalf("添加流失败: %v", err)
	}

	// 4. 发布持久化消息
	for i := 1; i <= 3; i++ {
		ack, err := js.Publish("events.user.1", []byte(fmt.Sprintf(`{"msg":"消息 #%d"}`, i)))
		if err != nil {
			log.Fatalf("发布消息失败: %v", err)
		}
		fmt.Printf("已发布消息: %s, 序号: %d\n", ack.Stream, ack.Sequence)
	}
	for i := 1; i <= 3; i++ {
		ack, err := js.Publish("events.user.2", []byte(fmt.Sprintf(`{"msg":"消息 #%d"}`, i)))
		if err != nil {
			log.Fatalf("发布消息失败: %v", err)
		}
		fmt.Printf("已发布消息: %s, 序号: %d\n", ack.Stream, ack.Sequence)
	}
	for i := 1; i <= 3; i++ {
		ack, err := js.Publish("events.admin.2", []byte(fmt.Sprintf(`{"msg":"消息 #%d"}`, i)))
		if err != nil {
			log.Fatalf("发布消息失败: %v", err)
		}
		fmt.Printf("已发布消息: %s, 序号: %d\n", ack.Stream, ack.Sequence)
	}
}
