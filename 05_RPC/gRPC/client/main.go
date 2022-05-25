package main
import (
	"fmt"
    "google.golang.org/grpc"
	"golang.org/x/net/context"

	"grpc_client/pb"
)

func main(){
    // 连接服务器
    conn, err := grpc.Dial(":8080", grpc.WithInsecure())
    if err != nil {
        fmt.Printf("连接服务端失败: %s", err)
        return
    }
    defer conn.Close()
    // 新建一个客户端
    c := pb.NewHelloClient(conn)

	//创建请求
	//API源自.proto文件自动生成
	req := &pb.HelloRequest{Name: "horika"}
    // 调用服务端函数
    r, err := c.SayHello(context.Background(), req)
    if err != nil {
        fmt.Printf("调用服务端代码失败: %s", err)
        return
    }
    fmt.Printf("调用成功: %s\n", r.Message)
}
