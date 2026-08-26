package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/adaravaks/URLshortener/internal/model"
)

var ErrNotFound = errors.New("link not found")

type LinkRepository interface {
	Create(ctx context.Context, link *model.Link) error
	GetByCode(ctx context.Context, code string) (*model.Link, error)
	IncrementClickCount(ctx context.Context, code string) error
}

type postgresLinkRepository struct {
	db *sql.DB
}

func NewPostgresLinkRepository(db *sql.DB) LinkRepository {
	return &postgresLinkRepository{db: db}
}

func (r *postgresLinkRepository) Create(ctx context.Context, link *model.Link) error {
	query := `
		INSERT INTO links (short_code, original_url)
		VALUES ($1, $2)
		RETURNING id, created_at, click_count
	`
	return r.db.QueryRowContext(ctx, query, link.ShortCode, link.OriginalURL).
		Scan(&link.ID, &link.CreatedAt, &link.ClickCount)
}

func (r *postgresLinkRepository) GetByCode(ctx context.Context, code string) (*model.Link, error) {
	query := `
	SELECT id, short_code, original_url, created_at, click_count
	FROM links
	WHERE short_code = $1
	`
	link := &model.Link{}
	err := r.db.QueryRowContext(ctx, query, code).
		Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt, &link.ClickCount)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *postgresLinkRepository) IncrementClickCount(ctx context.Context, code string) error {
	query := `
	UPDATE links SET click_count = click_count + 1 WHERE short_code = $1
	`
	result, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
