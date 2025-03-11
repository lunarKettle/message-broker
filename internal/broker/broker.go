package broker

import (
	"log/slog"
	"sync"
)

type MessageBroker struct {
	queues   map[string][]string
	mu       sync.Mutex
	commands chan *Command
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		mu:       sync.Mutex{},
		commands: make(chan *Command),
	}
}

func (mb *MessageBroker) Start() {
	go func() {
		for command := range mb.commands {
			slog.Info("Received command", "command", command)
		}
	}()
}

func (mb *MessageBroker) ExecuteCommand(command *Command) {
	mb.commands <- command
}
