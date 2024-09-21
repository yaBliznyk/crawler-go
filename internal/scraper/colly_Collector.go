package scraper

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"

	"github.com/DrewCyber/crawler-go/internal/crawler"
)

type Collector struct {
	cc  *colly.Collector
	svc crawler.Service
}

func NewCollector(svc crawler.Service) *Collector {
	cc := colly.NewCollector(
		colly.AllowedDomains("evo-lutio.livejournal.com"),
		colly.CacheDir("./colly_cache"),
		colly.URLFilters(
			regexp.MustCompile(`https://evo-lutio.livejournal.com/[\d]+.html$`), //https://user.livejournal.com/1697758.html
			regexp.MustCompile(`https://evo-lutio.livejournal.com(/[\d]+)+/$`),  //https://user.livejournal.com/2024/09/07/
		),
		// colly.Async(true),
	)
	return &Collector{
		cc:  cc,
		svc: svc,
	}
}

func (c *Collector) Init() {
	c.cc.Limit(&colly.LimitRule{
		Parallelism: 1,
		// RandomDelay: 5 * time.Second,
	})

	c.cc.OnResponse(func(r *colly.Response) {
		// fmt.Println("HTML string received: ", string(r.Body))
	})

	c.cc.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, "failed with response: ", r, "\nError: ", err)
	})

	// c.OnRequest(func(r *colly.Request) {
	// On every request, need to call Chromedp to handle the JavaScript.
	// It will load comments section
	// }

	// Parse single blog post
	c.cc.OnHTML(".b-singlepost", func(e *colly.HTMLElement) {
		blogPost := crawler.BlogPost{}
		blogPost.Title = e.ChildText("h1.b-singlepost-title")
		blogPost.Html, _ = e.DOM.Find("article.b-singlepost-body").Html()
		blogPost.DateTime = e.ChildText(".b-singlepost-author-date") // BUG. Why it's doubled? "2024-09-07 19:07:002024-09-07 19:07:00"
		// blogPost.DateTime = e.ChildText(".b-singlepost-author-userinfo-screen")
		blogPost.Tags = e.ChildText(".b-singlepost-tags-items") // Need to iterate?
		blogPost.URL = e.Request.URL.String()
		var err error
		blogPost.ID, err = GetPostIdFromUrl(blogPost.URL)
		if err != nil {
			fmt.Printf("Failed to get post ID: %s", err)
			panic(err)
		}
		err = c.svc.CreateBlogPost(blogPost)
		if err != nil {
			fmt.Printf("Failed to create blog post ID: %s", err)
			panic(err)
		}
	})

	// Get child urls
	c.cc.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, a *colly.HTMLElement) {
			// fmt.Printf("Found the URL: %s\n", a.Attr("href"))
			for _, pattern := range c.cc.URLFilters {
				if pattern.MatchString(a.Attr("href")) {
					_ = c.svc.CreateUrl(a.Attr("href"))
					fmt.Printf("Found the URL: %s\n", a.Attr("href"))
				}
			}
		})
	})
}

func (c *Collector) Start(url string) error {
	err := c.cc.Visit(url)
	if err != nil {
		return fmt.Errorf("cc.Visit: %w", err)
	}
	c.cc.Wait()
	err = c.svc.MarkUrlAsVisited(url)
	if err != nil {
		return fmt.Errorf("svc.MarkUrlAsVisited: %w", err)
	}

	for {
		nextUrl, err := c.svc.GetNextUrl(url)
		if err != nil {
			return fmt.Errorf("cc.Visit: %w", err)
		}
		if nextUrl == "" {
			return nil
		}

		err = c.cc.Visit(nextUrl)
		if err != nil {
			return fmt.Errorf("cc.Visit: %w", err)
		}
		c.cc.Wait()

		err = c.svc.MarkUrlAsVisited(nextUrl)
	}
}
