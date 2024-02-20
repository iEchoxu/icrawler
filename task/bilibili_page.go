package task

import (
	"context"
	"fmt"
	"github.com/antchfx/htmlquery"
	"icrawler/crawlers"
	"icrawler/crawlers/request"
	"log"
	"strings"
)

func BilibiliPage() {
	url := "https://www.bilibili.com/video/BV1Sv421k7pS"
	//url := "https://www.bilibili.com/video/BV18Z421m7yX"
	wb := request.NewExecAllocator(context.Background(), request.Opts)
	rootXpath := "document.querySelectorAll('#mirror-vdcon')[0]"
	response := request.Navigate(url, rootXpath, wb)
	HTMLContent, err := htmlquery.Parse(strings.NewReader(response))
	if err != nil {
		log.Println(err)
	}
	videoShareXpath := crawlers.EnvConfig.Selectors("video_share")
	videoShare := strings.TrimSpace(htmlquery.FindOne(HTMLContent, videoShareXpath).Data)

	AuthorXpath := crawlers.EnvConfig.Selectors("author")
	author, err := htmlquery.Query(HTMLContent, AuthorXpath)
	if err != nil {
		log.Println(err)
	}

	var authorResult string

	if author != nil {
		authorResult = strings.TrimSpace(author.Data)
	} else {
		authorXpathNew := "//div[@class='up-panel-container']//span[contains(text(),'UP主')]//../a/text()"
		authorResult = strings.TrimSpace(htmlquery.FindOne(HTMLContent, authorXpathNew).Data)
	}

	fmt.Println("转发量为:", videoShare, "UP 主为:", authorResult)

	err = wb.CloseBrowser(wb.BrowserCtx)
	if err != nil {
		log.Fatalln(err)
	}

}
