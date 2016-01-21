package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/skratchdot/open-golang/open"
)

var conn *net.UDPConn
var ServerAddr *net.UDPAddr

func repeatMessage(msg []byte) []byte {
	log.Println(string(msg))
	_, err := conn.Write(msg)
	CheckError(err)
	bufio.NewReader(conn).Read(msg)
	return msg
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func EchoServer(ws *websocket.Conn) {
	for {
		msg := make([]byte, 2048)
		ws.Read(msg)
		log.Println(string(msg))
		msg = repeatMessage(msg)
		ws.Write(msg)
	}

	//	io.Copy(ws, ws)
}

func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:12888")
	CheckError(err)

	ClientAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	conn, err = net.DialUDP("udp", ClientAddr, ServerAddr)
	CheckError(err)
	defer conn.Close()

	open.Run("http://localhost:8080")
	log.Println("joł")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
