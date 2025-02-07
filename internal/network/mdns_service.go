package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/yashpatil74/bloop/internal/utils"
)

type MdnsService struct {
	server    *zeroconf.Server
	entries   chan *zeroconf.ServiceEntry
	nodes     map[string]*zeroconf.ServiceEntry
	mu        sync.Mutex
	localNode string
}

func NewMdnsService() *MdnsService {
	randomName := "bloop-" + utils.GenerateRandomString(5)
	return &MdnsService{
		entries:   make(chan *zeroconf.ServiceEntry),
		nodes:     make(map[string]*zeroconf.ServiceEntry),
		localNode: randomName,
	}
}

func (m *MdnsService) Advertise(port int) error {
	txtRecords := []string{"nodeName=" + m.localNode}
	server, err := zeroconf.Register(m.localNode, "_bloop._tcp", "local.", port, txtRecords, nil)
	if err != nil {
		return err
	}
	m.server = server
	log.Printf("📡 Advertising node: %s on port %d", m.localNode, port)
	return nil
}

func (m *MdnsService) Discover(serviceType string) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalf("Failed to initialize resolver: %s", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		for entry := range entries {
			if entry.Instance == m.localNode {
				continue
			}

			m.mu.Lock()
			m.nodes[entry.Instance] = entry
			m.mu.Unlock()

			log.Printf("🔎 Discovered: %s at %v:%d", entry.Instance, entry.AddrIPv4, entry.Port)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err = resolver.Browse(ctx, serviceType, "local.", entries)
	if err != nil {
		log.Fatalf("Failed to browse: %s", err)
	}
}

func (m *MdnsService) FindNodeByName(nodeName string) *zeroconf.ServiceEntry {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, entry := range m.nodes {
		if entry.Instance == nodeName {
			return entry
		}
	}
	return nil
}

func (m *MdnsService) ValidateNodes() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for instance, entry := range m.nodes {
		address := net.JoinHostPort(entry.AddrIPv4[0].String(), fmt.Sprintf("%d", entry.Port))
		if !isNodeAlive(address) {
			log.Printf("❌ Node offline: %s", instance)
			delete(m.nodes, instance)
		}
	}
}

func isNodeAlive(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		log.Printf("❌ Node unreachable: %s (%v)", address, err)
		return false
	}
	conn.Close()
	return true
}

func (m *MdnsService) GetNodes() map[string]*zeroconf.ServiceEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.nodes
}

func (m *MdnsService) Stop() {
	if m.server != nil {
		m.server.Shutdown()
	}
	close(m.entries)
	log.Println("🛑 Stopped mDNS service")
}
