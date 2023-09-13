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
	g.ResponseNormal(convert(devices))
}

func convert(devices []models.IndoorDevice) map[string]interface{} {
	// 创建 sensor_data 切片
	sensorData := make([]map[string]interface{}, 0)
	// 创建 time 切片
	timeData := make([]string, 0)

	// 获取传感器类型列表
	sensorTypes := []string{"Tl16", "Tg1", "Ti", "BlowingAirTemp", "SetTemperature", "Fd", "IfRun", "Dt"}

	// 遍历传感器类型
	for _, sensorType := range sensorTypes {
		// 创建值切片
		values := make([]int, 0)

		// 遍历设备数据，收集对应传感器类型的值
		for _, device := range devices {
			switch sensorType {
			case "Tl16":
				values = append(values, device.Tl16)
			case "Tg1":
				values = append(values, device.Tg1)
			case "Ti":
				values = append(values, device.Ti)
			case "BlowingAirTemp":
				values = append(values, device.BlowingAirTemp)
			case "SetTemperature":
				values = append(values, device.SetTemperature)
			case "Fd":
				values = append(values, device.Fd)
			case "IfRun":
				values = append(values, device.IfRun)
			case "Dt":
				values = append(values, device.Dt)
			}
		}

		// 创建传感器数据 map
		sensor := map[string]interface{}{
			"sensor_type": sensorType,
			"value":       values,
		}

		// 添加到 sensor_data 切片
		sensorData = append(sensorData, sensor)
	}

	// 遍历设备数据，收集时间数据
	for _, device := range devices {
		timeData = append(timeData, device.Time)
	}

	// 创建最终结果 map
	result := map[string]interface{}{
		"sensor_data": sensorData,
		"time":        timeData,
	}
	return result
}
