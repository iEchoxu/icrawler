package crawlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/sync/errgroup"
	"icrawler/crawlers/request"
	"log"
	"strings"
)

func Register(name string, crawler ICrawler) {
	if _, exists := crawlerMap[name]; exists {
		log.Fatalln(name, "crawler already registered")
	}

	crawlerMap[name] = crawler
}

func Get(name string) (ICrawler, error) {
	if _, ok := crawlerMap[name]; !ok {
		return nil, errors.New("暂不支持该网站")
	}
	return crawlerMap[name], nil
}

func (c *Crawler) GetStartUrlResponse(wb *request.Browser) (response []string) {
	startURLResponse := request.Response(StartUrlsLink, StartUrlsRootXpath, wb)

	HTMLContent, err := htmlquery.Parse(strings.NewReader(startURLResponse))
	if err != nil {
		log.Println(err)
	}

	nodes := htmlquery.Find(HTMLContent, StartUrlsLinksNodeXpath)

	for _, v := range nodes {
		link := "https:" + strings.TrimSpace(htmlquery.InnerText(htmlquery.FindOne(v, StartUrlsLinksXpath)))
		response = append(response, link)
	}

	return
}

func (c *Crawler) GetLinks(result chan *Result) error {
	wb := request.NewExecAllocator(context.Background(), request.Opts)
	response := c.GetStartUrlResponse(wb)

	defer close(result)

	g := new(errgroup.Group)
	g.SetLimit(EnvConfig.MaxRequest)
	fmt.Println("总链接数据:", len(response))

	for _, link := range response {
		link := link
		g.Go(func() error {
			// 传入的 selectors 必须有 .length 方法，否则结果集为空
			singleVideoResponse := request.Response(link, ItemRootXpath, wb)
			// 没有获取到数据的 link 不添加进 channel
			if singleVideoResponse == "未找到标签" {
				log.Println(link, "未获取到数据...")
				return nil
			}
			result <- &Result{
				URL:         link,
				HTMLContent: &singleVideoResponse,
			}
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	err = wb.CloseBrowser(wb.BrowserCtx)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (c *Crawler) ParseItemOfLinks(item Item, Result chan *Result) {
	item.ParseItem(Result)
}

func ItemRegister(name string, item Item) {
	ItemMap[name] = item
}

func GetItem(name string) (Item, error) {
	if _, ok := ItemMap[name]; !ok {
		return nil, errors.New("暂不支持该网站")
	}
	return ItemMap[name], nil
}
