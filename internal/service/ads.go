package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	domainErrors "github.com/KotelCatherine/mini-avito/internal/domain/errors"

	"github.com/KotelCatherine/mini-avito/internal/model"
	"github.com/KotelCatherine/mini-avito/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdsServiceImpl struct {
	pool *pgxpool.Pool
	repo repository.AdsRepository
}

func NewAdsService(pool *pgxpool.Pool, repo repository.AdsRepository) AdsService {
	return &AdsServiceImpl{
		pool: pool,
		repo: repo,
	}
}

func (adService *AdsServiceImpl) CreateAd(ctx context.Context, req CreateAdRequest) (*model.Ad, error) {

	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	tx, err := adService.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	txRepo := repository.NewPosgresAdsRepository(tx)

	ad := &model.Ad{
		ID:           uuid.New().String(),
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		Category:     req.Category,
		ContactPhone: req.ContactPhone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = txRepo.Create(ctx, ad)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return ad, nil

}

func (adService *AdsServiceImpl) GetAd(ctx context.Context, id string) (*model.Ad, error) {

	if id == "" {
		return nil, domainErrors.ErrInvalidInput
	}

	ad, err := adService.repo.GetByID(ctx, id)
	if err != nil {

		if errors.Is(err, domainErrors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("get ad by id %s: %w", id, err)

	}

	return ad, nil

}

func (adService *AdsServiceImpl) ListAds(ctx context.Context, category string, cursor string, limit int) ([]model.Ad, string, error) {

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	ads, err := adService.repo.List(ctx, category, cursor, limit)
	if err != nil {
		return nil, "", fmt.Errorf("list ads: %w", err)
	}

	var nextCursor string

	if len(ads) == limit {
		nextCursor = ads[len(ads)-1].ID
	}

	return ads, nextCursor, nil

}

func (adService *AdsServiceImpl) UpdateAd(ctx context.Context, id string, req UpdateAdRequest) (*model.Ad, error) {

	if id == "" {
		return nil, domainErrors.ErrInvalidInput
	}

	tx, err := adService.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	txRepo := repository.NewPosgresAdsRepository(tx)

	ad, err := txRepo.GetById(ctx, id)
	if err != nil {

		if errors.Is(err, domainErrors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("get ad for update: %w", err)

	}

	if req.Title != nil {
		if len(*req.Title) < 3 || len(*req.Title) > 100 {
			return nil, domainErrors.ErrInvalidTitle
		}
		ad.Title = *req.Title
	}

	if req.Description != nil {
		ad.Description = *req.Description
	}

	if req.Price != nil {
		if *req.Price < 0 {
			return nil, domainErrors.ErrAdInvalidPrice
		}
		ad.Price = *req.Price
	}

	if req.Category != nil {
		if *req.Category == "" {
			return nil, domainErrors.ErrInvalidCategory
		}
		ad.Category = *req.Category
	}

	if req.ContactPhone != nil {
		ad.ContactPhone = *req.ContactPhone
	}

	ad.UpdatedAt = time.Now()

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return ad, nil

}

func (adService *AdsServiceImpl) DeleteAd(ctx context.Context, id string) error {

	if id == "" {
		return domainErrors.ErrInvalidInput
	}

	err := adService.repo.Delete(ctx, id)
	if err != nil {

		if errors.Is(err, domainErrors.ErrNotFound) {
			return err
		}

		return fmt.Errorf("delete ad: %w", err)

	}

	return nil

}

func validateCreateRequest(req CreateAdRequest) error {

	if len(req.Title) < 3 || len(req.Title) > 100 {
		return domainErrors.ErrInvalidTitle
	}

	if req.Price < 0 {
		return domainErrors.ErrAdInvalidPrice
	}

	if req.Category == "" {
		return domainErrors.ErrInvalidCategory
	}

	return nil

}
