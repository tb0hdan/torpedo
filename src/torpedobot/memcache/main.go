package memcache

import (
	"sync"
)


type MemCacheType struct {
	cache map[string][]string
	m     sync.RWMutex
}


func (mc *MemCacheType) Get(key string) (value []string, ok bool) {
	mc.m.RLock()
	value, ok = mc.cache[key]
	mc.m.RUnlock()
	return
}


func (mc *MemCacheType) Set(key string, value []string) {
	mc.m.Lock()
	mc.cache[key] = value
	mc.m.Unlock()
}


func (mc *MemCacheType) Len() (cache_size int) {
	cache_size = len(mc.cache)
	return
}


func (mc *MemCacheType) Cache() (cache map[string][]string) {
	cache = mc.cache
	return
}


func (mc *MemCacheType) Delete(key string) {
	delete(mc.cache, key)
}


func New() (memCache *MemCacheType) {
	memCache = &MemCacheType{cache: make(map[string][]string)}
	return
}
