package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gowthamd/go-crud-app/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(ctx context.Context, item *models.CreateItemDTO) (*models.Item, error) {
	query := `
		INSERT INTO items (name, price)
		VALUES ($1, $2)
		RETURNING id, name, price, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query, item.Name, item.Price)

	var i models.Item
	err := row.Scan(&i.ID, &i.Name, &i.Price, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return &i, nil
}

func (r *ItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	query := `
		SELECT id, name, price, created_at, updated_at
		FROM items
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var i models.Item
	err := row.Scan(&i.ID, &i.Name, &i.Price, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &i, nil
}

func (r *ItemRepository) GetAll(ctx context.Context) ([]models.Item, error) {
	query := `
		SELECT id, name, price, created_at, updated_at
		FROM items
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var i models.Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Price, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, i)
	}

	return items, nil
}

func (r *ItemRepository) Update(ctx context.Context, id uuid.UUID, updates *models.UpdateItemDTO) (*models.Item, error) {
	// Build dynamic update query
	query := "UPDATE items SET updated_at = NOW()"
	args := []interface{}{}
	argId := 1

	if updates.Name != nil {
		query += fmt.Sprintf(", name = $%d", argId)
		args = append(args, *updates.Name)
		argId++
	}
	if updates.Price != nil {
		query += fmt.Sprintf(", price = $%d", argId)
		args = append(args, *updates.Price)
		argId++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, price, created_at, updated_at", argId)
	args = append(args, id)

	row := r.db.QueryRow(ctx, query, args...)

	var i models.Item
	err := row.Scan(&i.ID, &i.Name, &i.Price, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return &i, nil
}

func (r *ItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM items WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}
