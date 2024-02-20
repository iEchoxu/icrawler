package crawlers

type ICrawler interface {
	StartRequests(chan *Result) error
	Parse(chan *Result)
}

type Item interface {
	ParseItem(chan *Result)
}

type Crawler struct {
	Name     string
	StartURL string
	Item     Item
}
