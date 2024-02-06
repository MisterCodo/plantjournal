package store

import (
	"context"
)

// Plant contains the details of a plant.
type Plant struct {
	ID          int32
	Name        string
	Lighting    string
	Watering    string
	Fertilizing string
	Toxicity    string
	Notes       string
	// TODO maintenance item.
	// TODO images.
}

// CreatePlant inserts a new plant in the database and returns the plant back with its ID.
func (s *Store) CreatePlant(ctx context.Context, p *Plant) (*Plant, error) {
	prep, err := s.db.Prepare("INSERT INTO plants (name, lighting, watering, fertilizing, toxicity, notes) VALUES (?,?,?,?,?,?) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	err = prep.QueryRowContext(ctx, p.Name, p.Lighting, p.Watering, p.Fertilizing, p.Toxicity, p.Notes).Scan(&p.ID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetPlantByID returns the plant with specified identifier.
func (s *Store) GetPlantByID(ctx context.Context, id int32) (*Plant, error) {
	plant := &Plant{}

	prep, err := s.db.Prepare("SELECT id, name, lighting, watering, fertilizing, toxicity, notes FROM plants WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	err = prep.QueryRowContext(ctx, id).Scan(&plant.ID, &plant.Name, &plant.Lighting, &plant.Watering, &plant.Fertilizing, &plant.Toxicity, &plant.Notes)
	if err != nil {
		return nil, err
	}

	return plant, nil
}
