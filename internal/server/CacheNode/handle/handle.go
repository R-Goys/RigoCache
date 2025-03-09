package handle

import (
	"context"
	Rigo "github.com/R-Goys/RigoCache/internal/core"
	"github.com/R-Goys/RigoCache/internal/rpc"
)

func (c CacheService) Get(_ context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	g := Rigo.GetGroup(request.Group)
	if g == nil {
		g = Rigo.NewGroup("score", nil, 10)
	}
	v, err := g.Get(request.Key)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		Value: v.ByteSlice(),
	}, nil
}
