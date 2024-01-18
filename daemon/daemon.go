package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/amqp"
	"github.com/RaymondCode/simple-demo/api"
	"github.com/RaymondCode/simple-demo/conf"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service/outdoor"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"math"
	"sync"
	"time"
)

var outDeviceUrlMap = map[int][2]string{
	1: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did106h_7_0100452c2f257624_0000",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did104h_7_0100452c2f257624_0000"},
	2: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did106h_7_0100452c2f257624_0001",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did104h_7_0100452c2f257624_0001"},
}

var innerDeviceUrlMap = map[int][2]string{
	1: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0001",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0001"},
	2: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0002",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0002"},
	3: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_5_0100452c2f257624_0003",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_5_0100452c2f257624_0003"},
	4: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0004",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0004"},
	5: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_69_0100452c2f257624_0005",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_69_0100452c2f257624_0005"},
	6: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0006",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0006"},
	7: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0007",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0007"},
	8: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0008",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0008"},
	9: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_5_0100452c2f257624_0009",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_5_0100452c2f257624_0009"},
	10: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0010",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0010"},
}

var outDoorDeviceList []*models.OutdoorDevice
var inDoorDeviceList []*models.IndoorDevice
var outdoorLock sync.RWMutex
var indoorLock sync.RWMutex

var point int
var standardDeviation []float64

func InitDaemon() {
	// Create a wait group to ensure all goroutines finish before exiting the program
	var wg sync.WaitGroup

	messageChan := make(chan string)

	// Start the scheduler to produce data
	wg.Add(1)
	go scheduler(&wg, messageChan)

	// Start the consumer to handle data in batches
	wg.Add(1)
	go consumer(&wg, messageChan)

	go amqp.HandlerMessage()

	// Wait for all goroutines to finish their jobs
	wg.Wait()
}

/*
1、获取企业数据 api.GetOutdoorDeviceInfo(urlArr[0], urlArr[1], token) api.GetIndoorDeviceInfo(urlArr[0], urlArr[1], token)
2、构建数据 func buildMessage(messageChan chan string)
2.1、过滤数据 func filter(pd float64) bool {}
3、三十条数据发送到模型端 func consumer(wg *sync.WaitGroup, messageChan chan string) {}
*/
func scheduler(wg *sync.WaitGroup, messageChan chan string) {
	// Stop the scheduler when the function finishes
	defer wg.Done()

	config := conf.OthersConfig
	// Create a ticker with 1-second interval
	log.Infof("collect interval: %vs", config.CollectInterval)
	ticker := time.NewTicker(time.Duration(config.CollectInterval) * time.Second)
	db := database.GetInstanceConnection().GetPrimaryDB()
	db.Logger = logger.Default.LogMode(logger.Silent)
	// Run the scheduler until it's stopped
	for {
		select {
		case <-ticker.C:
			//log.Info("start............")
			var requestWg sync.WaitGroup
			token := api.GetGlobalToken()
			requestWg.Add(len(outDeviceUrlMap) + len(innerDeviceUrlMap))

			for deviceId, urlArr := range outDeviceUrlMap {
				urlArr := urlArr
				deviceId := deviceId
				go func() {
					defer requestWg.Done()
					device, err := api.GetOutdoorDeviceInfo(urlArr[0], urlArr[1], token)
					if err != nil {
						log.Errorf("GetOutdoorDeviceInfo Error: %v", err)
					} else {
						device.DeviceId = deviceId
						unixTime, _ := time.Parse("2006-01-02 15:04:05", device.Time)
						device.TimeStamp = unixTime.Unix()
						err = device.Create(db)
						if err != nil {
							log.Errorf("create outdoor device Error: %v", err)
						}
					}
					outdoorLock.Lock()
					if device != nil {
						outDoorDeviceList = append(outDoorDeviceList, device)
					}
					outdoorLock.Unlock()
				}()
			}
			for deviceId, urlArr := range innerDeviceUrlMap {
				deviceId := deviceId
				urlArr := urlArr
				go func() {
					defer requestWg.Done()
					device, err := api.GetIndoorDeviceInfo(urlArr[0], urlArr[1], token)
					if err != nil {
						log.Errorf("GetIndoorDeviceInfo Error: %v", err)
					} else {
						device.DeviceId = deviceId
						unixTime, _ := time.Parse("2006-01-02 15:04:05", device.Time)
						device.TimeStamp = unixTime.Unix()
						device.Dt = device.Tg1 - device.Tl16
						err = device.Create(db)
						if err != nil {
							log.Errorf("create indoor device Error: %v", err)
						}
					}
					indoorLock.Lock()
					if device != nil {
						inDoorDeviceList = append(inDoorDeviceList, device)
					}
					indoorLock.Unlock()
				}()
			}

			requestWg.Wait()
			buildMessage(messageChan)
		}
	}
}

func reset() {
	outDoorDeviceList = []*models.OutdoorDevice{}
	inDoorDeviceList = []*models.IndoorDevice{}
}

func buildMessage(messageChan chan string) {
	// Create a map to hold the final JSON data
	//log.Info("build message.......")

	defer reset()

	if len(inDoorDeviceList) != len(innerDeviceUrlMap) || len(outDoorDeviceList) != len(outDeviceUrlMap) {
		//in, _ := json.Marshal(inDoorDeviceList)
		//out, _ := json.Marshal(outDoorDeviceList)
		log.Errorf("data loss: %v %v", len(inDoorDeviceList), len(outDoorDeviceList))
		return
	}

	jsonData2 := outdoor.MachineMask(outDoorDeviceList)
	config := conf.OthersConfig
	if config.Filter && filter(jsonData2["Pd"]) {
		return
	}

	standardDeviation = append(standardDeviation, jsonData2["Pd"])

	jsonData1 := make(map[string]interface{})
	jsonData1["startTime"] = outDoorDeviceList[0].TimeStamp
	jsonData1["startTimeStr"] = outDoorDeviceList[0].Time
	// Iterate through each IndoorDevice and add its fields to the map

	for i, device := range inDoorDeviceList {
		//log.Infof("indoor_%v time: %v", i, device.Time)
		jsonData1[fmt.Sprintf("iE_%d", i+1)] = device.DeviceId
		jsonData1[fmt.Sprintf("Tl_%d", i+1)] = device.Tl16
		jsonData1[fmt.Sprintf("Tg_%d", i+1)] = device.Tg1
		jsonData1[fmt.Sprintf("Ti_%d", i+1)] = device.Ti
		jsonData1[fmt.Sprintf("To_%d", i+1)] = device.BlowingAirTemp
		jsonData1[fmt.Sprintf("dT_%d", i+1)] = device.Tg1 - device.Tl16
		jsonData1[fmt.Sprintf("Ts_%d", i+1)] = device.SetTemperature
		jsonData1[fmt.Sprintf("fd_%d", i+1)] = device.Fd
		jsonData1[fmt.Sprintf("status_%d", i+1)] = device.IfRun
	}

	//log.Infof("outdoor_%v time: %v", 0, outDoorDeviceList[0].Time)
	//log.Infof("outdoor_%v time: %v", 1, outDoorDeviceList[1].Time)
	// Merge jsonData2 into jsonData1
	for key, value := range jsonData2 {
		jsonData1[key] = value
	}

	marshal, _ := json.Marshal(jsonData1)

	select {
	case messageChan <- string(marshal):
	default:
		log.Info("Queue is full. Skipping element.")
	}
}

// 数据过滤
func filter(pd float64) bool {
	if pd <= 0.0 {
		point = 0
		return true
	} else {
		point++
		if point > conf.OthersConfig.SkipPoint {
			return false
		}
		return true
	}
}

// 计算标准差的函数
func calculateStandardDeviation(data []float64) float64 {
	n := len(data)

	// 计算均值
	mean := 0.0
	for _, value := range data {
		mean += value
	}
	mean /= float64(n)

	// 计算方差
	variance := 0.0
	for _, value := range data {
		variance += (value - mean) * (value - mean)
	}
	variance /= float64(n)

	log.Infof("标准差：%v", variance)

	// 计算标准差
	return math.Sqrt(variance)
}

func consumer(wg *sync.WaitGroup, messageChan chan string) {
	config := conf.OthersConfig

	// Stop the consumer when the function finishes
	batchSize := config.SlideWindow

	defer wg.Done()

	// Create a list to store the elements
	var messageList []string
	for {
		select {
		case message := <-messageChan:
			messageList = append(messageList, message)
		}
		//log.Infof("message queue size: %d", len(messageList))
		// If the list reaches 30 elements, process it
		if len(messageList) == batchSize {
			if !config.Filter || calculateStandardDeviation(standardDeviation) <= config.DeviationThreshold {
				amqp.SendMessage(messageList)
			}
			messageList = nil
			standardDeviation = nil
		}
	}
}
