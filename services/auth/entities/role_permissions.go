package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type RolePermissions struct {
	ID           uuid.UUID `json:"id" db:"id"`
	RoleID       uuid.UUID `json:"role_id" db:"role_id"`
	PermissionID uuid.UUID `json:"permission_id" db:"permission_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (rp RolePermissions) String() string {
	parts := make([]string, 0, 5)

	if rp.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", rp.ID))
	}

	parts = append(parts, fmt.Sprintf("role_id: %s", rp.RoleID))
	parts = append(parts, fmt.Sprintf("permission_id: %s", rp.PermissionID))

	if !rp.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", rp.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", rp.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("RolePermissions{%s}", strings.Join(parts, ", "))
}

func (RolePermissions) TableName() string {
	return "role_permissions"
}

func (RolePermissions) EntityName() string {
	return "role_permissions"
}
