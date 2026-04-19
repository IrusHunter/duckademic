import { useState } from 'react'
import { useNavigate, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider, useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  getAcademicRanks, createAcademicRank, deleteAcademicRank,
  getAcademicDegrees, createAcademicDegree, deleteAcademicDegree,
  getEmployees, createEmployee, deleteEmployee,
  getTeachers, createTeacher, deleteTeacher,
  type AcademicRank, type AcademicDegree, type Employee, type Teacher
} from './api/employeeApi'

const queryClient = new QueryClient()

function DataTable<T extends { id: string }>({
  data,
  columns,
  onDelete,
}: {
  data: T[]
  columns: { key: keyof T; label: string }[]
  onDelete: (id: string) => void
}) {
  const [selected, setSelected] = useState<string[]>([])
  const [action, setAction] = useState('')

  const toggleSelect = (id: string) => {
    setSelected(prev => prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id])
  }

  const handleGo = () => {
    if (action === 'delete') {
      selected.forEach(id => onDelete(id))
      setSelected([])
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 12, display: 'flex', gap: 8, alignItems: 'center' }}>
        <span style={{ fontSize: 14 }}>Action:</span>
        <select value={action} onChange={e => setAction(e.target.value)} style={{ padding: '4px 8px', borderRadius: 4, border: '1px solid #ccc' }}>
          <option value="">---</option>
          <option value="delete">Delete selected</option>
        </select>
        <button onClick={handleGo} style={{ padding: '4px 12px', borderRadius: 4, border: '1px solid #ccc', cursor: 'pointer' }}>Go</button>
      </div>
      <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
        <thead>
          <tr style={{ background: '#f5f5f5' }}>
            <th style={{ padding: '8px 12px', textAlign: 'left', borderBottom: '2px solid #e0e0e0', width: 40 }}></th>
            {columns.map(col => (
              <th key={String(col.key)} style={{ padding: '8px 12px', textAlign: 'left', borderBottom: '2px solid #e0e0e0', color: '#555' }}>{col.label}</th>
            ))}
            <th style={{ padding: '8px 12px', borderBottom: '2px solid #e0e0e0', color: '#555' }}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {data.length === 0 && (
            <tr><td colSpan={columns.length + 2} style={{ padding: 24, textAlign: 'center', color: '#aaa' }}>Немає даних</td></tr>
          )}
          {data.map(row => (
            <tr key={row.id} style={{ background: selected.includes(row.id) ? '#e8f0fe' : 'white' }}>
              <td style={{ padding: '8px 12px', borderBottom: '1px solid #eee' }}>
                <input type="checkbox" checked={selected.includes(row.id)} onChange={() => toggleSelect(row.id)} />
              </td>
              {columns.map(col => (
                <td key={String(col.key)} style={{ padding: '8px 12px', borderBottom: '1px solid #eee' }}>{String(row[col.key] ?? '')}</td>
              ))}
              <td style={{ padding: '8px 12px', borderBottom: '1px solid #eee' }}>
                <button
                  onClick={() => onDelete(row.id)}
                  style={{ padding: '3px 10px', color: 'red', border: '1px solid red', borderRadius: 4, background: 'white', cursor: 'pointer', fontSize: 13 }}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function AddForm({ fields, onSubmit }: {
  fields: { name: string; label: string; required?: boolean }[]
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
            <div key={field.name} style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ fontSize: 12, color: '#666' }}>{field.label}</label>
              <input
                value={values[field.name] ?? ''}
                onChange={e => setValues(prev => ({ ...prev, [field.name]: e.target.value }))}
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

function TeachersPage() {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { data = [], isLoading } = useQuery({ queryKey: ['teachers'], queryFn: getTeachers })
  const deleteMutation = useMutation({
    mutationFn: deleteTeacher,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['teachers'] }),
    onError: (e) => console.error('deleteTeacher error:', e),
  })
  const createMutation = useMutation({ mutationFn: createTeacher, onSuccess: () => qc.invalidateQueries({ queryKey: ['teachers'] }) })

  return (
    <div style={{ padding: 24, flex: 1 }}>
      <button onClick={() => navigate('/admin/employee')} style={{ marginBottom: 16, padding: '6px 14px', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', background: 'white' }}>← Back</button>
      <h2 style={{ marginBottom: 20 }}>Teacher</h2>
      <AddForm
        fields={[
          { name: 'employee_id', label: 'Employee ID', required: true },
          { name: 'email', label: 'Email', required: true },
          { name: 'academic_rank_id', label: 'Academic Rank ID', required: true },
          { name: 'academic_degree_id', label: 'Academic Degree ID' },
        ]}
        onSubmit={(data) => createMutation.mutate(data as Omit<Teacher, 'id'>)}
      />
      {isLoading ? <p>Loading...</p> : (
        <DataTable<Teacher>
          data={data}
          columns={[
            { key: 'id', label: 'ID' },
            { key: 'email', label: 'Email' },
            { key: 'employee_id', label: 'Employee ID' },
            { key: 'academic_rank_id', label: 'Academic Rank' },
          ]}
          onDelete={(id) => deleteMutation.mutate(id)}
        />
      )}
    </div>
  )
}

function EmployeesPage() {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { data = [], isLoading } = useQuery({ queryKey: ['employees'], queryFn: getEmployees })
  const deleteMutation = useMutation({
    mutationFn: deleteEmployee,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['employees'] }),
    onError: (e) => console.error('deleteEmployee error:', e),
  })
  const createMutation = useMutation({ mutationFn: createEmployee, onSuccess: () => qc.invalidateQueries({ queryKey: ['employees'] }) })

  return (
    <div style={{ padding: 24, flex: 1 }}>
      <button onClick={() => navigate('/admin/employee')} style={{ marginBottom: 16, padding: '6px 14px', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', background: 'white' }}>← Back</button>
      <h2 style={{ marginBottom: 20 }}>Employee</h2>
      <AddForm
        fields={[
          { name: 'first_name', label: 'First Name', required: true },
          { name: 'last_name', label: 'Last Name', required: true },
          { name: 'middle_name', label: 'Middle Name' },
          { name: 'phone_number', label: 'Phone Number' },
        ]}
        onSubmit={(data) => createMutation.mutate(data as Omit<Employee, 'id'>)}
      />
      {isLoading ? <p>Loading...</p> : (
        <DataTable<Employee>
          data={data}
          columns={[
            { key: 'id', label: 'ID' },
            { key: 'first_name', label: 'First Name' },
            { key: 'last_name', label: 'Last Name' },
            { key: 'middle_name', label: 'Middle Name' },
            { key: 'phone_number', label: 'Phone' },
          ]}
          onDelete={(id) => deleteMutation.mutate(id)}
        />
      )}
    </div>
  )
}

function AcademicRanksPage() {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { data = [], isLoading } = useQuery({ queryKey: ['academic-ranks'], queryFn: getAcademicRanks })
  const deleteMutation = useMutation({
    mutationFn: deleteAcademicRank,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['academic-ranks'] }),
    onError: (e) => console.error('deleteAcademicRank error:', e),
  })
  const createMutation = useMutation({ mutationFn: createAcademicRank, onSuccess: () => qc.invalidateQueries({ queryKey: ['academic-ranks'] }) })

  return (
    <div style={{ padding: 24, flex: 1 }}>
      <button onClick={() => navigate('/admin/employee')} style={{ marginBottom: 16, padding: '6px 14px', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', background: 'white' }}>← Back</button>
      <h2 style={{ marginBottom: 20 }}>Academic Rank</h2>
      <AddForm
        fields={[{ name: 'title', label: 'Title', required: true }]}
        onSubmit={(data) => createMutation.mutate(data as { title: string })}
      />
      {isLoading ? <p>Loading...</p> : (
        <DataTable<AcademicRank>
          data={data}
          columns={[{ key: 'id', label: 'ID' }, { key: 'title', label: 'Title' }]}
          onDelete={(id) => deleteMutation.mutate(id)}
        />
      )}
    </div>
  )
}

function AcademicDegreesPage() {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { data = [], isLoading } = useQuery({ queryKey: ['academic-degrees'], queryFn: getAcademicDegrees })
  const deleteMutation = useMutation({
    mutationFn: deleteAcademicDegree,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['academic-degrees'] }),
    onError: (e) => console.error('deleteAcademicDegree error:', e),
  })
  const createMutation = useMutation({ mutationFn: createAcademicDegree, onSuccess: () => qc.invalidateQueries({ queryKey: ['academic-degrees'] }) })

  return (
    <div style={{ padding: 24, flex: 1 }}>
      <button onClick={() => navigate('/admin/employee')} style={{ marginBottom: 16, padding: '6px 14px', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', background: 'white' }}>← Back</button>
      <h2 style={{ marginBottom: 20 }}>Academic Degree</h2>
      <AddForm
        fields={[{ name: 'title', label: 'Title', required: true }]}
        onSubmit={(data) => createMutation.mutate(data as { title: string })}
      />
      {isLoading ? <p>Loading...</p> : (
        <DataTable<AcademicDegree>
          data={data}
          columns={[{ key: 'id', label: 'ID' }, { key: 'title', label: 'Title' }]}
          onDelete={(id) => deleteMutation.mutate(id)}
        />
      )}
    </div>
  )
}

const EMPLOYEE_TABLES = [
  { label: 'Teacher', path: '/admin/employee/teachers', description: 'Manage teachers' },
  { label: 'Employee', path: '/admin/employee/employees', description: 'Manage employees' },
  { label: 'Academic Rank', path: '/admin/employee/academic-ranks', description: 'Manage academic ranks' },
  { label: 'Academic Degree', path: '/admin/employee/academic-degrees', description: 'Manage academic degrees' },
]

function EmployeeServiceHome() {
  const navigate = useNavigate()
  return (
    <div style={{ padding: 24 }}>
      <h2 style={{ marginBottom: 20 }}>Employee Service</h2>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 10, maxWidth: 360 }}>
        {EMPLOYEE_TABLES.map(table => (
          <button
            key={table.path}
            onClick={() => navigate(table.path)}
            style={{ padding: '14px 18px', textAlign: 'left', border: '1px solid #e0e0e0', borderRadius: 8, cursor: 'pointer', background: 'white', fontSize: 15 }}
          >
            {table.label}
            <span style={{ display: 'block', fontSize: 12, color: '#999', marginTop: 2 }}>{table.description}</span>
          </button>
        ))}
      </div>
    </div>
  )
}

const SIDEBAR_SERVICES = [
  { label: 'Employee Service', path: '/admin/employee', icon: '⊞' },
  { label: 'Schedule Service', path: '/admin/schedule', icon: '💬' },
  { label: 'Courses Service', path: '/admin/courses', icon: '📅' },
]

function Sidebar() {
  const navigate = useNavigate()
  return (
    <aside style={{ width: 220, minHeight: '100vh', borderRight: '1px solid #e0e0e0', padding: '16px 0', background: '#fafafa', flexShrink: 0 }}>
      <p style={{ padding: '0 16px', fontSize: 11, color: '#aaa', textTransform: 'uppercase', letterSpacing: 1, marginBottom: 8 }}>Services</p>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        {SIDEBAR_SERVICES.map(item => (
          <button
            key={item.path}
            onClick={() => navigate(item.path)}
            style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '10px 16px', border: 'none', background: 'transparent', cursor: 'pointer', fontSize: 14, color: '#333', textAlign: 'left', borderRadius: 6, margin: '0 8px' }}
          >
            <span>{item.icon}</span>
            <span>{item.label}</span>
          </button>
        ))}
      </div>
    </aside>
  )
}

function AdminContent() {
  return (
    <div style={{ flex: 1 }}>
      <Routes>
        <Route path="/" element={<EmployeeServiceHome />} />
        <Route path="employee" element={<EmployeeServiceHome />} />
        <Route path="employee/teachers" element={<TeachersPage />} />
        <Route path="employee/employees" element={<EmployeesPage />} />
        <Route path="employee/academic-ranks" element={<AcademicRanksPage />} />
        <Route path="employee/academic-degrees" element={<AcademicDegreesPage />} />
        <Route path="schedule" element={<div style={{ padding: 24 }}>Schedule Service — coming soon</div>} />
        <Route path="courses" element={<div style={{ padding: 24 }}>Courses Service — coming soon</div>} />
      </Routes>
    </div>
  )
}

function AdminLayout() {
  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <Sidebar />
      <AdminContent />
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