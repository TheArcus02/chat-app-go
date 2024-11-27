package models

import "net"

type User struct {
	Name   string
	Conn   net.Conn
	Online bool
}
