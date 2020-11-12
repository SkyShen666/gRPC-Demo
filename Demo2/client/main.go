package main

import (
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	//这是引用编译好的protobuf
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// 建立到服务器的连接。
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	//延迟关闭连接
	defer conn.Close()

	//调用protobuf的函数创建客户端连接句柄
	c := pb.NewGreeterClient(conn)
	// 联系服务器并打印它的响应。
	  name := "ShenGuanchu"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	//调用protobuf的sayhello函数
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	//打印结果
	log.Printf("Greeting: %s", r.Message)
}
