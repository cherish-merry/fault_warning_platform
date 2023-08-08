package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type OutdoorDevice struct {
	Pd     float64 `json:"pd"`           // 高压压力测量值
	Ps     float64 `json:"ps"`           // 低压压力计算值
	Td1    float64 `json:"td1"`          // 压缩机顶部温度
	Te1    float64 `json:"te1"`          // 室外换热器液侧温度
	Ta     float64 `json:"ta"`           // 环境温度
	Tfin   float64 `json:"tfin1"`        // 变频散热片温度
	A12    float64 `json:"inv1a2"`       // 压缩机二次侧电流
	A1     float64 `json:"inv1a1"`       // 压缩机一次侧电流
	OE     float64 `json:"evo1"`         // 室外电子膨胀阀开度比例
	H1     float64 `json:"h1"`           // 压缩机运转频率
	Fo     float64 `json:"fo"`           // 室外风机运转风速等级
	TdSH   float64 `json:"tdsh"`         // 排气温度与饱和冷凝温度差值。TdSH = Td1-Tcond
	Info1  float64 `json:"tsc"`          // Te温度与饱和冷凝温度差值。TeSC =Tcond -Te
	Status int     `json:"status"`       // 运行状态
	Time   string  `json:"up_unix_time"` // 时间戳
}

type IndoorDevice struct {
	Ti             int64 `json:"ti"`
	BlowingAirTemp int   `json:"blowing_air_temp"`
	SetTemperature int   `json:"set_temperature"`
	IfRun          int   `json:"if_run"`
	Tl16           int   `json:"tl16"`
	Tg1            int   `json:"tg1"`
	Fd             int   `json:"fd"`
}

func getData(url string, token string) ([]byte, error) {
	method := "GET"

	// 创建HTTP请求
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return nil, err
	}

	// 设置请求头
	req.Header.Add("Authorization", token)

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

func mergeData(firstPartData, secondPartData string) (*IndoorDevice, error) {
	var myData IndoorDevice

	// Parse first part data
	var firstPartMap map[string]interface{}
	err := json.Unmarshal([]byte(firstPartData), &firstPartMap)
	if err != nil {
		return &myData, fmt.Errorf("error parsing first part data: %v", err)
	}

	// Extract desired fields from first part data
	myData.Ti = int64(firstPartMap["data"].(map[string]interface{})["ti"].(float64))
	myData.BlowingAirTemp = int(firstPartMap["data"].(map[string]interface{})["blowing_air_temp"].(float64))
	myData.SetTemperature = int(firstPartMap["data"].(map[string]interface{})["set_temperature"].(float64))
	myData.IfRun = int(firstPartMap["data"].(map[string]interface{})["if_run"].(float64))

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

	log.Info(myData)

	return &myData, nil
}

func GetOutdoorDeviceInfo(url string, token string) (*OutdoorDevice, error) {
	body, err := getData(url, token)
	if err != nil {
		log.Errorf("Get Data err: %v:", err)
		return nil, err
	}

	// 定义一个匿名结构体，仅包含data字段
	var response struct {
		Code int           `json:"code"`
		Msg  string        `json:"msg"`
		Data OutdoorDevice `json:"data"`
	}

	// 解析JSON响应到结构体
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Errorf("Error decoding JSON: %v", err)
		return nil, err
	}
	return &response.Data, nil
}

func GetIndoorDeviceInfo(url1 string, url2 string, token string) (*IndoorDevice, error) {
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
