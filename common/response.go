package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GeneralResponse 通用响应
type GeneralResponse struct {
	StatusCode int         `json:"status_code"` // 状态码
	StatusMsg  string      `json:"status_msg"`  // 消息
	Result     interface{} `json:"result"`      // 数据
}

type Gin struct {
	C *gin.Context
}

const (
	SUCCESS            = 200
	ERROR              = 500
	ErrorInvalidParams = 400
)

var MsgFlags = map[int]string{
	SUCCESS:            "success",
	ERROR:              "error",
	ErrorInvalidParams: "param error",
}

func GetMsg(code int) string {
	return MsgFlags[code]
}

func GetGin(c *gin.Context) Gin {
	return Gin{c}
}

func (g *Gin) ResponseFail() {
	g.C.JSON(http.StatusOK, gin.H{
		"status_code": ERROR,
		"status_msg":  GetMsg(ERROR),
	})
	return
}

func (g *Gin) ResponseNormal(data interface{}) {
	g.C.JSON(http.StatusOK, data)
	return
}

func (g *Gin) ResponseSuccess(data interface{}) {
	if data != nil {
		g.C.JSON(http.StatusOK, gin.H{
			"status_code": SUCCESS,
			"status_msg":  GetMsg(SUCCESS),
			"result":      data,
		})
		return
	} else {
		g.C.JSON(http.StatusOK, gin.H{
			"status_code": SUCCESS,
			"status_msg":  GetMsg(SUCCESS),
		})
		return
	}
}
