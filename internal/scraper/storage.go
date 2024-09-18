package scraper

import "github.com/DrewCyber/crawler-go/internal/crawler"

type Repo interface {
	AddBlogPost(crawler.BlogPost) error
	Close()
}
