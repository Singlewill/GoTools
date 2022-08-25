package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

var FuncList = []string{"Fork", "Execve", "Read", "Write"}

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

func createUploadDecl(t *ast.FuncDecl, key string) (node interface{}) {
	//更改函数名
	t.Name.Name = "UploadProcess" + key

	/////////////////////////////////////////////////
	//参数名
	//param_name := t.Type.Params.List[0].Names[0].Name

	//参数类型
	param_type, ok := t.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return nil
	}
	param_type_indet, ok := param_type.X.(*ast.Ident)
	if !ok {
		return nil
	}
	param_type_indet.Name = "Process" + key + "Info"
	/////////////////////////////////////////////////
	//body0 是err定义
	//body0:= t.Body.List[0]		//var err error
	//body1, 是if语句块
	body1, ok := t.Body.List[1].(*ast.IfStmt) //if 语句块
	if !ok {
		return nil
	}
	//if条件
	body1_cond, ok := body1.Cond.(*ast.UnaryExpr)
	if !ok {
		return nil
	}
	body1_cond_ident, ok := body1_cond.X.(*ast.Ident)
	if !ok {
		return nil
	}
	//if语句条件中的变量名	<<--------------------
	body1_cond_ident.Name = "streamProcess" + key + "_ok"
	//if语句body内容
	body1_body_assign, ok := body1.Body.List[0].(*ast.AssignStmt) //streamProcessFork, err = RemoteClient.UploadProcessFork(context.Background())
	if !ok {
		return nil
	}
	//左值
	body1_body_assign_lhs0, ok := body1_body_assign.Lhs[0].(*ast.Ident)
	if !ok {
		return nil
	}
	//streamProcessFork  << ----------------------
	body1_body_assign_lhs0.Name = "streamProcess" + key

	//右值
	body1_body_assign_rhs0, ok := body1_body_assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return nil
	}
	body1_body_rhs0_fun, ok := body1_body_assign_rhs0.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	//UploadProcessFork << ----------------------------
	body1_body_rhs0_fun.Sel.Name = "UploadProcess" + key

	//body2 是streamProcessFork_ok = true
	body2, ok := t.Body.List[2].(*ast.AssignStmt)
	if !ok {
		return nil
	}
	body2_lhs0, ok := body2.Lhs[0].(*ast.Ident)
	if !ok {
		return nil
	}
	//streamProcessFork_ok  << ---------------------------
	body2_lhs0.Name = "streamProcess" + key + "_ok"

	//body3是err = streamProcessFork.Send(info)
	body3, ok := t.Body.List[3].(*ast.AssignStmt)
	body3_rhs0, ok := body3.Rhs[0].(*ast.CallExpr)
	body3_rhs0_fun, ok := body3_rhs0.Fun.(*ast.SelectorExpr)
	body3_rhs0_fun_x, ok := body3_rhs0_fun.X.(*ast.Ident)
	//streamProcessFork
	body3_rhs0_fun_x.Name = "streamProcess" + key

	return t
}

func main() {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "template.go", nil, 0)
	if err != nil {
		panic(err)
	}
	//找到import语句
	import_spec := getImportSpec(f.Decls)
	if import_spec == nil {
		fmt.Println("import spec not found\n")
		return
	}
	//找到四个变量定义
	remote_conn_spec := getVarSpec(f.Decls, "RemoteConn")
	if remote_conn_spec == nil {
		fmt.Println("RemoteConn not found\n")
		return
	}
	remote_client_spec := getVarSpec(f.Decls, "RemoteClient")
	if remote_client_spec == nil {
		fmt.Println("RemoteClient not found\n")
	}
	stream_client_spec := getVarSpec(f.Decls, "streamProcessFork")
	if stream_client_spec == nil {
		fmt.Println("streamProcessFork not found\n")
	}
	stream_client_ok_spec := getVarSpec(f.Decls, "streamProcessFork_ok")
	if stream_client_ok_spec == nil {
		fmt.Println("streamProcessFork_ok not found\n")
	}
	//找到3个函数定义
	connect_func := getFuncDecl(f.Decls, "RemoteServerConnect")
	if connect_func == nil {
		fmt.Println("Disconnect not found\n")
	}
	disconnect_func := getFuncDecl(f.Decls, "RemoteServerDisconnect")
	if disconnect_func == nil {
		fmt.Println("RemoteServerDisconnect not found\n")
	}
	upload_func := getFuncDecl(f.Decls, "UploadProcessFork")
	if upload_func == nil {
		fmt.Println("UploadProcessFork not found\n")
	}
	//////////////////////////////////////////////////////////////////////////////
	//准备文件写入
	//0, 写入package语句和Import语句
	buf.Write([]byte("package remote_server\n"))
	err = astToGo(&buf, import_spec)
	if err != nil {
		return
	}
	//1, 写入开头的两个变量，原样写入
	err = astToGo(&buf, remote_conn_spec)
	if err != nil {
		return
	}
	err = astToGo(&buf, remote_client_spec)
	if err != nil {
		return
	}
	//2, 按照开头的后两个变量，生成其他服务变量声明
	for _, s := range FuncList {
		spec := createValueSpec(stream_client_spec, "streamProcess"+s, "RemoteServer_UploadProcess"+s+"Client")
		err = astToGo(&buf, spec)
		if err != nil {
			return
		}
		spec = createValueSpec(stream_client_ok_spec, "streamProcess"+s+"_ok", "bool")
		err = astToGo(&buf, spec)
		if err != nil {
			return
		}
	}

	//3, 生成其他upload函数
	for _, s := range FuncList {
		funcNode := createUploadDecl(upload_func, s)
		err = astToGo(&buf, funcNode)
		if err != nil {
			return
		}
	}

	//4, 生成disconnect函数
	for _, s := range FuncList {
		//定义新的表达式语句。类似A.B()形式
		new_exprStmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "streamProcess" + s,
					},
					Sel: &ast.Ident{
						Name: "CloseAndRecv",
					},
				},
			},
		}
		//将新表达式追加到原来的语句前面
		disconnect_func.Body.List = append([]ast.Stmt{new_exprStmt}, disconnect_func.Body.List...)

	}
	err = astToGo(&buf, disconnect_func)
	if err != nil {
		return
	}

	//5, 生成connect函数
	for _, s := range FuncList {
		/*
			body1 -->
			streamProcessFork, err := RemoteClient.UploadProcessFork(context.Background())
		*/
		body1 := &ast.AssignStmt{
			Lhs: []ast.Expr{ //两个左值
				0: &ast.Ident{
					Name: "streamProcess" + s,
					Obj: &ast.Object{
						Kind: ast.Var,
						//Name: "RemoteConn"
					},
				},
				1: &ast.Ident{
					Name: "err",
					Obj: &ast.Object{
						Kind: ast.Var,
						//Name: "err"
					},
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{ //1个右值RemoteClient.UploadProcessFork(context.Background())
				0: &ast.CallExpr{
					Fun: &ast.SelectorExpr{ // ast.SelectorExpr代表A.B
						X: &ast.Ident{
							Name: "RemoteClient",
						},
						Sel: &ast.Ident{
							Name: "UploadProcess" + s,
						},
					},
					Args: []ast.Expr{ //参数，又是一个A.B()形式
						0: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "context",
								},
								Sel: &ast.Ident{
									Name: "Background",
								},
							},
						},
					},
				},
			},
		}

		/*
			body2 -->
				if err != nil {
					streamProcessFork_ok = false
					return err
				}

		*/
		body2 := &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					0: &ast.AssignStmt{
						Lhs: []ast.Expr{
							0: &ast.Ident{
								Name: "streamProcess" + s + "_ok",
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							0: &ast.Ident{
								Name: "false",
							},
						},
					},
					1: &ast.ReturnStmt{
						Results: []ast.Expr{
							0: &ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		}

		/*
			body3 -->
				streamProcessFork_ok = true
		*/
		body3 := &ast.AssignStmt{
			Lhs: []ast.Expr{
				0: &ast.Ident{
					Name: "streamProcess" + s + "_ok",
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				0: &ast.Ident{
					Name: "true",
				},
			},
		}
		//将3个body追加到函数体后
		connect_func.Body.List = append(connect_func.Body.List, []ast.Stmt{body1, body2, body3}...)

	}
	err = astToGo(&buf, connect_func)
	if err != nil {
		return
	}

	ioutil.WriteFile("result.go", buf.Bytes(), 0666)
}
