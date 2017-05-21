package multibot

import (
	"torpedobot/memcache"
	"torpedobot/common"
)

func (tb *TorpedoBot) GetCreateCache(name string) (cache *memcache.MemCacheType) {
	value, success := tb.caches[name]
	if !success {
		cache = memcache.New()
		tb.caches[name] = cache
	} else {
		cache = value
	}
	return
}


func (tb *TorpedoBot) GetCachedItem(name string) (item string) {
	cache := *tb.GetCreateCache(name)
	if cache.Len() > 0 {
		tb.logger.Printf("\nUsing cached quote...%v\n", cache.Len())
		key := ""
		for key = range cache.Cache() {
			break
		}
		quote, _ := cache.Get(key)
		cache.Delete(key)
		item = quote[0]
	}
	return
}


func (tb *TorpedoBot) SetCachedItems(name string, items map[int]string) (item string) {
	cache := *tb.GetCreateCache(name)
	for idx := range items {
		message := common.MD5Hash(items[idx])
		_, ok := cache.Get(message)
		if !ok {
			values := make([]string, 1)
			values[0] = items[idx]
			cache.Set(message, values)
		}
	}

	item = items[0]
	message := common.MD5Hash(item)
	cache.Delete(message)
	return
}

