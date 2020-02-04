package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func readPackets(s socketio.Conn) {

	p := make([]byte, 1024)

	addr := net.UDPAddr{
		Port: 41181,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	prevData := ""
	for {
		n, remoteaddr, err := ser.ReadFromUDP(p)
		msg := strings.Split(string(p[:n]), ",")
		fmt.Printf("Prev:  %s %s \n", prevData, msg[1])
		if prevData != msg[1] {
			fmt.Printf("Read a message from %v %s \n", remoteaddr, p[:n])
			prevData = msg[1]
			s.Emit("field", string(p[:n]))

		}
		if err != nil {

		}

	}
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		readPackets(s)

		return nil
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)

		s.Emit("reply", "New Packet "+msg)
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset1")))
	log.Println("Serving at localhost:4000...")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
