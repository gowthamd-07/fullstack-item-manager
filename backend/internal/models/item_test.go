package models

import (
	"strings"
	"testing"
)

func TestCreateItemDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dto     CreateItemDTO
		wantErr string
	}{
		{
			name:    "valid item",
			dto:     CreateItemDTO{Name: "Widget", Price: 9.99},
			wantErr: "",
		},
		{
			name:    "empty name",
			dto:     CreateItemDTO{Name: "", Price: 9.99},
			wantErr: "name is required",
		},
		{
			name:    "whitespace-only name",
			dto:     CreateItemDTO{Name: "   ", Price: 9.99},
			wantErr: "name is required",
		},
		{
			name:    "name too long",
			dto:     CreateItemDTO{Name: strings.Repeat("a", 256), Price: 9.99},
			wantErr: "name must be 255 characters or less",
		},
		{
			name:    "negative price",
			dto:     CreateItemDTO{Name: "Widget", Price: -1.0},
			wantErr: "price must be non-negative",
		},
		{
			name:    "zero price is valid",
			dto:     CreateItemDTO{Name: "Free Item", Price: 0},
			wantErr: "",
		},
		{
			name:    "price too large",
			dto:     CreateItemDTO{Name: "Widget", Price: 100000000},
			wantErr: "price exceeds maximum allowed value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}

func TestUpdateItemDTO_Validate(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	floatPtr := func(f float64) *float64 { return &f }

	tests := []struct {
		name    string
		dto     UpdateItemDTO
		wantErr string
	}{
		{
			name:    "valid update name only",
			dto:     UpdateItemDTO{Name: strPtr("New Name")},
			wantErr: "",
		},
		{
			name:    "valid update price only",
			dto:     UpdateItemDTO{Price: floatPtr(19.99)},
			wantErr: "",
		},
		{
			name:    "empty name",
			dto:     UpdateItemDTO{Name: strPtr("")},
			wantErr: "name cannot be empty",
		},
		{
			name:    "whitespace name",
			dto:     UpdateItemDTO{Name: strPtr("   ")},
			wantErr: "name cannot be empty",
		},
		{
			name:    "name too long",
			dto:     UpdateItemDTO{Name: strPtr(strings.Repeat("a", 256))},
			wantErr: "name must be 255 characters or less",
		},
		{
			name:    "negative price",
			dto:     UpdateItemDTO{Price: floatPtr(-5.0)},
			wantErr: "price must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}
