package router

import (
	"github.com/gin-gonic/gin"
	c "github.com/vegacrypto/vegax_backend/middleware/http/app/api/controller"
)

func Routers(e *gin.RouterGroup) {
	userGroup := e.Group("/api")
	userGroup.POST("/conf/login", c.HandleLogin)
	userGroup.POST("/conf/reg", c.HandleRegister)
	userGroup.GET("/conf/tasks", c.HandleTaskSupp)

	userGroup.POST("/chat/save", c.HandleChatInput)
	userGroup.POST("/chat/history", c.HandleChatHistory)
	userGroup.POST("/chat/bychat", c.HandleChatsById)
	userGroup.POST("/chat/setop", c.HandleChatOp)
}
