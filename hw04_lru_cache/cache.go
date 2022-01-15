package hw04lrucache

import (
	"sync"
)

// Key is the type of key in the cache.
type Key string

// Cache is a cache interface.
type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// lruCache is a cache structure.
type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// Set adds an element to the cache.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Already existing key
	if listItem, ok := c.items[key]; ok {
		listItem.Value = cacheItem{key: key, value: value}
		c.queue.MoveToFront(listItem)
		return true
	}

	// Drop last element
	if c.queue.Len() >= c.capacity {
		lastListItem := c.queue.Back()
		c.queue.Remove(lastListItem)
		delete(c.items, lastListItem.Value.(cacheItem).key)
	}
	// New key
	listItem := c.queue.PushFront(
		cacheItem{
			key:   key,
			value: value,
		},
	)
	c.items[key] = listItem
	return false
}

// Get returns an element from the cache.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if listItem, ok := c.items[key]; ok {
		c.queue.MoveToFront(listItem)
		return listItem.Value.(cacheItem).value, true
	}
	return nil, false
}

// Clear clears the cache.
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[Key]*ListItem)
	c.queue = NewList()
}

// cacheItem is a structure for storing an item in the cache.
type cacheItem struct {
	key   Key
	value interface{}
}

// NewCache creates cache.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
