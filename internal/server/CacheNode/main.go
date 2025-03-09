package main

import (
	"github.com/R-Goys/RigoCache/internal/rpc"
	"github.com/R-Goys/RigoCache/internal/server/CacheNode/Initialize"
	"github.com/R-Goys/RigoCache/internal/server/CacheNode/handle"
	"github.com/bwmarrin/snowflake"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:10001")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	node, err := snowflake.NewNode(1)
	if err != nil {
		return
	}
	grpcserver := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterRigoCacheServer(grpcserver, &handle.CacheService{})
	Initialize.InitEtcd()
	Initialize.EtcdRegistry.ServiceRegister("cache"+strconv.FormatInt(int64(node.Generate()), 10), "127.0.0.1:10001")
	if err = grpcserver.Serve(listener); err != nil {
		panic(err)
	}
	defer grpcserver.GracefulStop()
	defer listener.Close()
}
