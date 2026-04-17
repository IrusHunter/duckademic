package components

import (
	"fmt"
	"strings"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
)

// ScheduleParameter defines a single component of the ScheduleFault.
// Each implementation provides its own fault calculation logic.
type ScheduleParameter interface {
	Fault() float64       // Calculates fault value based on the argument's internal calculation formula.
	GetArguments() string // Returns formatted, human-readable representation of the arguments
}

// NewSimpleScheduleParameter a new simpleScheduleParameter instance.
// The resulting fault value is computed as value (v) multiplied by its
// fault weight factor (f).
func NewSimpleScheduleParameter(v, f float64) ScheduleParameter {
	return &simpleScheduleParameter{value: v, fault: f}
}

type simpleScheduleParameter struct {
	fault float64
	value float64
}

func (ssp *simpleScheduleParameter) Fault() float64 {
	return ssp.fault * ssp.value
}
func (ssp *simpleScheduleParameter) GetArguments() string {
	return fmt.Sprintf("value: %f, fault %f", ssp.value, ssp.fault)
}

// ScheduleFault represents fault for generated schedule.
type ScheduleFault interface {
	AddParameter(name string, parameter ScheduleParameter) // Adds new ScheduleParameter or replace it with a new one.
	Fault() float64                                        // Returns sum of the ScheduleParameter Faults.
	GetParameters() string                                 // Returns formatted, human-readable representation of the parameters.
	GetResponse() responses.Fault
}

// NewScheduleFault creates a new ScheduleFault instance.
func NewScheduleFault() ScheduleFault {
	return &scheduleFault{parameters: map[string]ScheduleParameter{}}
}

type scheduleFault struct {
	parameters map[string]ScheduleParameter
}

func (sf *scheduleFault) AddParameter(name string, parameter ScheduleParameter) {
	sf.parameters[name] = parameter
}
func (sf *scheduleFault) Fault() (res float64) {
	for _, fault := range sf.parameters {
		res += fault.Fault()
	}
	return
}
func (sf *scheduleFault) GetParameters() string {
	if len(sf.parameters) == 0 {
		return ""
	}

	var b strings.Builder
	for key, value := range sf.parameters {
		if value.Fault() != 0 {
			fmt.Fprintf(&b, "%s: %s, fault: %f \n", key, value.GetArguments(), value.Fault())
		}
	}
	s := b.String()
	return s
}
func (sf *scheduleFault) GetResponse() responses.Fault {
	res := responses.Fault{
		TotalValue: sf.Fault(),
	}

	res.TeacherWindows = responses.FaultParam[responses.TeacherWindow]{
		Value: sf.parameters["teacher_windows"].Fault(),
	}
	res.StudentGroupWindows = responses.FaultParam[responses.StudentGroupWindow]{
		Value: sf.parameters["student_group_windows"].Fault(),
	}
	res.StudyLoadHoursDeficit = responses.FaultParam[responses.StudyLoadHoursDeficit]{
		Value: sf.parameters["study_load_hours_deficit"].Fault(),
	}
	res.TeacherLessonOverlapping = responses.FaultParam[responses.TeacherLessonOverlap]{
		Value: sf.parameters["teacher_lesson_overlapping"].Fault(),
	}

	res.StudentGroupLessonOverlapping = responses.FaultParam[responses.StudentGroupLessonOverlap]{
		Value: sf.parameters["student_group_lesson_overlapping"].Fault(),
	}

	res.ClassroomLessonOverlapping = responses.FaultParam[responses.ClassroomLessonOverlap]{
		Value: sf.parameters["classroom_lesson_overlapping"].Fault(),
	}

	res.StudentGroupOvertimeLessons = responses.FaultParam[responses.StudentGroupOvertimeLesson]{
		Value: sf.parameters["student_group_overtime_lessons"].Fault(),
	}

	res.StudentGroupInvalidLessons = responses.FaultParam[responses.StudentGroupInvalidLesson]{
		Value: sf.parameters["student_group_invalid_lessons_by_type"].Fault(),
	}

	res.LessonsWithoutClassroom = responses.FaultParam[responses.LessonWithoutClassroomF]{
		Value: sf.parameters["lessons_without_classroom"].Fault(),
	}

	res.ClassroomWithOverflow = responses.FaultParam[responses.ClassroomOverflow]{
		Value: sf.parameters["classroom_with_overflow"].Fault(),
	}

	return res
}
