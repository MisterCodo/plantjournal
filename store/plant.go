package store

import (
	"context"
	"fmt"
)

// Plant contains the details of a plant.
type Plant struct {
	ID          int
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
func (s *Store) GetPlantByID(ctx context.Context, id int) (*Plant, error) {
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

// GetPlants returns all plants.
func (s *Store) GetPlants(ctx context.Context) ([]*Plant, error) {
	plants := []*Plant{}

	prep, err := s.db.Prepare("SELECT id, name, lighting, watering, fertilizing, toxicity, notes FROM plants ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	rows, err := prep.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		plant := &Plant{}
		err := rows.Scan(&plant.ID, &plant.Name, &plant.Lighting, &plant.Watering, &plant.Fertilizing, &plant.Toxicity, &plant.Notes)
		if err != nil {
			return nil, err
		}
		plants = append(plants, plant)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return plants, nil
}

// UdatePlant updates the passed plant in the database.
func (s *Store) UpdatePlant(ctx context.Context, p *Plant) error {
	prep, err := s.db.Prepare("UPDATE plants SET name=?, lighting=?, watering=?, fertilizing=?, toxicity=?, notes=? WHERE id=?")
	if err != nil {
		return err
	}
	defer prep.Close()

	res, err := prep.ExecContext(ctx, p.Name, p.Lighting, p.Watering, p.Fertilizing, p.Toxicity, p.Notes, p.ID)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return fmt.Errorf("failed to update plant")
	}

	return nil
}
