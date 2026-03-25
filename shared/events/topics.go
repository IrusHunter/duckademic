package events

type RedisTopic string

const (
	AcademicRankRT         RedisTopic = "academic:ranks"
	TeacherRT              RedisTopic = "teachers"
	StudentRT              RedisTopic = "students"
	LessonTypeRT           RedisTopic = "lesson:types"
	DisciplineRT           RedisTopic = "disciplines"
	LessonTypeAssignmentRT RedisTopic = "lesson:type:assignments"
)

type EventType string

const (
	EntityCreated EventType = "ENTITY_CREATED"
	EntityUpdated EventType = "ENTITY_UPDATED"
	EntityDeleted EventType = "ENTITY_DELETED"
)
