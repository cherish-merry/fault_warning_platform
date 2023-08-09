package router

import (
	"github.com/RaymondCode/simple-demo/router/indoor"
	"github.com/RaymondCode/simple-demo/router/outdoor"
	"github.com/RaymondCode/simple-demo/router/user"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/api")
	apiRouter.POST("/user/register/", user.Register)

	apiRouter.GET("/indoor/get", indoor.GetIndoorDevice)
	apiRouter.GET("/outdoor/get", outdoor.GetOutdoorDevice)
}
