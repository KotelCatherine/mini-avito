package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KotelCatherine/mini-avito/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdsRepository struct {
	pool *pgxpool.Pool
}

func NewPosgresAdsRepository(pool *pgxpool.Pool) *PostgresAdsRepository {
	return &PostgresAdsRepository{pool: pool}
}

func (r *PostgresAdsRepository) Create(ctx context.Context, ad *model.Ad) error {

	query := `INSERT INTO ads (id, title, description, price, category, contact_phone, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.pool.Exec(ctx, query,
		ad.ID,
		ad.Title,
		ad.Description,
		ad.Price,
		ad.Category,
		ad.ContactPhone,
		ad.CreateAt,
		ad.UpdateAt,
	)

	if err != nil {
		return fmt.Errorf("exec insert: %w", err)
	}

	return nil

}

func (r *PostgresAdsRepository) GetById(ctx context.Context, id string) (*model.Ad, error) {

	query := `SELECT id, title, description, price, category, contact_phone, created_at, updated_at FROM ads WHERE id = $1`

	var ad model.Ad
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&ad.ID,
		&ad.Title,
		&ad.Description,
		&ad.Price,
		&ad.Category,
		&ad.ContactPhone,
		&ad.CreateAt,
		&ad.UpdateAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ad not found: %w", err)
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &ad, nil

}
