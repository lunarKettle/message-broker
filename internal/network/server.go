package network

import (
	"fmt"
	"github.com/lunarKettle/message-broker/internal/broker"
	"log/slog"
	"net"
)

type TCPServer struct {
	address  string
	listener net.Listener
	broker   *broker.MessageBroker
}

func NewTCPServer(addr string, broker *broker.MessageBroker) *TCPServer {
	return &TCPServer{
		address: addr,
		broker:  broker,
	}
}

func (s *TCPServer) Listen() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener

	// Условие?
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			slog.Warn("failed to accept connection", "error", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *TCPServer) Stop() error {
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("failed to close server: %w", err)
	}
	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn) {

}
