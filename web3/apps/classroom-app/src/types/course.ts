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