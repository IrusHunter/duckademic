import axios from 'axios'

const api = axios.create({
  baseURL: '/api/employee', // ← відносний шлях через proxy
  withCredentials: true,
})

// Типи
export type AcademicRank = { id: string; title: string }
export type AcademicDegree = { id: string; title: string }
export type Employee = {
  id: string
  first_name: string
  last_name: string
  middle_name?: string
  phone_number?: string
}
export type Teacher = {
  id: string
  employee_id: string
  email: string
  academic_rank_id: string
  academic_degree_id?: string
}

// Academic Ranks
export const getAcademicRanks = () => api.get<AcademicRank[]>('/academic-ranks').then(r => r.data)
export const createAcademicRank = (data: { title: string }) => api.post<AcademicRank>('/academic-ranks', data).then(r => r.data)
export const deleteAcademicRank = (id: string) => api.delete(`/academic-rank/${id}`).then(r => r.data)
export const updateAcademicRank = (id: string, data: { title: string }) => api.put<AcademicRank>(`/academic-rank/${id}`, data).then(r => r.data)

// Academic Degrees
export const getAcademicDegrees = () => api.get<AcademicDegree[]>('/academic-degrees').then(r => r.data)
export const createAcademicDegree = (data: { title: string }) => api.post<AcademicDegree>('/academic-degrees', data).then(r => r.data)
export const deleteAcademicDegree = (id: string) => api.delete(`/academic-degree/${id}`).then(r => r.data)
export const updateAcademicDegree = (id: string, data: { title: string }) => api.put<AcademicDegree>(`/academic-degree/${id}`, data).then(r => r.data)

// Employees
export const getEmployees = () => api.get<Employee[]>('/employees').then(r => r.data)
export const createEmployee = (data: Omit<Employee, 'id'>) => api.post<Employee>('/employees', data).then(r => r.data)
export const deleteEmployee = (id: string) => api.delete(`/employee/${id}`).then(r => r.data)
export const updateEmployee = (id: string, data: Omit<Employee, 'id'>) => api.put<Employee>(`/employee/${id}`, data).then(r => r.data)

// Teachers
export const getTeachers = () => api.get<Teacher[]>('/teachers').then(r => r.data)
export const createTeacher = (data: Omit<Teacher, 'id'>) => api.post<Teacher>('/teachers', data).then(r => r.data)
export const deleteTeacher = (id: string) => api.delete(`/teacher/${id}`).then(r => r.data)
export const updateTeacher = (id: string, data: Omit<Teacher, 'id'>) => api.put<Teacher>(`/teacher/${id}`, data).then(r => r.data)