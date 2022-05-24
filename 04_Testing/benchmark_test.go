package main

import (
	"testing"
)

func BenchmarkAdd(b *testing.B) {

	for i := 0; i < b.N; i++ { //use b.N for looping
		Add(10, 5)
	}
}

func BenchmarkTimeConsumingFunction(b *testing.B) {

	b.StopTimer() //调用该函数停止压力测试的时间计数

	//做一些初始化的工作，例如读取文件数据，数据库连接之类的，
	//这样这些时间不影响我们测试函数本身的性能

	b.StartTimer() //重新开始时间

	for i := 0; i < b.N; i++ {
		Add(10, 5)
	}
}
