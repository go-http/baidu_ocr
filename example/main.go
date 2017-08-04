package main

import (
	"fmt"
	"time"

	"flag"

	".."
)

func main() {
	var key, secret, file string

	flag.StringVar(&key, "key", "", "应用的API Key")
	flag.StringVar(&secret, "secret", "", "应用的Secret Key")
	flag.StringVar(&file, "file", "", "要识别的文件路径")

	flag.Parse()

	client := ocr.New(key, secret)

	//请求识别表单
	requestResult, err := client.FormOcrRequest(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("表单识别请求log_id: %d, request_id: %s\n", requestResult.LogId, requestResult.Result[0].RequestId)

	for i := 8; i > 0; i-- {
		fmt.Printf("\r等待%d秒", i)
		time.Sleep(time.Second)
	}

	maxRetry := 10
	var resultDataString string
	for i := 0; i < maxRetry; i++ {
		//查询识别结果
		resultResp, err := client.FormOcrResult(requestResult.Result[0].RequestId)
		if err != nil {
			fmt.Println(err)
			return
		}

		if resultResp.Result.RetCode == 3 {
			resultDataString = resultResp.Result.ResultData
			break
		}

		fmt.Printf("\r[%2d/%2d]当前状态: %s", i, maxRetry, resultResp.Result.RetMsg)
		time.Sleep(time.Second)
	}

	fmt.Println("")

	if resultDataString == "" {
		fmt.Println("超时")
		return
	}
	result, err := ocr.ParseFormResult(resultDataString)
	if err != nil {
		fmt.Println("表单识别结果解析失败:%s", err)
		return
	}

	result.FormatDump(12)
}
