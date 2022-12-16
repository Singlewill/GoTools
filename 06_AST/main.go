package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"log"
)

var FuncList = []string{"ProcessFork", "ProcessExit", "Execve", "FileOpen", "FileRead", "FileWrite", "TCPSendMsg", "TCPRecvMsg", "ConnectIpv4", "UDPSendMsg", "UDPRecvMsg"}

func astToGo(dst *bytes.Buffer, node interface{}) error {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			log.Panicln(err)
		}
	}

	//addNewline()

	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		return err
	}

	addNewline()

	return nil
}

func getImportSpec(decls []ast.Decl) *ast.GenDecl {
	for _, dl := range decls {
		t, ok := dl.(*ast.GenDecl)
		//虽然确定了GenDecl，但是GenDecl包含种类众多，比如import, const, type, var
		//因此要做二次判断
		if ok && t.Tok == token.IMPORT {
			return t
		}

	}
	return nil

}
func getVarSpec(decls []ast.Decl, key string) *ast.GenDecl {
	//func getClientDecl(decls []ast.Decl) {
	for _, dl := range decls {
		t, ok := dl.(*ast.GenDecl)
		//虽然确定了GenDecl，但是GenDecl包含种类众多，比如import, const, type, var
		//因此要做二次判断
		if ok {
			//找值定义ValueSpec
			spec, ok := t.Specs[0].(*ast.ValueSpec)
			if ok && spec.Names[0].Name == key {
				return t
			}
		}

	}
	return nil
}

func createValueSpec(t *ast.GenDecl, var_name string, type_name string) (node interface{}) {
	//断言变量定义
	spec, ok := t.Specs[0].(*ast.ValueSpec)
	if !ok {
		return nil
	}
	//断言类型定义
	spec_type, ok := spec.Type.(*ast.Ident)
	if !ok {
		return nil
	}
	spec.Names[0].Name = var_name
	spec_type.Name = type_name
	return t
}
func getFuncDecl(decls []ast.Decl, key string) *ast.FuncDecl {
	//func getClientDecl(decls []ast.Decl) {
	for _, dl := range decls {
		t, ok := dl.(*ast.FuncDecl)
		if ok && t.Name.Name == key {
			return t
		}

	}
	return nil
}

func main() {
	var buf bytes.Buffer
	var err error
	//准备文件写入
	//0, 写入package语句
	buf.Write([]byte("package remote_server\n"))
	//1,Import语句
	import_spec := genImportSpec()
	err = astToGo(&buf, import_spec)
	if err != nil {
		return
	}
	//2, 值定义语句
	//2-1:var RemoteConn *grpc.ClientConn
	value_spec := genValueSpec("RemoteConn", "*grpc.ClientConn")
	err = astToGo(&buf, value_spec)
	if err != nil {
		return
	}
	//2-2:var RemoteClient RemoteServerClient
	value_spec = genValueSpec("RemoteClient", "RemoteServerClient")
	err = astToGo(&buf, value_spec)
	if err != nil {
		return
	}
	//2-3:
	//	var StreamProcessFork RemoteServer_UploadProcessForkClient
	//	var StreamProcessFork_ok bool
	for _, s := range FuncList {
		value_spec = genValueSpec("Stream"+s, "RemoteServer_Upload"+s+"Client")
		err = astToGo(&buf, value_spec)
		if err != nil {
			return
		}
		value_spec = genValueSpec("Stream"+s+"_ok", "bool")
		err = astToGo(&buf, value_spec)
		if err != nil {
			return
		}
	}

	//3,Connect函数
	func_spec := genConnectFunc(FuncList)
	err = astToGo(&buf, func_spec)
	if err != nil {
		return
	}

	//4, Disconnect函数
	func_spec = genDisonnectFunc(FuncList)
	err = astToGo(&buf, func_spec)
	if err != nil {
		return
	}
	//5, Upload函数
	for _, key := range FuncList {
		func_spec = genUploadFunc(key)
		err = astToGo(&buf, func_spec)
		if err != nil {
			return
		}
	}

	ioutil.WriteFile("ast_generate.go", buf.Bytes(), 0666)
}
