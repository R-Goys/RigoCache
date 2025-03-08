package test

import (
	"fmt"
	Rigo "github.com/R-Goys/RigoCache/core"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
	"John": "8910",
}

func Test_group(T *testing.T) {
	var getter Rigo.GetterFunc = func(key string) ([]byte, error) {
		fmt.Println("SlowDB Search ", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, nil
	}
	g := Rigo.NewGroup("hajimi", getter, 20)
	for k, v := range db {
		if _, err := g.Get(k); err == nil {
			fmt.Println("cache miss", k, v)
		}
	}
	for k, v := range db {
		if _, err := g.Get(k); err == nil {
			fmt.Println("cache miss", k, v)
		}
	}
}
