package server

import (
	"HanHua/getNet"
	"HanHua/youdaoyunAPI"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Result 是对需要翻译的对象的抽象结果集
type NeedTrans struct {
	Key string
	Val string
}

func HanHuaServer(fileName string) {
	file, err := os.Open("./File/" + fileName)
	if err != nil {
		fmt.Printf("打开文件异常: < %v > \n", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("关闭文件 < %v > 异常 \n", file.Name())
		}
	}(file)

	// 存放待翻译字符
	var res []NeedTrans
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//分开后的值的对象 key为 "=" 前的值 val为等于号后的值
		entry := ParseLine(line)
		if entry != nil {
			res = append(res, *entry)
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Printf("读取文件异常: %v", err)
		return
	}

	DeepLTranslate(res)
}

// ParseLine 转换并拼接字符
func ParseLine(line string) *NeedTrans {
	split := strings.SplitN(line, "=", 2)
	if len(split) < 2 {
		// 处理没有 "=" 的情况，这里选择跳过
		return nil
	}
	return &NeedTrans{Key: split[0] + "=", Val: split[1]}
}

func DeepLTranslate(res []NeedTrans) {
	urlPre := "https://www.deepl.com/zh/translator#en/zh-hans/"
	var urls []string
	// 存储未翻译的索引以及字符
	noTrans := make(map[int]string)

	for i, re := range res {
		var url string
		// 转换标点 因为存在标点符号会让标签变更导致拿不到翻译
		if strings.Contains(re.Val, ". ") {
			re.Val = strings.ReplaceAll(re.Val, ". ", ".")
		}
		if strings.Contains(re.Val, ", ") {
			re.Val = strings.ReplaceAll(re.Val, ", ", ",")
		}
		if strings.Contains(re.Val, ": ") {
			re.Val = strings.ReplaceAll(re.Val, ": ", ":")
		}
		if strings.Contains(re.Val, "%") {
			re.Val = strings.ReplaceAll(re.Val, "%", "%25")
		}
		//TODO 无法直接翻译带 <br> 的字符
		if strings.Contains(re.Val, "<br>") {
			noTrans[i] = re.Val
			urls = append(urls, "")
			continue
		}

		url = urlPre + re.Val
		urls = append(urls, url)
	}
	// 获取翻译结果集
	results := getNet.TranslateDeepL(urls)

	/*for _, url := range urls {
		fmt.Println(url)
		time.Sleep(1 * time.Second)
	}*/
	for i, result := range results {
		if result.Err != nil {
			WriteTranslate(res[i].Key + "此行报错！！！！！！！！！！！！！！！！！！！！！: " + result.Err.Error())
			continue
		}
		if result.Translate == "" {
			if ele, exist := noTrans[result.Index]; exist {
				result.Translate = ele
			}
		}
		WriteTranslate(res[i].Key + result.Translate)
		//fmt.Println(res[i].Key + result.Translate)
		for index, ele := range noTrans {
			fmt.Printf("第 %v 行出现存在 <br> 需要手动翻译 原文为： %v \n", index, ele)
			WriteTranslate(res[i].Key + result.Translate)
		}
	}
}

// YoudaoTranslate 使用有道云API翻译 需要对应的 appKey 和 appSecret
func YoudaoTranslate(res []NeedTrans) {
	//遍历并且输出参数
	for _, entry := range res {
		if entry.Val == "" {
			WriteTranslate(entry.Key)
			fmt.Println(entry.Key)
			continue
		}
		result := youdaoyunAPI.TransYouDaoYun(entry.Val)
		WriteTranslate(entry.Key + result)
		fmt.Println(entry.Key + result)

		// 根据API限制调整睡眠时间 这里用 2 秒翻译一次
		time.Sleep(2000 * time.Millisecond)
	}
}

// WriteTranslate 输出到 zh_CN.lang 文件中
func WriteTranslate(text string) {
	//out, err := os.Create("./File/zh_CN.lang")
	out, err := os.OpenFile("./File/zh_CN.lang", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("创建输出文件异常: %v\n", err)
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			fmt.Printf("关闭文件 < %v > 异常 \n", out.Name())
		}
	}(out)

	writer := bufio.NewWriter(out)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			fmt.Printf("Flush 文件 < %v > 异常 \n", out.Name())
		}
	}(writer)

	// 输出到文件中
	_, err = writer.WriteString(text + "\n")
	if err != nil {
		fmt.Printf("写入文件异常: %v\n", err)
	}
}
