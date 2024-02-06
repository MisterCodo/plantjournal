package store

import "database/sql"

type Store struct {
	DB *sql.DB
}

func NewStore(filename string) (*Store, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	// TODO: Queries to create tables if they don't exist.
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS plants (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return nil, err
	}

	store := &Store{DB: db}
	return store, nil
}

// Close closes the sqlite database.
func (s *Store) Close() error {
	return s.DB.Close()
}
