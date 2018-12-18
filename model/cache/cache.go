package cache

import "sync"

var pNameCache *PropertyNameCache

// if updating wormhole version will cause property_id and
// property_name mismatched, the whcwallet needs restart
// to clear cache. Otherwise, the app based on this api will
// returned error result.
type PropertyNameCache struct {
	sync.Mutex
	cache map[int64]string
}

func (c *PropertyNameCache) StorePropertyName(pid int64, pname string) {
	c.Lock()
	if _, ok := c.cache[pid]; !ok {
		c.cache[pid] = pname
	}
	c.Unlock()
}

func (c *PropertyNameCache) GetPropertyName(pid int64) (string, bool) {
	c.Lock()
	name, ok := c.cache[pid]
	c.Unlock()

	return name, ok
}

func New() *PropertyNameCache {
	return &PropertyNameCache{
		cache: make(map[int64]string),
	}
}

func GetPNameCache() *PropertyNameCache {
	if pNameCache != nil {
		return pNameCache
	}

	pNameCache = New()
	return pNameCache
}
