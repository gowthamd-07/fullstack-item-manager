package models

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateItemDTO struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type UpdateItemDTO struct {
	Name  *string  `json:"name,omitempty"`
	Price *float64 `json:"price,omitempty"`
}
