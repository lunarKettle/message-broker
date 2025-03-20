package main

import (
	"github.com/lunarKettle/message-broker/internal/broker"
	"github.com/lunarKettle/message-broker/internal/network"
	"os"
)

const (
	ServerAddressEnv = "SERVER_ADDRESS"
)

func main() {
	serverAddress := os.Getenv(ServerAddressEnv)
	messageBroker := broker.NewMessageBroker()
	tcpServer := network.NewTCPServer(serverAddress, messageBroker)

	if err := tcpServer.Listen(); err != nil {
		panic(err)
	}
}
