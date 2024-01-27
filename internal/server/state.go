package server

import "net"

var client_ map[net.IPAddr]struct{}
