package storage

import (
	"database/sql"
	"errors"
	"fmt"

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
	sql    *sql.DB
	stmt   *sql.Stmt
	buffer []BlogPost
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
		sql:    db,
		stmt:   stmt,
		buffer: make([]BlogPost, 0, 10)}, nil
}

// Add stores a blogPost into the buffer. Once the buffer is full, the
// blogPosts are flushed to the database.
func (s *SqliteStore) AddBlogPost(bp BlogPost) error {
	if len(s.buffer) == cap(s.buffer) {
		return errors.New("store buffer is full")
	}
	s.buffer = append(s.buffer, bp)
	if len(s.buffer) == cap(s.buffer) {
		if err := s.flush(); err != nil {
			return fmt.Errorf("failed to flush store buffer: %w", err)
		}
	}
	return nil

}

func (s *SqliteStore) flush() error {
	tx, err := s.sql.Begin()
	if err != nil {
		return err
	}

	for _, bp := range s.buffer {
		_, err := s.stmt.Exec(bp.Id, bp.Url, bp.Title, bp.Html, bp.Tags)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	s.buffer = s.buffer[:0]
	return tx.Commit()
}

func (s *SqliteStore) Close() error {
	defer func() {
		s.sql.Close()
		s.stmt.Close()
	}()

	if err := s.flush(); err != nil {
		return err
	}
	return nil
}
