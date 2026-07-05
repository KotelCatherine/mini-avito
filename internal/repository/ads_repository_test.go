package repository

import (
	"context"
	"testing"
	"time"

	"github.com/KotelCatherine/mini-avito/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresAdsRepository_Create(t *testing.T) {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://avito_user:postgres@localhost:5432/avito?sslmode=disable")

	if err != nil {
		t.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer pool.Close()

	repo := NewPosgresAdsRepository(pool)

	ad := &model.Ad{
		ID:           uuid.New().String(),
		Title:        "Test Ad",
		Description:  "Test Description",
		Price:        100.50,
		Category:     "test Category",
		ContactPhone: "+79132554689",
		CreateAt:     time.Now(),
		UpdateAt:     time.Now(),
	}

	err = repo.Create(ctx, ad)
	if err != nil {
		t.Fatalf("Failed to create ad: %v", err)
	}

	retrieved, err := repo.GetById(ctx, ad.ID)
	if err != nil {
		t.Fatalf("Failed to get ad: %v", err)
	}

	if retrieved.Title != ad.Title {
		t.Errorf("Expected title %s, got %s", ad.Title, retrieved.Title)
	}

}
