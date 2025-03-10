package broker

import "sync"

type MessageBroker struct {
	queues   map[string][]string
	mu       sync.Mutex
	commands chan Command
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{}
}
