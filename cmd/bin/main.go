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
	handler.AddHandler("EXISTS", handler.Exists)
	handler.AddHandler("DEL", handler.Delete)
	handler.AddHandler("INCR", handler.Incr)
	handler.AddHandler("DECR", handler.Decr)
	handler.AddHandler("LRANGE", handler.LRange)
	handler.AddHandler("LPUSH", handler.Lpush)
	handler.AddHandler("RPUSH", handler.Rpush)

	redisServer.Start()
}
