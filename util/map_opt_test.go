package util

import (
	"testing"
	"time"
)

func TestAddCache(t *testing.T) {
	m := new(CacheMap)
	m.New()

	//a := Cache{total:1, time:time.Now().Unix()}

	var i int64
	for i = 0; i < 10; i++ {
		a := Cache{Total: 1, Time: time.Now().Unix()}
		go m.Add(i, &a)
	}

	t.Log(m.Get(9))
}
