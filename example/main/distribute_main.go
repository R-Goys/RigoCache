package main

import (
	"fmt"
	Rigo "github.com/R-Goys/RigoCache/core"
	RigoHTTP "github.com/R-Goys/RigoCache/http"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
	"John": "8910",
}

func CreateGroup(name string) *Rigo.Group {
	return Rigo.NewGroup(name, Rigo.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)
}

// 让我思考三分钟，首先是用http服务器的地址初始化节点选择器
// 同时他也能提供http服务，在函数的最后，会启动这个服务，用于通信
// 这里仅仅只有一个Group，但是存在多个节点，相当于这个Group的
// 不同kv键值对被缓存到了不同的节点上，于是就需要我们利用这个选择器来进行选择。
func startCacheServer(addr string, addrs []string, g *Rigo.Group) {
	peers := RigoHTTP.NewHttpPool(addr)
	peers.Set(addrs...)
	g.RegisterPeers(peers)
	log.Println("[Cache srv] server start in ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 这里和上面启动的http服务不一样，这里是暴露给用户的。
func startAPIServer(apiAddr string, g *Rigo.Group) {
	log.Println("[HTTP srv] is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], g.Peers.(*RigoHTTP.HttpPool)))
}

func main() {
	apiAddr := "http://localhost:8080"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := CreateGroup("score")
	for _, v := range addrMap {
		go startCacheServer(v, addrs, gee)
	}
	startAPIServer(apiAddr, gee)
}
