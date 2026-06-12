// Типи розкладу. Дзеркалять схему бекенду:
// docs/schemas.md → "Lesson Occurrence" (повна модель з вкладеним study_load).

export type LessonStatus = 'scheduled' | 'canceled' | 'completed'

export interface StudyLoad {
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

export interface LessonOccurrence {
  id: string
  study_load_id: string
  lesson_slot_id: string
  date: string // ISO timestamp — дата й час заняття
  status: LessonStatus
  classroom_id?: string | null
  study_load?: StudyLoad // присутній у "повній" моделі
}
