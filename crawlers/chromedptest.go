package crawlers

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"icrawler/crawlers/request"
	"log"
	"time"
)

func ResponseWithoutCheckElement(url, rootXpath string, wb *request.Browser) string {
	var resultHTML string

	tabCtx, cancel, err := wb.NewTab(wb.BrowserCtx)
	if err != nil {
		log.Println(err)
	}
	defer cancel()

	if err := chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		// 如果没有这个元素会导致程序卡住
		chromedp.OuterHTML(rootXpath, &resultHTML, chromedp.ByJSPath),
		chromedp.Sleep(2*time.Second), // TODO 通过配置文件设置
	); err != nil {
		fmt.Printf("unable to navigate to page: %v", err)
	}

	return resultHTML
}

// https://github.com/chromedp/chromedp/issues/120
// https://github.com/chromedp/chromedp/issues/568#issuecomment-836365750

func CheckElementIsExistWithNodes() {
	url := "https://www.bilibili.com/v/popular/rank/all"

	wb := request.NewExecAllocator(context.Background(), request.Opts)
	var resultHTML string

	tabCtx, cancel, err := wb.NewTab(wb.BrowserCtx)
	if err != nil {
		log.Println(err)
	}
	defer cancel()

	var nodes []*cdp.Node

	// 检测页面中元素是否存在，当 nodes 长度为 0 表示没有该元素
	// 该方式会导致地址栏里的 url 跳转两次
	chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		//chromedp.EvaluateAsDevTools(fmt.Sprintf(`document.querySelector("%v")`, visibleEleId), ""),
		chromedp.Nodes(`//*[@class="rank-list"]`, &nodes, chromedp.AtLeast(0)),
	)

	// 当页面中没有该元素时不执行后面的任务
	if len(nodes) == 0 {
		return
	}

	fmt.Println("node 数量为:", len(nodes))

	rootXpath := "document.getElementsByClassName('rank-list')"

	// 输出 HTML 内容
	if err := chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		chromedp.OuterHTML(rootXpath, &resultHTML, chromedp.ByJSPath),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		fmt.Printf("unable to navigate to page: %v", err)
	}

	fmt.Println(resultHTML)
}

func CheckElementIsExistWithJs() {
	url := "https://www.bilibili.com/v/popular/rank/all"

	wb := request.NewExecAllocator(context.Background(), request.Opts)
	var resultHTML string

	tabCtx, cancel, err := wb.NewTab(wb.BrowserCtx)
	if err != nil {
		log.Println(err)
	}
	defer cancel()

	rootXpath := "document.getElementsByClassName('rank-list')"

	jsText := request.GetElementValueByJS(rootXpath)

	if err := chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		chromedp.EvaluateAsDevTools(jsText, &resultHTML),
		//chromedp.EvaluateAsDevTools(`
		//    var sources = document.getElementsByClassName('rank-list');
		//    if (sources.length > 0) {
		//        sources[0].outerHTML;
		//    } else {
		//        '未找到标签';
		//    }
		//`, &resultHTML),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		fmt.Printf("unable to navigate to page: %v", err)
	}

	fmt.Println(resultHTML)
}

func Task(url string) chromedp.Tasks {
	var title, titleHTMLContent string

	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.EvaluateAsDevTools(`
		   var sources = document.getElementById('mirror-vdcon');
		   if (sources != null) {
		       sources;
		   } else {
		       '未找到<source>标签';
		   }
		`, &titleHTMLContent),
		chromedp.Title(&title),
		chromedp.OuterHTML(EnvConfig.StartUrl.HTMLContentXpath, &titleHTMLContent, chromedp.ByJSPath),
	}
}
