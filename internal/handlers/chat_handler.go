package handlers

import (
	"encoding/json"

	"github.com/TheArcus02/chat-app-go/internal/services"
	"github.com/TheArcus02/chat-app-go/internal/utils"
	"github.com/TheArcus02/chat-app-go/pkg/protocol"
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

func (ch *ChatHandler) HandleChatMessage(msg protocol.Message) {
	ch.Logger.Infof("Received chat message from %s to %s: %s", msg.SenderID, msg.RecieverID, msg.Content)

	chatResponse := protocol.Message{
		Type:    protocol.MessageTypeChat,
		SenderID:  msg.SenderID,
		RecieverID: msg.RecieverID,
		Content: msg.Content,
	}

	responseBytes, err := json.Marshal(chatResponse)
	if err != nil {
		ch.Logger.Errorf("Failed to serialize chat message: %v", err)
		return
	}

	err = ch.Pool.SendToUser(msg.RecieverID, string(responseBytes))
	if err != nil {
		ch.Logger.Errorf("Failed to send message to %s: %v", msg.RecieverID, err)
		return
	}
}
