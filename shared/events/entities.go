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
	Slug           string
	Name           string
	AcademicRankID uuid.UUID
}

func (t TeacherRE) String() string {
	parts := make([]string, 0, 4)
	parts = append(parts, fmt.Sprintf("event: %s", t.Event))
	parts = append(parts, fmt.Sprintf("id: %s", t.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", t.Slug))
	parts = append(parts, fmt.Sprintf("name: %s", t.Name))
	parts = append(parts, fmt.Sprintf("academic rank id: %s", t.AcademicRankID))
	return fmt.Sprintf("TeacherRE{%s}", strings.Join(parts, ", "))
}

type StudentRE struct {
	Event      EventType
	ID         uuid.UUID
	Slug       string
	Name       string
	SemesterID uuid.UUID
}

func (s StudentRE) String() string {
	parts := make([]string, 0, 4)
	parts = append(parts, fmt.Sprintf("event: %s", s.Event))
	parts = append(parts, fmt.Sprintf("id: %s", s.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", s.Slug))
	parts = append(parts, fmt.Sprintf("name: %s", s.Name))
	parts = append(parts, fmt.Sprintf("semester_id: %s", s.SemesterID))
	return fmt.Sprintf("StudentRE{%s}", strings.Join(parts, ", "))
}

type LessonTypeRE struct {
	Event      EventType
	ID         uuid.UUID
	Slug       string
	Name       string
	HoursValue int
}

func (l LessonTypeRE) String() string {
	parts := make([]string, 0, 5)
	parts = append(parts, fmt.Sprintf("event: %s", l.Event))
	parts = append(parts, fmt.Sprintf("id: %s", l.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", l.Slug))
	parts = append(parts, fmt.Sprintf("name: %s", l.Name))
	parts = append(parts, fmt.Sprintf("hours_value: %d", l.HoursValue))
	return fmt.Sprintf("LessonTypeRE{%s}", strings.Join(parts, ", "))
}

type DisciplineRE struct {
	Event EventType
	ID    uuid.UUID
	Slug  string
	Name  string
}

func (d DisciplineRE) String() string {
	parts := make([]string, 0, 4)
	parts = append(parts, fmt.Sprintf("event: %s", d.Event))
	parts = append(parts, fmt.Sprintf("id: %s", d.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", d.Slug))
	parts = append(parts, fmt.Sprintf("name: %s", d.Name))
	return fmt.Sprintf("DisciplineRE{%s}", strings.Join(parts, ", "))
}

type LessonTypeAssignmentRE struct {
	Event         EventType
	ID            uuid.UUID
	LessonTypeID  uuid.UUID
	DisciplineID  uuid.UUID
	RequiredHours int
}

func (lta LessonTypeAssignmentRE) String() string {
	parts := make([]string, 0, 5)
	parts = append(parts, fmt.Sprintf("event: %s", lta.Event))
	parts = append(parts, fmt.Sprintf("id: %s", lta.ID))
	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", lta.LessonTypeID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", lta.DisciplineID))
	parts = append(parts, fmt.Sprintf("required_hours: %d", lta.RequiredHours))
	return fmt.Sprintf("LessonTypeAssignmentRE{%s}", strings.Join(parts, ", "))
}

type SemesterRE struct {
	Event        EventType
	ID           uuid.UUID
	Slug         string
	CurriculumID uuid.UUID
	Number       int
}

func (sr SemesterRE) String() string {
	parts := make([]string, 0, 5)
	parts = append(parts, fmt.Sprintf("event: %s", sr.Event))
	parts = append(parts, fmt.Sprintf("id: %s", sr.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", sr.Slug))
	parts = append(parts, fmt.Sprintf("curriculum_id: %s", sr.CurriculumID))
	parts = append(parts, fmt.Sprintf("number: %d", sr.Number))
	return fmt.Sprintf("SemesterRE{%s}", strings.Join(parts, ", "))
}
