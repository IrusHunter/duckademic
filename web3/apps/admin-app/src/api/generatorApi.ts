import { makeApi } from './makeApi'

const scheduleApi = makeApi('/api/schedule')
const genApi = makeApi('/api/schedule-generator')

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

export function extractWorkloadsFromGenerator() {
  return scheduleApi.post('/extract-workloads-from-generator', {})
}

export function extractDataFromGenerator(startTime: string) {
  return scheduleApi.post('/extract-data-from-generator', { start_time: startTime })
}
