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
var serverMsg chan []byte

//// CLIENT ////////
func sendMessage(msg []byte) {
	log.Println(string(msg))
	_, err := conn.Write(msg)
	CheckError(err)

	//	return msgFromServer[:n]
}

func manageMessages() {
	for {
		msgFromServer := make([]byte, 2048)
		n, err := bufio.NewReader(conn).Read(msgFromServer)
		CheckError(err)
		serverMsg <- msgFromServer[:n]
	}
}

func manageWebSocket(ws *websocket.Conn) {
	for {
		serverMessage := <-serverMsg
		_, err := ws.Write(serverMessage)
		CheckError(err)
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	serverMsg = make(chan []byte)
	log.SetFlags(log.Lshortfile)
	ReadFromWebsocket := func(ws *websocket.Conn) {
		go manageWebSocket(ws)
	forLoop:
		for {
			msg := make([]byte, 100)
			n, err := ws.Read(msg)

			if err != nil {
				log.Println(err)
				break forLoop
			}

			switch {
			case strings.Contains(string(msg), "login:"):
				go sendMessage(msg[:n])
				//				log.Println(string(msgFromServer))
			default:
				msgToServer := msg[:n]
				go sendMessage(msgToServer)
				//				log.Println(len(msgFromServer))
				//				ws.Write(msgFromServer)
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

	go manageMessages()

	//	open.Run("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
