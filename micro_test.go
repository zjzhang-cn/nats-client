package nats_client

import (
	"context"
	"fmt"
	"log"
	"testing"

	services "github.com/nats-io/nats.go/micro"
	//nolint
)

func TestMicroSV(t *testing.T) {
	ctx := context.Background()

	nc, err := NewNATSConnect()
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer nc.Close()

	fmt.Println("Starting echo service")

	sv, err := services.AddService(nc, services.Config{
		Name:        "UserSV",
		Version:     "1.0.0",
		Description: "User management service",
		Metadata:    map[string]string{"MCP": "User management"},
		// Endpoints: []services.EndpointConfig{
		// base handler
		// Endpoint: &services.EndpointConfig{
		// 	Subject: "User.echo",
		// 	Handler: services.HandlerFunc(func(req services.Request) {
		// 		log.Printf("Received request: %s\n", string(req.Data()))
		// 		req.Respond(req.Data())
		// 	}),
		// },
	})
	sv.AddEndpoint("login",
		services.HandlerFunc(func(req services.Request) {
			log.Printf("Received request: %s\n", string(req.Subject()))
			log.Printf("Received request: %s\n", string(req.Data()))
			//req.Respond(req.Data())
			req.Error("400", "Bad Request", []byte("Invalid login credentials"))
		}),
		services.WithEndpointSubject("User.login"),
		services.WithEndpointMetadata(map[string]string{
			"description": "签入",
			"MCP":         "User management",
		}),
	)
	sv.AddEndpoint("logout", services.HandlerFunc(func(req services.Request) {
		log.Printf("Received request: %s\n", string(req.Subject()))
		log.Printf("Received request: %s\n", string(req.Data()))
		req.Respond(req.Data())
	}), services.WithEndpointSubject("User.logout"),
		services.WithEndpointMetadata(map[string]string{
			"description": "签出",
			"MCP":         "User management",
		}))
	sv.AddEndpoint("check", services.HandlerFunc(func(req services.Request) {
		log.Printf("Received request: %s\n", string(req.Subject()))
		log.Printf("Received request: %s\n", string(req.Data()))
		req.Respond(req.Data())
	}), services.WithEndpointSubject("User.check"),
		services.WithEndpointMetadata(map[string]string{
			"description": "检测",
			"MCP":         "User management",
		}))
	sv.AddEndpoint("create", services.HandlerFunc(func(req services.Request) {
		log.Printf("Received request: %s\n", string(req.Subject()))
		log.Printf("Received request: %s\n", string(req.Data()))
		req.Respond(req.Data())
	}), services.WithEndpointSubject("User.create"),
		services.WithEndpointMetadata(map[string]string{
			"description": "创建",
			"MCP":         "User management",
		}))
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
}
