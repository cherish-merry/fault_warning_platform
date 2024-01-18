package config

import (
	"github.com/RaymondCode/simple-demo/amqp"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/conf"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type updateConfigParam struct {
	SkipPoint          int     `form:"skipPoint"`
	DeviationThreshold float64 `form:"deviationThreshold"`
	SlideWindow        int     `form:"slideWindow"`
}

type updateWarnParams struct {
	WarnParams []string `json:"warnParams"`
}

func UpdatePreprocessing(ctx *gin.Context) {
	g := common.GetGin(ctx)
	param := updateConfigParam{}
	err := ctx.ShouldBind(&param)
	if err != nil {
		log.Errorf("bind param fail, err:%v", err)
		g.ResponseFail()
		return
	}
	config := conf.OthersConfig
	if param.SlideWindow != 0 {
		config.SlideWindow = param.SlideWindow
	}
	if param.SkipPoint != 0 {
		config.SkipPoint = param.SkipPoint
	}
	if param.DeviationThreshold != 0 {
		config.DeviationThreshold = param.DeviationThreshold
	}
	g.ResponseNormal("Success")
}

func UpdateWarnParams(ctx *gin.Context) {
	g := common.GetGin(ctx)
	uWP := updateWarnParams{}
	if err := ctx.ShouldBindJSON(&uWP); err != nil {
		log.Errorf("bind param fail, err:%v", err)
		g.ResponseFail()
		return
	}

	amqp.UpdateConfig(uWP.WarnParams)

	g.ResponseNormal("Success")
}
