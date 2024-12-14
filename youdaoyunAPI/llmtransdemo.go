package youdaoyunAPI

import (
	"fmt"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`"transFull":"([^"]*)"`)

// 您的应用ID
var appKey = ""

// 您的应用密钥
var appSecret = ""

func TransYouDaoYun(val string) string {
	// 添加请求参数
	paramsMap := createRequestParams(val)
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	// 添加鉴权相关参数
	AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	events := DoPostBySSE("https://openapi.youdao.com/llm_trans", header, paramsMap)
	for event := range events {
		// 处理接收到的事件
		//fmt.Println(event)
		if strings.Contains(event, "transFull") {
			//fmt.Println(event)
			matches := re.FindStringSubmatch(event)
			return matches[1]
		}

		if strings.Contains(event, "msg") {
			fmt.Printf("存在异常: %s \n", event)
		}
	}
	/*	// 请求api服务
		result := DoPost("https://openapi.youdao.com/api", header, paramsMap, "application/json")
		// 打印返回结果
		var resp RespTran
		err := json.Unmarshal(result, &resp)
		if err != nil {
			fmt.Println("json绑定异常")
		}
		if result != nil {
			//fmt.Print(resp.Translation[0])
			return resp.Translation[0]
		}*/
	return ""
}

func createRequestParams(val string) map[string][]string {
	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/trans/api/dmxfy/index.html
	*/
	i := val
	from := "en"
	to := "zh-CHS"
	polishOption := "8"
	expandOption := "0"
	streamType := "full"
	//domain := "game"

	return map[string][]string{
		"q":            {i},
		"i":            {i},
		"from":         {from},
		"to":           {to},
		"polishOption": {polishOption},
		"expandOption": {expandOption},
		"streamType":   {streamType},
		//"domain": {domain},
	}
}
