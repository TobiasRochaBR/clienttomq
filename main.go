package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	c := NewConsumer("general")
	c.Connect()
	c.Register()
	//c.PublishSomething()
	p, err := NewPublisher("general", true)
	if err != nil {
		panic(err)
	}
	jsonFile, err := os.Open("test.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	msgId, err1 := p.Publish(string(byteValue))
	if err1 != nil {
		panic(err1)
	}
	log.Println("published id", msgId)
	go c.Consume()
	for {
		time.Sleep(time.Second)
		for _, msg := range c.Messages {
			log.Println("READY TO GO", msg.id)
			time.Sleep(time.Second)
			worker(msg.body)
			c.Ack()
		}
	}
}

func worker(data []byte) {
	println("WorKER DATA", string(data))
	time.Sleep(120 * time.Millisecond)
	return
}
