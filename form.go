package ocr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

//请求识别表单
type FormOcrRequestResponse struct {
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`

	LogId  int64 `json:"log_id"`
	Result []struct {
		RequestId string `json:"request_id"`
	}
}

//请求识别表单
func (cli *Client) FormOcrRequest(filename string) (FormOcrRequestResponse, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return FormOcrRequestResponse{}, fmt.Errorf("读取图片错误: %s", err)
	}

	postData := url.Values{}
	postData.Set("image", base64.StdEncoding.EncodeToString(content))

	resp, err := cli.PostForm("/form_ocr/request", nil, postData)
	if err != nil {
		return FormOcrRequestResponse{}, fmt.Errorf("请求错误: %s", err)
	}
	defer resp.Body.Close()

	var respInfo FormOcrRequestResponse
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return FormOcrRequestResponse{}, fmt.Errorf("读取错误: %s", err)
	}

	if respInfo.ErrorCode != 0 {
		return FormOcrRequestResponse{}, fmt.Errorf("读取错误: [%d] %s", respInfo.ErrorCode, respInfo.ErrorMsg)
	}

	return respInfo, nil
}

//表单识别结果查询响应
type FormOcrResultResponse struct {
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`

	LogId  int64 `json:"log_id"`
	Result struct {
		ResultData string `json:"result_data"`
		Percent    int
		RequestId  string `json:"request_id"`
		RetCode    int    `json:"ret_code"`
		RetMsg     string `json:"ret_msg"`
	}
}

//表单识别结果查询
func (cli *Client) FormOcrResult(requestId string) (FormOcrResultResponse, error) {
	postData := url.Values{}
	postData.Set("request_id", requestId)
	postData.Set("result_type", "json")

	resp, err := cli.PostForm("/form_ocr/get_request_result", nil, postData)
	if err != nil {
		return FormOcrResultResponse{}, fmt.Errorf("请求错误: %s", err)
	}
	defer resp.Body.Close()

	var respInfo FormOcrResultResponse
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return FormOcrResultResponse{}, fmt.Errorf("读取错误: %s", err)
	}

	if respInfo.ErrorCode != "" {
		return FormOcrResultResponse{}, fmt.Errorf("读取错误: [%d] %s", respInfo.ErrorCode, respInfo.ErrorMsg)
	}

	return respInfo, nil
}
