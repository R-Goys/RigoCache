package consistenthash

import (
	"context"
	"github.com/R-Goys/RigoCache/internal/rpc"
)

type PeerGetter interface {
	Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error)
}

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}
