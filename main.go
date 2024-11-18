package main

import (
	"fashora-backend/config"
	"fashora-backend/controllers/auth_controller"
	"fashora-backend/controllers/image_controller"
	"fashora-backend/controllers/try_on_controller"
	"fashora-backend/middlewares"
	"fashora-backend/models"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Struct để nhận dữ liệu từ client
type Numbers struct {
	A int `json:"a"`
	B int `json:"b"`
	C int `json:"c"`
}

// Struct để gửi dữ liệu lại cho client
type Result struct {
	Sum int `json:"sum"`
}

type Message struct {
	SessionID string `json:"sessionID"`
	Content   string `json:"content"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Nâng cấp kết nối HTTP lên WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}
	defer conn.Close()

	// Lắng nghe tin nhắn từ client
	for {
		var msg Message
		// Read message from client
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading JSON:", err)
			break
		}
		if err != nil {
			log.Println("Error reading JSON:", err)
			break
		}

		// Tính tổng của a, b và c
		// Log and echo back the message with the sessionID
		log.Printf("Received message with sessionID %s: %s\n", msg.SessionID, msg.Content)
		response := Message{
			SessionID: msg.SessionID,
			Content:   "Server received: " + msg.Content,
		}

		err = conn.WriteJSON(response)
		if err != nil {
			log.Println("Error writing JSON:", err)
			break
		}
	}
}

func main() {
	// Cấu hình REST API server với Gin
	r := gin.Default()

	config.LoadConfig()
	models.ConnectDatabase()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},   // Allow methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow headers
		AllowCredentials: true,                                                // Allow credentials
		MaxAge:           12 * time.Hour,                                      // Max age for preflight requests
	}))

	// Định nghĩa các route cho REST API
	r.POST("/auth/register", auth_controller.Register)
	r.POST("/auth/login", auth_controller.Login)
	r.POST("/auth/check_phone", auth_controller.CheckPhoneNumberExists)
	r.POST("/image/push", image_controller.UploadImage)
	r.GET("/image/get", image_controller.GetImageURL)

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/try_on/push", try_on_controller.UploadImages)
		protected.POST("model/gen", image_controller.Upload3Images)
	}

	// Chạy REST API server trên cổng được cấu hình
	go func() {
		err := r.Run(fmt.Sprintf("%s:%s", config.AppConfig.HostServer, config.AppConfig.PortServer))
		if err != nil {
			log.Fatalf("Failed to start REST API server: %v", err)
		}
	}()

	// Thiết lập server cho WebSocket trên cổng khác (ví dụ: 8081)
	go func() {
		http.HandleFunc("/ws", handleWebSocket)
		log.Println("WebSocket server is running on port 8081...")
		err := http.ListenAndServe(":8082", nil)
		if err != nil {
			log.Fatalf("Failed to start WebSocket server: %v", err)
		}
	}()

	select {}
}
