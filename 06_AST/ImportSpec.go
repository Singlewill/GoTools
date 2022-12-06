package main

import (
	"go/ast"
	"go/token"
)

//生成Import语句,如下形式
/*	import (
		"context"
		"google.golang.org/grpc"
		"google.golang.org/grpc/credentials/insecure"
	)
*/
func genImportSpec() (node interface{}) {
	import_spec := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			0: &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: "\"context\"",
				},
			},
			1: &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: "\"google.golang.org/grpc\"",
				},
			},
			2: &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: "\"google.golang.org/grpc/credentials/insecure\"",
				},
			},
		},
	}
	return import_spec

}
