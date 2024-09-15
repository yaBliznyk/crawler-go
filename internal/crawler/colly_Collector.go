package crawler

import (
	"fmt"
	"regexp"

	"github.com/DrewCyber/crawler-go/internal/storage"
	"github.com/gocolly/colly"
)

var err error

func NewCollector(store storage.Repo) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("evo-lutio.livejournal.com"),
		colly.CacheDir("./colly_cache"),
		colly.URLFilters(
			regexp.MustCompile(`https://evo-lutio.livejournal.com/[\d]+.html$`), //https://user.livejournal.com/1697758.html
			regexp.MustCompile(`https://evo-lutio.livejournal.com(/[\d]+)+/$`),  //https://user.livejournal.com/2024/09/07/
		),
		// colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 1,
		// RandomDelay: 5 * time.Second,
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("HTML string received: ", string(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, "failed with response: ", r, "\nError: ", err)
	})

	// c.OnRequest(func(r *colly.Request) {
	// On every request, need to call Chromedp to handle the JavaScript.
	// It will load comments section
	// }

	// Parse single blog post
	c.OnHTML(".b-singlepost", func(e *colly.HTMLElement) {
		blogPost := storage.BlogPost{}
		blogPost.Title = e.ChildText("h1.b-singlepost-title")
		blogPost.Html, _ = e.DOM.Find("article.b-singlepost-body").Html()
		blogPost.DateTime = e.ChildText(".b-singlepost-author-date") // BUG. Why it's doubled? "2024-09-07 19:07:002024-09-07 19:07:00"
		// blogPost.DateTime = e.ChildText(".b-singlepost-author-userinfo-screen")
		blogPost.Tags = e.ChildText(".b-singlepost-tags-items") // Need to iterate?
		blogPost.Url = e.Request.URL.String()
		blogPost.Id, err = GetPostIdFromUrl(blogPost.Url)
		if err != nil {
			fmt.Printf("Failed to get post ID: %s", err)
			panic(err)
		}

		store.AddBlogPost(blogPost)
	})

	// Get child urls
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, a *colly.HTMLElement) {
			// fmt.Printf("Found the URL: %s\n", a.Attr("href"))
			for _, pattern := range c.URLFilters {
				re, err := regexp.Compile(pattern.String())
				if err != nil {
					fmt.Printf("Failed to compile regexp: %s\n", err)
					continue
				}
				if re.MatchString(a.Attr("href")) {
					fmt.Printf("Found the URL: %s\n", a.Attr("href"))
				}
			}

		})
	})

	return c
}
