export interface StudentCoursePage {
  id: string
  name: string
  description: string
  teacher_name: string
  average_mark: number
  assignments_count: number
  student_count: number
  upcoming_deadline: string
}

export interface TeacherCoursePage {
  id: string
  name: string
  description: string
  student_count: number
  assignments_count: number
}

export interface CourseInfo {
  id: string
  name: string
  description: string
  teacher_name?: string
  slug?: string
  colorIndex?: number 
}

export interface Task {
  id: string
  course_id: string
  slug: string
  title: string
  description: string
  max_mark: number
  deadline: string
  post_type: 'assignment' | 'announcement'
  attachment_url?: string
  created_at: string
  updated_at: string
}

export interface TaskStudent {
  id: string
  task_id: string
  student_id: string
  mark?: number
  submission_time?: string
  file_url?: string
  link_url?: string
  created_at: string
  updated_at: string
  task?: Task
  student_name?: string
}

export interface TaskComment {
  id: string
  task_id: string
  author_id: string
  author_name: string
  body: string
  is_private: boolean
  student_id?: string
  created_at: string
  updated_at: string
}