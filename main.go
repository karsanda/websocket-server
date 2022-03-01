package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/connect" {
		http.Error(rw, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(rw, "Method is not supported.", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Print("Upgrade failed:", err)
		return
	}

	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error when read:", err)
			break
		}

		log.Printf("Received: %s", message)

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error when write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/connect", handleWebsocket)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
