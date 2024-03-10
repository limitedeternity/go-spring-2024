//go:build !solution

package lrucache

import "container/list"

type CacheItem[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	items    *list.List
	kv       map[K]*list.Element
	capacity int
}

func (c *LRUCache[K, V]) update(item *list.Element, value V) {
	item.Value.(*CacheItem[K, V]).value = value
	c.items.MoveToFront(item)
}

func (c *LRUCache[K, V]) Init(capacity int) {
	c.items = list.New()
	c.kv = make(map[K]*list.Element, capacity)
	c.capacity = capacity
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	if ptr, ok := c.kv[key]; ok {
		c.update(ptr, ptr.Value.(*CacheItem[K, V]).value)
		return ptr.Value.(*CacheItem[K, V]).value, true
	}

	var placeholder V
	return placeholder, false
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	if ptr, ok := c.kv[key]; ok {
		c.update(ptr, value)
		return
	}

	if c.capacity < 1 {
		return
	}

	if len(c.kv) == c.capacity {
		node := c.items.Back()
		delete(c.kv, node.Value.(*CacheItem[K, V]).key)
		c.items.Remove(node)
	}

	c.kv[key] = c.items.PushFront(&CacheItem[K, V]{key, value})
}

func (c *LRUCache[K, V]) Range(f func(key K, value V) bool) {
	for elem := c.items.Back(); elem != nil; elem = elem.Prev() {
		item := elem.Value.(*CacheItem[K, V])
		if !f(item.key, item.value) {
			break
		}
	}
}

func (c *LRUCache[K, V]) Clear() {
	c.Init(c.capacity)
}

func New(cap int) Cache {
	cache := &LRUCache[int, int]{}
	cache.Init(cap)
	return cache
}
