package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"../engine"
)

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	_, err := conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		log.Printf("Couldn't send response %v", err)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	ser, err := net.ListenUDP("udp", ServerAddr)

	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	server := engine.NewServer(ser)
	go server.Listen()
	log.Println("listening on 12888...")
	go func() {
		for {
			server.SendAll()
			log.Print("sendAll")
			<-time.After(5 * time.Second)
		}
	}()
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
		default:
			server.ParseResponse(tmp[0], remoteaddr)
		}
	}
}
