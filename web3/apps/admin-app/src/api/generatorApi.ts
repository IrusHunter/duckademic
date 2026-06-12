import { makeApi } from './makeApi'

// schedule-сервіс: завантаження даних у генератор + витяг результату
const scheduleApi = makeApi('/api/schedule')
// сам генератор: кроки пайплайну
const genApi = makeApi('/api/schedule-generator')

// Усі ці маршрути на бекенді — HandleFunc (приймають POST), тому всюди POST.

export function loadDataIntoGenerator(semesterIds: string[]) {
  return scheduleApi.post('/load-data-into-generator', semesterIds)
}

export function loadClassroomsIntoGenerator(classroomIds: string[]) {
  return scheduleApi.post('/load-classrooms-into-generator', classroomIds)
}

export function submitAndGo() {
  return genApi.post('/submit-and-go', { ignore_warnings: true })
}

export function processStep(method: string) {
  return genApi.post('/process-step', { method })
}

export function extractDataFromGenerator(startTime: string) {
  return scheduleApi.post('/extract-data-from-generator', { start_time: startTime })
}
