package nats_client

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
)

func TestObjectPut(t *testing.T) {
	bucket := "my_object_store"
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
	if err != nil {
		log.Fatalf("创建 JetStream 客户端失败: %v", err)
	}

	obj, err := js.CreateOrUpdateObjectStore(ctx, jetstream.ObjectStoreConfig{
		Bucket:   bucket,
		Replicas: 2,
		Storage:  jetstream.FileStorage, // 持久化到磁盘
	})
	if err != nil {
		log.Fatalf("创建或更新对象存储失败: %v", err)
	}
	log.Printf("对象存储 '%s' 创建成功\n", bucket)
	fs, err := os.Open("nats-cli")
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}
	stat, _ := fs.Stat()
	log.Printf("打开文件成功: [%s]:%d", stat.Name(), stat.Size())

	// Calculate SHA-256 hash of the file
	hash := sha256.New()
	if _, err := io.Copy(hash, fs); err != nil {
		log.Fatalf("计算文件哈希值失败: %v", err)
	}
	sha256sum := fmt.Sprintf("%x", hash.Sum(nil))
	log.Printf("文件的 SHA-256: %s", sha256sum)

	// Reset file pointer to the beginning for the subsequent read
	if _, err := fs.Seek(0, 0); err != nil {
		log.Fatalf("重置文件指针失败: %v", err)
	}

	progressReader := &ProgressReader{
		Reader: fs,
		Total:  stat.Size(),
		OnProgress: func(readBytes int64, total int64) {
			fmt.Printf("\r进度: %.2f%%", float64(readBytes)/float64(total)*100)
		},
	}
	obj_info, err := obj.Put(ctx,
		jetstream.ObjectMeta{
			Name:        "nats-cli",
			Description: "NATS CLI Tool",
			Metadata: map[string]string{
				"version": "1.0",
				"author":  "NATS Team",
			},
		},
		progressReader,
	)
	//obj_info, err := obj.PutFile(ctx, "nats-cli")
	if err != nil {
		log.Fatalf("上传文件失败: %v", err)
	}
	log.Printf("文件上传成功: %s, size: %d bytes\n", obj_info.Name, obj_info.Size)
}
