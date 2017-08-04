package ocr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

//百度AI开放平台基础地址
const BaseUrl = "https://aip.baidubce.com"

//百度AI开放平台Restful接口基础地址
const BaseRestUrl = BaseUrl + "/rest/2.0/solution/v1"

//OCR识别请求发起用的客户端
type Client struct {
	AppID     string
	AppSecret string
}

func New(appId, appSecret string) *Client {
	return &Client{AppID: appId, AppSecret: appSecret}
}

type AccessTokenResponse struct {
	Error            string
	ErrorDescription string `json:"error_description"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string
	SessionKey       string `json:"session_key"`
	AccessToken      string `json:"access_token"`
	SessionSecret    string `json:"session_secret"`
}

//获取AccessToken，TODO:增加缓存机制
func (cli *Client) getAccessToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", cli.AppID)
	data.Set("client_secret", cli.AppSecret)

	resp, err := http.PostForm(BaseUrl+"/oauth/2.0/token", data)
	if err != nil {
		return "", fmt.Errorf("请求错误: %s", err)
	}
	defer resp.Body.Close()

	var respInfo AccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return "", fmt.Errorf("读取错误: %s", err)
	}

	if respInfo.Error != "" {
		return "", fmt.Errorf("失败: (%s) %s", respInfo.Error, respInfo.ErrorDescription)
	}

	return respInfo.AccessToken, nil
}

//执行PostForm请求(自动添加AccessToken)
func (cli *Client) PostForm(path string, urlParams, postData url.Values) (*http.Response, error) {
	token, err := cli.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取Token失败:%s", err)
	}

	if urlParams == nil {
		urlParams = url.Values{}
	}
	urlParams.Set("access_token", token)

	return http.PostForm(BaseRestUrl+path+"?"+urlParams.Encode(), postData)
}
