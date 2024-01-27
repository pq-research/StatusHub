package statushub

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"time"
)

type (
	ClientId int64

	StatusHub struct {
		clients map[ClientId]time.Time
	}
)

func NewStatusHub() StatusHub { return StatusHub{} }

func (sh *StatusHub) ListenForStatus() {
	listener, err := net.Listen("tcp", SERVER_ADDRESS)
	if err != nil {
		panic(err)
	}
	log.Println("Listening at", SERVER_ADDRESS)

	for {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			log.Println("Error while accepting connection")
			continue
		}

		go sh.deliverStatus(&conn, time.Now())
	}
}

func (sh *StatusHub) deliverStatus(conn *net.Conn, timeStamp time.Time) {
	var cid ClientId
	err := binary.Read(*conn, binary.LittleEndian, &cid)
	if err != nil {
		log.Println("Error while reading from the client")
	}
	sh.clients[cid] = timeStamp
	payload, _ := json.Marshal(sh.clients)
	(*conn).Write(payload)
	(*conn).Close()
}
