package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"testing"
	"time"
)

func TestFeedback_Create(t *testing.T) {
	go TestIndoor(t)
	go TestOutDoor(t)
}

func TestIndoor(t *testing.T) {

	dsn := "root:chd@163.com_1213@tcp(47.108.133.226:3389)/fault_warning_platform?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 获取当前时间
	currentTime := time.Now()

	// 获取当前时间的 Unix 时间戳（以秒为单位）
	timestamp := currentTime.Unix()

	for i := 0; i < 1000; i++ {
		t := currentTime.Add(time.Duration(i) * time.Second)
		indoor := IndoorDevice{
			DeviceId:       rand.Intn(10) + 1,
			Tl16:           rand.Intn(100),
			Tg1:            rand.Intn(100),
			Ti:             rand.Intn(100),
			BlowingAirTemp: rand.Intn(100),
			SetTemperature: rand.Intn(100),
			Fd:             rand.Intn(100),
			IfRun:          rand.Intn(1),
			Time:           t.Format("2006-01-02 15:04:05"),
			TimeStamp:      timestamp,
		}
		timestamp += 1
		indoor.Create(db)
	}
}

func TestOutDoor(t *testing.T) {
	dsn := "root:chd@163.com_1213@tcp(47.108.133.226:3389)/fault_warning_platform?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 获取当前时间
	currentTime := time.Now()

	// 获取当前时间的 Unix 时间戳（以秒为单位）
	timestamp := currentTime.Unix()

	for i := 0; i < 1000; i++ {
		t := currentTime.Add(time.Duration(i) * time.Second)
		outdoor := OutdoorDevice{
			DeviceId:  rand.Intn(2) + 1,
			Pd:        rand.Float64() * 100,
			Ps:        rand.Float64() * 100,
			Td1:       rand.Float64() * 100,
			Te1:       rand.Float64() * 100,
			Ta:        rand.Float64() * 100,
			Tfin:      rand.Float64() * 100,
			A12:       rand.Float64() * 100,
			A1:        rand.Float64() * 100,
			OE:        rand.Float64() * 100,
			H1:        rand.Float64() * 100,
			Fo:        rand.Float64() * 100,
			TdSH:      rand.Float64() * 100,
			Info1:     rand.Float64() * 100,
			Status:    rand.Float64() * 100,
			Time:      t.Format("2006-01-02 15:04:05"),
			TimeStamp: timestamp,
		}
		timestamp += 1
		outdoor.Create(db)
	}
}

func Test(u *testing.T) {
	// Unix时间戳
	timestamp := int64(1694161940)

	// 创建一个时间对象
	t := time.Unix(timestamp, 0)

	// 设置时区为北京时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tInCST := t.In(loc)

	fmt.Printf("Unix时间戳 %d 转换为北京时间：%s\n", timestamp, tInCST)
}
