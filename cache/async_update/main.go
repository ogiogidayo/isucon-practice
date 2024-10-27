package main

import (
	"log"
	"sync"
	"time"
)

const defaultValue = 100

func main() {
	mCache := NewCache()
	log.Println(mCache.Get(3))
	time.Sleep(2 * time.Second)
	log.Println(mCache.Get(3))
}

type Cache struct {
	mu    sync.Mutex
	items map[int]int
}

func NewCache() *Cache {
	m := make(map[int]int)
	c := &Cache{items: m}

	return c
}

func (c *Cache) Set(key int, value int) {
	c.mu.Lock()
	c.items[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key int) int {
	c.mu.Lock()
	v, ok := c.items[key]
	c.mu.Unlock()

	if ok {
		return v
	}

	go func() {
		// 非同期にキャッシュを更新
		v := GetFromDB(key)
		c.Set(key, v)
	}()

	return defaultValue
}

func GetFromDB(key int) int {
	time.Sleep(time.Second)
	return key * 2
}
