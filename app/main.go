// Package wallet
package main

import (
	"log"

	"wallet/app/queue"
	"wallet/app/server"
)

func main() {
	// Init nsq.
	nsq, err := queue.NewNSQ()
	if err != nil {
		log.Fatal(err)
	}
	defer nsq.Stop()

	// Init web server.
	s := server.New(nsq)
	s.SetupMiddleware()
	s.SetupApp()

	if err = s.Start(); err != nil {
		log.Fatal(err)
	}
}
