package handlers

import (
	"fmt"
	"net"

	"github.com/TheArcus02/chat-app-go/internal/models"
	"github.com/TheArcus02/chat-app-go/internal/services"
	"github.com/TheArcus02/chat-app-go/internal/utils"
	"github.com/TheArcus02/chat-app-go/pkg/protocol"

	"github.com/google/uuid"
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

	// Deserialize the message
	msg, err := protocol.DeserializeMessage(message)
	if err != nil {
		handler.Logger.Errorf(fmt.Sprintf("Failed to deserialize message: %v", err))
		return
	}

	switch msg.Type {
		case protocol.MessageTypeText:
			handler.Pool.Broadcast(msg.Content, msg.Sender)
		
		case protocol.MessageTypeConnect:
			userID := uuid.New().String()

			user := &models.User{
				ID:	   		userID,
				Username:	msg.Sender,
				Conn:     	*handler.Connection,
			}
			handler.Pool.AddUser <- user

			handler.Logger.Infof("User %s (ID: %s) connected", msg.Sender, userID)
	
		case protocol.MessageTypeDisconnect:
			user := &models.User{
				ID:	   		msg.Sender,
				Username:	msg.Sender,
				Conn:     	*handler.Connection,
			}
			handler.Pool.RemoveUser <- user

			handler.Logger.Infof("User %s disconnected", msg.Sender)
		default:
			handler.Logger.Errorf(fmt.Sprintf("Unknown message type: %s", msg.Type))
	}

	handler.Pool.Logger.Infof(fmt.Sprintf("Message handled: %s", string(message)))
}
