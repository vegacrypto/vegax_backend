package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	api "github.com/vegacrypto/vegax_backend/middleware/http/app/api/router"
	"github.com/vegacrypto/vegax_backend/middleware/http/interceptor"
)

type Option func(*gin.RouterGroup)

var options = []Option{}

func Include(opts ...Option) {
	options = append(options, opts...)
}

func Init() *gin.Engine {
	Include(api.Routers)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/index", helloHandler) //Default welcome api

	staticTplPath := "./middleware/http/templates/**/*"

	if configFilePathFromEnv := os.Getenv("DALINK_TPL_PATH"); configFilePathFromEnv != "" {
		staticTplPath = configFilePathFromEnv + "/**/*"
	}

	r.LoadHTMLGlob(staticTplPath)

	apiGroup := r.Group("/i", interceptor.HttpInterceptor())
	for _, opt := range options {
		opt(apiGroup)
	}
	r.Run(":18080")
	return r
}

func helloHandler(c *gin.Context) {
	log.Println("hello")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello dalink",
	})
}
