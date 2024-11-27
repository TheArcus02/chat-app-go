package services

import (
	"fmt"
	"net"
	"sync"

	"github.com/TheArcus02/chat-app-go/internal/models"
)

type ConnectionPool struct {
	mu    sync.Mutex
	users map[string]*models.User
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		users: make(map[string]*models.User),
	}
}

func (pool *ConnectionPool) HandleConnection(conn net.Conn) {
	defer conn.Close()

	var username string

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading username: %v\n", err)
		return
	}

	username = string(buffer[:n])
	user := &models.User{Name: username, Conn: conn, Online: true}

	pool.mu.Lock()
	pool.users[username] = user
	pool.mu.Unlock()

	fmt.Printf("User connected: %s\n", username)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("User %s disconnected.\n", username)
			break
		}

		message := string(buffer[:n])
		fmt.Printf("Message from %s: %s\n", username, message)

		conn.Write([]byte("Echo: " + message))
	}

	pool.mu.Lock()
	delete(pool.users, username)
	pool.mu.Unlock()
}
