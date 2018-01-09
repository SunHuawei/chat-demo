package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients []*websocket.Conn

func closeConn(conn *websocket.Conn) {
	conn.Close()

	for i, c := range clients {
		if c == conn {
			clients = append(clients[:i], clients[i + 1:]...)
		}
	}

}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	defer func() {
		closeConn(conn)
		log.Println("Clients:", clients)
	}()

	if err != nil {
		log.Println(err)
		return
	}

	clients = append(clients, conn)
	log.Println("Clients:", clients)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		broadcast(messageType, p)
	}
}

func broadcast(messageType int, message []byte) {
	for _, c := range clients {
		if err := c.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "chat.html")
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	err := http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
