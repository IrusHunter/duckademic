# API Documentation

# ApiGateway

## /seed

### ANY - seeds databases of all services.

204 NO CONTENT

500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

## /

### ANY - routes and proxies incoming requests to the appropriate backend service based on the request path.

- Employee Service (/employee)
  - Source of Truth
    - [/academic-ranks](#employee-academic-ranks) !employee.academic_rank
    - [/academic-rank/{id}](#employee-academic-rank-id) !employee.academic_rank
    - [/academic-degrees](#employee-academic-degrees) !employee.academic_degree
    - [/academic-degree/{id}](#employee-academic-degree-id) !employee.academic_degree
    - [/employees](#employee-employees) !employee.employee
    - [/employee/{id}](#employee-employee-id) !employee.employee
    - [/teachers](#employee-teachers) !employee.teacher
    - [/teacher/{id}](#employee-teacher-id) !employee.teacher

- Schedule Service (/schedule)
  - Enriched View
    - [/academic-ranks](#schedule-academic-ranks) !schedule.academic_rank
    - [/academic-rank/{id}](#schedule-academic-rank-id) !schedule.academic_rank
    - [/lesson-types](#schedule-lesson-types) !schedule.lesson_type
    - [/lesson-type/{id}](#schedule-lesson-type-id) !schedule.lesson_type
  - Mirror
    - [/teachers](#schedule-teachers) !schedule.teacher
    - [/disciplines](#schedule-disciplines) !schedule.discipline
    - [/lesson-type-assignments](#schedule-lesson-type-assignments) !schedule.lesson_type_assignment
    - [/students](#schedule-students) !schedule.student
    - [/student-groups](#schedule-student-groups) !schedule.student_group
    - [/group_members](#schedule-group-members) !schedule.group_member
    - [/teacher-loads](#schedule-teacher-loads) !schedule.teacher_load
    - [/group-cohorts](#schedule-group-cohorts) !schedule.group_cohort
    - [/group-cohort-assignments](#schedule-group-cohort-assignments) !schedule.group_cohort_assignment
    - [/classrooms](#schedule-classrooms) !schedule.classroom
    - [/semesters](#schedule-semesters) !schedule.semester
    - [/semester-disciplines](#schedule-semester-disciplines) !schedule.semester_discipline
  - Schedule Generator Integration
    - [/load-data-into-generator](#schedule-load-data-into-generator)
    - [/load-classrooms-into-generator](#schedule-load-classrooms-into-generator)
    - [/extract-data-from-generator](#schedule-extract-data-from-generator)
    - [/study-loads](#schedule-study-loads) !schedule.study_load
    - [/lesson-slots](#schedule-lesson-slots) !schedule.lesson_slot
    - [/lesson-occurrences](#schedule-lesson-occurrences) !schedule.lesson_occurrence
    - [/get-personal-schedule](#schedule-get-personal-schedule) !

- Schedule Generator Service (/schedule-generator)
  - Source of Truth
    - [/default-generator-config](#schedule-generator-default-generator-config)
  - Generation Pipeline
    - [/init](#schedule-generator-init)
    - [/submit-and-go](#schedule-generator-submit-and-go)
    - [/generate-days-for-lesson-types](#schedule-generator-generate-days-for-lesson-types)
    - [/generate-bone-lessons](#schedule-generator-generate-bone-lessons)
    - [/assign-classrooms-to-bone-lessons](#schedule-generator-assign-classrooms-to-bone-lessons)
    - [/build-schedule-skeleton](#schedule-generator-build-schedule-skeleton)
    - [/add-floating-lessons](#schedule-generator-add-floating-lessons)
    - [/assign-classrooms-to-floating-lessons](#schedule-generator-assign-classrooms-to-floating-lessons)
    - [/get-study-loads](#schedule-generator-get-study-loads)
    - [/get-lessons](#schedule-generator-get-lessons)
    - [/get-fault](#schedule-generator-get-fault)
  - Internal
    - [/set-teachers](#schedule-generator-set-teachers)
    - [/set-disciplines](#schedule-generator-set-disciplines)
    - [/set-lesson-types](#schedule-generator-set-lesson-types)
    - [/set-lesson-type-assignments](#schedule-generator-set-lesson-type-assignments)
    - [/set-student-groups](#schedule-generator-set-student-groups)
    - [/set-teacher-loads](#schedule-generator-set-teacher-loads)
    - [/set-classrooms](#schedule-generator-set-classrooms)

- Student Service (/student)
  - Source of Truth
    - [/students](#student-students) !student.student
    - [/student/{id}](#student-student-id) !student.student
  - Mirror
    - [/semesters](#student-semesters) !student.semester

- Student Group Service (/student-group)
  - Source of Truth
    - [/group-cohorts](#student-group-group-cohorts) !student_group.group_cohort
    - [/group-cohort/{id}](#student-group-group-cohort-id) !student_group.group_cohort
    - [/student-groups](#student-group-student-groups) !student_group.student_group
    - [/student-group/{id}](#student-group-student-group-id) !student_group.student_group
    - [/group-members](#student-group-group-members) !student_group.group_member
    - [/group-member/{id}](#student-group-group-member-id) !student_group.group_member
    - [/group-cohort-assignments](#student-group-group-cohort-assignments) !student_group.group_cohort_assignment
    - [/group-cohort-assignment/{id}](#student-group-group-cohort-assignment) !student_group.group_cohort_assignment
  - Mirror
    - [/semesters](#student-group-semesters) !student_group.semester
    - [/students](#student-group-students) !student_group.student
    - [/lesson-types](#student-group-lesson-types) !student_group.lesson_type
    - [/disciplines](#student-group-disciplines) !student_group.discipline

- Curriculum Service (/curriculum)
  - Source of Truth
    - [/curriculums](#curriculum-curriculum) !curriculum.curriculum
    - [/curriculum/{id}](#curriculum-curriculum-id) !curriculum.curriculum
    - [/semesters](#curriculum-semesters) !curriculum.semester
    - [/semester/{id}](#curriculum-semester-id) !curriculum.semester
    - [/lesson-types](#curriculum-lesson-types) !curriculum.lesson_type
    - [/lesson-type/{id}](#curriculum-lesson-type-id) !curriculum.lesson_type
    - [/disciplines](#curriculum-disciplines) !curriculum.discipline
    - [/discipline/{id}](#curriculum-discipline-id) !curriculum.discipline
    - [/lesson-type-assignments](#curriculum-lesson-type-assignments) !curriculum.lesson_type_assignment
    - [/lesson-type-assignment/{id}](#curriculum-lesson-type-assignment-id) !curriculum.lesson_type_assignment
    - [/semester-disciplines](#curriculum-semester-disciplines) !curriculum.semester_discipline
    - [/semester-discipline/{id}](#curriculum-semester-discipline-id) !curriculum.semester_discipline

- Teacher Load Service (/teacher-load)
  - Source of Truth
    - [/teacher-loads](#teacher-load-teacher-loads) !teacher_load.teacher_load
    - [/teacher-load/{id}](#teacher-load-teacher-load-id) !teacher_load.teacher_load
  - Mirror
    - [/teachers](#teacher-load-teachers) !teacher_load.teacher_load
    - [/lesson-types](#teacher-load-lesson-types) !teacher_load.teacher_load
    - [/disciplines](#teacher-load-disciplines) !teacher_load.teacher_load

- Asset Service (/asset)
  - Source of Truth
    - [/classrooms](#employee-classrooms) !asset.classroom
    - [/classroom/{id}](#employee-classroom) !asset.classroom

- Auth Service (/auth)
  - Source of Truth
    - [/permissions](#auth-permissions) !auth.permission
    - [/permission/{id}](#auth-permission-id) !auth.permission
    - [/roles](#auth-roles) !auth.role
    - [/role/{id}](#auth-role-id) !auth.role
    - [/role-permissions](#auth-role-permissions) !auth.role_permission
    - [/role-permission/{id}](#auth-role-permission-id) !auth.role_permission
    - [/services](#auth-services) !auth.service
    - [/service/{id}](#auth-service-id) !auth.service
    - [/service-permissions](#auth-service-permissions) !auth.service_permission
    - [/service-permission/{id}](#auth-service-permission-id) !auth.service_permission
  - Projection
    - [/users](#auth-users) !auth.user
    - [/user/{id}](#auth-user-id) !auth.user
  - Operations
    - [/login](#auth-login)
    - [/refresh](#auth-refresh)
    - [/reset-password/{id}](#auth-reset-password-id) !auth.user.reset_password
    - [/change-password](#auth-change-password)
- Course Service (/course)
  - Mirror
    - [/teachers](#course-teachers) !course.teacher
    - [/students](#course-students) !course.student
  - Projection
    - [/courses](#course-courses) !course.course
    - [/course/{id}](#course-course-id) !course.course
  - Source of Truth
    - [/student-courses](#course-student-courses) !course.student_course
    - [/student-course/{id}](#course-user-id) !course.student_course
    - [/teacher-courses](#course-teacher-courses) !course.teacher_course
    - [/teacher-course/{id}](#course-teacher-course-id) !course.teacher_course
    - [/tasks](#course-tasks) !course.task
    - [/task/{id}](#course-task-id) !course.task
    - [/task-students](#course-task-students) !course.task-student
    - [/task-student/{id}](#course-task-student-id) !course.task-student

400 BAD REQUEST or 500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

# Employee Service

<a id="employee-academic-ranks"></a>

## /academic-ranks

### GET (employee.academic_rank) - gets all academic ranks from the database

200 OK [=> AcademicRank[]](schemas.md#employee-academic-rank)

### POST (employee.academic_rank) - adds a new academic rank

```json
{
  "title": "string (human-readable name of the rank)"
}
```

200 OK [=> AcademicRank](schemas.md#employee-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-academic-rank-id"></a>

## /academic-rank/{id}

### GET (employee.academic_rank) - finds academic rank with an ID as an URL parameter

200 OK [=> AcademicRank](schemas.md#employee-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (employee.academic_rank) - deletes an academic rank by its ID provided in the URL path

200 OK [=> AcademicRank](schemas.md#employee-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (employee.academic_rank) - updates an academic rank by its ID with the data provided in the request body

```json
{
  "title": "string (human-readable name of the rank)"
}
```

200 OK [=> AcademicRank](schemas.md#employee-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-academic-degrees"></a>

## /academic-degrees

### GET (employee.academic_degree) - gets all academic degrees from the database

200 OK [=> AcademicDegree[]](schemas.md#employee-academic-degree)

### POST (employee.academic_degree) - adds a new academic degree

```json
{
  "title": "string (human-readable name of the degree)"
}
```

200 OK [=> AcademicDegree](schemas.md#employee-academic-degree)

400 BAD REQUEST\*\* [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-academic-degree-id"></a>

## /academic-degree/{id}

### GET (employee.academic_degree) - finds academic degree with an ID as a URL parameter

200 OK [=> AcademicDegree](schemas.md#employee-academic-degree)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (employee.academic_degree) - deletes an academic degree by its ID provided in the URL path

**200 OK** [=> AcademicDegree](schemas.md#employee-academic-degree)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (employee.academic_degree) - updates an academic degree by its ID with the data provided in the request body

```json
{
  "title": "string (human-readable name of the degree)"
}
```

200 OK [=> AcademicDegree](schemas.md#employee-academic-degree)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-employees"></a>

## /employees

### GET (employee.employee) - gets all employees from the database

**200 OK** [=> Employee[]](schemas.md#employee-employee)

### POST (employee.employee) - adds a new employee

```json
{
  "first_name": "string (employee's first name)",
  "last_name": "string (employee's last name)",

  // Optional fields
  "middle_name": "string (employee's middle name)",
  "phone_number": "string (contact phone number)"
}
```

**200 OK** [=> Employee](schemas.md#employee-employee)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-employee-id"></a>

## /employee/{id}

### GET (employee.employee) - finds employee with an ID as a URL parameter

**200 OK** [=> Employee](schemas.md#employee-employee)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (employee.employee) - deletes an employee by its ID provided in the URL path

**200 OK** [=> Employee](schemas.md#employee-employee)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

### PUT (employee.employee) - updates an employee by its ID with the data provided in the request body

```json
{
  "first_name": "string (employee's first name)",
  "last_name": "string (employee's last name)",

  // Optional fields
  "middle_name": "string (employee's middle name)",
  "phone_number": "string (contact phone number)"
}
```

**200 OK** [=> Employee](schemas.md#employee-employee)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-teachers"></a>

## /teachers

### GET (employee.teacher) - gets all teachers from the database

**200 OK** [=> Teacher[]](schemas.md#employee-teacher)

### POST (employee.teacher) - adds a new teacher

```json
{
  "employee_id": "uuid (reference to employee)",
  "email": "string (teacher's email address)",
  "academic_rank_id": "uuid (reference to academic rank)",

  // Optional fields
  "academic_degree_id": "uuid (reference to academic degree)"
}
```

**200 OK** [=> Teacher](schemas.md#employee-teacher)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-teacher-id"></a>

## /teacher/{id}

### GET (employee.teacher) - finds teacher with an ID as a URL parameter

**200 OK** [=> Teacher](schemas.md#employee-teacher)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (employee.teacher) - deletes a teacher by its ID provided in the URL path

**200 OK** [=> Teacher](schemas.md#employee-teacher)

**400 BAD REQUEST** [-> ErrorResponse](schemas.md#errorresponse)

### PUT (employee.teacher) - updates a teacher by its ID with the data provided in the request body

```json
{
  "employee_id": "uuid (reference to employee)",
  "email": "string (teacher's email address)",
  "academic_rank_id": "uuid (reference to academic rank)",

  // Optional fields
  "academic_degree_id": "uuid (reference to academic degree)"
}
```

**200 OK** [=> Teacher](schemas.md#employee-teacher)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

# Schedule Service

<a id="schedule-academic-ranks"></a>

## /academic-ranks

### GET (schedule.academic_rank) - gets all academic ranks from the database

200 OK [=> AcademicRank[]](schemas.md#schedule-academic-rank)

<a id="schedule-academic-rank-id"></a>

## /academic-rank/{id}

### GET (schedule.academic_rank) - finds academic rank with an ID as an URL parameter

200 OK [=> AcademicRank](schemas.md#schedule-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (schedule.academic_rank) - updates an academic rank by its ID with the data provided in the request body

```json
{
  "priority": "int (determines the rank's priority: higher value = higher rank)"
}
```

200 OK [=> AcademicRank](schemas.md#schedule-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-teachers"></a>

## /teachers

### GET (schedule.teacher) - gets all teachers from the database

**200 OK** [=> Teacher[]](schemas.md#schedule-teacher)

<a id="schedule-lesson-types"></a>

## /lesson-types

### GET (schedule.lesson_type) - gets all lesson types from the database

200 OK [=> LessonType[]](schemas.md#schedule-lesson-type)

<a id="schedule-lesson-type-id"></a>

## /lesson-type/{id}

### GET (schedule.lesson_type) - finds lesson type by ID (provided as a URL parameter)

200 OK [=> LessonType](schemas.md#schedule-lesson-type)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (schedule.lesson_type) - updates a lesson type by its ID using the request body

```json
{
  "reserved_weeks": "string (comma-separated week numbers where only this lesson type is allowed)"
}
```

200 OK [=> LessonType](schemas.md#schedule-lesson-type)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-disciplines"></a>

## /disciplines

### GET (schedule.discipline) – gets all disciplines from the database

200 OK [=> Discipline[]](schemas.md#schedule-discipline)

<a id="schedule-lesson-type-assignments"></a>

## /lesson-type-assignments

### GET (schedule.lesson_type_assignment) – gets all lesson type assignments from the database

200 OK [=> LessonTypeAssignment[]](schemas.md#schedule-lesson-type-assignment)

<a id="schedule-students"></a>

## /students

### GET (schedule.student) - gets all students from the database

200 OK [=> Student[]](schemas.md#schedule-student)

<a id="schedule-student-groups"></a>

## /student-groups

### GET (schedule.student_group) - gets all student groups from the database

200 OK [=> StudentGroup[]](schemas.md#schedule-student-group)

<a id="schedule-group-members"></a>

## /group-members

### GET (schedule.group_member) - gets all group members from the database

200 OK [=> GroupMember[]](schemas.md#schedule-group-member)

<a id="schedule-teacher-loads"></a>

## /teacher-loads

### GET (schedule.teacher_load) – gets all teacher loads from the database

200 OK [=> TeacherLoad[]](schemas.md#schedule-teacher-load)

<a id="schedule-group-cohorts"></a>

## /group-cohorts

### GET (schedule.group_cohort) - gets all group cohorts from the database

200 OK [=> GroupCohort[]](schemas.md#schedule-group-cohort)

<a id="schedule-group-cohort-assignments"></a>

## /group-cohort-assignments

### GET (schedule.group_cohort_assignment) - gets all group cohort assignments from the database

200 OK [=> GroupCohortAssignment[]](schemas.md#schedule-group-cohort-assignment)

<a id="schedule-classrooms"></a>

## /classrooms

### GET (schedule.classroom) - gets all classrooms from the database

200 OK [=> Classroom[]](schemas.md#schedule-classroom)

<a id="schedule-semesters"></a>

## /semesters

### GET (schedule.semester) - gets all semesters from the database

200 OK [=> Semester[]](schemas.md#schedule-semester)

<a id="schedule-semester-discipline"></a>

## /semester-disciplines

### GET (schedule.semester_discipline) - gets all semester-discipline relations from the database

200 OK [=> SemesterDiscipline[]](schemas.md#schedule-semester-discipline)

<a id="schedule-study-loads"></a>

## /study-loads

### GET (schedule.study_load) - gets all study loads from the database

200 OK [=> StudyLoad[]](schemas.md#schedule-study-load)

<a id="schedule-lesson-slots"></a>

## /lesson-slots

### GET (schedule.lesson_slot) - gets all lesson slots from the database

200 OK [=> LessonSlot[]](schemas.md#schedule-lesson-slot)

<a id="schedule-lesson-occurrences"></a>

## /lesson-occurrences

### GET (schedule.lesson_occurrence) - gets all lesson occurrences from the database

200 OK [=> LessonOccurrence[]](schemas.md#schedule-lesson-occurrence)

<a id="schedule-load-data-into-generator"></a>

## /load-data-into-generator

### GET - loads all required generating data to schedule generator service

```json
["uuid (id of the curriculum semester)"]
```

204 NO CONTENT

400 BAD REQUEST or 500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-load-classrooms-into-generator"></a>

## /load-classrooms-into-generator

### GET - loads selected classrooms into schedule generator service

```json
["uuid (id of the classroom)"]
```

204 NO CONTENT

400 BAD REQUEST or 500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-extract-data-from-generator"></a>

## /extract-data-from-generator

### GET - extracts study loads and lessons from schedule generator

200 OK

500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-get-personal-schedule"></a>

## /get-personal-schedule

### ANY () - returns schedule for authorize user

```json
{
  "start_time": "time (start of the requested schedule period)",
  "end_time": "time (end of the requested schedule period)"
}
```

200 OK [=> []LessonOccurrence (full)](schemas.md#schedule-lesson-occurrence)

400 BAD REQUEST or 401 UNAUTHORIZE or 500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

# Schedule Generator Service

<a id="schedule-generator-init"></a>

## /init

### ANY - creates a new schedule generator using configuration from the request body, validates the input, and initializes the generator if it does not already exist

[<= GeneratorConfig](schemas.md#schedule-generator-generator-config)

201 CREATED => null

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-set-teachers"></a>

## /set-teachers

### ANY - assigns the teachers to the schedule generator

[<= Teacher[]](schemas.md#schedule-generator-teaacher)

200 OK

```json
{
  "message": "n teachers assigned"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-set-disciplines"></a>

## /set-disciplines

### ANY - assigns the disciplines to the schedule generator

[<= Disciplines[]](schemas.md#schedule-generator-discipline)

200 OK

```json
{
  "message": "n disciplines assigned"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-set-lesson-types"></a>

## /set-lesson-types

### ANY - assigns the lesson types to the schedule generator

[<= LessonType[]](schemas.md#schedule-generator-lesson-type)

200 OK

```json
{
  "message": "n lesson types assigned"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-set-lesson-type-assignments"></a>

## /set-lesson-type-assignments

### ANY - assigns the lesson type assignments to the schedule generator

[<= LessonTypeAssignment[]](schemas.md#schedule-generator-lesson-type-assignment)

200 OK

```json
{
  "message": "n lesson type assignments assigned"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-set-student-groups"></a>

## /set-student-groups

### ANY – assigns the student groups to the schedule generator

[<= GroupCohort[]](schemas.md#schedule-generator-group-cohort) as "group_cohorts",<br>
[<= GroupCohortAssignment[]](schemas.md#schedule-generator-group-cohort-assignment) as "group_cohort_assignments"

200 OK

```json
{
  "message": "n group cohorts assigned, n assignments assigned"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-set-study-loads"></a>

## /set-study-loads

### ANY – assigns the teacher loads to the schedule generator

[<= TeacherLoad[]](schemas.md#schedule-generator-teahcer-load)

200 OK

```json
{
  "message": "n teacher loads assigned"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-set-classrooms"></a>

## /set-classrooms

### ANY - assigns the classrooms to the schedule generator

[<= Classroom[]](schemas.md#schedule-generator-classroom)

200 OK

```json
{
  "message": "n classrooms assigned"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-submit-and-go"></a>

## /submit-and-go

### ANY – submit changes and go to the next generating step

200 OK

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-generate-days-for-lesson-types"></a>

## /generate-days-for-lesson-types

### ANY – generate day binding in student group for lesson types

200 OK [=> DaysForLessonTypes](schemas.md#schedule-generator-days-for-lesson-types)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-generate-bone-lessons"></a>

## /generate-bone-lesson

### ANY – generate lessons for all study loads, but within the week

200 OK [=> BoneLessons](schemas.md#schedule-generator-bone-lessons)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-assign-classrooms-to-bone-lessons"></a>

## /assign-classrooms-to-bone-lessons

### ANY – assign all available classrooms to bone lessons

200 OK [=> BoneLessonsWithC](schemas.md#schedule-generator-bone-lessons-with-c)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-build-schedule-skeleton"></a>

## /build-schedule-skeleton

### ANY – distribute all bone lessons across the full schedule

200 OK [=> GeneratedLessons](schemas.md#schedule-generator-generated-lessons)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-add-floating-lessons"></a>

## /add-floating-lessons

### ANY – adds missing lessons from all study loads as floating lessons

200 OK [=> GeneratedLessons](schemas.md#schedule-generator-generated-lessons)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-assign-classrooms-to-floating-lessons"></a>

## /assign-classrooms-to-floating-lessons

### ANY – assign all available classrooms to floating lessons

200 OK [=> GeneratedLessonsWithC](schemas.md#schedule-generator-generated-lessons-with-c)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-get-study-loads"></a>

## /get-study-loads

### ANY – gets all study loads from generator

200 OK [=> []StudyLoad](schemas.md#schedule-generator-study-load)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-get-lessons"></a>

## /get-lessons

### ANY – gets all scheduled lessons from generator

200 OK [=> []Lesson](schemas.md#schedule-generator-lesson)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-generator-get-fault"></a>

## /get-fault

### ANY – gets calculated fault oh the generated schedule

200 OK [=> Fault](schemas.md#schedule-generator-fault)

<a id="schedule-generator-default-generator-config"></a>

## /default-generator-config

### GET - retrieves the default schedule generator configuration from the service

200 OK [=> GeneratorConfig](schemas.md#schedule-generator-generator-config)

500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

# Student Service

<a id="student-semesters"></a>

## /semesters

### GET (student.semester) - gets all semesters from the database

200 OK [=> Semester[]](schemas.md#student-semester)

<a id="student-students"></a>

## /students

### GET (student.student) - gets all students from the database

200 OK [=> Student[]](schemas.md#student-student)

### POST (student.student) - adds a new student

```json
{
  "first_name": "string (student's first name)",
  "last_name": "string (student's last name)",
  "email": "string (student's email address)",
  "semester_id": "uuid (identifier of the associated semester)",

  // Optional fields
  "middle_name": "string (student's middle name)",
  "phone_number": "string (contact phone number)"
}
```

200 OK [=> Student](schemas.md#student-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="student-student-id"></a>

## /student/{id}

### GET (student.student) - finds student with an ID as an URL parameter

200 OK [=> Student](schemas.md#student-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (student.student) - deletes an student by its ID provided in the URL path

**200 OK** [=> Student](schemas.md#student-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (student.student) - updates a student by its ID with the data provided in the request body

```json
{
  "first_name": "string (student's first name)",
  "last_name": "string (student's last name)",
  "email": "string (student's email address)",
  "semester_id": "uuid (identifier of the associated semester)",

  // Optional fields
  "middle_name": "string (student's middle name)",
  "phone_number": "string (contact phone number)"
}
```

200 OK [=> Student](schemas.md#student-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

# Student Group Service

<a id="student-group-semesters"></a>

## /semesters

### GET (student_group.semester) - gets all semesters from the database

200 OK [=> Semester[]](schemas.md#student-group-semester)

<a id="student-group-students"></a>

## /students

### GET (student_group.student) - gets all students from the database

200 OK [=> Student[]](schemas.md#student-group-student)

<a id="student-group-group-cohorts"></a>

## /group-cohorts

### GET (student_group.group_cohort) - gets all group cohorts from the database

200 OK [=> GroupCohort[]](schemas.md#student-group-group-cohort)

### POST (student_group.group_cohort) - adds a new group cohort

```json
{
  "name": "string (name of the cohort)",
  "semester_id": "uuid (identifier of the associated semester)"
}
```

200 OK [=> GroupCohort](schemas.md#student-group-group-cohort)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="student-group-group-cohort-id"></a>

## /group-cohort/{id}

### GET (student_group.group_cohort) - finds group cohort with an ID as a URL parameter

200 OK [=> GroupCohort](schemas.md#student-group-group-cohort)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (student_group.group_cohort) - deletes a group cohort by its ID provided in the URL path

200 OK [=> GroupCohort](schemas.md#student-group-group-cohort)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (student_group.group_cohort) - updates a group cohort by its ID with the data provided in the request body

```json
{
  "name": "string (name of the cohort)",
  "semester_id": "uuid (identifier of the associated semester)"
}
```

200 OK [=> GroupCohort](schemas.md#student-group-group-cohort)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="student-group-student-groups"></a>

## /student-groups

### GET (student_group.student_group) - gets all student groups from the database

200 OK [=> StudentGroup[]](schemas.md#student-group-student-group)

### POST (student_group.student_group) - adds a new student group

```json
{
  "name": "string (name of the group)",
  "group_cohort_id": "uuid (identifier of the associated group cohort)"
}
```

200 OK [=> StudentGroup](schemas.md#student-group-student-group)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="student-group-student-group-id"></a>

## /student-group/{id}

### GET (student_group.student_group) - finds a student group by its ID provided in the URL path

200 OK [=> StudentGroup](schemas.md#student-group-student-group)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (student_group.student_group) - deletes a student group by its ID provided in the URL path

200 OK [=> StudentGroup](schemas.md#student-group-student-group)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (student_group.student_group) - updates a student group by its ID with the data provided in the request body

```json
{
  "name": "string (name of the group)",
  "group_cohort_id": "uuid (identifier of the associated group cohort)"
}
```

200 OK [=> StudentGroup](schemas.md#student-group-student-group)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="student-group-group-members"></a>

## /group-members

### GET (student_group.group_member) - gets all group members from the database

200 OK [=> GroupMember[]](schemas.md#student-group-group-member)

### POST (student_group.group_member) - adds a new group member

```json
{
  "studentId": "uuid (identifier of the student)",
  "group_cohort_id": "uuid (identifier of the associated group cohort)",

  // Optional fields
  "student_group_id": "uuid (identifier of the student group)"
}
```

200 OK [=> GroupMember](schemas.md#student-group-group-member)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="student-group-group-member-id"></a>

## /group-members/{id}

### GET (student_group.group_member) - finds a group member by its ID provided in the URL path

200 OK [=> GroupMember](schemas.md#student-group-group-member)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (student_group.group_member) - deletes a group member by its ID provided in the URL path

200 OK [=> GroupMember](schemas.md#student-group-group-member)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (student_group.group_member) - updates a group member by its ID with the data provided in the request body

```json
{
  "studentId": "uuid (identifier of the student)",
  "group_cohort_id": "uuid (identifier of the associated group cohort)",

  // Optional fields
  "student_group_id": "uuid (identifier of the student group)"
}
```

200 OK [=> GroupMember](schemas.md#student-group-group-member)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="student-group-lesson-types"></a>

## /lesson-types

### GET (student_group.lesson_type) - gets all lesson types from the database

200 OK [=> LessonType[]](schemas.md#student-group-lesson-type)

<a id="student-group-disciplines"></a>

## /disciplines

### GET (student_group.discipline) – gets all disciplines from the database

200 OK [=> Discipline[]](schemas.md#student-group-discipline)

<a id="student-group-group-cohort-assignments"></a>

## /group-cohort-assignments

### GET (student_group.group_cohort_assignment) - gets all group cohort assignments from the database

200 OK [=> GroupCohortAssignment[]](schemas.md#student-group-group-cohort-assignment)

### POST (student_group.group_cohort_assignment) - adds a new group cohort assignment

```json
{
  "group_cohort_id": "uuid (identifier of the associated group cohort)",
  "discipline_id": "uuid (identifier of the discipline)",
  "lesson_type_id": "uuid (identifier of the lesson type)"
}
```

200 OK [=> GroupCohortAssignment](schemas.md#student-group-group-cohort-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="student-group-group-cohort-assignment-id"></a>

## /group-cohort-assignments/{id}

### GET (student_group.group_cohort_assignment) - finds a group cohort assignment by its ID provided in the URL path

200 OK [=> GroupCohortAssignment](schemas.md#student-group-group-cohort-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (student_group.group_cohort_assignment) - deletes a group cohort assignment by its ID provided in the URL path

200 OK [=> GroupCohortAssignment](schemas.md#student-group-group-cohort-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (student_group.group_cohort_assignment) - updates a group cohort assignment by its ID with the data provided in the request body

```json
{
  "group_cohort_id": "uuid (identifier of the associated group cohort)",
  "discipline_id": "uuid (identifier of the discipline)",
  "lesson_type_id": "uuid (identifier of the lesson type)"
}
```

200 OK [=> GroupCohortAssignment](schemas.md#student-group-group-cohort-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

# Curriculum Service

<a id="curriculum-curriculum"></a>

## /curriculums

### GET (curriculum.curriculum) - gets all curriculums from the database

200 OK [=> Curriculum[]](schemas.md#curriculum-curriculum)

### POST (curriculum.curriculum) - adds a new curriculum

```json
{
  "name": "string (curriculum name)",
  "duration_years": "integer (number of years for the curriculum)",
  "effective_from": "timestamp (when this curriculum version becomes effective)",

  // Optional fields
  "effective_to": "timestamp (when this curriculum version ends, nullable)"
}
```

200 OK [=> Curriculum](schemas.md#curriculum-curriculum)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-curriculum-id"></a>

## /curriculum/{id}

### GET (curriculum.curriculum) - finds curriculum with an ID as an URL parameter

200 OK [=> Curriculum](schemas.md#curriculum-curriculum)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (curriculum.curriculum) - deletes a curriculum by its ID provided in the URL path

200 OK [=> Curriculum](schemas.md#curriculum-curriculum)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (curriculum.curriculum) - updates a curriculum by its ID with the data provided in the request body

```json
{
  "name": "string (curriculum name)",
  "duration_years": "integer (number of years for the curriculum)",
  "effective_from": "timestamp (when this curriculum version becomes effective)",

  // Optional fields
  "effective_to": "timestamp (when this curriculum version ends, nullable)"
}
```

200 OK [=> Curriculum](schemas.md#curriculum-curriculum)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-semesters"></a>

## /semesters

### GET (curriculum.semester) - gets all semesters from the database

200 OK [=> Semester[]](schemas.md#curriculum-semester)

### POST (curriculum.semester) - adds a new semester

```json
{
  "curriculum_id": "uuid (identifier of the associated curriculum)",
  "number": "integer (semester number within the curriculum)"
}
```

200 OK [=> Semester](schemas.md#curriculum-semester)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-semester-id"></a>

## /semester/{id}

### GET (curriculum.semester) - finds semester with an ID as an URL parameter

200 OK [=> Semester](schemas.md#curriculum-semester)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (curriculum.semester) - deletes a semester by its ID provided in the URL path

200 OK [=> Semester](schemas.md#curriculum-semester)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (curriculum.semester) - updates a semester by its ID with the data provided in the request body

```json
{
  "curriculum_id": "uuid (identifier of the associated curriculum)",
  "number": "integer (semester number within the curriculum)"
}
```

200 OK [=> Semester](schemas.md#curriculum-semester)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-lesson-types"></a>

## /lesson-types

### GET (curriculum.lesson_type) - gets all lesson types from the database

200 OK [=> LessonType[]](schemas.md#curriculum-lesson-type)

### POST (curriculum.lesson_type) - adds a new lesson type

```json
{
  "name": "string (name of the lesson type)",
  "hours_value": "integer (number of hours per lesson)"
}
```

200 OK [=> LessonType](schemas.md#curriculum-lesson-type)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-lesson-type-id"></a>

## /lesson-type/{id}

### GET (curriculum.lesson_type) - finds lesson type with an ID as an URL parameter

200 OK [=> LessonType](schemas.md#curriculum-lesson-type)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (curriculum.lesson_type) - deletes a lesson type by its ID provided in the URL path

200 OK [=> LessonType](schemas.md#curriculum-lesson-type)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (curriculum.lesson_type) - updates a lesson type by its ID with the data provided in the request body

```json
{
  "name": "string (name of the lesson type)",
  "hours_value": "integer (number of hours per lesson)"
}
```

200 OK [=> LessonType](schemas.md#curriculum-lesson-type)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-disciplines"></a>

## /disciplines

### GET (curriculum.discipline) – gets all disciplines from the database

200 OK [=> Discipline[]](schemas.md#curriculum-discipline)

### POST (curriculum.discipline) – adds a new discipline

```json
{
  "name": "string (name of the discipline)"
}
```

200 OK [=> Discipline](schemas.md#curriculum-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-discipline-id"></a>

## /discipline/{id}

### GET (curriculum.discipline) – finds a discipline with an ID as a URL parameter

200 OK [=> Discipline](schemas.md#curriculum-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (curriculum.discipline) – deletes a discipline by its ID provided in the URL path

200 OK [=> Discipline](schemas.md#curriculum-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (curriculum.discipline) – updates a discipline by its ID with the data provided in the request body

```json
{
  "name": "string (name of the discipline)"
}
```

200 OK [=> Discipline](schemas.md#curriculum-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-lesson-type-assignments"></a>

## /lesson-type-assignments

### GET (curriculum.lesson_type_assignment) – gets all lesson type assignments from the database

200 OK [=> LessonTypeAssignment[]](schemas.md#curriculum-lesson-type-assignment)

### POST (curriculum.lesson_type_assignment) – adds a new lesson type assignment

```json
{
  "lesson_type_id": "uuid (identifier of the lesson type)",
  "discipline_id": "uuid (identifier of the discipline)",
  "required_hours": "integer (number of hours required for this lesson type in this discipline)"
}
```

200 OK [=> LessonTypeAssignment](schemas.md#lcurriculum-esson-type-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-lesson-type-assignment-id"></a>

## /lesson-type-assignment/{id}

### GET (curriculum.lesson_type_assignment) – finds a lesson type assignment with an ID as a URL parameter

200 OK [=> LessonTypeAssignment](schemas.md#curriculum-lesson-type-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (curriculum.lesson_type_assignment) – deletes a lesson type assignment by its ID provided in the URL path

200 OK [=> LessonTypeAssignment](schemas.md#curriculum-lesson-type-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (curriculum.lesson_type_assignment) – updates a lesson type assignment by its ID with the data provided in the request body

```json
{
  "lesson_type_id": "uuid (identifier of the lesson type)",
  "discipline_id": "uuid (identifier of the discipline)",
  "required_hours": "integer (number of hours required for this lesson type in this discipline)"
}
```

200 OK [=> LessonTypeAssignment](schemas.md#curriculum-lesson-type-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-semester-disciplines"></a>

## /semester-disciplines

### GET (curriculum.semester_discipline) – gets all semester discipline relations from the database

200 OK [=> SemesterDiscipline[]](schemas.md#curriculum-semester-discipline)

### POST (curriculum.semester_discipline) – adds a new semester discipline relation

```json
{
  "semester_id": "uuid (identifier of the semester)",
  "discipline_id": "uuid (identifier of the discipline)"
}
```

200 OK [=> SemesterDiscipline](schemas.md#curriculum-semester-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

---

<a id="curriculum-semester-discipline-id"></a>

## /semester-discipline/{id}

### GET (curriculum.semester_discipline) – finds a semester discipline relation by ID (URL parameter)

200 OK [=> SemesterDiscipline](schemas.md#curriculum-semester-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (curriculum.semester_discipline) – deletes a semester discipline relation by its ID (URL path)

200 OK [=> SemesterDiscipline](schemas.md#curriculum-semester-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

# Teacher Load Service

<a id="teacher-load-teachers"></a>

## /teachers

### GET (teacher_load.teacher) - gets all teachers from the database

200 OK [=> Teacher[]](schemas.md#teacher-load-teacher)

<a id="teacher-load-lesson-types"></a>

## /lesson-types

### GET (teacher_load.lesson_type) - gets all lesson types from the database

200 OK [=> LessonType[]](schemas.md#teacher-load-lesson-type)

<a id="teacher-load-disciplines"></a>

## /disciplines

### GET (teacher_load.discipline) – gets all disciplines from the database

200 OK [=> Discipline[]](schemas.md#teacher-load-discipline)

<a id="teacher-load-teacher-loads"></a>

## /teacher-loads

### GET (teacher_load.teacher_load) – gets all teacher loads from the database

200 OK [=> TeacherLoad[]](schemas.md#teacher-load-teacher-load)

### POST (teacher_load.teacher_load) – adds a new teacher load

```json
{
  "teacher_id": "uuid (unique identifier of the teacher)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)",
  "group_count": "integer (number of groups assigned for this load)"
}
```

200 OK [=> TeacherLoad](schemas.md#teacher-load-teacher-load)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="teacher-load-teacher-load-id"></a>

## /teacher-load/{id}

### GET (teacher_load.teacher_load) – finds a teacher load with an ID as an URL parameter

200 OK [=> TeacherLoad](schemas.md#teacher-load-teacher-load)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (teacher_load.teacher_load) – deletes a teacher load by its ID provided in the URL path

200 OK [=> TeacherLoad](schemas.md#teacher-load-teacher-load)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (teacher_load.teacher_load) – updates a teacher load by its ID with the data provided in the request body

```json
{
  "teacher_id": "uuid (unique identifier of the teacher)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)",
  "group_cohort_id": "uuid (unique identifier of the group cohort)",
  "group_count": "integer (number of groups assigned for this load)"
}
```

200 OK [=> TeacherLoad](schemas.md#teacher-load-teacher-load)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

# Asset Service

<a id="asset-classrooms"></a>

## /classrooms

### GET (asset.classroom) - gets all classrooms from the database

200 OK [=> Classroom[]](schemas.md#asset-classroom)

### POST (asset.classroom) - adds a new classroom

```json
{
  "number": "string (classroom number or label)",
  "capacity": "integer (maximum number of occupants in the classroom)"
}
```

200 OK [=> Classroom](schemas.md#asset-classroom)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="asset-classroom-id"></a>

## /classroom/{id}

### GET (asset.classroom) - finds classroom with an ID as an URL parameter

200 OK [=> Classroom](schemas.md#asset-classroom)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (asset.classroom) - deletes a classroom by its ID provided in the URL path

200 OK [=> Classroom](schemas.md#asset-classroom)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (asset.classroom) - updates a classroom by its ID with the data provided in the request body

```json
{
  "number": "string (classroom number or label)",
  "capacity": "integer (maximum number of occupants in the classroom)"
}
```

200 OK [=> Classroom](schemas.md#asset-classroom)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

# Auth Service

<a id="auth-permissions"></a>

## /permissions

### GET (auth.permission) - gets all permissions from the database

200 OK [=> Permission[]](schemas.md#auth-permission)

### POST (auth.permission) - adds a new permission

```json
{
  "name": "string (unique name of the permission)"
}
```

200 OK [=> Permission](schemas.md#auth-permission)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-permission-id"></a>

## /permission/{id}

### GET (auth.permission) - finds permission with an ID as a URL parameter

200 OK [=> Permission](schemas.md#auth-permission)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (auth.permission) - deletes a permission by its ID provided in the URL path

200 OK [=> Permission](schemas.md#auth-permission)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (auth.permission) - updates a permission by its ID with the data provided in the request body

```json
{
  "name": "string (unique name of the permission)"
}
```

200 OK [=> Permission](schemas.md#auth-permission)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-roles"></a>

## /roles

### GET (auth.role) - gets all roles from the database

200 OK [=> Role[]](schemas.md#auth-role)

### POST (auth.role) - adds a new role

```json
{
  "name": "string (unique name of the role)"
}
```

200 OK [=> Role](schemas.md#auth-role)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-role-id"></a>

## /role/{id}

### GET (auth.role) - finds role with an ID as a URL parameter

200 OK [=> Role](schemas.md#auth-role)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (auth.role) - deletes a role by its ID provided in the URL path

200 OK [=> Role](schemas.md#auth-role)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (auth.role) - updates a role by its ID with the data provided in the request body

```json
{
  "name": "string (unique name of the role)"
}
```

200 OK [=> Role](schemas.md#auth-role)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-role-permissions"></a>

## /role-permissions

### GET (auth.role_permission) - gets all role-permission links from the database

200 OK [=> RolePermissions[]](schemas.md#auth-role-permissions)

### POST (auth.role_permission) - adds a new role-permission link

```json
{
  "role_id": "uuid (associated role identifier)",
  "permission_id": "uuid (associated permission identifier)"
}
```

200 OK [=> RolePermissions](schemas.md#auth-role-permissions)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-role-permission-id"></a>

## /role-permission/{id}

### GET (auth.role_permission) - finds role-permission link by ID from URL parameter

200 OK [=> RolePermissions](schemas.md#auth-role-permissions)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (auth.role_permission) - deletes a role-permission link by ID from URL parameter

200 OK [=> RolePermissions](schemas.md#auth-role-permissions)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-services"></a>

## /services

### GET (auth.service) - gets all services from the database

200 OK [=> Service[]](schemas.md#auth-service)

### POST (auth.service) - adds a new service

```json
{
  "name": "string (unique name of the service)",
  "secrete": "string (service secret key)"
}
```

200 OK [=> Service](schemas.md#auth-service)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-service-id"></a>

## /service/{id}

### GET (auth.service) - finds service with an ID as a URL parameter

200 OK [=> Service](schemas.md#auth-service)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (auth.service) - deletes a service by its ID provided in the URL path

200 OK [=> Service](schemas.md#auth-service)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (auth.service) - updates a service by its ID with the data provided in the request body

```json
{
  "name": "string (unique name of the service)",
  "secrete": "string (service secret key)"
}
```

200 OK [=> Service](schemas.md#auth-service)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-service-permissions"></a>

## /service-permissions

### GET (auth.service_permission) - gets all service-permission links from the database

200 OK [=> ServicePermissions[]](schemas.md#auth-service-permissions)

### POST (auth.service_permission) - adds a new service-permission link

```json
{
  "service_id": "uuid (associated service identifier)",
  "permission_id": "uuid (associated permission identifier)"
}
```

200 OK [=> ServicePermissions](schemas.md#auth-service-permissions)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-service-permission-id"></a>

## /service-permission/{id}

### GET (auth.service_permission) - finds service-permission link by ID from URL parameter

200 OK [=> ServicePermissions](schemas.md#auth-service-permissions)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (auth.service_permission) - deletes a service-permission link by ID from URL parameter

200 OK [=> ServicePermissions](schemas.md#auth-service-permissions)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-users"></a>

## /users

### GET (auth.user) - gets all users from the database

200 OK [=> User[]](schemas.md#auth-user)

### POST (auth.user) - creates a new user

```json
{
  "login": "string (unique username for the user)",
  "role_id": "uuid (assigned role identifier)"
}
```

200 OK [=> User](schemas.md#auth-user)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-user-id"></a>

## /user/{id}

### GET (auth.user) - finds user by ID provided as a URL parameter

200 OK [=> User](schemas.md#auth-user)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (auth.user) - deletes a user by ID provided in the URL path

200 OK [=> User](schemas.md#auth-user)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (auth.user) - updates a user by ID with data provided in request body

```json
{
  "role_id": "uuid (assigned role identifier)",
  "is_default_password": "boolean (indicates whether the user is using a default password)"
}
```

200 OK [=> User](schemas.md#auth-user)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-login"></a>

## /login

### ANY - authenticates user and returns access credentials (tokens)

```json
{
  "login": "string (user login/username)",
  "password": "string (user password)"
}
```

200 OK

```json
{
  "id": "uuid (unique identifier of the user)",
  "login": "string (user login/username)",
  "role": "string (user role name)",
  "is_default_password": "boolean (indicates whether the user is using a default password)",
  "access_token": "string (JWT token used for authentication)",
  "refresh_token": "string (token used to refresh access token)"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-refresh"></a>

## /refresh

### ANY - issues a new access token using a valid refresh token

```json
{
  "access_token": "string (JWT token used for authentication)",
  "refresh_token": "string (token used to refresh access token)"
}
```

200 OK

```json
{
  "access_token": "string (JWT token used for authentication)",
  "refresh_token": "string (token used to refresh access token)"
}
```

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-reset-password-id"></a>

## /reset-password-id

### ANY (auth.user.reset_password) - resets the password for a user specified by ID in the URL path

204 NO CONTENT

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="auth-change-password"></a>

## /change-password

### ANY - changes the user's password

```json
{
  "id": "uuid (unique identifier of the user)",
  "login": "string (user login/username)",
  "password": "string (current user password, required for verification)",
  "new_password": "string (new password to be set for the user)"
}
```

204 NO CONTENT

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

# Course Service

<a id="course-teachers"></a>

## /teachers

### GET (course.teacher) - gets all teachers from the database

**200 OK** [=> Teacher[]](schemas.md#course-teacher)

<a id="course-students"></a>

## /students

### GET (course.student) - gets all students from the database

200 OK [=> Student[]](schemas.md#course-student)

<a id="course-courses"></a>

## /courses

### GET (course.course) - gets all courses from the database

200 OK [=> Course[]](schemas.md#course-course)

### POST (course.course) - adds a new course

```json
{
  "manager_id": "uuid | null (identifier of the course manager/teacher)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the course)",
  "description": "string (detailed course description)"
}
```

200 OK [=> Course](schemas.md#course-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-course-id"></a>

## /course/{id}

### GET (course.course) - finds a course by ID

200 OK [=> Course](schemas.md#course-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (course.course) - updates a course by ID

```json
{
  "manager_id": "uuid | null (identifier of the course manager/teacher)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the course)",
  "description": "string (detailed course description)"
}
```

200 OK [=> Course](schemas.md#course-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-student-courses"></a>

## /student-courses

### GET (course.student_course) - gets all student-course relations

200 OK [=> StudentCourse[]](schemas.md#course-student-course)

### POST (course.student_course) - adds a student to a course

```json
{
  "student_id": "uuid (identifier of the student)",
  "course_id": "uuid (identifier of the course)"
}
```

200 OK [=> StudentCourse](schemas.md#course-student-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-student-course-id"></a>

## /student-course/{id}

### GET (course.student_course) - finds a student-course relation by ID

200 OK [=> StudentCourse](schemas.md#course-student-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (course.student_course) - deletes a student-course relation by ID

200 OK [=> StudentCourse](schemas.md#course-student-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-teacher-courses"></a>

## /teacher-courses

### GET (course.teacher_course) - gets all teacher-course relations

200 OK [=> TeacherCourse[]](schemas.md#course-teacher-course)

### POST (course.teacher_course) - assigns a teacher to a course

```json
{
  "teacher_id": "uuid (identifier of the teacher)",
  "course_id": "uuid (identifier of the course)"
}
```

200 OK [=> TeacherCourse](schemas.md#course-teacher-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-teacher-course-id"></a>

## /teacher-course/{id}

### GET (course.teacher_course) - finds a teacher-course relation by ID

200 OK [=> TeacherCourse](schemas.md#course-teacher-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (course.teacher_course) - deletes a teacher-course relation by ID

200 OK [=> TeacherCourse](schemas.md#course-teacher-course)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-tasks"></a>

## /tasks

### GET (course.task) - gets all tasks from the database

200 OK [=> Task[]](schemas.md#course-task)

### POST (course.task) - adds a new task

```json id="task_post_01"
{
  "course_id": "uuid (identifier of the course)",
  "slug": "string (unique slug used internally)",
  "title": "string (title of the task)",
  "description": "string (detailed description of the task)",
  "max_mark": "float (maximum achievable mark)",
  "deadline": "timestamp (deadline for submission)"
}
```

200 OK [=> Task](schemas.md#course-task)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-task-id"></a>

## /task/{id}

### GET (course.task) - finds a task by ID

200 OK [=> Task](schemas.md#course-task)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (course.task) - deletes a task by ID

200 OK [=> Task](schemas.md#course-task)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT (course.task) - updates a task by ID

```json id="task_put_01"
{
  "course_id": "uuid (identifier of the course)",
  "slug": "string (unique slug used internally)",
  "title": "string (title of the task)",
  "description": "string (detailed description of the task)",
  "max_mark": "float (maximum achievable mark)",
  "deadline": "timestamp (deadline for submission)"
}
```

200 OK [=> Task](schemas.md#course-task)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-task-students"></a>

## /task-students

### GET (course.task_student) - gets all task-student relations

200 OK [=> TaskStudent[]](schemas.md#course-task-student)

### POST (course.task_student) - assigns a task to a student / submits a task

```json id="task_student_post_01"
{
  "task_id": "uuid (identifier of the task)",
  "student_id": "uuid (identifier of the student)",
  "mark": "float | null (assigned mark)",
  "submission_time": "timestamp | null (submission time)"
}
```

200 OK [=> TaskStudent](schemas.md#course-task-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="course-task-student-id"></a>

## /task-student/{id}

### GET (course.task_student) - finds a task-student relation by ID

200 OK [=> TaskStudent](schemas.md#course-task-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE (course.task_student) - deletes a task-student relation by ID

200 OK [=> TaskStudent](schemas.md#course-task-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)
