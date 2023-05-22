package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vegacrypto/vegax_backend/model"
	database "github.com/vegacrypto/vegax_backend/system"
)

func HandleChatsById(c *gin.Context) {
	p := c.Params
	user_id, get := p.Get("user_id")
	log.Println(user_id, get)

	userId, _ := strconv.ParseUint(user_id, 10, 64)

	chat_id := strings.Trim(c.PostForm("chat_id"), " ")
	chatId, _ := strconv.ParseUint(chat_id, 10, 64)

	var data interface{}

	if chatId == 0 {
		c.JSON(http.StatusOK, retObj("203", "parent chat id empty", data))
		return
	}

	var chats []model.Chat

	db := database.GetDb()
	db.Model(&model.Chat{}).Where("user_id = ? and chat_id = ?", userId, chatId).Find(&chats)

	c.JSON(http.StatusOK, retObj("100", "success", chats))
}

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

	var suppLLMs []model.SysConf
	db.Model(&model.SysConf{}).Where("conf_type = ?", "task_"+taskCode).Find(&suppLLMs)

	if len(suppLLMs) == 0 {
		c.JSON(http.StatusOK, retObj("202", "unsupport task type", data))
		return
	}
	chat := &model.Chat{
		BaseModel: model.BaseModel{
			AddTime:    time.Now(),
			UpdateTime: time.Now(),
		},
		UserId:   userId,
		Content:  prompt,
		Status:   0,
		Expect:   len(suppLLMs),
		TaskCode: taskCode,
	}
	db.Model(&model.Chat{}).Save(chat)

	c.JSON(http.StatusOK, retObj("100", "success", chat))

	//这里需要启动多线程去交互LLM
	go makeReqPlatforms(userId, chat.Id, prompt, suppLLMs)
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

/* private methods area for multi thread */
func makeReqPlatforms(userId, chatId uint64, prompt string, suppLLMs []model.SysConf) {
	url := "http://192.168.3.17:7050/chat"
	aiChannels := make([]chan string, len(suppLLMs))
	for i := range suppLLMs {
		aiChannels[i] = make(chan string)
		go func(ch chan string, mod *model.SysConf) {
			params := map[string]interface{}{
				"model_id": mod.ConfValue,
				"chat_id":  userId,
				"scene":    "",
				"chat":     prompt,
			}
			b, _ := json.Marshal(params)

			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/octet-stream")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			responseData := string(body)
			fmt.Println(responseData)
			ch <- responseData
		}(aiChannels[i], &suppLLMs[i])
	}

	resTime := time.Now()
	db := database.GetDb()

	var chatObj model.Chat
	db.Model(&model.Chat{}).Where("id = ?", chatId).Last(&chatObj)

	for i, ch := range aiChannels { //遍历切片，等待子协程结束
		retStr := <-ch

		retObj := map[string]interface{}{}
		err := json.Unmarshal([]byte(retStr), &retObj)

		// rpCode := retObj["code"].(int)
		rp := retObj["AI_response"].(string)
		rpStr, _ := zhToUnicode([]byte(rp))

		if err != nil {
			rpStr = err.Error()
		}

		chat := &model.Chat{
			BaseModel: model.BaseModel{
				AddTime:    resTime,
				UpdateTime: resTime,
			},
			UserId:   userId,
			ChatId:   chatId,
			Content:  rpStr,
			Status:   1,
			Expect:   1,
			TaskCode: chatObj.TaskCode,
			Source:   suppLLMs[i].ConfKey,
		}
		db.Model(&model.Chat{}).Save(chat)
	}
	db.Model(&chatObj).Update("status", len(aiChannels))
}

func zhToUnicode(raw []byte) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}
