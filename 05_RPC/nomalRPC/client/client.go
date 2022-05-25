package main

import (
	"fmt"
	"log"
	"net/rpc"
)

var (
	HelloServiceName = "path/to/pkg.HelloService"
)

type HelloServiceClient struct {
	*rpc.Client
}

var HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return HelloServiceClient{}, err
	}
	//return &HelloServiceClient{Client: c}, nil
	return HelloServiceClient{c}, nil
}

func (p *HelloServiceClient) Hello(request string, reply *string) error {
	//return p.Client.Call(HelloServiceName+".Hello", request, reply)
	return p.Call(HelloServiceName+".Hello", request, reply)
}

func main() {
	client, err := DialHelloService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Hello("hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
