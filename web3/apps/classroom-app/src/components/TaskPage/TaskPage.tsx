import { useState, useRef } from 'react'
import { useParams, useNavigate, useLocation } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { LuArrowLeft, LuPaperclip, LuLink, LuFileText, LuExternalLink } from 'react-icons/lu'
import {
  getTask,
  getTaskSubmissions,
  getMySubmissionsForCourse,
  submitTask,
  unsubmitTask,
  updateSubmission,
  gradeSubmission,
  uploadFile,
} from '../../services/courseService'
import { getAuthUser } from '../../utils/user'
import type { CourseInfo, TaskStudent } from '../../types/course'
import css from './TaskPage.module.css'

// ── helpers ─────────────────────────────────────────────────────────────────

function fmtDeadline(iso: string) {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function isImageUrl(url: string) {
  return /\.(png|jpe?g|gif|webp|svg)$/i.test(url)
}

const MAX_FILE_SIZE = 10 * 1024 * 1024 // 10 MB

// ── File Preview Modal ───────────────────────────────────────────────────────

function FileModal({ url, name, onClose }: { url: string; name: string; onClose: () => void }) {
  const fullUrl = url.startsWith('/uploads/') ? `/api/course${url}` : url

  return (
    <div className={css.overlay} onClick={e => e.target === e.currentTarget && onClose()}>
      <div className={css.modal}>
        <div className={css.modalHeader}>
          <span className={css.modalTitle}>{name}</span>
          <button className={css.modalClose} onClick={onClose}>×</button>
        </div>
        <div className={css.modalBody}>
          {isImageUrl(url) ? (
            <img src={fullUrl} alt={name} className={css.previewImage} />
          ) : (
            <div className={css.previewFallback}>
              <LuFileText size={48} color="#9ca3af" />
              <p>Перегляд недоступний для цього типу файлу</p>
              <a href={fullUrl} target="_blank" rel="noreferrer" className={css.btnDownload}>
                <LuExternalLink size={14} /> Відкрити файл
              </a>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

// ── Student Submission Panel ─────────────────────────────────────────────────

function StudentPanel({
  taskId,
  courseId,
  studentId,
  maxMark,
}: {
  taskId: string
  courseId: string
  studentId: string
  maxMark: number
}) {
  const qc = useQueryClient()
  const [file, setFile] = useState<File | null>(null)
  const [link, setLink] = useState('')
  const [fileError, setFileError] = useState('')
  const fileRef = useRef<HTMLInputElement>(null)

  const { data: submissions = [] } = useQuery({
    queryKey: ['my-submissions', courseId],
    queryFn: () => getMySubmissionsForCourse(courseId),
  })

  const submission = submissions.find(s => s.task_id === taskId)

  const invalidate = () => {
    qc.invalidateQueries({ queryKey: ['my-submissions', courseId] })
    qc.invalidateQueries({ queryKey: ['course-tasks', courseId] })
  }

  const submitMutation = useMutation({
    mutationFn: async () => {
      let fileUrl: string | undefined
      if (file) fileUrl = await uploadFile(file)
      return submitTask(taskId, studentId, fileUrl, link.trim() || undefined)
    },
    onSuccess: invalidate,
  })

  const updateMutation = useMutation({
    mutationFn: async () => {
      if (!submission) return
      let fileUrl = submission.file_url
      if (file) fileUrl = await uploadFile(file)
      return updateSubmission(submission.id, {
        file_url: fileUrl,
        link_url: link.trim() || submission.link_url,
      })
    },
    onSuccess: invalidate,
  })

  const unsubmitMutation = useMutation({
    mutationFn: () => unsubmitTask(submission!.id),
    onSuccess: invalidate,
  })

  function handleFile(e: React.ChangeEvent<HTMLInputElement>) {
    setFileError('')
    const f = e.target.files?.[0]
    if (!f) return
    if (f.size > MAX_FILE_SIZE) {
      setFileError('Файл перевищує 10 МБ')
      return
    }
    setFile(f)
  }

  const isGraded = submission?.mark !== undefined && submission.mark !== null
  const isChecking = !!submission && !isGraded
  const isPending = !submission

  return (
    <div className={css.panel}>
      <div className={css.panelTitle}>Ваша здача</div>

      <div className={css.statusRow}>
        <span className={css.statusLabel}>Статус:</span>
        {isPending  && <span className={css.badgePending}>Не здано</span>}
        {isChecking && <span className={css.badgeChecking}>Перевіряється</span>}
        {isGraded   && <span className={css.badgeGraded}>Оцінено</span>}
      </div>

      {isGraded && (
        <>
          <div className={css.score}>{submission!.mark} / {maxMark}</div>
          <div className={css.scoreLabel}>балів</div>
        </>
      )}

      {/* Show existing materials */}
      {submission?.file_url && (
        <div className={css.formGroup}>
          <span className={css.formLabel}>Прикріплений файл</span>
          <a
            href={submission.file_url.startsWith('/uploads/') ? `/api/course${submission.file_url}` : submission.file_url}
            target="_blank"
            rel="noreferrer"
            className={css.materialLink}
          >
            <LuPaperclip size={13} /> {submission.file_url.split('/').pop()}
          </a>
        </div>
      )}
      {submission?.link_url && (
        <div className={css.formGroup}>
          <span className={css.formLabel}>Посилання</span>
          <a href={submission.link_url} target="_blank" rel="noreferrer" className={css.materialLink}>
            <LuLink size={13} /> {submission.link_url}
          </a>
        </div>
      )}

      {/* Submission form */}
      {!isGraded && (
        <>
          {isChecking && <div className={css.panelTitle} style={{ marginTop: 12 }}>Оновити матеріали</div>}

          <div className={css.formGroup}>
            <label className={css.formLabel}>Файл</label>
            <input
              ref={fileRef}
              type="file"
              className={css.fileInput}
              onChange={handleFile}
              accept=".pdf,.doc,.docx,.txt,.png,.jpg,.jpeg,.gif,.zip,.pptx,.xlsx"
            />
            <button
              className={`${css.fileBtn} ${file ? css.fileSelected : ''}`}
              onClick={() => fileRef.current?.click()}
              type="button"
            >
              <LuPaperclip size={14} />
              {file ? file.name : 'Вибрати файл'}
            </button>
            {fileError && <div style={{ color: '#ef4444', fontSize: 12, marginTop: 4 }}>{fileError}</div>}
            <div className={css.fileSizeHint}>Макс. 10 МБ: PDF, DOC, PNG, JPG, ZIP…</div>
          </div>

          <div className={css.divider}>або</div>

          <div className={css.formGroup}>
            <label className={css.formLabel}>Посилання</label>
            <input
              className={css.linkInput}
              type="url"
              placeholder="https://..."
              value={link}
              onChange={e => setLink(e.target.value)}
            />
          </div>

          {isPending ? (
            <button
              className={css.btnPrimary}
              disabled={submitMutation.isPending || (!file && !link.trim())}
              onClick={() => submitMutation.mutate()}
            >
              {submitMutation.isPending ? 'Здаємо…' : 'Здати роботу'}
            </button>
          ) : (
            <button
              className={css.btnPrimary}
              disabled={updateMutation.isPending || (!file && !link.trim())}
              onClick={() => updateMutation.mutate()}
            >
              {updateMutation.isPending ? 'Зберігаємо…' : 'Оновити матеріали'}
            </button>
          )}

          {isChecking && (
            <button
              className={css.btnSecondary}
              disabled={unsubmitMutation.isPending}
              onClick={() => unsubmitMutation.mutate()}
            >
              {unsubmitMutation.isPending ? '…' : 'Скасувати здачу'}
            </button>
          )}
        </>
      )}
    </div>
  )
}

// ── Teacher Submission Card ──────────────────────────────────────────────────

function SubmissionCard({
  submission,
  maxMark,
  onFilePreview,
}: {
  submission: TaskStudent
  maxMark: number
  onFilePreview: (url: string, name: string) => void
}) {
  const qc = useQueryClient()
  const [markInput, setMarkInput] = useState(String(submission.mark ?? ''))

  const gradeMutation = useMutation({
    mutationFn: (mark: number) => gradeSubmission(submission.id, mark),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['task-submissions', submission.task_id] }),
  })

  const isGraded  = submission.mark !== undefined && submission.mark !== null
  const submitted = !!submission.submission_time

  return (
    <div className={css.submissionCard}>
      <div className={css.submissionHeader}>
        <span className={css.studentName}>{submission.student_name ?? 'Студент'}</span>
        <span className={`${css.submissionStatus} ${
          isGraded  ? css.statusGraded :
          submitted ? css.statusSubmitted :
                      css.statusNotSubmitted
        }`}>
          {isGraded ? `Оцінено: ${submission.mark}/${maxMark}` : submitted ? 'Перевіряється' : 'Не здано'}
        </span>
      </div>

      {/* Materials */}
      {(submission.file_url || submission.link_url) && (
        <div className={css.submissionMaterials}>
          {submission.file_url && (
            <button
              className={css.materialLink}
              style={{ background: 'none', border: 'none', cursor: 'pointer', padding: 0 }}
              onClick={() => onFilePreview(
                submission.file_url!,
                submission.file_url!.split('/').pop() ?? 'file',
              )}
            >
              <LuPaperclip size={13} /> {submission.file_url.split('/').pop()}
            </button>
          )}
          {submission.link_url && (
            <a href={submission.link_url} target="_blank" rel="noreferrer" className={css.materialLink}>
              <LuLink size={13} /> Посилання
            </a>
          )}
        </div>
      )}

      {/* Grade row */}
      {submitted && (
        <div className={css.gradeRow}>
          <input
            className={css.gradeInput}
            type="number"
            min="0"
            max={maxMark}
            step="0.5"
            value={markInput}
            onChange={e => setMarkInput(e.target.value)}
            placeholder="0"
          />
          <span className={css.gradeLabel}>/ {maxMark}</span>
          <button
            className={css.btnGrade}
            disabled={gradeMutation.isPending || markInput === ''}
            onClick={() => gradeMutation.mutate(parseFloat(markInput))}
          >
            {gradeMutation.isPending ? '…' : isGraded ? 'Змінити' : 'Оцінити'}
          </button>
        </div>
      )}
    </div>
  )
}

// ── Main ─────────────────────────────────────────────────────────────────────

export default function TaskPage() {
  const { courseId, taskId } = useParams<{ courseId: string; taskId: string }>()
  const navigate = useNavigate()
  const location = useLocation()
  const user = getAuthUser()
  const isTeacher = user?.role === 'teacher'

  const course = (location.state as CourseInfo | null)

  const [fileModal, setFileModal] = useState<{ url: string; name: string } | null>(null)

  const { data: task, isLoading } = useQuery({
    queryKey: ['task', taskId],
    queryFn: () => getTask(taskId!),
    enabled: !!taskId,
  })

  const { data: submissions = [] } = useQuery({
    queryKey: ['task-submissions', taskId],
    queryFn: () => getTaskSubmissions(taskId!),
    enabled: !!taskId && isTeacher,
  })

  if (isLoading) return <div className={css.page} style={{ padding: 28, color: '#6b7280' }}>Завантаження…</div>
  if (!task) return <div className={css.page} style={{ padding: 28, color: '#ef4444' }}>Завдання не знайдено</div>

  return (
    <div className={css.page}>
      {/* Header */}
      <div className={css.header}>
        <button className={css.backBtn} onClick={() => navigate(-1)}>
          <LuArrowLeft size={16} /> Назад
        </button>
        <div className={css.headerInfo}>
          <div className={css.taskTitle}>{task.title}</div>
          <div className={css.taskMeta}>
            {course && <span className={css.metaItem}>{course.name}</span>}
            {task.post_type === 'assignment' && (
              <>
                <span className={css.metaItem}>Дедлайн: {fmtDeadline(task.deadline)}</span>
                <span className={css.metaItem}>{task.max_mark} балів</span>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Body */}
      <div className={css.body}>
        {/* Left: description */}
        <div>
          {task.description && (
            <div className={css.descCard}>
              <div className={css.descTitle}>Опис завдання</div>
              <div className={css.descText}>{task.description}</div>
            </div>
          )}

          {/* Teacher: submissions list */}
          {isTeacher && (
            <div style={{ marginTop: 16 }}>
              <div className={css.descTitle} style={{ marginBottom: 12 }}>
                Здачі студентів ({submissions.length})
              </div>
              {submissions.length === 0 ? (
                <div className={css.descCard}>
                  <span style={{ color: '#9ca3af', fontSize: 14 }}>Ніхто ще не здав це завдання</span>
                </div>
              ) : (
                <div className={css.submissionsList}>
                  {submissions.map(s => (
                    <SubmissionCard
                      key={s.id}
                      submission={s}
                      maxMark={task.max_mark}
                      onFilePreview={(url, name) => setFileModal({ url, name })}
                    />
                  ))}
                </div>
              )}
            </div>
          )}
        </div>

        {/* Right: student panel */}
        {!isTeacher && courseId && user && (
          <div className={css.sidebar}>
            <StudentPanel
              taskId={taskId!}
              courseId={courseId}
              studentId={user.id}
              maxMark={task.max_mark}
            />
          </div>
        )}
      </div>

      {/* File modal */}
      {fileModal && (
        <FileModal
          url={fileModal.url}
          name={fileModal.name}
          onClose={() => setFileModal(null)}
        />
      )}
    </div>
  )
}
