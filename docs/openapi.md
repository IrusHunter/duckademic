# API Documentation

# ApiGateway

## /seed

### ANY - seeds databases of all services.

204 NO CONTENT

500 INTERNAL SERVER ERROR [=> ErrorResponse](schemas.md#errorresponse)

## /

### ANY - routes and proxies incoming requests to the appropriate backend service based on the request path.

- Employee Service
  - [/employee/academic-ranks](#employee-academic-ranks)
  - [/employee/academic-rank/{id}](#employee-academic-rank-id)
  - [/employee/academic-degrees](#employee-academic-degrees)
  - [/employee/academic-degree/{id}](#employee-academic-degree-id)
  - [/employee/employees](#employee-employees)
  - [/employee/employee/{id}](#employee-employee-id)
  - [/employee/teachers](#employee-teachers)
  - [/employee/teacher/{id}](#employee-teacher-id)

- Schedule Service
  - [/schedule/academic-ranks](#schedule-academic-ranks)
  - [/schedule/academic-rank/{id}](#schedule-academic-rank-id)
  - [/schedule/teachers](#schedule-teachers)
  - [/schedule/teacher/{id}](#schedule-teacher-id)

- Student Service
  - [/student/students](#student-students)
  - [/student/student/{id}](#student-student-id)

- Student Group Service
  - [/student-group/students](#student-group-students)
  - [/student-group/student/{id}](#student-group-student-id)

- Curriculum Service
  - [/curriculum/curriculums](#curriculum-curriculum)
  - [/curriculum/curriculum/{id}](#curriculum-curriculum-id)
  - [/curriculum/semesters](#curriculum-semesters)
  - [/curriculum/semester/{id}](#curriculum-semester-id)

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

**200 OK** [=> AcademicRank](schemas.md#employee-academic-rank)

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

<a id="schedule-teacher-id"></a>

## /teacher/{id}

### GET - finds teacher with an ID as a URL parameter

**200 OK** [=> Teacher](schemas.md#schedule-teacher)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a teacher by its ID with the data provided in the request body

```json
{}
```

**200 OK** [=> Teacher](schemas.md#schedule-teacher)

**400 BAD REQUEST** [=> ErrorResponse](schemas.md#errorresponse)

# Student Service

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

  // Optional fields
  "middle_name": "string (student's middle name)",
  "phone_number": "string (contact phone number)"
}
```

200 OK [=> Student](schemas.md#student-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

# Student Group Service

<a id="student-group-students"></a>

## /students

### GET - gets all students from the database

200 OK [=> Student[]](schemas.md#student-group-student)

<a id="student-group-student-id"></a>

## /student/{id}

### GET - finds student with an ID as an URL parameter

200 OK [=> Student](schemas.md#student-group-student)

400 BAD REQUEST [=> ErrorResponse](schemas.md#errorresponse)

### PUT - updates a student by its ID with the data provided in the request body

```json
{}
```

200 OK [=> Student](schemas.md#student-group-student)

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

---

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
