package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
)

// SendHttpGetRequest 发送GET请求
func SendHttpGetRequest(requestUrl string, requestHeaders map[string]string, isWriteLog bool) (string, error) {
	// 创建一个新的 GET 请求
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	if requestHeaders != nil {
		for key, value := range requestHeaders {
			req.Header.Set(key, value)
		}
	}

	// 使用默认的 HTTP 客户端发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logx.Errorf("SendHttpGetRequest发送请求失败: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBodyStr := string(respBody)

	if isWriteLog {
		logData := map[string]interface{}{
			"requestUrl":       requestUrl,
			"requestHeaders":   requestHeaders,
			"responseHttpCode": resp.StatusCode, //HTTP 状态码
			"responseData":     respBodyStr,
		}
		fmt.Println("SendHttpGetRequestLogData:", JsonMarshalIndent(logData))
	}
	return respBodyStr, nil
}

// SendHttpPostRequest 发送POST请求
func SendHttpPostRequest(requestUrl string, requestBody interface{}, requestHeaders map[string]string, isWriteLog bool) (string, error) {
	// 将请求体转换为 JSON
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	// 设置请求头，例如 Content-Type
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if requestHeaders != nil {
		for key, value := range requestHeaders {
			req.Header.Set(key, value)
		}
	}

	// 使用默认的 HTTP 客户端发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logx.Errorf("SendHttpPostRequest发送请求失败: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("SendHttpPostRequest读取响应内容失败: %v", err)
		return "", err
	}
	respBodyStr := string(respBody)

	if isWriteLog {
		logData := map[string]interface{}{
			"requestUrl":       requestUrl,
			"requestHeaders":   requestHeaders,
			"requestBody":      requestBody,
			"responseHttpCode": resp.StatusCode, //HTTP 状态码
			"responseData":     respBodyStr,
		}
		fmt.Println("SendHttpPostRequestLogData:", JsonMarshalIndent(logData))
	}
	return respBodyStr, nil
}
