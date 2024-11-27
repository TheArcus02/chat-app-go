package models

import (
	"net"
	"sync"
)

type User struct {
	ID       string      
	Username string      
	Conn     net.Conn    
	Mutex    sync.Mutex
}

func NewUser(id, username string, conn net.Conn) *User {
	return &User{
		ID:       id,
		Username: username,
		Conn:     conn,
	}
}

func (u *User) SendMessage(message string) error {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	_, err := u.Conn.Write([]byte(message + "\n"))
	return err
}
