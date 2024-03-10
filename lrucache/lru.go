//go:build !solution

package lrucache

import (
	"container/heap"
	"sort"
)

var (
	getNewPrio = func(seed int64) func() int64 {
		return func() int64 {
			seed += 1
			return seed
		}
	}(0)
)

type CacheItem[K comparable, V any] struct {
	key      K
	value    V
	priority int64
	index    int
}

type PriorityQueue[K comparable, V any] []*CacheItem[K, V]

func (pq PriorityQueue[K, V]) Len() int { return len(pq) }

func (pq PriorityQueue[K, V]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue[K, V]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *PriorityQueue[K, V]) Push(x any) {
	n := len(*pq)
	item := x.(*CacheItem[K, V])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[K, V]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue[K, V]) update(item *CacheItem[K, V], value V) {
	item.value = value
	item.priority = getNewPrio()
	heap.Fix(pq, item.index)
}

type LRUCache[K comparable, V any] struct {
	items PriorityQueue[K, V]
	kv    map[K]*CacheItem[K, V]
}

func (c *LRUCache[K, V]) Init(capacity int) {
	c.items = make(PriorityQueue[K, V], 0, capacity)
	c.kv = make(map[K]*CacheItem[K, V], capacity)
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	if ptr, ok := c.kv[key]; ok {
		c.items.update(ptr, ptr.value)
		return ptr.value, true
	}

	var placeholder V
	return placeholder, false
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	if ptr, ok := c.kv[key]; ok {
		c.items.update(ptr, value)
		return
	}

	if cap(c.items) < 1 {
		return
	}

	if len(c.items) == cap(c.items) {
		item := heap.Pop(&c.items).(*CacheItem[K, V])
		delete(c.kv, item.key)
	}

	c.kv[key] = &CacheItem[K, V]{key: key, value: value, priority: getNewPrio()}
	heap.Push(&c.items, c.kv[key])
}

func (c *LRUCache[K, V]) Range(f func(key K, value V) bool) {
	var iter PriorityQueue[K, V] = append([]*CacheItem[K, V](nil), c.items...)
	sort.Sort(iter)

	for _, ptr := range iter {
		if !f(ptr.key, ptr.value) {
			break
		}
	}
}

func (c *LRUCache[K, V]) Clear() {
	c.Init(cap(c.items))
}

func New(cap int) Cache {
	cache := &LRUCache[int, int]{}
	cache.Init(cap)
	return cache
}
