package events

type RedisTopic string

const (
	AcademicRankRT         RedisTopic = "academic:ranks"
	TeacherRT              RedisTopic = "teachers"
	StudentRT              RedisTopic = "students" // 1 delay
	LessonTypeRT           RedisTopic = "lesson:types"
	DisciplineRT           RedisTopic = "disciplines"
	LessonTypeAssignmentRT RedisTopic = "lesson:type:assignments"
	SemesterRT             RedisTopic = "semesters"
)

type EventType string

const (
	EntityCreated EventType = "ENTITY_CREATED"
	EntityUpdated EventType = "ENTITY_UPDATED"
	EntityDeleted EventType = "ENTITY_DELETED"
)
