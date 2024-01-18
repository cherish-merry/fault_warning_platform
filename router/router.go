package router

import (
	"github.com/RaymondCode/simple-demo/router/config"
	"github.com/RaymondCode/simple-demo/router/feedback"
	"github.com/RaymondCode/simple-demo/router/file"
	"github.com/RaymondCode/simple-demo/router/indoor"
	"github.com/RaymondCode/simple-demo/router/outdoor"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/api")
	apiRouter.Use(cors.Default())
	apiRouter.GET("/indoor/get", indoor.GetIndoorDevice)
	apiRouter.GET("/outdoor/get", outdoor.GetOutdoorDevice)
	apiRouter.GET("/feedback/get", feedback.GetFeedback)
	apiRouter.GET("/feedback/latest", feedback.GetLatestFeedBack)
	apiRouter.GET("/feedback/day", feedback.GetFeedBackByDay)
	apiRouter.POST("/upload", file.Upload)
	apiRouter.GET("/download", file.Download)
	apiRouter.POST("/config/preprocessing", config.UpdatePreprocessing)
	apiRouter.POST("/config/warnParams", config.UpdateWarnParams)
}
