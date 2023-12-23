package main

import (
	"github.com/iamvineettiwari/go-redis-server-lite/handler"
	"github.com/iamvineettiwari/go-redis-server-lite/server"
)

func main() {
	handler := handler.NewHandler()
	redisServer := server.NewRedisServer(":6379", handler)

	handler.AddHandler("PING", handler.Ping)
	handler.AddHandler("ECHO", handler.Echo)
	handler.AddHandler("SET", handler.Set)
	handler.AddHandler("GET", handler.Get)

	redisServer.Start()
}
