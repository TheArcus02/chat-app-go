package handlers

import (
	"fmt"
	"net"

	"github.com/TheArcus02/chat-app-go/internal/services"
	"github.com/TheArcus02/chat-app-go/internal/utils"
)

type ConnectionHandler struct {
	Logger     *utils.Logger
	Connection *net.Conn
	Pool       *services.ConnectionPool
}

func NewConnectionHandler(logger *utils.Logger, connection *net.Conn, pool *services.ConnectionPool) *ConnectionHandler {
	return &ConnectionHandler{
		Logger:     logger,
		Connection: connection,
		Pool:       pool,
	}
}

func (handler *ConnectionHandler) Handle() {
	buffer := make([]byte, 1024)
	for {
		n, err := (*handler.Connection).Read(buffer)
		if err != nil {
			handler.Logger.Errorf(fmt.Sprintf("Error reading from connection: %v", err))
			break
		}
		if n == 0 {
			continue
		}

		handler.handleMessage(buffer[:n])
	}
}

func (handler *ConnectionHandler) handleMessage(message []byte) {
	handler.Logger.Infof(fmt.Sprintf("Received message: %s", string(message)))

	// Here you can add message handling logic, such as routing messages to the correct user
	// Example: handle the message and send it to the correct user through the connection pool
	handler.Pool.Logger.Infof(fmt.Sprintf("Message handled: %s", string(message)))
}
