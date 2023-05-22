package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vegacrypto/vegax_backend/model"
	database "github.com/vegacrypto/vegax_backend/system"
)

func HandleChatInput(c *gin.Context) {
	p := c.Params
	user_id, get := p.Get("user_id")
	log.Println(user_id, get)
	userId, _ := strconv.ParseUint(user_id, 10, 64)

	var data interface{}

	prompt := strings.Trim(c.PostForm("prompt"), " ")
	if len(prompt) == 0 {
		c.JSON(http.StatusOK, retObj("201", "empty input", data))
		return
	}
	taskCode := strings.Trim(c.PostForm("task_code"), " ")

	db := database.GetDb()
	chat := &model.Chat{
		BaseModel: model.BaseModel{
			AddTime:    time.Now(),
			UpdateTime: time.Now(),
		},
		UserId:   userId,
		Content:  prompt,
		Status:   0,
		TaskCode: taskCode,
	}
	db.Model(&model.Chat{}).Save(chat)

	//这里需要启动多线程去交互LLM
	c.JSON(http.StatusOK, retObj("100", "success", chat))
}

func HandleChatHistory(c *gin.Context) {
	p := c.Params
	user_id, get := p.Get("user_id")

	userId, _ := strconv.ParseUint(user_id, 10, 64)

	log.Println(user_id, get, userId)

	var result []model.Chat

	db := database.GetDb()
	db.Model(&model.Chat{}).Where("user_id = ?", userId).Order("add_time desc").Find(&result)

	//这里需要启动多线程去交互LLM
	c.JSON(http.StatusOK, retObj("100", "success", result))
}
