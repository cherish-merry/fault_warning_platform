package user

import (
	"time"

	"github.com/RaymondCode/simple-demo/common"
	userModel "github.com/RaymondCode/simple-demo/models/user"
	userService "github.com/RaymondCode/simple-demo/service/user"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Register 注册用户
func Register(c *gin.Context) {
	var (
		g        = common.GetGin(c)
		userId   = util.NewShortIDString("user")
		username = c.Query("username")
		password = c.Query("password")
	)

	err := userService.Register(&userModel.User{
		Id:        userId,
		Username:  username,
		Password:  password,
		CreatedAt: time.Now().UTC(),
	})

	if err != nil {
		log.Errorf("register user fail err:%v", err)
		g.ResponseFail()
		return
	}
	g.ResponseSuccess(userId)
}
