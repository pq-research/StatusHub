package server

/* server actions
-when to remove client from storage
-who to send status updates to when a client goes [on/off]line

STATE MAINTAINED BY SERVER
-list of actively connected clients
-cached data from storage backing

HOW THE SERVER MANIPULATES THIS STATE

-when client sends a ping:
    -if client hasn't been seen before, then add client to storage with timestamp and notify all other clients about this new client's presence
    -if client has been seen before, then update timestamp in storage

-when client A "doesn't send" a ping (server hasn't heard from particular client in CLIENT_TIMEOUT
    -remove client from cache and storage
    -notify all active clients about the client A's inactiveness

HOW SERVER MAINTAINS THIS STATE
-background patrol process to calculate when each client was last heard from, if the elapsed time is more than CLIENT_TIMEOUT, then remove the client from storage.


-when client goes offline, maintain a "Last Seen" dictionary for each client with their timestamp.

*/

import (
	"fmt"
	"log"
	"net"
	"time"
)

type StatusHub struct {
	clients     map[string]string //time stored in RFC3339 format
	connections map[string]net.Conn
	lastSeen    map[string]string
}

func (sh *StatusHub) checkStatus() {
	now := time.Now()
	for clientId, lastPing := range sh.clients {
		lastPingTimestamp, _ := time.Parse(time.RFC3339, lastPing)
		elapsedTimeSinceLastPing := now.Sub(lastPingTimestamp)
		if elapsedTimeSinceLastPing > CLIENT_TIMEOUT {
			delete(sh.clients, clientId)
			sh.connections[clientId].Close()
			sh.broadcastStatus(clientId, OFFLINE)
		}
	}
}

func (sh *StatusHub) handlePing(conn *net.Conn) {
    pingTime := time.Now().Format(time.RFC3339) // string in RFC3339 format

    //get clientId from the connection
    clientId, err := sh.getClient(conn)
    if err != nil {
        log.Println(err)
        (*conn).Close()
    }

    //if client is already online, just update their ping timestamp 
    //else broadcast to others
	if _, ok := sh.clients[clientId]; ok {
		sh.clients[clientId] = pingTime
	} else {
		sh.clients[clientId] = pingTime
		sh.connections[clientId] = *conn
		sh.broadcastStatus(clientId, ONLINE)
	}
}

func (sh *StatusHub) broadcastStatus(clientOfInterestId string, status Status) {
	for clientId := range sh.clients {
		if clientId != clientOfInterestId {
			statusUpdate := fmt.Sprintf("%s=%s", clientOfInterestId, string(status))
			sh.connections[clientId].Write([]byte(statusUpdate))
		}
	}
}

func (sh *StatusHub) getClient(conn *net.Conn) (string, error) {
    //TODO: read in a while loop?
	buffer := make([]byte, 512)
	n, err := (*conn).Read(buffer)
	if err != nil {
		return "", fmt.Errorf("Error reading from connection: %v", err)
	}
	clientId := string(buffer[:n])
	return clientId, nil
}

// Starts an existing StatusHub instance, listening for 
// client pings and appropriately monitoring status. 
func (sh *StatusHub) ListenAndMonitor() {
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

    //periodically patrol for status
	go func() {
		for {
			sh.checkStatus()
            log.Println("Checking status...")
			time.Sleep(CLIENT_TIMEOUT)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
        
        //TODO: make goroutine?
		sh.handlePing(&conn)
	}
}

// Creates a new StatusHub instance
func New() StatusHub { return StatusHub{} }
