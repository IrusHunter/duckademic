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
