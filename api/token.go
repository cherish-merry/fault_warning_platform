package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var GlobalToken *AccessToken
var tokenLock sync.RWMutex

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func GetGlobalToken() string {
	tokenLock.RLock()
	defer tokenLock.RUnlock()
	return GlobalToken.AccessToken
}

func InitToken() error {
	err := RefreshToken()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Minute * 60)
	go func() {
		for range ticker.C {
			err = RefreshToken()
			if err != nil {
				log.Error(err)
			}
		}
	}()
	return nil
}

func RefreshToken() error {
	token, err := getAccessToken()
	if err != nil {
		return fmt.Errorf("get access token err: %v", err)

	}
	tokenLock.Lock()
	GlobalToken = token
	tokenLock.Unlock()
	return nil
}

func getAccessToken() (*AccessToken, error) {
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
