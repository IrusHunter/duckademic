export interface StudyLoadInfo {
  id: string
  teacher_id: string
  teacher_name: string
  student_group_id: string
  student_group_name: string
  discipline_id: string
  discipline_name: string
  lesson_type_id: string
  lesson_type_name: string
}

export interface ClassroomInfo {
  id: string
  slug: string
  number: string
  capacity: number
}

export interface LessonOccurrence {
  id: string
  study_load_id: string
  lesson_slot_id: string
  date: string  // ISO timestamp e.g. "2025-01-20T08:00:00Z"
  classroom_id?: string
  status: string
  study_load?: StudyLoadInfo
  classroom?: ClassroomInfo
}
