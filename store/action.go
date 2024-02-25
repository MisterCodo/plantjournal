package store

import "context"

// Action contains a plant maintenance action details.
type Action struct {
	ID         int
	PlantID    int
	Day        string // YYYY-MM-DD
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

	prep, err := s.db.Prepare("SELECT id, plant_id, day, watered, fertilized, notes FROM actions WHERE plant_id = ? ORDER BY day desc")
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
		err := rows.Scan(&action.ID, &action.PlantID, &action.Day, &action.Watered, &action.Fertilized, &action.Notes)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return actions, nil
}

// DeleteAction deletes the action with passed id.
func (s *Store) DeleteAction(ctx context.Context, id int) error {
	// TODO
	return nil
}
