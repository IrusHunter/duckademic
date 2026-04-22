package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type SemesterDiscipline struct {
	ID           uuid.UUID `db:"id" json:"id"`
	SemesterID   uuid.UUID `db:"semester_id" json:"semester_id"`
	DisciplineID uuid.UUID `db:"discipline_id" json:"discipline_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func (sd SemesterDiscipline) String() string {
	parts := make([]string, 0, 3)

	if sd.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", sd.ID))
	}

	parts = append(parts, fmt.Sprintf("semester_id: %s", sd.SemesterID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", sd.DisciplineID))

	if !sd.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", sd.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", sd.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("SemesterDiscipline{%s}", strings.Join(parts, ", "))
}

func (SemesterDiscipline) TableName() string {
	return "semester_discipline"
}
