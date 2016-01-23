package main

import (
	"strings"
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
	"../engine"
)

var conn *net.UDPConn
var ServerAddr *net.UDPAddr
var serverMsg chan []byte
var requestCounter int32 = 0
var waitingRequests map[int32][]byte
var client engine.Client

//// CLIENT ////////
func sendMessage(msg []byte) {
	log.Println(string(msg))
	_, err := conn.Write(msg)
	CheckError(err)
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
		serverMessageString := strings.Split(string(serverMessage), ";")
		log.Println(serverMessageString)
		log.Println(waitingRequests)
		if len(serverMessageString) >= 3 {
			switch serverMessageString[1] {
				case  "LOGIN":
				key, err := strconv.ParseInt(serverMessageString[3], 10, 32)
				log.Println(key)
				CheckError(err)
				if waitingRequests[int32(key)] != nil {
					waitingRequestsArray := strings.Split(string(waitingRequests[int32(key)]),";")
					id, err := strconv.Atoi(serverMessageString[2])
					CheckError(err)
					client.SetId(id)
					client.SetNick(waitingRequestsArray[1])
					delete(waitingRequests, int32(key))
				}
				case  "OK":
				key, err := strconv.ParseInt(serverMessageString[2], 10, 32)
				CheckError(err)
				if waitingRequests[int32(key)] != nil {
					delete(waitingRequests, int32(key))
				}
				default:
				
			}
		}
		
		
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
	client = engine.Client{}
	serverMsg = make(chan []byte)
	waitingRequests = make(map[int32][]byte)
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
			msgId := atomic.AddInt32(&requestCounter, 1)
			messageToSend := msg[:n]
			go func(messageToSend []byte, msgId int32) {
				waitingRequests[msgId] = messageToSend
				for waitingRequests[msgId] != nil {
					toSend := []byte(strconv.FormatInt(int64(msgId), 32) + ";")
					sendMessage(append(toSend, messageToSend...))
					<-time.After(time.Second * 2)
				}
			}(messageToSend, msgId)

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
