package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

var conn *net.UDPConn
var ServerAddr *net.UDPAddr

//// CLIENT ////////
func repeatMessage(msg []byte) []byte {
	log.Println(string(msg))
	_, err := conn.Write(msg)
	CheckError(err)
	msgFromServer := make([]byte, 2048)
	bufio.NewReader(conn).Read(msgFromServer)
	return msgFromServer
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	ReadFromWebsocket := func(ws *websocket.Conn) {
	forLoop:
		for {
			msg := make([]byte, 10)
			_, err := ws.Read(msg)
			if err != nil {
				break forLoop
			}
			msg = repeatMessage(msg)

			ws.Write(msg)
		}
	}
	http.Handle("/echo", websocket.Handler(ReadFromWebsocket))

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:12888")
	CheckError(err)

	ClientAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	conn, err = net.DialUDP("udp", ClientAddr, ServerAddr)
	CheckError(err)
	defer conn.Close()

	//	open.Run("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
