package events

type RedisTopic string

const (
	AcademicRankRT          RedisTopic = "academic:ranks"
	TeacherRT               RedisTopic = "teachers" // x2 delay
	StudentRT               RedisTopic = "students" // x2 delay
	LessonTypeRT            RedisTopic = "lesson:types"
	DisciplineRT            RedisTopic = "disciplines"
	SemesterDisciplineRT    RedisTopic = "semester:disciplines"
	LessonTypeAssignmentRT  RedisTopic = "lesson:type:assignments"
	SemesterRT              RedisTopic = "semesters"
	StudentGroupRT          RedisTopic = "student:groups"          // x2 delay
	GroupMemberRT           RedisTopic = "group:members"           // x3 delay
	GroupCohortRT           RedisTopic = "group:cohorts"           // x2 delay
	TeacherLoadRT           RedisTopic = "teacher:loads"           // x3 delay
	GroupCohortAssignmentRT RedisTopic = "group:cohort:assignment" // x2 delay
	ClassroomRT             RedisTopic = "classroom"
	AccessPermissionRT      RedisTopic = "access:permission"
)

type EventType string

const (
	EntityCreated EventType = "ENTITY_CREATED"
	EntityUpdated EventType = "ENTITY_UPDATED"
	EntityDeleted EventType = "ENTITY_DELETED"
)
