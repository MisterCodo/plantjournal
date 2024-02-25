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
	initDatabase := `
	CREATE TABLE IF NOT EXISTS plants (
		id INTEGER PRIMARY KEY,
		name TEXT,
		lighting TEXT,
		watering TEXT,
		fertilizing TEXT,
		toxicity TEXT,
		notes TEXT
	);
	CREATE TABLE IF NOT EXISTS actions (
		id INTEGER PRIMARY KEY,
		plant_id INTEGER NOT NULL,
	    day TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d', 'now')),
		watered INTEGER CHECK(watered IN(0, 1)),
		fertilized INTEGER CHECK(fertilized IN(0, 1)),
		notes TEXT,
		FOREIGN KEY (plant_id) REFERENCES plants(id)
	);`

	_, err := s.db.Exec(initDatabase)
	if err != nil {
		return err
	}

	return nil
}
