package data

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/iamvineettiwari/go-redis-server-lite/data/list"
	"github.com/iamvineettiwari/go-redis-server-lite/resp"
)

type Store struct {
	data map[string]interface{}
	wl   *sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]interface{}),
		wl:   &sync.RWMutex{},
	}
}

func (s *Store) Set(key string, value interface{}, expireCommand string, expireTime int) {
	if key == "" {
		return
	}

	s.setWithLock(key, value)

	if expireCommand != "" && expireTime != 0 {
		timeDuration := s.getTimeDuration(expireCommand, expireTime)

		if timeDuration != time.Duration(0) {
			go func() {
				<-time.After(timeDuration)
				s.deleteWithLock(key)
			}()
		}
	}
}

func (s *Store) Get(key string) (interface{}, bool, error) {
	if key == "" {
		return nil, false, nil
	}

	data, found := s.setLockAndGet(key)

	if !found {
		return data, found, nil
	}

	_, dataIsOfListType := data.(*list.List)

	if dataIsOfListType {
		return nil, found, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return data, found, nil
}

func (s *Store) Exists(key string) bool {
	if key == "" {
		return false
	}

	_, found, _ := s.Get(key)

	return found
}

func (s *Store) Delete(key string) bool {
	if key == "" {
		return false
	}

	exists := s.Exists(key)

	if !exists {
		return false
	}

	s.deleteWithLock(key)
	return true
}

func (s *Store) Incr(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Invalid operation")
	}

	data, exists, err := s.Get(key)

	if err != nil {
		return nil, err
	}

	if !exists {
		data = "1"
		s.setWithLock(key, data)
		return data, nil
	}

	data, isString := data.(string)

	if !isString {
		return nil, errors.New("ERR value is not an integer or out of range")
	}

	value, err := strconv.Atoi(data.(string))

	if err != nil {
		return nil, errors.New("ERR value is not an integer or out of range")
	}

	newValue := strconv.Itoa(value + 1)

	s.setWithLock(key, newValue)
	return newValue, nil
}

func (s *Store) Decr(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Invalid operation")
	}

	data, exists, err := s.Get(key)

	if err != nil {
		return nil, err
	}

	if !exists {
		data = "-1"
		s.setWithLock(key, data)
		return data, nil
	}

	data, isString := data.(string)

	if !isString {
		return nil, errors.New("ERR value is not an integer or out of range")
	}

	value, err := strconv.Atoi(data.(string))

	if err != nil {
		return nil, errors.New("ERR value is not an integer or out of range")
	}

	newValue := strconv.Itoa(value - 1)

	s.setWithLock(key, newValue)
	return newValue, nil
}

func (s *Store) LRange(key string, start, end int) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Invalid Operation")
	}

	data, found := s.setLockAndGet(key)

	if !found {
		return []resp.ArrayType{}, nil
	}

	list, isListType := data.(*list.List)

	if !isListType {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	value := list.GetValues()

	end = min(len(value), end+1)

	return value[start:end], nil
}

func (s *Store) Lpush(key string, val ...interface{}) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Invalid operation")
	}

	if len(val) < 1 {
		return nil, errors.New("ERR wrong number of arguments for 'lpush' command")
	}

	data, found := s.setLockAndGet(key)

	var existList *list.List
	var typeMatch bool

	if !found {
		existList = list.NewList()
	} else {
		existList, typeMatch = data.(*list.List)

		if !typeMatch {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	for _, item := range val {
		existList.InsertLast(item, resp.BULK_STRING)
	}

	s.setWithLock(key, existList)
	return existList.GetValues(), nil
}

func (s *Store) Rpush(key string, val ...interface{}) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Invalid operation")
	}

	if len(val) < 1 {
		return nil, errors.New("ERR wrong number of arguments for 'rpush' command")
	}

	data, found := s.setLockAndGet(key)

	var existList *list.List
	var typeMatch bool

	if !found {
		existList = list.NewList()
	} else {
		existList, typeMatch = data.(*list.List)

		if !typeMatch {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	for _, item := range val {
		existList.InsertFirst(item, resp.BULK_STRING)
	}

	s.setWithLock(key, existList)
	return existList.GetValues(), nil
}

func (s *Store) setLockAndGet(key string) (data interface{}, found bool) {
	s.wl.RLock()
	defer s.wl.RUnlock()
	data, found = s.data[key]
	return
}

func (s *Store) setWithLock(key string, value interface{}) {
	s.wl.Lock()
	defer s.wl.Unlock()
	s.data[key] = value
}

func (s *Store) deleteWithLock(key string) {
	s.wl.Lock()
	defer s.wl.Unlock()
	delete(s.data, key)
}

func (s *Store) getTimeDuration(expireCommand string, timeValue int) time.Duration {
	if expireCommand == "PX" {
		return time.Millisecond * time.Duration(timeValue)
	} else if expireCommand == "EX" {
		return time.Second * time.Duration(timeValue)
	}

	return time.Duration(0)
}
