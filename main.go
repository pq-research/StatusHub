package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	statushub "github.com/pq-research/StatusHub/src"
)

var usage string = `

    Start a StatusHub server:

    ./StatusHub server

    
    Start a StatusHub client:

    ./StatusHub client

`

func main() {

    logLevel := flag.String("loglevel", "none", "Set log level: debug, none")
    switch *logLevel {
	case "debug":
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Logging level set to DEBUG")
	case "none":
		log.SetOutput(io.Discard)
	default:
		fmt.Println("Invalid log level, using default (none)")
		log.SetOutput(io.Discard)
	}

    if len(os.Args) < 2 {
        fmt.Println(usage)
        return
    }

	role := os.Args[1]
	switch role {
	case "server":
		st := statushub.New()
		st.ListenForStatus()
	case "client":
		c := statushub.NewClient()
		c.Start()
	default:
        fmt.Println(usage)
	}

}
