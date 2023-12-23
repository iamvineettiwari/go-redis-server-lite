package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/iamvineettiwari/go-redis-server-lite/data"
	"github.com/iamvineettiwari/go-redis-server-lite/handler"
	"github.com/iamvineettiwari/go-redis-server-lite/resp"
)

type RedisServer struct {
	ListenAddr string
	Listener   net.Listener
	store      *data.Store
	connLock   chan struct{}
	handlers   *handler.Handler
}

func NewRedisServer(listenAddr string, handler *handler.Handler) *RedisServer {
	store := data.NewStore()
	handler.ConfigureStore(store)

	return &RedisServer{
		ListenAddr: listenAddr,
		connLock:   make(chan struct{}),
		store:      store,
		handlers:   handler,
	}
}

func (s *RedisServer) Start() error {
	listener, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return err
	}

	s.Listener = listener

	go s.accept()
	<-s.connLock
	return nil
}

func (s *RedisServer) accept() {
	for {
		conn, err := s.Listener.Accept()

		if err != nil {
			fmt.Println("Error while accepting connection : ", err)
			continue
		}

		fmt.Println("Accepted Connection : ", conn.RemoteAddr().String())

		go s.read(conn)
	}
}

func (s *RedisServer) closeConnection(conn net.Conn) {
	fmt.Println("Clossing connection : ", conn.RemoteAddr().String())
	conn.Close()
}

func (s *RedisServer) read(conn net.Conn) {
	defer s.closeConnection(conn)

	buffer := make([]byte, 6048)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Error while reading : ", err)
			continue
		}

		data := buffer[:n]

		request, requestType, err, _ := resp.Deserialize(data)

		s.handleRequest(conn, request, requestType)
	}
}

func (s *RedisServer) handleRequest(conn net.Conn, request any, requestType string) {
	switch requestType {
	case resp.ARRAY:
		items := request.([]resp.ArrayType)

		if len(items) < 1 {
			errorHelper(errors.New("Invalid operations"), conn)
			return
		}

		command := items[0].Value.(string)
		args := []any{}

		for _, argItem := range items[1:] {
			args = append(args, argItem.Value)
		}

		handlerFunc, handlerRegistered := s.handlers.ResolveHandler(strings.ToUpper(command))

		if !handlerRegistered {
			errorHelper(errors.New("Invalid operation"), conn)
			return
		}

		response, err := handlerFunc(args...)

		if err != nil {
			errorHelper(err, conn)
			return
		}

		conn.Write(response)

	default:
		errorHelper(errors.New("Operation not supported"), conn)
	}
}

func errorHelper(err error, conn net.Conn) {
	data, err := resp.Serialize(resp.ERROR, err.Error())

	if err != nil {
		fmt.Println("Error while serializing : ", err)
	}

	conn.Write(data)
}
