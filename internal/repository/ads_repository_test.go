package repository

import (
	"context"
	"fmt"
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

	ad := createModel()

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

func TestPostgresAdsRepository_Update(t *testing.T) {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://avito_user:postgres@localhost:5432/avito?sslmode=disable")

	if err != nil {
		t.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer pool.Close()

	repo := NewPosgresAdsRepository(pool)

	ad := createModel()
	err = repo.Create(ctx, ad)
	if err != nil {
		t.Fatalf("Failed to create test ad: %v", err)
	}

	updatedAd := *ad

	updatedAd.Title = "Updated Title"
	updatedAd.Description = "Updated Description"
	updatedAd.Price = 150.50
	updatedAd.Category = "Updated Category"
	updatedAd.UpdatedAt = time.Now()

	err = repo.Update(ctx, &updatedAd)
	if err != nil {
		t.Fatalf("Failed to update ad: %v", err)
	}

	retrievedUpdatedAd, err := repo.GetById(ctx, updatedAd.ID)

	if retrievedUpdatedAd.Title != "Updated Title" {
		t.Error("Expected title")
	}

	if retrievedUpdatedAd.Description != "Updated Description" {
		t.Error("Expected description")
	}

	if updatedAd.Price != 150.50 {
		t.Error("Expected price")
	}

}

func TestPostgresAdsRepository_List_Cursor(t *testing.T) {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://avito_user:postgres@localhost:5432/avito?sslmode=disable")
	if err != nil {
		t.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	repo := NewPosgresAdsRepository(pool)

	ids := []string{}
	for i := 0; i < 5; i++ {
		ad := createModel()
		ad.Title = fmt.Sprintf("Test Ad %d", i)
		err := repo.Create(ctx, ad)
		if err != nil {
			t.Fatalf("Failed to create test ad: %v", err)
		}
		ids = append(ids, ad.ID)
		defer repo.Delete(ctx, ad.ID)
	}

	pageFirst, err := repo.List(ctx, "", "", 2)
	if err != nil {
		t.Fatalf("Failed to list page 1: %v", err)
	}

	if len(pageFirst) != 2 {
		t.Errorf("Expected 2 ads on page 1, got %d", len(pageFirst))
	}

	cursor := pageFirst[len(pageFirst)-1].ID

	pageSecond, err := repo.List(ctx, "", cursor, 2)
	if err != nil {
		t.Errorf("Expected 2 ads on page 2, got %d", len(pageSecond))
	}

	for _, adFirst := range pageFirst {
		for _, adSecond := range pageSecond {
			if adFirst.ID == adSecond.ID {
				t.Errorf("Found duplicate ad")
			}
		}
	}

}

func createModel() *model.Ad {
	return &model.Ad{
		ID:           uuid.New().String(),
		Title:        "Test Ad",
		Description:  "Test Description",
		Price:        100.50,
		Category:     "test Category",
		ContactPhone: "+79132554689",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
