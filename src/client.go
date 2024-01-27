package statushub

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"time"
)

type StatusHubClient struct {
	id    ClientId
	peers map[ClientId]time.Time
}

func NewStatusHubClient() StatusHubClient {
	return StatusHubClient{id: ClientId(rand.Int63())}
}

func (c *StatusHubClient) Start() {
	for {
		c.pingStatusHub()
		time.Sleep(CLIENT_TIMEOUT / 2)
	}
}

func (c *StatusHubClient) pingStatusHub() {
	conn, err := net.Dial("tcp", SERVER_ADDRESS)
	defer conn.Close()
	if err != nil {
		log.Println("Could not dial StatusHub at", SERVER_ADDRESS)
	}
	err = binary.Write(conn, binary.LittleEndian, c.id)
	if err != nil {
		log.Println("err:", err)
	}

	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("Error reading:", err)
		return
	}

	if err := json.Unmarshal(data, &c.peers); err != nil {
		log.Println("Error unmarshaling:", err)
		return
	}

	log.Println("Received map:", c.peers)
}
