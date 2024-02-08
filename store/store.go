package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

// NewStore opens an existing sqlite3 database or creates a new one if it does not yet exists.
func NewStore(filename string) (*Store, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	s := &Store{db: db}
	err = s.Initialize()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Close closes the sqlite database.
func (s *Store) Close() error {
	return s.db.Close()
}

// Initialize sets database tables if needed.
func (s *Store) Initialize() error {
	statement, err := s.db.Prepare("CREATE TABLE IF NOT EXISTS plants (id INTEGER PRIMARY KEY, name TEXT, lighting TEXT, watering TEXT, fertilizing TEXT, toxicity TEXT, notes TEXT)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}
