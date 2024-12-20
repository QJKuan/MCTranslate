package getNet

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
)

// Trans 存储结果
type Trans struct {
	Index     int
	Translate string
	Err       error
}

func TranslateDeepL(urls []string) []Trans {
	// 创建一个浏览器上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 记录开始时间
	//startTime := time.Now().Unix()

	// 定义要获取文本的节点信息
	nodeInfo := "/html[1]/body[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[2]/div[1]/main[1]/div[2]/nav[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/section[1]/div[1]/div[2]/div[3]/section[1]/div[1]/d-textarea[1]/div[1]/p[1]/span[1]"

	var results []Trans
	// 遍历每个URL并执行任务
	for i, url := range urls {
		var res Trans
		res.Index = i

		//TODO 处理现阶段无法翻译的字符
		if url == "" {
			// 判空返回
			res.Translate = ""
			results = append(results, res)
			continue
		}
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Text(nodeInfo, &res.Translate, chromedp.BySearch),
		)
		if err != nil {
			res.Err = err
			fmt.Printf("第 %v 行翻译异常: %v\n", i, err)
			continue
		}
		fmt.Printf("翻译结果为: %s\n", res.Translate)
		results = append(results, res)
	}
	return results
	// 计算并输出执行时间
	//fmt.Printf("执行时间: %d 秒\n", time.Now().Unix()-startTime)
}
