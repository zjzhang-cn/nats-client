package nats_client

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"golang.org/x/net/proxy"
)

func NewNATSConnect() (*nats.Conn, error) {
	opts := []nats.Option{nats.Name("nats-client")}
	natsUrl := os.Getenv("NATS_URL")
	opts = append(opts, nats.SetCustomDialer(proxy.FromEnvironment()))
	opts = append(opts, nats.ConnectHandler(func(c *nats.Conn) {
		log.Printf(": [NATS] Connected\n")
	}))
	nc, err := nats.Connect(natsUrl, opts...)
	return nc, err
}
