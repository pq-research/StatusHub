package server

import "time"

/* TEMPORAL */

const (
	CLIENT_TIMEOUT = 3 * time.Minute
)

/* STATE */

type Status string

const (
	ONLINE  Status = "ONLINE"
	OFFLINE Status = "OFFLINE"
)
