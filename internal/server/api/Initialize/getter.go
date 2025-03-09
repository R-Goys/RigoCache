package Initialize

import (
	pb "github.com/R-Goys/RigoCache/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewClient(addr string) pb.RigoCacheClient {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(addr, opt)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil
	}
	Client := pb.NewRigoCacheClient(conn)
	return Client
}
