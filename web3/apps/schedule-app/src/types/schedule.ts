export interface LessonSlot {
  id: string
  slot: number
  weekday: number     // 1 = Mon … 6 = Sat
  start_time: number  // nanoseconds since midnight
  duration: number    // nanoseconds
}

export interface LessonOccurrence {
  id: string
  study_load_id: string
  teacher_id: string
  student_group_id: string
  lesson_slot_id: string
  classroom_id: string
  date: string   // "YYYY-MM-DD"
  status: string
}

export interface Teacher {
  id: string
  slug: string
  name: string
}

export interface Classroom {
  id: string
  slug: string
  number: string
}

export interface StudentGroup {
  id: string
  slug: string
  name: string
}
