package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type CacheStore struct {
	items map[string]*CacheItem
	mutex sync.RWMutex
}

var (
	singleton *CacheStore
	once      sync.Once
)

func NewCacheStore() *CacheStore {
	return &CacheStore{
		items: make(map[string]*CacheItem),
	}
}

func GetInstance() *CacheStore {
	once.Do(func() {
		singleton = NewCacheStore()
	})
	return singleton
}

func (cs *CacheStore) Set(key string, value interface{}) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.items[key] = &CacheItem{
		Value:      value,
		Expiration: time.Now().AddDate(0, 0, 1).UnixNano(),
	}
}

func (cs *CacheStore) Get(key string) (interface{}, bool) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	item, found := cs.items[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		delete(cs.items, key)
		return nil, false
	}

	return item.Value, true
}
