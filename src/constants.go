package statushub

import "time"

/* TEMPORAL */

const (
	//CLIENT_TIMEOUT = 3 * time.Minute
	CLIENT_TIMEOUT = 3 * time.Second
)

/* STATE */

type Status string

const (
	ONLINE  Status = "ONLINE"
	OFFLINE Status = "OFFLINE"
)

const SERVER_HOST = "localhost"
const SERVER_PORT = "9000"
const SERVER_ADDRESS = SERVER_HOST + ":" + SERVER_PORT
