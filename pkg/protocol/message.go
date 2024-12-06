package protocol

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Type    string `json:"type"`    
	SenderID  string `json:"senderID"`
	RecieverID string `json:"recieverID"`  
	Content string `json:"content"`
}


func SerializeMessage(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}

func DeserializeMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message: %w", err)
	}
	return &msg, nil
}
