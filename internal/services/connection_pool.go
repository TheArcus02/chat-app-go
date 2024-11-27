package services

import (
	"sync"

	"github.com/TheArcus02/chat-app-go/internal/models"
	"github.com/TheArcus02/chat-app-go/internal/utils"
)

type ConnectionPool struct {
	Users  map[string]*models.User
	Mutex  sync.Mutex
	Logger *utils.Logger
	AddUser    chan *models.User
	RemoveUser chan *models.User
}

func NewConnectionPool(logger *utils.Logger) *ConnectionPool {
	return &ConnectionPool{
		Users:  make(map[string]*models.User),
		Logger: logger,
		AddUser:    make(chan *models.User),
		RemoveUser: make(chan *models.User),
	}
}

func (cp *ConnectionPool) Run() {
	for {
		select {
		case user := <-cp.AddUser:
			cp.Mutex.Lock()
			cp.Users[user.ID] = user
			cp.Logger.Infof("User %s connected", user.Username)
			cp.Mutex.Unlock()
		case user := <-cp.RemoveUser:
			cp.Mutex.Lock()
			delete(cp.Users, user.ID)
			cp.Logger.Infof("User %s disconnected", user.Username)
			cp.Mutex.Unlock()
		}
	}
}

func (cp *ConnectionPool) Broadcast(message string, senderID string) {
	cp.Mutex.Lock()
	defer cp.Mutex.Unlock()

	for id, user := range cp.Users {
		if id != senderID {
			err := user.SendMessage(message)
			if err != nil {
				cp.Logger.Errorf("Failed to send message to %s: %v", user.Username, err)
			}
		}
	}
}

func (cp *ConnectionPool) Shutdown() {
	cp.Mutex.Lock()
	defer cp.Mutex.Unlock()
	for _, user := range cp.Users {
		user.Conn.Close()
	}
	cp.Logger.Infof("All connections closed.")
}