package main

import (
	"log"
	"net/http"

	"github.com/skratchdot/open-golang/open"
)

func main() {
	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	open.Run("http://localhost:12446")
	log.Fatal(http.ListenAndServe(":12446", nil))
}
