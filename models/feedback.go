package models

import "gorm.io/gorm"

type Feedback struct {
	Pk           uint   `gorm:"primary_key;auto_increment"` // 自增主键
	SystemResult string `json:"systemResult" gorm:"column:systemResult"`
	Iu1Result    string `json:"iu1Result" gorm:"column:iu1Result"`
	Iu2Result    string `json:"iu2Result" gorm:"column:iu2Result"`
	Iu3Result    string `json:"iu3Result" gorm:"column:iu3Result"`
	Iu4Result    string `json:"iu4Result" gorm:"column:iu4Result"`
	Iu5Result    string `json:"iu5Result" gorm:"column:iu5Result"`
	Iu6Result    string `json:"iu6Result" gorm:"column:iu6Result"`
	Iu7Result    string `json:"iu7Result" gorm:"column:iu7Result"`
	Iu8Result    string `json:"iu8Result" gorm:"column:iu8Result"`
	Iu9Result    string `json:"iu9Result" gorm:"column:iu9Result"`
	Iu10Result   string `json:"iu10Result" gorm:"column:iu10Result"`
	StartTime    int64  `json:"startTime" gorm:"column:startTime"`       // 开始时间
	StartTimeStr string `json:"startTimeStr" gorm:"column:startTimeStr"` // 开始时间
	EndTime      int64  `json:"endTime" gorm:"column:endTime"`           // 结束时间
	EndTimeStr   string `json:"endTimeStr" gorm:"column:endTimeStr"`     // 结束时间
}

func (i *Feedback) Create(db *gorm.DB) error {
	return db.Create(i).Error
}

func (i *Feedback) TableName() string {
	return "feedback"
}

func (i *Feedback) GetFeedbackOfTime(db *gorm.DB, startTime, endTime int64) (res []Feedback, err error) {
	err = db.Where("startTime >= ? and endTime <= ?", startTime, endTime).
		Find(&res).Error
	return
}

func (i *Feedback) GetLatestFeedback(db *gorm.DB, resultType string, result string, page int, limit int) (res []Feedback, err error) {
	err = db.Where(resultType+" <> ?", result).Order("pk DESC").Offset(page * limit).Limit(limit).Find(&res).Error
	return
}
