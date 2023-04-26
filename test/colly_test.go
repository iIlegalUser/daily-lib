package test

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"testing"
)

func TestCrawl(t *testing.T) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com"),
	)

	c.OnHTML("#s-top-left a", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	_ = c.Visit("https://www.baidu.com/")
}
