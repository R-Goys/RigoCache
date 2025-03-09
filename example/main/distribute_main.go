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
func startCacheServer(apiAddr string, addr string, addrs []string, g *Rigo.Group) {
	peers := RigoHTTP.NewHttpPool(apiAddr)
	peers.Set(addrs...)
	g.RegisterPeers(peers)
	log.Println("[Cache srv] server start in ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 这里和上面启动的http服务不一样，这里是暴露给用户的。
func startAPIServer(apiAddr string, g *Rigo.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := g.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println(" server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	//var port int
	//var api bool
	//flag.IntVar(&port, "port", 8001, "Geecache server port")
	//flag.BoolVar(&api, "api", false, "Start a api server?")
	//flag.Parse()

	apiAddr := "http://localhost:10004"
	addrMap := map[int]string{
		8001: "http://localhost:10001",
		8002: "http://localhost:10002",
		8003: "http://localhost:10003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	Rigo := CreateGroup("score")
	if true {
		go startAPIServer(apiAddr, Rigo)
	}
	startCacheServer(apiAddr, addrMap[8001], []string{addrMap[8001]}, Rigo)
}
