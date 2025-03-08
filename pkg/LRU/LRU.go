package LRU

import (
	"container/list"
)

type LRUCache struct {
	maxBytes  int64                         //最大容量
	usedBytes int64                         //已经使用的内存
	ll        *list.List                    //维护LRU的链表
	cache     map[string]*list.Element      //字典，方便插入和删除
	onEvicted func(key string, value Value) //回调函数
}
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *LRUCache {
	return &LRUCache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (c *LRUCache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *LRUCache) Put(key string, value Value) {
	if elem, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		elem := c.ll.PushFront(&entry{key, value})
		c.usedBytes += int64(value.Len())
		c.cache[key] = elem
	}
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}

}

func (c *LRUCache) RemoveOldest() {
	elem := c.ll.Back()
	if elem != nil {
		c.ll.Remove(elem)
		kv := elem.Value.(*entry)
		c.usedBytes -= int64(kv.value.Len())
		delete(c.cache, kv.key)
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *LRUCache) Len() int {
	return c.ll.Len()
}
