package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

	if len(args) > 2 && len(args) < 4 {
		return nil, errors.New("Invalid operation")
	}

	key := args[0].(string)
	value := args[1].(string)

	var expireCommand string
	var expireTime int

	if len(args) > 2 {
		com := strings.ToUpper(args[2].(string))

		if com != "PX" && com != "EX" {
			return nil, errors.New("Invalid operation")
		}

		val, err := strconv.Atoi(args[3].(string))

		if err != nil {
			return nil, errors.New("ERR invalid expire time")
		}

		expireTime = val
		expireCommand = com
	}

	h.store.Set(key, value, expireCommand, expireTime)

	data, err := resp.Serialize(resp.SIMPLE_STRING, "OK")

	return data, err
}

func (h *Handler) Get(args ...any) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Invalid Operation")
	}

	key := args[0].(string)

	value, _, err := h.store.Get(key)

	if err != nil {
		return nil, err
	}

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

func (h *Handler) Exists(args ...any) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Invalid Operation")
	}

	totalFound := 0

	for _, key := range args {
		if h.store.Exists(key.(string)) {
			totalFound++
		}
	}

	data, err := resp.Serialize(resp.INTEGER, totalFound)

	return data, err
}

func (h *Handler) Delete(args ...any) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Invalid operation")
	}

	totalDeleted := 0

	for _, key := range args {
		if h.store.Delete(key.(string)) {
			totalDeleted++
		}
	}

	data, err := resp.Serialize(resp.INTEGER, totalDeleted)
	return data, err
}

func (h *Handler) Incr(args ...any) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Invalid operation")
	}

	increment, err := h.store.Incr(args[0].(string))

	if err != nil {
		return nil, err
	}

	data, err := resp.Serialize(resp.INTEGER, increment)

	return data, err
}

func (h *Handler) Decr(args ...any) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Invalid operation")
	}

	increment, err := h.store.Decr(args[0].(string))

	if err != nil {
		return nil, err
	}

	data, err := resp.Serialize(resp.INTEGER, increment)

	return data, err
}

func (h *Handler) Lpush(args ...any) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("Invalid operation")
	}

	key := args[0].(string)
	val := args[1:]

	items, err := h.store.Lpush(key, val...)

	if err != nil {
		return nil, err
	}

	data, err := resp.Serialize(resp.ARRAY, items)

	return data, err
}

func (h *Handler) Rpush(args ...any) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("Invalid operation")
	}

	key := args[0].(string)
	val := args[1:]

	items, err := h.store.Rpush(key, val...)

	if err != nil {
		return nil, err
	}

	data, err := resp.Serialize(resp.ARRAY, items)

	return data, err
}

func (h *Handler) LRange(args ...any) ([]byte, error) {
	if len(args) < 3 {
		return nil, errors.New(fmt.Sprintf("wrong number of arguments (given %d, expected 3)", len(args)))
	}

	key := args[0].(string)
	start, err := strconv.Atoi(args[1].(string))

	if err != nil {
		return nil, errors.New("ERR value is not an integer or out of range")
	}

	end, err := strconv.Atoi(args[2].(string))

	if err != nil {
		return nil, errors.New("ERR value is not an integer or out of range")
	}

	items, err := h.store.LRange(key, start, end)

	if err != nil {
		return nil, err
	}

	data, err := resp.Serialize(resp.ARRAY, items)

	return data, err
}
