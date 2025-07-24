package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Accept connections from any origin for local Wi-Fi use
		return true
	},
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection established")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		log.Printf("[Received] %s\n", string(message))

		// Echo back the same message for this simple test
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}

	log.Println("WebSocket connection closed")
}
