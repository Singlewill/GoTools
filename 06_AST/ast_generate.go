package remote_server
import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
var RemoteConn *grpc.ClientConn
var RemoteClient RemoteServerClient
var StreamProcessFork RemoteServer_UploadProcessForkClient
var StreamProcessFork_ok bool
var StreamProcessExit RemoteServer_UploadProcessExitClient
var StreamProcessExit_ok bool
var StreamExecve RemoteServer_UploadExecveClient
var StreamExecve_ok bool
var StreamFileOpen RemoteServer_UploadFileOpenClient
var StreamFileOpen_ok bool
var StreamFileRead RemoteServer_UploadFileReadClient
var StreamFileRead_ok bool
var StreamFileWrite RemoteServer_UploadFileWriteClient
var StreamFileWrite_ok bool
var StreamTCPSendMsg RemoteServer_UploadTCPSendMsgClient
var StreamTCPSendMsg_ok bool
var StreamTCPRecvMsg RemoteServer_UploadTCPRecvMsgClient
var StreamTCPRecvMsg_ok bool
var StreamConnectIpv4 RemoteServer_UploadConnectIpv4Client
var StreamConnectIpv4_ok bool
var StreamUDPSendMsg RemoteServer_UploadUDPSendMsgClient
var StreamUDPSendMsg_ok bool
var StreamUDPRecvMsg RemoteServer_UploadUDPRecvMsgClient
var StreamUDPRecvMsg_ok bool
func RemoteServerConnect(addr string) error {
	var err error
	RemoteConn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	RemoteClient = NewRemoteServerClient(RemoteConn)
	StreamProcessFork, err = RemoteClient.UploadProcessFork(context.Background())
	if err != nil {
		StreamProcessFork_ok = false
		return err
	}
	StreamProcessFork_ok = true
	StreamProcessExit, err = RemoteClient.UploadProcessExit(context.Background())
	if err != nil {
		StreamProcessExit_ok = false
		return err
	}
	StreamProcessExit_ok = true
	StreamExecve, err = RemoteClient.UploadExecve(context.Background())
	if err != nil {
		StreamExecve_ok = false
		return err
	}
	StreamExecve_ok = true
	StreamFileOpen, err = RemoteClient.UploadFileOpen(context.Background())
	if err != nil {
		StreamFileOpen_ok = false
		return err
	}
	StreamFileOpen_ok = true
	StreamFileRead, err = RemoteClient.UploadFileRead(context.Background())
	if err != nil {
		StreamFileRead_ok = false
		return err
	}
	StreamFileRead_ok = true
	StreamFileWrite, err = RemoteClient.UploadFileWrite(context.Background())
	if err != nil {
		StreamFileWrite_ok = false
		return err
	}
	StreamFileWrite_ok = true
	StreamTCPSendMsg, err = RemoteClient.UploadTCPSendMsg(context.Background())
	if err != nil {
		StreamTCPSendMsg_ok = false
		return err
	}
	StreamTCPSendMsg_ok = true
	StreamTCPRecvMsg, err = RemoteClient.UploadTCPRecvMsg(context.Background())
	if err != nil {
		StreamTCPRecvMsg_ok = false
		return err
	}
	StreamTCPRecvMsg_ok = true
	StreamConnectIpv4, err = RemoteClient.UploadConnectIpv4(context.Background())
	if err != nil {
		StreamConnectIpv4_ok = false
		return err
	}
	StreamConnectIpv4_ok = true
	StreamUDPSendMsg, err = RemoteClient.UploadUDPSendMsg(context.Background())
	if err != nil {
		StreamUDPSendMsg_ok = false
		return err
	}
	StreamUDPSendMsg_ok = true
	StreamUDPRecvMsg, err = RemoteClient.UploadUDPRecvMsg(context.Background())
	if err != nil {
		StreamUDPRecvMsg_ok = false
		return err
	}
	StreamUDPRecvMsg_ok = true
	return err
}
func RemoteServerDisconnect() {
	StreamProcessFork.CloseAndRecv()
	StreamProcessExit.CloseAndRecv()
	StreamExecve.CloseAndRecv()
	StreamFileOpen.CloseAndRecv()
	StreamFileRead.CloseAndRecv()
	StreamFileWrite.CloseAndRecv()
	StreamTCPSendMsg.CloseAndRecv()
	StreamTCPRecvMsg.CloseAndRecv()
	StreamConnectIpv4.CloseAndRecv()
	StreamUDPSendMsg.CloseAndRecv()
	StreamUDPRecvMsg.CloseAndRecv()
	RemoteConn.Close()
}
func UploadProcessFork(info *ProcessForkInfo) {
	var err error
	if !StreamProcessFork_ok {
		StreamProcessFork, err = RemoteClient.UploadProcessFork(context.Background())
		if err != nil {
			return
		}
	}
	StreamProcessFork_ok = true
	err = StreamProcessFork.Send(info)
	if err != nil {
		StreamProcessFork_ok = false
	}
	return
}
func UploadProcessExit(info *ProcessExitInfo) {
	var err error
	if !StreamProcessExit_ok {
		StreamProcessExit, err = RemoteClient.UploadProcessExit(context.Background())
		if err != nil {
			return
		}
	}
	StreamProcessExit_ok = true
	err = StreamProcessExit.Send(info)
	if err != nil {
		StreamProcessExit_ok = false
	}
	return
}
func UploadExecve(info *ExecveInfo) {
	var err error
	if !StreamExecve_ok {
		StreamExecve, err = RemoteClient.UploadExecve(context.Background())
		if err != nil {
			return
		}
	}
	StreamExecve_ok = true
	err = StreamExecve.Send(info)
	if err != nil {
		StreamExecve_ok = false
	}
	return
}
func UploadFileOpen(info *FileOpenInfo) {
	var err error
	if !StreamFileOpen_ok {
		StreamFileOpen, err = RemoteClient.UploadFileOpen(context.Background())
		if err != nil {
			return
		}
	}
	StreamFileOpen_ok = true
	err = StreamFileOpen.Send(info)
	if err != nil {
		StreamFileOpen_ok = false
	}
	return
}
func UploadFileRead(info *FileReadInfo) {
	var err error
	if !StreamFileRead_ok {
		StreamFileRead, err = RemoteClient.UploadFileRead(context.Background())
		if err != nil {
			return
		}
	}
	StreamFileRead_ok = true
	err = StreamFileRead.Send(info)
	if err != nil {
		StreamFileRead_ok = false
	}
	return
}
func UploadFileWrite(info *FileWriteInfo) {
	var err error
	if !StreamFileWrite_ok {
		StreamFileWrite, err = RemoteClient.UploadFileWrite(context.Background())
		if err != nil {
			return
		}
	}
	StreamFileWrite_ok = true
	err = StreamFileWrite.Send(info)
	if err != nil {
		StreamFileWrite_ok = false
	}
	return
}
func UploadTCPSendMsg(info *TCPSendMsgInfo) {
	var err error
	if !StreamTCPSendMsg_ok {
		StreamTCPSendMsg, err = RemoteClient.UploadTCPSendMsg(context.Background())
		if err != nil {
			return
		}
	}
	StreamTCPSendMsg_ok = true
	err = StreamTCPSendMsg.Send(info)
	if err != nil {
		StreamTCPSendMsg_ok = false
	}
	return
}
func UploadTCPRecvMsg(info *TCPRecvMsgInfo) {
	var err error
	if !StreamTCPRecvMsg_ok {
		StreamTCPRecvMsg, err = RemoteClient.UploadTCPRecvMsg(context.Background())
		if err != nil {
			return
		}
	}
	StreamTCPRecvMsg_ok = true
	err = StreamTCPRecvMsg.Send(info)
	if err != nil {
		StreamTCPRecvMsg_ok = false
	}
	return
}
func UploadConnectIpv4(info *ConnectIpv4Info) {
	var err error
	if !StreamConnectIpv4_ok {
		StreamConnectIpv4, err = RemoteClient.UploadConnectIpv4(context.Background())
		if err != nil {
			return
		}
	}
	StreamConnectIpv4_ok = true
	err = StreamConnectIpv4.Send(info)
	if err != nil {
		StreamConnectIpv4_ok = false
	}
	return
}
func UploadUDPSendMsg(info *UDPSendMsgInfo) {
	var err error
	if !StreamUDPSendMsg_ok {
		StreamUDPSendMsg, err = RemoteClient.UploadUDPSendMsg(context.Background())
		if err != nil {
			return
		}
	}
	StreamUDPSendMsg_ok = true
	err = StreamUDPSendMsg.Send(info)
	if err != nil {
		StreamUDPSendMsg_ok = false
	}
	return
}
func UploadUDPRecvMsg(info *UDPRecvMsgInfo) {
	var err error
	if !StreamUDPRecvMsg_ok {
		StreamUDPRecvMsg, err = RemoteClient.UploadUDPRecvMsg(context.Background())
		if err != nil {
			return
		}
	}
	StreamUDPRecvMsg_ok = true
	err = StreamUDPRecvMsg.Send(info)
	if err != nil {
		StreamUDPRecvMsg_ok = false
	}
	return
}
