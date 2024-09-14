package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type BlogPost struct {
	Id       int
	Url      string
	DateTime string
	Title    string
	Html     string
	Tags     string
}

func main() {
	var blog []BlogPost
	store, err := NewSqliteStore("./blog.db")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// instantiate a new collector object
	c := colly.NewCollector(
		colly.AllowedDomains("evo-lutio.livejournal.com"),
		colly.CacheDir("./colly_cache"),
		// colly.URLFilters(
		// 	regexp.MustCompile(`https://evo-lutio.livejournal.com/[\d]+.html$`),
		// ),
		// colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 1,
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("HTML string received: ", string(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, "failed with response: ", r, "\nError: ", err)
	})

	// Parse single blog post
	c.OnHTML(".b-singlepost", func(e *colly.HTMLElement) {
		blogPost := BlogPost{}
		blogPost.Title = e.ChildText("h1.b-singlepost-title")
		blogPost.Html, _ = e.DOM.Find("article.b-singlepost-body").Html()
		blogPost.DateTime = e.ChildText(".b-singlepost-author-date") // Why it's doubled? "2024-09-07 19:07:002024-09-07 19:07:00"
		// blogPost.DateTime = e.ChildText(".b-singlepost-author-userinfo-screen")
		blogPost.Tags = e.ChildText(".b-singlepost-tags-items") // Need to iterate
		blogPost.Url = e.Request.URL.String()
		blogPost.Id, err = GetPostId(blogPost.Url)
		if err != nil {
			fmt.Printf("Failed to get post ID: %s", err)
			panic(err)
		}
		blog = append(blog, blogPost)
	})

	// open the target URL
	c.Visit("https://evo-lutio.livejournal.com/1699420.html")
	c.Wait()

	for i, blogPost := range blog {
		fmt.Println("Post html >>>:", i, blogPost.Html)
		fmt.Println("Post title >>>:", i, blogPost.Title)
		fmt.Println("Post url >>>:", i, blogPost.Url)
		fmt.Println("Post date >>>:", i, blogPost.DateTime)
		fmt.Println("Post tags >>>:", i, blogPost.Tags)
		store.AddBlogPost(blogPost)
	}

}
