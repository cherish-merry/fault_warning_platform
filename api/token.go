package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func GetAssessToken() (*AccessToken, error) {
	url := "https://cloudmaster.hisensehitachi.com/auth/oauth/token?username=p_dcfyzd_shijingfeng&password=W%2Bx6Ljdj7ZlLz6wDkpju3w%3D%3D&grant_type=client_credentials&scope=server"
	method := "GET"

	// 创建HTTP请求
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Errorf("Error creating request:%v", err)
		return nil, err
	}

	// 设置请求头
	req.Header.Add("X-His-Brand", "hitachi")
	req.Header.Add("Authorization", "Basic aGhsaW5rOmhobGluaw==")

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
	// 打印响应内容
	log.Info(string(body))

	var response AccessToken
	err = json.Unmarshal([]byte(string(body)), &response)
	if err != nil {
		log.Errorf("Error decoding JSON: %v", err)
		return nil, err
	}
	return &response, nil
}
