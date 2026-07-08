package model

import "time"

type Ad struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	Category     string    `json:"category"`
	ContactPhone string    `json:"contact_phone"`
	CreatedAt    time.Time `json:"create_at"`
	UpdatedAt    time.Time `json:"update_at"`
}
