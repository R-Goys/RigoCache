package handle

import (
	"context"
	Rigo "github.com/R-Goys/RigoCache/internal/core"
	pb "github.com/R-Goys/RigoCache/internal/rpc"
	"github.com/R-Goys/RigoCache/pkg/consistenthash"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultCacheBytes = 10
)

type HttpGetter struct {
	Client pb.RigoCacheClient
}

func (h *HttpGetter) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	res, err := h.Client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

var Getters = &GetterPool{}

type GetterPool struct {
	BasePath    string
	Self        string
	Mu          sync.Mutex
	Peers       *consistenthash.Map
	HttpGetters map[string]*HttpGetter
}

func (p *GetterPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.BasePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	parts := strings.SplitN(r.URL.Path[len(p.BasePath)+1:], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := Rigo.GetGroup(groupName)

	if group == nil {
		group = Rigo.NewGroup(groupName, nil, defaultCacheBytes)
		group.RegisterPeers(p)
	}
	peer, ok := group.Peers.PickPeer(key)
	if !ok {
		http.Error(w, "no such peer: "+key, http.StatusNotFound)
		return
	}
	req := &pb.GetRequest{
		Group: groupName,
		Key:   key,
	}
	res, err := peer.Get(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(res.Value)
}

func (p *GetterPool) PickPeer(key string) (consistenthash.PeerGetter, bool) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if peer := p.Peers.Get(key); peer != "" && peer != p.Self {
		log.Printf("[Server %s] Pick peer %s", p.Self, peer)
		return p.HttpGetters[peer], true
	}
	return nil, false
}
