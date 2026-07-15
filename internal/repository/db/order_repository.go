package db

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"context"
	"database/sql"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		`INSERT INTO orders (order_code, user_id, total, status) VALUES (?, ?, ?, ?)`,
		order.OrderCode, order.UserID, order.Total, order.Status,
	)
	if err != nil {
		return err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = int(orderID)

	for i, item := range order.Items {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO order_items (order_id, service_id, service_name, price) VALUES (?, ?, ?, ?)`,
			order.ID, item.ServiceID, item.ServiceName, item.Price,
		)
		if err != nil {
			return err
		}
		order.Items[i].OrderID = order.ID
	}

	return tx.Commit()
}

func (r *orderRepository) FindByUser(ctx context.Context, userID int) ([]entity.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, order_code, user_id, total, status, created_at FROM orders WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.OrderCode, &o.UserID, &o.Total, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range orders {
		items, err := r.findItems(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (r *orderRepository) findItems(ctx context.Context, orderID int) ([]entity.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, order_id, service_id, service_name, price FROM order_items WHERE order_id = ?`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.OrderItem
	for rows.Next() {
		var it entity.OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ServiceID, &it.ServiceName, &it.Price); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	return items, rows.Err()
}
