package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vegacrypto/vegax_backend/model"
	database "github.com/vegacrypto/vegax_backend/system"
)

func retObj(code, msg string, data interface{}) gin.H {
	return gin.H{"code": code, "msg": msg, "data": data}
}

func HandleLogin(c *gin.Context) {
	// p := c.Params
	// tokenid, get := p.Get("user_id")
	// log.Println(tokenid, get)

	email := c.PostForm("email")
	passwd := c.PostForm("passwd")

	var result []model.User
	db := database.GetDb()
	db.Model(&model.User{}).Where("email = ?", email).Find(&result)

	var data interface{}
	if len(result) == 0 {
		c.JSON(http.StatusOK, retObj("102", "user not found", data))
	} else if len(result) > 1 {
		c.JSON(http.StatusOK, retObj("103", "user config error", data))
	} else {
		mu := result[0]

		if passwd != mu.Password {
			c.JSON(http.StatusOK, retObj("104", "user password error", data))
		} else {
			_map := map[string]interface{}{}
			_map["token"] = mu.Id
			_map["avarta"] = mu.Avarta
			_map["nick"] = mu.NickName
			c.JSON(http.StatusOK, retObj("100", "success", _map))
		}
	}
}
