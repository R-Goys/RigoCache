package Rigo

import (
	"context"
	"github.com/R-Goys/RigoCache/internal/rpc"
	"github.com/R-Goys/RigoCache/pkg/consistenthash"
	"log"
	"sync"
)

// Group 是一切的起点，在进行一系列封装之后
// 我们只需要考虑初始化Group即可
type Group struct {
	name      string
	mainCache cache
	getter    Getter
	Peers     consistenthash.PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

var DataBase = map[string]string{
	"John":     "111",
	"Doe":      "222",
	"Steve":    "333",
	"Sam":      "444",
	"Lance":    "555",
	"lanLance": "666",
	"Noname":   "777",
}

var fn GetterFunc = func(key string) ([]byte, error) {
	log.Println("[SlowDB] Search ", key)
	if v, ok := DataBase[key]; ok {
		return []byte(v), nil
	}
	log.Println("[SlowDB] Not found ", key)
	return nil, nil
}

func NewGroup(name string, getter Getter, cacheBytes int64) *Group {
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}
	if g.getter == nil {
		g.getter = fn
	}
	mu.Lock()
	groups[name] = g
	mu.Unlock()
	return g
}

// RegisterPeers 为group注册注册一个peer
func (g *Group) RegisterPeers(peers consistenthash.PeerPicker) {
	if g.Peers != nil {
		log.Fatalf("[Server %s] Register repeated", g.name)
		return
	}
	g.Peers = peers
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, nil
	}
	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[Cache] Hit")
		return v, nil
	}

	return g.GetLocally(key)
}

func (g *Group) Pick(key string) (ByteView, error) {
	if g.Peers != nil {
		log.Printf("[Server %s] Loading %s", g.name, key)
		//这里先选择一个可供使用的客户端PeerGetter
		if peer, ok := g.Peers.PickPeer(key); ok {
			//利用选择的客户端拿数据
			value, err := g.getFromPeer(peer, key)
			if err == nil {
				log.Println("利用客户端拿到了数据")
				return value, nil
			}
			log.Printf("[Server %s] Failed to get from peers %s\n", g.name, err.Error())
		}
	}
	log.Printf("在本地找")
	return g.GetLocally(key)
}

// GetLocally 从本地拿取数据
func (g *Group) GetLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) getFromPeer(peer consistenthash.PeerGetter, key string) (ByteView, error) {
	req := &pb.GetRequest{
		Group: g.name,
		Key:   key,
	}
	res, err := peer.Get(context.Background(), req)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, err
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Put(key, value)
}
