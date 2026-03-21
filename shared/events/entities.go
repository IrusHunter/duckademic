package events

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type AcademicRankRE struct {
	Event EventType
	ID    uuid.UUID
	Slug  string
	Title string
}

// String returns a human-readable representation of the AcademicRankRE.
// Includes event type, ID, slug and title.
func (ar AcademicRankRE) String() string {
	parts := make([]string, 0, 4)
	parts = append(parts, fmt.Sprintf("event: %s", ar.Event))
	parts = append(parts, fmt.Sprintf("id: %s", ar.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", ar.Slug))
	parts = append(parts, fmt.Sprintf("title: %s", ar.Title))
	return fmt.Sprintf("AcademicRankRE{%s}", strings.Join(parts, ", "))
}

type TeacherRE struct {
	Event          EventType
	ID             uuid.UUID
	Name           string
	AcademicRankID uuid.UUID
}

func (t TeacherRE) String() string {
	parts := make([]string, 0, 4)
	parts = append(parts, fmt.Sprintf("event: %s", t.Event))
	parts = append(parts, fmt.Sprintf("id: %s", t.ID))
	parts = append(parts, fmt.Sprintf("name: %s", t.Name))
	parts = append(parts, fmt.Sprintf("academic rank id: %s", t.AcademicRankID))
	return fmt.Sprintf("TeacherRE{%s}", strings.Join(parts, ", "))
}
