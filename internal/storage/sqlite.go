package storage

import (
	"database/sql"

	"github.com/DrewCyber/crawler-go/internal/crawler"
)

const (
	insertSQL = `INSERT INTO blog (id, url, title, html, tags) VALUES (?, ?, ?, ?, ?);`
)

type SqliteStore struct {
	sql  *sql.DB
	stmt *sql.Stmt
}

func NewRepo(db *sql.DB) (*SqliteStore, error) {
	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return nil, err
	}

	return &SqliteStore{
		sql:  db,
		stmt: stmt,
	}, nil
}

func (s *SqliteStore) AddBlogPost(blogPost crawler.BlogPost) error {
	_, err := s.stmt.Exec(blogPost.ID, blogPost.URL, blogPost.Title, blogPost.Html, blogPost.Tags)
	return err
}

func (s *SqliteStore) Close() {
	s.sql.Close()
	s.stmt.Close()
}
