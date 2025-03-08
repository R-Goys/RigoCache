package Rigo

import (
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

func NewGroup(name string, getter Getter, cacheBytes int64) *Group {
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
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

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	if g.Peers != nil {
		//这里先选择一个可供使用的客户端PeerGetter
		if peer, ok := g.Peers.PickPeer(key); ok {
			//利用选择的客户端拿数据
			if value, err := g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Fatalf("[Server %s] Failed to get from peers", g.name)
		}
	}
	return g.getLocally(key)
}

// 从本地拿取数据
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) getFromPeer(peer consistenthash.PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: cloneBytes(bytes)}, err
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Put(key, value)
}
