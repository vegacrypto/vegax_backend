package interceptor

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
			userid, valid := judgeToken(c)
			if valid {
				c.AddParam("user_id", strconv.FormatUint(userid, 10))
				c.Next()
			} else {
				var data interface{}
				c.AbortWithStatusJSON(http.StatusNotAcceptable,
					errObj("101", "need login", data))
			}
		}
	}
}

func judgeToken(c *gin.Context) (uint64, bool) {
	token := c.PostForm("token")
	if len(token) == 0 {
		token = c.DefaultQuery("token", "")
	}
	if len(token) > 0 {
		e := tool.DesToken{}
		realToken, success := e.Decrypt(token)
		if !success {
			return 0, false
		}
		arr := strings.Split(realToken, ",")
		if len(arr) == 2 {
			userid, err := strconv.Atoi(arr[0])
			if err == nil {
				ts, err := strconv.Atoi(arr[1])
				if err == nil && time.Now().Unix()-int64(ts) < 30*24*60*60 {
					return uint64(userid), true
				}
			}
		}
	}
	return 0, false
}

func errObj(code string, msg string, data interface{}) gin.H {
	return gin.H{"code": code, "msg": msg, "data": data}
}
