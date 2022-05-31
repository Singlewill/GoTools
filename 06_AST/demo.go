package main

import (
	"context"
	"fmt"
)

// Foo 结构体
type Foo struct {
	i int
}

// Bar 接口
type Bar interface {
	Do(ctx context.Context) error
}

// main方法
func main() {
	a := 1
	fmt.Println(a)
}
