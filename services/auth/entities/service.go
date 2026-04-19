package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Service struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Secrete   string    `json:"secrete" db:"secrete"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (s Service) String() string {
	parts := make([]string, 0, 5)

	if s.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", s.ID))
	}

	parts = append(parts, fmt.Sprintf("name: %s", s.Name))
	parts = append(parts, fmt.Sprintf("secrete: %s", s.Secrete))

	if !s.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", s.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", s.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Service{%s}", strings.Join(parts, ", "))
}
func (Service) TableName() string {
	return "services"
}

func (Service) EntityName() string {
	return "service"
}
