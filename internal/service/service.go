package service

import (
	"context"

	"github.com/KotelCatherine/mini-avito/internal/model"
)

type CreateAdRequest struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Category     string  `json:"category"`
	ContactPhone string  `json:"contact_phone"`
}

type UpdateAdRequest struct {
	Title        *string  `json:"title,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	Category     *string  `json:"category,omitempty"`
	ContactPhone *string  `json:"contact_phone,omitempty"`
}

type AdsService interface {
	CreateAd(ctx context.Context, req CreateAdRequest) (*model.Ad, error)
	GetAd(ctx context.Context, id string) (*model.Ad, error)
	ListAds(ctx context.Context, category string, cursor string, limit int) ([]model.Ad, string, error)
	UpdateAd(ctx context.Context, id string, req UpdateAdRequest) (*model.Ad, error)
	DeleteAd(ctx context.Context, id string) error
}
