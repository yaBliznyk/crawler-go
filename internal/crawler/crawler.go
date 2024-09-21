package crawler

// DOMAIN
type BlogPost struct {
	ID       int
	URL      string
	DateTime string
	Title    string
	Html     string
	Tags     string
	Comments []string
}

func (bp *BlogPost) AddComments(comments []string) error {
	bp.Comments = append(bp.Comments, comments...)
	return nil
}

// APPLICATION
type Repo interface {
	CreateBlogPost(*BlogPost) error
	GetBlogPostByID(id int) (*BlogPost, error)
	UpdateBlogPost(*BlogPost) error
}

type Service interface {
	CreateBlogPost(BlogPost) error
	SaveBlogComments(id int, comments []string)
}

// ----------------------

// Crawler gets pages recursively and queues them
type Crawler struct {
	repo Repo // Storage stores data
}

func NewCrawler(repo Repo) *Crawler {
	return &Crawler{
		repo: repo,
	}
}
