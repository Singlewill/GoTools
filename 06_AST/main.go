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

func astToGo(dst *bytes.Buffer, node interface{}) error {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			log.Panicln(err)
		}
	}

	addNewline()

	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		return err
	}

	addNewline()

	return nil
}

func getValueSpec(decls []ast.Decl) *ast.GenDecl {
	//func getClientDecl(decls []ast.Decl) {
	for _, dl := range decls {
		t, ok := dl.(*ast.GenDecl)
		//虽然确定了GenDecl，但是GenDecl包含种类众多，比如import, const, type, var
		//因此要做二次判断
		if ok {
			//找值定义ValueSpec
			_, ok := t.Specs[0].(*ast.ValueSpec)
			if ok {
				return t
			}
		}

	}
	return nil
}

func createValueSpec(t *ast.GenDecl, key string) (node interface{}) {
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
	spec.Names[0].Name = "streamProcess" + key
	spec_type.Name = "RemoteServer_UploadProcess" + key + "Client"
	return t
}
func getFuncDecl(decls []ast.Decl) *ast.FuncDecl {
	//func getClientDecl(decls []ast.Decl) {
	for _, dl := range decls {
		t, ok := dl.(*ast.FuncDecl)
		if ok {
			return t
		}

	}
	return nil
}

func createFuncDecl(t *ast.FuncDecl, key string) (node interface{}) {
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
	fset := token.NewFileSet()
	//f, err := parser.ParseFile(fset, "remote_comm.go", nil, parser.ParseComments)
	f, err := parser.ParseFile(fset, "remote_comm.go", nil, 0)
	if err != nil {
		panic(err)
	}
	//ast.Print(fset, f)
	//找到指定的变量定义
	var buf bytes.Buffer
	spec := getValueSpec(f.Decls)
	if spec == nil {
		fmt.Println("getValueSpace Failed\n")
	}
	templ := getFuncDecl(f.Decls)
	if templ == nil {
		fmt.Println("getFuncDecl Failed\n")
	}

	spec_new := createValueSpec(spec, "Fork")
	err = astToGo(&buf, spec_new)
	if err != nil {
		return
	}

	func_new := createFuncDecl(templ, "Fork")
	err = astToGo(&buf, func_new)
	if err != nil {
		return
	}
	ioutil.WriteFile("tmp.txt", buf.Bytes(), 0666)

	/*
		var buf bytes.Buffer
		err = astToGo(&buf, client)
		if err != nil {
			return
		}
		ioutil.WriteFile("tmp.txt", buf.Bytes(), 0666)
	*/
	/*
		fmt.Println(f.Name.Name)
		var buf bytes.Buffer
		for _, fn := range f.Decls {
			err := astToGo(&buf, fn)
			if err != nil {
				return
			}
		}

		ioutil.WriteFile("tmp.txt", buf.Bytes(), 0666)
	*/
	/*
		var buf bytes.Buffer
		err = format.Node(&buf, fset, f)
		err = ioutil.WriteFile("tmp.txt", buf.Bytes(), 0666)
		fmt.Println(err)
	*/
}
