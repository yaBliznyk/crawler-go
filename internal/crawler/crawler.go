package crawler

type BlogPost struct {
	Id       int
	Url      string
	DateTime string
	Title    string
	Html     string
	Tags     string
}

// Crawler gets pages recursively and queues them
type Crawler struct {
	scraper Scraper // Scraper gets pages, extract new urls and blogPost from them
	storage Storage // Storage stores data
}

type Scraper interface {
	Visit(string) error
	// Parse(HTMLElement)
	Wait()
}

type Storage interface {
	AddBlogPost(BlogPost) error
	Close()
}

func NewCrawler(scraper Scraper, storage Storage) *Crawler {
	return &Crawler{
		scraper: scraper,
		storage: storage,
	}
}

func (c *Crawler) Start(url string) error {
	c.scraper.Visit(url)
	c.scraper.Wait()
	return nil
}
