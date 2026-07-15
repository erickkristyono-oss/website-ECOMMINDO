package db

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"context"
	"database/sql"
	"errors"
)

type serviceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) domain.ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) FindAll(ctx context.Context) ([]entity.Service, error) {
	query := `
		SELECT id, name, short_desc, description, price, icon
		FROM services
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []entity.Service
	for rows.Next() {
		var s entity.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.ShortDesc, &s.Description, &s.Price, &s.Icon); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, rows.Err()
}

func (r *serviceRepository) FindByID(ctx context.Context, id string) (*entity.Service, error) {
	query := `
		SELECT id, name, short_desc, description, price, icon
		FROM services
		WHERE id = ?
	`

	var s entity.Service
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.ShortDesc, &s.Description, &s.Price, &s.Icon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}

	return &s, nil
}
