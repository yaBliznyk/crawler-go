package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/DrewCyber/crawler-go/internal/crawler"
	"github.com/DrewCyber/crawler-go/internal/scraper"
	"github.com/DrewCyber/crawler-go/internal/storage"
)

func main() {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo, err := storage.NewRepo(db)
	if err != nil {
		panic(err)
	}

	svc := crawler.NewCrawler(repo)

	scraper := scraper.NewCollector(svc)
	scraper.Init()

	// open the target URL
	scraper.Start("https://evo-lutio.livejournal.com/1699420.html")
	// collector.Visit("https://evo-lutio.livejournal.com/2024/07/14/")
}
