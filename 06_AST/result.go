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
var streamProcessExecve RemoteServer_UploadProcessExecveClient
var streamProcessExecve_ok bool
var streamProcessRead RemoteServer_UploadProcessReadClient
var streamProcessRead_ok bool
var streamProcessWrite RemoteServer_UploadProcessWriteClient
var streamProcessWrite_ok bool
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
func UploadProcessExecve(info *ProcessExecveInfo) {
	var err error
	if !streamProcessExecve_ok {
		streamProcessExecve, err = RemoteClient.UploadProcessExecve(context.Background())
		if err != nil {
			return
		}
	}
	streamProcessExecve_ok = true
	err = streamProcessExecve.Send(info)
	if err != nil {
		streamProcessFork_ok = false
	}
}
func UploadProcessRead(info *ProcessReadInfo) {
	var err error
	if !streamProcessRead_ok {
		streamProcessRead, err = RemoteClient.UploadProcessRead(context.Background())
		if err != nil {
			return
		}
	}
	streamProcessRead_ok = true
	err = streamProcessRead.Send(info)
	if err != nil {
		streamProcessFork_ok = false
	}
}
func UploadProcessWrite(info *ProcessWriteInfo) {
	var err error
	if !streamProcessWrite_ok {
		streamProcessWrite, err = RemoteClient.UploadProcessWrite(context.Background())
		if err != nil {
			return
		}
	}
	streamProcessWrite_ok = true
	err = streamProcessWrite.Send(info)
	if err != nil {
		streamProcessFork_ok = false
	}
}
func RemoteServerDisconnect() {
	streamProcessWrite.CloseAndRecv()
	streamProcessRead.CloseAndRecv()
	streamProcessExecve.CloseAndRecv()
	streamProcessFork.CloseAndRecv()
	RemoteConn.Close()
}
func RemoteServerConnect(addr string) error {
	RemoteConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	RemoteClient = NewRemoteServerClient(RemoteConn)
	streamProcessFork, err := RemoteClient.UploadProcessFork(context.Background())
	if err != nil {
		streamProcessFork_ok = false
		return err
	}
	streamProcessFork_ok = true
	streamProcessExecve, err := RemoteClient.UploadProcessExecve(context.Background())
	if err != nil {
		streamProcessExecve_ok = false
		return err
	}
	streamProcessExecve_ok = true
	streamProcessRead, err := RemoteClient.UploadProcessRead(context.Background())
	if err != nil {
		streamProcessRead_ok = false
		return err
	}
	streamProcessRead_ok = true
	streamProcessWrite, err := RemoteClient.UploadProcessWrite(context.Background())
	if err != nil {
		streamProcessWrite_ok = false
		return err
	}
	streamProcessWrite_ok = true
}
