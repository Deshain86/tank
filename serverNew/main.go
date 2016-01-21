package main

import (
	"fmt"
	"log"
	"net"
)

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, msg []byte) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		log.Printf("Couldn't send response %v", err)
	}
}

func main() {
	msg := make([]byte, 2048)
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
		_, remoteaddr, err := ser.ReadFromUDP(msg)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, msg)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(ser, remoteaddr, msg)
	}
}
