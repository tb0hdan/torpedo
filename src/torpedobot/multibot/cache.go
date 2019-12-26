package multibot

import (
	log "github.com/sirupsen/logrus"
	"github.com/tb0hdan/memcache"
	common "github.com/tb0hdan/torpedo_common"
)

func (tb *TorpedoBot) GetCreateCache(name string) (cache *memcache.CacheType) {
	value, success := tb.caches[name]
	if !success {
		cache = memcache.New(log.New())
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
		value, _ := cache.Get(key)
		cache.Delete(key)
		quote := value.([]string)
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
