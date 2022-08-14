package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net"

	"github.com/DarthPestilane/easytcp"
	"github.com/sirupsen/logrus"
)

type publisher struct {
	config     configuration
	queue      string
	connection net.Conn
	packer     *easytcp.DefaultPacker
	log        *logrus.Logger
}

type Publisher interface {
	Connect() error
	Publish(data string) error
}

func (p *publisher) Connect() error {
	err := p.config.readEnv()
	if err != nil {
		return err
	}
	log.Println("connection string", p.config.getConnectionString())
	conn, err := net.Dial("tcp", p.config.getConnectionString())
	if err != nil {
		return err
	}
	log := logrus.New()
	packer := easytcp.NewDefaultPacker()
	p.connection = conn
	p.packer = packer

	p.log = log
	return nil
}

func NewPublisher(queue string, autoConnect bool) (publisher, error) {
	p := publisher{
		config:     &config{},
		queue:      queue,
		connection: nil,
		packer:     &easytcp.DefaultPacker{},
		log:        &logrus.Logger{},
	}
	if autoConnect {
		err := p.Connect()
		if err != nil {
			return p, err
		}

	}
	return p, nil
}

func (p *publisher) send(tcpMessage *easytcp.Message) error {
	packedBytes, err := p.packer.Pack(tcpMessage)
	if err != nil {
		return err
	}
	if _, err := p.connection.Write(packedBytes); err != nil {
		return err
	}
	p.log.Infof("[pbulisher] snd >>> | id:(%d) size:(%d) data: %s", tcpMessage.ID(), len(tcpMessage.Data()), tcpMessage.Data())
	return nil
}

func (p *publisher) Publish(data string) (string, error) {
	payload := make([]byte, 0)
	payload = append(payload, []byte(p.queue)...)
	payload = append(payload, []byte(" ")...)
	b64Data := base64.StdEncoding.EncodeToString([]byte(data))
	payload = append(payload, []byte(b64Data)...)
	msg := easytcp.NewMessage(MsgPublishTcpReq, []byte(payload))
	err := p.send(msg)
	if err != nil {
		return "", err
	}
	// wait for message id
	tcpmsg, err := p.packer.Unpack(p.connection)
	if err != nil {
		panic(err)
	}
	if tcpmsg.ID() == MsgPublishTcpAck && string(tcpmsg.Data()) != "error" {
		return string(tcpmsg.Data()), nil
	}
	defer p.connection.Close()
	return "", errors.New("impossible to publish this message")
}
