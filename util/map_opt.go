package util

import (
	"errors"
	"sync"
)

type CacheMap struct {
	lock   sync.RWMutex
	Caches map[int64]Cache
}

type Cache struct {
	Total int64
	Time  int64
}

func (m *CacheMap) New() {
	m.Caches = make(map[int64]Cache)
}

func (m *CacheMap) Add(key int64, value *Cache) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.Caches[key] = *value
}

func (m *CacheMap) Get(key int64) (*Cache, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	cache, ok := m.Caches[key]
	if !ok {
		return nil, errors.New("not exist")
	}

	return &cache, nil
}
