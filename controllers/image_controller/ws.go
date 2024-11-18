package image_controller

//func(c *gin.Context){
//	userId := c.Query("userId")
//
//	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		log.Println("WebSocket upgrade failed:", err)
//		return
//	}
//	defer func() {
//		conn.Close()
//		delete(clients, userId)
//		delete(websocket_conn.UserChannels, userId) // Clean up user channel on disconnect
//		log.Println("WebSocket connection closed for user:", userId)
//	}()
//
//	clients[userId] = conn
//	if _, exists := websocket_conn.UserChannels[userId]; !exists {
//		websocket_conn.UserChannels[userId] = make(chan string)
//	}
//
//	go func() {
//		for {
//			processedImageData, ok := <-websocket_conn.UserChannels[userId]
//			if !ok {
//				log.Println("Channel closed for user:", userId)
//				return
//			}
//
//			// Send processed image data to WebSocket
//			msg := Message{
//				Status:    "completed",
//				ImageData: processedImageData,
//			}
//			if err := conn.WriteJSON(msg); err != nil {
//				log.Println("Error sending WebSocket message:", err)
//				return
//			}
//		}
//	}()
//
//	for {
//		_, _, err := conn.ReadMessage()
//		if err != nil {
//			log.Println("WebSocket closed for user:", userId, "Error:", err)
//			return
//		}
//	}
//}
