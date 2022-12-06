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
func RemoteServerConnect(addr string) error {
	RemoteConn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	RemoteClient = NewRemoteServerClient(RemoteConn)
	StreamProcessFork, err = RemoteClient.UploadProcessFork(context.Backgroud())
	if err != nil {
		StreamProcessFork_ok = false
		return err
	}
	StreamProcessFork_ok = true
	StreamProcessExit, err = RemoteClient.UploadProcessExit(context.Backgroud())
	if err != nil {
		StreamProcessExit_ok = false
		return err
	}
	StreamProcessExit_ok = true
	StreamExecve, err = RemoteClient.UploadExecve(context.Backgroud())
	if err != nil {
		StreamExecve_ok = false
		return err
	}
	StreamExecve_ok = true
	StreamFileOpen, err = RemoteClient.UploadFileOpen(context.Backgroud())
	if err != nil {
		StreamFileOpen_ok = false
		return err
	}
	StreamFileOpen_ok = true
	StreamFileRead, err = RemoteClient.UploadFileRead(context.Backgroud())
	if err != nil {
		StreamFileRead_ok = false
		return err
	}
	StreamFileRead_ok = true
	StreamFileWrite, err = RemoteClient.UploadFileWrite(context.Backgroud())
	if err != nil {
		StreamFileWrite_ok = false
		return err
	}
	StreamFileWrite_ok = true
}
func RemoteServerDisconnect() {
	StreamProcessFork.CloseAndRecv()
	StreamProcessExit.CloseAndRecv()
	StreamExecve.CloseAndRecv()
	StreamFileOpen.CloseAndRecv()
	StreamFileRead.CloseAndRecv()
	StreamFileWrite.CloseAndRecv()
	RemoteConn.Close()
}
func UploadProcessFork(info *ProcessForkInfo) {
	var err error
	if !StreamProcessFork_ok {
		StreamProcessFork, err = RemoteClient.UploadProcessFork(context.Backgroud())
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
		StreamProcessExit, err = RemoteClient.UploadProcessExit(context.Backgroud())
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
		StreamExecve, err = RemoteClient.UploadExecve(context.Backgroud())
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
		StreamFileOpen, err = RemoteClient.UploadFileOpen(context.Backgroud())
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
		StreamFileRead, err = RemoteClient.UploadFileRead(context.Backgroud())
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
		StreamFileWrite, err = RemoteClient.UploadFileWrite(context.Backgroud())
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
