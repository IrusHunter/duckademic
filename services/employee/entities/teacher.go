package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Teacher struct {
	EmployeeID       uuid.UUID  `db:"employee_id" json:"employee_id"`
	Email            string     `db:"email" json:"email"`
	AcademicDegreeID *uuid.UUID `db:"academic_degree_id" json:"academic_degree_id,omitempty"`
	AcademicRankID   uuid.UUID  `db:"academic_rank_id" json:"academic_rank_id"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`           // Record creation timestamp.
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`           // Record last update timestamp.
	DeletedAt        *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // Record deleted timestamp.

	Employee       *Employee       `db:"employee" json:"employee,omitempty"`
	AcademicDegree *AcademicDegree `db:"academic_degree" json:"academic_degree,omitempty"`
	AcademicRank   *AcademicRank   `db:"academic_rank" json:"academic_rank,omitempty"`
}

func (t Teacher) String() string {
	parts := make([]string, 0, 10)
	if t.EmployeeID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", t.EmployeeID))
	}
	parts = append(parts, fmt.Sprintf("email: %s", t.Email))
	parts = append(parts, fmt.Sprintf("academic rank id: %s", t.AcademicRankID))
	if t.AcademicDegreeID != nil {
		parts = append(parts, fmt.Sprintf("academic degree id: %s", t.AcademicDegreeID))
	}
	if !t.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", t.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", t.UpdatedAt.Format(db.TimeFormat)))
	}
	if t.DeletedAt != nil {
		parts = append(parts, fmt.Sprintf("deleted_at: %s", t.DeletedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Teacher{%s}", strings.Join(parts, ", "))
}

func (t *Teacher) ValidateEmail() error {
	if t.Email == "" {
		return fmt.Errorf("email required")
	}

	return nil
}

func (Teacher) TableName() string {
	return "teachers"
}
