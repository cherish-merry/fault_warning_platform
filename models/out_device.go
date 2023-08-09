package models

import "gorm.io/gorm"

type OutdoorDevice struct {
	Pk        uint    `json:"-" gorm:"primary_key;auto_increment"` // 自增主键
	DeviceId  int     `json:"device_id" gorm:"column:device_id"`   // 设备id
	Pd        float64 `json:"pd" gorm:"column:pd"`                 // 高压压力测量值
	Ps        float64 `json:"ps" gorm:"column:ps"`                 // 低压压力计算值
	Td1       float64 `json:"td1" gorm:"column:td1"`               // 压缩机顶部温度
	Te1       float64 `json:"te1" gorm:"column:te1"`               // 室外换热器液侧温度
	Ta        float64 `json:"ta" gorm:"column:ta"`                 // 环境温度
	Tfin      float64 `json:"tfin1" gorm:"column:tfin1"`           // 变频散热片温度
	A12       float64 `json:"inv1a2" gorm:"column:inv1a2"`         // 压缩机二次侧电流
	A1        float64 `json:"inv1a1" gorm:"column:inv1a1"`         // 压缩机一次侧电流
	OE        float64 `json:"evo1" gorm:"column:evo1"`             // 室外电子膨胀阀开度比例
	H1        float64 `json:"h1" gorm:"column:h1"`                 // 压缩机运转频率
	Fo        float64 `json:"fo" gorm:"column:fo"`                 // 室外风机运转风速等级
	TdSH      float64 `json:"tdsh" gorm:"column:tdsh"`             // 排气温度与饱和冷凝温度差值。TdSH = Td1-Tcond
	Info1     float64 `json:"tsc" gorm:"column:tsc"`               // Te温度与饱和冷凝温度差值。TeSC =Tcond -Te
	Status    float64 `json:"ou_off" gorm:"column:ou_off"`         // 运行状态
	Time      string  `json:"up_unix_time,omitempty" gorm:"-"`     // 时间字符串
	TimeStamp int64   `json:"time" gorm:"column:up_unix_time"`     // 时间戳
}

func (o *OutdoorDevice) Create(db *gorm.DB) error {
	return db.Create(o).Error
}

func (o *OutdoorDevice) TableName() string {
	return "outdoor_device"
}

func (o *OutdoorDevice) GetDeviceInfoOfTime(db *gorm.DB, deviceId, startTime, endTime int64) (res []OutdoorDevice, err error) {
	err = db.Where("device_id = ? and up_unix_time >= ? and up_unix_time <= ?", deviceId, startTime, endTime).
		Find(&res).Error
	return
}
