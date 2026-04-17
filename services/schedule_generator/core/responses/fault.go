package responses

type Fault struct {
	TotalValue float64 `json:"total_value"`

	TeacherWindows                FaultParam[TeacherWindow]              `json:"teacher_windows"`
	StudentGroupWindows           FaultParam[StudentGroupWindow]         `json:"student_group_windows"`
	StudyLoadHoursDeficit         FaultParam[StudyLoadHoursDeficit]      `json:"study_load_hours_deficit"`
	TeacherLessonOverlapping      FaultParam[TeacherLessonOverlap]       `json:"teacher_lesson_overlapping"`
	StudentGroupLessonOverlapping FaultParam[StudentGroupLessonOverlap]  `json:"student_group_lesson_overlapping"`
	ClassroomLessonOverlapping    FaultParam[ClassroomLessonOverlap]     `json:"classroom_lesson_overlapping"`
	StudentGroupOvertimeLessons   FaultParam[StudentGroupOvertimeLesson] `json:"student_group_overtime_lessons"`
	StudentGroupInvalidLessons    FaultParam[StudentGroupInvalidLesson]  `json:"student_group_invalid_lessons_by_type"`
	LessonsWithoutClassroom       FaultParam[LessonWithoutClassroomF]    `json:"lessons_without_classroom"`
	ClassroomWithOverflow         FaultParam[ClassroomOverflow]          `json:"classroom_with_overflow"`
}

type FaultParam[T any] struct {
	Value      float64 `json:"value"`
	Defections []T     `json:"defections"`
}

type TeacherWindow struct {
	TeacherID int `json:"teacher_id"`
	Day       int `json:"day"`
	StartSlot int `json:"start_slot"`
	EndSlot   int `json:"end_slot"`
}

type TeacherLessonOverlap struct {
	TeacherID int   `json:"teacher_id"`
	Lessons   []int `json:"lessons"`
	Day       int   `json:"day"`
	Slot      int   `json:"slot"`
}

type StudentGroupWindow struct {
	GroupID   int `json:"group_id"`
	Day       int `json:"day"`
	StartSlot int `json:"start_slot"`
	EndSlot   int `json:"end_slot"`
}

type StudentGroupLessonOverlap struct {
	GroupID int   `json:"group_id"`
	Lessons []int `json:"lessons"`
	Day     int   `json:"day"`
	Slot    int   `json:"slot"`
}

type StudentGroupOvertimeLesson struct {
	GroupID  int `json:"group_id"`
	LessonID int `json:"lesson_id"`
	Day      int `json:"day"`
	Slot     int `json:"slot"`
}

type StudentGroupInvalidLesson struct {
	GroupID  int    `json:"group_id"`
	LessonID int    `json:"lesson_id"`
	Type     string `json:"type"`
}

type StudyLoadHoursDeficit struct {
	GroupID int     `json:"group_id"`
	Missing float64 `json:"missing_hours"`
}

type ClassroomLessonOverlap struct {
	ClassroomID int   `json:"classroom_id"`
	Lessons     []int `json:"lessons"`
	Day         int   `json:"day"`
	Slot        int   `json:"slot"`
}

type ClassroomOverflow struct {
	ClassroomID int `json:"classroom_id"`
	LessonID    int `json:"lesson_id"`
	Capacity    int `json:"capacity"`
	Required    int `json:"required"`
}

type LessonWithoutClassroomF struct {
	LessonID int `json:"lesson_id"`
	Day      int `json:"day"`
	Slot     int `json:"slot"`
}
