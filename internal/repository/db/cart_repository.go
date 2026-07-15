package db

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"context"
	"database/sql"
)

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) domain.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Add(ctx context.Context, userID int, serviceID string) error {
	query := `
		INSERT IGNORE INTO cart_items (user_id, service_id)
		VALUES (?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, userID, serviceID)
	return err
}

func (r *cartRepository) Remove(ctx context.Context, userID int, serviceID string) error {
	query := `DELETE FROM cart_items WHERE user_id = ? AND service_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, serviceID)
	return err
}

func (r *cartRepository) FindByUser(ctx context.Context, userID int) ([]entity.CartItem, error) {
	query := `SELECT id, user_id, service_id FROM cart_items WHERE user_id = ? ORDER BY created_at`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.CartItem
	for rows.Next() {
		var item entity.CartItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.ServiceID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *cartRepository) Clear(ctx context.Context, userID int) error {
	query := `DELETE FROM cart_items WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
