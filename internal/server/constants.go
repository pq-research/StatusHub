package server

import "time"

/* TEMPORAL */

const (
	CLIENT_TIMEOUT = 3 * time.Minute
)

/* STATE */

type Status bool

const (
	ONLINE  Status = true
	OFFLINE Status = false
)
