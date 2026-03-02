package services

import (
	"fmt"
	"slices"

	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
)

// StudyLoadService aggregates and manages study loads that the generator works with.
type StudyLoadService interface {
	GetAll() []*entities.StudyLoad                      // Returns a slice with all study loads as pointers.
	CountHoursDeficit() int                             // Returns the number of missing study hours.
	Find(entities.UnassignedLesson) *entities.StudyLoad // Returns a pointer to the teacher with the given data.
}

// NewStudyLoadService creates a new StudyLoadService basic instance.
//
// It requires an array of database study loads (sl), teacher, student group, discipline,
// and lesson type services (ts, sgs, ds, and lts).
//
// Returns an error if any study load is an invalid model.
func NewStudyLoadService(
	sl []types.StudyLoad,
	ts TeacherService,
	sgs StudentGroupService,
	ds DisciplineService,
	lts LessonTypeService,
) (StudyLoadService, error) {
	sls := &studyLoadService{}

	for _, studyLoad := range sl {

		teacher := ts.Find(studyLoad.TeacherID)
		if teacher == nil {
			return nil, fmt.Errorf("teacher %s not found", studyLoad.TeacherID)
		}

		for _, disciplineLoad := range studyLoad.Disciplines {
			discipline := ds.Find(disciplineLoad.DisciplineID)
			if discipline == nil {
				return nil, fmt.Errorf("discipline %s not found", disciplineLoad.DisciplineID)
			}
			lessonType := lts.Find(disciplineLoad.LessonTypeID)
			if lessonType == nil {
				return nil, fmt.Errorf("lesson type %s not found", disciplineLoad.LessonTypeID)
			}

			discipline.AddLoad(lessonType, disciplineLoad.Hours) //TODO: do it at discipline service

			studentGroups := make([]*entities.StudentGroup, len(disciplineLoad.GroupsID))
			for j, studentGroupID := range disciplineLoad.GroupsID {
				studentGroup := sgs.Find(studentGroupID)
				if studentGroup == nil {
					return nil, fmt.Errorf("student group %s not found", studentGroupID)
				}
				for week := range lessonType.Weeks {
					studentGroup.BindWeek(lessonType, week)
				}

				studentGroups[j] = studentGroup

				studyLoad := entities.NewStudyLoad(*entities.NewUnassignedLesson(
					lessonType, teacher, studentGroup, discipline),
				)
				sls.loads = append(sls.loads, studyLoad)
				teacher.AddLoad(studyLoad)
				studentGroup.AddLoad(studyLoad)
			}

			// if err := discipline.AddLoad(teacher, disciplineLoad.Hours, studentGroups, lessonType); err != nil {
			// 	return err
			// }
		}
	}

	return sls, nil
}

// studyLoadService is the basic implementation of the StudyLoadService interface.
type studyLoadService struct {
	loads []*entities.StudyLoad
}

func (s *studyLoadService) GetAll() []*entities.StudyLoad {
	return s.loads
}
func (s *studyLoadService) CountHoursDeficit() (result int) {
	for _, load := range s.loads {
		result += load.CountHoursDeficit()
	}
	return
}
func (s *studyLoadService) Find(ul entities.UnassignedLesson) *entities.StudyLoad {
	ind := slices.IndexFunc(s.loads, func(load *entities.StudyLoad) bool {
		return ul == load.UnassignedLesson
	})

	if ind == -1 {
		return nil
	}
	return s.loads[ind]
}
