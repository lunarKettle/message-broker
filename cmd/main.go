package main

import (
	"github.com/lunarKettle/message-broker/internal/broker"
	"github.com/lunarKettle/message-broker/internal/network"
	"log/slog"
	"os"
)

const (
	ServerAddressEnv = "SERVER_ADDRESS"
)

func main() {
	serverAddress := os.Getenv(ServerAddressEnv)
	messageBroker := broker.NewMessageBroker()
	tcpServer := network.NewTCPServer(serverAddress, messageBroker)

	slog.Info("Starting server", "Address", serverAddress)
	if err := tcpServer.Listen(); err != nil {
		panic(err)
	}
}
