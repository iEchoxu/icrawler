package bilibili

import (
	_ "icrawler/bilibili/item"
	"icrawler/crawlers"
	"log"
)

type bilibili struct {
	crawlers.Crawler
}

func init() {
	bilibiliCrawler := New()
	crawlers.Register("bilibili", bilibiliCrawler)
	log.Println("已注册 bilibili crawler")
}

func New() crawlers.ICrawler {
	return &bilibili{}
}

func (b *bilibili) StartRequests(results chan *crawlers.Result) error {
	err := b.GetLinks(results)
	if err != nil {
		return err
	}
	return nil
}

func (b *bilibili) Parse(results chan *crawlers.Result) {
	item := crawlers.EnvConfig.Item
	itemInstance, err := crawlers.GetItem(item)
	if err != nil {
		log.Println(err)
	}
	b.ParseItemOfLinks(itemInstance, results)

}
