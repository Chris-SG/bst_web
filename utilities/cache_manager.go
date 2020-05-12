package utilities

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	caches map[string]*cache.Cache
)

func CreateCaches() {
	caches["users"] = cache.New(15*time.Minute, 20*time.Minute)
}

func ClearCacheValue(cacheName string, key string) {
	if cacheObject, exists := caches[cacheName]; exists && cacheObject != nil {
		cacheObject.Delete(key)
	}
}

func GetCacheValue(cacheName string, key string) interface{} {
	if cacheObject, exists := caches[cacheName]; exists && cacheObject != nil {
		if value, found := cacheObject.Get(key); found {
			return value
		}
	}
	return nil
}

func SetCacheValue(cacheName string, key string, value interface{}) bool {
	if cacheObject, exists := caches[cacheName]; exists && cacheObject != nil {
		cacheObject.Set(key, value, cache.DefaultExpiration)
		return true
	}

	return false
}