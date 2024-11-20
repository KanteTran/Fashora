package websocket_conn

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var Semaphore = make(chan struct{}, 2)

type Message struct {
	Status    string `json:"status"`
	ImageData string `json:"imageData,omitempty"`
}

var Clients = make(map[string]*websocket.Conn)
var UserChannels = make(map[string]chan string)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Cho phép tất cả các nguồn, điều chỉnh lại khi vào production
	},
}

func AddUser(userID string, conn *websocket.Conn) {
	Clients[userID] = conn
	UserChannels[userID] = make(chan string)
}
func sendMessageToUser(userID string, message string) {
	if ch, ok := UserChannels[userID]; ok {
		ch <- message
	}
}

func listenToUserChannel(userID string) {
	conn := Clients[userID]
	ch := UserChannels[userID]

	go func() {
		for message := range ch {
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Error writing to WebSocket:", err)
				return
			}
		}
	}()
}
func removeUser(userID string) {
	if conn, ok := Clients[userID]; ok {
		conn.Close()
		delete(Clients, userID)
	}
	if ch, ok := UserChannels[userID]; ok {
		close(ch)
		delete(UserChannels, userID)
	}
}
