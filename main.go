package main

import (
	"log"
	"net/http"

	"./tank"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := tank.NewServer("/entry")
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
