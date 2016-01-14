package main

import (
	"log"
	"net/http"
	"time"

	"./tank"
)

var refreshrate float32 = 60
var serverRefreshrate float32 = 15

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := tank.NewServer("/entry", (refreshrate / serverRefreshrate))
	go server.Listen()

	ticker := time.NewTicker(time.Second / 15)
	go server.RunInterval(ticker)

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
