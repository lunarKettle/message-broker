package main

import (
	"github.com/lunarKettle/message-broker/internal/broker"
	"github.com/lunarKettle/message-broker/internal/network"
)

func main() {
	messageBroker := broker.NewMessageBroker()
	messageBroker.Start()

	tcpServer := network.NewTCPServer(":8080", messageBroker)
	if err := tcpServer.Listen(); err != nil {
		panic(err)
	}
}
