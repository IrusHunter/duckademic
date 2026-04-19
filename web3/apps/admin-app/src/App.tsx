import { useState } from 'react'
import { useNavigate, useParams, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider, useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'

const queryClient = new QueryClient()

// ─── ТИПИ ────────────────────────────────────────────────────────────────────

type FieldDef = {
  key: string
  label: string
  required?: boolean
}

type TableDef = {
  key: string               // унікальний ключ таблиці
  label: string             // назва в UI
  listEndpoint: string      // GET всіх записів, напр. '/employees'
  itemEndpoint: string      // DELETE/PUT одного запису, напр. '/employee'  (без /{id})
  fields: FieldDef[]        // поля для форми додавання
  columns: FieldDef[]       // колонки таблиці (включаючи id)
  readOnly?: boolean        // якщо true — тільки перегляд, без Add/Delete
}

type ServiceDef = {
  key: string
  label: string
  icon: string
  baseURL: string           // /api/employee, /api/student тощо
  tables: TableDef[]
}

// ─── КОНФІГ СЕРВІСІВ ─────────────────────────────────────────────────────────

const SERVICES: ServiceDef[] = [
  {
    key: 'employee',
    label: 'Employee Service',
    icon: '👤',
    baseURL: '/api/employee',
    tables: [
      {
        key: 'employees',
        label: 'Employees',
        listEndpoint: '/employees',
        itemEndpoint: '/employee',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'first_name', label: 'First Name' },
          { key: 'last_name', label: 'Last Name' },
          { key: 'middle_name', label: 'Middle Name' },
          { key: 'phone_number', label: 'Phone' },
        ],
        fields: [
          { key: 'first_name', label: 'First Name', required: true },
          { key: 'last_name', label: 'Last Name', required: true },
          { key: 'middle_name', label: 'Middle Name' },
          { key: 'phone_number', label: 'Phone Number' },
        ],
      },
      {
        key: 'teachers',
        label: 'Teachers',
        listEndpoint: '/teachers',
        itemEndpoint: '/teacher',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'employee_id', label: 'Employee ID' },
          { key: 'email', label: 'Email' },
          { key: 'academic_rank_id', label: 'Academic Rank ID' },
          { key: 'academic_degree_id', label: 'Academic Degree ID' },
        ],
        fields: [
          { key: 'employee_id', label: 'Employee ID', required: true },
          { key: 'email', label: 'Email', required: true },
          { key: 'academic_rank_id', label: 'Academic Rank ID', required: true },
          { key: 'academic_degree_id', label: 'Academic Degree ID' },
        ],
      },
      {
        key: 'academic-ranks',
        label: 'Academic Ranks',
        listEndpoint: '/academic-ranks',
        itemEndpoint: '/academic-rank',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'title', label: 'Title' },
        ],
        fields: [
          { key: 'title', label: 'Title', required: true },
        ],
      },
      {
        key: 'academic-degrees',
        label: 'Academic Degrees',
        listEndpoint: '/academic-degrees',
        itemEndpoint: '/academic-degree',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'title', label: 'Title' },
        ],
        fields: [
          { key: 'title', label: 'Title', required: true },
        ],
      },
    ],
  },
  {
    key: 'student',
    label: 'Student Service',
    icon: '🎓',
    baseURL: '/api/student',
    tables: [
      {
        key: 'students',
        label: 'Students',
        listEndpoint: '/students',
        itemEndpoint: '/student',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'first_name', label: 'First Name' },
          { key: 'last_name', label: 'Last Name' },
          { key: 'middle_name', label: 'Middle Name' },
          { key: 'email', label: 'Email' },
          { key: 'phone_number', label: 'Phone' },
          { key: 'semester_id', label: 'Semester ID' },
        ],
        fields: [
          { key: 'first_name', label: 'First Name', required: true },
          { key: 'last_name', label: 'Last Name', required: true },
          { key: 'email', label: 'Email', required: true },
          { key: 'semester_id', label: 'Semester ID', required: true },
          { key: 'middle_name', label: 'Middle Name' },
          { key: 'phone_number', label: 'Phone Number' },
        ],
      },
      {
        key: 'semesters',
        label: 'Semesters',
        listEndpoint: '/semesters',
        itemEndpoint: '/semester',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'number', label: 'Number' },
          { key: 'curriculum_id', label: 'Curriculum ID' },
        ],
        fields: [],
      },
    ],
  },
  {
    key: 'student-group',
    label: 'Student Group Service',
    icon: '👥',
    baseURL: '/api/student-group',
    tables: [
      {
        key: 'group-cohorts',
        label: 'Group Cohorts',
        listEndpoint: '/group-cohorts',
        itemEndpoint: '/group-cohort',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'name', label: 'Name' },
          { key: 'semester_id', label: 'Semester ID' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
          { key: 'semester_id', label: 'Semester ID', required: true },
        ],
      },
      {
        key: 'student-groups',
        label: 'Student Groups',
        listEndpoint: '/student-groups',
        itemEndpoint: '/student-group',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'name', label: 'Name' },
          { key: 'group_cohort_id', label: 'Group Cohort ID' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
          { key: 'group_cohort_id', label: 'Group Cohort ID', required: true },
        ],
      },
      {
        key: 'group-members',
        label: 'Group Members',
        listEndpoint: '/group-members',
        itemEndpoint: '/group-members',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'studentId', label: 'Student ID' },
          { key: 'group_cohort_id', label: 'Group Cohort ID' },
          { key: 'student_group_id', label: 'Student Group ID' },
        ],
        fields: [
          { key: 'studentId', label: 'Student ID', required: true },
          { key: 'group_cohort_id', label: 'Group Cohort ID', required: true },
          { key: 'student_group_id', label: 'Student Group ID' },
        ],
      },
      {
        key: 'group-cohort-assignments',
        label: 'Cohort Assignments',
        listEndpoint: '/group-cohort-assignments',
        itemEndpoint: '/group-cohort-assignments',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'group_cohort_id', label: 'Group Cohort ID' },
          { key: 'discipline_id', label: 'Discipline ID' },
          { key: 'lesson_type_id', label: 'Lesson Type ID' },
        ],
        fields: [
          { key: 'group_cohort_id', label: 'Group Cohort ID', required: true },
          { key: 'discipline_id', label: 'Discipline ID', required: true },
          { key: 'lesson_type_id', label: 'Lesson Type ID', required: true },
        ],
      },
      {
        key: 'semesters',
        label: 'Semesters',
        listEndpoint: '/semesters',
        itemEndpoint: '/semester',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'number', label: 'Number' }],
        fields: [],
      },
      {
        key: 'students',
        label: 'Students (read)',
        listEndpoint: '/students',
        itemEndpoint: '/student',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'first_name', label: 'First Name' }, { key: 'last_name', label: 'Last Name' }],
        fields: [],
      },
      {
        key: 'lesson-types',
        label: 'Lesson Types (read)',
        listEndpoint: '/lesson-types',
        itemEndpoint: '/lesson-type',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'name', label: 'Name' }],
        fields: [],
      },
      {
        key: 'disciplines',
        label: 'Disciplines (read)',
        listEndpoint: '/disciplines',
        itemEndpoint: '/discipline',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'name', label: 'Name' }],
        fields: [],
      },
    ],
  },
  {
    key: 'curriculum',
    label: 'Curriculum Service',
    icon: '📚',
    baseURL: '/api/curriculum',
    tables: [
      {
        key: 'curriculums',
        label: 'Curriculums',
        listEndpoint: '/curriculums',
        itemEndpoint: '/curriculum',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'name', label: 'Name' },
          { key: 'duration_years', label: 'Duration (years)' },
          { key: 'effective_from', label: 'Effective From' },
          { key: 'effective_to', label: 'Effective To' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
          { key: 'duration_years', label: 'Duration (years)', required: true },
          { key: 'effective_from', label: 'Effective From', required: true },
          { key: 'effective_to', label: 'Effective To' },
        ],
      },
      {
        key: 'semesters',
        label: 'Semesters',
        listEndpoint: '/semesters',
        itemEndpoint: '/semester',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'curriculum_id', label: 'Curriculum ID' },
          { key: 'number', label: 'Number' },
        ],
        fields: [
          { key: 'curriculum_id', label: 'Curriculum ID', required: true },
          { key: 'number', label: 'Number', required: true },
        ],
      },
      {
        key: 'lesson-types',
        label: 'Lesson Types',
        listEndpoint: '/lesson-types',
        itemEndpoint: '/lesson-type',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'name', label: 'Name' },
          { key: 'hours_value', label: 'Hours' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
          { key: 'hours_value', label: 'Hours per lesson', required: true },
        ],
      },
      {
        key: 'disciplines',
        label: 'Disciplines',
        listEndpoint: '/disciplines',
        itemEndpoint: '/discipline',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'name', label: 'Name' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
        ],
      },
      {
        key: 'lesson-type-assignments',
        label: 'Lesson Type Assignments',
        listEndpoint: '/lesson-type-assignments',
        itemEndpoint: '/lesson-type-assignment',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'lesson_type_id', label: 'Lesson Type ID' },
          { key: 'discipline_id', label: 'Discipline ID' },
          { key: 'required_hours', label: 'Required Hours' },
        ],
        fields: [
          { key: 'lesson_type_id', label: 'Lesson Type ID', required: true },
          { key: 'discipline_id', label: 'Discipline ID', required: true },
          { key: 'required_hours', label: 'Required Hours', required: true },
        ],
      },
      {
        key: 'semester-disciplines',
        label: 'Semester Disciplines',
        listEndpoint: '/semester-disciplines',
        itemEndpoint: '/semester-discipline',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'semester_id', label: 'Semester ID' },
          { key: 'discipline_id', label: 'Discipline ID' },
        ],
        fields: [
          { key: 'semester_id', label: 'Semester ID', required: true },
          { key: 'discipline_id', label: 'Discipline ID', required: true },
        ],
      },
    ],
  },
  {
    key: 'teacher-load',
    label: 'Teacher Load Service',
    icon: '📋',
    baseURL: '/api/teacher-load',
    tables: [
      {
        key: 'teacher-loads',
        label: 'Teacher Loads',
        listEndpoint: '/teacher-loads',
        itemEndpoint: '/teacher-load',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'teacher_id', label: 'Teacher ID' },
          { key: 'discipline_id', label: 'Discipline ID' },
          { key: 'lesson_type_id', label: 'Lesson Type ID' },
          { key: 'group_count', label: 'Group Count' },
        ],
        fields: [
          { key: 'teacher_id', label: 'Teacher ID', required: true },
          { key: 'discipline_id', label: 'Discipline ID', required: true },
          { key: 'lesson_type_id', label: 'Lesson Type ID', required: true },
          { key: 'group_count', label: 'Group Count', required: true },
        ],
      },
      {
        key: 'teachers',
        label: 'Teachers (read)',
        listEndpoint: '/teachers',
        itemEndpoint: '/teacher',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'email', label: 'Email' }],
        fields: [],
      },
      {
        key: 'disciplines',
        label: 'Disciplines (read)',
        listEndpoint: '/disciplines',
        itemEndpoint: '/discipline',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'name', label: 'Name' }],
        fields: [],
      },
      {
        key: 'lesson-types',
        label: 'Lesson Types (read)',
        listEndpoint: '/lesson-types',
        itemEndpoint: '/lesson-type',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'name', label: 'Name' }],
        fields: [],
      },
    ],
  },
  {
    key: 'asset',
    label: 'Asset Service',
    icon: '🏫',
    baseURL: '/api/asset',
    tables: [
      {
        key: 'classrooms',
        label: 'Classrooms',
        listEndpoint: '/classrooms',
        itemEndpoint: '/classroom',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'number', label: 'Number' },
          { key: 'capacity', label: 'Capacity' },
        ],
        fields: [
          { key: 'number', label: 'Classroom Number', required: true },
          { key: 'capacity', label: 'Capacity', required: true },
        ],
      },
    ],
  },
  {
    key: 'schedule',
    label: 'Schedule Service',
    icon: '📅',
    baseURL: '/api/schedule',
    tables: [
      {
        key: 'academic-ranks',
        label: 'Academic Ranks (read)',
        listEndpoint: '/academic-ranks',
        itemEndpoint: '/academic-rank',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'title', label: 'Title' }, { key: 'priority', label: 'Priority' }],
        fields: [],
      },
      {
        key: 'lesson-types',
        label: 'Lesson Types (read)',
        listEndpoint: '/lesson-types',
        itemEndpoint: '/lesson-type',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'name', label: 'Name' }, { key: 'reserved_weeks', label: 'Reserved Weeks' }],
        fields: [],
      },
      {
        key: 'teachers',
        label: 'Teachers (read)',
        listEndpoint: '/teachers',
        itemEndpoint: '/teacher',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'email', label: 'Email' }],
        fields: [],
      },
      {
        key: 'disciplines',
        label: 'Disciplines (read)',
        listEndpoint: '/disciplines',
        itemEndpoint: '/discipline',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'name', label: 'Name' }],
        fields: [],
      },
      {
        key: 'students',
        label: 'Students (read)',
        listEndpoint: '/students',
        itemEndpoint: '/student',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'first_name', label: 'First Name' }, { key: 'last_name', label: 'Last Name' }],
        fields: [],
      },
      {
        key: 'student-groups',
        label: 'Student Groups (read)',
        listEndpoint: '/student-groups',
        itemEndpoint: '/student-group',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'name', label: 'Name' }],
        fields: [],
      },
      {
        key: 'classrooms',
        label: 'Classrooms (read)',
        listEndpoint: '/classrooms',
        itemEndpoint: '/classroom',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'number', label: 'Number' }, { key: 'capacity', label: 'Capacity' }],
        fields: [],
      },
      {
        key: 'lesson-occurrences',
        label: 'Lesson Occurrences (read)',
        listEndpoint: '/lesson-occurrences',
        itemEndpoint: '/lesson-occurrence',
        readOnly: true,
        columns: [{ key: 'id', label: 'ID' }, { key: 'lesson_slot_id', label: 'Slot ID' }, { key: 'classroom_id', label: 'Classroom ID' }],
        fields: [],
      },
    ],
  },
]

// ─── API ХЕЛПЕР ──────────────────────────────────────────────────────────────

function makeApi(baseURL: string) {
  return axios.create({ baseURL, withCredentials: true })
}

// ─── КОМПОНЕНТИ ──────────────────────────────────────────────────────────────

function AddForm({ fields, onSubmit }: {
  fields: FieldDef[]
  onSubmit: (data: Record<string, string>) => void
}) {
  const [values, setValues] = useState<Record<string, string>>({})
  const [open, setOpen] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(values)
    setValues({})
    setOpen(false)
  }

  return (
    <div style={{ marginBottom: 20 }}>
      <button
        onClick={() => setOpen(!open)}
        style={{ padding: '6px 16px', background: '#4A6CF7', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer', marginBottom: 10 }}
      >
        {open ? '✕ Cancel' : '+ Add'}
      </button>
      {open && (
        <form onSubmit={handleSubmit} style={{ display: 'flex', gap: 12, flexWrap: 'wrap', padding: 16, background: '#f9f9f9', borderRadius: 8 }}>
          {fields.map(field => (
            <div key={field.key} style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ fontSize: 12, color: '#666' }}>{field.label}{field.required && ' *'}</label>
              <input
                value={values[field.key] ?? ''}
                onChange={e => setValues(prev => ({ ...prev, [field.key]: e.target.value }))}
                required={field.required}
                style={{ padding: '6px 8px', borderRadius: 4, border: '1px solid #ccc', width: 180 }}
              />
            </div>
          ))}
          <div style={{ display: 'flex', alignItems: 'flex-end' }}>
            <button type="submit" style={{ padding: '6px 16px', background: '#333', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>Save</button>
          </div>
        </form>
      )}
    </div>
  )
}

function DataTable({ data, columns, onDelete, readOnly }: {
  data: Record<string, unknown>[]
  columns: FieldDef[]
  onDelete: (id: string) => void
  readOnly?: boolean
}) {
  const [selected, setSelected] = useState<string[]>([])
  const [action, setAction] = useState('')

  const toggleSelect = (id: string) =>
    setSelected(prev => prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id])

  const handleGo = () => {
    if (action === 'delete') {
      selected.forEach(id => onDelete(id))
      setSelected([])
    }
  }

  return (
    <div>
      {!readOnly && (
        <div style={{ marginBottom: 12, display: 'flex', gap: 8, alignItems: 'center' }}>
          <span style={{ fontSize: 14 }}>Action:</span>
          <select value={action} onChange={e => setAction(e.target.value)} style={{ padding: '4px 8px', borderRadius: 4, border: '1px solid #ccc' }}>
            <option value="">---</option>
            <option value="delete">Delete selected</option>
          </select>
          <button onClick={handleGo} style={{ padding: '4px 12px', borderRadius: 4, border: '1px solid #ccc', cursor: 'pointer' }}>Go</button>
        </div>
      )}
      <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
        <thead>
          <tr style={{ background: '#f5f5f5' }}>
            {!readOnly && <th style={{ padding: '8px 12px', borderBottom: '2px solid #e0e0e0', width: 40 }} />}
            {columns.map(col => (
              <th key={col.key} style={{ padding: '8px 12px', textAlign: 'left', borderBottom: '2px solid #e0e0e0', color: '#555' }}>{col.label}</th>
            ))}
            {!readOnly && <th style={{ padding: '8px 12px', borderBottom: '2px solid #e0e0e0', color: '#555' }}>Actions</th>}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 && (
            <tr><td colSpan={columns.length + (readOnly ? 0 : 2)} style={{ padding: 24, textAlign: 'center', color: '#aaa' }}>No data</td></tr>
          )}
          {data.map(row => {
            const id = String(row.id ?? '')
            return (
              <tr key={id} style={{ background: selected.includes(id) ? '#e8f0fe' : 'white' }}>
                {!readOnly && (
                  <td style={{ padding: '8px 12px', borderBottom: '1px solid #eee' }}>
                    <input type="checkbox" checked={selected.includes(id)} onChange={() => toggleSelect(id)} />
                  </td>
                )}
                {columns.map(col => (
                  <td key={col.key} style={{ padding: '8px 12px', borderBottom: '1px solid #eee', maxWidth: 220, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {String(row[col.key] ?? '')}
                  </td>
                ))}
                {!readOnly && (
                  <td style={{ padding: '8px 12px', borderBottom: '1px solid #eee' }}>
                    <button
                      onClick={() => onDelete(id)}
                      style={{ padding: '3px 10px', color: 'red', border: '1px solid red', borderRadius: 4, background: 'white', cursor: 'pointer', fontSize: 13 }}
                    >
                      Delete
                    </button>
                  </td>
                )}
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}

// ─── ДИНАМІЧНА СТОРІНКА ТАБЛИЦІ ──────────────────────────────────────────────

function TablePage({ service, table }: { service: ServiceDef; table: TableDef }) {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const api = makeApi(service.baseURL)
  const queryKey = [service.key, table.key]

  const { data: rawData, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
  const r = await api.get(table.listEndpoint)
  const result = r.data
  // якщо масив — повертаємо як є
  if (Array.isArray(result)) return result
  // якщо об'єкт з полем data/items/results — беремо його
  if (result && Array.isArray(result.data)) return result.data
  if (result && Array.isArray(result.items)) return result.items
  if (result && Array.isArray(result.results)) return result.results
  // якщо один об'єкт — обгортаємо в масив
  if (result && typeof result === 'object') return [result]
  return []
},
  })
  const data = Array.isArray(rawData) ? rawData : []

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.delete(`${table.itemEndpoint}/${id}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey }),
    onError: (e) => console.error('Delete error:', e),
  })

  const createMutation = useMutation({
    mutationFn: (body: Record<string, string>) => api.post(table.listEndpoint, body).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey }),
    onError: (e) => console.error('Create error:', e),
  })

  return (
    <div style={{ padding: 24, flex: 1 }}>
      <button
        onClick={() => navigate(`/admin/${service.key}`)}
        style={{ marginBottom: 16, padding: '6px 14px', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', background: 'white' }}
      >
        ← Back
      </button>
      <h2 style={{ marginBottom: 4 }}>{table.label}</h2>
      <p style={{ color: '#999', fontSize: 13, marginBottom: 20 }}>{service.label} · {service.baseURL}{table.listEndpoint}</p>

      {!table.readOnly && table.fields.length > 0 && (
        <AddForm fields={table.fields} onSubmit={data => createMutation.mutate(data)} />
      )}

      {isLoading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>Error loading data</p>}
      {!isLoading && !error && (
        <DataTable
          data={data}
          columns={table.columns}
          onDelete={id => deleteMutation.mutate(id)}
          readOnly={table.readOnly}
        />
      )}
    </div>
  )
}

// ─── ДОМАШНЯ СТОРІНКА СЕРВІСУ ─────────────────────────────────────────────────

function ServiceHome({ service }: { service: ServiceDef }) {
  const navigate = useNavigate()
  return (
    <div style={{ padding: 24 }}>
      <h2 style={{ marginBottom: 4 }}>{service.label}</h2>
      <p style={{ color: '#999', fontSize: 13, marginBottom: 24 }}>{service.baseURL}</p>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 10, maxWidth: 400 }}>
        {service.tables.map(table => (
          <button
            key={table.key}
            onClick={() => navigate(`/admin/${service.key}/${table.key}`)}
            style={{ padding: '14px 18px', textAlign: 'left', border: '1px solid #e0e0e0', borderRadius: 8, cursor: 'pointer', background: 'white', fontSize: 15 }}
          >
            {table.label}
            {table.readOnly && <span style={{ marginLeft: 8, fontSize: 11, color: '#aaa', background: '#f0f0f0', padding: '2px 6px', borderRadius: 4 }}>read-only</span>}
          </button>
        ))}
      </div>
    </div>
  )
}

// ─── ДИНАМІЧНИЙ РОУТЕР ────────────────────────────────────────────────────────

function DynamicServiceRoute() {
  const { serviceKey, tableKey } = useParams<{ serviceKey: string; tableKey: string }>()
  const service = SERVICES.find(s => s.key === serviceKey)
  if (!service) return <div style={{ padding: 24 }}>Service not found</div>
  if (!tableKey) return <ServiceHome service={service} />
  const table = service.tables.find(t => t.key === tableKey)
  if (!table) return <div style={{ padding: 24 }}>Table not found</div>
  return <TablePage service={service} table={table} />
}

// ─── САЙДБАР ─────────────────────────────────────────────────────────────────

function Sidebar() {
  const navigate = useNavigate()
  const currentPath = window.location.pathname

  return (
    <aside style={{ width: 220, minHeight: '100vh', borderRight: '1px solid #e0e0e0', padding: '16px 0', background: '#fafafa', flexShrink: 0 }}>
      <p style={{ padding: '0 16px', fontSize: 11, color: '#aaa', textTransform: 'uppercase', letterSpacing: 1, marginBottom: 8 }}>Services</p>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        {SERVICES.map(service => {
          const isActive = currentPath.includes(`/admin/${service.key}`)
          return (
            <button
              key={service.key}
              onClick={() => navigate(`/admin/${service.key}`)}
              style={{
                display: 'flex', alignItems: 'center', gap: 10,
                padding: '10px 16px', border: 'none',
                background: isActive ? '#EEF2FF' : 'transparent',
                color: isActive ? '#4A6CF7' : '#333',
                cursor: 'pointer', fontSize: 14, textAlign: 'left',
                borderRadius: 6, margin: '0 8px',
              }}
            >
              <span>{service.icon}</span>
              <span>{service.label}</span>
            </button>
          )
        })}
      </div>
    </aside>
  )
}

// ─── LAYOUT ───────────────────────────────────────────────────────────────────

function AdminLayout() {
  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <Sidebar />
      <div style={{ flex: 1 }}>
        <Routes>
          <Route path="/" element={<ServiceHome service={SERVICES[0]} />} />
          <Route path=":serviceKey" element={<DynamicServiceRoute />} />
          <Route path=":serviceKey/:tableKey" element={<DynamicServiceRoute />} />
        </Routes>
      </div>
    </div>
  )
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AdminLayout />
    </QueryClientProvider>
  )
}
