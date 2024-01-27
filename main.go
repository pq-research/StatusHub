package main

import "github.com/pq-research/StatusHub/src"

func main() {
    sh := statushub.NewStatusHub()
    sh.ListenForStatus()
}
