package main

import (
	"go/ast"
	"go/token"
)

/*
生成值定义，如下形式：
	var StreamProcessFork RemoteServer_UploadProcessForkClient
	var StreamProcessFork_ok bool
*/

func genValueSpec(var_name string, type_name string) (node interface{}) {
	value_space := &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			0: &ast.ValueSpec{
				Names: []*ast.Ident{
					0: &ast.Ident{
						Name: var_name,
					},
				},
				Type: &ast.Ident{
					Name: type_name,
				},
			},
		},
	}
	return value_space

}
