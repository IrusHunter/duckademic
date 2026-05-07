package responses

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/google/uuid"
)

type StudyLoad struct {
	ID               uuid.UUID `json:"id"`
	TeacherID        uuid.UUID `json:"teacher_id"`
	TeacherName      string    `json:"teacher_name"`
	StudentGroupID   uuid.UUID `json:"student_group_id"`
	StudentGroupName string    `json:"student_group_name"`
	DisciplineID     uuid.UUID `json:"discipline_id"`
	DisciplineName   string    `json:"discipline_name"`
	LessonTypeID     uuid.UUID `json:"lesson_type_id"`
	LessonTypeName   string    `json:"lesson_type_name"`
}

func FormStudyLoad(studyLoad *entities.StudyLoad) StudyLoad {
	return StudyLoad{
		ID:               studyLoad.ID,
		TeacherID:        studyLoad.Teacher.ID,
		TeacherName:      studyLoad.Teacher.UserName,
		StudentGroupID:   studyLoad.StudentGroup.ID,
		StudentGroupName: studyLoad.StudentGroup.Name,
		DisciplineID:     studyLoad.Discipline.ID,
		DisciplineName:   studyLoad.Discipline.Name,
		LessonTypeID:     studyLoad.Type.ID,
		LessonTypeName:   studyLoad.Type.Name,
	}
}

func FormStudyLoads(studyLoads []*entities.StudyLoad) []StudyLoad {
	result := make([]StudyLoad, 0, len(studyLoads))

	for _, studyLoad := range studyLoads {
		result = append(result, FormStudyLoad(studyLoad))
	}

	return result
}
