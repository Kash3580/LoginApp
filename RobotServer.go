package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type TodoPageData struct {
	PageTitle string
}

var data []TodoPageData

func main() {

	router := mux.NewRouter()

	readsocket()
	router.HandleFunc("/index", indexHandler).Methods("GET")
	router.HandleFunc("/", ReadData).Methods("GET")
	log.Println("Server started on :4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func printSlice(s []TodoPageData) {
	fmt.Printf("len=%d  %v\n", len(s), s)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index called")
	t, _ := template.ParseFiles("view/welcome.html")
	//t:=template.Must(template.New("html-tmpl").ParseFiles("view/welcome.html"))
	t.Execute(w, "nodata")

}
func ReadData(w http.ResponseWriter, r *http.Request) {
	//var remoteip *net.UDPAddr

	log.Println("started reading")
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
	for {
		n, remoteaddr, err := ser.ReadFromUDP(p)

		fmt.Printf("Read a message from %v %s \n", remoteaddr, p[:n])
		//	fmt.Printf("%s", p)

		w.Write([]byte(p[:n]))
		///	sendResponse("\n this is rasing alarm: " + remoteip.String())
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		//sendResponse("\n this is rasing alarm: " + remoteaddr.String())
	}

}

func readsocket() {

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
	for {
		n, remoteaddr, err := ser.ReadFromUDP(p)

		//	fmt.Printf("Read a message from %v %s \n", remoteaddr, p[:n])
		raw := TodoPageData{
			PageTitle: string([]byte(p[:n])),
		}
		data = append(data, raw)
		sendResponse(raw)

		if err != nil {
			fmt.Printf("Some error  %v %v", err, remoteaddr)
			continue
		}

	}
}

func sendResponse(Messages TodoPageData) {

	conn, err := net.Dial("udp", "127.0.0.1:4000")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	//for i, message := range Messages {
	//	fmt.Println("sending on 4001: %d", i)

	fmt.Printf("%q\n", bytes.Split([]byte(Messages.PageTitle), []byte(","))[2])
	conn.Write([]byte(string(Messages.PageTitle)))

	//}

}
