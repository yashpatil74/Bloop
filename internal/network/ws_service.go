package network

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

type Message struct {
	Type     string `json:"type"`
	From     string `json:"from"`
	To       string `json:"to,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FileSize int    `json:"file_size,omitempty"`
	FilePath string `json:"file_path,omitempty"` // ✅ Added field
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (ws *WebSocketService) AddClient(client *websocket.Conn) {
	ws.mu.Lock()
	ws.clients[client] = true
	ws.mu.Unlock()
	log.Println("✅ New WebSocket client connected")
}

func (ws *WebSocketService) RemoveClient(client *websocket.Conn) {
	ws.mu.Lock()
	delete(ws.clients, client)
	ws.mu.Unlock()
	log.Println("❌ WebSocket client disconnected")
}

func (ws *WebSocketService) Broadcast(message interface{}) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Println("❌ Error marshaling WebSocket message:", err)
		return
	}

	for client := range ws.clients {
		err := client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("❌ Error sending message:", err)
			client.Close()
			delete(ws.clients, client)
		}
	}
}

func (ws *WebSocketService) SendToClient(targetIP string, message interface{}) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Println("❌ Error marshaling WebSocket message:", err)
		return
	}

	for client := range ws.clients {
		remoteAddr := client.RemoteAddr().String()
		if strings.Contains(remoteAddr, targetIP) {
			err := client.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("❌ Error sending message to", targetIP, ":", err)
				client.Close()
				delete(ws.clients, client)
			} else {
				log.Println("📤 Sent message to", targetIP)
			}
		}
	}
}
