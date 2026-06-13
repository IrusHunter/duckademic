import { api, asArray } from './api'
import type { LessonOccurrence, LessonSlot, Teacher, Classroom, StudentGroup } from '../types/schedule'

export async function getLessonOccurrences(): Promise<LessonOccurrence[]> {
  const res = await api.get('/api/schedule/lesson-occurrences')
  return asArray<LessonOccurrence>(res.data)
}

export async function getLessonSlots(): Promise<LessonSlot[]> {
  const res = await api.get('/api/schedule/lesson-slots')
  return asArray<LessonSlot>(res.data)
}

export async function getTeachers(): Promise<Teacher[]> {
  const res = await api.get('/api/schedule/teachers')
  return asArray<Teacher>(res.data)
}

export async function getClassrooms(): Promise<Classroom[]> {
  const res = await api.get('/api/schedule/classrooms')
  return asArray<Classroom>(res.data)
}

export async function getStudentGroups(): Promise<StudentGroup[]> {
  const res = await api.get('/api/schedule/student-groups')
  return asArray<StudentGroup>(res.data)
}
