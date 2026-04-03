# API Documentation

# ApiGateway

## /seed

### ANY - seeds databases of all services.

204 NO CONTENT

500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

## /

### ANY - routes and proxies incoming requests to the appropriate backend service based on the request path.

- Employee Service (/employee)
  - [/academic-ranks](#employee-academic-ranks)
  - [/academic-rank/{id}](#employee-academic-rank-id)
  - [/academic-degrees](#employee-academic-degrees)
  - [/academic-degree/{id}](#employee-academic-degree-id)
  - [/employees](#employee-employees)
  - [/employee/{id}](#employee-employee-id)
  - [/teachers](#employee-teachers)
  - [/teacher/{id}](#employee-teacher-id)

- Schedule Service (/schedule)
  - [/academic-ranks](#schedule-academic-ranks)
  - [/academic-rank/{id}](#schedule-academic-rank-id)
  - [/lesson-types](#schedule-lesson-types)
  - [/teachers](#schedule-teachers)
  - [/disciplines](#schedule-disciplines)
  - [/lesson-type-assignments](#schedule-lesson-type-assignments)
  - [/students](#schedule-students)
  - [/student-groups](#schedule-student-groups)
  - [/group_members](#schedule-group-members)
  - [/teacher-loads](#schedule-teacher-loads)

- Schedule Generator Service (/schedule-generator)
  - [/init](#schedule-generator-init)
  - [/set-teachers](#schedule-generator-set-teachers)
  - [/default-generator-config](#schedule-generator-default-generator-config)

- Student Service (/student)
  - [/semesters](#student-semesters)
  - [/students](#student-students)
  - [/student/{id}](#student-student-id)

- Student Group Service (/student-group)
  - [/semesters](#student-group-semesters)
  - [/students](#student-group-students)
  - [/group-cohorts](#student-group-group-cohorts)
  - [/group-cohort/{id}](#student-group-group-cohort-id)
  - [/student-groups](#student-group-student-groups)
  - [/student-group/{id}](#student-group-student-group-id)
  - [/group_members](#student-group-group-members)
  - [/group_member/{id}](#student-group-group-member-id)

- Curriculum Service (/curriculum)
  - [/curriculums](#curriculum-curriculum)
  - [/curriculum/{id}](#curriculum-curriculum-id)
  - [/semesters](#curriculum-semesters)
  - [/semester/{id}](#curriculum-semester-id)
  - [/lesson-types](#curriculum-lesson-types)
  - [/lesson-type/{id}](#curriculum-lesson-type-id)
  - [/disciplines](#curriculum-disciplines)
  - [/discipline/{id}](#curriculum-discipline-id)
  - [/lesson-type-assignments](#curriculum-lesson-type-assignments)
  - [/lesson-type-assignment/{id}](#curriculum-lesson-type-assignment-id)
  - [/semester-disciplines](#curriculum-semester-disciplines)
  - [/semester-discipline/{id}](#curriculum-semester-discipline-id)

- Teacher Load Service (/teacher-load)
  - [/teachers](#teacher-load-teachers)
  - [/group-cohorts](#teacher-load-group-cohorts)
  - [/lesson-types](#teacher-load-lesson-types)
  - [/disciplines](#teacher-load-disciplines)
  - [/teacher-loads](#teacher-load-teacher-loads)
  - [/teacher-load/{id}](#teacher-load-teacher-load-id)

400 BAD REQUEST or 500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

# Employee Service

<a id="employee-academic-ranks"></a>

## /academic-ranks

### GET - gets all academic ranks from the database

200 OK [=> AcademicRank[]](schemas.md#employee-academic-rank)

### POST - adds a new academic rank

```json
{
  "title": "string (human-readable name of the rank)"
}
```

200 OK [=> AcademicRank](schemas.md#employee-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-academic-rank-id"></a>

## /academic-rank/{id}

### GET - finds academic rank with an ID as an URL parameter

200 OK [=> AcademicRank](schemas.md#employee-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes an academic rank by its ID provided in the URL path

200 OK [=> AcademicRank](schemas.md#employee-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates an academic rank by its ID with the data provided in the request body

```json
{
  "title": "string (human-readable name of the rank)"
}
```

200 OK [=> AcademicRank](schemas.md#employee-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-academic-degrees"></a>

## /academic-degrees

### GET - gets all academic degrees from the database

200 OK [=> AcademicDegree[]](schemas.md#employee-academic-degree)

### POST - adds a new academic degree

```json
{
  "title": "string (human-readable name of the degree)"
}
```

200 OK [=> AcademicDegree](schemas.md#employee-academic-degree)

400 BAD REQUEST\*\* [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-academic-degree-id"></a>

## /academic-degree/{id}

### GET - finds academic degree with an ID as a URL parameter

200 OK [=> AcademicDegree](schemas.md#employee-academic-degree)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes an academic degree by its ID provided in the URL path

**200 OK** [=> AcademicDegree](schemas.md#employee-academic-degree)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates an academic degree by its ID with the data provided in the request body

```json
{
  "title": "string (human-readable name of the degree)"
}
```

200 OK [=> AcademicDegree](schemas.md#employee-academic-degree)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

<a id="employee-employees"></a>

## /employees

### GET - gets all employees from the database

**200 OK** [=> Employee[]](schemas.md#employee-employee)

### POST - adds a new employee

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

### GET - finds employee with an ID as a URL parameter

**200 OK** [=> Employee](schemas.md#employee-employee)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes an employee by its ID provided in the URL path

**200 OK** [=> Employee](schemas.md#employee-employee)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates an employee by its ID with the data provided in the request body

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

### GET - gets all teachers from the database

**200 OK** [=> Teacher[]](schemas.md#employee-teacher)

### POST - adds a new teacher

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

### GET - finds teacher with an ID as a URL parameter

**200 OK** [=> Teacher](schemas.md#employee-teacher)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes a teacher by its ID provided in the URL path

**200 OK** [=> Teacher](schemas.md#employee-teacher)

**400 BAD REQUEST** [-> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a teacher by its ID with the data provided in the request body

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

### GET - gets all academic ranks from the database

200 OK [=> AcademicRank[]](schemas.md#schedule-academic-rank)

<a id="schedule-academic-rank-id"></a>

## /academic-rank/{id}

### GET - finds academic rank with an ID as an URL parameter

200 OK [=> AcademicRank](schemas.md#schedule-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates an academic rank by its ID with the data provided in the request body

```json
{
  "priority": "int (determines the rank's priority: higher value = higher rank)"
}
```

200 OK [=> AcademicRank](schemas.md#schedule-academic-rank)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="schedule-teachers"></a>

## /teachers

### GET - gets all teachers from the database

**200 OK** [=> Teacher[]](schemas.md#schedule-teacher)

<a id="schedule-lesson-types"></a>

## /lesson-types

### GET - gets all lesson types from the database

200 OK [=> LessonType[]](schemas.md#schedule-lesson-type)

<a id="schedule-disciplines"></a>

## /disciplines

### GET – gets all disciplines from the database

200 OK [=> Discipline[]](schemas.md#schedule-discipline)

<a id="schedule-lesson-type-assignments"></a>

## /lesson-type-assignments

### GET – gets all lesson type assignments from the database

200 OK [=> LessonTypeAssignment[]](schemas.md#schedule-lesson-type-assignment)

<a id="schedule-students"></a>

## /students

### GET - gets all students from the database

200 OK [=> Student[]](schemas.md#schedule-student)

<a id="schedule-student-groups"></a>

## /student-groups

### GET - gets all student groups from the database

200 OK [=> StudentGroup[]](schemas.md#schedule-student-group)

<a id="schedule-group-members"></a>

## /group-members

### GET - gets all group members from the database

200 OK [=> GroupMember[]](schemas.md#schedule-group-member)

<a id="schedule-teacher-loads"></a>

## /teacher-loads

### GET – gets all teacher loads from the database

200 OK [=> TeacherLoad[]](schemas.md#schedule-teacher-load)

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

<a id="schedule-generator-default-generator-config"></a>

## /default-generator-config

### GET - retrieves the default schedule generator configuration from the service

200 OK [=> GeneratorConfig](schemas.md#schedule-generator-generator-config)

500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

# Student Service

<a id="student-semesters"></a>

## /semesters

### GET - gets all semesters from the database

200 OK [=> Semester[]](schemas.md#student-semester)

<a id="student-students"></a>

## /students

### GET - gets all students from the database

200 OK [=> Student[]](schemas.md#student-student)

### POST - adds a new student

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

### GET - finds student with an ID as an URL parameter

200 OK [=> Student](schemas.md#student-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes an student by its ID provided in the URL path

**200 OK** [=> Student](schemas.md#student-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a student by its ID with the data provided in the request body

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

### GET - gets all semesters from the database

200 OK [=> Semester[]](schemas.md#student-group-semester)

<a id="student-group-students"></a>

## /students

### GET - gets all students from the database

200 OK [=> Student[]](schemas.md#student-group-student)

<a id="student-group-group-cohorts"></a>

## /group-cohorts

### GET - gets all group cohorts from the database

200 OK [=> GroupCohort[]](schemas.md#student-group-group-cohort)

### POST - adds a new group cohort

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

### GET - finds group cohort with an ID as a URL parameter

200 OK [=> GroupCohort](schemas.md#student-group-group-cohort)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes a group cohort by its ID provided in the URL path

200 OK [=> GroupCohort](schemas.md#student-group-group-cohort)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a group cohort by its ID with the data provided in the request body

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

### GET - gets all student groups from the database

200 OK [=> StudentGroup[]](schemas.md#student-group-student-group)

### POST - adds a new student group

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

### GET - finds a student group by its ID provided in the URL path

200 OK [=> StudentGroup](schemas.md#student-group-student-group)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes a student group by its ID provided in the URL path

200 OK [=> StudentGroup](schemas.md#student-group-student-group)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a student group by its ID with the data provided in the request body

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

### GET - gets all group members from the database

200 OK [=> GroupMember[]](schemas.md#student-group-group-member)

### POST - adds a new group member

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

### GET - finds a group member by its ID provided in the URL path

200 OK [=> GroupMember](schemas.md#student-group-group-member)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes a group member by its ID provided in the URL path

200 OK [=> GroupMember](schemas.md#student-group-group-member)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a group member by its ID with the data provided in the request body

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

# Curriculum Service

<a id="curriculum-curriculum"></a>

## /curriculums

### GET - gets all curriculums from the database

200 OK [=> Curriculum[]](schemas.md#curriculum-curriculum)

### POST - adds a new curriculum

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

### GET - finds curriculum with an ID as an URL parameter

200 OK [=> Curriculum](schemas.md#curriculum-curriculum)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes a curriculum by its ID provided in the URL path

200 OK [=> Curriculum](schemas.md#curriculum-curriculum)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a curriculum by its ID with the data provided in the request body

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

### GET - gets all semesters from the database

200 OK [=> Semester[]](schemas.md#curriculum-semester)

### POST - adds a new semester

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

### GET - finds semester with an ID as an URL parameter

200 OK [=> Semester](schemas.md#curriculum-semester)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes a semester by its ID provided in the URL path

200 OK [=> Semester](schemas.md#curriculum-semester)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a semester by its ID with the data provided in the request body

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

### GET - gets all lesson types from the database

200 OK [=> LessonType[]](schemas.md#curriculum-lesson-type)

### POST - adds a new lesson type

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

### GET - finds lesson type with an ID as an URL parameter

200 OK [=> LessonType](schemas.md#curriculum-lesson-type)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE - deletes a lesson type by its ID provided in the URL path

200 OK [=> LessonType](schemas.md#curriculum-lesson-type)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a lesson type by its ID with the data provided in the request body

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

### GET – gets all disciplines from the database

200 OK [=> Discipline[]](schemas.md#curriculum-discipline)

### POST – adds a new discipline

```json
{
  "name": "string (name of the discipline)"
}
```

200 OK [=> Discipline](schemas.md#curriculum-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-discipline-id"></a>

## /discipline/{id}

### GET – finds a discipline with an ID as a URL parameter

200 OK [=> Discipline](schemas.md#curriculum-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE – deletes a discipline by its ID provided in the URL path

200 OK [=> Discipline](schemas.md#curriculum-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT – updates a discipline by its ID with the data provided in the request body

```json
{
  "name": "string (name of the discipline)"
}
```

200 OK [=> Discipline](schemas.md#curriculum-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

<a id="curriculum-lesson-type-assignments"></a>

## /lesson-type-assignments

### GET – gets all lesson type assignments from the database

200 OK [=> LessonTypeAssignment[]](schemas.md#curriculum-lesson-type-assignment)

### POST – adds a new lesson type assignment

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

### GET – finds a lesson type assignment with an ID as a URL parameter

200 OK [=> LessonTypeAssignment](schemas.md#curriculum-lesson-type-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE – deletes a lesson type assignment by its ID provided in the URL path

200 OK [=> LessonTypeAssignment](schemas.md#curriculum-lesson-type-assignment)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT – updates a lesson type assignment by its ID with the data provided in the request body

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

### GET – gets all semester discipline relations from the database

200 OK [=> SemesterDiscipline[]](schemas.md#curriculum-semester-discipline)

### POST – adds a new semester discipline relation

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

### GET – finds a semester discipline relation by ID (URL parameter)

200 OK [=> SemesterDiscipline](schemas.md#curriculum-semester-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE – deletes a semester discipline relation by its ID (URL path)

200 OK [=> SemesterDiscipline](schemas.md#curriculum-semester-discipline)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

# Teacher Load Service

<a id="teacher-load-teachers"></a>

## /teachers

### GET - gets all teachers from the database

**200 OK** [=> Teacher[]](schemas.md#teacher-load-teacher)

<a id="teacher-load-group-cohorts"></a>

## /group-cohorts

### GET - gets all group cohorts from the database

200 OK [=> GroupCohort[]](schemas.md#teacher-load-group-cohort)

<a id="teacher-load-lesson-types"></a>

## /lesson-types

### GET - gets all lesson types from the database

200 OK [=> LessonType[]](schemas.md#teacher-load-lesson-type)

<a id="teacher-load-disciplines"></a>

## /disciplines

### GET – gets all disciplines from the database

200 OK [=> Discipline[]](schemas.md#teacher-load-discipline)

<a id="teacher-load-teacher-loads"></a>

## /teacher-loads

### GET – gets all teacher loads from the database

200 OK [=> TeacherLoad[]](schemas.md#teacher-load-teacher-load)

### POST – adds a new teacher load

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

<a id="teacher-load-teacher-load-id"></a>

## /teacher-load/{id}

### GET – finds a teacher load with an ID as an URL parameter

200 OK [=> TeacherLoad](schemas.md#teacher-load-teacher-load)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### DELETE – deletes a teacher load by its ID provided in the URL path

200 OK [=> TeacherLoad](schemas.md#teacher-load-teacher-load)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT – updates a teacher load by its ID with the data provided in the request body

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
