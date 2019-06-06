package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan Message)           //broadcast channel

var upgrader = websocket.Upgrader{}

type Message struct {
	Email    string `json:"email"`
	username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessage()

	//start the server on local host port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	//upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}
	//make sure we cloase the connection when the function returns
	defer ws.Close()
	clients[ws] = true

	for {
		var msg Message // Message is struct

		//read in a new msg as Json and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// send the newly received meassage to the broadcast
		broadcast <- msg
	}

}

func handleMessage() {
	for {
		//grab the next message from the broadcast channel
		msg := <-broadcast
		// send it out to every client that is currently connected

		for client := range clients {
			err := client.WriteJSON(msg)

			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
