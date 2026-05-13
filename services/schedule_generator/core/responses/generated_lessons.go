package responses

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/google/uuid"
)

type GeneratedLessons struct {
	Lessons []GeneratedLesson  `json:"lessons"`
	Errors  []UnassignedLesson `json:"errors"`
}

type GeneratedLesson struct {
	StudyLoad
	Days      []int         `json:"days"`
	Slot      int           `json:"slot"`
	Classroom *CommonEntity `json:"classroom,omitempty"`
}

func FormGeneratedLessons(lessons []*entities.Lesson) []GeneratedLesson {
	type key struct {
		TeacherID    uuid.UUID
		GroupID      uuid.UUID
		DisciplineID uuid.UUID
		LessonTypeID uuid.UUID
		Weekday      int
		Slot         int
		ClassroomID  *uuid.UUID
	}

	grouped := make(map[key]*GeneratedLesson)

	for _, lesson := range lessons {
		var classroomID *uuid.UUID
		if lesson.Classroom != nil {
			classroomID = &lesson.Classroom.ID
		}

		k := key{
			TeacherID:    lesson.Teacher.ID,
			GroupID:      lesson.StudentGroup.ID,
			DisciplineID: lesson.Discipline.ID,
			LessonTypeID: lesson.Type.ID,
			Slot:         lesson.Slot,
			Weekday:      lesson.Day % 7,
			ClassroomID:  classroomID,
		}

		if existing, ok := grouped[k]; ok {
			existing.Days = append(existing.Days, lesson.Day)
			continue
		}

		grouped[k] = &GeneratedLesson{
			StudyLoad: FormStudyLoad(lesson.StudyLoad),
			Days:      []int{lesson.Day},
			Slot:      lesson.Slot,
			Classroom: func() *CommonEntity {
				if lesson.Classroom == nil {
					return nil
				}
				return &CommonEntity{
					ID:   lesson.Classroom.ID,
					Name: lesson.Classroom.RoomNumber,
				}
			}(),
		}
	}

	generatedLessons := make([]GeneratedLesson, 0, len(grouped))
	for _, lesson := range grouped {
		generatedLessons = append(generatedLessons, *lesson)
	}

	return generatedLessons
}
