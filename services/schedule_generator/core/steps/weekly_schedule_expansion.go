package steps

import (
	"errors"
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
)

type weeklyScheduleExpansionStep struct {
	*GeneratorContext
	canGoToTheNextStep bool
}

func NewWeeklyScheduleExpansionStep(c *GeneratorContext) PipelineStep {
	s := weeklyScheduleExpansionStep{GeneratorContext: c}

	return &s
}

func (s *weeklyScheduleExpansionStep) GetNextStep(c *GeneratorContext) PipelineStep {
	return NewFullTimeSlotAssignmentStep(c)
}
func (s *weeklyScheduleExpansionStep) CanGoToTheNextStep() error {
	if !s.canGoToTheNextStep {
		return fmt.Errorf("required action was not applied")
	}

	return nil
}
func (s *weeklyScheduleExpansionStep) InsertData(data any) error {
	return nil
}
func (s *weeklyScheduleExpansionStep) Process(components.ComponentIdentifier) (any, error) {
	s.canGoToTheNextStep = true

	lessons := s.weekData.lessonService.GetAll()
	for _, lesson := range lessons {
		teacher := s.fullData.teacherService.Find(lesson.Teacher.ID)
		studentGroup := s.fullData.studentGroupService.Find(lesson.StudentGroup.ID)

		// binding lesson type to day for actual student groups
		for weekday := range 7 {
			weekLT := lesson.StudentGroup.GetTypeOfDay(weekday)
			if weekLT != nil {
				lt := studentGroup.GetTypeOfDay(weekday)
				if lt == nil {
					lt := s.fullData.lessonTypeService.Find(weekLT.ID)
					slots := 0
					if studentGroup.GetReservedSlotsForLT(lt) == 0 {
						slots = lesson.StudentGroup.GetReservedSlotsForLT(weekLT)
					}
					studentGroup.BindWeekday(lt, weekday, slots)
				}
			}
		}

		discipline := s.fullData.disciplineService.Find(lesson.Discipline.ID)
		lessonType := s.fullData.lessonTypeService.Find(lesson.Type.ID)
		studyLoad := s.fullData.studyLoadService.Find(*entities.NewUnassignedLesson(
			lessonType, teacher, studentGroup, discipline,
		))
		classroom := func(weekC *entities.Classroom) *entities.Classroom {
			if weekC == nil {
				return nil
			}
			return s.fullData.classroomService.Find(weekC.ID)
		}(lesson.Classroom)

		// copy week lesson to all weeks
		currentWeek := 0
		outOfGrid := false
		for !outOfGrid {
			err := s.fullData.lessonService.AssignLesson(studyLoad,
				entities.NewLessonSlot(lesson.Day+currentWeek*7, lesson.Slot),
			)
			if err == nil && classroom != nil {
				fullL := s.fullData.lessonService.Select().Sort().Last()
				fullL.SetClassroom(classroom)
			}

			var dayErr *entities.DayOutError
			if errors.As(err, &dayErr) {
				outOfGrid = true
			}
			currentWeek++
		}
	}

	s.canGoToTheNextStep = true

	res := responses.GeneratedLessons{Lessons: responses.FormGeneratedLessons(s.fullData.lessonService.GetAll())}
	return res, nil
}
func (s *weeklyScheduleExpansionStep) ApplyManualChange(data map[string]string) error {
	return fmt.Errorf("you can't manually change something at weekly schedule expansion step")
}
