package RigoHTTP

import (
	"fmt"
	Rigo "github.com/R-Goys/RigoCache/core"
	"github.com/R-Goys/RigoCache/pkg/consistenthash"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultbasepath string = "/RigoCache/"
	defaultReplicas        = 50
)

// HttpPool 实现了节点选择的接口，作为大脑，去选择，可以去选择节点
type HttpPool struct {
	basePath    string
	self        string
	mu          sync.Mutex
	peers       *consistenthash.Map
	HttpGetters map[string]*HttpGetter
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		basePath:    defaultbasepath,
		self:        self,
		HttpGetters: make(map[string]*HttpGetter),
	}
}
func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 这里也相当于是实现了接口，url格式固定为:"/{basePath}/{GroupName}/{Key}"
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log(r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	group := Rigo.GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

// Set 设置节点
func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.peers == nil {
		p.peers = consistenthash.New(defaultReplicas, nil)
	}
	p.peers.Add(peers...)
	if len(p.HttpGetters) == 0 {
		p.HttpGetters = make(map[string]*HttpGetter)
	}
	for _, peer := range peers {
		p.HttpGetters[peer] = &HttpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer 选择一个节点，返回一个实现了能够拿去键值对的接口，从而使的调用方能够拿取数据。
func (p *HttpPool) PickPeer(key string) (consistenthash.PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		log.Printf("[Server %s] Pick peer %s", p.self, peer)
		return p.HttpGetters[peer], true
	}
	return nil, false
}
