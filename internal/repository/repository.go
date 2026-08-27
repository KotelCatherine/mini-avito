package repository

import (
	"context"

	"github.com/KotelCatherine/mini-avito/internal/model"
)

type AdsRepository interface {
	Create(ctx context.Context, ad *model.Ad) error
	GetByID(ctx context.Context, id string) (*model.Ad, error)
	List(ctx context.Context, category string, cursor string, limit int) ([]model.Ad, error)
	Update(ctx context.Context, ad *model.Ad) error
	Delete(ctx context.Context, id string) error
}
