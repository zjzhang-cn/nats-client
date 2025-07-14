package nats_client

import (
	"context"
	"log"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
)

func TestKvWatch(t *testing.T) {
	bucket := "my_bucket"
	nc, err := NewNATSConnect()
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer nc.Close()

	// nc.QueueSubscribe(subject, "queue", func(msg *nats.Msg) {
	// 	log.Printf("Received a message on subject %s: %s", msg.Subject, string(msg.Data))
	// })

	ctx := context.WithoutCancel(context.Background())
	js, err := jetstream.NewWithDomain(nc, "hub")

	// js.Stream("EVENTS", nats.StreamConfig{
	// 	Name:     "EVENTS",
	// 	Subjects: []string{"events.*"},
	//   Storage:  nats.FileStorage, // 持久化到��盘
	// })
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}

	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		log.Fatalf("Error Open KeyValue: %v", err)
	}

	// 启动watcher，监控所有key的变化
	watcher, err := kv.WatchAll(ctx)
	if err != nil {
		log.Fatalf("启动watcher失败: %v", err)
	}
	defer watcher.Stop()

	log.Println("正在监控 KeyValue 变化...")

	// 监控变化并打印
	for entry := range watcher.Updates() {
		if entry == nil {
			continue // watch 启动时会有 nil
		}
		switch entry.Operation() {
		case jetstream.KeyValuePut:
			log.Printf("Key [%s] 被写入/更新, Value: %s, Revision: %d\n", entry.Key(), entry.Value(), entry.Revision())
		case jetstream.KeyValueDelete:
			log.Printf("Key [%s] 被删除, Revision: %d\n", entry.Key(), entry.Revision())
		case jetstream.KeyValuePurge:
			log.Printf("Key [%s] 被彻底清除, Revision: %d\n", entry.Key(), entry.Revision())
		}
		if entry.Key() == "foo" && entry.Value() != nil && string(entry.Value()) == "baz" {
			// 监控到最后一次修改后退出
			break
		}
	}
	log.Println("监控结束。")
}
