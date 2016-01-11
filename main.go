package main

import (
	"log"
	"net/http"
	"time"

	"./tank"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := tank.NewServer("/entry")
	go server.Listen()
	ticker := time.NewTicker(time.Second / 30)
	go server.RunInterval(ticker)

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
