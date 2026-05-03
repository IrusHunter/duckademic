package responses

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
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

func FormLessons(lessons []*entities.Lesson) []Lesson {
	result := make([]Lesson, 0, len(lessons))

	for _, lesson := range lessons {
		result = append(result, Lesson{
			ID:             lesson.ID,
			StudyLoadID:    lesson.StudyLoad.ID,
			TeacherID:      lesson.Teacher.ID,
			StudentGroupID: lesson.StudentGroup.ID,
			Slot:           lesson.Slot,
			Day:            lesson.Day,
			ClassroomID: func(c *entities.Classroom) *uuid.UUID {
				if c == nil {
					return nil
				}
				return &c.ID
			}(lesson.Classroom),
		})
	}

	return result
}
