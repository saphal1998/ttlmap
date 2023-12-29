package ttlcache

import (
	"log"
	"time"
)

type TTLCache interface {
	Add(key string, value any, ttl uint64) // Add any items to the cache with a time to live, overrides the key if it exists previously, so we do not need a different update method
	Get(key string) any                    // Get any value for a key
	Delete(key string)                     // Delete the key from the cache
}

type addRequestArgs struct {
	item any
	key  string
	ttl  uint64
}

type element struct {
	item     any
	refcount uint64
}

type Cache struct {
	store                map[string]element
	addRequestChannel    chan addRequestArgs
	getRequestChannel    chan string
	getResponseChannel   chan any
	deleteRequestChannel chan string
	detonateChannel      chan string
}

func (c *Cache) run() {
	for {
		select {
		case addRequest := <-c.addRequestChannel:
			c.store[addRequest.key] = element{
				item:     addRequest.item,
				refcount: 1,
			}
			go c.detonate(addRequest.key, addRequest.ttl)
		case key := <-c.getRequestChannel:
			val, ok := c.store[key]
			if !ok || val.refcount <= 0 {
				c.getResponseChannel <- nil
				continue
			}
			c.getResponseChannel <- val.item
		case key := <-c.deleteRequestChannel:
			val, ok := c.store[key]
			if !ok {
				// Don't do anything if trying to delete a key which does not exist in the cache
				continue
			}
			// Since the element already is in the store, it has a ref out there, once that ref ticks, we can actually remove it from the cache
			val.refcount = 0
		case key := <-c.detonateChannel:
			val, ok := c.store[key]
			if !ok {
				log.Fatalf("Trying to detonate %v but this isn't even present in the cache", key)
				continue
			}
			val.refcount -= 1
			if val.refcount <= 0 {
				delete(c.store, key)
			}
		}
	}
}

func (c *Cache) detonate(key string, ttl uint64) {
	time.Sleep(time.Duration(ttl))
	c.detonateChannel <- key
}

func NewCache() Cache {
	cache := Cache{
		store:                make(map[string]element),
		addRequestChannel:    make(chan addRequestArgs),
		getRequestChannel:    make(chan string),
		getResponseChannel:   make(chan any),
		deleteRequestChannel: make(chan string),
		detonateChannel:      make(chan string),
	}

	go cache.run()

	return cache
}
