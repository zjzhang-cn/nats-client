package nats_client

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	"golang.org/x/net/proxy"
)

// TestNATSConnection 测试 NATS 连接功能
func TestNATSConnection(t *testing.T) {
	// 准备测试环境
	nc, err := NewNATSConnect()
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer nc.Close()

	// 验证连接状态
	if !nc.IsConnected() {
		t.Fatal("NATS连接状态异常")
	}
}

// TestJetStreamContext 测试 JetStream 上下文获取
func TestJetStreamContext(t *testing.T) {
	nc := setupTestConnection(t)
	defer nc.Drain()

	// 获取 JetStream 上下文
	js, err := nc.JetStream(nats.Domain("hub"))
	if err != nil {
		t.Fatalf("获取JetStream失败: %v", err)
	}

	// 验证 JetStream 上下文
	if js == nil {
		t.Fatal("JetStream上下文为空")
	}

	t.Log("JetStream上下文获取成功")
}

// TestQueueSubscribe 测试队列订阅功能
func TestQueueSubscribe(t *testing.T) {
	nc := setupTestConnection(t)
	defer nc.Drain()

	js, err := nc.JetStream(nats.Domain("hub"))
	if err != nil {
		t.Fatalf("获取JetStream失败: %v", err)
	}

	// 创建测试用的消息接收通道
	msgReceived := make(chan *nats.Msg, 1)

	// 设置队列订阅
	sub, err := js.QueueSubscribe("events.test", "test_consumer", func(msg *nats.Msg) {
		log.Printf("[TEST] 收到消息: %s %s", msg.Subject, string(msg.Data))
		msgReceived <- msg
		msg.Ack()
	}, nats.Durable("test-worker-group"))

	if err != nil {
		t.Fatalf("队列订阅失败: %v", err)
	}
	defer sub.Unsubscribe()

	// 发布测试消息
	testData := []byte("test message data")
	_, err = js.Publish("events.test", testData)
	if err != nil {
		t.Fatalf("发布测试消息失败: %v", err)
	}

	// 等待接收消息
	select {
	case msg := <-msgReceived:
		if string(msg.Data) != string(testData) {
			t.Errorf("接收到的消息内容不匹配: 期望=%s, 实际=%s", string(testData), string(msg.Data))
		}
		if msg.Subject != "events.test" {
			t.Errorf("接收到的消息主题不匹配: 期望=events.test, 实际=%s", msg.Subject)
		}
		t.Log("队列订阅测试成功")
	case <-time.After(5 * time.Second):
		t.Fatal("等待消息超时")
	}
}

// TestMultipleMessages 测试处理多个消息
func TestMultipleMessages(t *testing.T) {
	nc := setupTestConnection(t)
	defer nc.Drain()

	js, err := nc.JetStream(nats.Domain("hub"))
	if err != nil {
		t.Fatalf("获取JetStream失败: %v", err)
	}

	const messageCount = 5
	msgReceived := make(chan *nats.Msg, messageCount)

	// 设置订阅
	sub, err := js.QueueSubscribe("events.multi", "multi_consumer", func(msg *nats.Msg) {
		log.Printf("[TEST] 收到消息 #%s: %s", msg.Subject, string(msg.Data))
		msgReceived <- msg
		msg.Ack()
	}, nats.Durable("test-multi-group"))

	if err != nil {
		t.Fatalf("队列订阅失败: %v", err)
	}
	defer sub.Unsubscribe()

	// 发布多个测试消息
	for i := 0; i < messageCount; i++ {
		testData := []byte("test message " + string(rune('0'+i)))
		_, err = js.Publish("events.multi", testData)
		if err != nil {
			t.Fatalf("发布测试消息 %d 失败: %v", i, err)
		}
	}

	// 验证所有消息都被接收
	receivedCount := 0
	timeout := time.After(10 * time.Second)

	for receivedCount < messageCount {
		select {
		case msg := <-msgReceived:
			receivedCount++
			t.Logf("接收到消息 %d/%d: %s", receivedCount, messageCount, string(msg.Data))
		case <-timeout:
			t.Fatalf("接收消息超时，只收到 %d/%d 条消息", receivedCount, messageCount)
		}
	}

	t.Logf("成功接收所有 %d 条消息", messageCount)
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 测试错误的连接URL
	_, err := nats.Connect("nats://invalid-url:4222")
	if err == nil {
		t.Error("期望连接失败，但连接成功了")
	}
	t.Logf("正确处理了连接错误: %v", err)

	// 测试有效连接但无效的JetStream域
	nc := setupTestConnection(t)
	defer nc.Drain()

	_, err = nc.JetStream(nats.Domain("invalid-domain"))
	if err != nil {
		t.Logf("正确处理了JetStream域错误: %v", err)
	}
}

// BenchmarkQueueSubscribe 性能基准测试
func BenchmarkQueueSubscribe(b *testing.B) {
	nc := setupBenchConnection(b)
	defer nc.Drain()

	js, err := nc.JetStream(nats.Domain("hub"))
	if err != nil {
		b.Fatalf("获取JetStream失败: %v", err)
	}

	msgReceived := make(chan *nats.Msg, b.N)

	sub, err := js.QueueSubscribe("events.bench", "bench_consumer", func(msg *nats.Msg) {
		msgReceived <- msg
		msg.Ack()
	}, nats.Durable("bench-group"))

	if err != nil {
		b.Fatalf("队列订阅失败: %v", err)
	}
	defer sub.Unsubscribe()

	testData := []byte("benchmark test data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = js.Publish("events.bench", testData)
		if err != nil {
			b.Fatalf("发布消息失败: %v", err)
		}
	}

	// 等待所有消息被处理
	for i := 0; i < b.N; i++ {
		select {
		case <-msgReceived:
		case <-time.After(time.Second):
			b.Fatalf("基准测试超时，处理了 %d/%d 条消息", i, b.N)
		}
	}
}

// setupTestConnection 设置测试连接的辅助函数
func setupTestConnection(t *testing.T) *nats.Conn {
	opts := []nats.Option{nats.Name("nats-client-test")}

	opts = append(opts, nats.SetCustomDialer(proxy.FromEnvironment()))

	natsUrl := os.Getenv("NATS_TEST_URL")

	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		t.Fatalf("NATS连接失败: %v", err)
	}

	return nc
}

// setupBenchConnection 设置基准测试连接的辅助函数
func setupBenchConnection(b *testing.B) *nats.Conn {
	opts := []nats.Option{nats.Name("nats-client-bench")}
	opts = append(opts, nats.SetCustomDialer(proxy.FromEnvironment()))

	natsUrl := os.Getenv("NATS_TEST_URL")

	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		b.Fatalf("NATS连接失败: %v", err)
	}

	return nc
}

// TestMain 测试入口点，可以进行全局的测试设置和清理
func TestMain(m *testing.M) {
	// 测试前的全局设置
	log.Println("开始运行 NATS 客户端测试...")

	// 运行测试
	code := m.Run()

	// 测试后的全局清理
	log.Println("NATS 客户端测试完成")

	os.Exit(code)
}
