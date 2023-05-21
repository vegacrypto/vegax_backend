package interceptor

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vegacrypto/vegax_backend/tool"
)

// http 请求拦截器
func HttpInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if strings.HasPrefix(path, "/i/api/conf/") {
			c.Next()
		} else {
			token, valid := judgeToken(c)
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
		e := tool.DesToken{}
		realToken, success := e.Decrypt(token)
		if !success {
			return "", false
		}
		return realToken, true
	}
	return "", false
}

func errObj(code string, msg string, data interface{}) gin.H {
	return gin.H{"code": code, "msg": msg, "data": data}
}
