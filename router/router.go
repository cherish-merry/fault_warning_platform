package router

import (
	"github.com/RaymondCode/simple-demo/router/user"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/api")
	apiRouter.POST("/user/register/", user.Register)
}
