package main

import (
	"fmt"
	"grpc_server/edr_pb"
	"io"
	"net"

	"google.golang.org/grpc"
)

type process_ops_server struct{}

func (server *process_ops_server) UploadProcessInfo(stream edr_pb.RemoteServer_UploadProcessInfoServer) error {
	for {
		tem, err := stream.Recv()
		//客户端调用stream.CloseAndRecv()
		if err == io.EOF {
			fmt.Println("Got EOF")
			return stream.SendAndClose(&edr_pb.ServerAck{Ack: 0})
		}
		//客户端连接断开
		if err != nil {
			fmt.Println("Got error")
			return err
		}
		fmt.Println(tem)
	}
	return nil
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
	edr_pb.RegisterRemoteServerServer(s, &process_ops_server{})
	//reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}
}
