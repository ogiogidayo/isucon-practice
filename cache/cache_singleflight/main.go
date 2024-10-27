package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"math/rand"
	"sync"
	"time"
)

var group singleflight.Group

func main() {
	rand.Seed(time.Now().UnixNano())

	mCache := NewCache()

	for i := 0; i < 20; i++ {
		key := rand.Intn(5) + 1
		value := mCache.Get(key)
		log.Println(value)
	}
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

	vv, err, _ := group.Do(fmt.Sprintf("cacheGet_%d", key), func() (interface{}, error) {
		value := GetFromDB(key)
		c.Set(key, value)
		return value, nil
	})

	if err != nil {
		panic(err)
	}

	return vv.(int)

	return v
}

func GetFromDB(key int) int {
	time.Sleep(time.Second)
	return key * 2
}
