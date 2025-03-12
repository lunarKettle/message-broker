package network

import (
	"bufio"
	"fmt"
	"github.com/lunarKettle/message-broker/internal/broker"
	"log/slog"
	"net"
	"strings"
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
	clientAddr := conn.RemoteAddr().String()
	slog.Info("Client connected", "Address", clientAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		go func() {
			line := scanner.Text()
			slog.Info("Message received", "Address", clientAddr, "Message", line)

			parts := strings.SplitN(line, " ", 3)

			action := parts[0]
			queue := parts[1]
			message := ""
			if len(parts) == 3 {
				message = parts[2]
			}

			responseChan := make(chan string)
			command := &broker.Command{
				ClientID: clientAddr,
				Action:   action,
				Queue:    queue,
				Message:  message,
				Response: responseChan,
			}
			s.broker.ExecuteCommand(command)

			for response := range command.Response {
				if _, err := conn.Write([]byte(response + "\n")); err != nil {
					slog.Warn("failed to send response", "error", err)
				}
			}
		}()
	}

	if err := conn.Close(); err != nil {
		slog.Warn("failed to close connection", "error", err)
	}
	slog.Info("Client disconnected", "Address", clientAddr)
}
