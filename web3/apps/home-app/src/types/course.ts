// Тип завдання. Дзеркалить схему бекенду: docs/schemas.md → "Task".
// /api/course/get-upcoming-events повертає масив Task — найближчі дедлайни.

export interface Task {
  id: string
  course_id: string
  slug: string
  title: string
  description: string
  max_mark: number
  deadline: string // ISO timestamp
  created_at: string
  updated_at: string
}
