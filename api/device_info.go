package api

import (
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func getData(url string, token string) ([]byte, error) {
	method := "GET"

	// 创建HTTP请求
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return nil, err
	}
	// 设置请求头
	req.Header.Add("Authorization", "bearer "+token)

	// 发送HTTP请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading response: %v", err)
		return nil, err
	}

	return body, nil
}

func mergeData(firstPartData, secondPartData string) (*models.IndoorDevice, error) {
	var myData models.IndoorDevice

	// Parse first part data
	var firstPartMap map[string]interface{}
	err := json.Unmarshal([]byte(firstPartData), &firstPartMap)
	if err != nil {
		return &myData, fmt.Errorf("error parsing first part data: %v", err)
	}

	// Extract desired fields from first part data
	myData.Ti = int(firstPartMap["data"].(map[string]interface{})["ti"].(float64))
	myData.BlowingAirTemp = int(firstPartMap["data"].(map[string]interface{})["blowing_air_temp"].(float64))
	myData.SetTemperature = int(firstPartMap["data"].(map[string]interface{})["set_temperature"].(float64))
	myData.IfRun = int(firstPartMap["data"].(map[string]interface{})["if_run"].(float64))
	myData.Time = firstPartMap["data"].(map[string]interface{})["up_unix_time"].(string)

	// Parse second part data
	var secondPartMap map[string]interface{}
	err = json.Unmarshal([]byte(secondPartData), &secondPartMap)
	if err != nil {
		return &myData, fmt.Errorf("error parsing second part data: %v", err)
	}

	// Extract desired fields from second part data
	myData.Tl16 = int(secondPartMap["data"].(map[string]interface{})["tl16"].(float64))
	myData.Tg1 = int(secondPartMap["data"].(map[string]interface{})["tg1"].(float64))
	myData.Fd = int(secondPartMap["data"].(map[string]interface{})["fd"].(float64))
	return &myData, nil
}

func GetOutdoorDeviceInfo(url1, url2 string, token string) (*models.OutdoorDevice, error) {
	// 第一个url获取除status外的其他字段
	body, err := getData(url1, token)
	if err != nil {
		log.Errorf("Get Data err: %v:", err)
		return nil, err
	}

	// 定义一个匿名结构体，仅包含data字段
	var response struct {
		Code int                  `json:"code"`
		Msg  string               `json:"msg"`
		Data models.OutdoorDevice `json:"data"`
	}

	// 解析JSON响应到结构体
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Errorf("Error decoding JSON: %v", err)
		return nil, err
	}

	// 第二个url，获取status
	body, err = getData(url2, token)
	if err != nil {
		log.Errorf("Get Data err: %v:", err)
		return nil, err
	}
	var dataMap map[string]interface{}
	err = json.Unmarshal(body, &dataMap)
	if err != nil {
		return nil, fmt.Errorf("error parsing status info: %v", err)
	}
	response.Data.Status = dataMap["data"].(map[string]interface{})["ou_off"].(float64)
	return &response.Data, nil
}

func GetIndoorDeviceInfo(url1 string, url2 string, token string) (*models.IndoorDevice, error) {
	part1, err := getData(url1, token)
	if err != nil {
		log.Errorf("Get Data err: %v:", err)
		return nil, err
	}
	part2, err := getData(url2, token)
	if err != nil {
		log.Errorf("Get Data err: %v:", err)
		return nil, err
	}
	data, err := mergeData(string(part1), string(part2))
	if err != nil {
		log.Errorf("MergeData err: %v:", err)
		return nil, err
	}
	return data, nil
}
