package main

import (
	"log"
	"net/http"
	"time"

	"./tank"
)

var refreshrate int = 30
var serverRefreshrate int = 10

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := tank.NewServer("/entry", (float32(refreshrate) / float32(serverRefreshrate)))
	go server.Listen()

	ticker := time.NewTicker(time.Second / time.Duration(serverRefreshrate))
	go server.RunInterval(ticker)

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
