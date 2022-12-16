#!/bin/sh
go build main.go ImportSpec.go ValueSpec.go FuncDecl.go
./main
cp ast_generate.go ../../remote_server/
