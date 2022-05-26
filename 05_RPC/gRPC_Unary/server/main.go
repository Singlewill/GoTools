package main

import (
	"fmt"
	"grpc_server/pb"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (this *server) SayHello(ctx context.Context, in *pb.HelloRequest) (out *pb.HelloResponse, err error) {
	return &pb.HelloResponse{Message: in.Name + " : hello"}, nil
}

func main() {
	// 监听本地端口
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("监听端口失败: %s", err)
		return
	}
	// 创建gRPC服务器
	s := grpc.NewServer()
	// 注册服务
	//API源自.proto文件自动生成
	pb.RegisterHelloServer(s, &server{})
	//reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}
}
