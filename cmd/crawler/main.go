package main

import (
	crawler "github.com/DrewCyber/crawler-go/internal/crawler"
	storage "github.com/DrewCyber/crawler-go/internal/storage"
)

func main() {
	store, err := storage.NewSqliteStore("../../blog.db")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	collector := crawler.NewCollector(store)

	// open the target URL
	collector.Visit("https://evo-lutio.livejournal.com/1699420.html")
	// collector.Visit("https://evo-lutio.livejournal.com/2024/07/14/")
	collector.Wait()

}
