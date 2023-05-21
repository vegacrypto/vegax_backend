package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vegacrypto/vegax_backend/model"
	database "github.com/vegacrypto/vegax_backend/system"
	"github.com/vegacrypto/vegax_backend/tool"
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
			e := tool.DesToken{}

			token, success := e.Encrypt(strconv.FormatUint(mu.Id, 10) + "," + strconv.FormatInt(time.Now().Unix(), 10))
			if !success {
				c.JSON(http.StatusOK, retObj("105", "user token build error", data))
			} else {
				_map := map[string]interface{}{}
				_map["token"] = token
				_map["avarta"] = mu.Avarta
				_map["nick"] = mu.NickName
				c.JSON(http.StatusOK, retObj("100", "success", _map))
			}
		}
	}
}

func HandleRegister(c *gin.Context) {
	email := c.PostForm("email")
	passwd := c.PostForm("passwd")
	code := c.PostForm("code")

	var result []model.User
	db := database.GetDb()
	db.Model(&model.User{}).Where("email = ?", email).Find(&result)

	var data interface{}

	if len(result) > 0 {
		c.JSON(http.StatusOK, retObj("106", "user exists", data))
		return
	}

	if code != "1234" {
		c.JSON(http.StatusOK, retObj("107", "email code error", data))
		return
	}
	mou := &model.User{
		BaseModel: model.BaseModel{
			AddTime:    time.Now(),
			UpdateTime: time.Now(),
		},
		Email:    email,
		Password: passwd,
		Flag:     0,
	}
	db.Save(mou)

	c.JSON(http.StatusOK, retObj("100", "success", data))
}
