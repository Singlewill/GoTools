package main

import (
	"fmt"
)

func Add(a int, b int) int {
	return a + b
}

func Mul(a int, b int) int {
	return a * b
}

func main() {
	fmt.Println(Add(10, 20))
}
