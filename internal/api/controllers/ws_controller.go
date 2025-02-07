package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yashpatil74/bloop/internal/network"
)

type WebSocketController struct {
	wsService   *network.WebSocketService
	mdnsService *network.MdnsService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWSController(wsService *network.WebSocketService, mdnsService *network.MdnsService) *WebSocketController {
	return &WebSocketController{
		wsService:   wsService,
		mdnsService: mdnsService,
	}
}

func (ws *WebSocketController) WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	ws.wsService.AddClient(conn)
	defer ws.wsService.RemoveClient(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ WebSocket connection closed:", err)
			break
		}

		var request network.Message
		err = json.Unmarshal(msg, &request)
		if err != nil {
			log.Printf("❌ Invalid JSON: %s", err)
			continue
		}

		switch request.Type {
		case "file_request":
			log.Printf("📩 File Request: %s wants to send %s to %s", request.From, request.FileName, request.To)
			recipient := ws.mdnsService.FindNodeByName(request.To)
			if recipient != nil {
				ws.wsService.SendToClient(recipient.AddrIPv4[0].String(), request)
			}
		case "file_accept":
			log.Printf("✅ File Accepted: %s accepted %s", request.To, request.FileName)
			go network.SendFile(request.From, request.FilePath)
		case "file_decline":
			log.Printf("❌ File Declined: %s declined %s", request.To, request.FileName)
			ws.wsService.Broadcast(request)
		}
	}
}
