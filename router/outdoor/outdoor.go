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
	g.ResponseNormal(convert(devices))
}

func convert(devices []models.OutdoorDevice) map[string]interface{} {
	// 创建 sensor_data 切片
	sensorData := make([]map[string]interface{}, 0)
	// 创建 time 切片
	timeData := make([]string, 0)

	// 获取传感器类型列表
	sensorTypes := []string{"Pd", "Ps", "Td1", "Te1", "Ta", "Tfin", "A12", "A1", "OE", "H1", "Fo", "Tdsh", "Info1", "Status"}

	// 遍历传感器类型
	for _, sensorType := range sensorTypes {
		// 创建值切片
		values := make([]float64, 0)

		// 遍历设备数据，收集对应传感器类型的值
		for _, device := range devices {
			switch sensorType {
			case "Pd":
				values = append(values, device.Pd)
			case "Ps":
				values = append(values, device.Ps)
			case "Td1":
				values = append(values, device.Td1)
			case "Te1":
				values = append(values, device.Te1)
			case "Ta":
				values = append(values, device.Ta)
			case "Tfin":
				values = append(values, device.Tfin)
			case "A12":
				values = append(values, device.A12)
			case "A1":
				values = append(values, device.A1)
			case "OE":
				values = append(values, device.OE)
			case "H1":
				values = append(values, device.H1)
			case "Fo":
				values = append(values, device.Fo)
			case "Tdsh":
				values = append(values, device.TdSH)
			case "Info1":
				values = append(values, device.Info1)
			case "Status":
				values = append(values, device.Status)
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
