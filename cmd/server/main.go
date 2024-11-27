package main

import (
	"fmt"
	"log"
	"net"

	"github.com/TheArcus02/chat-app-go/config"
	"github.com/TheArcus02/chat-app-go/internal/services"
)

func main() {
	cfg := config.LoadConfig()
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
	defer listener.Close()

	fmt.Printf("Server running on %s\n", address)

	pool := services.NewConnectionPool()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go pool.HandleConnection(conn)
	}
}
