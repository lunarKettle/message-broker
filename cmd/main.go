package main

import (
	"github.com/lunarKettle/message-broker/internal/broker"
	"github.com/lunarKettle/message-broker/internal/network"
	"log/slog"
)

func main() {
	messageBroker := broker.NewMessageBroker()
	tcpServer := network.NewTCPServer(":8080", messageBroker)

	if tcpServer.Listen() == nil {
		slog.Error("Error starting TCP server")
	}
}
