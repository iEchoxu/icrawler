package crawlers

import (
	"icrawler/crawlers/configs"
)

var (
	EnvConfig               = configs.Load()
	StartUrlsLink           = EnvConfig.StartUrl.Url
	StartUrlsRootXpath      = EnvConfig.StartUrl.HTMLContentXpath
	StartUrlsLinksNodeXpath = EnvConfig.StartUrl.RestrictXpath
	StartUrlsLinksXpath     = EnvConfig.StartUrl.LinksXpath
	ItemRootXpath           = EnvConfig.Selectors("root")
)

type Result struct {
	URL         string
	HTMLContent *string
}

var crawlerMap = make(map[string]ICrawler)

var ItemMap = make(map[string]Item)
