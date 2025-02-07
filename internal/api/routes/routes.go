package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yashpatil74/bloop/internal/api/controllers"
)

func RegisterRoutes(router *gin.Engine, wsController *controllers.WebSocketController) {
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/ws", wsController.WebSocketHandler)
	}
}
