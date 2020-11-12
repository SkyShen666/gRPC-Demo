package main

import (
	"context"
	"fmt"
	pt "grpcDemo3/myproto"
	"net"

	"google.golang.org/grpc"
)

const (
	post = "127.0.0.1:18881"
)

// 对象要和proto内定义的服务一样
type server struct{}

// SayHello 实现RPC SayHello 接口,对应helloServer.proto文件生成的helloServer.pb.go文件
func (s *server) SayHello(ctx context.Context, in *pt.HelloRequest) (*pt.HelloReply, error) {
	return &pt.HelloReply{Message: "hello " + in.Name}, nil
}

//GetHelloMsg 实现RPC GetHelloMsg 接口
func (s *server) GetHelloMsg(ctx context.Context, in *pt.HelloRequest) (*pt.HelloMessage, error) {
	return &pt.HelloMessage{Msg: "Hi " + in.Name + " this is from server HAHA!"}, nil
}

func main() {
	//监听网络
	ln, err := net.Listen("tcp", post)
	if err != nil {
		fmt.Println("网络异常", err)
	}

	// 创建一个grpc的句柄
	srv := grpc.NewServer()

	//将server结构体注册到 grpc服务中
	pt.RegisterHelloWorldServiceServer(srv, &server{})

	//监听grpc服务
	err = srv.Serve(ln)
	if err != nil {
		fmt.Println("网络启动异常", err)
	}
}
