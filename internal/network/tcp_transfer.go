package network

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

const FILE_TRANSFER_PORT = 9090

// 📤 Send file over TCP
func SendFile(targetIP string, filePath string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", targetIP, FILE_TRANSFER_PORT))
	if err != nil {
		return fmt.Errorf("❌ Could not connect to target: %v", err)
	}
	defer conn.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("❌ Could not open file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(conn, file)
	if err != nil {
		return fmt.Errorf("❌ Error sending file: %v", err)
	}

	log.Println("✅ File sent successfully!")
	return nil
}

// 📥 Receive file
func StartFileReceiver() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", FILE_TRANSFER_PORT))
	if err != nil {
		log.Fatalf("❌ Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("📥 File Receiver listening on port %d", FILE_TRANSFER_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("❌ Error accepting connection:", err)
			continue
		}
		go handleIncomingFile(conn)
	}
}

// 📁 Saves file to ~/Downloads/bloop/
func handleIncomingFile(conn net.Conn) {
	defer conn.Close()

	downloadPath := filepath.Join(os.Getenv("HOME"), "Downloads", "bloop")
	os.MkdirAll(downloadPath, os.ModePerm)

	filePath := filepath.Join(downloadPath, "received_file")
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("❌ Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		log.Println("❌ Error receiving file:", err)
	} else {
		log.Println("✅ File received successfully:", filePath)
	}
}
