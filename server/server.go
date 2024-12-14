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

// NeedTrans 是对需要翻译的对象的抽象结果集
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
		if entry := ParseLine(line); entry != nil {
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
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		// 处理没有 "=" 的情况，这里选择跳过
		return nil
	}
	return &NeedTrans{Key: parts[0] + "=", Val: parts[1]}
}

func DeepLTranslate(res []NeedTrans) {
	// 需要转换的符号
	replacements := []string{
		". ", ",",
		",", " ",
		": ", ":",
		"! ", "!",
		"? ", "?",
		"%", "%25",
		"<br>", "*QJK*",
	}

	// 打开文件流
	fw, err := NewFileWriter("./File/zh_CN.lang")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = fw.Close(); err != nil {
			fmt.Printf("关闭文件失败: %v\n", err)
		}
	}()

	urlPre := "https://www.deepl.com/zh/translator#en/zh-hans/"
	var urls []string

	// 拼接所有待翻译的 url 存入数组中
	for _, re := range res {
		// 转换标点 因为存在标点符号会让标签变更导致拿不到翻译
		re.Val = strings.NewReplacer(replacements...).Replace(re.Val)

		// 拼接URL
		urls = append(urls, urlPre+re.Val)
	}

	// 获取翻译结果集
	results := getNet.TranslateDeepL(urls)

	for i, result := range results {
		if result.Err != nil {
			err = fw.WriteLine(res[i].Key + "此行报错！！！！！！！！！！！！！！！！！！！！！: " + result.Err.Error())
			if err != nil {
				fmt.Printf("第 %v 行的 %v 写入异常: %v\n", i, res[i].Key, err)
			}
			continue
		}

		if strings.Contains(result.Translate, "*QJK*") {
			result.Translate = strings.ReplaceAll(result.Translate, "*QJK*", "<br>")
		}

		// 修正翻译中的百分号问题，确保不会错误地将 "%%" 转换为 "%"
		if strings.Contains(result.Translate, "% ") {
			result.Translate = strings.ReplaceAll(result.Translate, "% ", "%% ")
			result.Translate = strings.ReplaceAll(result.Translate, "%%%", "%%")
		}
		err = fw.WriteLine(res[i].Key + result.Translate)
		if err != nil {
			fmt.Printf("第 %v 行的 %v 写入异常: %v\n", i, res[i].Key, err)
		}
	}
}

// YoudaoTranslate 使用有道云API翻译 需要对应的 appKey 和 appSecret
func YoudaoTranslate(res []NeedTrans) {
	filePath := "./File/zh_CN.lang"
	fw, err := NewFileWriter(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(fw *FileWriter) {
		if err != nil {
			fmt.Printf("关闭文件失败: %v\n", err)
		}
	}(fw)

	//遍历并且输出参数
	for _, entry := range res {
		if entry.Val == "" {
			if err = fw.WriteLine(entry.Key); err != nil {
				fmt.Printf("写入文件异常: %v\n", err)
			}
			fmt.Println(entry.Key)
			continue
		}
		result := youdaoyunAPI.TransYouDaoYun(entry.Val)
		if err = fw.WriteLine(entry.Key + result); err != nil {
			fmt.Printf("写入文件异常: %v\n", err)
		}
		fmt.Println(entry.Key + result)

		// 根据API限制调整睡眠时间 这里用 2 秒翻译一次
		time.Sleep(2000 * time.Millisecond)
	}
}
