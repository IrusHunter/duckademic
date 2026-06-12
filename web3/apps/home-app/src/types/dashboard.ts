// Прогрес курсу — реальний ендпоінт /api/course/get-courses-progress
export interface CourseProgress {
  id: string
  name: string
  complete_rate: number      // 0..1 — частка проходження
  complete_accuracy: number  // 0..1 — частка правильних
}

// Завдання — реальний ендпоінт /api/course/get-upcoming-events
export interface Task {
  id: string
  course_id: string
  slug: string
  title: string
  description: string
  max_mark: number
  deadline: string
  created_at: string
  updated_at: string
}

// Дані залогіненого користувача (з localStorage, кладе shell)
export interface AuthUser {
  id: string
  email: string
  role: 'admin' | 'teacher' | 'student'
  is_default_password?: boolean
}

// Пост стрічки — поки моки (бекенд ще не реалізував стрічку)
export interface Post {
  id: string
  avatar: string
  author: string
  role: string
  time: string
  content: string
  likes: number
  comments: number
  shares: number
}
