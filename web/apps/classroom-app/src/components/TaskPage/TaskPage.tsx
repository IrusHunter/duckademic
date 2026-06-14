import { useState, useRef } from 'react'
import { useParams, useNavigate, useLocation } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  LuArrowLeft, LuPaperclip, LuLink, LuFileText,
  LuExternalLink, LuClipboardList, LuMegaphone,
  LuCalendar, LuStar, LuClock, LuChevronLeft, LuSend,
} from 'react-icons/lu'
import {
  getTask,
  getTaskSubmissions,
  getMySubmissionsForCourse,
  submitTask,
  unsubmitTask,
  gradeSubmission,
  uploadFile,
  getTaskComments,
  addTaskComment,
  getPrivateComments,
  addPrivateComment,
} from '../../services/courseService'
import { getAuthUser } from '../../utils/user'
import type { CourseInfo, TaskStudent } from '../../types/course'
import css from './TaskPage.module.css'

// ── helpers ───────────────────────────────────────────────────────────────────

function fmtDate(iso: string) {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleString('en-US', { month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function fmtDateTime(iso: string) {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function timeAgo(iso: string) {
  const diff = Date.now() - new Date(iso).getTime()
  const h = Math.floor(diff / 3_600_000)
  if (h < 1) return 'Just now'
  if (h < 24) return `${h}h ago`
  const d = Math.floor(h / 24)
  return `${d} day${d > 1 ? 's' : ''} ago`
}

function isImageUrl(url: string) { return /\.(png|jpe?g|gif|webp|svg)$/i.test(url) }
function isOverdue(iso: string)  { return new Date(iso) < new Date() }

function nameAbbr(name: string): string {
  const parts = name.trim().split(/\s+/)
  if (parts.length < 2) return name
  return parts[0] + ' ' + parts.slice(1).map(p => p.charAt(0).toUpperCase() + '.').join('')
}

const MAX_FILE_SIZE = 10 * 1024 * 1024

// ── File preview modal ────────────────────────────────────────────────────────

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
              <p>Preview not available for this file type</p>
              <a href={fullUrl} target="_blank" rel="noreferrer" className={css.btnDownload}>
                <LuExternalLink size={14} /> Open file
              </a>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

// ── Comments section ──────────────────────────────────────────────────────────

function Comments({
  placeholder,
  taskId,
  isPrivate = false,
  studentId,
}: {
  placeholder: string
  taskId: string
  isPrivate?: boolean
  studentId?: string
}) {
  const qc = useQueryClient()
  const [input, setInput] = useState('')

  const queryKey = isPrivate
    ? ['comments-private', taskId, studentId]
    : ['comments-class', taskId]

  const { data: comments = [] } = useQuery({
    queryKey,
    queryFn: () =>
      isPrivate ? getPrivateComments(taskId, studentId) : getTaskComments(taskId),
    enabled: !!taskId,
  })

  const addMutation = useMutation({
    mutationFn: (body: string) =>
      isPrivate
        ? addPrivateComment(taskId, body, studentId)
        : addTaskComment(taskId, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey })
      setInput('')
    },
  })

  function send() {
    const trimmed = input.trim()
    if (!trimmed || addMutation.isPending) return
    addMutation.mutate(trimmed)
  }

  return (
    <div>
      {comments.map(c => (
        <div key={c.id} className={css.commentItem}>
          <span className={css.commentAuthor}>{c.author_name}</span>
          <span className={css.commentBody}>{c.body}</span>
        </div>
      ))}
      <div className={css.commentRow}>
        <input
          className={css.commentInput}
          placeholder={placeholder}
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={e => e.key === 'Enter' && send()}
        />
        <button
          className={css.commentSendBtn}
          disabled={!input.trim() || addMutation.isPending}
          onClick={send}
        >
          <LuSend size={13} />
        </button>
      </div>
    </div>
  )
}

// ── Student: My Work card ─────────────────────────────────────────────────────

function MyWorkCard({
  taskId, courseId, studentId, maxMark, deadline, onFilePreview,
}: {
  taskId: string; courseId: string; studentId: string
  maxMark: number; deadline: string
  onFilePreview: (url: string, name: string) => void
}) {
  const qc = useQueryClient()
  const [file, setFile]           = useState<File | null>(null)
  const [link, setLink]           = useState('')
  const [showLink, setShowLink]   = useState(false)
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
    onSuccess: () => { invalidate(); setFile(null); setLink(''); setShowLink(false) },
  })

  const unsubmitMutation = useMutation({
    mutationFn: () => unsubmitTask(submission!.id),
    onSuccess: invalidate,
  })

  function handleFile(e: React.ChangeEvent<HTMLInputElement>) {
    setFileError('')
    const f = e.target.files?.[0]
    if (!f) return
    if (f.size > MAX_FILE_SIZE) { setFileError('File exceeds 10 MB'); return }
    setFile(f)
  }

  const isGraded   = submission?.mark !== undefined && submission.mark !== null
  const isChecking = !!submission && !isGraded
  const isPending  = !submission
  const overdue    = isOverdue(deadline)
  const isLate     = !!submission?.submission_time && new Date(submission.submission_time) > new Date(deadline)
  const locked     = isChecking || isGraded

  return (
    <div className={css.myWorkCard}>
      {/* Header */}
      <div className={css.mwHeader}>
        <span className={css.mwTitle}>My work</span>
        {isPending && overdue  && <span className={`${css.badge} ${css.badgeMissed}`}>Missing</span>}
        {isPending && !overdue && <span className={`${css.badge} ${css.badgePending}`}>Not submitted</span>}
        {isChecking            && <span className={`${css.badge} ${css.badgeChecking}`}>Submitted</span>}
        {isGraded              && <span className={`${css.badge} ${css.badgeGraded}`}>Graded</span>}
      </div>

      {/* Grade */}
      {isGraded && (
        <div className={css.gradeDisplay}>
          <span className={css.gradeScore}>{submission!.mark}</span>
          <span className={css.gradeOf}>/{maxMark}</span>
        </div>
      )}

      {/* Submitted attachments */}
      {submission?.file_url && (
        <button
          className={css.attachedItem}
          onClick={() => onFilePreview(submission.file_url!, submission.file_url!.split('/').pop() ?? 'file')}
        >
          <LuPaperclip size={13} />
          <span className={css.attachedItemName}>{submission.file_url.split('/').pop()}</span>
          <LuExternalLink size={11} style={{ marginLeft: 'auto', opacity: 0.4 }} />
        </button>
      )}
      {submission?.link_url && (
        <a href={submission.link_url} target="_blank" rel="noreferrer" className={css.attachedItem}>
          <LuLink size={13} />
          <span className={css.attachedItemUrl}>{submission.link_url}</span>
          <LuExternalLink size={11} style={{ marginLeft: 'auto', opacity: 0.4 }} />
        </a>
      )}

      {/* Pending file chip */}
      {isPending && file && (
        <div className={css.fileChip}>
          <LuPaperclip size={12} />
          <span>{file.name}</span>
          <button className={css.fileChipRemove} onClick={() => setFile(null)}>×</button>
        </div>
      )}

      {/* Hidden file input */}
      <input
        ref={fileRef} type="file" style={{ display: 'none' }} onChange={handleFile}
        accept=".pdf,.doc,.docx,.txt,.png,.jpg,.jpeg,.gif,.zip,.pptx,.xlsx"
      />

      {/* Attach buttons */}
      {!locked && (
        <div className={css.attachBtns}>
          <button className={css.attachBtn} onClick={() => fileRef.current?.click()}>
            <LuPaperclip size={14} />
            {file ? 'Change file' : 'Attach file'}
          </button>
          <button
            className={`${css.attachBtn} ${showLink ? css.attachBtnActive : ''}`}
            onClick={() => setShowLink(v => !v)}
          >
            <LuLink size={14} /> Add link
          </button>
        </div>
      )}

      {/* Link input (toggled) */}
      {!locked && showLink && (
        <input
          className={css.linkInput}
          type="url"
          placeholder="Paste link…"
          value={link}
          onChange={e => setLink(e.target.value)}
          autoFocus
        />
      )}

      {fileError && <p className={css.fileError}>{fileError}</p>}

      {/* Action button */}
      {!isGraded && (
        <div className={css.submitRow}>
          {isPending ? (
            <button
              className={css.btnPrimary}
              disabled={submitMutation.isPending || (!file && !link.trim())}
              onClick={() => submitMutation.mutate()}
            >
              {submitMutation.isPending ? 'Submitting…' : overdue ? 'Turn in late' : 'Turn in'}
            </button>
          ) : (
            <button
              className={css.btnSecondary}
              disabled={unsubmitMutation.isPending}
              onClick={() => unsubmitMutation.mutate()}
            >
              {unsubmitMutation.isPending ? '…' : 'Cancel submission'}
            </button>
          )}
        </div>
      )}

      {/* Status text */}
      <div className={css.statusLine}>
        {isPending && !overdue && (
          <span className={css.statusPending}>Assignment pending · Due {fmtDate(deadline)}</span>
        )}
        {isPending && overdue && (
          <span className={css.statusMissed}>Missed deadline · Was due {fmtDate(deadline)}</span>
        )}
        {(isChecking || isGraded) && !isLate && (
          <span className={css.statusOk}>Submitted · {fmtDate(submission!.submission_time!)}</span>
        )}
        {(isChecking || isGraded) && isLate && (
          <span className={css.statusLate}>Submitted late · {fmtDate(submission!.submission_time!)}</span>
        )}
      </div>

      {/* Private comments */}
      <div className={css.privateSep}>
        <div className={css.commentsTitle}>Private comments</div>
        <Comments
          placeholder="Add private comment…"
          taskId={taskId}
          isPrivate
          studentId={studentId}
        />
      </div>
    </div>
  )
}

// ── Teacher: Submission detail ────────────────────────────────────────────────

function SubmissionDetail({
  submission, maxMark, deadline, onBack, onFilePreview,
}: {
  submission: TaskStudent; maxMark: number; deadline: string
  onBack: () => void
  onFilePreview: (url: string, name: string) => void
}) {
  const qc = useQueryClient()
  const [markInput, setMarkInput]       = useState(String(submission.mark ?? ''))
  const [returnComment, setReturnComment] = useState('')

  const gradeMutation = useMutation({
    mutationFn: (mark: number) => gradeSubmission(submission.id, mark),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['task-submissions', submission.task_id] }),
  })

  const returnMutation = useMutation({
    mutationFn: async () => {
      if (returnComment.trim()) {
        await addPrivateComment(submission.task_id, returnComment.trim(), submission.student_id)
      }
      await unsubmitTask(submission.id)
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['task-submissions', submission.task_id] })
      qc.invalidateQueries({ queryKey: ['comments-private', submission.task_id, submission.student_id] })
      onBack()
    },
  })

  const isGraded    = submission.mark !== undefined && submission.mark !== null
  const submittedAt = submission.submission_time ? new Date(submission.submission_time) : null
  const onTime      = submittedAt ? submittedAt <= new Date(deadline) : null

  return (
    <div className={css.myWorkCard}>
      <button className={css.backToList} onClick={onBack}>
        <LuChevronLeft size={14} /> All submissions
      </button>

      <div className={css.detailHeader}>
        <span className={css.studentName}>{submission.student_name ?? 'Student'}</span>
        <span className={`${css.badge} ${isGraded ? css.badgeGraded : submission.submission_time ? css.badgeChecking : css.badgeMissed}`}>
          {isGraded ? `${submission.mark}/${maxMark}` : submission.submission_time ? 'Submitted' : 'Missing'}
        </span>
      </div>

      {submittedAt && (
        <div className={`${css.timeChip} ${onTime ? css.timeOnTime : css.timeLate}`}>
          <LuClock size={12} />
          {onTime ? 'On time' : 'Late'} · {fmtDateTime(submission.submission_time!)}
        </div>
      )}

      {(submission.file_url || submission.link_url) && (
        <div className={css.detailAttachments}>
          <div className={css.sectionLabel}>Attachments</div>
          {submission.file_url && (
            <button
              className={css.attachedItem}
              onClick={() => onFilePreview(submission.file_url!, submission.file_url!.split('/').pop() ?? 'file')}
            >
              <LuPaperclip size={13} />
              <span className={css.attachedItemName}>{submission.file_url.split('/').pop()}</span>
              <LuExternalLink size={11} style={{ marginLeft: 'auto', opacity: 0.4 }} />
            </button>
          )}
          {submission.link_url && (
            <a href={submission.link_url} target="_blank" rel="noreferrer" className={css.attachedItem}>
              <LuLink size={13} />
              <span className={css.attachedItemUrl}>{submission.link_url}</span>
              <LuExternalLink size={11} style={{ marginLeft: 'auto', opacity: 0.4 }} />
            </a>
          )}
        </div>
      )}

      {submission.submission_time && (
        <>
          <div className={css.commentsTitle} style={{ marginTop: 14 }}>Grade</div>
          <div className={css.gradeRow}>
            <input
              className={css.gradeInput}
              type="number" min="0" max={maxMark} step="0.5"
              value={markInput} onChange={e => setMarkInput(e.target.value)}
              placeholder="0"
            />
            <span className={css.gradeSlash}>/ {maxMark}</span>
            <button
              className={css.btnPrimary}
              style={{ flex: 1, padding: '7px 12px', fontSize: 13 }}
              disabled={gradeMutation.isPending || markInput === ''}
              onClick={() => gradeMutation.mutate(parseFloat(markInput))}
            >
              {gradeMutation.isPending ? '…' : isGraded ? 'Update' : 'Grade'}
            </button>
          </div>

          <div className={css.returnSection}>
            <div className={css.commentsTitle}>Return with comment</div>
            <textarea
              className={css.returnTextarea}
              placeholder="Write a comment for the student…"
              value={returnComment}
              onChange={e => setReturnComment(e.target.value)}
              rows={3}
            />
            <button
              className={css.btnSecondary}
              disabled={returnMutation.isPending}
              onClick={() => returnMutation.mutate()}
            >
              {returnMutation.isPending ? '…' : 'Return assignment'}
            </button>
          </div>
        </>
      )}

      <div className={css.privateSep}>
        <div className={css.commentsTitle}>Private comments</div>
        <Comments
          placeholder="Write a private comment…"
          taskId={submission.task_id}
          isPrivate
          studentId={submission.student_id}
        />
      </div>
    </div>
  )
}

// ── Main ──────────────────────────────────────────────────────────────────────

export default function TaskPage() {
  const { courseId, taskId } = useParams<{ courseId: string; taskId: string }>()
  const navigate  = useNavigate()
  const location  = useLocation()
  const user      = getAuthUser()
  const isTeacher = user?.role === 'teacher'
  const course    = location.state as CourseInfo | null

  const [fileModal, setFileModal] = useState<{ url: string; name: string } | null>(null)
  const [selectedSub, setSelectedSub] = useState<TaskStudent | null>(null)

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

  const openPreview = (url: string, name: string) => setFileModal({ url, name })

  if (isLoading) return <div className={css.page}><p className={css.stateMsg}>Loading…</p></div>
  if (!task)     return <div className={css.page}><p className={css.stateMsg} style={{ color: '#ef4444' }}>Task not found</p></div>

  const isAnnouncement = task.post_type === 'announcement'
  const turnedIn = submissions.filter(s => !!s.submission_time).length
  const graded   = submissions.filter(s => s.mark !== undefined && s.mark !== null).length
  const hasRight = !isAnnouncement

  return (
    <div className={css.page}>
      <div className={css.inner}>
        {/* Back */}
        <button className={css.backBtn} onClick={() => navigate(-1)}>
          <LuArrowLeft size={14} />
          {course?.name ?? 'Back'}
        </button>

        {/* Main card */}
        <div className={`${css.mainCard} ${hasRight ? css.mainCardWithSidebar : ''}`}>
          {/* ── Left column ── */}
          <div className={css.mainLeft}>
            {/* Task title */}
            <div className={css.taskTitleRow}>
              {isAnnouncement
                ? <LuMegaphone size={18} className={css.taskIcon} />
                : <LuClipboardList size={18} className={css.taskIcon} />
              }
              <h1 className={css.taskTitle}>{task.title}</h1>
            </div>

            {/* Meta */}
            <div className={css.taskMeta}>
              {course?.teacher_name && <span>{course.teacher_name}</span>}
              <span>{timeAgo(task.created_at)}</span>
              {!isAnnouncement && (
                <>
                  <span className={css.metaChip}><LuStar size={12} /> {task.max_mark} pts</span>
                  <span className={css.metaChip}><LuCalendar size={12} /> Due {fmtDate(task.deadline)}</span>
                </>
              )}
            </div>

            {/* Teacher stats */}
            {isTeacher && !isAnnouncement && (
              <div className={css.statsRow}>
                <div className={css.statItem}><b>{turnedIn}</b> submitted</div>
                <div className={css.statItem}><b>{graded}</b> graded</div>
                <div className={css.statItem}><b>{submissions.length}</b> total</div>
              </div>
            )}

            <hr className={css.divider} />

            {/* Description */}
            {task.description ? (
              <div className={css.descText}>{task.description}</div>
            ) : (
              <p className={css.emptyDesc}>No description provided.</p>
            )}

            {/* Teacher: submissions list */}
            {isTeacher && !isAnnouncement && (
              <>
                <hr className={css.divider} />
                <div className={css.subListHeader}>
                  Student submissions
                  {submissions.length > 0 && <span className={css.countPill}>{submissions.length}</span>}
                </div>
                {submissions.length === 0 ? (
                  <p className={css.emptyDesc}>No submissions yet.</p>
                ) : (
                  <div className={css.subList}>
                    {submissions.map(s => {
                      const g = s.mark !== undefined && s.mark !== null
                      const sub = !!s.submission_time
                      return (
                        <button
                          key={s.id}
                          className={`${css.subRow} ${selectedSub?.id === s.id ? css.subRowActive : ''}`}
                          onClick={() => setSelectedSub(s)}
                        >
                          <div className={css.subInfo}>
                            <span className={css.subName}>{nameAbbr(s.student_name ?? 'Student')}</span>
                            {sub && <span className={css.subMeta}>{fmtDateTime(s.submission_time!)}</span>}
                            <span className={css.subMeta}>
                              {g ? `${s.mark} / ${task.max_mark} pts` : sub ? '—' : 'Not submitted'}
                            </span>
                          </div>
                          {(s.file_url || s.link_url) && <LuPaperclip size={12} style={{ color: '#9ca3af', flexShrink: 0 }} />}
                        </button>
                      )
                    })}
                  </div>
                )}
              </>
            )}

            <hr className={css.divider} />

            {/* Class comments */}
            <div className={css.commentsSection}>
              <div className={css.commentsTitle}>Class comments</div>
              <Comments placeholder="Add a class comment…" taskId={taskId!} />
            </div>
          </div>

          {/* ── Right column ── */}
          {hasRight && (
            <div className={css.mainRight}>
              {!isTeacher && courseId && user && (
                <MyWorkCard
                  taskId={taskId!}
                  courseId={courseId}
                  studentId={user.id}
                  maxMark={task.max_mark}
                  deadline={task.deadline}
                  onFilePreview={openPreview}
                />
              )}
              {isTeacher && (
                selectedSub ? (
                  <SubmissionDetail
                    submission={selectedSub}
                    maxMark={task.max_mark}
                    deadline={task.deadline}
                    onBack={() => setSelectedSub(null)}
                    onFilePreview={openPreview}
                  />
                ) : (
                  <div className={css.myWorkCard} style={{ textAlign: 'center', padding: '28px 16px' }}>
                    <p style={{ fontSize: 13, color: '#9ca3af', margin: 0 }}>
                      Select a submission from the list to review it.
                    </p>
                  </div>
                )
              )}
            </div>
          )}
        </div>
      </div>

      {fileModal && (
        <FileModal url={fileModal.url} name={fileModal.name} onClose={() => setFileModal(null)} />
      )}
    </div>
  )
}
