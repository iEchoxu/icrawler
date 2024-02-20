package task

import (
	_ "icrawler/bilibili"
	"icrawler/crawlers"
	"log"
	"sync"
)

func Run() {
	result := make(chan *crawlers.Result, crawlers.EnvConfig.MaxRequest)

	var wg sync.WaitGroup
	wg.Add(2)

	crawlerName := crawlers.EnvConfig.ProjectName
	crawler, err := crawlers.Get(crawlerName)
	if err != nil {
		log.Println(err)
	}

	go func() {
		defer wg.Done()
		if err = crawler.StartRequests(result); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		crawler.Parse(result)
	}()

	wg.Wait()
}
