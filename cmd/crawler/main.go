package main

import (
	"github.com/DrewCyber/crawler-go/internal/crawler"
	"github.com/DrewCyber/crawler-go/internal/scraper"
	storage "github.com/DrewCyber/crawler-go/internal/sqlite_storage"
)

func main() {
	store, err := storage.NewSqliteStore("./blog.db")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	scraper := scraper.NewCollector(store)

	crawler := crawler.NewCrawler(scraper, store)

	// open the target URL
	crawler.Start("https://evo-lutio.livejournal.com/1699420.html")
	// collector.Visit("https://evo-lutio.livejournal.com/2024/07/14/")

}
