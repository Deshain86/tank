package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"../engine"
)

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	_, err := conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		log.Printf("Couldn't send response %v", err)
	}
}

func main() {
	server := engine.NewServer()
	go server.Listen()
	addr := net.UDPAddr{
		Port: 12888,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	log.Println("listening on 12888...")

	for {
		msg := make([]byte, 100)
		n, remoteaddr, err := ser.ReadFromUDP(msg)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, msg)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		tmp := strings.Split(string(msg[:n]), ":")
		switch string(tmp[0]) {
		case "login":
			client := server.NewClient(remoteaddr)
			server.Add(client)
			sendResponse(ser, remoteaddr, strconv.Itoa(client.GetId()))
		default:
			server.ParseResponse(tmp[0], remoteaddr.String())
		}
	}
}
