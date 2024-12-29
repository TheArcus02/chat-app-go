# Go Chat Server

This repository contains the implementation of a real-time chat server written in Go. The server handles user connections, message broadcasting, and user management efficiently, leveraging Go's concurrency model.

---

## Features

- Real-time messaging using TCP sockets.
- User connection and disconnection handling.
- Broadcast messages to all connected clients.
- Lightweight and efficient server design.
- Docker support for containerized deployment.

---

## Requirements

- Go 1.18 or later.
- Docker (optional, for containerized deployment).

---

## Installation and Usage

### Clone the Repository

```bash
git clone https://github.com/TheArcus02/chat-app-go.git
cd chat-app-go
```

### Build the Server

#### Locally

1. Navigate to the directory containing the server code (e.g., `cmd/server`).
2. Run the server:

      ```bash
      go run cmd/server/main.go
      ```

#### With Docker

1. Build the Docker image:

      ```bash
      docker build -t chat-app-go .
      ```

2. Run the container:

      ```bash
      docker run -p 8080:8080 chat-app-go
      ```

3. The server will be accessible at `localhost:8080`.

---

## Configuration

- The server's configuration (e.g., host, port) can be adjusted via environment variables or configuration files.

---

## Directory Structure

```
chat-app-go/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point of the application
├── config/
│   └── config.go                   # Configuration setup (e.g., server port)
├── internal/
│   ├── handlers/
│   │   ├── connection.go           # Logic for handling individual connections
│   │   └── chat_handler.go         # Chat message handling
│   ├── models/
│   │   └── user.go                 # User model definition
│   ├── services/
│   │   ├── connection_pool.go      # Managing connected users
│   └── utils/
│       └── logger.go               # Logging utility
├── pkg/
│   └── protocol/
│       ├── message.go              # Structs and utilities for message parsing
│       └── constants.go            # Message types or protocol constants
└── go.mod                          # Go module file
```

---

## Troubleshooting

### Common Issues

- **Cannot connect to server:** Ensure the server is running and accessible at the specified host and port.
- **Docker networking issues:** Use `docker inspect` to check the container's IP address and ensure proper port mapping.

### Logs

- Logs are printed to the console and can be redirected to a file if needed.

---

## Contributions

Contributions are welcome! Feel free to open issues or submit pull requests to improve the server or add new features.
