package generatorResponses

import (
	"github.com/google/uuid"
)

type Lesson struct {
	ID             uuid.UUID  `json:"id"`
	StudyLoadID    uuid.UUID  `json:"study_load_id"`
	TeacherID      uuid.UUID  `json:"teacher_id"`
	StudentGroupID uuid.UUID  `json:"student_group_id"`
	Slot           int        `json:"slot"`
	Day            int        `json:"day"`
	ClassroomID    *uuid.UUID `json:"classroom_id,omitempty"`
}
