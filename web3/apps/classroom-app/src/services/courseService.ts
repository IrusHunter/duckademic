import { api, asArray } from './api'
import type { StudentCoursePage } from '../types/course'

export async function getStudentCourses(): Promise<StudentCoursePage[]> {
  const res = await api.get('/api/course/courses/student')
  return asArray<StudentCoursePage>(res.data)
}