package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// 服务器用于实现helloworld.GreeterServer。
type server struct {
	// 需要实现UnimplementedGreeterServer，方法，否则会报错
	// cannot use &server literal (type *server) as type helloworld.GreeterServer in argument to helloworld.RegisterGreeterServer:
	// *server does not implement helloworld.GreeterServer (missing helloworld.mustEmbedUnimplementedGreeterServer method)
	pb.UnimplementedGreeterServer
}

// SayHello实现helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	//监听
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// new服务对象
	s := grpc.NewServer()
	// 注册服务
	pb.RegisterGreeterServer(s, &server{})
	// 在gRPC服务器上注册反射服务。
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
