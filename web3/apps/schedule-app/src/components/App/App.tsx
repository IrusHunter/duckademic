import { useEffect, useMemo, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  getLessonOccurrences,
  getLessonSlots,
  getTeachers,
  getClassrooms,
  getStudentGroups,
} from '../../services/scheduleService'
import { nsToTime } from '../../utils/time'
import css from './App.module.css'

const DAY_NAMES = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
const VIEW_OPTIONS = ['General', 'Date'] as const
type ViewOption = (typeof VIEW_OPTIONS)[number]

const ROW_COLORS = [css.rowBlue, css.rowGreen, css.rowRose, css.rowAmber]

function cx(...parts: Array<string | false | undefined | null>) {
  return parts.filter(Boolean).join(' ')
}

function todayIndex(): number {
  const d = new Date().getDay() // 0=Sun, 1=Mon … 6=Sat
  return d === 0 ? 0 : d - 1   // Mon→0 … Sat→5; Sun defaults to Mon
}

export default function App() {
  const [dayIndex, setDayIndex] = useState(todayIndex)
  const [view, setView] = useState<ViewOption>('General')
  const [viewOpen, setViewOpen] = useState(false)

  const occurrences = useQuery({ queryKey: ['lesson-occurrences'], queryFn: getLessonOccurrences })
  const slots = useQuery({ queryKey: ['lesson-slots'], queryFn: getLessonSlots })
  const teachers = useQuery({ queryKey: ['schedule-teachers'], queryFn: getTeachers })
  const classrooms = useQuery({ queryKey: ['schedule-classrooms'], queryFn: getClassrooms })
  const groups = useQuery({ queryKey: ['schedule-groups'], queryFn: getStudentGroups })

  const slotMap = useMemo(
    () => Object.fromEntries((slots.data ?? []).map((s) => [s.id, s])),
    [slots.data],
  )
  const teacherMap = useMemo(
    () => Object.fromEntries((teachers.data ?? []).map((t) => [t.id, t])),
    [teachers.data],
  )
  const classroomMap = useMemo(
    () => Object.fromEntries((classrooms.data ?? []).map((c) => [c.id, c])),
    [classrooms.data],
  )
  const groupMap = useMemo(
    () => Object.fromEntries((groups.data ?? []).map((g) => [g.id, g])),
    [groups.data],
  )

  // API weekday is 1-based (1=Mon … 6=Sat)
  const weekday = dayIndex + 1

  const dayLessons = useMemo(() => {
    return (occurrences.data ?? [])
      .filter((o) => slotMap[o.lesson_slot_id]?.weekday === weekday)
      .sort((a, b) => (slotMap[a.lesson_slot_id]?.slot ?? 0) - (slotMap[b.lesson_slot_id]?.slot ?? 0))
  }, [occurrences.data, slotMap, weekday])

  useEffect(() => {
    function onDocClick() { setViewOpen(false) }
    function onKeyDown(e: KeyboardEvent) { if (e.key === 'Escape') setViewOpen(false) }
    document.addEventListener('click', onDocClick)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('click', onDocClick)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [])

  const isLoading = occurrences.isLoading || slots.isLoading
  const isError = occurrences.isError || slots.isError

  return (
    <main className={css.page}>
      <section className={css.layout}>
        <h1 className={css.title}>Schedule</h1>

        <div className={css.card}>
          <header className={css.mainHeader}>
            <div className={css.mainLeft}>
              <svg width="20" height="20">
                <use href="/img/icons.svg#icon-SVG-11" />
              </svg>
              <span className={css.heading}>Today&apos;s Classes</span>
            </div>

            <div className={css.controls}>
              <div className={css.dayNavRow}>
                <button
                  className={css.navBtn}
                  type="button"
                  onClick={() => setDayIndex((i) => (i - 1 + DAY_NAMES.length) % DAY_NAMES.length)}
                >
                  <svg width="14" height="20">
                    <use href="/img/icons.svg#icon-vector-left" />
                  </svg>
                </button>

                <button className={css.dayBtn} type="button">
                  <span className={css.dayLabel}>{DAY_NAMES[dayIndex]}</span>
                </button>

                <button
                  className={css.navBtn}
                  type="button"
                  onClick={() => setDayIndex((i) => (i + 1) % DAY_NAMES.length)}
                >
                  <svg width="14" height="20">
                    <use href="/img/icons.svg#icon-vector-right" />
                  </svg>
                </button>
              </div>

              <div className={cx(css.viewDd, viewOpen && css.open)}>
                <button
                  className={css.viewSelect}
                  type="button"
                  aria-expanded={viewOpen}
                  onClick={(e) => {
                    e.stopPropagation()
                    setViewOpen((v) => !v)
                  }}
                >
                  <span className={css.viewLabel}>{view}</span>
                  <span className={css.viewCaret}>
                    <svg width="18" height="13">
                      <use href="/img/icons.svg#icon-vector-down" />
                    </svg>
                  </span>
                </button>

                <div className={css.dropdown}>
                  {VIEW_OPTIONS.map((opt) => (
                    <button
                      key={opt}
                      className={cx(css.dropdownItem, opt === view && css.dropdownItemActive)}
                      type="button"
                      onClick={() => {
                        setView(opt)
                        setViewOpen(false)
                      }}
                    >
                      {opt}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </header>

          <ul className={css.classes}>
            {isLoading && (
              <li className={css.classRow}>
                <div className={css.classTime} />
                <article className={css.classCard}>
                  <p className={css.classTeacher}>Loading…</p>
                </article>
              </li>
            )}

            {isError && (
              <li className={css.classRow}>
                <div className={css.classTime} />
                <article className={css.classCard}>
                  <p className={css.classTeacher}>Failed to load schedule.</p>
                </article>
              </li>
            )}

            {!isLoading && !isError && dayLessons.length === 0 && (
              <li className={css.classRow}>
                <div className={css.classTime} />
                <article className={css.classCard}>
                  <p className={css.classTeacher}>No classes on this day.</p>
                </article>
              </li>
            )}

            {dayLessons.map((occ, idx) => {
              const slot = slotMap[occ.lesson_slot_id]
              const teacher = teacherMap[occ.teacher_id]
              const classroom = classroomMap[occ.classroom_id]
              const group = groupMap[occ.student_group_id]

              const startTime = slot ? nsToTime(slot.start_time) : '—'
              const endTime = slot ? nsToTime(slot.start_time + slot.duration) : '—'
              const colorClass = ROW_COLORS[idx % ROW_COLORS.length]

              return (
                <li key={occ.id} className={cx(css.classRow, colorClass)}>
                  <div className={css.classTime}>{startTime} – {endTime}</div>

                  <article className={css.classCard}>
                    <header className={css.classCardHeader}>
                      <h2 className={css.classTitle}>
                        {teacher?.name ?? `Slot ${slot?.slot ?? '—'}`}
                      </h2>
                      {slot && <span className={css.classTag}>#{slot.slot}</span>}
                    </header>

                    {group && <p className={css.classTeacher}>{group.name}</p>}

                    <div className={css.classMeta}>
                      {classroom && (
                        <div className={css.classMetaItem}>
                          Room {classroom.number}
                        </div>
                      )}
                      {view === 'Date' && occ.date && (
                        <div className={css.classMetaDates}>{occ.date}</div>
                      )}
                    </div>
                  </article>
                </li>
              )
            })}
          </ul>
        </div>
      </section>
    </main>
  )
}
