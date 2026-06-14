import { api, asArray } from './api'
import type { StudentCoursePage, TeacherCoursePage, Task, TaskStudent, TaskComment } from '../types/course'

export async function getStudentCourses(): Promise<StudentCoursePage[]> {
  const res = await api.get('/api/course/courses/student')
  return asArray<StudentCoursePage>(res.data)
}

export async function getTeacherCourses(): Promise<TeacherCoursePage[]> {
  const res = await api.get('/api/course/courses/teacher')
  return asArray<TeacherCoursePage>(res.data)
}

export async function getTasksByCourse(courseId: string): Promise<Task[]> {
  const res = await api.get(`/api/course/course/${courseId}/tasks`)
  return asArray<Task>(res.data)
}

export async function getTask(taskId: string): Promise<Task> {
  const res = await api.get(`/api/course/task/${taskId}`)
  return res.data as Task
}

export async function getMySubmissionsForCourse(courseId: string): Promise<TaskStudent[]> {
  const res = await api.get(`/api/course/course/${courseId}/student-tasks`)
  return asArray<TaskStudent>(res.data)
}

export async function getTaskSubmissions(taskId: string): Promise<TaskStudent[]> {
  const res = await api.get(`/api/course/task/${taskId}/submissions`)
  return asArray<TaskStudent>(res.data)
}

export async function createTask(task: Pick<Task, 'course_id' | 'title' | 'description' | 'max_mark' | 'deadline'>): Promise<Task> {
  const res = await api.post('/api/course/tasks', task)
  return res.data as Task
}

export async function createAnnouncement(data: Pick<Task, 'course_id' | 'title' | 'description'> & { attachment_url?: string }): Promise<Task> {
  const res = await api.post('/api/course/tasks', {
    ...data,
    post_type: 'announcement',
    max_mark: 0,
    deadline: new Date(Date.now() + 365 * 86400000).toISOString(),
  })
  return res.data as Task
}

export async function submitTask(
  taskId: string,
  studentId: string,
  fileUrl?: string,
  linkUrl?: string,
): Promise<TaskStudent> {
  const res = await api.post('/api/course/task-students', {
    task_id: taskId,
    student_id: studentId,
    submission_time: new Date().toISOString(),
    file_url: fileUrl ?? null,
    link_url: linkUrl ?? null,
  })
  return res.data as TaskStudent
}

export async function updateSubmission(
  submissionId: string,
  data: { mark?: number; file_url?: string; link_url?: string },
): Promise<TaskStudent> {
  const res = await api.put(`/api/course/task-student/${submissionId}`, data)
  return res.data as TaskStudent
}

export async function gradeSubmission(submissionId: string, mark: number): Promise<TaskStudent> {
  const res = await api.put(`/api/course/task-student/${submissionId}`, { mark })
  return res.data as TaskStudent
}

export async function unsubmitTask(submissionId: string): Promise<void> {
  await api.delete(`/api/course/task-student/${submissionId}`)
}

export async function deleteTask(taskId: string): Promise<void> {
  await api.delete(`/api/course/task/${taskId}`)
}

export async function uploadFile(file: File): Promise<string> {
  const form = new FormData()
  form.append('file', file)
  const res = await api.post('/api/course/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return (res.data as { url: string }).url
}

export async function getTaskComments(taskId: string): Promise<TaskComment[]> {
  const res = await api.get(`/api/course/task/${taskId}/comments`)
  return asArray<TaskComment>(res.data)
}

export async function addTaskComment(taskId: string, body: string): Promise<TaskComment> {
  const res = await api.post(`/api/course/task/${taskId}/comments`, { body })
  return res.data as TaskComment
}

export async function getPrivateComments(taskId: string, studentId?: string): Promise<TaskComment[]> {
  const params = studentId ? `?student_id=${studentId}` : ''
  const res = await api.get(`/api/course/task/${taskId}/comments/private${params}`)
  return asArray<TaskComment>(res.data)
}

export async function addPrivateComment(
  taskId: string,
  body: string,
  studentId?: string,
): Promise<TaskComment> {
  const res = await api.post(`/api/course/task/${taskId}/comments/private`, {
    body,
    student_id: studentId ?? undefined,
  })
  return res.data as TaskComment
}
