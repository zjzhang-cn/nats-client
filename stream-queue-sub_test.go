package nats_client

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestStreamQueueSub(t *testing.T) {
	// 定义测试用的主题、队列和持久化名称
	const (
		testSubject = "events.user.test"
		testQueue   = "test-queue"
		testDurable = "test-durable"
		testDomain  = "hub"
	)

	// 1. 连接到 NATS 服务器
	nc, err := NewNATSConnect()
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer nc.Close()

	// 2. 获取 JetStream 上下文
	js, err := nc.JetStream(nats.Domain(testDomain))
	if err != nil {
		t.Fatalf("获取JetStream失败: %v", err)
	}

	// 3. 创建一个 channel 来接收消息和错误
	msgCh := make(chan *nats.Msg, 1)
	errCh := make(chan error, 1)

	// 4. 队列订阅
	sub, err := js.QueueSubscribe(testSubject, testQueue, func(msg *nats.Msg) {
		t.Logf("收到消息: %s %s", msg.Subject, string(msg.Data))
		msg.Ack()
		msgCh <- msg
	}, nats.Durable(testDurable))
	if err != nil {
		t.Fatalf("队列订阅失败: %v", err)
	}
	defer sub.Unsubscribe()

	// 5. 发布一条测试消息
	const testMessage = "hello world"
	_, err = js.Publish(testSubject, []byte(testMessage))
	if err != nil {
		t.Fatalf("发布消息失败: %v", err)
	}

	// 6. 等待消息或超时
	select {
	case msg := <-msgCh:
		if string(msg.Data) != testMessage {
			t.Errorf("收到的消息不匹配: got %q, want %q", string(msg.Data), testMessage)
		}
	case err := <-errCh:
		t.Fatalf("订阅出错: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("等待消息超时")
	}
}
