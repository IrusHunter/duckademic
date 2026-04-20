package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type ServicePermissions struct {
	ID           uuid.UUID `json:"id" db:"id"`
	ServiceID    uuid.UUID `json:"service_id" db:"service_id"`
	PermissionID uuid.UUID `json:"permission_id" db:"permission_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (sp ServicePermissions) String() string {
	parts := make([]string, 0, 5)

	if sp.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", sp.ID))
	}

	parts = append(parts, fmt.Sprintf("service_id: %s", sp.ServiceID))
	parts = append(parts, fmt.Sprintf("permission_id: %s", sp.PermissionID))

	if !sp.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", sp.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", sp.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("ServicePermissions{%s}", strings.Join(parts, ", "))
}

func (ServicePermissions) TableName() string {
	return "service_permissions"
}

func (ServicePermissions) EntityName() string {
	return "service_permission"
}
