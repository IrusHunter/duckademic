package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/jmoiron/sqlx"
)

// AcademicRankRepository represents a storage for academic rank entities.
type AcademicRankRepository interface {
	// Add inserts a new AcademicRank into the database and returns it, or an error if it fails.
	Add(context.Context, entities.AcademicRank) (entities.AcademicRank, error)
	Clear(context.Context) // Clear removes all academic ranks from the repository.
	FindBySlug(context.Context, string) *entities.AcademicRank
	// GetAll returns a slice with all academic ranks.
	GetAll(context.Context) []entities.AcademicRank
}

// NewAcademicRankRepository creates a new AcademicRankRepository instance.
//
// It requires a database connection (db).
func NewAcademicRankRepository(db *sqlx.DB) AcademicRankRepository {
	return &academicRankRepository{db: db}
}

type academicRankRepository struct {
	db *sqlx.DB
}

func (r *academicRankRepository) Add(
	ctx context.Context, academicRank entities.AcademicRank,
) (entities.AcademicRank, error) {
	rows, err := r.db.NamedQueryContext(
		ctx,
		` INSERT INTO academic_ranks
		(id, slug, title)
		VALUES
		(:id, :slug, :title)
		RETURNING created_at, updated_at
		`,
		academicRank,
	)

	if err != nil {
		return entities.AcademicRank{}, fmt.Errorf("failed to insert %s: %w", academicRank.String(), err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&academicRank.CreatedAt, &academicRank.UpdatedAt); err != nil {
			return entities.AcademicRank{},
				fmt.Errorf("failed to scan database row for %s: %w", academicRank.String(), err)
		}
	}

	return academicRank, nil
}
func (r *academicRankRepository) Clear(ctx context.Context) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM academic_ranks`)
	if err != nil {
		log.Println("Can't truncate table academic_ranks: " + err.Error())
	}
}
func (r *academicRankRepository) FindBySlug(ctx context.Context, slug string) *entities.AcademicRank {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, slug, title, created_at, updated_at FROM academic_ranks
		WHERE slug=$1`,
		slug,
	)
	var academicRank entities.AcademicRank
	if err := row.Scan(
		&academicRank.ID,
		&academicRank.Slug,
		&academicRank.Title,
		&academicRank.CreatedAt,
		&academicRank.UpdatedAt,
	); err != nil {
		log.Printf("Can't scan database row for slug %q: %s \n", slug, err.Error())
		return nil
	}

	return &academicRank
}
func (r *academicRankRepository) GetAll(ctx context.Context) []entities.AcademicRank {
	ranks := []entities.AcademicRank{}
	err := r.db.SelectContext(
		ctx,
		&ranks,
		`SELECT id, slug, title, created_at, updated_at FROM academic_ranks`,
	)
	if err != nil {
		log.Println("failed to get academic ranks: " + err.Error())
	}

	return ranks
}
