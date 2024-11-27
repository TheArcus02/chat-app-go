package handlers

import (
	"fmt"

	"github.com/TheArcus02/chat-app-go/internal/models"
	"github.com/TheArcus02/chat-app-go/internal/services"
	"github.com/TheArcus02/chat-app-go/internal/utils"
)

type ChatHandler struct {
	Logger *utils.Logger
	Pool   *services.ConnectionPool
}

func NewChatHandler(logger *utils.Logger, pool *services.ConnectionPool) *ChatHandler {
	return &ChatHandler{
		Logger: logger,
		Pool:   pool,
	}
}

func (handler *ChatHandler) SendMessage(sender *models.User, recipientID string, message string) error {
	recipient := handler.Pool.Users[recipientID]
	if recipient == nil {
		return fmt.Errorf("recipient not found")
	}

	// Here, we're assuming that the recipient is connected and ready to receive messages
	// You could implement more sophisticated message routing here

	handler.Logger.Infof(fmt.Sprintf("Sending message from %s to %s: %s", sender.Username, recipient.Username, message))

	// In a real scenario, you'd push this message to the recipient's socket connection
	// For now, let's simulate that the message was "sent"
	return nil
}
