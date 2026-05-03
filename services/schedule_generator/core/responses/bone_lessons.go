package responses

import "github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"

type BoneLessons struct {
	Lessons []BoneLesson       `json:"bone_lessons"`
	Errors  []UnassignedLesson `json:"errors"`
}

type BoneLesson struct {
	CommonLesson
	Day       int           `json:"day"`
	Slot      int           `json:"slot"`
	Classroom *CommonEntity `json:"classroom,omitempty"`
}

func FormBoneLessons(lessons []*entities.Lesson) []BoneLesson {
	boneLessons := make([]BoneLesson, 0, len(lessons))
	for _, lesson := range lessons {
		boneLessons = append(boneLessons, BoneLesson{
			CommonLesson: CommonLesson{
				Teacher: CommonEntity{
					ID:   lesson.Teacher.ID,
					Name: lesson.Teacher.UserName,
				},
				StudentGroup: CommonEntity{
					ID:   lesson.StudentGroup.ID,
					Name: lesson.StudentGroup.Name,
				},
				Discipline: CommonEntity{
					ID:   lesson.Discipline.ID,
					Name: lesson.Discipline.Name,
				},
				LessonType: CommonEntity{
					ID:   lesson.Type.ID,
					Name: lesson.Type.Name,
				},
			},
			Day:  lesson.Day,
			Slot: lesson.Slot,
			Classroom: func() *CommonEntity {
				if lesson.Classroom == nil {
					return nil
				}
				return &CommonEntity{
					ID:   lesson.Classroom.ID,
					Name: lesson.Classroom.RoomNumber,
				}
			}(),
		})
	}

	return boneLessons
}
