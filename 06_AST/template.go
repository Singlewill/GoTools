package remote_server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var RemoteConn *grpc.ClientConn
var RemoteClient RemoteServerClient

var streamProcessFork RemoteServer_UploadProcessForkClient
var streamProcessFork_ok bool

func RemoteServerConnect(addr string) error {
	RemoteConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	RemoteClient = NewRemoteServerClient(RemoteConn)

}

func RemoteServerDisconnect() {
	RemoteConn.Close()
}

//Fork信息上传
//这里非阻塞，不判断返回值
//如果发送失败，将重新建立一次流连接,同样不管建立成功与否
func UploadProcessFork(info *ProcessForkInfo) {
	var err error
	if !streamProcessFork_ok {
		streamProcessFork, err = RemoteClient.UploadProcessFork(context.Background())
		if err != nil {
			return
		}
	}
	streamProcessFork_ok = true
	err = streamProcessFork.Send(info)
	if err != nil {
		streamProcessFork_ok = false
	}
}