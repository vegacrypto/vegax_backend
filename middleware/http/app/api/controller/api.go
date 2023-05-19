package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vegacrypto/vegax_backend/middleware/http/common"
)

func GetEventsByContract(c *gin.Context) {
	res := common.Response{}

	res.Code = 0
	res.Msg = "success"
	// res.Data = data
	c.JSON(http.StatusOK, res)
}
