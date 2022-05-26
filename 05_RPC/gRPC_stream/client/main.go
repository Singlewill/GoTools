package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"grpc_client/edr_pb"
)

func main() {
	// 连接服务器
	//conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("连接服务端失败: %s", err)
		return
	}
	defer conn.Close()
	// 新建一个client对象
	c := edr_pb.NewRemoteServerClient(conn)

	//client推送流
	stream, _ := c.UploadProcessInfo(context.Background())

	i := 0
	for {
		i++
		info := &edr_pb.ProcessInfo{
			Name:     "horika",
			PidCur:   int64(i),
			PidChild: int64(i),
		}
		err := stream.Send(info)
		//服务器连接断开
		if err != nil {
			break
		}
		time.Sleep(1 * time.Second)
		if i > 10 {
			break
		}
	}
	for {
	}

	res, _ := stream.CloseAndRecv()
	fmt.Println(res.Ack)

}
