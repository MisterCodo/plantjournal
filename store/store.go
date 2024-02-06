package store

import (
	"context"
	"database/sql"
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

	// TODO: Delete this code.
	_, err = s.CreatePlant(context.Background(), &Plant{Name: "Spider", Lighting: "Low", Watering: "Let soil dry", Fertilizing: "Rarely", Toxicity: "Oh no", Notes: "Test"})
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
