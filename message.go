package main

const (
	MESSAGE_WAITING_NACK = 1
)

type message struct {
	id    string
	queue string
	body  []byte
}

type Message interface {
	Id() string
	Queue() string
	Body() []byte
}

func (m *message) Id() string {
	return m.id
}
func (m *message) Queue() string {
	return m.queue
}
func (m *message) Body() []byte {
	return m.body
}

func newMessage(id string, queueName string, body []byte) *message {
	m := message{
		id:    id,
		queue: queueName,
		body:  body,
	}
	return &m
}
