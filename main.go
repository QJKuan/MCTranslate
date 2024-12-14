package main

import (
	"HanHua/youdaoyunAPI"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Result struct {
	Key string
	Val string
}

func main() {
	file, err := os.Open("./File/en_US.lang")
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

	var res []Result
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//分开后的值的对象 key为 "=" 前的值 val为等于号后的值
		entry := parseLine(line)
		if entry != nil {
			res = append(res, *entry)
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Printf("读取文件异常: %v", err)
		return
	}

	translateEntries(res)
}

func parseLine(line string) *Result {
	split := strings.SplitN(line, "=", 2)
	if len(split) < 2 {
		// 处理没有 "=" 的情况，可以选择跳过或记录错误
		return nil
	}
	return &Result{Key: split[0] + "=", Val: split[1]}
}

func translateEntries(res []Result) {
	out, err := os.Create("./File/zh_CNccc.lang")
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

	for _, entry := range res {
		if entry.Val == "" {
			writer.WriteString(entry.Key + "\n")
			fmt.Println(entry.Key)
			continue
		}
		result := youdaoyunAPI.TransYouDaoYun(entry.Val)
		_, err := writer.WriteString(entry.Key + result + "\n")
		if err != nil {
			fmt.Printf("写入文件异常: %v\n", err)
		}
		fmt.Println(entry.Key + result)

		// 根据API限制调整睡眠时间
		time.Sleep(2000 * time.Millisecond)
	}
}
