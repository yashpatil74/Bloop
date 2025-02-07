package services

import (
	"encoding/json"

	"github.com/grandcat/zeroconf"
	"github.com/yashpatil74/bloop/internal/network"
)

type WSWebService struct {
	wsService   *network.WebSocketService
	mdnsService *network.MdnsService
}

func NewWSWebService(wsService *network.WebSocketService, mdnsService *network.MdnsService) *WSWebService {
	return &WSWebService{
		wsService:   wsService,
		mdnsService: mdnsService,
	}
}

func (ws *WSWebService) BroadcastNodes() {
	nodes := ws.mdnsService.GetNodes()
	nodesJSON := formatNodes(nodes)

	ws.wsService.Broadcast(nodesJSON)
}

func formatNodes(nodes map[string]*zeroconf.ServiceEntry) []byte {
	formattedNodes := make([]map[string]interface{}, 0)
	for _, entry := range nodes {
		formattedNodes = append(formattedNodes, map[string]interface{}{
			"name": entry.Instance,
			"ip":   entry.AddrIPv4,
			"port": entry.Port,
		})
	}
	data, err := json.Marshal(formattedNodes)
	if err != nil {
		return []byte("[]")
	}
	return data
}
