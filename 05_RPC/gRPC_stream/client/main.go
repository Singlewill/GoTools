package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"grpc_client/edr_pb"
)

var RemoteConn *grpc.ClientConn
var RemoteClient edr_pb.RemoteServerClient

var stream edr_pb.RemoteServer_UploadProcessInfoClient
var stream_ok bool

func RemoteServerDisonnect() {
	fmt.Println("Disconnect")
	if stream_ok {
		stream.CloseAndRecv()
	}
	RemoteConn.Close()
}

func RemoteServerConnect(addr string) error {
	RemoteConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("连接服务端失败: %s", err)
		return err
	}
	// 新建一个client对象
	RemoteClient = edr_pb.NewRemoteServerClient(RemoteConn)

	//client推送流
	stream, err = RemoteClient.UploadProcessInfo(context.Background())
	if err != nil {
		stream_ok = false
		return err
	}
	stream_ok = true
	return nil
}

func send_msg(info *edr_pb.ProcessInfo) error {
	var err error
	if !stream_ok {
		stream, err = RemoteClient.UploadProcessInfo(context.Background())
		if err != nil {
			fmt.Println("Create stream failed")
			return err
		}
	}
	stream_ok = true
	err = stream.Send(info)
	if err != nil {
		stream_ok = false
	}
	return err

}

func main() {
	var err error
	for {
		err = RemoteServerConnect("127.0.0.1:8080")
		if err == nil {
			break
		}
		fmt.Println("RemoteServer Connect failed")
		time.Sleep(5 * time.Second)
	}
	fmt.Println("RemoteServer Connect success")
	defer RemoteServerDisonnect()

	i := 0
	for {
		i++
		info := &edr_pb.ProcessInfo{
			Name:     "horika",
			PidCur:   int64(i),
			PidChild: int64(i),
		}
		send_msg(info)
		time.Sleep(1 * time.Second)
		if i > 100 {
			break
		}
	}

}
