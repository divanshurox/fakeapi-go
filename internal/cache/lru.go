package cache

import "errors"

type CacheEntry struct {
	key   int
	value float32
	prev  *CacheEntry
	next  *CacheEntry
}

type LRUCache struct {
	size  int
	cache map[int]*CacheEntry
	head  *CacheEntry
	tail  *CacheEntry
}

type LRUCacheInterface interface {
	Get(key int) (float32, error)
	Set(key int, value float32)
}

func NewLRUCache(size int) *LRUCache {
	head := &CacheEntry{}
	tail := &CacheEntry{}
	head.next = tail
	tail.prev = head
	return &LRUCache{
		head:  head,
		tail:  tail,
		size:  size,
		cache: make(map[int]*CacheEntry),
	}
}

func (l *LRUCache) Get(key int) (float32, error) {
	if _, ok := l.cache[key]; ok {
		node := l.cache[key]
		l.removeNode(node)
		l.moveToHead(node)
		return node.value, nil
	}
	return -1, errors.New("cache miss")
}

func (l *LRUCache) Set(key int, value float32) {
	if _, ok := l.cache[key]; !ok {
		l.cache[key] = &CacheEntry{key: key, value: value}
		l.moveToHead(l.cache[key])
		if len(l.cache) > l.size {
			tail := l.popTail(l.tail)
			delete(l.cache, tail.key)
		}
	} else {
		l.cache[key].value = value
		l.removeNode(l.cache[key])
		l.moveToHead(l.cache[key])
	}
}

func (l *LRUCache) removeNode(node *CacheEntry) {
	prevNode := node.prev
	nextNode := node.next
	prevNode.next = nextNode
	nextNode.prev = prevNode
}

func (l *LRUCache) moveToHead(node *CacheEntry) {
	node.prev = l.head
	node.next = l.head.next
	l.head.next.prev = node
	l.head.next = node
}

func (l *LRUCache) popTail(node *CacheEntry) *CacheEntry {
	tail := node.prev
	l.removeNode(tail)
	return tail
}
