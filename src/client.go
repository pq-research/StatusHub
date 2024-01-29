package statushub

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/ubhattac/when"
)

type StatusHubClient struct {
	id    ClientId
	peers map[ClientId]int64
}

func NewClient() StatusHubClient {
	return StatusHubClient{id: ClientId(rand.Int63()), peers: make(map[ClientId]int64)}
}

func (c *StatusHubClient) Start() {
	log.Println("Started client", c.GetClientId())
	for {
		c.pingStatusHub()
		time.Sleep(5 * time.Second)
	}
}

func (c *StatusHubClient) pingStatusHub() {
	conn, err := net.Dial("tcp", SERVER_ADDRESS)
	defer conn.Close()
	if err != nil {
		log.Println("Failed to create connection with StatusHub at", SERVER_ADDRESS)
	}
	err = binary.Write(conn, binary.LittleEndian, c.id)
	if err != nil {
		log.Println("err:", err)
	}

	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&c.peers); err != nil {
		log.Println("error while decoding:", err)
		return
	}

	peerStatus := c.calculatePeerStatus()
	printPeerStatus(&peerStatus)
}

func (c *StatusHubClient) GetClientId() ClientId { return c.id }

func (c *StatusHubClient) calculatePeerStatus() map[ClientId]string {
	peerStatus := make(map[ClientId]string)
	for peerId, lastPing := range c.peers {
		peerStatus[peerId] = c.getLastSeen(lastPing)
	}
	return peerStatus
}

func (c *StatusHubClient) getLastSeen(lastPing int64) string {
	lastSeen, err := when.When(fmt.Sprintf("%d", lastPing))
	if err != nil {
		log.Println(err)
		return "OFFLINE"
	}
	return lastSeen
}

func printPeerStatus(ps *map[ClientId]string) {
	fmt.Println("PEERS")
	for peerId, status := range *ps {
		fmt.Printf("Peer %d, Last seen %s\n", peerId, status)
	}
}
