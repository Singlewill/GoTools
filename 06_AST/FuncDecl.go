package main

import (
	"go/ast"
	"go/token"
)

/*
genUploadFunc(): 生成函数定义，如下形式

	func UploadProcessFork(info *ProcessForkInfo) {
		var err error		//body 0
		if !StreamProcessFork_ok {		//body 1
			StreamProcessFork, err = RemoteClient.UploadProcessFork(context.Backgroud())
			if err != nil {
				return
			}
		}
		StreamProcessFork_ok = true		//body 2
		err = StreamProcessFork.Send(info)	//body 3
		if err != nil {		//body 4
			StreamProcessFork_ok = false
		}
		return			//body5
	}
*/
func genUploadFunc(key string) (node interface{}) {
	/////////////////////////////////////////////////
	c := &ast.FuncDecl{
		Name: &ast.Ident{ //函数名称"UploadProcessFork"
			Name: "Upload" + key,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					0: {
						Names: []*ast.Ident{
							0: {
								Name: "info",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: key + "Info",
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{},
	}
	//body0是err定义var err error
	body0 := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				0: &ast.ValueSpec{
					Names: []*ast.Ident{
						0: &ast.Ident{
							Name: "err",
						},
					},
					Type: &ast.Ident{
						Name: "error",
					},
				},
			},
		},
	}

	//body1, 是if语句块
	//if !StreamProcessFork_ok {
	//	StreamProcessFork, err = RemoteClient.UploadProcessFork(context.Background())
	//	if err != nil {
	//		return
	//	}
	//}
	body1 := &ast.IfStmt{
		Cond: &ast.UnaryExpr{ //UnaryExpr 一元表达式
			X: &ast.Ident{
				Name: "Stream" + key + "_ok",
			},
			Op: token.NOT,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				0: &ast.AssignStmt{ //body内第一个表达式 StreamProcessFork, err = RemoteClient.UploadProcessFork(context.Background())
					Lhs: []ast.Expr{
						0: &ast.Ident{
							Name: "Stream" + key,
						},
						1: &ast.Ident{
							Name: "err",
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						0: &ast.CallExpr{ //右值是函数调用
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "RemoteClient",
								},
								Sel: &ast.Ident{
									Name: "Upload" + key,
								},
							}, //函数是a.b形式
							Args: []ast.Expr{
								0: &ast.CallExpr{ //参数也是a.b()形式函数调用
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "context",
										},
										Sel: &ast.Ident{
											Name: "Backgroud",
										},
									},
									//无参数
								},
							},
						},
					},
				},
				//	if err != nil {
				//		return
				//	}
				1: &ast.IfStmt{
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
							0: &ast.ReturnStmt{},
						},
					},
				},
			},
		},
	}

	//body2 是streamProcessFork_ok = true
	body2 := &ast.AssignStmt{
		Lhs: []ast.Expr{ //左值 err
			0: &ast.Ident{
				Name: "Stream" + key + "_ok",
			},
		},
		Tok: token.ASSIGN, // "="
		Rhs: []ast.Expr{ //右值streamProcessFork.Send(info)
			0: &ast.Ident{
				Name: "true",
			},
		},
	}

	//body3是err = streamProcessFork.Send(info)
	body3 := &ast.AssignStmt{
		Lhs: []ast.Expr{ //左值 err
			0: &ast.Ident{
				Name: "err",
			},
		},
		Tok: token.ASSIGN, // "="
		Rhs: []ast.Expr{ //右值streamProcessFork.Send(info)
			0: &ast.CallExpr{ //右值是表达式
				Fun: &ast.SelectorExpr{ //函数属于a.b结构
					X: &ast.Ident{ //StreamProcessFork.Send
						Name: "Stream" + key,
					},
					Sel: &ast.Ident{
						Name: "Send",
					},
				},
				Args: []ast.Expr{ //参数info
					0: &ast.Ident{
						Name: "info",
					},
				},
			},
		},
	}

	//body4 是if语句
	//	if err != nil {
	//	streamProcessFork_ok = false
	//}
	//return t
	body4 := &ast.IfStmt{
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
							Name: "Stream" + key + "_ok",
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						0: &ast.Ident{
							Name: "false",
						},
					},
				},
			},
		},
	}

	body5 := &ast.ReturnStmt{}

	c.Body.List = append(c.Body.List, []ast.Stmt{body0, body1, body2, body3, body4, body5}...)
	return c
}

func genConnectFunc(keys []string) (node interface{}) {
	c := &ast.FuncDecl{
		//函数名称 "RemoteServerConnect",
		Name: &ast.Ident{
			Name: "RemoteServerConnect",
		},
		//函数类型，包括参数表，和返回值
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					0: {
						Names: []*ast.Ident{
							0: {
								Name: "addr",
							},
						},
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					0: {
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		//函数体，初始化状态为：
		/*
			RemoteConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}
			RemoteClient = NewRemoteServerClient(RemoteConn)
		*/
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				//body 0 : 即第一行
				0: &ast.AssignStmt{
					Lhs: []ast.Expr{
						//第一个左值
						0: &ast.Ident{
							Name: "RemoteConn",
						},
						//第二个左值
						1: &ast.Ident{
							Name: "err",
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						//第一个右值, 是函数调用CallExpr
						0: &ast.CallExpr{
							//函数名是a.b()形式
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "grpc",
								},
								Sel: &ast.Ident{
									Name: "Dial",
								},
							},
							Args: []ast.Expr{
								//参数0
								0: &ast.Ident{
									Name: "addr",
								},
								//参数1
								1: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "grpc",
										},
										Sel: &ast.Ident{
											Name: "WithTransportCredentials",
										},
									},
									Args: []ast.Expr{
										0: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "insecure",
												},
												Sel: &ast.Ident{
													Name: "NewCredentials",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				/*
					if err != nil {
						return err
					}
				*/
				1: &ast.IfStmt{
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
							0: &ast.ReturnStmt{
								Results: []ast.Expr{
									0: &ast.Ident{
										Name: "err",
									},
								},
							},
						},
					},
				},
				/*
					RemoteClient = NewRemoteServerClient(RemoteConn)
				*/
				2: &ast.AssignStmt{
					Lhs: []ast.Expr{
						0: &ast.Ident{
							Name: "RemoteClient",
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						0: &ast.CallExpr{ //右值是函数调用
							Fun: &ast.Ident{
								Name: "NewRemoteServerClient",
							},
							Args: []ast.Expr{
								0: &ast.Ident{ //参数也是a.b()形式函数调用
									Name: "RemoteConn",
								},
							},
						},
					},
				},
			},
		},
	}

	/*
		针对每一个功能，添加如下形式函数体
		StreamProcessFork, err := RemoteClient.UploadProcessFork(context.Background())
		if err != nil {
			StreamProcessFork_ok = false
			return err
		}
		StreamProcessFork_ok = true
	*/
	for _, key := range keys {
		body1 := &ast.AssignStmt{
			Lhs: []ast.Expr{
				//第一个左值
				0: &ast.Ident{
					Name: "Stream" + key,
				},
				//第二个左值
				1: &ast.Ident{
					Name: "err",
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				//第一个右值, 是函数调用CallExpr
				0: &ast.CallExpr{
					//函数名是a.b()形式
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "RemoteClient",
						},
						Sel: &ast.Ident{
							Name: "Upload" + key,
						},
					},
					Args: []ast.Expr{
						//参数0
						0: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "context",
								},
								Sel: &ast.Ident{
									Name: "Backgroud",
								},
							},
						},
					},
				},
			},
		}
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
								Name: "Stream" + key + "_ok",
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
		body3 := &ast.AssignStmt{
			Lhs: []ast.Expr{
				0: &ast.Ident{
					Name: "Stream" + key + "_ok",
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				0: &ast.Ident{
					Name: "true",
				},
			},
		}
		c.Body.List = append(c.Body.List, []ast.Stmt{body1, body2, body3}...)

	}

	return c
}
func genDisonnectFunc(keys []string) (node interface{}) {
	c := &ast.FuncDecl{
		//函数名称 "RemoteServerConnect",
		Name: &ast.Ident{
			Name: "RemoteServerDisconnect",
		},
		//参数表为空
		Type: &ast.FuncType{},
		//函数体为空
		Body: &ast.BlockStmt{},
	}
	for _, key := range keys {
		body := &ast.ExprStmt{
			X: &ast.CallExpr{ //右值是函数调用
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "Stream" + key,
					},
					Sel: &ast.Ident{
						Name: "CloseAndRecv",
					},
				}, //函数是a.b形式
			},
		}

		c.Body.List = append(c.Body.List, []ast.Stmt{body}...)
	}
	//添加结尾的 RemoteConn.Close()
	body_end := &ast.ExprStmt{
		X: &ast.CallExpr{ //右值是函数调用
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "RemoteConn",
				},
				Sel: &ast.Ident{
					Name: "Close",
				},
			}, //函数是a.b形式
		},
	}
	c.Body.List = append(c.Body.List, []ast.Stmt{body_end}...)
	return c

}
