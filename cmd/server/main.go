package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheArcus02/chat-app-go/config"
	"github.com/TheArcus02/chat-app-go/internal/handlers"
	"github.com/TheArcus02/chat-app-go/internal/services"
	"github.com/TheArcus02/chat-app-go/internal/utils"
)

func main() {
	cfg := config.LoadConfig()
	logger := utils.NewLogger()

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
	
	defer listener.Close()

	logger.Infof("Server running on %s", address)

	connectionPool := services.NewConnectionPool(logger)
	go connectionPool.Run()

	go handleShutdown(listener, connectionPool, logger)

	chatHandler := handlers.NewChatHandler(logger, connectionPool)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("Failed to accept connection: %v", err)
			continue
		}
		logger.Infof("New connection from %s", conn.RemoteAddr())
		connectionHandler := handlers.NewConnectionHandler(logger, &conn, connectionPool, chatHandler)
		go connectionHandler.Handle()
	}
}

func handleShutdown(listener net.Listener, pool *services.ConnectionPool, logger *utils.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Infof("Shutting down server...")

	pool.Shutdown()  
	listener.Close() 
	os.Exit(0)
}
