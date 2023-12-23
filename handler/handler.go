package handler

import (
	"errors"

	"github.com/iamvineettiwari/go-redis-server-lite/data"
	"github.com/iamvineettiwari/go-redis-server-lite/resp"
)

type Handler struct {
	handlers map[string]func(args ...any) ([]byte, error)
	store    *data.Store
}

func NewHandler() *Handler {
	return &Handler{
		handlers: make(map[string]func(args ...any) ([]byte, error)),
	}
}

func (h *Handler) ResolveHandler(path string) (func(args ...any) ([]byte, error), bool) {
	handlerFunc, found := h.handlers[path]
	return handlerFunc, found
}

func (h *Handler) AddHandler(path string, handlerFunc func(args ...any) ([]byte, error)) {
	h.handlers[path] = handlerFunc
}

func (h *Handler) ConfigureStore(store *data.Store) {
	h.store = store
}

func (h *Handler) Ping(args ...any) ([]byte, error) {
	data, err := resp.Serialize(resp.SIMPLE_STRING, "PONG")
	return data, err
}

func (h *Handler) Set(args ...any) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("Invalid operation")
	}

	key := args[0].(string)

	value := ""

	for _, v := range args[1:] {
		value += (v.(string))
	}

	h.store.Set(key, value)

	data, err := resp.Serialize(resp.SIMPLE_STRING, "OK")

	return data, err
}

func (h *Handler) Get(args ...any) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Invalid Operation")
	}

	key := args[0].(string)

	value, _ := h.store.Get(key)

	data, err := resp.Serialize(resp.BULK_STRING, value)

	return data, err
}

func (h *Handler) Echo(args ...any) ([]byte, error) {
	echoString := ""

	for _, item := range args {
		echoString += (item.(string))
	}

	data, err := resp.Serialize(resp.SIMPLE_STRING, echoString)

	return data, err
}