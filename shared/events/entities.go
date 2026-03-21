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
func (e AcademicRankRE) String() string {
	parts := make([]string, 0, 4)
	parts = append(parts, fmt.Sprintf("event: %s", e.Event))
	parts = append(parts, fmt.Sprintf("id: %s", e.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", e.Slug))
	parts = append(parts, fmt.Sprintf("title: %s", e.Title))
	return fmt.Sprintf("AcademicRankRE{%s}", strings.Join(parts, ", "))
}
