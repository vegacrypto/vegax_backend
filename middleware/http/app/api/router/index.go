package router

import (
	"github.com/gin-gonic/gin"
	c "github.com/vegacrypto/vegax_backend/middleware/http/app/api/controller"
)

func Routers(e *gin.RouterGroup) {
	userGroup := e.Group("/api")
	userGroup.GET("/events_by_contract", c.GetEventsByContract)
}
