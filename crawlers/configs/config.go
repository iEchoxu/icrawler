package configs

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

var (
	seeds *Config
	once  sync.Once
)

type Config struct {
	ProjectName  string                 `json:"project_name"`
	SpiderName   string                 `json:"spider_name"`
	StartUrl     *StartUrl              `json:"start_urls"`
	Item         string                 `json:"item"`
	Selector     map[string]interface{} `json:"selectors"`
	MaxRequest   int                    `json:"max_request"`
	TimeoutCount time.Duration          `json:"timeout"`
}

type StartUrl struct {
	Url              string `json:"url"`
	Follow           bool   `json:"follow"`
	HTMLContentXpath string `json:"html_content_xpath"`
	RestrictXpath    string `json:"restrict_xpath"`
	LinksXpath       string `json:"links_xpath"`
}

func Load() *Config {
	once.Do(func() {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalln("打开配置文件出错", err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		if err = json.NewDecoder(file).Decode(&seeds); err != nil {
			log.Fatalln("读取 json 文件出错", err)
		}
	})

	return seeds
}

func (c *Config) Selectors(key string) (value string) {
	for i, v := range c.Selector {
		if i == key {
			value = v.(string)
		}
	}
	return
}
