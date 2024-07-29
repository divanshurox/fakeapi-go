package main

import (
	"FakeAPI/internal/server"
	"log"
)

func main() {
	port := "8080"
	newServer := server.NewServer(port)
	err := newServer.Prestart()
	if err != nil {
		log.Fatal(err)
	}
	done, err := newServer.Start()
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
