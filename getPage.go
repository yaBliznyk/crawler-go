package main

func main() {
	store, err := NewSqliteStore("./blog.db")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	collector := NewCollector(store)

	// open the target URL
	collector.Visit("https://evo-lutio.livejournal.com/1699420.html")
	// collector.Visit("https://evo-lutio.livejournal.com/2024/07/14/")
	collector.Wait()

}
