package server

import (
	"fmt"
	"net"
	"time"
)

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

type StatusHub struct {
	clients     map[string]string
	connections map[string]net.Conn
}

func (st *StatusHub) checkPresence() {
	now := time.Now()
	for clientId, lastPing := range st.clients {
		lastPingTimestamp, _ := time.Parse(time.RFC3339, lastPing)
		elapsedTimeSinceLastPing := now.Sub(lastPingTimestamp)
		if elapsedTimeSinceLastPing > CLIENT_TIMEOUT {
			delete(st.clients, clientId)
			st.connections[clientId].Close()
			st.broadcastStatus(clientId, OFFLINE)
		}
	}
}

func (st *StatusHub) handlePing(clientId string, pingTime string) {
	if _, ok := st.clients[clientId]; ok {
		st.clients[clientId] = pingTime
	} else {
		st.clients[clientId] = pingTime
		st.broadcastStatus(clientId, ONLINE)
	}
}

func (st *StatusHub) broadcastStatus(clientOfInterestId string, status Status) {
	for clientId, _ := range st.clients {
		if clientId != clientOfInterestId {
			statusUpdate := fmt.Sprintf("%s=%s", clientOfInterestId, status)
			st.connections[clientId].Write([]byte(statusUpdate))
		}
	}
}
