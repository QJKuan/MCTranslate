package getNet

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"time"
)

func GetTranslate() {
	// 创建一个浏览器上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 记录开始时间
	startTime := time.Now().Unix()

	// 定义要访问的URL列表
	urls := []string{
		"https://www.deepl.com/zh/translator#en/zh-hans/Construct of Terror",
		"https://www.deepl.com/zh/translator#en/zh-hans/DeathRay Shot",
		"https://www.deepl.com/zh/translator#en/zh-hans/Dungeon Keeper",
		"https://www.deepl.com/zh/translator#en/zh-hans/Explosives Expert",
		"https://www.deepl.com/zh/translator#en/zh-hans/Non-Lethal Attacks Create Lightning on Impact",
		"https://www.deepl.com/zh/translator#en/zh-hans/Launches The Holder in The Direction They're Facing",
		"https://www.deepl.com/zh/translator#en/zh-hans/Only those who are real artists can wear this armor.",
		"https://www.deepl.com/zh/translator#en/zh-hans/Dungeon Keeper say You will get a return crystal upon entry This crystal will let you end your run at any time.",
	}

	// 定义要获取文本的节点信息
	nodeInfo := "/html[1]/body[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[2]/div[1]/main[1]/div[2]/nav[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/section[1]/div[1]/div[2]/div[3]/section[1]/div[1]/d-textarea[1]/div[1]/p[1]/span[1]"

	// 存储结果
	//var results []string
	type result struct {
		Index     int
		Translate string
		Err       error
	}

	var results []result
	// 遍历每个URL并执行任务
	for i, url := range urls {
		var res result
		res.Index = i
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Text(nodeInfo, &res.Translate, chromedp.BySearch),
		)
		if err != nil {
			res.Err = err
			//fmt.Printf("第 %v 行翻译异常: %v\n", i, err)
			continue
		}
		fmt.Printf("翻译结果为: %s\n", res.Translate)
		results = append(results, res)
	}

	// 计算并输出执行时间
	fmt.Printf("执行时间: %d 秒\n", time.Now().Unix()-startTime)
}
