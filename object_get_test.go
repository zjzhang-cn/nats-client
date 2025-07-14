package nats_client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
)

func TestObjectGet(t *testing.T) {
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

	obj, err := js.ObjectStore(ctx, bucket)
	if err != nil {
		log.Fatalf("对象存储打开失败: %v", err)
	}
	log.Printf("对象存储 '%s' 打开成功\n", bucket)

	fs, err := os.OpenFile("out/nats-cli-tmp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}

	obj_result, err := obj.Get(ctx, "nats-cli")
	//obj_info, err := obj.PutFile(ctx, "nats-cli")
	if err != nil {
		log.Fatalf("获取文件失败: %v", err)
	}
	log.Printf("开始获取文件 'nats-cli' \n")
	info, err := obj_result.Info()
	if err != nil {
		log.Fatalf("获取文件信息失败: %v", err)
	}
	log.Printf("File Name: %s", info.Name)
	log.Printf("File Size: %d bytes", info.Size)
	log.Printf("Modified: %v", info.ModTime)
	for k, v := range info.Metadata {
		log.Printf("Metadata %s: %s", k, v)
	}
	progressReader := &ProgressReader{
		Reader: obj_result,
		Total:  int64(info.Size),
		OnProgress: func(readBytes int64, total int64) {
			fmt.Printf("\r进度: %.2f%%", float64(readBytes)/float64(total)*100)
		},
	}
	n, err := io.Copy(fs, progressReader)
	if err != nil {
		log.Fatalf("复制文件内容失败: %v", err)
	}
	log.Printf("文件 '%s' 内容复制成功, 复制了 %d 字节\n", "nats-cli", n)
}
