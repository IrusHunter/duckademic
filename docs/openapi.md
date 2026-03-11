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

# Employee Service

## /employee/academic-ranks

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
