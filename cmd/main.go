package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yashpatil74/bloop/internal/api/controllers"
	"github.com/yashpatil74/bloop/internal/api/routes"
	"github.com/yashpatil74/bloop/internal/api/services"
	"github.com/yashpatil74/bloop/internal/network"
	"github.com/yashpatil74/bloop/internal/utils"
)

const HTTP_PORT = 8780

func main() {
	log.Println("🚀 Starting Bloop")
	defer log.Println("🛑 Shutting down Bloop")

	//  Network Services
	mdnsService := network.NewMdnsService()
	wsService := network.NewWebSocketService()

	// Web Services
	wsWebService := services.NewWSWebService(wsService, mdnsService)

	// Controllers
	wsController := controllers.NewWSController(wsService, mdnsService)

	err := mdnsService.Advertise(HTTP_PORT)
	if err != nil {
		log.Fatalf("❌ Failed to start mDNS: %v", err)
	}
	defer mdnsService.Stop()

	go func() {
		for {
			log.Println("🔍 Scanning for nearby devices...")
			mdnsService.Discover("_bloop._tcp")
			mdnsService.ValidateNodes()
			wsWebService.BroadcastNodes()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	router := gin.Default()
	routes.RegisterRoutes(router, wsController)
	router.StaticFS("/app", http.Dir("../web/next/out"))

	go func() {
		log.Printf("🚀 Bloop UI available at: http://localhost:%d\n", HTTP_PORT)
		if err := router.Run(fmt.Sprintf(":%d", HTTP_PORT)); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)
	utils.OpenBrowser(fmt.Sprintf("http://localhost:%d/app", HTTP_PORT))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
