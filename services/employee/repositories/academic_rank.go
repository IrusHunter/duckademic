package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AcademicRankRepository represents a storage for academic rank entities.
type AcademicRankRepository interface {
	// Add inserts a new AcademicRank into the database and returns it, or an error if it fails.
	Add(context.Context, entities.AcademicRank) (entities.AcademicRank, error)
	Clear(context.Context) // Clear removes all academic ranks from the database.
	// FindBySlug returns a pointer to the academic rank from database with the given slug.
	FindBySlug(context.Context, string) *entities.AcademicRank
	// FindByID returns a pointer to the academic rank from database with the given id.
	FindByID(context.Context, uuid.UUID) *entities.AcademicRank
	// GetAll returns a slice with all academic ranks from database.
	GetAll(context.Context) []entities.AcademicRank
	// Delete removes the AcademicRank with the specified ID from the database.
	Delete(context.Context, uuid.UUID) error
	// Update updates the AcademicRank with the specified ID and returns the updated entity.
	Update(context.Context, uuid.UUID, entities.AcademicRank) (entities.AcademicRank, error)
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
		RETURNING created_at, updated_at`,
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
func (r *academicRankRepository) FindByID(ctx context.Context, id uuid.UUID) *entities.AcademicRank {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, slug, title, created_at, updated_at FROM academic_ranks
		WHERE id=$1`,
		id,
	)
	var academicRank entities.AcademicRank
	if err := row.Scan(
		&academicRank.ID,
		&academicRank.Slug,
		&academicRank.Title,
		&academicRank.CreatedAt,
		&academicRank.UpdatedAt,
	); err != nil {
		log.Printf("Can't scan database row for id %q: %s \n", id, err.Error())
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
func (r *academicRankRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM academic_ranks WHERE id=$1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete academic rank %q: %w", id, err)
	}

	return nil
}
func (r *academicRankRepository) Update(ctx context.Context, id uuid.UUID, rank entities.AcademicRank,
) (entities.AcademicRank, error) {
	rank.ID = id
	rows, err := r.db.NamedQueryContext(
		ctx,
		`UPDATE academic_ranks SET
		title= :title
		WHERE id= :id
		RETURNING slug, created_at, updated_at`,
		rank,
	)
	if err != nil {
		return entities.AcademicRank{}, fmt.Errorf("failed to update %s: %w", rank.String(), err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&rank.Slug, &rank.CreatedAt, &rank.UpdatedAt); err != nil {
			return entities.AcademicRank{}, fmt.Errorf("failed to scan database row for %s: %w", rank.String(), err)
		}
	}

	return rank, nil
}
