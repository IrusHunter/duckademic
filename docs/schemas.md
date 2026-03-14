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
