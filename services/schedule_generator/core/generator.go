package core

import (
	"errors"
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/types"
)

type generatorData struct {
	busyGrid            [][]float32
	teacherService      services.TeacherService
	studentGroupService services.StudentGroupService
	lessonService       services.LessonService
	disciplineService   services.DisciplineService
	lessonTypeService   services.LessonTypeService
	studyLoadService    services.StudyLoadService
}

// 0 - teacher, 1 - student group, 2 - discipline, 3 - lesson type service.
func (g *generatorData) CheckServices(services []bool) error {
	checks := append(services, make([]bool, 4-len(services))...)

	if checks[0] && g.teacherService == nil {
		return fmt.Errorf("teachers not set")
	}

	if checks[1] && g.studentGroupService == nil {
		return fmt.Errorf("student groups not set")
	}

	if checks[2] && g.disciplineService == nil {
		return fmt.Errorf("disciplines not set")
	}

	if checks[3] && g.lessonTypeService == nil {
		return fmt.Errorf("lesson types not set")
	}

	return nil
}

type ScheduleGenerator struct {
	externalEntities.ScheduleGeneratorConfig
	generatorData
	errorService     components.ErrorService
	weekData         generatorData
	fullData         generatorData
	classroomService services.ClassroomService
}

func NewScheduleGenerator(cfg externalEntities.ScheduleGeneratorConfig) (*ScheduleGenerator, error) {
	scheduleGenerator := ScheduleGenerator{
		ScheduleGeneratorConfig: cfg,
	}

	index := 0
	fullBusyGrid := [][]float32{}
	for range cfg.StartDate.Weekday() {
		fullBusyGrid = append(fullBusyGrid, []float32{})
		index++
	}
	for date := cfg.StartDate; !date.After(cfg.EndDate); date = date.AddDate(0, 0, 1) {
		fullBusyGrid = append(fullBusyGrid, make([]float32, len(cfg.SlotPreference[date.Weekday()])))
		copy(fullBusyGrid[index], cfg.SlotPreference[date.Weekday()])
		index++
	}
	for range 6 - cfg.EndDate.Weekday() {
		fullBusyGrid = append(scheduleGenerator.busyGrid, []float32{})
	}
	scheduleGenerator.fullData.busyGrid = fullBusyGrid

	for i := range 7 {
		scheduleGenerator.weekData.busyGrid = append(scheduleGenerator.weekData.busyGrid,
			make([]float32, len(cfg.SlotPreference[i])))
		copy(scheduleGenerator.weekData.busyGrid[i], cfg.SlotPreference[i])
	}

	ls, err := services.NewLessonService(2) // TODO: remove and get it value from lesson type
	if err != nil {
		return nil, err
	}

	weekLS, err := services.NewLessonService(2) // TODO: remove and get it value from lesson type
	if err != nil {
		return nil, err
	}

	scheduleGenerator.fullData.lessonService = ls
	scheduleGenerator.weekData.lessonService = weekLS

	scheduleGenerator.errorService = components.NewErrorService()

	return &scheduleGenerator, nil
}

func (g *ScheduleGenerator) SetTeachers(teachers []externalEntities.Teacher) error {
	if g.teacherService != nil {
		return fmt.Errorf("teachers already set")
	}

	ts, err := services.NewTeacherService(teachers, g.busyGrid)
	if err != nil {
		return err
	}

	weekTS, err := services.NewTeacherService(teachers, g.weekData.busyGrid)
	if err != nil {
		return err
	}

	g.fullData.teacherService = ts
	g.weekData.teacherService = weekTS
	return nil
}

func (g *ScheduleGenerator) SetStudentGroups(studentGroups []types.StudentGroup) error {
	sgs, err := services.NewStudentGroupService(studentGroups, g.MaxDailyStudentLoad, g.busyGrid)
	if err != nil {
		return err
	}

	weekSGS, err := services.NewStudentGroupService(studentGroups, g.MaxDailyStudentLoad, g.weekData.busyGrid)
	if err != nil {
		return err
	}

	g.studentGroupService = sgs
	g.weekData.studentGroupService = weekSGS
	return nil
}

func (g *ScheduleGenerator) SetDisciplines(disciplines []externalEntities.Discipline) error {
	if g.disciplineService != nil {
		return fmt.Errorf("disciplines already set")
	}

	ds, err := services.NewDisciplineService(disciplines)
	if err != nil {
		return err
	}

	weekDS, err := services.NewDisciplineService(disciplines)
	if err != nil {
		return err
	}

	g.fullData.disciplineService = ds
	g.weekData.disciplineService = weekDS
	return nil
}

func (g *ScheduleGenerator) SetLessonTypes(lTypes []externalEntities.LessonType) error {
	lts, err := services.NewLessonTypeService(lTypes)
	if err != nil {
		return err
	}

	weekLTS, err := services.NewLessonTypeService(lTypes)
	if err != nil {
		return err
	}

	g.fullData.lessonTypeService = lts
	g.weekData.lessonTypeService = weekLTS
	return nil
}

func (g *ScheduleGenerator) SetLessonTypeAssignments(ltAssignments []externalEntities.LessonTypeAssignment) error {
	helperFunc := func(
		ds services.DisciplineService,
		lts services.LessonTypeService,
		assignment externalEntities.LessonTypeAssignment,
	) error {
		lessonType := lts.Find(assignment.LessonTypeID)
		if lessonType == nil {
			return fmt.Errorf("lesson type with id %q not found", assignment.LessonTypeID)
		}
		discipline := ds.Find(assignment.DisciplineID)
		if discipline == nil {
			return fmt.Errorf("discipline with id %q not found", assignment.DisciplineID)
		}

		err := discipline.AddLoad(lessonType, assignment.RequiredHours)
		if err != nil {
			return fmt.Errorf("failed to add load with id %q: %w", assignment.ID, err)
		}

		return nil
	}

	for _, assignment := range ltAssignments {
		if err := helperFunc(g.fullData.disciplineService, g.fullData.lessonTypeService, assignment); err != nil {
			return err
		}

		assignment.RequiredHours = 2
		if err := helperFunc(g.weekData.disciplineService, g.weekData.lessonTypeService, assignment); err != nil {
			return err
		}
	}

	return nil
}

func (g *ScheduleGenerator) SetStudyLoads(studyLoads []types.StudyLoad) error {
	if err := g.CheckServices([]bool{true, true, true, true}); err != nil {
		return err
	}
	if err := g.weekData.CheckServices([]bool{true, true, true, true}); err != nil {
		return err
	}

	sls, err := services.NewStudyLoadService(studyLoads, g.teacherService, g.studentGroupService,
		g.disciplineService, g.lessonTypeService)
	if err != nil {
		return err
	}

	weekSLS, err := services.NewStudyLoadService(studyLoads, g.weekData.teacherService, g.weekData.studentGroupService,
		g.weekData.disciplineService, g.weekData.lessonTypeService)
	if err != nil {
		return err
	}
	g.weekData.studentGroupService.UnbindWeeks()

	g.studyLoadService = sls
	g.weekData.studyLoadService = weekSLS
	return nil
}

func (g *ScheduleGenerator) SetClassrooms(classrooms []types.Classroom) error {
	cs, err := services.NewClassroomService(classrooms, g.busyGrid, float32(g.ClassroomOccupancy))
	if err != nil {
		return fmt.Errorf("classroom service creation fails: %w", err)
	}

	g.classroomService = cs
	return nil
}

// main function
func (g *ScheduleGenerator) GenerateSchedule() error {
	if g.studyLoadService == nil {
		return fmt.Errorf("study loads not set")
	}
	if g.weekData.studyLoadService == nil {
		return fmt.Errorf("study loads not set")
	}

	components.NewDayBlocker(g.weekData.studentGroupService.GetAll(), g.errorService).SetDayTypes()

	components.NewBoneGenerator(g.errorService, g.weekData.studyLoadService.GetAll(), g.weekData.lessonService).GenerateBoneLessons()
	g.buildLessonCarcass()

	// components.NewMissingLessonAdder(g.errorService, g.studyLoadService.GetAll(), g.lessonService).AddMissingLessons()

	// improver := components.NewImprover(g.lessonService)
	// improver.SubmitChanges() // CRUNCH - sets start slots for first lesson
	// result := true
	// currentFault := g.ScheduleFault()
	// for result {
	// 	if currentFault.Fault() <= 0.000001 {
	// 		break
	// 	}
	// 	fault := g.ScheduleFault()
	// 	if fault.Fault() < currentFault.Fault() {
	// 		improver.SubmitChanges()
	// 	}
	// 	result = improver.ImproveToNext()
	// }

	classroomAssigner := components.NewClassroomAssigner(g.classroomService.GetAll(),
		g.lessonService.Sort(g.lessonService.GetAll(),
			g.lessonService.ByLessonSlot(1,
				g.lessonService.Equal,
			),
		), g.errorService,
	)
	if err := classroomAssigner.CheckAvailability(); err != nil {
		g.errorService.AddError(components.NewUnexpectedError(
			"can't assign classrooms", "Generator", "GenerateSchedule", err,
		))
	}
	classroomAssigner.AssignClassrooms()

	if !g.errorService.IsClear() {
		return g.errorService
	}
	return nil
}

func (g *ScheduleGenerator) buildLessonCarcass() {
	lessons := g.weekData.lessonService.GetAll()
	for _, lesson := range lessons {
		teacher := g.teacherService.Find(lesson.Teacher.ID)
		studentGroup := g.studentGroupService.Find(lesson.StudentGroup.ID)
		for weekday := range 7 {
			weekLT := lesson.StudentGroup.GetTypeOfDay(weekday)
			if weekLT != nil {
				lt := studentGroup.GetTypeOfDay(weekday)
				if lt == nil {
					lt := g.lessonTypeService.Find(weekLT.ID)
					err := studentGroup.BindWeekday(lt, weekday)
					if err != nil {
						g.errorService.AddError(components.NewUnexpectedError("can't bind the lesson type to the day",
							"ScheduleGenerator", "buildLessonCarcass", err))
					}
				}
			}
		}
		discipline := g.disciplineService.Find(lesson.Discipline.ID)
		lessonType := g.lessonTypeService.Find(lesson.Type.ID)
		studyLoad := g.studyLoadService.Find(*entities.NewUnassignedLesson(
			lessonType, teacher, studentGroup, discipline,
		))

		currentWeek := 0
		outOfGrid := false
		for !outOfGrid {
			err := g.lessonService.AssignLesson(studyLoad,
				entities.NewLessonSlot(lesson.Day+currentWeek*7, lesson.Slot),
			)

			var dayErr *entities.DayOutError
			if errors.As(err, &dayErr) {
				outOfGrid = true
			}
			currentWeek++
		}
	}
}

// Rates schedule fault. Returns ScheduleFault as a result.
// Returns an empty ScheduleFault if an not enough data.
func (g *ScheduleGenerator) ScheduleFault() (result components.ScheduleFault) {
	result = components.NewScheduleFault()

	err := g.CheckServices([]bool{true, true})
	if err != nil {
		return
	}

	result.AddParameter("teacher_windows", components.NewSimpleScheduleParameter(
		float64(g.teacherService.CountWindows()), 0.1,
	))
	result.AddParameter("student_group_windows", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountWindows()), 1000,
	))
	result.AddParameter("study_load_hours_deficit", components.NewSimpleScheduleParameter(
		float64(g.studyLoadService.CountHoursDeficit()), 10,
	))
	result.AddParameter("teacher_lesson_overlapping", components.NewSimpleScheduleParameter(
		float64(g.teacherService.CountLessonOverlapping()), 1000,
	))
	result.AddParameter("student_group_lesson_overlapping", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountLessonOverlapping()), 1000,
	))
	result.AddParameter("classroom_lesson_overlapping", components.NewSimpleScheduleParameter(
		float64(g.classroomService.CountLessonOverlapping()), 1000,
	))
	result.AddParameter("student_group_overtime_lessons", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountOvertimeLessons()), 10,
	))
	result.AddParameter("student_group_invalid_lessons_by_type", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountInvalidLessonsByType()), 10,
	))
	result.AddParameter("lessons_without_classroom", components.NewSimpleScheduleParameter(
		float64(g.lessonService.CountLessonsWithoutClassroom(g.lessonService.GetAll())), 10,
	))
	result.AddParameter("classroom_with_overflow", components.NewSimpleScheduleParameter(
		float64(g.classroomService.CountOverflowLessons()), 10,
	))

	return
}

func (g *ScheduleGenerator) WriteSchedule() {
	tSchedule := make(map[*entities.Teacher]*entities.PersonalSchedule, len(g.teacherService.GetAll()))
	for i := range g.teacherService.GetAll() {
		t := g.teacherService.GetAll()[i]
		tSchedule[t] = &entities.PersonalSchedule{
			BusyGrid: &t.BusyGrid,
			Out:      "schedule/" + t.UserName + ".txt",
		}
	}

	sgSchedule := make(map[*entities.StudentGroup]*entities.PersonalSchedule, len(g.studentGroupService.GetAll()))
	for i := range g.studentGroupService.GetAll() {
		sg := g.studentGroupService.GetAll()[i]
		sgSchedule[sg] = &entities.PersonalSchedule{
			BusyGrid: &sg.BusyGrid,
			Out:      "schedule/" + sg.Name + ".txt",
		}
	}

	for _, l := range g.lessonService.GetAll() {
		tSchedule[l.Teacher].InsertLesson(l)
		sgSchedule[l.StudentGroup].InsertLesson(l)
	}

	for _, ps := range tSchedule {
		ps.WritePS(func(l *entities.Lesson) string {
			classroomStr := ""
			if l.Classroom != nil {
				classroomStr = fmt.Sprintf(", аудиторія: %s", l.Classroom.RoomNumber)
			}
			return fmt.Sprintf("дисципліна: %s, тип: %s, група: %s%s",
				l.Discipline.Name, l.Type.Name, l.StudentGroup.Name, classroomStr,
			)
		})
	}
	for _, ps := range sgSchedule {
		ps.WritePS(func(l *entities.Lesson) string {
			classroomStr := ""
			if l.Classroom != nil {
				classroomStr = fmt.Sprintf(", аудиторія: %s", l.Classroom.RoomNumber)
			}
			return fmt.Sprintf("дисципліна: %s, тип: %s, викладач: %s%s",
				l.Discipline.Name, l.Type.Name, l.Teacher.UserName, classroomStr,
			)
		})
	}
}
