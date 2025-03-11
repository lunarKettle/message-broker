package broker

import (
	"log/slog"
	"sync"
)

type MessageBroker struct {
	queues   map[string][]string
	subs     map[string][]*Subscriber
	mu       sync.Mutex
	commands chan *Command
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		queues:   make(map[string][]string),
		subs:     make(map[string][]*Subscriber),
		commands: make(chan *Command),
	}
}

func (mb *MessageBroker) Start() {
	go func() {
		for command := range mb.commands {
			slog.Info("Received command", "command", command)
			switch command.Action {
			case "PUBLISH":
				mb.publish(command.Queue, command.Message, command.Response)
			case "SUBSCRIBE":
				mb.subscribe(command.ClientID, command.Queue, command.Response)
			case "CONSUME":
				mb.consume(command.Queue, command.Response)
			}
		}
	}()
}

func (mb *MessageBroker) ExecuteCommand(command *Command) {
	mb.commands <- command
}

func (mb *MessageBroker) publish(queue string, message string, responseCh chan string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.queues[queue] = append(mb.queues[queue], message)

	if subs, exists := mb.subs[queue]; exists {
		for _, sub := range subs {
			go func(s *Subscriber) {
				s.Ch <- message
			}(sub)
		}
	}
	responseCh <- "OK"
	close(responseCh)
}

func (mb *MessageBroker) subscribe(clientID string, queue string, responseCh chan string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	sub := &Subscriber{
		ID:    clientID,
		Queue: queue,
		Ch:    responseCh,
	}

	mb.subs[queue] = append(mb.subs[queue], sub)
	responseCh <- "OK"
}

func (mb *MessageBroker) consume(queue string, responseCh chan string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	msg := mb.queues[queue][0]
	mb.queues[queue] = mb.queues[queue][1:]
	responseCh <- msg
	close(responseCh)
}
