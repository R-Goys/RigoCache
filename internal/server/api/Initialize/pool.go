package Initialize

import (
	"github.com/R-Goys/RigoCache/internal/server/api/handle"
	"github.com/R-Goys/RigoCache/pkg/consistenthash"
	"sync"
)

const (
	defaultbasepath string = "/RigoCache"
	defaultReplicas        = 50
)

func InitPool() {
	handle.Getters = &handle.GetterPool{
		BasePath:    defaultbasepath,
		Self:        "score",
		HttpGetters: make(map[string]*handle.HttpGetter),
		Mu:          sync.Mutex{},
		Peers:       consistenthash.New(defaultReplicas, nil),
	}
}
