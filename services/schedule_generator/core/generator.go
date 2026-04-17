package core

import (
	"errors"
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"
	"github.com/google/uuid"
)

type GeneratorStep string

const (
	Setup                               GeneratorStep = "SETUP"
	DayBlocking                         GeneratorStep = "DAY_BLOCKING"
	BoneLessonBuilding                  GeneratorStep = "BONE_LESSON_BUILDING"
	ToBoneLessonsClassroomAssigning     GeneratorStep = "TO_BONE_LESSONS_CLASSROOM_ASSIGNING"
	LessonSkeletonBuilding              GeneratorStep = "LESSON_SKELETON_BUILDING"
	FloatingLessonAdding                GeneratorStep = "FLOATING_LESSON_ADDING"
	ToFloatingLessonsClassroomAssigning GeneratorStep = "TO_FLOATING_LESSONS_CLASSROOM_ASSIGNING"
	Extraction                          GeneratorStep = "EXTRACTION"
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
	weekData              generatorData
	fullData              generatorData
	floatingLessonService services.LessonService
	currentStep           GeneratorStep
	canGoToTheNextStep    bool
}

func NewScheduleGenerator(cfg externalEntities.ScheduleGeneratorConfig) (*ScheduleGenerator, error) {
	scheduleGenerator := ScheduleGenerator{
		ScheduleGeneratorConfig: cfg,
		currentStep:             Setup,
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
		fullBusyGrid = append(fullBusyGrid, []float32{})
	}
	scheduleGenerator.fullData.busyGrid = fullBusyGrid

	for i := range 7 {
		scheduleGenerator.weekData.busyGrid = append(scheduleGenerator.weekData.busyGrid,
			make([]float32, len(cfg.SlotPreference[i])))
		copy(scheduleGenerator.weekData.busyGrid[i], cfg.SlotPreference[i])
	}

	ls, err := services.NewLessonService()
	if err != nil {
		return nil, err
	}
	weekLS, _ := services.NewLessonService()
	scheduleGenerator.fullData.lessonService = ls
	scheduleGenerator.weekData.lessonService = weekLS

	return &scheduleGenerator, nil
}

func (g *ScheduleGenerator) SetTeachers(teachers []externalEntities.Teacher) error {
	if g.fullData.teacherService != nil {
		return fmt.Errorf("teachers already set")
	}

	ts, err := services.NewTeacherService(teachers, g.fullData.busyGrid)
	if err != nil {
		return err
	}
	g.fullData.teacherService = ts
	g.weekData.teacherService, _ = services.NewTeacherService(teachers, g.weekData.busyGrid)

	return nil
}

func (g *ScheduleGenerator) SetStudentGroups(
	groupCohorts []externalEntities.GroupCohort,
	groupCohortAssignments []externalEntities.GroupCohortAssignment,
) error {
	if g.fullData.studentGroupService != nil {
		return fmt.Errorf("student groups already set")
	}
	if g.fullData.disciplineService == nil {
		return fmt.Errorf("disciplines did't set")
	}
	if g.fullData.lessonTypeService == nil {
		return fmt.Errorf("lesson types did't set")
	}

	groupCohortsMap := make(map[uuid.UUID]externalEntities.GroupCohort, len(groupCohorts))
	studentGroups := []externalEntities.StudentGroup{}

	for _, groupCohort := range groupCohorts {
		groupCohortsMap[groupCohort.ID] = groupCohort
		for _, studentGroup := range groupCohort.Groups {
			studentGroups = append(studentGroups, studentGroup)
		}
	}

	sgs, err := services.NewStudentGroupService(studentGroups, g.MaxDailyStudentLoad, g.fullData.busyGrid)
	if err != nil {
		return err
	}
	weekSGS, _ := services.NewStudentGroupService(studentGroups, g.MaxDailyStudentLoad, g.weekData.busyGrid)

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
				return nil, fmt.Errorf("discipline with id %q not found", groupCohortAssignment.DisciplineID)
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

	studyLoads, err := helper(g.fullData.disciplineService, g.fullData.lessonTypeService, sgs)
	if err != nil {
		return err
	}
	weekSL, _ := helper(g.weekData.disciplineService, g.weekData.lessonTypeService, weekSGS)

	g.fullData.studentGroupService = sgs
	weekSGS.UnbindWeeks()
	g.weekData.studentGroupService = weekSGS
	g.fullData.studyLoadService, _ = services.NewStudyLoadService(studyLoads)
	g.weekData.studyLoadService, _ = services.NewStudyLoadService(weekSL)
	g.weekData.disciplineService.CutLoadTo(2)
	return nil
}

func (g *ScheduleGenerator) SetDisciplines(disciplines []externalEntities.Discipline) error {
	if g.fullData.disciplineService != nil {
		return fmt.Errorf("disciplines already set")
	}

	ds, err := services.NewDisciplineService(disciplines)
	if err != nil {
		return err
	}

	g.fullData.disciplineService = ds
	g.weekData.disciplineService, _ = services.NewDisciplineService(disciplines)
	return nil
}

func (g *ScheduleGenerator) SetLessonTypes(lTypes []externalEntities.LessonType) error {
	if g.fullData.lessonTypeService != nil {
		return fmt.Errorf("lesson types already set")
	}

	lts, err := services.NewLessonTypeService(lTypes)
	if err != nil {
		return err
	}

	g.fullData.lessonTypeService = lts
	g.weekData.lessonTypeService, _ = services.NewLessonTypeService(lTypes)
	return nil
}

func (g *ScheduleGenerator) SetLessonTypeAssignments(ltAssignments []externalEntities.LessonTypeAssignment) error {
	if g.fullData.lessonTypeService == nil {
		return fmt.Errorf("lesson types didn't set")
	}
	if g.fullData.disciplineService == nil {
		return fmt.Errorf("disciplines didn't set")
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

	if err := helper(g.fullData.disciplineService, g.fullData.lessonTypeService); err != nil {
		return err
	}
	helper(g.weekData.disciplineService, g.weekData.lessonTypeService)

	return nil
}

func (g *ScheduleGenerator) SetStudyLoads(teacherLoads []externalEntities.TeacherLoad) error {
	if err := g.fullData.CheckServices([]bool{true, true, true, true}); err != nil {
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

	if err := helper(&g.fullData); err != nil {
		return err
	}
	helper(&g.weekData)

	g.canGoToTheNextStep = true
	return nil
}

func (g *ScheduleGenerator) SetClassrooms(classrooms []externalEntities.Classroom) error {
	if g.fullData.classroomService != nil {
		return fmt.Errorf("classrooms already set")
	}

	cs, err := services.NewClassroomService(classrooms, g.fullData.busyGrid, float32(g.ClassroomOccupancy))
	if err != nil {
		return fmt.Errorf("classroom service creation fails: %w", err)
	}

	g.fullData.classroomService = cs
	g.weekData.classroomService, _ = services.NewClassroomService(
		classrooms, g.weekData.busyGrid, float32(g.ClassroomOccupancy))
	return nil
}

func (g *ScheduleGenerator) SubmitAndGoToTheNextStep() (GeneratorStep, error) {
	if !g.canGoToTheNextStep {
		return g.currentStep, fmt.Errorf("can't go to the next step, %q unfinished", g.currentStep)
	}

	switch g.currentStep {
	case Setup:
		g.currentStep = DayBlocking
	case DayBlocking:
		g.currentStep = BoneLessonBuilding
	case BoneLessonBuilding:
		g.currentStep = ToBoneLessonsClassroomAssigning
	case ToBoneLessonsClassroomAssigning:
		g.currentStep = LessonSkeletonBuilding
	case LessonSkeletonBuilding:
		g.currentStep = FloatingLessonAdding
	case FloatingLessonAdding:
		g.currentStep = ToFloatingLessonsClassroomAssigning
	case ToFloatingLessonsClassroomAssigning:
		g.currentStep = Extraction
	default:
		return g.currentStep, fmt.Errorf("tou are at the last step")
	}

	g.canGoToTheNextStep = false
	return g.currentStep, nil
}

func (g *ScheduleGenerator) SetDaysForLessonTypes() (responses.DaysForLessonTypes, error) {
	if g.currentStep != DayBlocking {
		return responses.DaysForLessonTypes{},
			fmt.Errorf("invalid method: current step is %s instead of %s", g.currentStep, DayBlocking)
	}

	dayBlocker := components.NewDayBlocker(
		g.weekData.studentGroupService.GetAll(),
		components.NewErrorService[responses.LessonTypeDayDebt, *components.SetDayTypeError](),
		g.fullData.studentGroupService.GetAll()[0].GetFullWeekCount(), g.LessonFillRate,
	)
	dayBlocker.SetDayTypes()

	g.canGoToTheNextStep = true

	res := g.formDaysForLessonTypes()
	res.Errors = dayBlocker.GetErrorService().GetGeneratorResponseErrors()
	return res, nil
}

func (g *ScheduleGenerator) formDaysForLessonTypes() responses.DaysForLessonTypes {
	studentGroups := g.weekData.studentGroupService.GetAll()
	result := responses.DaysForLessonTypes{
		StudentGroups: make([]responses.StudentGroupWithLTypeDays, 0, len(studentGroups)),
	}

	for _, studentGroup := range studentGroups {
		resultSG := responses.StudentGroupWithLTypeDays{
			CommonEntity: responses.CommonEntity{
				ID:   studentGroup.ID,
				Name: studentGroup.Name,
			},
			WeekdayLessonTypes: make([]responses.LessonTypeWeekdayBinding, 0),
		}

		for j := range 7 {
			lessonType := studentGroup.GetTypeOfDay(j)
			if lessonType != nil {
				resultSG.WeekdayLessonTypes = append(resultSG.WeekdayLessonTypes,
					responses.LessonTypeWeekdayBinding{
						CommonEntity: responses.CommonEntity{
							ID:   lessonType.ID,
							Name: lessonType.Name,
						},
						Weekday: j,
					})
			}
		}

		result.StudentGroups = append(result.StudentGroups, resultSG)
	}

	return result
}

func (g *ScheduleGenerator) GenerateBoneLessons() (responses.BoneLessons, error) {
	if g.currentStep != BoneLessonBuilding {
		return responses.BoneLessons{},
			fmt.Errorf("invalid method: current step is %s instead of %s", g.currentStep, BoneLessonBuilding)
	}

	boneGenerator := components.NewBoneGenerator(
		components.NewErrorService[responses.UnassignedLesson, *components.BoneWeekError](),
		g.weekData.studyLoadService.GetAll(),
		g.weekData.lessonService,
	)
	boneGenerator.GenerateBoneLessons()

	g.canGoToTheNextStep = true

	res := responses.BoneLessons{
		Lessons: g.formBoneLessons(),
	}
	res.Errors = boneGenerator.GetErrorService().GetGeneratorResponseErrors()
	return res, nil
}

func (g *ScheduleGenerator) formBoneLessons() []responses.BoneLesson {
	lessons := g.weekData.lessonService.GetAll()
	boneLessons := make([]responses.BoneLesson, 0, len(lessons))
	for _, lesson := range lessons {
		boneLessons = append(boneLessons, responses.BoneLesson{
			CommonLesson: responses.CommonLesson{
				Teacher: responses.CommonEntity{
					ID:   lesson.Teacher.ID,
					Name: lesson.Teacher.UserName,
				},
				StudentGroup: responses.CommonEntity{
					ID:   lesson.StudentGroup.ID,
					Name: lesson.StudentGroup.Name,
				},
				Discipline: responses.CommonEntity{
					ID:   lesson.Discipline.ID,
					Name: lesson.Discipline.Name,
				},
				LessonType: responses.CommonEntity{
					ID:   lesson.Type.ID,
					Name: lesson.Type.Name,
				},
			},
			Day:  lesson.Day,
			Slot: lesson.Slot,
			Classroom: func() *responses.CommonEntity {
				if lesson.Classroom == nil {
					return nil
				}
				return &responses.CommonEntity{
					ID:   lesson.Classroom.ID,
					Name: lesson.Classroom.RoomNumber,
				}
			}(),
		})
	}

	return boneLessons
}

func (g *ScheduleGenerator) AssignClassroomsToBoneLessons() (responses.BoneLessonsWithC, error) {
	if g.currentStep != ToBoneLessonsClassroomAssigning {
		return responses.BoneLessonsWithC{},
			fmt.Errorf("invalid method: current step is %s instead of %s", g.currentStep, ToBoneLessonsClassroomAssigning)
	}

	classroomAssigner := components.NewClassroomAssigner(
		g.weekData.classroomService.GetAll(),
		g.weekData.lessonService.GetAll(),
		components.NewErrorService[responses.LessonWithoutClassroom, *components.ClassroomAssignError](),
	)
	if err := classroomAssigner.CheckAvailability(); err != nil {
		return responses.BoneLessonsWithC{}, fmt.Errorf("can't assign classrooms: %w", err)
	}
	classroomAssigner.AssignClassrooms()

	g.canGoToTheNextStep = true

	res := responses.BoneLessonsWithC{LessonsWithClassroom: g.formBoneLessons()}
	res.LessonsWithoutClassroom = classroomAssigner.GetErrorService().GetGeneratorResponseErrors()
	return res, nil
}

func (g *ScheduleGenerator) BuildScheduleSkeleton() (responses.GeneratedLessons, error) {
	if g.currentStep != LessonSkeletonBuilding {
		return responses.GeneratedLessons{},
			fmt.Errorf("invalid method: current step is %s instead of %s", g.currentStep, LessonSkeletonBuilding)
	}

	lessons := g.weekData.lessonService.GetAll()
	for _, lesson := range lessons {
		teacher := g.fullData.teacherService.Find(lesson.Teacher.ID)
		studentGroup := g.fullData.studentGroupService.Find(lesson.StudentGroup.ID)

		// binding lesson type to day for actual student groups
		for weekday := range 7 {
			weekLT := lesson.StudentGroup.GetTypeOfDay(weekday)
			if weekLT != nil {
				lt := studentGroup.GetTypeOfDay(weekday)
				if lt == nil {
					lt := g.fullData.lessonTypeService.Find(weekLT.ID)
					studentGroup.BindWeekday(lt, weekday)
				}
			}
		}

		discipline := g.fullData.disciplineService.Find(lesson.Discipline.ID)
		lessonType := g.fullData.lessonTypeService.Find(lesson.Type.ID)
		studyLoad := g.fullData.studyLoadService.Find(*entities.NewUnassignedLesson(
			lessonType, teacher, studentGroup, discipline,
		))
		classroom := func(weekC *entities.Classroom) *entities.Classroom {
			if weekC == nil {
				return nil
			}
			return g.fullData.classroomService.Find(weekC.ID)
		}(lesson.Classroom)

		// copy week lesson to all weeks
		currentWeek := 0
		outOfGrid := false
		for !outOfGrid {
			err := g.fullData.lessonService.AssignLesson(studyLoad,
				entities.NewLessonSlot(lesson.Day+currentWeek*7, lesson.Slot),
			)
			if err == nil && classroom != nil {
				fullL := g.fullData.lessonService.Select().Sort().Last()
				fullL.SetClassroom(classroom)
			}

			var dayErr *entities.DayOutError
			if errors.As(err, &dayErr) {
				outOfGrid = true
			}
			currentWeek++
		}
	}

	g.canGoToTheNextStep = true

	res := responses.GeneratedLessons{Lessons: g.formLessons(g.fullData.lessonService.GetAll())}
	return res, nil
}

func (g *ScheduleGenerator) formLessons(lessons []*entities.Lesson) []responses.GeneratedLesson {
	type key struct {
		TeacherID    uuid.UUID
		GroupID      uuid.UUID
		DisciplineID uuid.UUID
		LessonTypeID uuid.UUID
		Weekday      int
		Slot         int
		ClassroomID  *uuid.UUID
	}

	grouped := make(map[key]*responses.GeneratedLesson)

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

		grouped[k] = &responses.GeneratedLesson{
			CommonLesson: responses.CommonLesson{
				Teacher: responses.CommonEntity{
					ID:   lesson.Teacher.ID,
					Name: lesson.Teacher.UserName,
				},
				StudentGroup: responses.CommonEntity{
					ID:   lesson.StudentGroup.ID,
					Name: lesson.StudentGroup.Name,
				},
				Discipline: responses.CommonEntity{
					ID:   lesson.Discipline.ID,
					Name: lesson.Discipline.Name,
				},
				LessonType: responses.CommonEntity{
					ID:   lesson.Type.ID,
					Name: lesson.Type.Name,
				},
			},
			Days: []int{lesson.Day},
			Slot: lesson.Slot,
			Classroom: func() *responses.CommonEntity {
				if lesson.Classroom == nil {
					return nil
				}
				return &responses.CommonEntity{
					ID:   lesson.Classroom.ID,
					Name: lesson.Classroom.RoomNumber,
				}
			}(),
		}
	}

	generatedLessons := make([]responses.GeneratedLesson, 0, len(grouped))
	for _, lesson := range grouped {
		generatedLessons = append(generatedLessons, *lesson)
	}

	return generatedLessons
}

func (g *ScheduleGenerator) AddFloatingLessons() (responses.GeneratedLessons, error) {
	if g.currentStep != FloatingLessonAdding {
		return responses.GeneratedLessons{},
			fmt.Errorf("invalid method: current step is %s instead of %s", g.currentStep, FloatingLessonAdding)
	}

	g.floatingLessonService, _ = services.NewLessonService()

	missingLessonAssigner := components.NewMissingLessonAdder(
		components.NewErrorService[responses.UnassignedLesson, *components.MissingLessonsAdderError](),
		g.fullData.studyLoadService.GetAll(),
		g.floatingLessonService)
	missingLessonAssigner.AddMissingLessons()

	g.canGoToTheNextStep = true

	res := responses.GeneratedLessons{Lessons: g.formLessons(g.floatingLessonService.GetAll())}
	res.Errors = missingLessonAssigner.GetErrorService().GetGeneratorResponseErrors()
	return res, nil
}

func (g *ScheduleGenerator) AssignClassroomsToFloatingLessons() (responses.GeneratedLessonsWithC, error) {
	if g.currentStep != ToFloatingLessonsClassroomAssigning {
		return responses.GeneratedLessonsWithC{},
			fmt.Errorf("invalid method: current step is %s instead of %s", g.currentStep, ToFloatingLessonsClassroomAssigning)
	}

	classroomAssigner := components.NewClassroomAssigner(
		g.fullData.classroomService.GetAll(),
		g.floatingLessonService.GetAll(),
		components.NewErrorService[responses.LessonWithoutClassroom, *components.ClassroomAssignError](),
	)
	if err := classroomAssigner.CheckAvailability(); err != nil {
		return responses.GeneratedLessonsWithC{}, fmt.Errorf("can't assign classrooms: %w", err)
	}
	classroomAssigner.AssignClassrooms()

	g.canGoToTheNextStep = true

	res := responses.GeneratedLessonsWithC{LessonsWithClassroom: g.formLessons(g.floatingLessonService.GetAll())}
	res.LessonsWithoutClassroom = classroomAssigner.GetErrorService().GetGeneratorResponseErrors()
	return res, nil
}

func (g *ScheduleGenerator) ExtractStudyLoads() ([]responses.StudyLoad, error) {
	if g.currentStep != Extraction {
		return nil, fmt.Errorf("invalid method: current step is %s instead of %s", g.currentStep, Extraction)
	}

	studyLoads := g.fullData.studyLoadService.GetAll()
	result := make([]responses.StudyLoad, 0, len(studyLoads))

	for _, studyLoad := range studyLoads {
		result = append(result, responses.StudyLoad{
			ID:             studyLoad.ID,
			TeacherID:      studyLoad.Teacher.ID,
			StudentGroupID: studyLoad.StudentGroup.ID,
			DisciplineID:   studyLoad.Discipline.ID,
			LessonTypeID:   studyLoad.Type.ID,
		})
	}

	return result, nil
}

func (g *ScheduleGenerator) ExtractLessons() ([]responses.Lesson, error) {
	if g.currentStep != Extraction {
		return nil, fmt.Errorf("invalid method: current step is %s instead of %s", g.currentStep, Extraction)
	}

	boneLessons := g.fullData.lessonService.GetAll()
	floatingLessons := g.floatingLessonService.GetAll()
	allLessons := append(boneLessons, floatingLessons...)
	result := make([]responses.Lesson, 0, len(allLessons))

	for _, lesson := range allLessons {
		result = append(result, responses.Lesson{
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

	return result, nil
}

// main function
func (g *ScheduleGenerator) GenerateSchedule() error {
	if g.fullData.studyLoadService == nil {
		return fmt.Errorf("study loads not set")
	}
	if g.weekData.studyLoadService == nil {
		return fmt.Errorf("study loads not set")
	}

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

	// if !g.errorService.IsClear() {
	// 	return g.errorService
	// }
	return nil
}

// Rates schedule fault. Returns empty if the generator is at the wrong step.
func (g *ScheduleGenerator) GetFault() responses.Fault {
	// if g.currentStep != Extraction {
	// 	return responses.Fault{}
	// }

	scheduleFault := components.NewScheduleFault()

	scheduleFault.AddParameter("teacher_windows", components.NewSimpleScheduleParameter(
		float64(g.fullData.teacherService.CountWindows()), 0.1,
	))
	scheduleFault.AddParameter("student_group_windows", components.NewSimpleScheduleParameter(
		float64(g.fullData.studentGroupService.CountWindows()), 1000,
	))
	scheduleFault.AddParameter("study_load_hours_deficit", components.NewSimpleScheduleParameter(
		float64(g.fullData.studyLoadService.CountHoursDeficit()), 10,
	))
	scheduleFault.AddParameter("teacher_lesson_overlapping", components.NewSimpleScheduleParameter(
		float64(g.fullData.teacherService.CountLessonOverlapping()), 1000,
	))
	scheduleFault.AddParameter("student_group_lesson_overlapping", components.NewSimpleScheduleParameter(
		float64(g.fullData.studentGroupService.CountLessonOverlapping()), 1000,
	))
	scheduleFault.AddParameter("classroom_lesson_overlapping", components.NewSimpleScheduleParameter(
		float64(g.fullData.classroomService.CountLessonOverlapping()), 1000,
	))
	scheduleFault.AddParameter("student_group_overtime_lessons", components.NewSimpleScheduleParameter(
		float64(g.fullData.studentGroupService.CountOvertimeLessons()), 10,
	))
	scheduleFault.AddParameter("student_group_invalid_lessons_by_type", components.NewSimpleScheduleParameter(
		float64(g.fullData.studentGroupService.CountInvalidLessonsByType()), 10,
	))
	scheduleFault.AddParameter("lessons_without_classroom", components.NewSimpleScheduleParameter(
		float64(g.fullData.lessonService.CountLessonsWithoutClassroom(g.fullData.lessonService.GetAll())), 10,
	))
	scheduleFault.AddParameter("classroom_with_overflow", components.NewSimpleScheduleParameter(
		float64(g.fullData.classroomService.CountOverflowLessons()), 10,
	))

	g.WriteSchedule()

	return scheduleFault.GetResponse()
}

func (g *ScheduleGenerator) WriteSchedule() {
	tSchedule := make(map[*entities.Teacher]*entities.PersonalSchedule, len(g.fullData.teacherService.GetAll()))
	for i := range g.fullData.teacherService.GetAll() {
		t := g.fullData.teacherService.GetAll()[i]
		tSchedule[t] = &entities.PersonalSchedule{
			BusyGrid: &t.BusyGrid,
			Out:      "schedule-output/" + t.UserName + ".txt",
		}
	}

	sgSchedule := make(map[*entities.StudentGroup]*entities.PersonalSchedule, len(g.fullData.studentGroupService.GetAll()))
	for i := range g.fullData.studentGroupService.GetAll() {
		sg := g.fullData.studentGroupService.GetAll()[i]
		sgSchedule[sg] = &entities.PersonalSchedule{
			BusyGrid: &sg.BusyGrid,
			Out:      "schedule-output/" + sg.Name + ".txt",
		}
	}

	lessons := append(g.fullData.lessonService.GetAll(), g.floatingLessonService.GetAll()...)
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
