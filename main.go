package main

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"log"
	"os"
	"strings"
)

func main() {

	//Instantiate default collector
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
		colly.CacheDir("./cache"),
	)

	c.OnHTML(".j-chapter-link", func(e *colly.HTMLElement) {
		log.Println(e.Attr("data-hreflink"))
		title := strings.Replace(e.Text, " ", "", -1)
		title = strings.Replace(title, "\n", "", -1)
		log.Println(title)
		if !strings.Contains(title, "第") {
			return
		}
		path := "./jinji/" + title
		if !Exists(path) {
			os.MkdirAll(path, os.ModePerm)
		}

		d := c.Clone()
		d.OnRequest(func(request *colly.Request) {
			request.Ctx.Put("path", path) //保存路径
		})
		d.OnHTML(".rd-article__pic", func(element *colly.HTMLElement) {
			imgurl := element.ChildAttr("img", "data-src")

			m := c.Clone()
			m.OnRequest(func(request *colly.Request) {
				request.Ctx.Put("img_file", element.Response.Ctx.Get("path")+"/"+element.Attr("data-page_id")) //保存路径
			})
			m.OnResponse(func(response *colly.Response) {
				if strings.Index(response.Headers.Get("Content-Type"), "image") > -1 {
					response.Save(response.Ctx.Get("img_file") + ".jpg")
					return
				}
			})
			imgurl = strings.Replace(imgurl, "page-800", "page-1200", 1)
			m.Visit(imgurl)
		})
		d.Visit("https://www.mkzhan.com/" + e.Attr("data-hreflink"))

	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.mkzhan.com/211879/")
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
