package store

import (
	"context"
	"time"
)

// Action contains a plant maintenance action details.
type Action struct {
	Day        string // YYYY-MM-DD
	PlantID    int
	Watered    bool
	Fertilized bool
	Notes      string
}

// CreateAction inserts a new action in the database and returns the plant back with its ID.
func (s *Store) CreateAction(ctx context.Context, a *Action) (*Action, error) {
	// TODO
	return nil, nil
}

// GetActionsByPlantID returns all actions for a plant with passed id.
func (s *Store) GetActionsByPlantID(ctx context.Context, plantID int) ([]*Action, error) {
	actions := []*Action{}

	prep, err := s.db.Prepare("SELECT day, plant_id, watered, fertilized, notes FROM actions WHERE plant_id = ? ORDER BY day desc")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	rows, err := prep.QueryContext(ctx, plantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		action := &Action{}
		err := rows.Scan(&action.Day, &action.PlantID, &action.Watered, &action.Fertilized, &action.Notes)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// DELETE THIS
	actions = append(actions, &Action{Day: "2024-02-24", PlantID: plantID, Watered: true, Fertilized: false, Notes: "debugging a"})
	actions = append(actions, &Action{Day: "2024-02-23", PlantID: plantID, Watered: true, Fertilized: true, Notes: "debugging b"})
	actions = append(actions, &Action{Day: "2024-02-21", PlantID: plantID, Watered: false, Fertilized: false, Notes: "debugging c"})

	return actions, nil
}

// DeleteAction deletes the action with passed day.
func (s *Store) DeleteAction(ctx context.Context, day string) error {
	// TODO
	return nil
}

// WaterPlant upserts an action for the day, indicating plant was watered.
func (s *Store) WaterPlant(ctx context.Context, plantID int) (*Action, error) {
	prep, err := s.db.Prepare("INSERT INTO actions (day, plant_id, watered, fertilized, notes) VALUES (?,?,?,?,?) ON CONFLICT(day, plant_id) DO UPDATE SET watered = ? RETURNING day, plant_id, watered, fertilized, notes")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	action := &Action{}
	day := time.Now().Format("2006-01-02")
	err = prep.QueryRowContext(ctx, day, plantID, 1, 0, "", 1).Scan(&action.Day, &action.PlantID, &action.Watered, &action.Fertilized, &action.Notes)
	if err != nil {
		return nil, err
	}

	return action, nil
}

// FertilizePlant upserts an action for the day, indicating plant was fertilized.
func (s *Store) FertilizePlant(ctx context.Context, plantID int) (*Action, error) {
	prep, err := s.db.Prepare("INSERT INTO actions (day, plant_id, watered, fertilized, notes) VALUES (?,?,?,?,?) ON CONFLICT(day, plant_id) DO UPDATE SET fertilized = ? RETURNING day, plant_id, watered, fertilized, notes")
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	action := &Action{}
	day := time.Now().Format("2006-01-02")
	err = prep.QueryRowContext(ctx, day, plantID, 0, 1, "", 1).Scan(&action.Day, &action.PlantID, &action.Watered, &action.Fertilized, &action.Notes)
	if err != nil {
		return nil, err
	}

	return action, nil
}
