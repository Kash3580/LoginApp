package main

import (
	"fmt"
	"net"
)

func main() {

	p := make([]byte, 1024)

	addr := net.UDPAddr{
		Port: 4002,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		n, remoteaddr, err := ser.ReadFromUDP(p)

		fmt.Printf("Read a message from %v %s \n", remoteaddr, p[:n])

		if err != nil {
		}

	}
}
