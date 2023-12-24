package main

import (
	"github.com/iamvineettiwari/go-redis-server-lite/handler"
	"github.com/iamvineettiwari/go-redis-server-lite/server"
)

func main() {
	handlerInstance := handler.NewHandler()
	redisServer := server.NewRedisServer(":6379", handlerInstance)

	handlerInstance.AddHandler(handler.PING, handlerInstance.Ping)
	handlerInstance.AddHandler(handler.ECHO, handlerInstance.Echo)
	handlerInstance.AddHandler(handler.SET, handlerInstance.Set)
	handlerInstance.AddHandler(handler.GET, handlerInstance.Get)
	handlerInstance.AddHandler(handler.EXISTS, handlerInstance.Exists)
	handlerInstance.AddHandler(handler.DEL, handlerInstance.Delete)
	handlerInstance.AddHandler(handler.INCR, handlerInstance.Incr)
	handlerInstance.AddHandler(handler.DECR, handlerInstance.Decr)
	handlerInstance.AddHandler(handler.LRANGE, handlerInstance.LRange)
	handlerInstance.AddHandler(handler.LPUSH, handlerInstance.Lpush)
	handlerInstance.AddHandler(handler.RPUSH, handlerInstance.Rpush)

	redisServer.Start()
}
