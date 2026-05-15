package steps

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"
	"github.com/google/uuid"
)

type ServiceEnum int

const (
	TeacherServiceE ServiceEnum = iota
	StudentGroupServiceE
	LessonServiceE
	DisciplineServiceE
	LessonTypeServiceE
	StudyLoadServiceE
	ClassroomServiceE
)

type generatorData struct {
	busyGrid            [][]float32
	teacherService      services.TeacherService
	studentGroupService services.StudentGroupService
	lessonService       services.LessonService
	disciplineService   services.DisciplineService
	lessonTypeService   services.LessonTypeService
	studyLoadService    services.StudyLoadService
	classroomService    services.ClassroomService
}

type GeneratorContext struct {
	config             externalEntities.ScheduleGeneratorConfig
	weekData           generatorData
	fullData           generatorData
	floatLessonService services.LessonService
	allDataReceived    bool
}

func NewGeneratorContext(cfg externalEntities.ScheduleGeneratorConfig) (*GeneratorContext, error) {
	c := &GeneratorContext{config: cfg}

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
		fullBusyGrid = append(fullBusyGrid, []float32{})
	}
	c.fullData.busyGrid = fullBusyGrid

	for i := range 7 {
		c.weekData.busyGrid = append(c.weekData.busyGrid,
			make([]float32, len(cfg.SlotPreference[i])))
		copy(c.weekData.busyGrid[i], cfg.SlotPreference[i])
	}

	ls := services.NewLessonService([]*entities.Lesson{})
	weekLS := services.NewLessonService([]*entities.Lesson{})
	c.fullData.lessonService = ls
	c.weekData.lessonService = weekLS
	c.floatLessonService = services.NewLessonService([]*entities.Lesson{})

	return c, nil
}

func (c *GeneratorContext) CheckServices(ss []ServiceEnum) error {
	for _, s := range ss {
		switch s {
		case TeacherServiceE:
			if c.weekData.teacherService == nil {
				return fmt.Errorf("week teacher service didn't set")
			}
			if c.fullData.teacherService == nil {
				return fmt.Errorf("full teacher service didn't set")
			}

		case StudentGroupServiceE:
			if c.weekData.studentGroupService == nil {
				return fmt.Errorf("week student group service didn't set")
			}
			if c.fullData.studentGroupService == nil {
				return fmt.Errorf("full student group service didn't set")
			}

		case LessonServiceE:
			if c.weekData.lessonService == nil {
				return fmt.Errorf("week lesson service didn't set")
			}
			if c.fullData.lessonService == nil {
				return fmt.Errorf("full lesson service didn't set")
			}

		case DisciplineServiceE:
			if c.weekData.disciplineService == nil {
				return fmt.Errorf("week discipline service didn't set")
			}
			if c.fullData.disciplineService == nil {
				return fmt.Errorf("full discipline service didn't set")
			}

		case LessonTypeServiceE:
			if c.weekData.lessonTypeService == nil {
				return fmt.Errorf("week lesson type service didn't set")
			}
			if c.fullData.lessonTypeService == nil {
				return fmt.Errorf("full lesson type service didn't set")
			}

		case StudyLoadServiceE:
			if c.weekData.studyLoadService == nil {
				return fmt.Errorf("week study load service didn't set")
			}
			if c.fullData.studyLoadService == nil {
				return fmt.Errorf("full study load service didn't set")
			}

		case ClassroomServiceE:
			if c.weekData.classroomService == nil {
				return fmt.Errorf("week classroom service didn't set")
			}
			if c.fullData.classroomService == nil {
				return fmt.Errorf("full classroom service didn't set")
			}

		default:
			return fmt.Errorf("unknown service enum")
		}
	}
	return nil
}

func (c *GeneratorContext) SetTeachers(teachers []externalEntities.Teacher) error {
	if c.fullData.teacherService != nil {
		return fmt.Errorf("teachers already set")
	}

	helper := func(busyGrid [][]float32) ([]*entities.Teacher, error) {
		res := make([]*entities.Teacher, 0, len(teachers))
		for _, t := range teachers {
			teacher := entities.NewDefaultTeacher(t.ID, t.Name, t.Priority, entities.NewBusyGrid(busyGrid))

			for day, priorities := range t.SlotsPriorities {
				for slot, priority := range priorities {
					if priority < 0.000_000_001 {
						timeSlot := entities.NewLessonSlot(day, slot)
						err := teacher.BlockWeekdaySlotForAllWeeks(timeSlot)
						if err != nil {
							return nil, fmt.Errorf("invalid %s at %s: %w", timeSlot, t.Name, err)
						}
					}
				}
			}

			res = append(res, teacher)
		}

		return res, nil
	}

	fullDataT, err := helper(c.fullData.busyGrid)
	if err != nil {
		return err
	}
	weekDataT, err := helper(c.weekData.busyGrid)
	if err != nil {
		return err
	}

	c.fullData.teacherService = services.NewTeacherService(fullDataT)
	c.weekData.teacherService = services.NewTeacherService(weekDataT)

	return nil
}

func (c *GeneratorContext) SetStudentGroups(
	groupCohorts []externalEntities.GroupCohort,
	groupCohortAssignments []externalEntities.GroupCohortAssignment,
) error {
	if c.fullData.studentGroupService != nil {
		return fmt.Errorf("student groups already set")
	}
	if err := c.CheckServices([]ServiceEnum{DisciplineServiceE, LessonTypeServiceE}); err != nil {
		return err
	}

	groupCohortsMap := make(map[uuid.UUID]externalEntities.GroupCohort, len(groupCohorts))
	studentGroups := []externalEntities.StudentGroup{}

	for _, groupCohort := range groupCohorts {
		groupCohortsMap[groupCohort.ID] = groupCohort
		for _, studentGroup := range groupCohort.Groups {
			studentGroups = append(studentGroups, studentGroup)
		}
	}

	sgs, err := services.NewStudentGroupService(studentGroups, c.config.MaxDailyStudentLoad, c.fullData.busyGrid)
	if err != nil {
		return err
	}
	weekSGS, _ := services.NewStudentGroupService(studentGroups, c.config.MaxDailyStudentLoad, c.weekData.busyGrid)

	helper := func(ds services.DisciplineService, lts services.LessonTypeService, sgs services.StudentGroupService,
	) ([]*entities.StudyLoad, error) {
		studyLoads := []*entities.StudyLoad{}

		for _, groupCohortAssignment := range groupCohortAssignments {
			lessonType := lts.Find(groupCohortAssignment.LessonTypeID)
			if lessonType == nil {
				return nil, fmt.Errorf("lesson type with id %q not found", groupCohortAssignment.LessonTypeID)
			}

			discipline := ds.Find(groupCohortAssignment.DisciplineID)
			if discipline == nil {
				return nil, fmt.Errorf("at group cohort with id %q discipline with id %q not found",
					groupCohortAssignment.GroupCohortID, groupCohortAssignment.DisciplineID)
			}

			groupCohort, ok := groupCohortsMap[groupCohortAssignment.GroupCohortID]
			if !ok {
				return nil, fmt.Errorf("group cohort with id %q not found", groupCohortAssignment.GroupCohortID)
			}

			for _, studentGroup := range groupCohort.Groups {
				studentGroup := sgs.Find(studentGroup.ID)
				if studentGroup == nil {
					panic("student group already set but not found")
				}

				for week := range lessonType.Weeks {
					studentGroup.BindWeek(lessonType, week)
				}

				studyLoad := entities.NewStudyLoad(
					*entities.NewUnassignedLesson(lessonType, nil, studentGroup, discipline),
				)
				studyLoads = append(studyLoads, studyLoad)
				studentGroup.AddLoad(studyLoad)
			}
		}
		return studyLoads, nil
	}

	studyLoads, err := helper(c.fullData.disciplineService, c.fullData.lessonTypeService, sgs)
	if err != nil {
		return err
	}
	weekSL, _ := helper(c.weekData.disciplineService, c.weekData.lessonTypeService, weekSGS)

	c.fullData.studentGroupService = sgs
	weekSGS.UnbindWeeks()
	c.weekData.studentGroupService = weekSGS
	c.fullData.studyLoadService, _ = services.NewStudyLoadService(studyLoads)
	c.weekData.studyLoadService, _ = services.NewStudyLoadService(weekSL)
	return nil
}

func (c *GeneratorContext) SetDisciplines(disciplines []externalEntities.Discipline) error {
	if c.fullData.disciplineService != nil {
		return fmt.Errorf("disciplines already set")
	}

	ds, err := services.NewDisciplineService(disciplines)
	if err != nil {
		return err
	}

	c.fullData.disciplineService = ds
	c.weekData.disciplineService, _ = services.NewDisciplineService(disciplines)
	return nil
}

func (c *GeneratorContext) SetLessonTypes(lTypes []externalEntities.LessonType) error {
	if c.fullData.lessonTypeService != nil {
		return fmt.Errorf("lesson types already set")
	}

	lts, err := services.NewLessonTypeService(lTypes)
	if err != nil {
		return err
	}

	c.fullData.lessonTypeService = lts
	c.weekData.lessonTypeService, _ = services.NewLessonTypeService(lTypes)
	return nil
}

func (c *GeneratorContext) SetLessonTypeAssignments(ltAssignments []externalEntities.LessonTypeAssignment) error {
	if err := c.CheckServices([]ServiceEnum{DisciplineServiceE, LessonServiceE}); err != nil {
		return err
	}

	helper := func(ds services.DisciplineService, lts services.LessonTypeService) error {
		for _, assignment := range ltAssignments {
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
		}
		return nil
	}

	if err := helper(c.fullData.disciplineService, c.fullData.lessonTypeService); err != nil {
		return err
	}
	helper(c.weekData.disciplineService, c.weekData.lessonTypeService)

	return nil
}

func (c *GeneratorContext) SetStudyLoads(teacherLoads []externalEntities.TeacherLoad) error {
	if err := c.CheckServices(
		[]ServiceEnum{DisciplineServiceE, LessonServiceE, StudyLoadServiceE, StudentGroupServiceE, TeacherServiceE},
	); err != nil {
		return err
	}

	type key struct {
		LessonTypeID uuid.UUID
		DisciplineID uuid.UUID
	}

	helper := func(gd *generatorData) error {
		teacherLoadsMap := map[key][]externalEntities.TeacherLoad{}
		for _, teacherLoad := range teacherLoads {
			key := key{
				LessonTypeID: teacherLoad.LessonTypeID,
				DisciplineID: teacherLoad.DisciplineID,
			}
			_, ok := teacherLoadsMap[key]
			if !ok {
				teacherLoadsMap[key] = []externalEntities.TeacherLoad{}
			}
			teacherLoadsMap[key] = append(teacherLoadsMap[key], teacherLoad)
		}

		studyLoads := gd.studyLoadService.GetAll()
		for _, studyLoad := range studyLoads {
			key := key{
				LessonTypeID: studyLoad.Type.ID,
				DisciplineID: studyLoad.Discipline.ID,
			}

			tLoads, ok := teacherLoadsMap[key]
			if !ok {
				return fmt.Errorf("teacher load for %s %s not found", studyLoad.Discipline.Name, studyLoad.Type.Name)
			}
			for {
				if len(tLoads) == 0 {
					return fmt.Errorf("not enough groups in teacher loads for %s %s", studyLoad.Discipline.Name, studyLoad.Type.Name)
				}
				tLoad := tLoads[0]
				if tLoad.GroupCount == 0 {
					tLoads = tLoads[1:]
					continue
				}
				break
			}
			tLoads[0].GroupCount -= 1

			teacher := gd.teacherService.Find(tLoads[0].TeacherID)
			if teacher == nil {
				return fmt.Errorf("teacher with id %q not found", tLoads[0].TeacherID)
			}
			studyLoad.Teacher = teacher
			teacher.AddLoad(studyLoad)

			teacherLoadsMap[key] = tLoads
		}
		return nil
	}

	if err := helper(&c.fullData); err != nil {
		return err
	}
	helper(&c.weekData)
	c.allDataReceived = true

	return nil
}

func (c *GeneratorContext) SetClassrooms(classrooms []externalEntities.Classroom) error {
	if c.fullData.classroomService != nil {
		return fmt.Errorf("classrooms already set")
	}

	cs, err := services.NewClassroomService(classrooms, c.fullData.busyGrid, float32(c.config.ClassroomOccupancy))
	if err != nil {
		return fmt.Errorf("classroom service creation fails: %w", err)
	}

	c.fullData.classroomService = cs
	c.weekData.classroomService, _ = services.NewClassroomService(
		classrooms, c.weekData.busyGrid, float32(c.config.ClassroomOccupancy))
	return nil
}

func (c *GeneratorContext) AllDataReceived() bool {
	if c.CheckServices([]ServiceEnum{ClassroomServiceE}) != nil {
		return false
	}

	return c.allDataReceived
}

func (c *GeneratorContext) ExtractStudyLoads() ([]responses.StudyLoad, error) {
	return responses.FormStudyLoads(c.fullData.studyLoadService.GetAll()), nil
}

func (c *GeneratorContext) ExtractLessons() ([]responses.Lesson, error) {
	return responses.FormLessons(append(c.fullData.lessonService.GetAll(), c.floatLessonService.GetAll()...)), nil
}

func (c *GeneratorContext) GetFault() responses.Fault {
	fault := responses.Fault{}

	teacherWindows, teacherWindowsC := c.fullData.teacherService.GetWindows()
	fault.TeacherWindows = responses.FaultParam[responses.TeacherWindow]{
		Value:      float64(teacherWindowsC),
		Defections: teacherWindows,
	}
	fault.TotalValue += fault.TeacherWindows.Value

	studentGroupWindows, studentGroupWindowsC := c.fullData.studentGroupService.GetWindows()
	fault.StudentGroupWindows = responses.FaultParam[responses.StudentGroupWindow]{
		Value:      float64(studentGroupWindowsC),
		Defections: studentGroupWindows,
	}
	fault.TotalValue += fault.StudentGroupWindows.Value

	teacherLessonOverlapping, teacherLessonOverlappingC := c.fullData.teacherService.GetLessonOverlapping()
	fault.TeacherLessonOverlapping = responses.FaultParam[responses.TeacherLessonOverlap]{
		Value:      float64(teacherLessonOverlappingC),
		Defections: teacherLessonOverlapping,
	}
	fault.TotalValue += fault.TeacherLessonOverlapping.Value

	studentGroupLessonOverlapping, studentGroupLessonOverlappingC := c.fullData.studentGroupService.GetLessonOverlapping()
	fault.StudentGroupLessonOverlapping = responses.FaultParam[responses.StudentGroupLessonOverlap]{
		Value:      float64(studentGroupLessonOverlappingC),
		Defections: studentGroupLessonOverlapping,
	}
	fault.TotalValue += fault.StudentGroupLessonOverlapping.Value

	classroomLessonOverlapping, classroomLessonOverlappingC := c.fullData.classroomService.GetLessonOverlapping()
	fault.ClassroomLessonOverlapping = responses.FaultParam[responses.ClassroomLessonOverlap]{
		Value:      float64(classroomLessonOverlappingC),
		Defections: classroomLessonOverlapping,
	}
	fault.TotalValue += fault.ClassroomLessonOverlapping.Value

	invalidLessonByType, invalidLessonByTypeC := c.fullData.studentGroupService.GetInvalidLessonsByType()
	fault.StudentGroupInvalidLessons = responses.FaultParam[responses.StudentGroupInvalidLesson]{
		Value:      float64(invalidLessonByTypeC),
		Defections: invalidLessonByType,
	}
	fault.TotalValue += fault.StudentGroupInvalidLessons.Value

	studyLoadHoursDeficits, studyLoadHoursDeficitsC := c.fullData.studyLoadService.GetHoursDeficit()
	fault.StudyLoadHoursDeficit = responses.FaultParam[responses.StudyLoadHoursDeficit]{
		Value:      float64(studyLoadHoursDeficitsC),
		Defections: studyLoadHoursDeficits,
	}
	fault.TotalValue += fault.StudyLoadHoursDeficit.Value

	classroomOverflow, classroomOverflowC := c.fullData.classroomService.GetOverflowLessons()
	fault.ClassroomWithOverflow = responses.FaultParam[responses.ClassroomOverflow]{
		Value:      float64(classroomOverflowC),
		Defections: classroomOverflow,
	}
	fault.TotalValue += fault.ClassroomWithOverflow.Value

	lessonsWithoutClassroom := c.fullData.lessonService.GetLessonsWithoutClassroom()
	floatingLessonsWithoutC := c.floatLessonService.GetLessonsWithoutClassroom()
	allLessonsWithoutC := append(lessonsWithoutClassroom, floatingLessonsWithoutC...)
	fault.LessonsWithoutClassroom = responses.FaultParam[responses.GeneratedLesson]{
		Value:      float64(len(allLessonsWithoutC)),
		Defections: responses.FormGeneratedLessons(allLessonsWithoutC),
	}
	fault.TotalValue += fault.LessonsWithoutClassroom.Value

	c.WriteSchedule()

	return fault
}

func (c *GeneratorContext) WriteSchedule() {
	tSchedule := make(map[*entities.Teacher]*entities.PersonalSchedule, len(c.fullData.teacherService.GetAll()))
	for i := range c.fullData.teacherService.GetAll() {
		t := c.fullData.teacherService.GetAll()[i]
		tSchedule[t] = &entities.PersonalSchedule{
			BusyGrid: &t.BusyGrid,
			Out:      "schedule-output/" + t.UserName + ".txt",
		}
	}

	sgSchedule := make(map[*entities.StudentGroup]*entities.PersonalSchedule, len(c.fullData.studentGroupService.GetAll()))
	for i := range c.fullData.studentGroupService.GetAll() {
		sg := c.fullData.studentGroupService.GetAll()[i]
		sgSchedule[sg] = &entities.PersonalSchedule{
			BusyGrid: &sg.BusyGrid,
			Out:      "schedule-output/" + sg.Name + ".txt",
		}
	}

	lessons := append(c.fullData.lessonService.GetAll(), c.floatLessonService.GetAll()...)
	for _, l := range lessons {
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
