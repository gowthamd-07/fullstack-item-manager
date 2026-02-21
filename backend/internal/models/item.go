package models

import (
	"fmt"
	"strings"
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

func (d *CreateItemDTO) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(d.Name) > 255 {
		return fmt.Errorf("name must be 255 characters or less")
	}
	if d.Price < 0 {
		return fmt.Errorf("price must be non-negative")
	}
	if d.Price > 99999999.99 {
		return fmt.Errorf("price exceeds maximum allowed value")
	}
	return nil
}

type UpdateItemDTO struct {
	Name  *string  `json:"name,omitempty"`
	Price *float64 `json:"price,omitempty"`
}

func (d *UpdateItemDTO) Validate() error {
	if d.Name != nil {
		trimmed := strings.TrimSpace(*d.Name)
		d.Name = &trimmed
		if trimmed == "" {
			return fmt.Errorf("name cannot be empty")
		}
		if len(trimmed) > 255 {
			return fmt.Errorf("name must be 255 characters or less")
		}
	}
	if d.Price != nil {
		if *d.Price < 0 {
			return fmt.Errorf("price must be non-negative")
		}
		if *d.Price > 99999999.99 {
			return fmt.Errorf("price exceeds maximum allowed value")
		}
	}
	return nil
}

type PaginatedResponse struct {
	Items  []Item `json:"items"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
