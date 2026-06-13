import { makeApi } from './makeApi'

const scheduleApi = makeApi('/api/schedule')
const genApi = makeApi('/api/schedule-generator')
const curriculumApi = makeApi('/api/curriculum')
const assetApi = makeApi('/api/asset')

export interface Semester {
  id: string
  slug: string
  curriculum_id: string
  number: number
}

export interface Classroom {
  id: string
  slug: string
  number: string
  capacity: number
}

export interface GeneratorConfig {
  start_date: string
  end_date: string
  slot_preference: number[][]
  max_daily_student_load: number
  lesson_fill_rate: number
  classroom_occupancy: number
}

export function fetchSemesters() {
  return curriculumApi.get<Semester[]>('/semesters')
}

export function fetchClassrooms() {
  return assetApi.get<Classroom[]>('/classrooms')
}

export function initGenerator(config: GeneratorConfig) {
  return genApi.post('/init', config)
}

export function loadDataIntoGenerator(semesterIds: string[]) {
  return scheduleApi.post('/load-data-into-generator', semesterIds)
}

export function loadClassroomsIntoGenerator(classroomIds: string[]) {
  return scheduleApi.post('/load-classrooms-into-generator', classroomIds)
}

export function submitAndGo(ignoreWarnings = false) {
  return genApi.post<{ warnings?: string[] }>('/submit-and-go', { ignore_warnings: ignoreWarnings })
}

export function processStep(method: string) {
  return genApi.post<unknown>('/process-step', { method })
}

export function extractDataFromGenerator(startTime: string) {
  return scheduleApi.post('/extract-data-from-generator', { start_time: startTime })
}
