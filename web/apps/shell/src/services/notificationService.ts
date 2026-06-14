import axios from 'axios'

export interface Notification {
  id: string
  recipient_id: string
  type: 'new_task' | 'submission' | 'grade'
  title: string
  task_id: string
  course_id: string
  task_student_id?: string
  is_read: boolean
  created_at: string
}

export async function getNotifications(): Promise<Notification[]> {
  const res = await axios.get('/api/course/notifications')
  const data = res.data
  if (Array.isArray(data)) return data
  if (data && Array.isArray(data.data)) return data.data
  return []
}

export async function markAllNotificationsRead(): Promise<void> {
  await axios.post('/api/course/notifications/read-all')
}
