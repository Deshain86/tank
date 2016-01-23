package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

var conn *net.UDPConn
var ServerAddr *net.UDPAddr

//// CLIENT ////////
func sendMessage(msg []byte) []byte {
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
			msg := make([]byte, 100)
			_, err := ws.Read(msg)
			if err != nil {
				log.Println(err)
				break forLoop
			}

			switch {
			case strings.Contains(string(msg), "nick:"):
				nick := "login:" + string(msg)[5:]
				msgFromServer := sendMessage([]byte(nick))
				log.Println(string(msgFromServer))
			default:
				msgToServer := msg[0:10]
				msgFromServer := sendMessage(msgToServer)
				ws.Write(msgFromServer)
			}

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
