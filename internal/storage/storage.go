package storage

type BlogPost struct {
	Id       int
	Url      string
	DateTime string
	Title    string
	Html     string
	Tags     string
}
type Repo interface {
	AddBlogPost(BlogPost) error
	Close() error
}
