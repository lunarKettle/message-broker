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
	slog.Info("Received command", "command", command)
	mb.commands <- command
}

func (mb *MessageBroker) publish(topic string, message string, responseCh chan<- string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	defer close(responseCh)
	mb.queues[topic] = append(mb.queues[topic], message)

	if subs, exists := mb.subs[topic]; exists {
		for _, sub := range subs {
			go func(s *Subscriber) {
				s.Ch <- message
			}(sub)
		}
	}
	responseCh <- "OK"
}

func (mb *MessageBroker) subscribe(clientID string, topic string, responseCh chan<- string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	sub := &Subscriber{
		ID:    clientID,
		Queue: topic,
		Ch:    responseCh,
	}

	mb.subs[topic] = append(mb.subs[topic], sub)
	responseCh <- "OK"
}

func (mb *MessageBroker) consume(topic string, responseCh chan<- string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	defer close(responseCh)

	var msg string
	if len(mb.queues[topic]) == 0 {
		msg = "Empty topic"
	} else {
		msg = mb.queues[topic][0]
		mb.queues[topic] = mb.queues[topic][1:]
	}

	responseCh <- msg
}
