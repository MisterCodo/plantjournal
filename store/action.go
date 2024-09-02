package store

import (
	"context"
	"fmt"
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

	return actions, nil
}

// DeleteAction deletes the action with passed pland id and action day.
func (s *Store) DeleteAction(ctx context.Context, plantID int, day string) error {
	prep, err := s.db.Prepare("DELETE FROM actions WHERE plant_id=? and day=?")
	if err != nil {
		return err
	}
	defer prep.Close()

	// Validate day value is in expected format.
	_, err = time.Parse("2006-01-02", day)
	if err != nil {
		return err
	}

	res, err := prep.ExecContext(ctx, plantID, day)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return fmt.Errorf("failed to delete plant action")
	}

	return nil
}

// WaterPlant upserts an action for the day, indicating plant was watered.
func (s *Store) WaterPlant(ctx context.Context, plantID int) error {
	prep, err := s.db.Prepare("INSERT INTO actions (day, plant_id, watered, fertilized, notes) VALUES (?,?,?,?,?) ON CONFLICT(day, plant_id) DO UPDATE SET watered = ?")
	if err != nil {
		return err
	}
	defer prep.Close()

	day := time.Now().Format("2006-01-02")
	res, err := prep.ExecContext(ctx, day, plantID, 1, 0, "", 1)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return fmt.Errorf("failed to update watered value for action")
	}

	return nil
}

// FertilizePlant upserts an action for the day, indicating plant was fertilized.
func (s *Store) FertilizePlant(ctx context.Context, plantID int) error {
	prep, err := s.db.Prepare("INSERT INTO actions (day, plant_id, watered, fertilized, notes) VALUES (?,?,?,?,?) ON CONFLICT(day, plant_id) DO UPDATE SET fertilized = ?")
	if err != nil {
		return err
	}
	defer prep.Close()

	day := time.Now().Format("2006-01-02")
	res, err := prep.ExecContext(ctx, day, plantID, 0, 1, "", 1)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return fmt.Errorf("failed to update fertilized value for action")
	}

	return nil
}

// AddNoteToPlant inserts an action for the day if it does not already exists. Watered and fertilized default to 0 while notes defaults to empty string.
func (s *Store) AddNoteToPlant(ctx context.Context, plantID int) error {
	prep, err := s.db.Prepare("INSERT INTO actions (day, plant_id, watered, fertilized, notes) VALUES (?,?,?,?,?) ON CONFLICT(day, plant_id) DO NOTHING")
	if err != nil {
		return err
	}
	defer prep.Close()

	day := time.Now().Format("2006-01-02")
	_, err = prep.ExecContext(ctx, day, plantID, 0, 0, "")
	if err != nil {
		return err
	}

	// OK if affected rows is 0, so don't check

	return nil
}

// UdateAction updates the passed action in the database.
func (s *Store) UpdateAction(ctx context.Context, a *Action) error {
	// TODO: for now only updates notes, fix for watered and fertilized.
	prep, err := s.db.Prepare("UPDATE actions SET notes=? WHERE day=? and plant_id=?")
	if err != nil {
		return err
	}
	defer prep.Close()

	res, err := prep.ExecContext(ctx, a.Notes, a.Day, a.PlantID)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return fmt.Errorf("failed to update action")
	}

	return nil
}
