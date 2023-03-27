package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/wonderivan/logger"
)

func sub_process(ch chan bool) {
	fmt.Printf("gorouting count = %d\n", runtime.NumGoroutine())
	//随机1~5秒
	time.Sleep(time.Duration((rand.Intn(5) + 1)) * time.Second)
	logger.Info("subProcess out")

	//任务结束
	<-ch
}

// 演示了一个控制gorouting最大数量的方法
func main() {
	//log := logger.NewLogger()
	//logger.SetLogger(`{"Console": {"level": "Warn", "color" : true}}`)
	//logr.SetLogger("file", `{"filename" : "123.log","Level": "Info", "append" : true, "permit" : "0666"}`)
	logger.SetLogger(logger.AdapterFile, `{"File": {"filename" : "123.log","LogLevel": "LevelInfo", "append" : true, "permit" : "0666"}}`)
	logger.Warn("subProcess out")
	//最大允许5个gorouting(包含main gorouting)
	/*
		ch := make(chan bool, 5)

		for {
			ch <- true
			go sub_process(ch)
		}
	*/
}
