package data

import (
	"errors"
	"strconv"
	"sync"
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

func (s *Store) Set(key string, value interface{}) {
	if key == "" {
		return
	}

	s.setWithLock(key, value)
}

func (s *Store) Get(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	data, found := s.setLockAndGet(key)

	return data, found
}

func (s *Store) Exists(key string) bool {
	if key == "" {
		return false
	}

	_, found := s.Get(key)

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

	data, exists := s.Get(key)

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

	data, exists := s.Get(key)

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
