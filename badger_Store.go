// https://dgraph.io/docs/badger/get-started/
package main

import (
	badger "github.com/dgraph-io/badger/v4"
)

type BadgerStore struct {
	database *badger.DB
}

func NewBadgerStore(file string) (*BadgerStore, error) {
	db, err := badger.Open(badger.DefaultOptions(file))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return &BadgerStore{database: db}, nil
}

func (s *BadgerStore) AddBlogPost(bp BlogPost) error {
	return s.database.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(bp.Url), []byte(bp.Html))
	})
}
