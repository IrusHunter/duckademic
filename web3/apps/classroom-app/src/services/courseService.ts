import { api, asArray } from './api'
import type { StudentCoursePage, TeacherCoursePage } from '../types/course'

export async function getStudentCourses(): Promise<StudentCoursePage[]> {
  const res = await api.get('/api/course/courses/student')
  return asArray<StudentCoursePage>(res.data)
}

export async function getTeacherCourses(): Promise<TeacherCoursePage[]> {
  const res = await api.get('/api/course/courses/teacher')
  return asArray<TeacherCoursePage>(res.data)
}