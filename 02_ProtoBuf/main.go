package main

import (
	"fmt"
	"kalo/hello"

	"github.com/golang/protobuf/proto"
)

func main() {
	p1 := &hello.People{
		Name: "kalo",
		Age:  32,
		Num:  1,
		Info : []byte("123456"),
	}
	fmt.Println(p1.Info)
	data, _ := proto.Marshal(p1)
	fmt.Println(data)

	p2 := &hello.People{}
	proto.Unmarshal(data, p2)
	fmt.Println(p2.Info)
	fmt.Println(len(p2.Info))

}
