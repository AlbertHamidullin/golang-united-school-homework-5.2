package cache

import (
	"time"
)

type (
	item struct {
		value         string
		isDeadlineSet bool
		deadline      time.Time
	}

	Cache struct {
		m map[string]item
	}
)

func newCacheItem(value string) item {
	return item{
		value:         value,
		isDeadlineSet: false,
		deadline:      time.Time{},
	}
}

func newCacheItemWithDeadline(value string, deadline time.Time) item {
	return item{
		value:         value,
		isDeadlineSet: true,
		deadline:      deadline,
	}
}

func (i item) isExpired() bool {
	return i.isDeadlineSet && time.Until(i.deadline) <= 0
}

func NewCache() Cache {
	return Cache{m: make(map[string]item)}
}

func (cache *Cache) DeleteIfExpired(key string) bool {
	i, ok := cache.m[key]

	if !ok {
		return false
	}

	if i.isExpired() {
		delete(cache.m, key)
		return true
	} else {
		return false
	}
}

func (cache *Cache) Get(key string) (string, bool) {
	i, ok := cache.m[key]

	if !ok {
		return "", false
	}

	// if cache.DeleteIfExpired(key) {
	// 	return "", false
	// }
	if i.isExpired() {
		delete(cache.m, key)
		return "", false
	}

	return i.value, true
}

func (cache *Cache) Put(key, value string) {
	cache.m[key] = newCacheItem(value)
}

func (cache *Cache) Keys() []string {
	r := make([]string, 0, len(cache.m))
	d := make([]string, 0, len(cache.m))

	for k, i := range cache.m {
		if !i.isExpired() {
			r = append(r, k)
		} else {
			d = append(d, k)
		}
	}

	for _, k := range d {
		//cache.DeleteIfExpired(k)
		delete(cache.m, k)
	}

	return r
}

func (cache *Cache) PutTill(key, value string, deadline time.Time) {
	cache.m[key] = newCacheItemWithDeadline(value, deadline)
}
