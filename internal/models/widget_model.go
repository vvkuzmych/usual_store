package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Widget represents a widget in the system.
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// GetWidget retrieves a widget by its ID.
func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at FROM widgets WHERE id=$1`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var widget Widget
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return widget, fmt.Errorf("widget not found")
		}
		return widget, err
	}
	return widget, nil
}
