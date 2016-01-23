package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"../engine"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

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
	log.Println("listening on 8081...")
	go func() {
		for {
			server.SendAll()
			// log.Print("sendAll")
			<-time.After(time.Second / 20)
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
		tmp := strings.Split(string(msg[:n]), ";")
		if len(tmp) > 1 {
			switch string(tmp[1]) {
			case "login":
				if len(tmp) == 3 {
					client, clientId := server.NewClient(remoteaddr, tmp[2], tmp[0])
					if clientId == 0 {
						server.Add(client, tmp[0])
					}
				}
			default:
				server.ParseResponse(tmp[0], tmp[1], remoteaddr)
			}
		}
	}
}
