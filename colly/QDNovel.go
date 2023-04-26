package colly

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

type Novel struct {
	category   string
	rank       int
	title      string
	author     string
	words      float64
	tags       []string
	url        string
	lastUpdate string
}

func (n Novel) String() string {
	return fmt.Sprintf("排名: %d\r\n"+
		"标题: %s\r\n"+
		"作者: %s\r\n"+
		"字数: %.2f万字\r\n"+
		"标签: %s\r\n"+
		"链接: %s\r\n"+
		"更新时间: %s\r\n", n.rank, n.title, n.author, n.words, n.tags, n.url, n.lastUpdate)
}

func Crawl() {
	novels := make([]*Novel, 0, 100)

	c1 := colly.NewCollector(
		colly.AllowedDomains("www.qidian.com", "book.qidian.com"),
		colly.Async(true),
	)

	// 获取所有分类链接
	c1.OnHTML("#classify-list a", func(e *colly.HTMLElement) {
		category := e.Attr("title")
		href := e.Attr("href")
		if href == "" {
			return
		}
		c2 := c1.Clone()

		// 获取榜单书籍链接
		c2.OnHTML("div.popular-serial + div li", func(e *colly.HTMLElement) {
			rank, _ := strconv.Atoi(strings.TrimSpace(e.Attr("data-rid")))
			infoUrl := e.Request.AbsoluteURL(e.ChildAttr("a.link, a.name", "href"))
			c3 := c1.Clone()

			// 获取书籍详细信息
			c3.OnHTML("div.book-info", func(e *colly.HTMLElement) {
				title := e.ChildText("h1 em")
				author := e.ChildText("h1 a.writer")
				lastUpdate := string([]rune(e.ChildText("h1 span.book-update-time"))[5:])
				tags := make([]string, 0)
				tags = append(tags, e.ChildText("p.tag > a:nth-child(4)"), e.ChildText("p.tag > a:nth-child(5)"))
				if tag := e.ChildText("p.tag > a:nth-child(6)"); tag != "" {
					tags = append(tags, tag)
				}
				words, _ := strconv.ParseFloat(e.ChildText("p:nth-child(4) > em:nth-child(1)"), 64)

				novel := &Novel{
					category:   category,
					rank:       rank,
					title:      title,
					author:     author,
					words:      words,
					tags:       tags,
					url:        infoUrl,
					lastUpdate: lastUpdate,
				}
				//fmt.Println(novel)
				novels = append(novels, novel)
			})
			_ = c3.Visit(infoUrl)
			// 因为是异步，所以需要等待完成
			c3.Wait()
		})
		_ = c2.Visit(e.Request.AbsoluteURL(href))
		c2.Wait()
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("visiting ", r.Request.URL, "failed: ", err)
	})

	_ = c1.Visit("https://www.qidian.com/")
	c1.Wait()

	//将结果写入文件
	file, err := os.OpenFile("res.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0744)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_, _ = file.WriteString(time.Now().Format("2006-01-02 15:04:05") + "\r\n\n")
	if err != nil {
		return
	}
	if err != nil {
		fmt.Println(err)
	}
	for _, novel := range novels {
		if novel.rank%10 == 1 {
			_, _ = file.WriteString("---------------------" + novel.category + "---------------------\r\n")
		}
		_, _ = file.WriteString(novel.String() + "\r\n")
	}

}
