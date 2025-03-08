package test

import (
	"fmt"
	"github.com/R-Goys/RigoCache/pkg/LRU"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func Test_lru(T *testing.T) {
	lru := LRU.New(12, func(key string, value LRU.Value) {
		fmt.Println(key, value, "被删除啦！！！！")
	})
	lru.Put("key1", String("value1"))
	lru.Put("key2", String("value2"))
	lru.Put("key3", String("value3"))
	lru.Put("key4", String("value4"))
	if v, ok := lru.Get("key1"); ok {
		fmt.Println(v)
		return
	}
	if v, ok := lru.Get("key2"); ok {
		fmt.Println(v)
	}
	if v, ok := lru.Get("key3"); ok {
		fmt.Println(v)
	}
	if v, ok := lru.Get("key4"); ok {
		fmt.Println(v)
	}
}
