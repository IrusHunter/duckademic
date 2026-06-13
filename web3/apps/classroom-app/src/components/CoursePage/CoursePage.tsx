import { useState, useMemo } from 'react'
import { useParams, useLocation, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { LuArrowLeft, LuStar, LuMegaphone } from "react-icons/lu";
import {
  getTasksByCourse,
  getMySubmissionsForCourse,
  createTask,
  createAnnouncement,
  deleteTask,
} from '../../services/courseService'
import { getAuthUser } from '../../utils/user'
import type { CourseInfo, Task } from '../../types/course'
import css from './CoursePage.module.css'

// ── helpers ────────────────────────────────────────────────────────────────

const PALETTE = [
  '#DB4437','#E67E22','#F4B400','#0F9D58','#16A085',
  '#4285F4','#2980B9','#8E44AD','#9B59B6','#E91E63',
]

function hashStr(s: string): number {
  let h = 0
  for (let i = 0; i < s.length; i++) h = (Math.imul(31, h) + s.charCodeAt(i)) | 0
  return Math.abs(h)
}

function colorFor(name: string) { return PALETTE[hashStr(name || '?') % PALETTE.length] }

function initials(name: string) {
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return name.slice(0, 2).toUpperCase()
}

function timeAgo(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime()
  const h = Math.floor(diff / 3_600_000)
  if (h < 1) return 'Just now'
  if (h < 24) return `${h} hour${h > 1 ? 's' : ''} ago`
  const d = Math.floor(h / 24)
  return `${d} day${d > 1 ? 's' : ''} ago`
}

function fmtDeadline(iso: string) {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function isOverdue(iso: string) { return new Date(iso) < new Date() }

// ── AvatarCircle ────────────────────────────────────────────────────────────

function AvatarCircle({ name, size = 36 }: { name: string; size?: number }) {
  return (
    <div
      className={css.avatar}
      style={{ width: size, height: size, backgroundColor: colorFor(name), fontSize: Math.round(size * 0.38) }}
    >
      {initials(name)}
    </div>
  )
}

// ── Add-task modal ──────────────────────────────────────────────────────────

interface AddTaskModalProps {
  courseId: string
  onClose: () => void
}

function AddTaskModal({ courseId, onClose }: AddTaskModalProps) {
  const qc = useQueryClient()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [maxMark, setMaxMark] = useState('100')
  const [deadline, setDeadline] = useState('')

  const mutation = useMutation({
    mutationFn: () => createTask({
      course_id: courseId,
      title: title.trim(),
      description: description.trim(),
      max_mark: parseFloat(maxMark) || 100,
      deadline: deadline ? new Date(deadline).toISOString() : new Date(Date.now() + 7 * 86400000).toISOString(),
    }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['course-tasks', courseId] })
      onClose()
    },
  })

  return (
    <div className={css.overlay} onClick={(e) => e.target === e.currentTarget && onClose()}>
      <div className={css.modal}>
        <h2 className={css.modalTitle}>Add Assignment</h2>

        <div className={css.formGroup}>
          <label className={css.formLabel}>Title</label>
          <input className={css.formInput} value={title} onChange={e => setTitle(e.target.value)} placeholder="Assignment title" />
        </div>

        <div className={css.formGroup}>
          <label className={css.formLabel}>Description</label>
          <textarea className={css.formTextarea} value={description} onChange={e => setDescription(e.target.value)} placeholder="Describe the task…" />
        </div>

        <div className={css.formRow}>
          <div className={css.formGroup}>
            <label className={css.formLabel}>Max Mark</label>
            <input className={css.formInput} type="number" min="1" value={maxMark} onChange={e => setMaxMark(e.target.value)} />
          </div>
          <div className={css.formGroup}>
            <label className={css.formLabel}>Deadline</label>
            <input className={css.formInput} type="datetime-local" value={deadline} onChange={e => setDeadline(e.target.value)} />
          </div>
        </div>

        <div className={css.modalActions}>
          <button className={css.btnSecondary} onClick={onClose}>Cancel</button>
          <button
            className={css.btnAdd}
            disabled={!title.trim() || mutation.isPending}
            onClick={() => mutation.mutate()}
          >
            {mutation.isPending ? 'Saving…' : 'Add Assignment'}
          </button>
        </div>
      </div>
    </div>
  )
}

// ── Add-post modal ──────────────────────────────────────────────────────────

function AddPostModal({ courseId, onClose }: { courseId: string; onClose: () => void }) {
  const qc = useQueryClient()
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')

  const mutation = useMutation({
    mutationFn: () => createAnnouncement({
      course_id: courseId,
      title: title.trim(),
      description: content.trim(),
    }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['course-tasks', courseId] })
      onClose()
    },
  })

  return (
    <div className={css.overlay} onClick={e => e.target === e.currentTarget && onClose()}>
      <div className={css.modal}>
        <h2 className={css.modalTitle}>Новий пост</h2>

        <div className={css.formGroup}>
          <label className={css.formLabel}>Заголовок</label>
          <input className={css.formInput} value={title} onChange={e => setTitle(e.target.value)} placeholder="Заголовок поста" />
        </div>

        <div className={css.formGroup}>
          <label className={css.formLabel}>Повідомлення</label>
          <textarea className={css.formTextarea} value={content} onChange={e => setContent(e.target.value)} placeholder="Текст повідомлення або матеріали…" style={{ minHeight: 120 }} />
        </div>

        <div className={css.modalActions}>
          <button className={css.btnSecondary} onClick={onClose}>Скасувати</button>
          <button
            className={css.btnAdd}
            disabled={!title.trim() || mutation.isPending}
            onClick={() => mutation.mutate()}
          >
            {mutation.isPending ? 'Публікуємо…' : 'Опублікувати'}
          </button>
        </div>
      </div>
    </div>
  )
}

// ── Main ────────────────────────────────────────────────────────────────────

type Tab = 'stream' | 'classwork' | 'people' | 'grades'

export default function CoursePage() {
  const { courseId } = useParams<{ courseId: string }>()
  const location = useLocation()
  const navigate = useNavigate()
  const user = getAuthUser()
  const isTeacher = user?.role === 'teacher'

  const course = (location.state as CourseInfo | null) ?? { id: courseId ?? '', name: 'Course', description: '' }

  const [activeTab, setActiveTab] = useState<Tab>('stream')
  const [showModal, setShowModal] = useState(false)
  const [showPostModal, setShowPostModal] = useState(false)

  const qc = useQueryClient()

  const { data: tasks = [] } = useQuery({
    queryKey: ['course-tasks', courseId],
    queryFn: () => getTasksByCourse(courseId!),
    enabled: !!courseId,
  })

  const { data: submissions = [] } = useQuery({
    queryKey: ['my-submissions', courseId],
    queryFn: () => getMySubmissionsForCourse(courseId!),
    enabled: !!courseId && !isTeacher,
  })

  const submissionMap = useMemo(() => {
    const m = new Map<string, typeof submissions[0]>()
    for (const s of submissions) m.set(s.task_id, s)
    return m
  }, [submissions])

  const deleteTaskMutation = useMutation({
    mutationFn: (taskId: string) => deleteTask(taskId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['course-tasks', courseId] }),
  })

  const teacherName = course.teacher_name ?? 'Instructor'
  const slug = course.id.split('-')[0].toUpperCase()

  const BANNER_CLASSES = [css.bannerBlue, css.bannerGreen, css.bannerPink]
  const bannerClass = BANNER_CLASSES[
    course.colorIndex !== undefined ? course.colorIndex % 3 : hashStr(course.id) % 3
  ]

  // upcoming tasks sorted by deadline
  const upcoming = useMemo(() =>
    [...tasks]
      .filter(t => !isOverdue(t.deadline) && !submissionMap.has(t.id))
      .sort((a, b) => new Date(a.deadline).getTime() - new Date(b.deadline).getTime())
      .slice(0, 3),
    [tasks, submissionMap]
  )

  // stream: sorted newest first
  const streamTasks = useMemo(() =>
    [...tasks].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()),
    [tasks]
  )

  // classwork: sorted by deadline asc
  const classworkTasks = useMemo(() =>
    [...tasks].sort((a, b) => new Date(a.deadline).getTime() - new Date(b.deadline).getTime()),
    [tasks]
  )

  function goToTask(taskId: string) {
    navigate(`task/${taskId}`, { state: course })
  }

  function renderTaskCard(task: Task, mode: 'stream' | 'classwork') {
    const submission = submissionMap.get(task.id)
    const isAnnouncement = task.post_type === 'announcement'
    const clickable = !isAnnouncement

    return (
      <div
        key={task.id}
        className={`${mode === 'stream' ? css.postCard : css.classworkCard} ${clickable ? css.cardClickable : ''}`}
        onClick={clickable ? () => goToTask(task.id) : undefined}
        role={clickable ? 'button' : undefined}
        tabIndex={clickable ? 0 : undefined}
        onKeyDown={clickable ? (e) => e.key === 'Enter' && goToTask(task.id) : undefined}
      >
        {mode === 'stream' ? (
          <>
            <div className={css.postHeader}>
              <div className={css.postAuthor}>
                <AvatarCircle name={teacherName} />
                <div>
                  <div className={css.authorName}>{teacherName}</div>
                  <div className={css.authorTime}>{timeAgo(task.created_at)}</div>
                </div>
              </div>
              {isTeacher && (
                <button
                  className={css.btnDanger}
                  onClick={(e) => {
                    e.stopPropagation()
                    if (confirm(`Delete "${task.title}"?`)) deleteTaskMutation.mutate(task.id)
                  }}
                >
                  Delete
                </button>
              )}
            </div>

            {isAnnouncement ? (
              <>
                <div className={css.postTitle}>
                  <LuMegaphone size={16} style={{ marginRight: 6, color: '#6b7280', verticalAlign: 'middle' }} />
                  {task.title}
                </div>
                {task.description && <div className={css.postDesc}>{task.description}</div>}
              </>
            ) : (
              <>
                <div className={css.postTitle}>{task.title}</div>
                {task.description && <div className={css.postDesc}>{task.description}</div>}

                <div className={css.postMeta}>
                  <div className={css.postMetaLeft}>
                    <span className={css.metaBadge}>Deadline: {fmtDeadline(task.deadline)}</span>
                    <span className={css.metaBadge}><LuStar /> {task.max_mark} pts</span>
                    {!isTeacher && !submission && <span className={css.badgeAvailable}>available</span>}
                    {!isTeacher && submission && !submission.mark && <span className={css.badgeSubmitted}>Перевіряється</span>}
                    {!isTeacher && submission?.mark !== undefined && submission.mark !== null && (
                      <span className={css.badgeGraded}>Score: {submission.mark}/{task.max_mark}</span>
                    )}
                  </div>
                </div>
              </>
            )}

            <div className={css.postFooter}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
              </svg>
              0 comments
            </div>
          </>
        ) : (
          // classwork card
          <>
            <div className={css.classworkIcon}>📋</div>
            <div className={css.classworkInfo}>
              <div className={css.classworkTitle}>{task.title}</div>
              <div className={css.classworkMeta}>
                <span>Deadline: {fmtDeadline(task.deadline)}</span>
                <span>{task.max_mark} pts</span>
                {!isTeacher && submission && !submission.mark && <span style={{ color: '#2563eb' }}>Перевіряється</span>}
                {!isTeacher && submission?.mark !== undefined && submission.mark !== null && (
                  <span style={{ color: '#16a34a' }}>Score: {submission.mark}/{task.max_mark}</span>
                )}
              </div>
            </div>
            {isTeacher && (
              <div className={css.classworkActions}>
                <button
                  className={css.btnDanger}
                  onClick={(e) => {
                    e.stopPropagation()
                    if (confirm(`Delete "${task.title}"?`)) deleteTaskMutation.mutate(task.id)
                  }}
                >
                  Delete
                </button>
              </div>
            )}
          </>
        )}
      </div>
    )
  }

  return (
    <div className={css.page}>
      {/* Banner */}
      <div className={`${css.banner} ${bannerClass}`}>
        <button className={css.backLink} onClick={() => navigate(-1)}>
          <LuArrowLeft />
          Back to Courses
        </button>
        <div className={css.bannerContent}>
          <h1 className={css.courseName}>{course.name}</h1>
          {course.description && <p className={css.courseDescription}>{course.description}</p>}
          <span className={css.courseMeta}>{slug} • {teacherName}</span>
        </div>
      </div>

      {/* Tabs */}
      <div className={css.tabs}>
        {(['stream', 'classwork', 'people', 'grades'] as Tab[]).map(tab => (
          <button
            key={tab}
            className={`${css.tabBtn} ${activeTab === tab ? css.tabActive : ''}`}
            onClick={() => setActiveTab(tab)}
          >
            {tab.charAt(0).toUpperCase() + tab.slice(1)}
          </button>
        ))}
      </div>

      {/* Body */}
      <div className={css.body}>
        {/* Left feed */}
        <div className={css.feed}>

          {/* Stream tab */}
          {activeTab === 'stream' && (
            <>
              {isTeacher && (
                <div className={css.actionBar}>
                  <button className={css.btnAdd} onClick={() => setShowModal(true)}>
                    + Add Assignment
                  </button>
                  <button className={css.btnAdd} onClick={() => setShowPostModal(true)}>
                    <LuMegaphone size={14} /> Add Post
                  </button>
                </div>
              )}
              {streamTasks.length === 0
                ? <p className={css.empty}>No assignments yet.</p>
                : streamTasks.map(t => renderTaskCard(t, 'stream'))
              }
            </>
          )}

          {/* Classwork tab */}
          {activeTab === 'classwork' && (
            <>
              {isTeacher && (
                <div className={css.actionBar}>
                  <button className={css.btnAdd} onClick={() => setShowModal(true)}>
                    + Add Assignment
                  </button>
                  <button className={css.btnAdd} onClick={() => setShowPostModal(true)}>
                    <LuMegaphone size={14} /> Add Post
                  </button>
                </div>
              )}
              {classworkTasks.length === 0
                ? <p className={css.empty}>No assignments yet.</p>
                : classworkTasks.map(t => renderTaskCard(t, 'classwork'))
              }
            </>
          )}

          {/* People tab */}
          {activeTab === 'people' && (
            <div className={css.peopleSection}>
              <div className={css.peopleSectionTitle}>Teachers</div>
              <div className={css.peopleTeacher}>
                <AvatarCircle name={teacherName} size={40} />
                <div>
                  <div className={css.peopleName}>{teacherName}</div>
                  <div className={css.peopleRole}>Instructor</div>
                </div>
              </div>
            </div>
          )}

          {/* Grades tab */}
          {activeTab === 'grades' && (
            <div className={css.gradesTable}>
              <div className={`${css.gradesRow} ${css.gradesHeader}`}>
                <span>Assignment</span>
                <span>Deadline</span>
                <span>Grade</span>
              </div>
              {tasks.length === 0
                ? <p className={css.empty}>No assignments yet.</p>
                : tasks.map(task => {
                    const sub = submissionMap.get(task.id)
                    return (
                      <div key={task.id} className={css.gradesRow}>
                        <span className={css.gradesTaskName}>{task.title}</span>
                        <span className={css.gradesDeadline}>{fmtDeadline(task.deadline)}</span>
                        <span className={`${css.gradesMark} ${sub?.mark !== undefined && sub.mark !== null ? css.gradesMarkGreen : css.gradesMarkGray}`}>
                          {sub?.mark !== undefined && sub.mark !== null
                            ? `${sub.mark} / ${task.max_mark}`
                            : sub ? '—' : 'Not submitted'}
                        </span>
                      </div>
                    )
                  })
              }
            </div>
          )}
        </div>

        {/* Sidebar */}
        <div className={css.sidebar}>
          <div className={css.sideCard}>
            <div className={css.sideCardTitle}>Course Code</div>
            <div className={css.sideCardValue}>{slug}</div>
          </div>

          {upcoming.length > 0 && (
            <div className={css.sideCard}>
              <div className={css.sideCardTitle}>Upcoming</div>
              <ul className={css.upcomingList}>
                {upcoming.map((t, i) => (
                  <li key={t.id} className={css.upcomingItem}>
                    <div className={`${css.upcomingDot} ${i === 0 ? css.upcomingDotOrange : css.upcomingDotBlue}`} />
                    <div>
                      <div className={css.upcomingTitle}>{t.title}</div>
                      <div className={css.upcomingDate}>{fmtDeadline(t.deadline)}</div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </div>

      {showModal && courseId && (
        <AddTaskModal courseId={courseId} onClose={() => setShowModal(false)} />
      )}
      {showPostModal && courseId && (
        <AddPostModal courseId={courseId} onClose={() => setShowPostModal(false)} />
      )}
    </div>
  )
}
