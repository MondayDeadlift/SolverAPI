package postgres

import (
	"SolverAPI/internal/model"
	"SolverAPI/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type KataRepo struct {
	db *sql.DB
}

func NewKataRepository(db *sql.DB) repository.KataRepository {
	return &KataRepo{db: db}
}

func (r *KataRepo) SaveKata(ctx context.Context, kata *model.Kata) error {
	tagsJSON, _ := json.Marshal(kata.Tags)
	languagesJSON, _ := json.Marshal(kata.Languages)

	query := `
        INSERT INTO katas (id, name, slug, url, tags, languages, added_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            slug = EXCLUDED.slug,
            url = EXCLUDED.url,
            tags = EXCLUDED.tags,
            languages = EXCLUDED.languages
    `
	_, err := r.db.ExecContext(ctx, query,
		kata.ID,
		kata.Name,
		kata.Slug,
		kata.URL,
		tagsJSON,
		languagesJSON,
		kata.AddedAt,
	)
	return err
}

func (r *KataRepo) GetRandomKata(ctx context.Context) (*model.Kata, error) {
	query := `
        SELECT id, name, slug, url, tags, languages, added_at 
        FROM katas
        ORDER BY RANDOM()
        LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, query)

	var kata model.Kata
	var tagsJSON, languagesJSON []byte

	err := row.Scan(
		&kata.ID,
		&kata.Name,
		&kata.Slug,
		&kata.URL,
		&tagsJSON,
		&languagesJSON,
		&kata.AddedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan kata: %w", err)
	}

	json.Unmarshal(tagsJSON, &kata.Tags)
	json.Unmarshal(languagesJSON, &kata.Languages)

	return &kata, nil
}
