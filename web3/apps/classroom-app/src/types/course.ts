// Студентські ендпоінти course-сервісу (дозволені ролі student) — docs/schemas.md

// GET /api/course/get-courses-progress → мої курси + прогрес
export interface CourseProgress {
  id: string
  name: string
  complete_rate: number      // 0..1 — частка проходження
  complete_accuracy: number  // 0..1 — частка правильних
}

// GET /api/course/get-upcoming-events → мої найближчі завдання
export interface Task {
  id: string
  course_id: string
  slug: string
  title: string
  description: string
  max_mark: number
  deadline: string
  created_at?: string
  updated_at?: string
}
