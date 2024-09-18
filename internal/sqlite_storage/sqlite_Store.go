package storage

import (
	"database/sql"

	"github.com/DrewCyber/crawler-go/internal/crawler"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schemaSQL = `
CREATE TABLE IF NOT EXISTS blog (
  id INTEGER PRIMARY KEY,
  url TEXT,
  title TEXT,
  html TEXT,
  tags TEXT
);
`
	insertSQL = `INSERT INTO blog (id, url, title, html, tags) VALUES (?, ?, ?, ?, ?);`
)

type SqliteStore struct {
	sql  *sql.DB
	stmt *sql.Stmt
}

func NewSqliteStore(file string) (*SqliteStore, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

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
	_, err := s.stmt.Exec(blogPost.Id, blogPost.Url, blogPost.Title, blogPost.Html, blogPost.Tags)
	return err
}

func (s *SqliteStore) Close() {
	s.sql.Close()
	s.stmt.Close()
}
