package models

import "gorm.io/gorm"

type IndoorDevice struct {
	Pk             uint   `json:"-" gorm:"primary_key;auto_increment"`             // 自增主键
	DeviceId       int    `json:"device_id" gorm:"column:device_id"`               // 设备id
	Tl16           int    `gorm:"column:tl16" json:"tl16"`                         // 室内机液管温度
	Tg1            int    `gorm:"column:tg1" json:"tg1"`                           // 室内机气管温度
	Ti             int    `gorm:"column:ti" json:"ti"`                             // 室内回风温度
	BlowingAirTemp int    `gorm:"column:blowing_air_temp" json:"blowing_air_temp"` // 室内出风温度
	SetTemperature int    `gorm:"column:set_temperature" json:"set_temperature"`   // 设定温度
	Fd             int    `gorm:"column:fd" json:"fd"`                             // 内机期望压机功率
	IfRun          int    `gorm:"column:if_run" json:"if_run"`                     // 开机状态
	Dt             int    `json:"dt" gorm:"column:dt"`
	Time           string `json:"up_unix_time,omitempty" gorm:"-"` // 时间戳
	TimeStamp      int64  `json:"time" gorm:"column:up_unix_time"` // 时间戳
}

func (i *IndoorDevice) Create(db *gorm.DB) error {
	return db.Create(i).Error
}

func (i *IndoorDevice) TableName() string {
	return "indoor_device"
}

func (i *IndoorDevice) GetDeviceInfoOfTime(db *gorm.DB, deviceId, startTime, endTime int64) (res []IndoorDevice, err error) {
	err = db.Where("device_id = ? and up_unix_time >= ? and up_unix_time <= ?", deviceId, startTime, endTime).
		Find(&res).Error
	return
}
