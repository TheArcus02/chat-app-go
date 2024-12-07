package handlers

import (
	"encoding/json"
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
	ChatHandler *ChatHandler
}

func NewConnectionHandler(logger *utils.Logger, connection *net.Conn, pool *services.ConnectionPool, chatHandler *ChatHandler) *ConnectionHandler {
	return &ConnectionHandler{
		Logger:     logger,
		Connection: connection,
		Pool:       pool,
		ChatHandler: chatHandler,
	}
}

func (handler *ConnectionHandler) Handle() {
	defer (*handler.Connection).Close()

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
	msg, err := protocol.DeserializeMessage(message)
	if err != nil {
		handler.Logger.Errorf(fmt.Sprintf("Failed to deserialize message: %v", err))
		return
	}

	switch msg.Type {
		case protocol.MessageTypeChat:
			handler.forwardToChatHandler(*msg)
		
		case protocol.MessageTypeConnect:
			userID := uuid.New().String()

			user := &models.User{
				ID:	   		userID,
				Username:	msg.Content,
				Conn:     	*handler.Connection,
			}
			handler.Pool.AddUser <- user

			if err := handler.sendConnectResponse(user); err != nil {
				handler.Logger.Errorf("Error sending connect response: %v", err)
				return
			}
		
			if err := handler.broadcastUserListUpdate(userID); err != nil {
				handler.Logger.Errorf("Error broadcasting user list update: %v", err)
			}

		case protocol.MessageTypeDisconnect:
			user := &models.User{
				ID:	   		msg.SenderID,
				Username:	msg.Content,
				Conn:     	*handler.Connection,
			}
			handler.Pool.RemoveUser <- user
			handler.Logger.Infof("User %s disconnected", msg.SenderID)

			if err := handler.broadcastUserListUpdate(user.ID); err != nil {
				handler.Logger.Errorf("Error broadcasting user list update: %v", err)
			}

		default:
			handler.Logger.Errorf(fmt.Sprintf("Unknown message type: %s", msg.Type))
	}

	handler.Pool.Logger.Infof(fmt.Sprintf("Message handled: %s", string(message)))
}

func (handler *ConnectionHandler) sendConnectResponse(user *models.User) error {
	connectResponse := protocol.Message{
		Type:   "connect_response",
		SenderID: "server",
		Content: string(func() []byte {
			content, err := json.Marshal(map[string]interface{}{
				"userID":   user.ID,
				"username": user.Username,
				"userList": handler.Pool.GetUsersList(),
			})
			if err != nil {
				handler.Logger.Errorf("Failed to marshal connect response content: %v", err)
			}
			return content
		}()),
	}

	responseJSON, err := json.Marshal(connectResponse)
	if err != nil {
		handler.Logger.Errorf("Failed to marshal connect response: %v", err)
		return err
	}

	err = user.SendMessage(string(responseJSON))
	if err != nil {
		handler.Logger.Errorf("Failed to send connect response to user %s: %v", user.Username, err)
		return err
	}

	handler.Logger.Infof("Connect response sent to user %s (ID: %s)", user.Username, user.ID)
	return nil
}


func (handler *ConnectionHandler) broadcastUserListUpdate(excludeUserID string) error {
	broadcastMessage := protocol.Message{
		Type:   "user_list_update",
		SenderID: "server",
		Content: string(func() []byte {
			content, err := json.Marshal(map[string]interface{}{
				"userList": handler.Pool.GetUsersList(),
			})
			if err != nil {
				handler.Logger.Errorf("Failed to marshal user list update content: %v", err)
			}
			return content
		}()),
	}

	broadcastJSON, err := json.Marshal(broadcastMessage)
	if err != nil {
		handler.Logger.Errorf("Failed to marshal broadcast message: %v", err)
		return err
	}

	handler.Pool.Broadcast(string(broadcastJSON), excludeUserID)
	handler.Logger.Infof("Broadcasted updated user list excluding user ID: %s", excludeUserID)
	return nil
}

func (handler *ConnectionHandler) forwardToChatHandler(msg protocol.Message) {
	go handler.ChatHandler.HandleChatMessage(msg)
}