package events

type RedisTopic string

const (
	AcademicRankRT          RedisTopic = "academic:ranks"
	TeacherRT               RedisTopic = "teachers" // 1 delay
	StudentRT               RedisTopic = "students" // 1 delay
	LessonTypeRT            RedisTopic = "lesson:types"
	DisciplineRT            RedisTopic = "disciplines"
	LessonTypeAssignmentRT  RedisTopic = "lesson:type:assignments"
	SemesterRT              RedisTopic = "semesters"
	StudentGroupRT          RedisTopic = "student:groups"          // 1 delay
	GroupMemberRT           RedisTopic = "group:members"           // 2 delay
	GroupCohortRT           RedisTopic = "group:cohorts"           // 1 delay
	TeacherLoadRT           RedisTopic = "teacher:loads"           // 2 delay
	GroupCohortAssignmentRT RedisTopic = "group:cohort:assignment" // 1 delay
	ClassroomRT             RedisTopic = "classroom"
	AccessPermissionRT      RedisTopic = "access:permission"
)

type EventType string

const (
	EntityCreated EventType = "ENTITY_CREATED"
	EntityUpdated EventType = "ENTITY_UPDATED"
	EntityDeleted EventType = "ENTITY_DELETED"
)
