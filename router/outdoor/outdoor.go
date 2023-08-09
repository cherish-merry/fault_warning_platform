package outdoor

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type getOutdoorDeviceParam struct {
	DeviceId  int64 `form:"deviceId" binding:"required"`
	StartTime int64 `form:"startTime" binding:"required"`
	EndTime   int64 `form:"endTime" binding:"required"`
}

func GetOutdoorDevice(ctx *gin.Context) {
	g := common.GetGin(ctx)
	param := getOutdoorDeviceParam{}
	err := ctx.ShouldBindQuery(&param)
	if err != nil {
		log.Errorf("get param fail, err:%v", err)
		g.ResponseFail()
		return
	}
	outdoorDevice := models.OutdoorDevice{}
	db := database.GetInstanceConnection().GetPrimaryDB()
	devices, err := outdoorDevice.GetDeviceInfoOfTime(db, param.DeviceId, param.StartTime, param.EndTime)
	if err != nil {
		log.Errorf("get indoor device info fail, err:%v", err)
		g.ResponseFail()
		return
	}
	g.ResponseSuccess(devices)
}
