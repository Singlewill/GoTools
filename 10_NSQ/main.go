package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

//生产者
func startProducer() {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", cfg)
	if err != nil {
		log.Fatal("NewProducer failed")
	}
	defer producer.Stop()
	//发布消息
	for {
		if err := producer.Publish("test", []byte("test message")); err != nil {
			log.Fatal("Publish failed")
		}
		fmt.Println("publis test")
		time.Sleep(10 * time.Second)
	}

}

//消费者
func startConsumer() {
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("test", "sensor01", cfg)
	if err != nil {
		log.Fatal("NewConsumer failed")
	}

	consumer2, err := nsq.NewConsumer("test", "sensor02", cfg)
	if err != nil {
		log.Fatal("NewConsumer failed")
	}
	//设置消息处理函数
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Println("01-----", string(message.Body))
		return nil
	}))
	consumer2.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Println("02------", string(message.Body))
		return nil
	}))
	/*
		// 连接到单例nsqd
		if err := consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
			log.Fatal("Consumer connect failed failed")
		}
	*/
	//连接到单例nsqlookupd
	if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		log.Fatal("Consumer connect failed failed")
	}
	if err := consumer2.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		log.Fatal("Consumer connect failed failed")
	}
	<-consumer.StopChan
	<-consumer2.StopChan
	fmt.Println("consumer out !")
}
func main() {
	go startProducer()
	time.Sleep(1 * time.Second)
	startConsumer()
}
