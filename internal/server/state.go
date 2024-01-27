package server

import "net"

type StatusHub struct {
	clients     map[string]string
	connections map[string]net.Conn
    lastSeen map[string]string
}

