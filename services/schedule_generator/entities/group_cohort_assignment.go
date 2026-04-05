package entities

import (
	"github.com/google/uuid"
)

type GroupCohortAssignment struct {
	ID            uuid.UUID `json:"id"`
	GroupCohortID uuid.UUID `json:"group_cohort_id"`
	DisciplineID  uuid.UUID `json:"discipline_id"`
	LessonTypeID  uuid.UUID `json:"lesson_type_id"`
}
