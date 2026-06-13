import { api, asArray } from './api'
import type { StudentCoursePage, TeacherCoursePage, Task, TaskStudent } from '../types/course'

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

export async function getMySubmissionsForCourse(courseId: string): Promise<TaskStudent[]> {
  const res = await api.get(`/api/course/course/${courseId}/student-tasks`)
  return asArray<TaskStudent>(res.data)
}

export async function createTask(task: Pick<Task, 'course_id' | 'title' | 'description' | 'max_mark' | 'deadline'>): Promise<Task> {
  const res = await api.post('/api/course/tasks', task)
  return res.data as Task
}

export async function submitTask(taskId: string, studentId: string): Promise<TaskStudent> {
  const res = await api.post('/api/course/task-students', {
    task_id: taskId,
    student_id: studentId,
    submission_time: new Date().toISOString(),
  })
  return res.data as TaskStudent
}

export async function unsubmitTask(submissionId: string): Promise<void> {
  await api.delete(`/api/course/task-student/${submissionId}`)
}

export async function deleteTask(taskId: string): Promise<void> {
  await api.delete(`/api/course/task/${taskId}`)
}