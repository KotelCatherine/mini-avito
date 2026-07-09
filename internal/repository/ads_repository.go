package repository

import (
	"context"
	"errors"
	"fmt"

	domainErrors "github.com/KotelCatherine/mini-avito/internal/domain/errors"
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

func (repo *PostgresAdsRepository) Create(ctx context.Context, ad *model.Ad) error {

	query := `INSERT INTO ads (id, title, description, price, category, contact_phone, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := repo.pool.Exec(ctx, query,
		ad.ID,
		ad.Title,
		ad.Description,
		ad.Price,
		ad.Category,
		ad.ContactPhone,
		ad.CreatedAt,
		ad.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("exec insert: %w", err)
	}

	return nil

}

func (repo *PostgresAdsRepository) GetById(ctx context.Context, id string) (*model.Ad, error) {

	query := `SELECT id, title, description, price, category, contact_phone, created_at, updated_at FROM ads WHERE id = $1`

	var ad model.Ad
	err := repo.pool.QueryRow(ctx, query, id).Scan(
		&ad.ID,
		&ad.Title,
		&ad.Description,
		&ad.Price,
		&ad.Category,
		&ad.ContactPhone,
		&ad.CreatedAt,
		&ad.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.NewNotFoundError("ad", id)
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &ad, nil

}

func (repo *PostgresAdsRepository) Update(ctx context.Context, ad *model.Ad) error {

	query := `UPDATE ads SET title = $1, description = $2, price = $3, category = $4, contact_phone = $5, updated_at = $6
	WHERE id = $7`

	result, err := repo.pool.Exec(ctx, query,
		ad.Title,
		ad.Description,
		ad.Price,
		ad.Category,
		ad.ContactPhone,
		ad.UpdatedAt,
		ad.ID,
	)

	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domainErrors.ErrAdNotFound
	}

	return nil

}

func (repo *PostgresAdsRepository) Delete(ctx context.Context, id string) error {

	query := `DELETE FROM ads WHERE id = $1`

	result, err := repo.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("exec delete: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domainErrors.NewNotFoundError("ad", id)
	}

	return nil

}

func (repo *PostgresAdsRepository) List(ctx context.Context, category string, cursor string, limit int) ([]model.Ad, error) {

	query := `SELECT id, title, description, price, category, contact_phone, created_at, updated_at FROM ads`

	args := []interface{}{}
	conditions := []string{}
	argIndex := 1

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, category)
		argIndex++
	}

	if cursor != "" {
		conditions = append(conditions, fmt.Sprintf("id < $%d", argIndex))
		args = append(args, cursor)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + joinConditions(conditions)
	}

	query += " ORDER BY id DESC"
	query += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, limit)

	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var ads []model.Ad
	for rows.Next() {
		var ad model.Ad
		err := rows.Scan(
			&ad.ID,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Category,
			&ad.ContactPhone,
			&ad.CreatedAt,
			&ad.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		ads = append(ads, ad)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iterate: %w", err)
	}

	return ads, nil

}

func joinConditions(conditions []string) string {

	result := ""
	for i, condition := range conditions {
		if i > 0 {
			result += " AND "
		}
		result += condition
	}

	return result

}
