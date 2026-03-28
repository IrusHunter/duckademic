# Schemas

### ErrorResponse

```json
{
  "error": "detailed error description"
}
```

## Employee Service

<a id="employee-academic-rank"></a>

### AcademicRank

```json
{
  "id": "uuid (identifier of the academic rank)",
  "slug": "string (unique slug used internally)",
  "title": "string (human-readable name of the rank)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record update timestamp)"
}
```

<a id="employee-academic-degree"></a>

### AcademicDegree

```json
{
  "id": "uuid (identifier of the academic degree)",
  "slug": "string (unique slug used internally)",
  "title": "string (human-readable name of the degree)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record update timestamp)"
}
```

<a id="employee-employee"></a>

### Employee

```json
{
  "id": "uuid (unique identifier of the employee)",
  "slug": "string (unique slug used internally)",
  "first_name": "string (employee's first name)",
  "last_name": "string (employee's last name)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",

  // Optional fields
  "middle_name": "string (employee's middle name)",
  "phone_number": "string (contact phone number)",
  "deleted_at": "timestamp (soft delete timestamp)"
}
```

<a id="employee-teacher"></a>

### Teacher

```json
{
  "employee_id": "uuid (reference to employee)",
  "email": "string (teacher's email address)",
  "academic_rank_id": "uuid (reference to academic rank)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",

  // Optional fields
  "academic_degree_id": "uuid (reference to academic degree)",
  "deleted_at": "timestamp (soft delete timestamp)"
}
```

## Schedule Service

<a id="schedule-academic-rank"></a>

### AcademicRank

```json
{
  "id": "uuid (identifier of the academic rank)",
  "slug": "string (unique slug used internally)",
  "title": "string (human-readable name of the rank)",
  "priority": "int (determines the rank's priority: higher value = higher rank)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record update timestamp)"
}
```

<a id="schedule-teacher"></a>

### Teacher

```json
{
  "id": "uuid (identifier of the teacher)",
  "name": "string (short full name of the teacher)",
  "academic_rank_id": "uuid (reference to academic rank)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",

  // Optional fields
  "deleted_at": "timestamp (soft delete timestamp)"
}
```

<a id="schedule-discipline"></a>

### Discipline

```json
{
  "id": "uuid (unique identifier of the discipline)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the discipline)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="schedule-lesson-type"></a>

### Lesson Type

```json
{
  "id": "uuid (unique identifier of the lesson type)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the lesson type)",
  "hours_value": "integer (number of hours per lesson)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="schedule-lesson-type-assignment"></a>

### Lesson Type Assignment

```json
{
  "id": "uuid (unique identifier of this lesson type assignment)",
  "lesson_type_id": "uuid (identifier of the associated lesson type)",
  "discipline_id": "uuid (identifier of the associated discipline)",
  "required_hours": "integer (number of hours required for this lesson type in this discipline)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="schedule-student"></a>

### Student

```json
{
  "id": "uuid (unique identifier of the student)",
  "slug": "string (unique slug used internally)",
  "name": "string (student's short full name)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",

  // Optional fields
  "deleted_at": "timestamp (soft delete timestamp)"
}
```

<a id="schedule-student-group"></a>

### Student Group

```json
{
  "id": "uuid (unique identifier of the student group)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the group)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="schedule-group-member"></a>

### Group Member

```json
{
  "id": "uuid (unique identifier of the group member record)",
  "studentId": "uuid (identifier of the student)",
  "createdAt": "timestamp (record creation timestamp)",
  "updatedAt": "timestamp (record last update timestamp)",

  // Optional fields
  "student_group_id": "uuid (identifier of the student group)"
}
```

<a id="schedule-teacher-load"></a>

### Teacher Load

```json
{
  "id": "uuid (unique identifier of the teacher load record)",
  "teacher_id": "uuid (unique identifier of the teacher)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)",
  "group_cohort_id": "uuid (unique identifier of the group cohort)",
  "group_count": "integer (number of groups assigned for this load)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

## Student Service

<a id="student-semester"></a>

### Semester

```json
{
  "id": "uuid (unique identifier of the semester)",
  "slug": "string (unique slug used internally)",
  "curriculum_id": "uuid (identifier of the associated curriculum)",
  "number": "integer (semester number within the curriculum)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="student-student"></a>

### Student

```json
{
  "id": "uuid (unique identifier of the student)",
  "slug": "string (unique slug used internally)",
  "first_name": "string (student's first name)",
  "last_name": "string (student's last name)",
  "email": "string (student's email address)",
  "semester_id": "uuid (identifier of the associated semester)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",

  // Optional fields
  "middle_name": "string (student's middle name)",
  "phone_number": "string (contact phone number)",
  "deleted_at": "timestamp (soft delete timestamp)"
}
```

## Student Group Service

<a id="student-group-semester"></a>

### Semester

```json
{
  "id": "uuid (unique identifier of the semester)",
  "slug": "string (unique slug used internally)",
  "curriculum_id": "uuid (identifier of the associated curriculum)",
  "number": "integer (semester number within the curriculum)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="student-group-student"></a>

### Student

```json
{
  "id": "uuid (unique identifier of the student)",
  "slug": "string (unique slug used internally)",
  "name": "string (student's short full name)",
  "semester_id": "uuid (identifier of the associated semester)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",

  // Optional fields
  "deleted_at": "timestamp (soft delete timestamp)"
}
```

<a id="student-group-group-cohort"></a>

### Group Cohort

```json
{
  "id": "uuid (unique identifier of the group cohort)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the cohort)",
  "semester_id": "uuid (identifier of the associated semester)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="student-group-student-group"></a>

### Student Group

```json
{
  "id": "uuid (unique identifier of the student group)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the group)",
  "group_cohort_id": "uuid (identifier of the associated student group)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="student-group-group-member"></a>

### Group Member

```json
{
  "id": "uuid (unique identifier of the group member record)",
  "studentId": "uuid (identifier of the student)",
  "group_cohort_id": "uuid (identifier of the associated group cohort)",
  "createdAt": "timestamp (record creation timestamp)",
  "updatedAt": "timestamp (record last update timestamp)",

  // Optional fields
  "student_group_id": "uuid (identifier of the student group)"
}
```

## Curriculum Service

<a id="curriculum-curriculum"></a>

### Curriculum

```json
{
  "id": "uuid (unique identifier of the curriculum)",
  "slug": "string (unique slug used internally)",
  "name": "string (curriculum name)",
  "duration_years": "integer (number of years for the curriculum)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",
  "effective_from": "timestamp (when this curriculum becomes effective)",

  // Optional fields
  "effective_to": "timestamp (when this curriculum ends, nullable)"
}
```

<a id="curriculum-semester"></a>

### Semester

```json
{
  "id": "uuid (unique identifier of the semester)",
  "slug": "string (unique slug used internally)",
  "curriculum_id": "uuid (identifier of the associated curriculum)",
  "number": "integer (semester number within the curriculum)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="curriculum-lesson-type"></a>

### Lesson Type

```json
{
  "id": "uuid (unique identifier of the lesson type)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the lesson type)",
  "hours_value": "integer (number of hours per lesson)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="curriculum-discipline"></a>

### Discipline

```json
{
  "id": "uuid (unique identifier of the discipline)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the discipline)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="curriculum-lesson-type-assignment"></a>

### Lesson Type Assignment

```json
{
  "id": "uuid (unique identifier of this lesson type assignment)",
  "lesson_type_id": "uuid (identifier of the associated lesson type)",
  "discipline_id": "uuid (identifier of the associated discipline)",
  "required_hours": "integer (number of hours required for this lesson type in this discipline)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="curriculum-semester-discipline"></a>

### Semester Discipline Relation

```json
{
  "id": "uuid (unique identifier of this semester discipline record)",
  "semester_id": "uuid (identifier of the associated semester)",
  "discipline_id": "uuid (identifier of the associated discipline)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

## Teacher Load Service

<a id="teacher-load-teacher"></a>

### Teacher

```json
{
  "id": "uuid (identifier of the teacher)",
  "name": "string (short full name of the teacher)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="teacher-load-group-cohort"></a>

### Group Cohort

```json
{
  "id": "uuid (unique identifier of the group cohort)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the cohort)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="teacher-load-discipline"></a>

### Discipline

```json
{
  "id": "uuid (unique identifier of the discipline)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the discipline)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="teacher-load-lesson-type"></a>

### Lesson Type

```json
{
  "id": "uuid (unique identifier of the lesson type)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the lesson type)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="teacher-load-teacher-load"></a>

### Teacher Load

```json
{
  "id": "uuid (unique identifier of the teacher load record)",
  "teacher_id": "uuid (unique identifier of the teacher)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)",
  "group_cohort_id": "uuid (unique identifier of the group cohort)",
  "group_count": "integer (number of groups assigned for this load)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```
