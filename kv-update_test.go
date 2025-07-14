package nats_client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func TestKvUpdate(t *testing.T) {
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
	// js, err := jetstream.New(nc)
	js, err := jetstream.NewWithDomain(nc, "hub")
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}

	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  bucket,
		TTL:     10 * time.Second,
		History: 10,
	})
	if err != nil {
		log.Printf("Error creating KeyValueBucket: %v", err)
		if kv, err = js.KeyValue(ctx, bucket); err != nil {
			log.Fatalf("Error open KeyValueBucket: %v", err)
		}
	}
	log.Printf("Created KeyValue Bucket: %s", bucket)
	//删除 Bucket
	defer js.DeleteKeyValue(ctx, bucket)
	revision, err := kv.Create(ctx, "key1", []byte("Hello, NATS! "))
	if err != nil {
		log.Printf("Error Creating KeyValue Entry: %v", err)
		if kve, err := kv.Get(ctx, "key1"); err == nil {
			revision = kve.Revision()
		} else {
			log.Fatalf("Error getting KeyValue Entry: %v", err)
		}
		log.Printf("Opened KeyValue Entry: %s, Revision: %d", "key1", revision)
	} else {
		log.Printf("Created KeyValue Entry: %s, Revision: %d", "key1", revision)
	}
	// Publish messages to the subject
	i := 0
	for i < 10 {

		revision, _ = kv.Update(ctx, "key1", []byte(fmt.Sprintf("Hello, NATS! Message number %d", i+1)), revision)
		kvEntry, _ := kv.Get(ctx, "key1")

		if kvEntry != nil {
			log.Printf("Key: %s, Value: %s, Revision: %d", kvEntry.Key(), string(kvEntry.Value()), kvEntry.Revision())
		} else {
			log.Printf("Key not found or expired")
		}
		time.Sleep(time.Second)
		i++
	}
	kv.Purge(ctx, "key1")

}
