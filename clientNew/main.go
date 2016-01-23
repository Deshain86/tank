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
		n,err :=bufio.NewReader(conn).Read(msgFromServer)
		CheckError(err)
		serverMsg <- msgFromServer[:n]
	}
}

func manageWebSocket(ws *websocket.Conn){
	for {
			serverMessage := <-serverMsg
			log.Println(string(serverMessage))
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
			default:
				msgToServer := msg[:n]
				go sendMessage(msgToServer)
			}
		}
		
	}
	http.Handle("/echo", websocket.Handler(ReadFromWebsocket))

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	ServerAddr, err := net.ResolveUDPAddr("udp", "89.72.59.44:8081")
	CheckError(err)

	ClientAddr, err := net.ResolveUDPAddr("udp", ":0")
	CheckError(err)

	conn, err = net.DialUDP("udp", ClientAddr, ServerAddr)
	CheckError(err)
	defer conn.Close()
	
	go manageMessages()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
