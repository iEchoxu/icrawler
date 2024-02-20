package item

import (
	"icrawler/crawlers"
)

var (
	titleXpath         = crawlers.EnvConfig.Selectors("title")
	AuthorXpath        = crawlers.EnvConfig.Selectors("author")
	ViewCountsXpath    = crawlers.EnvConfig.Selectors("view_counts")
	DanmuCountsXpath   = crawlers.EnvConfig.Selectors("danmu_counts")
	VideoLikeXpath     = crawlers.EnvConfig.Selectors("video_like")
	VideoCoinXpath     = crawlers.EnvConfig.Selectors("video_coin")
	VideoFavoriteXpath = crawlers.EnvConfig.Selectors("video_favorite")
	VideoShareXpath    = crawlers.EnvConfig.Selectors("video_share")
)
