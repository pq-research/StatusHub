package statushub

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type (
	ClientId int64

	StatusHub struct {
		clients map[ClientId]int64
	}
)

func New() StatusHub { return StatusHub{clients: make(map[ClientId]int64)} }

func (sh *StatusHub) ListenForStatus() {
	listener, err := net.Listen("tcp", SERVER_ADDRESS)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening at", SERVER_ADDRESS)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error while accepting connection")
			continue
		}
		log.Println("Accepted connection from", conn.RemoteAddr())

		go sh.deliverStatus(&conn, time.Now().Unix())
	}
}

func (sh *StatusHub) deliverStatus(conn *net.Conn, timeStamp int64) {
	defer (*conn).Close()
	var cid ClientId
	err := binary.Read(*conn, binary.LittleEndian, &cid)
	if err != nil {
		log.Println("Error while reading from the client")
	}
	sh.clients[cid] = timeStamp
	log.Println(sh.clients)
	payload, _ := json.Marshal(sh.clients)
	(*conn).Write(payload)
}
