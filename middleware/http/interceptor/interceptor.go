package interceptor

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var exception = []string{""}

// http 请求拦截器
func HttpInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if strings.HasPrefix(path, "/i/api/conf/") {
			c.Next()
		} else {
			token, valid := judgeToken(c)
			log.Println(token)

			if valid {
				c.AddParam("user_id", token)
				c.Next()
			} else {
				var data interface{}
				c.AbortWithStatusJSON(http.StatusNotAcceptable,
					errObj("101", "need login", data))
			}
		}
	}
}

func judgeToken(c *gin.Context) (string, bool) {
	token := c.PostForm("token")
	if len(token) == 0 {
		token = c.DefaultQuery("token", "")
	}
	if len(token) > 0 {
		return token, true
	}
	return "", false
}

func errObj(code string, msg string, data interface{}) gin.H {
	return gin.H{"code": code, "msg": msg, "data": data}
}
