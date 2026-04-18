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

type StudentGroupRE struct {
	Event         EventType
	ID            uuid.UUID
	Slug          string
	Name          string
	GroupCohortID uuid.UUID
}

func (sr StudentGroupRE) String() string {
	parts := make([]string, 0, 5)
	parts = append(parts, fmt.Sprintf("event: %s", sr.Event))
	parts = append(parts, fmt.Sprintf("id: %s", sr.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", sr.Slug))
	parts = append(parts, fmt.Sprintf("name: %s", sr.Name))
	parts = append(parts, fmt.Sprintf("group_cohort_id: %s", sr.GroupCohortID))
	return fmt.Sprintf("StudentGroupRE{%s}", strings.Join(parts, ", "))
}

type GroupMemberRE struct {
	Event          EventType
	ID             uuid.UUID
	StudentID      uuid.UUID
	GroupCohortID  uuid.UUID
	StudentGroupID *uuid.UUID
}

func (gr GroupMemberRE) String() string {
	parts := make([]string, 0, 5)

	parts = append(parts, fmt.Sprintf("event: %s", gr.Event))
	parts = append(parts, fmt.Sprintf("id: %s", gr.ID))
	parts = append(parts, fmt.Sprintf("student_id: %s", gr.StudentID))
	parts = append(parts, fmt.Sprintf("group_cohort_id: %s", gr.GroupCohortID))

	if gr.StudentGroupID != nil {
		parts = append(parts, fmt.Sprintf("student_group_id: %s", *gr.StudentGroupID))
	} else {
		parts = append(parts, "student_group_id: <nil>")
	}

	return fmt.Sprintf("GroupMemberRE{%s}", strings.Join(parts, ", "))
}

type GroupCohortRE struct {
	Event      EventType
	ID         uuid.UUID
	Slug       string
	Name       string
	SemesterID uuid.UUID
}

func (gc GroupCohortRE) String() string {
	parts := make([]string, 0, 5)

	parts = append(parts, fmt.Sprintf("event: %s", gc.Event))
	parts = append(parts, fmt.Sprintf("id: %s", gc.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", gc.Slug))
	parts = append(parts, fmt.Sprintf("name: %s", gc.Name))
	parts = append(parts, fmt.Sprintf("semester_id: %s", gc.SemesterID))

	return fmt.Sprintf("GroupCohortRE{%s}", strings.Join(parts, ", "))
}

type TeacherLoadRE struct {
	Event        EventType
	ID           uuid.UUID
	TeacherID    uuid.UUID
	DisciplineID uuid.UUID
	LessonTypeID uuid.UUID
	GroupCount   int
}

func (tl TeacherLoadRE) String() string {
	parts := make([]string, 0, 7)

	parts = append(parts, fmt.Sprintf("event: %s", tl.Event))
	parts = append(parts, fmt.Sprintf("id: %s", tl.ID))
	parts = append(parts, fmt.Sprintf("teacher_id: %s", tl.TeacherID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", tl.DisciplineID))
	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", tl.LessonTypeID))
	parts = append(parts, fmt.Sprintf("group_count: %d", tl.GroupCount))

	return fmt.Sprintf("TeacherLoadRE{%s}", strings.Join(parts, ", "))
}

type GroupCohortAssignmentRE struct {
	Event         EventType
	ID            uuid.UUID
	GroupCohortID uuid.UUID
	DisciplineID  uuid.UUID
	LessonTypeID  uuid.UUID
}

func (gca GroupCohortAssignmentRE) String() string {
	parts := make([]string, 0, 4)

	parts = append(parts, fmt.Sprintf("event: %s", gca.Event))
	parts = append(parts, fmt.Sprintf("id: %s", gca.ID))
	parts = append(parts, fmt.Sprintf("group_cohort_id: %s", gca.GroupCohortID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", gca.DisciplineID))
	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", gca.LessonTypeID))

	return fmt.Sprintf("GroupCohortAssignmentRE{%s}", strings.Join(parts, ", "))
}

type ClassroomRE struct {
	Event    EventType
	ID       uuid.UUID
	Slug     string
	Number   string
	Capacity int
}

func (c ClassroomRE) String() string {
	parts := make([]string, 0, 6)

	parts = append(parts, fmt.Sprintf("event: %s", c.Event))
	parts = append(parts, fmt.Sprintf("id: %s", c.ID))
	parts = append(parts, fmt.Sprintf("slug: %s", c.Slug))
	parts = append(parts, fmt.Sprintf("number: %s", c.Number))
	parts = append(parts, fmt.Sprintf("capacity: %d", c.Capacity))

	return fmt.Sprintf("ClassroomRE{%s}", strings.Join(parts, ", "))
}

type AccessPermissionRE struct {
	Name string
}

func (a AccessPermissionRE) String() string {
	return fmt.Sprintf("AccessPermission{%s}", a.Name)
}
