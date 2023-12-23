package data

import "sync"

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
	s.setWithLock(key, value)
}

func (s *Store) Get(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	data, found := s.data[key]
	return data, found
}

func (s *Store) setWithLock(key string, value interface{}) {
	if key == "" {
		return
	}

	s.wl.Lock()
	defer s.wl.Unlock()
	s.data[key] = value
}
