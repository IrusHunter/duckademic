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
  "slug": "string (unique slug used internally)",
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
  "reserved_weeks": "string (comma-separated week numbers where only this lesson type is allowed)",
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

<a id="schedule-group-cohort"></a>

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

<a id="schedule-group-cohort-assignment"></a>

### Group Cohort Assignment

```json
{
  "id": "uuid (unique identifier of the assignment)",
  "group_cohort_id": "uuid (unique identifier of the group cohort)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="schedule-classroom"></a>

### Classroom

```json
{
  "id": "uuid (unique identifier of the classroom)",
  "slug": "string (unique slug used internally)",
  "number": "string (classroom number or label)",
  "capacity": "integer (maximum number of visitors in the classroom)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="schedule-study-load"></a>

### Study Load

```json
{
  "id": "uuid (unique identifier of the study load)",
  "teacher_id": "uuid (unique identifier of the teacher)",
  "student_group_id": "uuid (unique identifier of the student group)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="schedule-study-lesson-slot"></a>

### Lesson Slot

```json
{
  "id": "uuid (unique identifier of the lesson slot)",
  "slot": "integer (index of the slot within a day, starting from 0)",
  "weekday": "integer (day of the week, 0 = Sunday, 6 = Saturday)",
  "start_time": "duration (lesson start time as duration since midnight, in nanoseconds)",
  "duration": "duration (length of the lesson, in nanoseconds)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="schedule-generator-lesson-occurrense"></a>

### Lesson Occurrence

```json
{
  "id": "uuid (unique identifier of the lesson occurrence)",
  "study_load_id": "uuid (reference to study load)",
  "teacher_id": "uuid (unique identifier of the teacher)",
  "student_group_id": "uuid (unique identifier of the student group)",
  "lesson_slot_id": "uuid (unique identifier of the lesson slot)",
  "date": "timestamp (date and time of the lesson occurrence)",
  "status": "string (lesson status: scheduled, canceled, completed)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",

  // Optional fields
  "classroom_id": "uuid | null (optional classroom assignment)",
  "moved_to_id": "uuid | null (reference to rescheduled occurrence, if moved)",
  "moved_from_id": "uuid | null (reference to original occurrence if moved)"
}
```

## Schedule Generator Service

<a id="schedule-generator-generator-config"></a>

### Generator Config

```json
{
  "start_date": "timestamp (start date and time of the study period)",
  "end_date": "timestamp (end date and time of the study period, inclusive of the last day)",
  "slot_preference": [
    [
      "float (preference coefficient for a time slot; ordered by days starting from Sunday: Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday)"
    ]
  ],
  "max_daily_student_load": "integer (maximum number of classes per student per day)",
  "lesson_fill_rate": "float (percentage of lesson type utilization used to determine the number of study days)",
  "classroom_occupancy": "float (percentage of classroom occupancy)"
}
```

<a id="schedule-generator-teacher"></a>

### Teacher

```json
{
  "id": "uuid (identifier of the teacher)",
  "name": "string (short full name of the teacher)",
  "priority": "int (determines the rank's priority: higher value = higher rank)"
}
```

<a id="schedule-generator-discipline"></a>

### Discipline

```json
{
  "id": "uuid (unique identifier of the discipline)",
  "name": "string (name of the discipline)"
}
```

<a id="schedule-generator-lesson-type"></a>

### Lesson Type

```json
{
  "id": "uuid (unique identifier of the lesson type)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the lesson type)",
  "hours_value": "integer (number of hours per lesson)",
  "reserved_weeks": "string (comma-separated week numbers where only this lesson type is allowed)"
}
```

<a id="schedule-generator-lesson-type-assignment"></a>

### Lesson Type Assignment

```json
{
  "id": "uuid (unique identifier of this lesson type assignment)",
  "lesson_type_id": "uuid (identifier of the associated lesson type)",
  "discipline_id": "uuid (identifier of the associated discipline)",
  "required_hours": "integer (number of hours required for this lesson type in this discipline)"
}
```

<a id="schedule-generator-group-cohort"></a>

### Group Cohort

```json
{
  "id": "uuid (unique identifier of the group cohort)",
  "slug": "string (unique slug used internally)",
  "name": "string (name of the cohort)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)",
  "groups": [
    {
      "id": "uuid (unique identifier of the student group)",
      "name": "string (name of the student group)",
      "connected_groups": ["uuid (ID of connected student group)"],
      "student_count": "int (number of students in this group)"
    }
  ]
}
```

<a id="schedule-generator-group-cohort-assignment"></a>

### Group Cohort Assignment

```json
{
  "id": "uuid (unique identifier of the assignment)",
  "group_cohort_id": "uuid (unique identifier of the group cohort)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)"
}
```

<a id="schedule-generator-teacher-load"></a>

### Teacher Load

```json
{
  "id": "uuid (unique identifier of the teacher load record)",
  "teacher_id": "uuid (unique identifier of the teacher)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)",
  "group_count": "integer (number of groups assigned for this load)"
}
```

<a id="schedule-generator-classroom"></a>

### Classroom

```json
{
  "id": "uuid (unique identifier of the classroom)",
  "number": "string (classroom number or label)",
  "capacity": "integer (maximum number of visitors in the classroom)"
}
```

<a id="schedule-generator-days-for-lesson-types"></a>

### Days for Lesson Types

```json
{
  "student_groups": [
    {
      "id": "uuid (unique identifier of the student group)",
      "name": "string (name of the student group)",
      "weekday_lesson_types": [
        {
          "id": "uuid (unique identifier of the lesson type)",
          "name": "string (name of the lesson type)",
          "weekday": "integer (weekday number, e.g., 1 = Monday)"
        }
      ]
    }
  ],
  "errors": [
    {
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "slots_dept": "number (float value representing slots debt)"
    }
  ]
}
```

<a id="schedule-generator-bone-lessons"></a>

### Bone Lessons

```json
{
  "bone_lessons": [
    {
      "teacher": {
        "id": "uuid (unique identifier of the teacher)",
        "name": "string (name of the teacher)"
      },
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "discipline": {
        "id": "uuid (unique identifier of the discipline)",
        "name": "string (name of the discipline)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "day": "integer (weekday number, e.g., 1 = Monday)",
      "slot": "integer (lesson slot number in the day)"
    }
  ],
  "errors": [
    {
      "teacher": {
        "id": "uuid (unique identifier of the teacher)",
        "name": "string (name of the teacher)"
      },
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "discipline": {
        "id": "uuid (unique identifier of the discipline)",
        "name": "string (name of the discipline)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "count": "integer (number of unassigned lessons)"
    }
  ]
}
```

<a id="schedule-generator-bone-lessons-with-c"></a>

### Bone Lessons With C

```json
{
  "lessons_with_classroom": [
    {
      "teacher": {
        "id": "uuid (unique identifier of the teacher)",
        "name": "string (name of the teacher)"
      },
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "discipline": {
        "id": "uuid (unique identifier of the discipline)",
        "name": "string (name of the discipline)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "classroom": {
        "id": "uuid (unique identifier of the classroom)",
        "name": "string (name/number of the classroom)"
      },
      "day": "integer (weekday number, e.g., 1 = Monday)",
      "slot": "integer (lesson slot number in the day)"
    }
  ],
  "lessons_without_classroom": [
    {
      "teacher": {
        "id": "uuid (unique identifier of the teacher)",
        "name": "string (name of the teacher)"
      },
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "discipline": {
        "id": "uuid (unique identifier of the discipline)",
        "name": "string (name of the discipline)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "day": "integer (weekday number, e.g., 1 = Monday)",
      "slot": "integer (lesson slot number in the day)"
    }
  ]
}
```

<a id="schedule-generator-generated-lessons"></a>

### Generated Lessons

```json
{
  "lessons": [
    {
      "teacher": {
        "id": "uuid (unique identifier of the teacher)",
        "name": "string (name of the teacher)"
      },
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "discipline": {
        "id": "uuid (unique identifier of the discipline)",
        "name": "string (name of the discipline)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "days": ["integer (day number of similar lessons at same weekday)"],
      "slot": "integer (lesson slot number in the day)",

      // Optional fields
      "classroom": {
        "id": "uuid (unique identifier of the classroom)",
        "name": "string (name of the classroom)"
      }
    }
  ],
  "errors": [
    {
      "teacher": {
        "id": "uuid (unique identifier of the teacher)",
        "name": "string (name of the teacher)"
      },
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "discipline": {
        "id": "uuid (unique identifier of the discipline)",
        "name": "string (name of the discipline)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "count": "integer (number of unassigned lessons)"
    }
  ]
}
```

<a id="schedule-generator-generated-lessons-with-c"></a>

### Generated Lessons With C

```json
{
  "lessons_with_classroom": [
    {
      "teacher": {
        "id": "uuid (unique identifier of the teacher)",
        "name": "string (name of the teacher)"
      },
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "discipline": {
        "id": "uuid (unique identifier of the discipline)",
        "name": "string (name of the discipline)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "classroom": {
        "id": "uuid (unique identifier of the classroom)",
        "name": "string (name/number of the classroom)"
      },
      "days": ["integer (day number of similar lessons at same weekday)"],
      "slot": "integer (lesson slot number in the day)"
    }
  ],
  "lessons_without_classroom": [
    {
      "teacher": {
        "id": "uuid (unique identifier of the teacher)",
        "name": "string (name of the teacher)"
      },
      "student_group": {
        "id": "uuid (unique identifier of the student group)",
        "name": "string (name of the student group)"
      },
      "discipline": {
        "id": "uuid (unique identifier of the discipline)",
        "name": "string (name of the discipline)"
      },
      "lesson_type": {
        "id": "uuid (unique identifier of the lesson type)",
        "name": "string (name of the lesson type)"
      },
      "day": "integer (weekday number, e.g., 1 = Monday)",
      "slot": "integer (lesson slot number in the day)"
    }
  ]
}
```

<a id="schedule-generator-study-load"></a>

### Study Load

```json
{
  "id": "uuid (unique identifier of the study load)",
  "teacher_id": "uuid (unique identifier of the teacher)",
  "student_group_id": "uuid (unique identifier of the student group)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)"
}
```

<a id="schedule-generator-lesson"></a>

### Lesson

```json
{
  "id": "uuid (unique identifier of the lesson)",
  "study_load_id": "uuid (unique identifier of the study load)",
  "teacher_id": "uuid (unique identifier of the teacher)",
  "student_group_id": "uuid (unique identifier of the student group)",
  "slot": "integer (lesson slot number in the day)",
  "day": "integer (number of days after start)",

  // Optional field
  "classroom_id": "uuid (unique identifier of the classroom)"
}
```

<a id="schedule-generator-fault"></a>

### Fault

```json
{
  "total_value": "double (sum of the generator's fault values)"
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

<a id="student-group-discipline"></a>

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

<a id="student-group-lesson-type"></a>

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

<a id="student-group-group-cohort-assignment"></a>

### Group Cohort Assignment

```json
{
  "id": "uuid (unique identifier of the assignment)",
  "group_cohort_id": "uuid (unique identifier of the group cohort)",
  "discipline_id": "uuid (unique identifier of the discipline)",
  "lesson_type_id": "uuid (unique identifier of the lesson type)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
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
  "slug": "string (unique slug used internally)",
  "name": "string (short full name of the teacher)",
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
  "group_count": "integer (number of groups assigned for this load)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

## Asset Service

<a id="asset-classroom"></a>

### Classroom

```json
{
  "id": "uuid (unique identifier of the classroom)",
  "slug": "string (unique slug used internally)",
  "number": "string (classroom number or label)",
  "capacity": "integer (maximum number of visitors in the classroom)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

## Auth Service

<a id="auth-permission"></a>

### Permission

```json
{
  "id": "uuid (unique identifier of the classroom)",
  "name": "string (unique name of the permission)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="auth-role"></a>

### Role

```json
{
  "id": "uuid (unique identifier of the role)",
  "name": "string (unique name of the role)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="auth-service"></a>

### Service

```json
{
  "id": "uuid (unique identifier of the service)",
  "name": "string (unique name of the service)",
  "secrete": "string (service secret key)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="auth-role-permissions"></a>

### Role Permissions

```json
{
  "id": "uuid (unique identifier of the role permission link)",
  "role_id": "uuid (associated role identifier)",
  "permission_id": "uuid (associated permission identifier)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```

<a id="auth-service-permissions"></a>

### Service Permissions

```json
{
  "id": "uuid (unique identifier of the service permission link)",
  "service_id": "uuid (associated service identifier)",
  "permission_id": "uuid (associated permission identifier)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record last update timestamp)"
}
```
