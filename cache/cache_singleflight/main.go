package cache_singleflight

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

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

	// もしなかったらDBにアクセス
	v = GetFromDB(key)
	// キャッシュに保存
	c.Set(key, v)

	return v
}

func GetFromDB(key int) int {
	time.Sleep(time.Second)
	return key * 2
}
