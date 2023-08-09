package indoor

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type getIndoorDeviceParam struct {
	DeviceId  int64 `form:"deviceId" binding:"required"`
	StartTime int64 `form:"startTime" binding:"required"`
	EndTime   int64 `form:"endTime" binding:"required"`
}

func GetIndoorDevice(ctx *gin.Context) {
	g := common.GetGin(ctx)
	param := getIndoorDeviceParam{}
	err := ctx.ShouldBindQuery(&param)
	if err != nil {
		log.Errorf("get param fail, err:%v", err)
		g.ResponseFail()
		return
	}
	indoorDevice := models.IndoorDevice{}
	db := database.GetInstanceConnection().GetPrimaryDB()
	devices, err := indoorDevice.GetDeviceInfoOfTime(db, param.DeviceId, param.StartTime, param.EndTime)
	if err != nil {
		log.Errorf("get indoor device info fail, err:%v", err)
		g.ResponseFail()
		return
	}
	g.ResponseSuccess(devices)
}
