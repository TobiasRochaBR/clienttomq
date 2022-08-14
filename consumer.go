package main

import (
	"bytes"
	"encoding/base64"
	"log"
	"net"

	"github.com/DarthPestilane/easytcp"
	"github.com/sirupsen/logrus"
)

type consumer struct {
	config         configuration
	id             string
	status         int
	queue          string
	tcpMessages    []*easytcp.Message
	Messages       []*message
	connection     net.Conn
	packer         *easytcp.DefaultPacker
	log            *logrus.Logger
	workingMessage *message
}

type Consumer interface {
	Connect() error
	Register() error
	ReceiveMessage(tcpMsgId int, data []byte)
	Consume()
	Ack() error
	NAck() error
	Reject()
	send(tcpMessage *easytcp.Message) error
}

func (c *consumer) Consume() {
	go func() {
		// read loop
		for {
			msg, err := c.packer.Unpack(c.connection)
			if err != nil {
				panic(err)
			}
			if msg.ID() == MsgDistributeTcpReq && string(msg.Data()) != "error" {
				c.ReceiveMessage(msg.Data())

			}
			//c.log.Infof("[consumer] rec <<< | id:(%d) size:(%d) data: %s", msg.ID(), len(msg.Data()), msg.Data())
		}
	}()
	select {}
}

func (c *consumer) send(tcpMessage *easytcp.Message) error {
	packedBytes, err := c.packer.Pack(tcpMessage)
	if err != nil {
		return err
	}
	if _, err := c.connection.Write(packedBytes); err != nil {
		return err
	}
	//c.log.Infof("[consumer] snd >>> | id:(%d) size:(%d) data: %s", tcpMessage.ID(), len(tcpMessage.Data()), tcpMessage.Data())
	return nil
}

func (c *consumer) Connect() error {
	err := c.config.readEnv()
	if err != nil {
		return err
	}
	log.Println("connection string", c.config.getConnectionString())
	conn, err := net.Dial("tcp", c.config.getConnectionString())
	if err != nil {
		return err
	}
	log := logrus.New()
	packer := easytcp.NewDefaultPacker()
	c.connection = conn
	c.packer = packer

	c.log = log
	return nil
}

func remove(slice []*message, s int) []*message {
	return append(slice[:s], slice[s+1:]...)
}

func (c *consumer) Ack() error {
	msg := c.workingMessage
	payload := []byte(c.queue)
	payload = append(payload, []byte(" ")...)
	payload = append(payload, []byte(msg.id)...)
	tcpMessage := easytcp.NewMessage(MsgAckTcpReq, payload)
	c.Messages = remove(c.Messages, 0)
	c.workingMessage = nil
	return c.send(tcpMessage)

}

func (c *consumer) NAck() error {
	msg := c.workingMessage
	payload := []byte(c.queue)
	payload = append(payload, []byte(" ")...)
	payload = append(payload, []byte(msg.id)...)
	tcpMessage := easytcp.NewMessage(MsgNAckTcpReq, payload)

	return c.send(tcpMessage)

}

func NewConsumer(queue string) consumer {
	c := consumer{
		config:      &config{},
		id:          "",
		status:      0,
		queue:       queue,
		tcpMessages: []*easytcp.Message{},
	}
	return c
}

func (c *consumer) Register() error {
	msg := easytcp.NewMessage(ConsumerRegisterTcpReq, []byte(c.queue))
	return c.send(msg)
}

func parseMessageData(data []byte, qMsgId *string, qMsgData *[]byte) bool {
	// for future implementations
	return true
}

func (c *consumer) ReceiveMessage(data []byte) {
	rawDecodedText, err := base64.StdEncoding.DecodeString(string(bytes.Split(data, []byte(" "))[1]))
	if err != nil {
		panic(err)
	}
	var qMsgId string
	var qMsgData []byte
	qMsgId = string(bytes.TrimSpace(bytes.Split(data, []byte(" "))[0]))
	qMsgData, err = base64.StdEncoding.DecodeString(string(rawDecodedText))
	if err != nil {
		panic(err)
	}
	qMsg := newMessage(qMsgId, c.queue, qMsgData)
	if parseMessageData(data, &qMsgId, &qMsgData) {

		c.Messages = append(c.Messages, qMsg)
		c.workingMessage = qMsg
		c.NAck()

	} else {

		c.workingMessage = qMsg
		c.NAck()
	}

}
