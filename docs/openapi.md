# API Documentation

# ApiGateway

## /seed

### ANY

**Summary:** seed databases of all services.

#### Responses

**204 NO CONTENT**

If all done correctly.

**500 INTERNAL SERVER ERROR**

Something important failed.

```json
{
  "error": "detailed error description"
}
```

## /

### ANY

#### Available routes

- [/employee/academic-ranks](#employee-academic-ranks)
- [/employee/academic-ranks/{id}](#employee-academic-ranks-id)

**Summary:** routes and proxies incoming requests to the appropriate backend service based on the request path.

#### Responses

**400 BAD REQUEST**

No service found for the requested path.

```json
{
  "error": "detailed error description"
}
```

**500 INTERNAL SERVER ERROR**

Something important failed.

```json
{
  "error": "detailed error description"
}
```

# Employee Service

<a id="employee-academic-ranks"></a>
## /academic-ranks

### GET

**Summary:** gets all academic ranks from the database

#### Responses

**200 OK**

```json
[
  {
    "id": "uuid (identifier of the academic rank)",
    "slug": "string (unique slug used internally)",
    "title": "string (human-readable name of the rank)",
    "created_at": "timestamp (record creation timestamp)",
    "updated_at": "timestamp (record update timestamp)"
  }
]
```

### POST

**Summary:** adds a new academic rank

#### Request

```json
{
  "title": "string (human-readable name of the rank)"
}
```

#### Responses

**200 OK**

```json
{
  "id": "uuid (identifier of the academic rank)",
  "slug": "string (unique slug used internally)",
  "title": "string (human-readable name of the rank)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record update timestamp)"
}
```

**400 BAD REQUEST**

```json
{
  "error": "detailed error description"
}
```

<a id="employee-academic-ranks-id"></a>
## /academic-rank/{id}

### GET

**Summary:** finds academic rank with an ID as a URL parameter

#### Responses

**200 OK**

```json
{
  "id": "uuid (identifier of the academic rank)",
  "slug": "string (unique slug used internally)",
  "title": "string (human-readable name of the rank)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record update timestamp)"
}
```

**400 BAD REQUEST**

When the ID is invalid, or the academic rank is not found.

```json
{
  "error": "detailed error description"
}
```

### DELETE

**Summary:** deletes an academic rank by its ID provided in the URL path

#### Request

No body required; the ID should be provided as a URL parameter:

```
DELETE /academic-ranks/{id}
```

#### Responses

**204 NO CONTENT**

Returned when the academic rank is successfully deleted.

```json
null
```

**400 BAD REQUEST**

Returned when the ID is invalid or deletion fails.

```json
{
  "error": "detailed error description"
}
```

### PUT

**Summary:** updates an academic rank by its ID with the data provided in the request body

#### Request

```json
{
  "title": "string (human-readable name of the rank, optional)"
}
```

#### Responses

**200 OK**

Returned when the academic rank is successfully updated.

```json
{
  "id": "uuid (identifier of the academic rank)",
  "slug": "string (unique slug used internally)",
  "title": "string (human-readable name of the rank)",
  "created_at": "timestamp (record creation timestamp)",
  "updated_at": "timestamp (record update timestamp)"
}
```

**400 BAD REQUEST**

Returned when the ID is invalid, the request body is malformed, or the update fails.

```json
{
  "error": "detailed error description"
}
```

# Service name

## url path

### method

**Summary:** what operation do

#### Request

```json

```

#### Responses

**Code DESCRIPTION**

Condition when returned

```json

```
