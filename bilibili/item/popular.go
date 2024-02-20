package item

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"icrawler/crawlers"
	"log"
	"strings"
)

func init() {
	popularRankItem := New()
	crawlers.ItemRegister("PopularRank", popularRankItem)
	log.Println("已注册 popular item")
}

type PopularRank struct {
	Title       string
	Link        string
	Poster      string
	Author      string
	ViewCounts  string
	DanMuCounts string
}

func New() crawlers.Item {
	return &PopularRank{}
}

func (p *PopularRank) ParseItem(results chan *crawlers.Result) {
	successCount := 0
	for i := range results {
		i := i
		//fmt.Println(i.URL)
		HTMLContent, err := htmlquery.Parse(strings.NewReader(*i.HTMLContent))
		if err != nil {
			log.Println(err)
		}
		title := htmlquery.FindOne(HTMLContent, titleXpath).Data
		// https://www.bilibili.com/video/BV1Wm411Q7vL/ 这个会导致 author 为空
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
			log.Println("重新获取到了 UP 主信息。。。")
		}

		viewCounts := strings.TrimSpace(htmlquery.FindOne(HTMLContent, ViewCountsXpath).Data)
		danMuCounts := strings.TrimSpace(htmlquery.FindOne(HTMLContent, DanmuCountsXpath).Data)
		videoLike := strings.TrimSpace(htmlquery.FindOne(HTMLContent, VideoLikeXpath).Data)
		videoCoin := strings.TrimSpace(htmlquery.FindOne(HTMLContent, VideoCoinXpath).Data)
		videoFav := strings.TrimSpace(htmlquery.FindOne(HTMLContent, VideoFavoriteXpath).Data)
		videoShare := strings.TrimSpace(htmlquery.FindOne(HTMLContent, VideoShareXpath).Data)
		fmt.Printf(
			"%s 获取到的 Title 为: %s, UP 主为: %s,播放量为: %s,弹幕量为: %s,点赞数为: %s,投币数为: %s,收藏量为: %s,转发量为: %s\n",
			i.URL, title, authorResult, viewCounts, danMuCounts, videoLike, videoCoin, videoFav, videoShare,
		)

		// TODO 调用存储接口
		//p.Save()
		successCount++
	}
	fmt.Println("总共执行了:", successCount)
}
